package powerplatform

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"

	clients "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/clients"
	common "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/common"
	connectors "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/connectors"
	dlp_policy "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/dlp_policy"
	environment "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/environment"
	powerapps "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/powerapps"
	solution "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/solution"
	tenant_settings "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/tenant_settings"
)

var _ provider.Provider = &PowerPlatformProvider{}

type PowerPlatformProvider struct {
	Config           *common.ProviderConfig
	BapiApi          *clients.BapiClient
	DataverseApi     *clients.DataverseClient
	PowerPlatformApi *clients.PowerPlatoformApiClient
}

func NewPowerPlatformProvider() func() provider.Provider {
	return func() provider.Provider {

		cred := common.ProviderCredentials{}
		config := common.ProviderConfig{
			Credentials: &cred,
			Urls: common.ProviderConfigUrls{
				BapiUrl:          "api.bap.microsoft.com",
				PowerAppsUrl:     "api.powerapps.com",
				PowerPlatformUrl: "api.powerplatform.com",
			},
		}

		baseAuthBapi := api.NewAuthBase(&config)
		bapiAuth := api.NewBapiAuth(baseAuthBapi)
		baseApiForBapi := api.NewApiClientBase(&config, baseAuthBapi)
		bapiClientApi := api.NewBapiClientApi(baseApiForBapi, bapiAuth, nil)
		bapiClient := clients.NewBapiClient(bapiAuth, bapiClientApi)

		baseAuthPowerPlatform := api.NewAuthBase(&config)
		powerplatformAuth := api.NewPowerPlatformAuth(baseAuthPowerPlatform)
		baseApiForPpApi := api.NewApiClientBase(&config, baseAuthPowerPlatform)
		powerplatformClientApi := api.NewPowerPlatformClientApi(baseApiForPpApi, powerplatformAuth)
		powerplatformClient := clients.NewPowerPlatformApiClient(powerplatformAuth, powerplatformClientApi)

		baseAuthDataverse := api.NewAuthBase(&config)
		dataverseAuth := api.NewDataverseAuth(baseAuthDataverse)
		baseApiForDataverse := api.NewApiClientBase(&config, baseAuthDataverse)
		dataverseClientApi := api.NewDataverseClientApi(baseApiForDataverse, dataverseAuth)
		dataverseClient := clients.NewDataverseClient(dataverseAuth, dataverseClientApi)

		bapiClient.Client.SetDataverseClient(dataverseClient.Client)
		dataverseClient.Client.SetBapiClient(bapiClient.Client)

		p := &PowerPlatformProvider{
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
	var config common.ProviderCredentialsModel

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
		p.Config.Credentials.TenantId = tenantId
		p.Config.Credentials.ClientId = clientId
		p.Config.Credentials.Secret = secret
	} else if username != "" && password != "" && tenantId != "" {
		p.Config.Credentials.TenantId = tenantId
		p.Config.Credentials.Username = username
		p.Config.Credentials.Password = password
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

	providerClient := clients.ProviderClient{
		Config:           p.Config,
		BapiApi:          p.BapiApi,
		DataverseApi:     p.DataverseApi,
		PowerPlatformApi: p.PowerPlatformApi,
	}
	resp.DataSourceData = &providerClient
	resp.ResourceData = &providerClient

	tflog.Info(ctx, "Configured API client", map[string]any{"success": true})
}

func (p *PowerPlatformProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return environment.NewEnvironmentResource() },
		func() resource.Resource { return dlp_policy.NewDataLossPreventionPolicyResource() },
		func() resource.Resource { return solution.NewSolutionResource() },
		func() resource.Resource { return tenant_settings.NewTenantSettingsResource() },
	}
}

func (p *PowerPlatformProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return connectors.NewConnectorsDataSource() },
		func() datasource.DataSource { return powerapps.NewPowerAppsDataSource() },
		func() datasource.DataSource { return environment.NewEnvironmentsDataSource() },
		func() datasource.DataSource { return solution.NewSolutionsDataSource() },
		func() datasource.DataSource { return dlp_policy.NewDataLossPreventionPolicyDataSource() },
		func() datasource.DataSource { return tenant_settings.NewTenantSettingsDataSource() },
	}
}
