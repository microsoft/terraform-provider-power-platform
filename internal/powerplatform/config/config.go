// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform_config

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProviderConfig struct {
	Credentials     *ProviderCredentials
	Urls            ProviderConfigUrls
	TelemetryOptout bool
	Cloud           cloud.Configuration
}

type ProviderConfigUrls struct {
	BapiUrl            string
	PowerAppsUrl       string
	PowerAppsScope     string
	PowerPlatformUrl   string
	PowerPlatformScope string
}

type ProviderCredentials struct {
	TestMode bool
	UseCli   bool
	UseOidc  bool

	TenantId     string
	ClientId     string
	ClientSecret string

	ClientCertificatePassword string
	ClientCertificateRaw      string

	OidcRequestToken  string
	OidcRequestUrl    string
	OidcToken         string
	OidcTokenFilePath string
}

func (model *ProviderCredentials) IsClientSecretCredentialsProvided() bool {
	return model.ClientId != "" && model.ClientSecret != "" && model.TenantId != ""
}

func (model *ProviderCredentials) IsClientCertificateCredentialsProvided() bool {
	return model.ClientCertificateRaw != ""
}

func (model *ProviderCredentials) IsCliProvided() bool {
	return model.UseCli
}

func (model *ProviderCredentials) IsOidcProvided() bool {
	return model.UseOidc
}

type ProviderCredentialsModel struct {
	UseCli  types.Bool `tfsdk:"use_cli"`
	UseOidc types.Bool `tfsdk:"use_oidc"`

	Cloud           types.String `tfsdk:"cloud"`
	TelemetryOptout types.Bool   `tfsdk:"telemetry_optout"`

	TenantId     types.String `tfsdk:"tenant_id"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`

	ClientCertificateFilePath types.String `tfsdk:"client_certificate_file_path"`
	ClientCertificate         types.String `tfsdk:"client_certificate"`
	ClientCertificatePassword types.String `tfsdk:"client_certificate_password"`

	OidcRequestToken  types.String `tfsdk:"oidc_request_token"`
	OidcRequestUrl    types.String `tfsdk:"oidc_request_url"`
	OidcToken         types.String `tfsdk:"oidc_token"`
	OidcTokenFilePath types.String `tfsdk:"oidc_token_file_path"`
}
