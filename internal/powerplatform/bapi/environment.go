package powerplatform_bapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

	apiResponse, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}

	envArray := models.EnvironmentDtoArray{}
	err = apiResponse.MarshallTo(&envArray)
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

	apiResponse, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}

	env := models.EnvironmentDto{}
	err = apiResponse.MarshallTo(&env)
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

	apiResponse, err := client.doRequest(request)
	if err != nil {
		return err
	}

	err = apiResponse.ValidateStatusCode(http.StatusAccepted)
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

	apiResponse, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}

	tflog.Debug(ctx, "Environment Creation Opeartion HTTP Status: '"+apiResponse.Response.Status+"'")

	createdEnvironmentId := ""
	if apiResponse.Response.StatusCode == http.StatusAccepted {

		locationHeader := apiResponse.GetHeader("Location")
		tflog.Debug(ctx, "Location Header: "+locationHeader)

		_, err = url.Parse(locationHeader)
		if err != nil {
			tflog.Error(ctx, "Error parsing location header: "+err.Error())
		}

		retryHeader := apiResponse.GetHeader("Retry-After")
		tflog.Debug(ctx, "Retry Header: "+retryHeader)
		retryAfter, err := time.ParseDuration(retryHeader)
		if err != nil {
			retryAfter = time.Duration(5) * time.Second
		} else {
			retryAfter = retryAfter * time.Second
		}

		for {
			request, err = http.NewRequestWithContext(ctx, "GET", locationHeader, bytes.NewReader(body))
			if err != nil {
				return nil, err
			}

			apiResponse, err = client.doRequest(request)
			if err != nil {
				return nil, err
			}

			lifecycleResponse := models.EnvironmentLifecycleDto{}
			err = apiResponse.MarshallTo(&lifecycleResponse)
			if err != nil {
				return nil, err
			}

			time.Sleep(retryAfter)

			tflog.Debug(ctx, "Environment Creation Opeartion State: '"+lifecycleResponse.State.Id+"'")
			tflog.Debug(ctx, "Environment Creation Opeartion HTTP Status: '"+apiResponse.Response.Status+"'")

			if lifecycleResponse.State.Id == "Succeeded" {
				parts := strings.Split(lifecycleResponse.Links.Environment.Path, "/")
				if len(parts) > 0 {
					createdEnvironmentId = parts[len(parts)-1]
				} else {
					return nil, errors.New("can't parse environment id from response " + lifecycleResponse.Links.Environment.Path)
				}
				tflog.Debug(ctx, "Created Environment Id: "+createdEnvironmentId)
				break
			}
		}
	} else if apiResponse.Response.StatusCode == http.StatusCreated {
		envCreatedResponse := models.EnvironmentLifecycleCreatedDto{}
		apiResponse.MarshallTo(&envCreatedResponse)
		if envCreatedResponse.Properties.ProvisioningState != "Succeeded" {
			return nil, errors.New("environment creation failed. privisioning state: " + envCreatedResponse.Properties.ProvisioningState)
		}
		createdEnvironmentId = envCreatedResponse.Name
	}

	env, err := client.GetEnvironment(ctx, createdEnvironmentId)
	if err != nil {
		return &models.EnvironmentDto{}, errors.New("environment not found")
	}
	return env, err
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
	apiResponse, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}
	err = apiResponse.ValidateStatusCode(http.StatusAccepted)
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
