// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package admin_management_application

import (
	"context"
	
	"fmt"
	"net/http"
	"net/url"
	
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
)

func NewAdminManagementApplicationClient(api *api.ApiClient) AdminManagementApplicationClient {
	return AdminManagementApplicationClient{
		Api: api,
	}
}

type AdminManagementApplicationClient struct {	
	Api *api.ApiClient
}

func (client *AdminManagementApplicationClient) GetAdminApplication(ctx context.Context, clientId string) (*AdminManagementApplicationDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:  fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/adminApplications/%s", clientId),
		RawQuery: url.Values{
			"api-version": []string{"2020-10-01"},
		}.Encode(),
	}

	var adminApp AdminManagementApplicationDto
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &adminApp)

	return &adminApp, err
}

func (client *AdminManagementApplicationClient) RegisterAdminApplication(ctx context.Context, clientId string) (*AdminManagementApplicationDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:  fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/adminApplications/%s", clientId),
		RawQuery: url.Values{
			"api-version": []string{"2020-10-01"},
		}.Encode(),
	}

	var adminApp AdminManagementApplicationDto
	_, err := client.Api.Execute(ctx, "PUT", apiUrl.String(), nil, nil, []int{http.StatusOK}, &adminApp)

	return &adminApp, err
}

func (client *AdminManagementApplicationClient) UnregisterAdminApplication(ctx context.Context, clientId string) (error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:  fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/adminApplications/%s", clientId),
		RawQuery: url.Values{
			"api-version": []string{"2020-10-01"},
		}.Encode(),
	}

	_, err := client.Api.Execute(ctx, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNoContent}, nil)
	
	return err
}
