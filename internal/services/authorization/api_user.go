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
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers/array"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment"
)

func newUserClient(apiClient *api.Client) client {
	return client{
		Api:               apiClient,
		environmentClient: environment.NewEnvironmentClient(apiClient),
	}
}

type client struct {
	Api               *api.Client
	environmentClient environment.Client
}

func (client *client) EnvironmentHasDataverse(ctx context.Context, environmentId string) (bool, error) {
	env, err := client.environmentClient.GetEnvironment(ctx, environmentId)
	if err != nil {
		return false, err
	}
	return env.Properties.LinkedEnvironmentMetadata != nil, nil
}

func (client *client) DataverseExists(ctx context.Context, environmentId string) (bool, error) {
	env, err := client.getEnvironment(ctx, environmentId)
	if err != nil {
		return false, err
	}
	return env.Properties.LinkedEnvironmentMetadata.InstanceURL != "", nil
}

func (client *client) GetDataverseUsers(ctx context.Context, environmentId string) ([]userDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   "/api/data/v9.2/systemusers",
	}
	userArray := userArrayDto{}
	_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &userArray)
	if err != nil {
		return nil, err
	}
	return userArray.Value, nil
}

func (client *client) GetDataverseUserBySystemUserId(ctx context.Context, environmentId, systemUserId string) (*userDto, error) {
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

func (client *client) GetEnvironmentUserByAadObjectId(ctx context.Context, environmentId, aadObjectId string) (*userDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s/roleAssignments", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()

	respObj := EnvironmentUserGetResponseDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &respObj)
	if err != nil {
		return nil, err
	}
	user := &userDto{
		SecurityRoles: []securityRoleDto{},
	}

	for _, roleAssignment := range respObj.Value {
		if roleAssignment.Properties.Principal.Id == aadObjectId {
			isAdminRole := roleAssignment.Properties.RoleDefinition.Name == "EnvironmentAdmin"
			isMakerRole := roleAssignment.Properties.RoleDefinition.Name == "EnvironmentMaker"

			user.Id = roleAssignment.Properties.Principal.Id
			user.AadObjectId = roleAssignment.Properties.Principal.Id
			user.DomainName = roleAssignment.Properties.Principal.DisplayName
			if isAdminRole {
				user.SecurityRoles = append(user.SecurityRoles, securityRoleDto{
					Name:   roleAssignment.Name,
					RoleId: ROLE_ENVIRONMENT_ADMIN,
				})
			} else if isMakerRole {
				user.SecurityRoles = append(user.SecurityRoles, securityRoleDto{
					Name:   roleAssignment.Name,
					RoleId: ROLE_ENVIRONMENT_MAKER,
				})
			}
		}
	}

	return user, nil
}

func (client *client) GetDataverseUserByAadObjectId(ctx context.Context, environmentId, aadObjectId string) (*userDto, error) {
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

	user := userArrayDto{}
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

func (client *client) RemoveEnvironmentUserSecurityRoles(ctx context.Context, environmentId, aadObjectId string, securityRoles []string, savedRoles []securityRoleDto) (*userDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s/modifyRoleAssignments", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()

	userRead, err := client.GetEnvironmentUserByAadObjectId(ctx, environmentId, aadObjectId)
	if err != nil {
		return nil, err
	}
	userDisplayName := userRead.DomainName

	remove := EnvironmentUserRequestDto{
		Remove: []RoleDefinitionDto{},
	}

	for _, role := range securityRoles {
		savedRoleData := array.Find(savedRoles, func(roleDto securityRoleDto) bool {
			return roleDto.RoleId == role
		})

		remove.Remove = append(remove.Remove, RoleDefinitionDto{
			Id: fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s/roleAssignments/%s", environmentId, savedRoleData.Name),
		})
	}

	_, err = client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, remove, []int{http.StatusOK}, nil)
	if err != nil {
		return nil, err
	}

	userRead, err = client.GetEnvironmentUserByAadObjectId(ctx, environmentId, aadObjectId)
	if err != nil {
		return nil, err
	}

	user := userDto{
		Id:            aadObjectId,
		AadObjectId:   aadObjectId,
		DomainName:    userDisplayName,
		SecurityRoles: userRead.SecurityRoles,
	}
	return &user, nil
}

func (client *client) AddEnvironmentUserSecurityRoles(ctx context.Context, environmentId, aadObjectId string, securityRoles []string) (*userDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s/modifyRoleAssignments", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()

	add := EnvironmentUserRequestDto{
		Add: []AddItemRequestDto{},
	}

	for _, role := range securityRoles {
		role = strings.ToLower(strings.ReplaceAll(role, " ", ""))
		add.Add = append(add.Add, AddItemRequestDto{
			Properties: PropertiesDto{
				RoleDefinition: RoleDefinitionDto{
					Id: fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s/roleDefinitions/%s", environmentId, role),
				},
				Principal: PrincipalDto{
					Id:   aadObjectId,
					Type: "User",
				},
			},
		})
	}

	respObj := EnvironmentUserResponseDto{}
	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, add, []int{http.StatusOK}, &respObj)
	if err != nil {
		return nil, err
	}

	user := userDto{
		Id:            aadObjectId,
		AadObjectId:   aadObjectId,
		DomainName:    respObj.Add[0].RoleAssignment.Properties.Principal.DisplayName,
		SecurityRoles: []securityRoleDto{},
	}

	for _, roleAssignment := range respObj.Add {
		user.SecurityRoles = append(user.SecurityRoles, securityRoleDto{
			RoleId: roleAssignment.RoleAssignment.Properties.RoleDefinition.Name,
			Name:   roleAssignment.RoleAssignment.Name,
		})
	}
	return &user, nil
}

func (client *client) CreateEnvironmentUser(ctx context.Context, environmentId, aadObjectId string, securityRoles []string) (*userDto, error) {
	// adding security roles and creating user is the same API call.
	return client.AddEnvironmentUserSecurityRoles(ctx, environmentId, aadObjectId, securityRoles)
}

func (client *client) CreateDataverseUser(ctx context.Context, environmentId, aadObjectId string) (*userDto, error) {
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

	user, err := client.GetDataverseUserByAadObjectId(ctx, environmentId, aadObjectId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (client *client) UpdateDataverseUser(ctx context.Context, environmentId, systemUserId string, userUpdate *userDto) (*userDto, error) {
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

	user, err := client.GetDataverseUserBySystemUserId(ctx, environmentId, systemUserId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (client *client) DeleteDataverseUser(ctx context.Context, environmentId, systemUserId string) error {
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

func (client *client) RemoveDataverseSecurityRoles(ctx context.Context, environmentId, systemUserId string, securityRolesIds []string) (*userDto, error) {
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
			if strings.Contains(err.Error(), "0x80060888") && strings.Contains(err.Error(), roleId) {
				return nil, fmt.Errorf("Role with id '%s' is not valid", roleId)
			}
			return nil, err
		}
	}

	user, err := client.GetDataverseUserBySystemUserId(ctx, environmentId, systemUserId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (client *client) AddDataverseSecurityRoles(ctx context.Context, environmentId, systemUserId string, securityRolesIds []string) (*userDto, error) {
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
			if strings.Contains(err.Error(), "0x80060888") && strings.Contains(err.Error(), roleId) {
				return nil, fmt.Errorf("Role with id '%s' is not valid", roleId)
			}
			return nil, err
		}
	}
	user, err := client.GetDataverseUserBySystemUserId(ctx, environmentId, systemUserId)
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

func (client *client) GetDataverseSecurityRoles(ctx context.Context, environmentId, businessUnitId string) ([]securityRoleDto, error) {
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
	securityRoleArray := securityRoleArrayDto{}
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
