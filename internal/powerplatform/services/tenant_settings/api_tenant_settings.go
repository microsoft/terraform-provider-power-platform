package powerplatform

import (
	"context"
	"net/http"
	"net/url"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

func NewTenantSettingsClient(bapi *api.BapiClientApi) TenantSettingsClient {
	return TenantSettingsClient{
		//bapiClient: bapi,
		api: nil,
	}
}

type TenantSettingsClient struct {
	//bapiClient *api.BapiClientApi
	api *api.ApiClientBase
}

func (client *TenantSettingsClient) GetTenantSettings(ctx context.Context) (*TenantSettingsDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		//Host:   client.bapiClient.GetConfig().Urls.BapiUrl,
		Host: client.api.GetConfig().Urls.BapiUrl,
		Path: "/providers/Microsoft.BusinessAppPlatform/listTenantSettings",
	}

	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	tenantSettings := TenantSettingsDto{}
	//_, err := client.bapiClient.Execute(ctx, "POST", apiUrl.String(), nil, nil, []int{http.StatusOK}, &tenantSettings)
	_, err := client.api.ExecuteBase(ctx, "TOKEN!!!!", "POST", apiUrl.String(), nil, nil, []int{http.StatusOK}, &tenantSettings)
	if err != nil {
		return nil, err
	}
	return &tenantSettings, nil
}

func (client *TenantSettingsClient) UpdateTenantSettings(ctx context.Context, tenantSettings TenantSettingsDto) (*TenantSettingsDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		//Host:   client.bapiClient.GetConfig().Urls.BapiUrl,
		Host: client.api.GetConfig().Urls.BapiUrl,
		Path: "/providers/Microsoft.BusinessAppPlatform/scopes/admin/updateTenantSettings",
	}

	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	_, err := client.bapiClient.Execute(ctx, "POST", apiUrl.String(), nil, tenantSettings, []int{http.StatusOK}, &tenantSettings)
	if err != nil {
		return nil, err
	}
	return &tenantSettings, nil
}
