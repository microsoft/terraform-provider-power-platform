package powerplatform

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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
	values.Add("$expand", "systemuserroles_association($select=roleid,name)")
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
	values.Add("$expand", "systemuserroles_association($select=roleid,name)")
	apiUrl.RawQuery = values.Encode()

	//https://org3bf4a377.crm17.dynamics.com/api/data/v9.2/systemusers?$select=firstname,lastname&$filter=azureactivedirectoryobjectid%20eq%20ad7b0121-6fca-440b-99ae-0d54d89a3ac7&$expand=systemuserroles_association($select=roleid,name)

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

	_, err := client.Api.Execute(ctx, "POST", apiUrl.String(), nil, userToCreate, []int{http.StatusOK}, nil)
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

func (client *UserClient) AssignSecurityRoles(ctx context.Context, environmentId, systemUserId string, securityRolesIds []string) (*UserDto, error) {
	//https://org3bf4a377.crm17.dynamics.com//api/data/v9.0/systemusers(debdcfd6-8adf-ee11-904c-002248f3fcc4)%2Fsystemuserroles_association%2F%24ref
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
			//"https://org3bf4a377.crm17.dynamics.com/api/data/v9.2/roles(a1801436-efd6-e811-a96e-000d3a3ab886)",
		}
		//user := UserDto{}
		_, err = client.Api.Execute(ctx, "POST", apiUrl.String(), nil, roleToassociate, []int{http.StatusNoContent}, nil)
		if err != nil {
			return nil, err
		}
	}

	//https://org3bf4a377.crm17.dynamics.com/api/data/v9.2/systemusers(debdcfd6-8adf-ee11-904c-002248f3fcc4)?$select=systemuserid,fullname&$expand=systemuserroles_association($select=roleid,name)
	//https://org3bf4a377.crm17.dynamics.com/api/data/v9.2/systemusers?$select=firstname,lastname&$filter=azureactivedirectoryobjectid%20eq%20ad7b0121-6fca-440b-99ae-0d54d89a3ac7&$expand=systemuserroles_association($select=roleid,name)

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
