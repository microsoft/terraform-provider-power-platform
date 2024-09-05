// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/constants"
)

func NewApplicationClient(api *api.ApiClient) ApplicationClient {
	return ApplicationClient{
		Api: api,
	}
}

type ApplicationClient struct {
	Api *api.ApiClient
}

func (client *ApplicationClient) DataverseExists(ctx context.Context, environmentId string) (bool, error) {

	env, err := client.getEnvironment(ctx, environmentId)
	if err != nil {
		return false, err
	}
	return env.Properties.LinkedEnvironmentMetadata.InstanceURL != "", nil
}

func (client *ApplicationClient) getEnvironment(ctx context.Context, environmentId string) (*EnvironmentIdDto, error) {

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	env := EnvironmentIdDto{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}

func (client *ApplicationClient) GetTenantApplications(ctx context.Context) ([]TenantApplicationDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   "/appmanagement/applicationPackages",
	}
	values := url.Values{
		"api-version": []string{"2022-03-01-preview"},
	}
	apiUrl.RawQuery = values.Encode()

	application := TenantApplicationArrayDto{}

	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &application)
	if err != nil {
		return nil, err
	}

	return application.Value, nil
}

func (client *ApplicationClient) GetApplicationsByEnvironmentId(ctx context.Context, environmentId string) ([]EnvironmentApplicationDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/appmanagement/environments/%s/applicationPackages", environmentId),
	}
	values := url.Values{
		"api-version": []string{"2022-03-01-preview"},
	}
	apiUrl.RawQuery = values.Encode()

	application := EnvironmentApplicationArrayDto{}

	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &application)
	if err != nil {
		return nil, err
	}

	return application.Value, nil
}

func (client *ApplicationClient) InstallApplicationInEnvironment(ctx context.Context, environmentId string, uniqueName string) (string, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/appmanagement/environments/%s/applicationPackages/%s/install", environmentId, uniqueName),
	}
	values := url.Values{
		"api-version": []string{"2022-03-01-preview"},
	}
	apiUrl.RawQuery = values.Encode()

	response, err := client.Api.Execute(ctx, "POST", apiUrl.String(), nil, nil, []int{http.StatusAccepted}, nil)
	if err != nil {
		return "", err
	}

	applicationId := ""
	if response.Response.StatusCode == http.StatusAccepted {
		operationLocationHeader := response.GetHeader(constants.HEADER_OPERATION_LOCATION)
		tflog.Debug(ctx, "Opeartion Location Header: "+operationLocationHeader)

		_, err = url.Parse(operationLocationHeader)
		if err != nil {
			tflog.Error(ctx, "Error parsing location header: "+err.Error())
		}

		for {
			lifecycleResponse := EnvironmentApplicationLifecycleDto{}
			_, err = client.Api.Execute(ctx, "GET", operationLocationHeader, nil, nil, []int{http.StatusOK}, &lifecycleResponse)
			if err != nil {
				return "", err
			}

			if lifecycleResponse.Status == "Succeeded" {
				parts := strings.Split(lifecycleResponse.CreatedDateTime, "/")
				if len(parts) > 0 {
					applicationId = parts[len(parts)-1]
				} else {
					return "", errors.New("can't parse application id from response " + lifecycleResponse.CreatedDateTime)
				}
				tflog.Debug(ctx, "Created Application Id: "+applicationId)
				break
			} else if lifecycleResponse.Status == "Failed" {
				return "", errors.New("application installation failed. status message: " + lifecycleResponse.Error.Message)
			}
		}
	} else if response.Response.StatusCode == http.StatusCreated {
		appCreatedResponse := EnvironmentApplicationLifecycleCreatedDto{}
		err = response.MarshallTo(&appCreatedResponse)
		if err != nil {
			return "", err
		}
		if appCreatedResponse.Properties.ProvisioningState != "Succeeded" {
			return "", errors.New("application installation failed. provisioning state: " + appCreatedResponse.Properties.ProvisioningState)
		}
		applicationId = appCreatedResponse.Name
	}

	return applicationId, nil
}
