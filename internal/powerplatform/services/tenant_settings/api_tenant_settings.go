// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant_settings

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"reflect"

	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/constants"
)

func NewTenantSettingsClient(api *api.ApiClient) TenantSettingsClient {
	return TenantSettingsClient{
		Api: api,
	}
}

type TenantSettingsClient struct {
	Api *api.ApiClient
}

func (client *TenantSettingsClient) GetTenant(ctx context.Context) (*TenantDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/tenant",
	}

	values := url.Values{}
	values.Add("api-version", "2020-08-01")
	apiUrl.RawQuery = values.Encode()

	tenant := TenantDto{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &tenant)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (client *TenantSettingsClient) GetTenantSettings(ctx context.Context) (*TenantSettingsDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/listTenantSettings",
	}

	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	tenantSettings := TenantSettingsDto{}
	_, err := client.Api.Execute(ctx, "POST", apiUrl.String(), nil, nil, []int{http.StatusOK}, &tenantSettings)
	if err != nil {
		return nil, err
	}
	return &tenantSettings, nil
}

func (client *TenantSettingsClient) UpdateTenantSettings(ctx context.Context, tenantSettings TenantSettingsDto) (*TenantSettingsDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/scopes/admin/updateTenantSettings",
	}

	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	var backendSettings TenantSettingsDto
	_, err := client.Api.Execute(ctx, "POST", apiUrl.String(), nil, tenantSettings, []int{http.StatusOK}, &backendSettings)
	if err != nil {
		return nil, err
	}
	return &backendSettings, nil
}

func applyCorrections(planned TenantSettingsDto, actual TenantSettingsDto) *TenantSettingsDto {
	corrected := filterDto(planned, actual).(*TenantSettingsDto)
	if planned.PowerPlatform != nil && planned.PowerPlatform.Governance != nil {
		if planned.PowerPlatform.Governance.EnvironmentRoutingTargetSecurityGroupId != nil && *planned.PowerPlatform.Governance.EnvironmentRoutingTargetSecurityGroupId == constants.ZERO_UUID && corrected.PowerPlatform.Governance.EnvironmentRoutingTargetSecurityGroupId == nil {
			zu := constants.ZERO_UUID
			corrected.PowerPlatform.Governance.EnvironmentRoutingTargetSecurityGroupId = &zu
		}

		if planned.PowerPlatform.Governance.EnvironmentRoutingTargetEnvironmentGroupId != nil && *planned.PowerPlatform.Governance.EnvironmentRoutingTargetEnvironmentGroupId == constants.ZERO_UUID && corrected.PowerPlatform.Governance.EnvironmentRoutingTargetEnvironmentGroupId == nil {
			zu := constants.ZERO_UUID
			corrected.PowerPlatform.Governance.EnvironmentRoutingTargetEnvironmentGroupId = &zu
		}
	}
	return corrected
}

// This function is used to filter out the fields that are not opted in to configuration
// The backend always returns all properties, but Terraform can only handle the properties that are opted in
func filterDto(configuredSettings interface{}, backendSettings interface{}) interface{} {
	configuredType := reflect.TypeOf(configuredSettings)
	backendType := reflect.TypeOf(backendSettings)
	if configuredType != backendType {
		return nil
	}

	output := reflect.New(configuredType).Interface()

	visibleFields := reflect.VisibleFields(configuredType)

	configuredValue := reflect.ValueOf(configuredSettings)
	backendValue := reflect.ValueOf(backendSettings)

	for fieldIndex, fieldInfo := range visibleFields {
		log.Default().Printf("Field: %s", fieldInfo.Name)

		configuredFieldValue := configuredValue.Field(fieldIndex)
		backendFieldValue := backendValue.Field(fieldIndex)
		outputField := reflect.ValueOf(output).Elem().Field(fieldIndex)

		if !configuredFieldValue.IsNil() && !backendFieldValue.IsNil() && backendFieldValue.IsValid() && outputField.CanSet() {
			if configuredFieldValue.Kind() == reflect.Pointer && configuredFieldValue.Elem().Kind() == reflect.Struct {
				outputStruct := filterDto(configuredFieldValue.Elem().Interface(), backendFieldValue.Elem().Interface())
				outputField.Set(reflect.ValueOf(outputStruct))
			} else if configuredFieldValue.Kind() == reflect.Pointer && configuredFieldValue.Elem().Kind() == reflect.Bool {
				boolValue := backendFieldValue.Elem().Bool()
				newBool := bool(boolValue)
				outputField.Set(reflect.ValueOf(&newBool))
			} else if configuredFieldValue.Kind() == reflect.Pointer && configuredFieldValue.Elem().Kind() == reflect.String {
				stringValue := backendFieldValue.Elem().String()
				newString := string(stringValue)
				outputField.Set(reflect.ValueOf(&newString))
			} else if configuredFieldValue.Kind() == reflect.Pointer && configuredFieldValue.Elem().Kind() == reflect.Int64 {
				int64Value := backendFieldValue.Elem().Int()
				newInt64 := int64(int64Value)
				outputField.Set(reflect.ValueOf(&newInt64))
			} else {
				log.Default().Printf("Skipping unknown field type %s", configuredFieldValue.Kind())
			}
		}
	}

	return output
}
