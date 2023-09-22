package powerplatform_bapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

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
