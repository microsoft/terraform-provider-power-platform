// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package publisher

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment"
)

type client struct {
	Api               *api.Client
	environmentClient environment.Client
}

func newPublisherClient(apiClient *api.Client) client {
	return client{
		Api:               apiClient,
		environmentClient: environment.NewEnvironmentClient(apiClient),
	}
}

func (client *client) CreatePublisher(ctx context.Context, environmentId string, model *ResourceModel) (*publisherDto, error) {
	environmentHost, err := client.environmentClient.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := helpers.BuildDataverseApiUrl(environmentHost, "/api/data/v9.2/publishers", nil)
	resp, err := client.Api.Execute(ctx, nil, http.MethodPost, apiUrl, nil, publisherBodyFromModel(model), []int{http.StatusCreated, http.StatusNoContent, http.StatusForbidden}, nil)
	if err != nil {
		return nil, err
	}
	if err := client.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}

	publisherId, err := getPublisherIdFromResponse(resp)
	if err != nil {
		return nil, err
	}

	return client.GetPublisherById(ctx, environmentId, publisherId)
}

func (client *client) GetPublisherById(ctx context.Context, environmentId, publisherId string) (*publisherDto, error) {
	environmentHost, err := client.environmentClient.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := helpers.BuildDataverseApiUrl(environmentHost, fmt.Sprintf("/api/data/v9.2/publishers(%s)", publisherId), nil)
	publisher := publisherDto{}
	resp, err := client.Api.Execute(ctx, nil, http.MethodGet, apiUrl, nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &publisher)
	if err != nil {
		if resp != nil && resp.HttpResponse.StatusCode == http.StatusNotFound {
			return nil, customerrors.WrapIntoProviderError(err, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("publisher '%s' not found", publisherId))
		}
		return nil, err
	}
	if err := client.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if err := client.Api.HandleNotFoundResponse(resp); err != nil {
		return nil, customerrors.WrapIntoProviderError(err, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("publisher '%s' not found", publisherId))
	}

	return &publisher, nil
}

func (client *client) GetPublisherByUniqueName(ctx context.Context, environmentId, uniqueName string) (*publisherDto, error) {
	environmentHost, err := client.environmentClient.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Add("$filter", fmt.Sprintf("uniquename eq '%s'", escapeODataString(uniqueName)))
	apiUrl := helpers.BuildDataverseApiUrl(environmentHost, "/api/data/v9.2/publishers", values)

	publishers := publishersDto{}
	resp, err := client.Api.Execute(ctx, nil, http.MethodGet, apiUrl, nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &publishers)
	if err != nil {
		return nil, err
	}
	if err := client.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if err := client.Api.HandleNotFoundResponse(resp); err != nil {
		return nil, customerrors.WrapIntoProviderError(err, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("publisher '%s' not found", uniqueName))
	}

	if len(publishers.Value) == 0 {
		return nil, customerrors.WrapIntoProviderError(nil, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("publisher '%s' not found", uniqueName))
	}

	return &publishers.Value[0], nil
}

func (client *client) UpdatePublisher(ctx context.Context, environmentId, publisherId string, model *ResourceModel) (*publisherDto, error) {
	environmentHost, err := client.environmentClient.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := helpers.BuildDataverseApiUrl(environmentHost, fmt.Sprintf("/api/data/v9.2/publishers(%s)", publisherId), nil)
	resp, err := client.Api.Execute(ctx, nil, http.MethodPatch, apiUrl, nil, publisherBodyFromModel(model), []int{http.StatusNoContent, http.StatusForbidden}, nil)
	if err != nil {
		return nil, err
	}
	if err := client.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}

	return client.GetPublisherById(ctx, environmentId, publisherId)
}

func (client *client) DeletePublisher(ctx context.Context, environmentId, publisherId string) error {
	environmentHost, err := client.environmentClient.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return err
	}

	apiUrl := helpers.BuildDataverseApiUrl(environmentHost, fmt.Sprintf("/api/data/v9.2/publishers(%s)", publisherId), nil)
	resp, err := client.Api.Execute(ctx, nil, http.MethodDelete, apiUrl, nil, nil, []int{http.StatusNoContent, http.StatusNotFound, http.StatusForbidden}, nil)
	if err != nil {
		return err
	}
	if err := client.Api.HandleForbiddenResponse(resp); err != nil {
		return err
	}
	if resp.HttpResponse.StatusCode == http.StatusNotFound {
		return customerrors.WrapIntoProviderError(nil, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("publisher '%s' not found", publisherId))
	}

	return nil
}

func escapeODataString(input string) string {
	return strings.ReplaceAll(input, "'", "''")
}
