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

	TenantId     string
	ClientId     string
	ClientSecret string
}

type ProviderCredentialsModel struct {
	UseCli types.Bool `tfsdk:"use_cli"`

	TenantId     types.String `tfsdk:"tenant_id"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

func (model *ProviderCredentials) IsClientSecretCredentialsProvided() bool {
	return model.ClientId != "" && model.ClientSecret != "" && model.TenantId != ""
}

func (model *ProviderCredentials) IsCliProvided() bool {
	return model.UseCli
}
