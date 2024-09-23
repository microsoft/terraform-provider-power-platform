// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_groups

import (
	"context"
	"net/http"
	"net/url"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

func newEnvironmentGroupClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}

type client struct {
	Api *api.Client
}

func (client *client) CreateEnvironmentGroup(ctx context.Context, environmentGroup environmentGroupDto) (*environmentGroupDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/environmentGroups",
	}

	values := url.Values{}
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()

	newEnvironmentGroup := environmentGroupDto{}
	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, environmentGroup, []int{http.StatusCreated}, &newEnvironmentGroup)
	if err != nil {
		return nil, err
	}

	return &newEnvironmentGroup, nil
}

// DeleteEnvironmentGroup deletes an environment group.
func (client *client) DeleteEnvironmentGroup(ctx context.Context, environmentGroupId string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/environmentGroups/" + environmentGroupId,
	}

	values := url.Values{}
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()

	_, err := client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return err
	}

	return nil
}

// updateEnvironmentGroup updates an environment group.
func (client *client) UpdateEnvironmentGroup(ctx context.Context, environmentGroupId string, environmentGroup environmentGroupDto) (*environmentGroupDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/environmentGroups/" + environmentGroupId,
	}

	values := url.Values{}
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()

	updatedEnvironmentGroup := environmentGroupDto{}
	_, err := client.Api.Execute(ctx, nil, "PUT", apiUrl.String(), nil, environmentGroup, []int{http.StatusOK}, &updatedEnvironmentGroup)
	if err != nil {
		return nil, err
	}

	return &updatedEnvironmentGroup, nil
}

// GetEnvironmentGroup gets an environment group.
func (client *client) GetEnvironmentGroup(ctx context.Context, environmentGroupId string) (*environmentGroupDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/environmentGroups/" + environmentGroupId,
	}

	values := url.Values{}
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()

	environmentGroup := environmentGroupDto{}
	httpResponse, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNotFound}, &environmentGroup)
	if httpResponse.HttpResponse.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &environmentGroup, nil
}
