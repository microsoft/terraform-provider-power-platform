package powerplatform_bapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi/models"
)

func (client *ApiClient) GetEnvironments(ctx context.Context) ([]models.EnvironmentDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   "api.bap.microsoft.com",
		Path:   "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments",
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()
	request, err := http.NewRequestWithContext(ctx, "GET", apiUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}

	envArray := models.EnvironmentDtoArray{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&envArray)
	if err != nil {
		return nil, err
	}

	return envArray.Value, nil
}

func (client *ApiClient) GetEnvironment(ctx context.Context, environmentId string) (*models.EnvironmentDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   "api.bap.microsoft.com",
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("$expand", "permissions,properties.capacity")
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()
	request, err := http.NewRequestWithContext(ctx, "GET", apiUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}

	env := models.EnvironmentDto{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&env)
	if err != nil {
		return nil, err
	}

	if env.Properties.LinkedEnvironmentMetadata.SecurityGroupId == "" {
		env.Properties.LinkedEnvironmentMetadata.SecurityGroupId = "00000000-0000-0000-0000-000000000000"
	}

	return &env, nil
}

func (client *ApiClient) DeleteEnvironment(ctx context.Context, environmentId string) error {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   "api.bap.microsoft.com",
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	environmentDelete := models.EnvironmentDeleteDto{
		Code:    "7", //Application
		Message: "Deleted using Terraform Provider for Power Platform",
	}
	body, err := json.Marshal(environmentDelete)
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, "DELETE", apiUrl.String(), bytes.NewReader(body))

	if err != nil {
		return err
	}

	_, err = client.doRequest(request)
	if err != nil {
		return err
	}

	return nil
}

func (client *ApiClient) CreateEnvironment(ctx context.Context, environment models.EnvironmentCreateDto) (*models.EnvironmentDto, error) {
	body, err := json.Marshal(environment)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   "api.bap.microsoft.com",
		Path:   "/providers/Microsoft.BusinessAppPlatform/environments",
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()
	request, err := http.NewRequestWithContext(ctx, "POST", apiUrl.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	_, err = client.doRequest(request)
	if err != nil {
		return nil, err
	}

	time.Sleep(10 * time.Second)

	environments, err := client.GetEnvironments(ctx)
	if err != nil {
		return nil, err
	}

	for _, env := range environments {
		if env.Location == environment.Location && env.Properties.DisplayName == environment.Properties.DisplayName {
			for {
				createdEnv, err := client.GetEnvironment(ctx, env.Name)
				if err != nil {
					return nil, err
				}
				tflog.Info(ctx, "Environment State: '"+createdEnv.Properties.States.Management.Id+"'")
				time.Sleep(1 * time.Second)
				if createdEnv.Properties.States.Management.Id != "Running" {
					return createdEnv, nil
				}

			}
		}
	}
	return &models.EnvironmentDto{}, errors.New("environment not found")
}

func (client *ApiClient) UpdateEnvironment(ctx context.Context, environmentId string, environment models.EnvironmentDto) (*models.EnvironmentDto, error) {
	body, err := json.Marshal(environment)
	if err != nil {
		return nil, err
	}
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   "api.bap.microsoft.com",
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2022-05-01")
	apiUrl.RawQuery = values.Encode()
	request, err := http.NewRequestWithContext(ctx, "PATCH", apiUrl.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	_, err = client.doRequest(request)
	if err != nil {
		return nil, err
	}

	time.Sleep(10 * time.Second)

	environments, err := client.GetEnvironments(ctx)
	if err != nil {
		return nil, err
	}

	for _, env := range environments {
		if env.Name == environmentId {
			for {
				createdEnv, err := client.GetEnvironment(ctx, env.Name)
				if err != nil {
					return nil, err
				}
				tflog.Info(ctx, "Environment State: '"+createdEnv.Properties.States.Management.Id+"'")
				time.Sleep(3 * time.Second)
				if createdEnv.Properties.States.Management.Id == "Ready" {

					return createdEnv, nil
				}

			}
		}
	}

	return nil, errors.New("environment not found")
}
