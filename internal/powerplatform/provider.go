package powerplatform

import (
	"context"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	powerplatform_bapi "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi"
)

var _ provider.Provider = &PowerPlatformProvider{}
var _ powerplatform_bapi.ApiClientInterface = &powerplatform_bapi.ApiClient{}

type ProviderCredentialsModel struct {
	TenantId types.String `tfsdk:"tenant_id"`
	ClientId types.String `tfsdk:"client_id"`
	Secret   types.String `tfsdk:"secret"`

	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

type PowerPlatformProvider struct {
	old_bapiClient powerplatform_bapi.ApiClientInterface

	Config           *ProviderConfig
	BapiApi          *BapiClient
	DataverseApi     *DataverseClient
	PowerPlatformApi *PowerPlatoformApiClient
}

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

func (model *ProviderCredentials) IsUserPassCredentialsProvided() bool {
	return model.Username != "" || model.Password != "" || model.TenantId != ""
}

func (model *ProviderCredentials) IsClientSecretCredentialsProvided() bool {
	return model.ClientId != "" || model.Secret != "" || model.TenantId != ""
}

type BapiClient struct {
	Auth   BapiAuthInterface
	Client BapiClientInterface
}

type DataverseClient struct {
	Auth   DataverseAuthInterface
	Client DataverseClientInterface
}

type PowerPlatoformApiClient struct {
	Auth   PowerPlatformAuthInterface
	Client PowerPlatformClientInterface
}

func NewPowerPlatformProvider() func() provider.Provider {
	return func() provider.Provider {

		cred := ProviderCredentials{}
		config := ProviderConfig{
			Credentials: &cred,
			Urls: ProviderConfigUrls{
				BapiUrl:          "api.bap.microsoft.com",
				PowerAppsUrl:     "api.powerapps.com",
				PowerPlatformUrl: "api.powerplatform.com",
			},
		}

		//bapi
		baseAuthBapi := &AuthImplementation{
			Config: config,
		}
		bapiAuth := &BapiAuthImplementation{
			BaseAuth: baseAuthBapi,
		}
		baseApiForBapi := &ApiClientImplementation{
			Config:   config,
			BaseAuth: baseAuthBapi,
		}
		bapiClient := &BapiClient{
			Auth: bapiAuth,
			Client: &BapiClientImplementation{
				BaseApi: baseApiForBapi,
				Auth:    bapiAuth,
			},
		}
		bapiClient.Client.GetBase().SetAuth(bapiAuth)
		//

		//powerplatform
		baseAuthPowerPlatform := &AuthImplementation{
			Config: config,
		}
		powerplatformAuth := &PowerPlatformAuthImplementation{
			BaseAuth: baseAuthPowerPlatform,
		}

		baseApiForPpApi := &ApiClientImplementation{
			Config:   config,
			BaseAuth: baseAuthPowerPlatform,
		}
		powerplatformClient := &PowerPlatoformApiClient{
			Auth: powerplatformAuth,
			Client: &PowerPlatformClientImplementation{
				BaseApi: baseApiForPpApi,
				Auth:    powerplatformAuth,
			},
		}
		powerplatformClient.Client.GetBase().SetAuth(powerplatformAuth)
		//

		//dataverse
		baseAuthDataverse := &AuthImplementation{
			Config: config,
		}
		dataverseAuth := &DataverseAuthImplementation{
			BaseAuth: baseAuthDataverse,
		}
		baseApiForDataverse := &ApiClientImplementation{
			Config:   config,
			BaseAuth: baseAuthDataverse,
		}
		dataverseClient := &DataverseClient{
			Auth: dataverseAuth,
			Client: &DataverseClientImplementation{
				BaseApi:    baseApiForDataverse,
				Auth:       dataverseAuth,
				BapiClient: bapiClient.Client,
			},
		}
		//

		p := &PowerPlatformProvider{
			//todo to be removed
			old_bapiClient: &powerplatform_bapi.ApiClient{
				HttpClient:       http.DefaultClient,
				Provider:         &powerplatform_bapi.Provider{},
				DataverseAuthMap: make(map[string]*powerplatform_bapi.AuthResponse),
			},
			//

			Config:           &config,
			BapiApi:          bapiClient,
			DataverseApi:     dataverseClient,
			PowerPlatformApi: powerplatformClient,
		}
		return p
	}
}

func (p *PowerPlatformProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "powerplatform"
}

func (p *PowerPlatformProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {

	tflog.Debug(ctx, "Schema request received")

	resp.Schema = schema.Schema{

		Description:         "The Power Platform Terraform Provider allows managing environments and other resources within Power Platform",
		MarkdownDescription: "The Power Platform Provider allows managing environments and other resources within [Power Platform](https://powerplatform.microsoft.com/)",
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Description:         "The id of the AAD tenant that Power Platform API uses to authenticate with",
				MarkdownDescription: "The id of the AAD tenant that Power Platform API uses to authenticate with",
				Optional:            true,
			},
			"client_id": schema.StringAttribute{
				Description:         "The client id of the Power Platform API app registration",
				MarkdownDescription: "The client id of the Power Platform API app registration",
				Optional:            true,
			},
			"secret": schema.StringAttribute{
				Description:         "The secret of the Power Platform API app registration",
				MarkdownDescription: "The secret of the Power Platform API app registration",
				Optional:            true,
				Sensitive:           true,
			},

			"username": schema.StringAttribute{
				Description:         "The username of the Power Platform API in user@domain format",
				MarkdownDescription: "The username of the Power Platform API in user@domain format",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				Description:         "The password of the Power Platform API use",
				MarkdownDescription: "The password of the Power Platform API use",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *PowerPlatformProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config ProviderCredentialsModel

	tflog.Debug(ctx, "Configure request received")

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := ""
	envTenantId := os.Getenv("POWER_PLATFORM_TENANT_ID")
	if config.TenantId.IsNull() {
		tenantId = envTenantId
	} else {
		tenantId = config.TenantId.ValueString()
	}

	username := ""
	envUsername := os.Getenv("POWER_PLATFORM_USERNAME")
	if config.Username.IsNull() {
		username = envUsername
	} else {
		username = config.Username.ValueString()
	}

	password := ""
	envPassword := os.Getenv("POWER_PLATFORM_PASSWORD")
	if config.Password.IsNull() {
		password = envPassword
	} else {
		password = config.Password.ValueString()
	}

	clientId := ""
	envClientId := os.Getenv("POWER_PLATFORM_CLIENT_ID")
	if config.ClientId.IsNull() {
		clientId = envClientId
	} else {
		clientId = config.ClientId.ValueString()
	}

	secret := ""
	envSecret := os.Getenv("POWER_PLATFORM_SECRET")
	if config.Secret.IsNull() {
		secret = envSecret
	} else {
		secret = config.Secret.ValueString()
	}

	ctx = tflog.SetField(ctx, "power_platform_tenant_id", tenantId)
	ctx = tflog.SetField(ctx, "power_platform_username", username)
	ctx = tflog.SetField(ctx, "power_platform_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "power_platform_password")
	ctx = tflog.SetField(ctx, "power_platform_client_id", clientId)
	ctx = tflog.SetField(ctx, "power_platform_secret", secret)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "power_platform_secret")

	if clientId != "" && secret != "" && tenantId != "" {
		//todo refarctor remove
		authResp, err := p.old_bapiClient.DoAuthClientSecret(ctx, tenantId, clientId, secret)

		p.Config.Credentials.TenantId = tenantId
		p.Config.Credentials.ClientId = clientId
		p.Config.Credentials.Secret = secret

		if err != nil {
			resp.Diagnostics.AddError("Provider client's authentication has failed.", err.Error())
		} else {
			tflog.Info(ctx, "Authentication response token", map[string]any{"Token": authResp.Token})
		}

	} else if username != "" && password != "" && tenantId != "" {
		//todo refarctor remove
		authResp, err := p.old_bapiClient.DoAuthUsernamePassword(ctx, tenantId, username, password)

		p.Config.Credentials.TenantId = tenantId
		p.Config.Credentials.Username = username
		p.Config.Credentials.Password = password

		if err != nil {
			resp.Diagnostics.AddError("Provider client's authentication has failed.", err.Error())
		} else {
			tflog.Info(ctx, "Authentication response token", map[string]any{"Token": authResp.Token})
		}
	} else {
		if tenantId == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("tenant_id"),
				"Unknown API tenant id",
				"The provider cannot create the API client as there is an unknown configuration value for the tenant id. "+
					"Either target apply the source of the value first, set the value statically in the configuration, or use the POWER_PLATFORM_TENANT_ID environment variable.",
			)
		}
		if username == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("username"),
				"Unknown username",
				"The provider cannot create the API client as there is an unknown configuration value for the username. "+
					"Either target apply the source of the value first, set the value statically in the configuration, or use the POWER_PLATFORM_USERNAME environment variable.",
			)
		}
		if password == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("password"),
				"Unknown password",
				"The provider cannot create the API client as there is an unknown configuration value for the password. "+
					"Either target apply the source of the value first, set the value statically in the configuration, or use the POWER_PLATFORM_PASSWORD environment variable.",
			)
		}
		if clientId == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("client_id"),
				"Unknown client id",
				"The provider cannot create the API client as there is an unknown configuration value for the client id. "+
					"Either target apply the source of the value first, set the value statically in the configuration, or use the POWER_PLATFORM_CLIENT_ID environment variable.",
			)
		}
		if secret == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("secret"),
				"Unknown secret",
				"The provider cannot create the API client as there is an unknown configuration value for the secret. "+
					"Either target apply the source of the value first, set the value statically in the configuration, or use the POWER_PLATFORM_SECRET environment variable.",
			)
		}
	}

	resp.DataSourceData = p
	resp.ResourceData = p

	tflog.Info(ctx, "Configured API client", map[string]any{"success": true})
}

func (p *PowerPlatformProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return NewEnvironmentResource() },
		func() resource.Resource { return NewDataLossPreventionPolicyResource() },
		func() resource.Resource { return NewSolutionResource() },
	}
}

func (p *PowerPlatformProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return NewConnectorsDataSource() },
		func() datasource.DataSource { return NewPowerAppsDataSource() },
		func() datasource.DataSource { return NewEnvironmentsDataSource() },
		func() datasource.DataSource { return NewSolutionsDataSource() },
	}
}
