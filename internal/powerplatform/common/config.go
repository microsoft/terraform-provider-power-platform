package powerplatform_common

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
	TenantId string
	ClientId string
	Secret   string

	Username string
	Password string
}

type ProviderCredentialsModel struct {
	TenantId types.String `tfsdk:"tenant_id"`
	ClientId types.String `tfsdk:"client_id"`
	Secret   types.String `tfsdk:"secret"`

	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (model *ProviderCredentials) IsUserPassCredentialsProvided() bool {
	return model.Username != "" || model.Password != "" || model.TenantId != ""
}

func (model *ProviderCredentials) IsClientSecretCredentialsProvided() bool {
	return model.ClientId != "" || model.Secret != "" || model.TenantId != ""
}
