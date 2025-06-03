// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package admin_management_application

import (
	"context"

	"fmt"
	"net/http"
	"net/url"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

func newAdminManagementApplicationClient(clientApi *api.Client) client {
	return client{
		Api: clientApi,
	}
}

type client struct {
	Api *api.Client
}

func (client *client) GetAdminApplication(ctx context.Context, clientId string) (*adminManagementApplicationDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/adminApplications/%s", clientId),
		RawQuery: url.Values{
			constants.API_VERSION_PARAM: []string{"2020-10-01"},
		}.Encode(),
	}

	var adminApp adminManagementApplicationDto
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &adminApp)

	// Returning pointer to local variable is acceptable for small DTOs
	return &adminApp, err
}

func (client *client) RegisterAdminApplication(ctx context.Context, clientId string) (*adminManagementApplicationDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/adminApplications/%s", clientId),
		RawQuery: url.Values{
			"api-version": []string{"2020-10-01"},
		}.Encode(),
	}

	var adminApp adminManagementApplicationDto
	_, err := client.Api.Execute(ctx, nil, "PUT", apiUrl.String(), nil, nil, []int{http.StatusOK}, &adminApp)

	// Returning pointer to local variable is acceptable for small DTOs
	return &adminApp, err
}

func (client *client) UnregisterAdminApplication(ctx context.Context, clientId string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/adminApplications/%s", clientId),
		RawQuery: url.Values{
			"api-version": []string{"2020-10-01"},
		}.Encode(),
	}

	_, err := client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNoContent}, nil)

	return err
}
