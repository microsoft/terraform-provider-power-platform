package powerplatform

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

func NewUserClient(api *api.ApiClient) UserClient {
	return UserClient{
		Api: api,
	}
}

type UserClient struct {
	Api *api.ApiClient
}

func (client *UserClient) GetUsers(ctx context.Context, environmentId string) ([]UserDto, error) {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return nil, err
	}
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   "/api/data/v9.2/systemusers",
	}
	userArray := UserDtoArray{}
	_, err = client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &userArray)
	if err != nil {
		return nil, err
	}
	return userArray.Value, nil
}

func (client *UserClient) GetUserBySystemUserId(ctx context.Context, environmentId, systemUserId string) (*UserDto, error) {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return nil, err
	}
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   "/api/data/v9.2/systemusers(" + systemUserId + ")",
	}
	values := url.Values{}
	values.Add("$expand", "systemuserroles_association($select=roleid,name,ismanaged,_businessunitid_value)")
	apiUrl.RawQuery = values.Encode()

	user := UserDto{}
	_, err = client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (client *UserClient) GetUserByAadObjectId(ctx context.Context, environmentId, aadObjectId string) (*UserDto, error) {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return nil, err
	}
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   "/api/data/v9.2/systemusers",
	}
	values := url.Values{}
	values.Add("$filter", fmt.Sprintf("azureactivedirectoryobjectid eq %s", aadObjectId))
	values.Add("$expand", "systemuserroles_association($select=roleid,name,ismanaged,_businessunitid_value)")
	apiUrl.RawQuery = values.Encode()

	user := UserDtoArray{}
	_, err = client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &user)
	if err != nil {
		return nil, err
	}
	return &user.Value[0], nil
}

func (client *UserClient) CreateUser(ctx context.Context, environmentId, aadObjectId string) (*UserDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s/addUser", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	userToCreate := map[string]interface{}{
		"objectId": aadObjectId,
	}

	retryCount := 6 * 9 // 9 minutes of retries
	err := fmt.Errorf("")
	for retryCount > 0 {
		_, err = client.Api.Execute(ctx, "POST", apiUrl.String(), nil, userToCreate, []int{http.StatusOK}, nil)
		//the license assigment in Entra is async, so we need to wait for that to happen if a user is created in the same terraform run
		if err == nil || !strings.Contains(err.Error(), "userNotLicensed") {
			break
		}
		tflog.Debug(ctx, fmt.Sprintf("Error creating user: %s", err.Error()))
		time.Sleep(10 * time.Second)
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

func (client *UserClient) UpdateUser(ctx context.Context, environmentId, systemUserId string, userUpdate *UserDto) (*UserDto, error) {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return nil, err
	}
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   "/api/data/v9.2/systemusers(" + systemUserId + ")",
	}

	_, err = client.Api.Execute(ctx, "PATCH", apiUrl.String(), nil, userUpdate, []int{http.StatusOK}, nil)
	if err != nil {
		return nil, err
	}

	user, err := client.GetUserBySystemUserId(ctx, environmentId, systemUserId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (client *UserClient) RemoveSecurityRoles(ctx context.Context, environmentId, systemUserId string, securityRolesIds []string) (*UserDto, error) {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	for _, roleId := range securityRolesIds {
		apiUrl := &url.URL{
			Scheme: "https",
			Host:   strings.TrimPrefix(environmentUrl, "https://"),
			Path:   "/api/data/v9.2/systemusers(" + systemUserId + ")/systemuserroles_association/$ref",
		}
		values := url.Values{}
		values.Add("$id", fmt.Sprintf("%s/api/data/v9.2/roles(%s)", environmentUrl, roleId))
		apiUrl.RawQuery = values.Encode()

		_, err = client.Api.Execute(ctx, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusNoContent}, nil)
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

func (client *UserClient) AddSecurityRoles(ctx context.Context, environmentId, systemUserId string, securityRolesIds []string) (*UserDto, error) {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return nil, err
	}
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   "/api/data/v9.2/systemusers(" + systemUserId + ")/systemuserroles_association/$ref",
	}

	for _, roleId := range securityRolesIds {
		roleToassociate := map[string]interface{}{
			"@odata.id": fmt.Sprintf("%s/api/data/v9.2/roles(%s)", environmentUrl, roleId),
		}
		_, err = client.Api.Execute(ctx, "POST", apiUrl.String(), nil, roleToassociate, []int{http.StatusNoContent}, nil)
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

func (client *UserClient) GetEnvironmentUrlById(ctx context.Context, environmentId string) (string, error) {
	env, err := client.getEnvironment(ctx, environmentId)
	if err != nil {
		return "", err
	}
	environmentUrl := strings.TrimSuffix(env.Properties.LinkedEnvironmentMetadata.InstanceURL, "/")
	return environmentUrl, nil
}

func (client *UserClient) getEnvironment(ctx context.Context, environmentId string) (*EnvironmentIdDto, error) {

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("$expand", "permissions,properties.capacity,properties/billingPolicy")
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	env := EnvironmentIdDto{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}

func (client *UserClient) GetSecurityRoles(ctx context.Context, environmentId string) ([]SecurityRoleDto, error) {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return nil, err
	}
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   "/api/data/v9.2/roles",
	}
	//var values = url.Values{}
	//values.Add("$filter", fmt.Sprintf("_businessunitid_value eq %s", "84b8a367-62de-ee11-904d-002248c9a0e9"))
	//apiUrl.RawQuery = values.Encode()

	securityRoleArray := SecurityRoleDtoArray{}
	_, err = client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &securityRoleArray)
	if err != nil {
		return nil, err
	}
	return securityRoleArray.Value, nil
}
