// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform_config

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProviderConfig struct {
	Credentials *ProviderCredentials
	Urls        ProviderConfigUrls
}

type ProviderConfigUrls struct {
	BapiUrl          string
	PowerAppsUrl     string
	PowerPlatformUrl string
}

type ProviderCredentials struct {
	TestMode bool
	UseCli   bool
	UseOidc  bool

	TelemetryOptout bool

	TenantId     string
	ClientId     string
	ClientSecret string

	OidcRequestToken  string
	OidcRequestUrl    string
	OidcToken         string
	OidcTokenFilePath string
}

type ProviderCredentialsModel struct {
	UseCli  types.Bool `tfsdk:"use_cli"`
	UseOidc types.Bool `tfsdk:"use_oidc"`

	TelemetryOptout types.Bool `tfsdk:"telemetry_optout"`

	TenantId     types.String `tfsdk:"tenant_id"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`

	OidcRequestToken  types.String `tfsdk:"oidc_request_token"`
	OidcRequestUrl    types.String `tfsdk:"oidc_request_url"`
	OidcToken         types.String `tfsdk:"oidc_token"`
	OidcTokenFilePath types.String `tfsdk:"oidc_token_file_path"`
}

func (model *ProviderCredentials) IsTelemetryOprout() bool {
	return model.TelemetryOptout
}

func (model *ProviderCredentials) IsClientSecretCredentialsProvided() bool {
	return model.ClientId != "" && model.ClientSecret != "" && model.TenantId != ""
}

func (model *ProviderCredentials) IsCliProvided() bool {
	return model.UseCli
}

func (model *ProviderCredentials) IsOidcProvided() bool {
	return model.UseOidc
}
