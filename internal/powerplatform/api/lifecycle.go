// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/constants"
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

func (client *ApiClient) DoWaitForLifecycleOperationStatus(ctx context.Context, response *ApiHttpResponse) (*LifecycleDto, error) {

	locationHeader := response.GetHeader(constants.HEADER_LOCATION)
	tflog.Debug(ctx, "Location Header: "+locationHeader)

	_, err := url.Parse(locationHeader)
	if err != nil {
		tflog.Error(ctx, "Error parsing location header: "+err.Error())
	}

	retryHeader := response.GetHeader(constants.HEADER_RETRY_AFTER)
	tflog.Debug(ctx, "Retry Header: "+retryHeader)
	retryAfter, err := time.ParseDuration(retryHeader)

	if err != nil {
		retryAfter = time.Duration(5) * time.Second
	} else {
		retryAfter = retryAfter * time.Second
	}

	for {
		lifecycleResponse := LifecycleDto{}
		response, err = client.Execute(ctx, "GET", locationHeader, nil, nil, []int{http.StatusOK}, &lifecycleResponse)
		if err != nil {
			return nil, err
		}

		err = client.SleepWithContext(ctx, retryAfter)
		if err != nil {
			return nil, err
		}

		tflog.Debug(ctx, "Environment Creation Operation State: '"+lifecycleResponse.State.Id+"'")
		tflog.Debug(ctx, "Environment Creation Operation HTTP Status: '"+response.Response.Status+"'")

		if lifecycleResponse.State.Id == "Succeeded" {
			return &lifecycleResponse, nil
		} else if lifecycleResponse.State.Id == "Failed" {
			return &lifecycleResponse, errors.New("environment creation failed. provisioning state: " + lifecycleResponse.State.Id)
		}
	}
}
