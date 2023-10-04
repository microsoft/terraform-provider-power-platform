package powerplatform_api_bapi

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
	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/models"
)

var _ api.BapiClientInterface = &BapiClientApi{}

type BapiClientApi struct {
	BaseApi         api.ApiClientInterface
	Auth            BapiAuthInterface
	DataverseClient api.DataverseClientInterface
}

func (client *BapiClientApi) SetDataverseClient(dataverseClient api.DataverseClientInterface) {
	client.DataverseClient = dataverseClient
}

func (client *BapiClientApi) GetBase() api.ApiClientInterface {
	return client.BaseApi
}

func (client *BapiClientApi) Execute(ctx context.Context, method string, url string, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*api.ApiHttpResponse, error) {
	token, err := client.BaseApi.InitializeBase(ctx)
	if err != nil {
		return nil, err
	}
	return client.BaseApi.ExecuteBase(ctx, token, method, url, body, acceptableStatusCodes, responseObj)
}

func (client *BapiClientApi) GetEnvironment(ctx context.Context, environmentId string) (*models.EnvironmentDto, error) {

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BaseApi.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("$expand", "permissions,properties.capacity")
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	env := models.EnvironmentDto{}
	_, err := client.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &env)
	if err != nil {
		return nil, err
	}

	if env.Properties.LinkedEnvironmentMetadata.SecurityGroupId == "" {
		env.Properties.LinkedEnvironmentMetadata.SecurityGroupId = "00000000-0000-0000-0000-000000000000"
	}

	return &env, nil
}

func (client *BapiClientApi) GetEnvironments(ctx context.Context) ([]models.EnvironmentDto, error) {

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BaseApi.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments",
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	envArray := models.EnvironmentDtoArray{}
	_, err := client.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &envArray)
	if err != nil {
		return nil, err
	}

	return envArray.Value, nil
}

func (client *BapiClientApi) DeleteEnvironment(ctx context.Context, environmentId string) error {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BaseApi.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	environmentDelete := models.EnvironmentDeleteDto{
		Code:    "7", //Application
		Message: "Deleted using Terraform Provider for Power Platform",
	}

	_, err := client.Execute(ctx, "DELETE", apiUrl.String(), environmentDelete, []int{http.StatusAccepted}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (client *BapiClientApi) CreateEnvironment(ctx context.Context, environment models.EnvironmentCreateDto) (*models.EnvironmentDto, error) {

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BaseApi.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/environments",
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()
	apiResponse, err := client.Execute(ctx, "POST", apiUrl.String(), environment, []int{http.StatusAccepted, http.StatusCreated}, nil)
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

			lifecycleResponse := models.EnvironmentLifecycleDto{}
			apiResponse, err = client.Execute(ctx, "GET", locationHeader, nil, []int{http.StatusOK}, &lifecycleResponse)
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
		envCreatedResponse := models.EnvironmentLifecycleCreatedDto{}
		apiResponse.MarshallTo(&envCreatedResponse)
		if envCreatedResponse.Properties.ProvisioningState != "Succeeded" {
			return nil, errors.New("environment creation failed. provisioning state: " + envCreatedResponse.Properties.ProvisioningState)
		}
		createdEnvironmentId = envCreatedResponse.Name
	}

	env, err := client.GetEnvironment(ctx, createdEnvironmentId)
	if err != nil {
		return &models.EnvironmentDto{}, errors.New("environment not found")
	}
	return env, err
}

func (client *BapiClientApi) UpdateEnvironment(ctx context.Context, environmentId string, environment models.EnvironmentDto) (*models.EnvironmentDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BaseApi.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2022-05-01")
	apiUrl.RawQuery = values.Encode()
	_, err := client.Execute(ctx, "PATCH", apiUrl.String(), environment, []int{http.StatusAccepted}, nil)
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

func (client *BapiClientApi) GetPowerApps(ctx context.Context, environmentId string) ([]models.PowerAppBapi, error) {
	envs, err := client.GetEnvironments(ctx)
	if err != nil {
		return nil, err
	}
	apps := make([]models.PowerAppBapi, 0)
	for _, env := range envs {
		apiUrl := &url.URL{
			Scheme: "https",
			Host:   client.BaseApi.GetConfig().Urls.PowerAppsUrl,
			Path:   fmt.Sprintf("/providers/Microsoft.PowerApps/scopes/admin/environments/%s/apps", env.Name),
		}
		values := url.Values{}
		values.Add("api-version", "2023-06-01")
		apiUrl.RawQuery = values.Encode()

		appsArray := models.PowerAppDtoArray{}
		_, err := client.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &appsArray)
		if err != nil {
			return nil, err
		}
		apps = append(apps, appsArray.Value...)

	}
	return apps, nil
}

func (client *BapiClientApi) GetConnectors(ctx context.Context) ([]models.ConnectorDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BaseApi.GetConfig().Urls.PowerAppsUrl,
		Path:   "/providers/Microsoft.PowerApps/apis",
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	values.Add("showApisWithToS", "true")
	values.Add("hideDlpExemptApis", "true")
	values.Add("showAllDlpEnforceableApis", "true")
	values.Add("$filter", "environment eq '~Default'")
	apiUrl.RawQuery = values.Encode()

	connectorArray := models.ConnectorDtoArray{}
	_, err := client.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &connectorArray)
	if err != nil {
		return nil, err
	}

	apiUrl = &url.URL{
		Scheme: "https",
		Host:   client.BaseApi.GetConfig().Urls.BapiUrl,
		Path:   "/providers/PowerPlatform.Governance/v1/connectors/metadata/unblockable",
	}
	unblockableConnectorArray := []models.UnblockableConnectorDto{}
	_, err = client.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &unblockableConnectorArray)
	if err != nil {
		return nil, err
	}

	for inx, connector := range connectorArray.Value {
		for _, unblockableConnector := range unblockableConnectorArray {
			if connector.Id == unblockableConnector.Id {
				connectorArray.Value[inx].Properties.Unblockable = unblockableConnector.Metadata.Unblockable
			}
		}
	}

	apiUrl = &url.URL{
		Scheme: "https",
		Host:   client.BaseApi.GetConfig().Urls.BapiUrl,
		Path:   "/providers/PowerPlatform.Governance/v1/connectors/metadata/virtual",
	}
	virtualConnectorArray := []models.VirtualConnectorDto{}
	_, err = client.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &virtualConnectorArray)
	if err != nil {
		return nil, err
	}
	for _, virutualConnector := range virtualConnectorArray {
		connectorArray.Value = append(connectorArray.Value, models.ConnectorDto{
			Id:   virutualConnector.Id,
			Name: virutualConnector.Metadata.Name,
			Type: virutualConnector.Metadata.Type,
			Properties: models.ConnectorPropertiesDto{
				DisplayName: virutualConnector.Metadata.DisplayName,
				Unblockable: false,
				Tier:        "Built-in",
				Publisher:   "Microsoft",
				Description: "",
			},
		})
	}

	for inx, connector := range connectorArray.Value {
		nameSplit := strings.Split(connector.Id, "/")
		connectorArray.Value[inx].Name = nameSplit[len(nameSplit)-1]
	}

	return connectorArray.Value, nil
}
