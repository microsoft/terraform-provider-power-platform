// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package config

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/customtypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type CloudType string

const (
	CloudTypePublic  CloudType = "public"
	CloudTypeGcc     CloudType = "gcc"
	CloudTypeGccHigh CloudType = "gcchigh"
	CloudTypeDod     CloudType = "dod"
	CloudTypeChina   CloudType = "china"
	CloudTypeEx      CloudType = "ex"
	CloudTypeRx      CloudType = "rx"
)

type CloudTypeConfigurationKey string

const (
	FirstReleaseClusterName CloudTypeConfigurationKey = "release_cycle"
)

type ProviderConfig struct {
	UseCli  bool
	UseOidc bool
	UseMsi  bool

	TenantId           string
	AuxiliaryTenantIDs []string
	ClientId           string
	ClientSecret       string

	ClientCertificatePassword string
	ClientCertificateRaw      string

	OidcRequestToken  string
	OidcRequestUrl    string
	OidcToken         string
	OidcTokenFilePath string

	CloudType               CloudType
	AzDOServiceConnectionID string

	// CAE-related configuration
	EnableContinuousAccessEvaluation bool

	// internal runtime configuration values
	TestMode                  bool
	Urls                      ProviderConfigUrls
	TelemetryOptout           bool
	PartnerId                 string
	DisableTerraformPartnerId bool
	Cloud                     cloud.Configuration
	TerraformVersion          string
}

type ProviderConfigUrls struct {
	AdminPowerPlatformUrl string
	BapiUrl               string
	PowerAppsUrl          string
	PowerAppsScope        string
	PowerPlatformUrl      string
	PowerPlatformScope    string
	LicensingUrl          string
	PowerAppsAdvisor      string
	PowerAppsAdvisorScope string
	AnalyticsScope        string
}

func (model *ProviderConfig) GetCurrentCloudConfiguration(key CloudTypeConfigurationKey) *string {
	configuration := map[string]map[string]*string{
		string(CloudTypePublic): {
			string(FirstReleaseClusterName): helpers.StringPtr("FirstRelease"),
			// Add more cloud specific configurations here
		},
		string(CloudTypeGcc): {
			string(FirstReleaseClusterName): helpers.StringPtr("GovFR"),
		},
		string(CloudTypeGccHigh): {
			string(FirstReleaseClusterName): nil,
		},
		string(CloudTypeDod): {
			string(FirstReleaseClusterName): nil,
		},
		string(CloudTypeChina): {
			string(FirstReleaseClusterName): nil,
		},
		string(CloudTypeEx): {
			string(FirstReleaseClusterName): nil,
		},
		string(CloudTypeRx): {
			string(FirstReleaseClusterName): nil,
		},
	}

	return configuration[string(model.CloudType)][string(key)]
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

	Cloud                     types.String     `tfsdk:"cloud"`
	TelemetryOptout           types.Bool       `tfsdk:"telemetry_optout"`
	PartnerId                 customtypes.UUID `tfsdk:"partner_id"`
	DisableTerraformPartnerId types.Bool       `tfsdk:"disable_terraform_partner_id"`

	TenantId           types.String `tfsdk:"tenant_id"`
	AuxiliaryTenantIDs types.List   `tfsdk:"auxiliary_tenant_ids"`
	ClientId           types.String `tfsdk:"client_id"`
	ClientSecret       types.String `tfsdk:"client_secret"`

	ClientCertificateFilePath types.String `tfsdk:"client_certificate_file_path"`
	ClientCertificate         types.String `tfsdk:"client_certificate"`
	ClientCertificatePassword types.String `tfsdk:"client_certificate_password"`

	OidcRequestToken  types.String `tfsdk:"oidc_request_token"`
	OidcRequestUrl    types.String `tfsdk:"oidc_request_url"`
	OidcToken         types.String `tfsdk:"oidc_token"`
	OidcTokenFilePath types.String `tfsdk:"oidc_token_file_path"`

	AzDOServiceConnectionID types.String `tfsdk:"azdo_service_connection_id"`

	// CAE-related configuration
	EnableContinuousAccessEvaluation types.Bool `tfsdk:"enable_continuous_access_evaluation"`
}
