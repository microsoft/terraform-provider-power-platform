package powerplatform

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	solution "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/solution"
)

func NewEnvironmentClient(bapi *api.BapiClientApi, dv *api.DataverseClientApi) EnvironmentClient {
	return EnvironmentClient{
		bapiClient:      bapi,
		dataverseClient: dv,
		solutionClient:  solution.NewSolutionClient(bapi, dv),
	}
}

type EnvironmentClient struct {
	bapiClient      *api.BapiClientApi
	dataverseClient *api.DataverseClientApi
	solutionClient  solution.SolutionClient
}

func (client *EnvironmentClient) GetEnvironmentUrlById(ctx context.Context, environmentId string) (string, error) {
	env, err := client.GetEnvironment(ctx, environmentId)
	if err != nil {
		return "", err
	}
	environmentUrl := strings.TrimSuffix(env.Properties.LinkedEnvironmentMetadata.InstanceURL, "/")
	return environmentUrl, nil
}

func (client *EnvironmentClient) GetEnvironment(ctx context.Context, environmentId string) (*EnvironmentDto, error) {

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.bapiClient.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("$expand", "permissions,properties.capacity")
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	env := EnvironmentDto{}
	_, err := client.bapiClient.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &env)
	if err != nil {
		return nil, err
	}

	if env.Properties.LinkedEnvironmentMetadata.SecurityGroupId == "" {
		env.Properties.LinkedEnvironmentMetadata.SecurityGroupId = "00000000-0000-0000-0000-000000000000"
	}

	return &env, nil
}

func (client *EnvironmentClient) DeleteEnvironment(ctx context.Context, environmentId string) error {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.bapiClient.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	environmentDelete := EnvironmentDeleteDto{
		Code:    "7", //Application
		Message: "Deleted using Terraform Provider for Power Platform",
	}

	_, err := client.bapiClient.Execute(ctx, "DELETE", apiUrl.String(), environmentDelete, []int{http.StatusAccepted}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (client *EnvironmentClient) CreateEnvironment(ctx context.Context, environment EnvironmentCreateDto) (*EnvironmentDto, error) {

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.bapiClient.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/environments",
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()
	apiResponse, err := client.bapiClient.Execute(ctx, "POST", apiUrl.String(), environment, []int{http.StatusAccepted, http.StatusCreated}, nil)
	if err != nil {
		return nil, err
	}

	tflog.Debug(ctx, "Environment Creation Operation HTTP Status: '"+apiResponse.Response.Status+"'")

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

			lifecycleResponse := EnvironmentLifecycleDto{}
			apiResponse, err = client.bapiClient.Execute(ctx, "GET", locationHeader, nil, []int{http.StatusOK}, &lifecycleResponse)
			if err != nil {
				return nil, err
			}

			time.Sleep(retryAfter)

			tflog.Debug(ctx, "Environment Creation Operation State: '"+lifecycleResponse.State.Id+"'")
			tflog.Debug(ctx, "Environment Creation Operation HTTP Status: '"+apiResponse.Response.Status+"'")

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
		envCreatedResponse := EnvironmentLifecycleCreatedDto{}
		apiResponse.MarshallTo(&envCreatedResponse)
		if envCreatedResponse.Properties.ProvisioningState != "Succeeded" {
			return nil, errors.New("environment creation failed. provisioning state: " + envCreatedResponse.Properties.ProvisioningState)
		}
		createdEnvironmentId = envCreatedResponse.Name
	}

	env, err := client.GetEnvironment(ctx, createdEnvironmentId)
	if err != nil {
		return &EnvironmentDto{}, fmt.Errorf("environment '%s' not found. '%s'", createdEnvironmentId, err)
	}
	return env, err
}

func (client *EnvironmentClient) UpdateEnvironment(ctx context.Context, environmentId string, environment EnvironmentDto) (*EnvironmentDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.bapiClient.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2022-05-01")
	apiUrl.RawQuery = values.Encode()
	_, err := client.bapiClient.Execute(ctx, "PATCH", apiUrl.String(), environment, []int{http.StatusAccepted}, nil)
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

	return nil, fmt.Errorf("environment '%s' not found", environmentId)
}

func (client *EnvironmentClient) GetEnvironments(ctx context.Context) ([]EnvironmentDto, error) {

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.bapiClient.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments",
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	envArray := EnvironmentDtoArray{}
	_, err := client.bapiClient.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &envArray)
	if err != nil {
		return nil, err
	}

	return envArray.Value, nil
}

func (client *EnvironmentClient) GetDefaultCurrencyForEnvironment(ctx context.Context, environmentId string) (*TransactionCurrencyDto, error) {
	orgSettings := OrganizationSettingsArrayDto{}
	err := client.solutionClient.GetTableData(ctx, environmentId, "organizations", "", &orgSettings)
	if err != nil {
		return nil, err
	} else {
		values := url.Values{}
		values.Add("$filter", "transactioncurrencyid eq "+orgSettings.Value[0].BaseCurrencyId)

		currencies := TransactionCurrencyArrayDto{}
		err := client.solutionClient.GetTableData(ctx, environmentId, "transactioncurrencies", values.Encode(), &currencies)
		if err != nil {
			return nil, err
		} else {
			if currencies.Value != nil && len(currencies.Value) >= 1 {
				return &currencies.Value[0], nil
			} else {
				return nil, fmt.Errorf("no default currency found for environment %s", environmentId)
			}
		}
	}
}
