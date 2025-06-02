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
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

func newApplicationClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}

type client struct {
	Api *api.Client
}

func (client *client) DataverseExists(ctx context.Context, environmentId string) (bool, error) {
	env, err := client.getEnvironment(ctx, environmentId)
	if err != nil {
		return false, err
	}
	return env.Properties.LinkedEnvironmentMetadata.InstanceURL != "", nil
}

func (client *client) AddApplicationUser(ctx context.Context, environmentId string, applicationId string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s/addAppUser", environmentId),
	}
	values := url.Values{
		"api-version": []string{"2020-10-01"},
	}
	apiUrl.RawQuery = values.Encode()

	// Create the request body
	requestBody := map[string]string{
		"servicePrincipalAppId": applicationId,
	}

	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, requestBody, []int{http.StatusOK}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (client *client) GetEnvironmentHostById(ctx context.Context, environmentId string) (string, error) {
	env, err := client.getEnvironment(ctx, environmentId)
	if err != nil {
		return "", err
	}
	if env.Properties.LinkedEnvironmentMetadata.InstanceURL == "" {
		return "", fmt.Errorf("environment %s does not have Dataverse", environmentId)
	}

	// Parse the instance URL to get the host
	instanceURL := env.Properties.LinkedEnvironmentMetadata.InstanceURL
	instanceURLParsed, err := url.Parse(instanceURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse instance URL %s: %v", instanceURL, err)
	}

	return instanceURLParsed.Host, nil
}

func (client *client) ApplicationUserExists(ctx context.Context, environmentId string, applicationId string) (bool, error) {
	// Get the environment host
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return false, err
	}

	// Create the Dataverse Web API URL to query for application users
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   "/api/data/v9.2/systemusers",
	}
	values := url.Values{}
	values.Add("$filter", fmt.Sprintf("applicationid eq %s", applicationId))
	apiUrl.RawQuery = values.Encode()

	// Make the request
	var response applicationUsersResponseDto
	resp, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNotFound, http.StatusForbidden}, &response)
	if err != nil {
		return false, err
	}

	// Handle forbidden or not found cases
	if resp.HttpResponse.StatusCode == http.StatusForbidden || resp.HttpResponse.StatusCode == http.StatusNotFound {
		tflog.Debug(ctx, fmt.Sprintf("Failed to query application users. Status: %d", resp.HttpResponse.StatusCode))
		return false, nil
	}

	// Check if the application user exists
	return len(response.Value) > 0, nil
}

func (client *client) GetApplicationUserSystemId(ctx context.Context, environmentId string, applicationId string) (string, error) {
	// Get the environment host
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return "", err
	}

	// Create the Dataverse Web API URL to query for application users
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   "/api/data/v9.2/systemusers",
	}
	values := url.Values{}
	values.Add("$select", "systemuserid")
	values.Add("$filter", fmt.Sprintf("applicationid eq %s", applicationId))
	apiUrl.RawQuery = values.Encode()

	// Make the request
	var response applicationUsersResponseDto
	resp, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNotFound, http.StatusForbidden}, &response)
	if err != nil {
		return "", err
	}

	// Handle forbidden or not found cases
	if resp.HttpResponse.StatusCode == http.StatusForbidden || resp.HttpResponse.StatusCode == http.StatusNotFound {
		tflog.Debug(ctx, fmt.Sprintf("Failed to query application users. Status: %d", resp.HttpResponse.StatusCode))
		return "", errors.New("failed to query application users")
	}

	// Check if the application user exists
	if len(response.Value) == 0 {
		return "", errors.New("application user not found")
	}

	return response.Value[0].SystemUserId, nil
}

func (client *client) DeactivateSystemUser(ctx context.Context, environmentId string, systemUserId string) error {
	// Get the environment host
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return err
	}

	// Get the application user to find the application ID
	appUser, err := client.getApplicationUserBySystemId(ctx, environmentId, systemUserId)
	if err != nil {
		return err
	}

	// Create the Dataverse Web API URL to deactivate the system user
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/v9.0/systemusers(%s)", systemUserId),
	}

	// The request body to disable the user
	requestBody := map[string]any{
		"isdisabled":    true,
		"applicationid": appUser.ApplicationId,
	}

	// Make the PATCH request to disable the user
	_, err = client.Api.Execute(ctx, nil, "PATCH", apiUrl.String(), nil, requestBody, []int{http.StatusNoContent, http.StatusOK}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (client *client) DeleteSystemUser(ctx context.Context, environmentId string, systemUserId string) error {
	// Get the environment host
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return err
	}

	// Create the Dataverse Web API URL to delete the system user
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/v9.2/systemusers(%s)", systemUserId),
	}

	// Make the request
	_, err = client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusNoContent, http.StatusOK}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (client *client) getEnvironment(ctx context.Context, environmentId string) (*environmentIdDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	env := environmentIdDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}

func (client *client) GetTenantApplications(ctx context.Context) ([]tenantApplicationDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   "/appmanagement/applicationPackages",
	}
	values := url.Values{
		"api-version": []string{"2022-03-01-preview"},
	}
	apiUrl.RawQuery = values.Encode()

	application := tenantApplicationArrayDto{}

	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &application)
	if err != nil {
		return nil, err
	}

	return application.Value, nil
}

func (client *client) GetApplicationsByEnvironmentId(ctx context.Context, environmentId string) ([]environmentApplicationDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/appmanagement/environments/%s/applicationPackages", environmentId),
	}
	values := url.Values{
		"api-version": []string{"2022-03-01-preview"},
	}
	apiUrl.RawQuery = values.Encode()

	application := environmentApplicationArrayDto{}

	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &application)
	if err != nil {
		return nil, err
	}

	return application.Value, nil
}

func (client *client) InstallApplicationInEnvironment(ctx context.Context, environmentId string, uniqueName string) (string, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/appmanagement/environments/%s/applicationPackages/%s/install", environmentId, uniqueName),
	}
	values := url.Values{
		"api-version": []string{"2022-03-01-preview"},
	}
	apiUrl.RawQuery = values.Encode()

	response, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, nil, []int{http.StatusAccepted}, nil)
	if err != nil {
		return "", err
	}

	applicationId := ""
	if response.HttpResponse.StatusCode == http.StatusAccepted {
		operationLocationHeader := response.GetHeader(constants.HEADER_OPERATION_LOCATION)
		tflog.Debug(ctx, "Opeartion Location Header: "+operationLocationHeader)

		_, err = url.Parse(operationLocationHeader)
		if err != nil {
			tflog.Error(ctx, "Error parsing location header: "+err.Error())
		}

		for {
			lifecycleResponse := environmentApplicationLifecycleDto{}
			response, err := client.Api.Execute(ctx, nil, "GET", operationLocationHeader, nil, nil, []int{http.StatusOK, http.StatusConflict}, &lifecycleResponse)
			if err != nil {
				return "", err
			}

			if response.HttpResponse.StatusCode == http.StatusConflict {
				tflog.Debug(ctx, "Lifecycle Operation HTTP Status: '"+response.HttpResponse.Status+"'")
				continue
			}

			if lifecycleResponse.Status == "Succeeded" {
				parts := strings.Split(lifecycleResponse.CreatedDateTime, "/")
				if len(parts) == 0 {
					return "", errors.New("can't parse application id from response " + lifecycleResponse.CreatedDateTime)
				}
				applicationId = parts[len(parts)-1]
				tflog.Debug(ctx, "Created Application Id: "+applicationId)
				break
			} else if lifecycleResponse.Status == "Failed" {
				return "", errors.New("application installation failed. status message: " + lifecycleResponse.Error.Message)
			}
		}
	} else if response.HttpResponse.StatusCode == http.StatusCreated {
		appCreatedResponse := environmentApplicationLifecycleCreatedDto{}
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

func (client *client) getApplicationUserBySystemId(ctx context.Context, environmentId string, systemUserId string) (*applicationUserDto, error) {
	// Get the environment host
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	// Create the Dataverse Web API URL to query for application users by system user ID
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   "/api/data/v9.0/systemusers",
	}
	values := url.Values{}
	values.Add("$select", "applicationid,systemuserid,applicationuserid,fullname")
	values.Add("$filter", fmt.Sprintf("systemuserid eq %s", systemUserId))
	apiUrl.RawQuery = values.Encode()

	// Make the request
	var response applicationUsersResponseDto
	resp, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNotFound, http.StatusForbidden}, &response)
	if err != nil {
		return nil, err
	}

	// Handle forbidden or not found cases
	if resp.HttpResponse.StatusCode == http.StatusForbidden || resp.HttpResponse.StatusCode == http.StatusNotFound {
		tflog.Debug(ctx, fmt.Sprintf("Failed to query application user by system ID. Status: %d", resp.HttpResponse.StatusCode))
		return nil, fmt.Errorf("application user not found for system user ID %s", systemUserId)
	}

	// Check if the application user exists
	if len(response.Value) == 0 {
		return nil, fmt.Errorf("application user not found for system user ID %s", systemUserId)
	}

	return &response.Value[0], nil
}
