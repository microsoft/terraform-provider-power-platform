// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package authorization

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
)

func newUserClient(apiClient *api.Client) client {
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

func (client *client) GetUsers(ctx context.Context, environmentId string) ([]userDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   "/api/data/v9.2/systemusers",
	}
	userArray := userDtoArray{}
	_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &userArray)
	if err != nil {
		return nil, err
	}
	return userArray.Value, nil
}

func (client *client) GetUserBySystemUserId(ctx context.Context, environmentId, systemUserId string) (*userDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   "/api/data/v9.2/systemusers(" + systemUserId + ")",
	}
	values := url.Values{}
	values.Add("$expand", "systemuserroles_association($select=roleid,name,ismanaged,_businessunitid_value)")
	apiUrl.RawQuery = values.Encode()

 user := userDto{}
 main
	_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &user)
	if err != nil {
		var unexpectedError *customerrors.UnexpectedHttpStatusCodeError
		if errors.As(err, &unexpectedError) && unexpectedError.StatusCode == http.StatusNotFound {
			return nil, customerrors.WrapIntoProviderError(err, customerrors.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("User with systemUserId %s not found", systemUserId))
		}
		return nil, err
	}
	return &user, nil
}

func (client *client) GetUserByAadObjectId(ctx context.Context, environmentId, aadObjectId string) (*userDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   "/api/data/v9.2/systemusers",
	}
	values := url.Values{}
	values.Add("$filter", fmt.Sprintf("azureactivedirectoryobjectid eq %s", aadObjectId))
	values.Add("$expand", "systemuserroles_association($select=roleid,name,ismanaged,_businessunitid_value)")
	apiUrl.RawQuery = values.Encode()

	user := userDtoArray{}
	_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &user)
	if err != nil {
		var httpError *customerrors.UnexpectedHttpStatusCodeError
		if errors.As(err, &httpError) && httpError.StatusCode == http.StatusNotFound {
			return nil, customerrors.WrapIntoProviderError(err, customerrors.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("User with aadObjectId %s not found", aadObjectId))
		}

		return nil, err
	}
	return &user.Value[0], nil
}

func (client *client) CreateUser(ctx context.Context, environmentId, aadObjectId string) (*userDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s/addUser", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	userToCreate := map[string]any{
		"objectId": aadObjectId,
	}

	// 9 minutes of retries.
	retryCount := 6 * 9
	err := fmt.Errorf("")
	for retryCount > 0 {
		_, err = client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, userToCreate, []int{http.StatusOK}, nil)
		// the license assignment in Entra is async, so we need to wait for that to happen if a user is created in the same terraform run.
		if err == nil || !strings.Contains(err.Error(), "userNotLicensed") {
			break
		}
		tflog.Debug(ctx, fmt.Sprintf("Error creating user: %s", err.Error()))
		err = client.Api.SleepWithContext(ctx, 10*time.Second)
		if err != nil {
			return nil, err
		}

		retryCount--
	}
	if err != nil {
		return nil, err
	}

	user, err := client.GetUserByAadObjectId(ctx, environmentId, aadObjectId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (client *client) UpdateUser(ctx context.Context, environmentId, systemUserId string, userUpdate *userDto) (*userDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   "/api/data/v9.2/systemusers(" + systemUserId + ")",
	}

	_, err = client.Api.Execute(ctx, nil, "PATCH", apiUrl.String(), nil, userUpdate, []int{http.StatusOK}, nil)
	if err != nil {
		return nil, err
	}

	user, err := client.GetUserBySystemUserId(ctx, environmentId, systemUserId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (client *client) DeleteUser(ctx context.Context, environmentId, systemUserId string) error {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return err
	}
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   "/api/data/v9.2/systemusers(" + systemUserId + ")",
	}

	_, err = client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusNoContent}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (client *client) RemoveSecurityRoles(ctx context.Context, environmentId, systemUserId string, securityRolesIds []string) (*userDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	for _, roleId := range securityRolesIds {
		apiUrl := &url.URL{
			Scheme: constants.HTTPS,
			Host:   environmentHost,
			Path:   "/api/data/v9.2/systemusers(" + systemUserId + ")/systemuserroles_association/$ref",
		}
		values := url.Values{}
		values.Add("$id", fmt.Sprintf("https://%s/api/data/v9.2/roles(%s)", environmentHost, roleId))
		apiUrl.RawQuery = values.Encode()

		_, err = client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusNoContent}, nil)
		if err != nil {
			return nil, err
		}
	}

	user, err := client.GetUserBySystemUserId(ctx, environmentId, systemUserId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (client *client) AddSecurityRoles(ctx context.Context, environmentId, systemUserId string, securityRolesIds []string) (*userDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   "/api/data/v9.2/systemusers(" + systemUserId + ")/systemuserroles_association/$ref",
	}

	for _, roleId := range securityRolesIds {
		roleToassociate := map[string]any{
			"@odata.id": fmt.Sprintf("https://%s/api/data/v9.2/roles(%s)", environmentHost, roleId),
		}
		_, err = client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, roleToassociate, []int{http.StatusNoContent}, nil)
		if err != nil {
			return nil, err
		}
	}
	user, err := client.GetUserBySystemUserId(ctx, environmentId, systemUserId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (client *client) GetEnvironmentHostById(ctx context.Context, environmentId string) (string, error) {
	env, err := client.getEnvironment(ctx, environmentId)
	if err != nil {
		return "", err
	}
	environmentUrl := strings.TrimSuffix(env.Properties.LinkedEnvironmentMetadata.InstanceURL, "/")
	if environmentUrl == "" {
		return "", customerrors.WrapIntoProviderError(nil, customerrors.ERROR_ENVIRONMENT_URL_NOT_FOUND, "environment url not found, please check if the environment has dataverse linked")
	}
	envUrl, err := url.Parse(environmentUrl)
	if err != nil {
		return "", err
	}
	return envUrl.Host, nil
}

func (client *client) getEnvironment(ctx context.Context, environmentId string) (*environmentIdDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("$expand", "permissions,properties.capacity,properties/billingPolicy")
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	env := environmentIdDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
	if err != nil {
		var httpError *customerrors.UnexpectedHttpStatusCodeError
		if errors.As(err, &httpError) && httpError.StatusCode == http.StatusNotFound {
			return nil, customerrors.WrapIntoProviderError(err, customerrors.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("environment %s not found", environmentId))
		}
		return nil, err
	}

	return &env, nil
}

func (client *client) GetSecurityRoles(ctx context.Context, environmentId, businessUnitId string) ([]securityRoleDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   "/api/data/v9.2/roles",
	}
	if businessUnitId != "" {
		var values = url.Values{}
		values.Add("$filter", fmt.Sprintf("_businessunitid_value eq %s", businessUnitId))
		apiUrl.RawQuery = values.Encode()
	}
	securityRoleArray := securityRoleDtoArray{}
	_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &securityRoleArray)
	if err != nil {
		var httpError *customerrors.UnexpectedHttpStatusCodeError
		if errors.As(err, &httpError) && httpError.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, fmt.Sprintf("Error getting security roles: %s", err.Error()))
			return nil, customerrors.WrapIntoProviderError(err, customerrors.ERROR_OBJECT_NOT_FOUND, "security roles not found")
		}
		return nil, err
	}
	return securityRoleArray.Value, nil
}
