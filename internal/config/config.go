// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package config

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProviderConfig struct {
	UseCli  bool
	UseOidc bool
	UseMsi  bool

	TenantId     string
	ClientId     string
	ClientSecret string

	ClientCertificatePassword string
	ClientCertificateRaw      string

	OidcRequestToken  string
	OidcRequestUrl    string
	OidcToken         string
	OidcTokenFilePath string

	AzDOServiceConnectionID string

	// internal runtime configuration values
	TestMode         bool
	Urls             ProviderConfigUrls
	TelemetryOptout  bool
	Cloud            cloud.Configuration
	TerraformVersion string
}

type ProviderConfigUrls struct {
	BapiUrl               string
	PowerAppsUrl          string
	PowerAppsScope        string
	PowerPlatformUrl      string
	PowerPlatformScope    string
	LicensingUrl          string
	PowerAppsAdvisor      string
	PowerAppsAdvisorScope string
}

func (model *ProviderConfig) IsUserManagedIdentityProvided() bool {
	return model.UseMsi && model.ClientId != ""
}

func (model *ProviderConfig) IsSystemManagedIdentityProvided() bool {
	return model.UseMsi && model.ClientId == "" // The switch that consumes this could be structured to avoid the second check, but we don't have a guarantee of what's consuming this.
}

func (model *ProviderConfig) IsAzDOWorkloadIdentityFederationProvided() bool {
	return model.UseOidc && model.AzDOServiceConnectionID != ""
}

func (model *ProviderConfig) IsClientSecretCredentialsProvided() bool {
	return model.ClientId != "" && model.ClientSecret != "" && model.TenantId != ""
}

func (model *ProviderConfig) IsClientCertificateCredentialsProvided() bool {
	return model.ClientCertificateRaw != ""
}

func (model *ProviderConfig) IsCliProvided() bool {
	return model.UseCli
}

func (model *ProviderConfig) IsOidcProvided() bool {
	return model.UseOidc
}

// ProviderConfigModel is a model for the provider configuration.
type ProviderConfigModel struct {
	UseCli  types.Bool `tfsdk:"use_cli"`
	UseOidc types.Bool `tfsdk:"use_oidc"`
	UseMsi  types.Bool `tfsdk:"use_msi"`

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

	AzDOServiceConnectionID types.String `tfsdk:"azdo_service_connection_id"`
}
