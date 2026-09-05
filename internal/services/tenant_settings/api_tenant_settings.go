// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant_settings

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

func newTenantSettingsClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}

type client struct {
	Api *api.Client
}

func (client *client) GetTenant(ctx context.Context) (*tenantDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/tenant",
	}

	values := url.Values{}
	values.Add(constants.API_VERSION_PARAM, constants.TENANT_SETTINGS_API_VERSION)
	apiUrl.RawQuery = values.Encode()

	tenant := tenantDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &tenant)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (client *client) GetTenantSettings(ctx context.Context) (*tenantSettingsDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/listTenantSettings",
	}

	values := url.Values{}
	values.Add(constants.API_VERSION_PARAM, constants.BAP_API_VERSION)
	apiUrl.RawQuery = values.Encode()

	tenantSettings := tenantSettingsDto{}
	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, nil, []int{http.StatusOK}, &tenantSettings)
	if err != nil {
		return nil, err
	}
	return &tenantSettings, nil
}

func (client *client) UpdateTenantSettings(ctx context.Context, tenantSettings tenantSettingsDto) (*tenantSettingsDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/scopes/admin/updateTenantSettings",
	}

	values := url.Values{}
	values.Add(constants.API_VERSION_PARAM, constants.BAP_API_VERSION)
	apiUrl.RawQuery = values.Encode()

	var backendSettings tenantSettingsDto
	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, tenantSettings, []int{http.StatusOK}, &backendSettings)
	if err != nil {
		return nil, err
	}
	return &backendSettings, nil
}

func applyCorrections(ctx context.Context, planned tenantSettingsDto, actual tenantSettingsDto) (*tenantSettingsDto, error) {
	correctedFilter := filterDto(ctx, planned, actual)
	corrected, ok := correctedFilter.(*tenantSettingsDto)
	if !ok {
		tflog.Error(ctx, "Type assertion failed in applyCorrections")
		return nil, errors.New("type assertion to *tenantSettingsDto failed in applyCorrections")
	}

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
	if planned.PowerPlatform != nil && planned.PowerPlatform.Intelligence != nil && corrected.PowerPlatform != nil && corrected.PowerPlatform.Intelligence != nil {
		if planned.PowerPlatform.Intelligence.CopilotStudioAuthorsSecurityGroupId != nil && *planned.PowerPlatform.Intelligence.CopilotStudioAuthorsSecurityGroupId == constants.ZERO_UUID && corrected.PowerPlatform.Intelligence.CopilotStudioAuthorsSecurityGroupId == nil {
			zu := constants.ZERO_UUID
			corrected.PowerPlatform.Intelligence.CopilotStudioAuthorsSecurityGroupId = &zu
		}
	}

	return corrected, nil
}

// This function is used to filter out the fields that are not opted in to configuration.
// The backend always returns all properties, but Terraform can only handle the properties that are opted in.
func filterDto(ctx context.Context, configuredSettings any, backendSettings any) any {
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
		tflog.Debug(ctx, fmt.Sprintf("Field: %s", fieldInfo.Name))

		configuredFieldValue := configuredValue.Field(fieldIndex)
		backendFieldValue := backendValue.Field(fieldIndex)
		outputField := reflect.ValueOf(output).Elem().Field(fieldIndex)

		if !configuredFieldValue.IsNil() && !backendFieldValue.IsNil() && backendFieldValue.IsValid() && outputField.CanSet() {
			if configuredFieldValue.Kind() == reflect.Pointer && configuredFieldValue.Elem().Kind() == reflect.Struct {
				outputStruct := filterDto(ctx, configuredFieldValue.Elem().Interface(), backendFieldValue.Elem().Interface())
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
				tflog.Debug(ctx, fmt.Sprintf("Skipping unknown field type %s", configuredFieldValue.Kind()))
			}
		}
	}

	return output
}
