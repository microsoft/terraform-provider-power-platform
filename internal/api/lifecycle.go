// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

import (
	"context"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

type LifecycleDto struct {
	Id                 string                  `json:"id"`
	Links              LifecycleLinksDto       `json:"links"`
	State              LifecycleStateDto       `json:"state"`
	Type               LifecycleStateDto       `json:"type"`
	CreatedDateTime    string                  `json:"createdDateTime"`
	LastActionDateTime string                  `json:"lastActionDateTime"`
	RequestedBy        LifecycleRequestedByDto `json:"requestedBy"`
	Stages             []LifecycleStageDto     `json:"stages"`
}

type LifecycleStageDto struct {
	Id                  string            `json:"id"`
	Name                string            `json:"name"`
	State               LifecycleStateDto `json:"state"`
	FirstActionDateTime string            `json:"firstActionDateTime"`
	LastActionDateTime  string            `json:"lastActionDateTime"`
}

type LifecycleLinksDto struct {
	Self        LifecycleLinkDto `json:"self"`
	Environment LifecycleLinkDto `json:"environment"`
}

type LifecycleLinkDto struct {
	Path string `json:"path"`
}

type LifecycleStateDto struct {
	Id string `json:"id"`
}

type LifecycleRequestedByDto struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	Type        string `json:"type"`
}

func (client *Client) DoWaitForLifecycleOperationStatus(ctx context.Context, response *Response) (*LifecycleDto, error) {
	locationHeader := response.GetHeader(constants.HEADER_LOCATION)
	if locationHeader == "" {
		locationHeader = response.GetHeader(constants.HEADER_OPERATION_LOCATION)
	}
	tflog.Debug(ctx, "Location Header: "+locationHeader)

	_, err := url.Parse(locationHeader)
	if err != nil {
		tflog.Error(ctx, "Error parsing location header: "+err.Error())
	}

	waitFor := retryAfter(ctx, response.HttpResponse)

	for {
		lifecycleResponse := LifecycleDto{}
		response, err = client.Execute(ctx, nil, "GET", locationHeader, nil, nil, []int{http.StatusOK, http.StatusConflict}, &lifecycleResponse)
		if err != nil {
			return nil, err
		}

		err = client.SleepWithContext(ctx, waitFor)
		if err != nil {
			return nil, err
		}

		tflog.Debug(ctx, "Lifecycle Operation HTTP Status: '"+response.HttpResponse.Status+"'")
		if response.HttpResponse.StatusCode == http.StatusConflict {
			continue
		}
		tflog.Debug(ctx, "Lifecycle Operation State: '"+lifecycleResponse.State.Id+"'")

		if lifecycleResponse.State.Id == "Succeeded" || lifecycleResponse.State.Id == "Failed" {
			return &lifecycleResponse, nil
		}
	}
}
