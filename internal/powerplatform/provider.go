// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

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
	constants "github.com/microsoft/terraform-provider-power-platform/constants"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	config "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/config"
	application "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/application"
	auth "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/authorization"
	connectors "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/connectors"
	dlp_policy "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/dlp_policy"
	environment "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/environment"
	licensing "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/licensing"
	managed_environment "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/managed_environment"
	powerapps "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/powerapps"
	solution "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/solution"
	tenant_settings "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/tenant_settings"
)

var _ provider.Provider = &PowerPlatformProvider{}

type PowerPlatformProvider struct {
	Config *config.ProviderConfig
	Api    *api.ApiClient
}

func NewPowerPlatformProvider(ctx context.Context, testModeEnabled ...bool) func() provider.Provider {
	cred := config.ProviderCredentials{}
	config := config.ProviderConfig{
		Credentials: &cred,
		Urls: config.ProviderConfigUrls{
			BapiUrl:          constants.BAPI_DOMAIN,
			PowerAppsUrl:     constants.POWERAPPS_API_DOMAIN,
			PowerPlatformUrl: constants.POWERPLATFORM_API_DOMAIN,
		},
	}

	if len(testModeEnabled) > 0 && testModeEnabled[0] {
		tflog.Warn(ctx, "Test mode enabled. Authentication requests will not be sent to the backend APIs.")
		config.Credentials.TestMode = true
	}

	return func() provider.Provider {
		p := &PowerPlatformProvider{
			Config: &config,
			Api:    api.NewApiClientBase(&config, api.NewAuthBase(&config)),
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
			"use_cli": schema.BoolAttribute{
				Description:         "Flag to indicate whether to use the CLI for authentication",
				MarkdownDescription: "Flag to indicate whether to use the CLI for authentication. ",
				Optional:            true,
			},
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
			"client_secret": schema.StringAttribute{
				Description:         "The secret of the Power Platform API app registration",
				MarkdownDescription: "The secret of the Power Platform API app registration",
				Optional:            true,
				Sensitive:           true,
			},
			"use_oidc": schema.BoolAttribute{
				Description:         "Allow OpenID Connect to be used for authentication",
				MarkdownDescription: "Allow OpenID Connect to be used for authentication",
				Optional:            true,
			},
			"oidc_request_token": schema.StringAttribute{
				Description: "The bearer token for the request to the OIDC provider. For use When authenticating as a Service Principal using OpenID Connect.",
				Optional:    true,
			},
			"oidc_request_url": schema.StringAttribute{
				Description: "The URL for the OIDC provider from which to request an ID token. For use When authenticating as a Service Principal using OpenID Connect.",
				Optional:    true,
			},
			"oidc_token": schema.StringAttribute{
				Description: "The OIDC ID token for use when authenticating as a Service Principal using OpenID Connect.",
				Optional:    true,
			},
			"oidc_token_file_path": schema.StringAttribute{
				Description: "The path to a file containing an OIDC ID token for use when authenticating as a Service Principal using OpenID Connect.",
				Optional:    true,
			},
			"telemetry_optout": schema.BoolAttribute{
				Description:         "Flag to indicate whether to opt out of telemetry. Default is `false`",
				MarkdownDescription: "Flag to indicate whether to opt out of telemetry. Default is `false`",
				Optional:            true,
			},
		},
	}
}

func (p *PowerPlatformProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config config.ProviderCredentialsModel

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

	clientId := ""
	envClientId := os.Getenv("POWER_PLATFORM_CLIENT_ID")
	if config.ClientId.IsNull() {
		clientId = envClientId
	} else {
		clientId = config.ClientId.ValueString()
	}

	clientSecret := ""
	envSecret := os.Getenv("POWER_PLATFORM_CLIENT_SECRET")
	if config.ClientSecret.IsNull() {
		clientSecret = envSecret
	} else {
		clientSecret = config.ClientSecret.ValueString()
	}

	//Check for AzDO and GitHub environment variables
	oidcRequestUrl := ""
	envOidcRequestUrl := MultiEnvDefaultFunc([]string{"ARM_OIDC_REQUEST_URL", "ACTIONS_ID_TOKEN_REQUEST_URL"})
	if config.OidcRequestUrl.IsNull() {
		tflog.Debug(ctx, "OIDC request URL environment variable is null")
		oidcRequestUrl = envOidcRequestUrl
	} else {
		oidcRequestUrl = config.OidcRequestUrl.ValueString()
	}

	oidcRequestToken := ""
	envOidcRequestToken := MultiEnvDefaultFunc([]string{"ARM_OIDC_REQUEST_TOKEN", "ACTIONS_ID_TOKEN_REQUEST_TOKEN"})
	if config.OidcRequestToken.IsNull() {
		tflog.Debug(ctx, "OIDC request token environment variable is null")
		oidcRequestToken = envOidcRequestToken
	} else {
		oidcRequestToken = config.OidcRequestToken.ValueString()
	}

	oidcToken := ""
	envOidcToken := EnvDefaultFunc("ARM_OIDC_TOKEN", "")
	if config.OidcToken.IsNull() {
		tflog.Debug(ctx, "OIDC token environment variable is null")
		oidcToken = envOidcToken
	} else {
		oidcToken = config.OidcToken.ValueString()
	}

	oidcTokenFilePath := ""
	envOidcTokenFilePath := EnvDefaultFunc("ARM_OIDC_TOKEN_FILE_PATH", "")
	if config.OidcTokenFilePath.IsNull() {
		oidcTokenFilePath = envOidcTokenFilePath
	} else {
		oidcTokenFilePath = config.OidcTokenFilePath.ValueString()
	}

	ctx = tflog.SetField(ctx, "use_cli", config.UseCli.ValueBool())
	ctx = tflog.SetField(ctx, "use_oidc", config.UseOidc.ValueBool())
	ctx = tflog.SetField(ctx, "power_platform_tenant_id", tenantId)
	ctx = tflog.SetField(ctx, "power_platform_client_id", clientId)
	ctx = tflog.SetField(ctx, "power_platform_client_secret", clientSecret)
	ctx = tflog.SetField(ctx, "oidc_request_url", oidcRequestUrl)
	ctx = tflog.SetField(ctx, "oidc_request_token", oidcRequestToken)
	ctx = tflog.SetField(ctx, "oidc_token", oidcToken)
	ctx = tflog.SetField(ctx, "oidc_token_file_path", oidcTokenFilePath)
	ctx = tflog.SetField(ctx, "telemetry_optout", config.TelemetryOptout.ValueBool())
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "power_platform_client_secret")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "oidc_request_token")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "oidc_token")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "oidc_token_file_path")

	if config.UseCli.ValueBool() {
		p.Config.Credentials.UseCli = true
	} else if config.UseOidc.ValueBool() {
		p.Config.Credentials.UseOidc = true
		p.Config.Credentials.TenantId = tenantId
		p.Config.Credentials.ClientId = clientId
		p.Config.Credentials.OidcRequestToken = oidcRequestToken
		p.Config.Credentials.OidcRequestUrl = oidcRequestUrl
		p.Config.Credentials.OidcToken = oidcToken
		p.Config.Credentials.OidcTokenFilePath = oidcTokenFilePath

	} else {

		if clientId != "" && clientSecret != "" && tenantId != "" {
			p.Config.Credentials.TenantId = tenantId
			p.Config.Credentials.ClientId = clientId
			p.Config.Credentials.ClientSecret = clientSecret
		} else {
			if tenantId == "" {
				resp.Diagnostics.AddAttributeError(
					path.Root("tenant_id"),
					"Unknown API tenant id",
					"The provider cannot create the API client as there is an unknown configuration value for the tenant id. "+
						"Either target apply the source of the value first, set the value statically in the configuration, or use the POWER_PLATFORM_TENANT_ID environment variable.",
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
			if clientSecret == "" {
				resp.Diagnostics.AddAttributeError(
					path.Root("client_secret"),
					"Unknown client secret",
					"The provider cannot create the API client as there is an unknown configuration value for the client secret. "+
						"Either target apply the source of the value first, set the value statically in the configuration, or use the POWER_PLATFORM_CLIENT_SECRET environment variable.",
				)
			}
		}
	}
	p.Config.Credentials.TelemetryOptout = config.TelemetryOptout.ValueBool()

	providerClient := api.ProviderClient{
		Config: p.Config,
		Api:    p.Api,
	}
	resp.DataSourceData = &providerClient
	resp.ResourceData = &providerClient

	tflog.Info(ctx, "Configured API client", map[string]any{"success": true})
}

func (p *PowerPlatformProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return environment.NewEnvironmentResource() },
		func() resource.Resource { return application.NewEnvironmentApplicationPackageInstallResource() },
		func() resource.Resource { return dlp_policy.NewDataLossPreventionPolicyResource() },
		func() resource.Resource { return solution.NewSolutionResource() },
		func() resource.Resource { return tenant_settings.NewTenantSettingsResource() },
		func() resource.Resource { return managed_environment.NewManagedEnvironmentResource() },
		func() resource.Resource { return licensing.NewBillingPolicyEnvironmentResource() },
		func() resource.Resource { return licensing.NewBillingPolicyResource() },
		func() resource.Resource { return auth.NewUserResource() },
	}
}

func (p *PowerPlatformProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return connectors.NewConnectorsDataSource() },
		func() datasource.DataSource { return application.NewEnvironmentApplicationPackagesDataSource() },
		func() datasource.DataSource { return powerapps.NewEnvironmentPowerAppsDataSource() },
		func() datasource.DataSource { return environment.NewEnvironmentsDataSource() },
		func() datasource.DataSource { return solution.NewSolutionsDataSource() },
		func() datasource.DataSource { return dlp_policy.NewDataLossPreventionPolicyDataSource() },
		func() datasource.DataSource { return tenant_settings.NewTenantSettingsDataSource() },
		func() datasource.DataSource { return licensing.NewBillingPoliciesDataSource() },
		func() datasource.DataSource { return licensing.NewBillingPoliciesEnvironmetsDataSource() },
		func() datasource.DataSource { return auth.NewSecurityRolesDataSource() },
	}
}

//TODO figure out how to return these defaultfuncs to their former interface-based glory

// MultiEnvDefaultFunc is a helper function that returns the value of the first
// environment variable in the given list that returns a non-empty value. If
// none of the environment variables return a value, the default value is
// returned.
func MultiEnvDefaultFunc(ks []string) string {

	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

// EnvDefaultFunc is a helper function that returns the value of the
// given environment variable, if one exists, or the default value
// otherwise.
func EnvDefaultFunc(k string, dv interface{}) string {

	if v := os.Getenv(k); v != "" {
		return v
	}
	return ""
}
