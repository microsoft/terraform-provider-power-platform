// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"

	azcloud "github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	constants "github.com/microsoft/terraform-provider-power-platform/constants"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	config "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/config"
	helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	application "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/application"
	auth "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/authorization"
	connection "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/connection"
	connections "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/connection"
	connectors "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/connectors"
	currencies "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/currencies"
	data_record "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/data_record"
	dlp_policy "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/dlp_policy"
	environment "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/environment"
	env_settings "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/environment_settings"
	environment_templates "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/environment_templates"
	languages "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/languages"
	licensing "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/licensing"
	locations "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/locations"
	managed_environment "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/managed_environment"
	powerapps "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/powerapps"
	rest "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/rest"
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
			BapiUrl:            constants.PUBLIC_BAPI_DOMAIN,
			PowerAppsUrl:       constants.PUBLIC_POWERAPPS_API_DOMAIN,
			PowerAppsScope:     constants.PUBLIC_POWERAPPS_SCOPE,
			PowerPlatformUrl:   constants.PUBLIC_POWERPLATFORM_API_DOMAIN,
			PowerPlatformScope: constants.PUBLIC_POWERPLATFORM_API_SCOPE,
		},
		Cloud: azcloud.AzurePublic,
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
			"client_certificate": schema.StringAttribute{
				MarkdownDescription: "Base64 encoded PKCS#12 certificate bundle. For use when authenticating as a Service Principal using a Client Certificate.",
				Optional:            true,
				Sensitive:           true,
			},
			"client_certificate_file_path": schema.StringAttribute{
				MarkdownDescription: "The path to the Client Certificate associated with the Service Principal for use when authenticating as a Service Principal using a Client Certificate.",
				Optional:            true,
			},
			"client_certificate_password": schema.StringAttribute{
				MarkdownDescription: "The password associated with the Client Certificate. For use when authenticating as a Service Principal using a Client Certificate.",
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
			"cloud": schema.StringAttribute{
				Description:         "The cloud to use for authentication and Power Platform API requests. Default is `public`. Valid values are `public`, `gcc`, `gcchigh`, `china`, `dod`, `ex`, `rx`",
				MarkdownDescription: "The cloud to use for authentication and Power Platform API requests. Default is `public`. Valid values are `public`, `gcc`, `gcchigh`, `china`, `dod`, `ex`, `rx`",
				Optional:            true,
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

	cloud := "public"
	envCloud := os.Getenv("POWER_PLATFORM_CLOUD")
	if config.Cloud.IsNull() && envCloud != "" {
		cloud = envCloud
	} else if !config.Cloud.IsNull() {
		cloud = config.Cloud.ValueString()
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

	useOidc := false
	_, envUseOidc := os.LookupEnv("POWER_PLATFORM_USE_OIDC")
	tflog.Debug(ctx, fmt.Sprintf("Use OIDC env value: %v", envUseOidc))
	if config.UseOidc.IsNull() {
		tflog.Debug(ctx, fmt.Sprintf("Use OIDC is null. Using env value: %v", envUseOidc))
		useOidc = envUseOidc
	} else {
		tflog.Debug(ctx, fmt.Sprintf("Use OIDC is not null. Using config value: %v", config.UseOidc.ValueBool()))
		useOidc = config.UseOidc.ValueBool()
	}

	useCli := false
	_, envUseCli := os.LookupEnv("POWER_PLATFORM_USE_CLI")
	if config.UseCli.IsNull() {
		useCli = envUseCli
	} else {
		useCli = config.UseCli.ValueBool()
	}

	clientCertificate := ""
	envClientCertificate := os.Getenv("POWER_PLATFORM_CLIENT_CERTIFICATE")
	if config.ClientCertificate.IsNull() {
		clientCertificate = envClientCertificate
	} else {
		clientCertificate = config.ClientCertificate.ValueString()
	}

	clientCertificateFilePath := ""
	envClientCertificateFilePath := os.Getenv("POWER_PLATFORM_CLIENT_CERTIFICATE_FILE_PATH")
	if config.ClientCertificateFilePath.IsNull() {
		clientCertificateFilePath = envClientCertificateFilePath
	} else {
		clientCertificateFilePath = config.ClientCertificateFilePath.ValueString()
	}

	clientCertificatePassword := ""
	envClientCertificatePassword := os.Getenv("POWER_PLATFORM_CLIENT_CERTIFICATE_PASSWORD")
	if config.ClientCertificatePassword.IsNull() {
		clientCertificatePassword = envClientCertificatePassword
	} else {
		clientCertificatePassword = config.ClientCertificatePassword.ValueString()
	}

	//Check for AzDO and GitHub environment variables
	oidcRequestUrl := ""
	envOidcRequestUrl := MultiEnvDefaultFunc([]string{"ARM_OIDC_REQUEST_URL", "ACTIONS_ID_TOKEN_REQUEST_URL"})
	if config.OidcRequestUrl.IsNull() {
		oidcRequestUrl = envOidcRequestUrl
	} else {
		oidcRequestUrl = config.OidcRequestUrl.ValueString()
	}

	oidcRequestToken := ""
	envOidcRequestToken := MultiEnvDefaultFunc([]string{"ARM_OIDC_REQUEST_TOKEN", "ACTIONS_ID_TOKEN_REQUEST_TOKEN"})
	if config.OidcRequestToken.IsNull() {
		oidcRequestToken = envOidcRequestToken
	} else {
		oidcRequestToken = config.OidcRequestToken.ValueString()
	}

	oidcToken := ""
	envOidcToken := EnvDefaultFunc("ARM_OIDC_TOKEN", "")
	if config.OidcToken.IsNull() {
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

	ctx = tflog.SetField(ctx, "telemetry_optout", strconv.FormatBool(config.TelemetryOptout.ValueBool())+"\n")
	ctx = tflog.SetField(ctx, "use_oidc", strconv.FormatBool(useOidc)+"\n")
	ctx = tflog.SetField(ctx, "use_cli", strconv.FormatBool(useCli)+"\n")
	ctx = tflog.SetField(ctx, "cloud", cloud+"\n")

	ctx = tflog.SetField(ctx, "power_platform_tenant_id", tenantId+"\n")
	ctx = tflog.SetField(ctx, "power_platform_client_id", clientId+"\n")
	ctx = tflog.SetField(ctx, "power_platform_client_secret", clientSecret+"\n")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "power_platform_client_secret\n")

	ctx = tflog.SetField(ctx, "client_certificate_file_path", clientCertificateFilePath+"\n")
	ctx = tflog.SetField(ctx, "client_certificate", clientCertificate+"\n")
	ctx = tflog.SetField(ctx, "client_certificate_password", clientCertificatePassword+"\n")
	ctx = tflog.MaskAllFieldValuesRegexes(ctx, regexp.MustCompile(`(?i)client_certificate`))

	ctx = tflog.SetField(ctx, "oidc_request_url", oidcRequestUrl+"\n")
	ctx = tflog.SetField(ctx, "oidc_request_token", oidcRequestToken+"\n")
	ctx = tflog.SetField(ctx, "oidc_token", oidcToken+"\n")
	ctx = tflog.SetField(ctx, "oidc_token_file_path", oidcTokenFilePath+"\n")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "oidc_request_token\n")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "oidc_token\n")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "oidc_token_file_path\n")

	if useCli {
		tflog.Info(ctx, "Using CLI for authentication")
		p.Config.Credentials.UseCli = true
	} else if useOidc {
		tflog.Info(ctx, "Using OpenID Connect for authentication")
		ValidateProviderAttribute(resp, path.Root("tenant_id"), "tenant id", tenantId, "POWER_PLATFORM_TENANT_ID")
		ValidateProviderAttribute(resp, path.Root("client_id"), "client id", clientId, "POWER_PLATFORM_CLIENT_ID")

		p.Config.Credentials.UseOidc = true
		p.Config.Credentials.TenantId = tenantId
		p.Config.Credentials.ClientId = clientId
		p.Config.Credentials.OidcRequestToken = oidcRequestToken
		p.Config.Credentials.OidcRequestUrl = oidcRequestUrl
		p.Config.Credentials.OidcToken = oidcToken
		p.Config.Credentials.OidcTokenFilePath = oidcTokenFilePath
	} else if clientCertificatePassword != "" && (clientCertificate != "" || clientCertificateFilePath != "") {
		tflog.Info(ctx, "Using client certificate for authentication")
		ValidateProviderAttribute(resp, path.Root("tenant_id"), "tenant id", tenantId, "POWER_PLATFORM_TENANT_ID")
		ValidateProviderAttribute(resp, path.Root("client_id"), "client id", clientId, "POWER_PLATFORM_CLIENT_ID")

		cert, err := helpers.GetCertificateRawFromCertOrFilePath(clientCertificate, clientCertificateFilePath)
		if err != nil {
			resp.Diagnostics.AddAttributeError(path.Root("client_certificate"), "Error getting certificate", err.Error())
		}
		p.Config.Credentials.ClientCertificateRaw = cert
		p.Config.Credentials.ClientCertificatePassword = clientCertificatePassword
		p.Config.Credentials.TenantId = tenantId
		p.Config.Credentials.ClientId = clientId
	} else {
		tflog.Info(ctx, "Using client id and secret for authentication")
		if tenantId != "" && clientId != "" && clientSecret != "" {
			p.Config.Credentials.TenantId = tenantId
			p.Config.Credentials.ClientId = clientId
			p.Config.Credentials.ClientSecret = clientSecret
		} else {
			ValidateProviderAttribute(resp, path.Root("tenant_id"), "tenant id", tenantId, "POWER_PLATFORM_TENANT_ID")
			ValidateProviderAttribute(resp, path.Root("client_id"), "client id", clientId, "POWER_PLATFORM_CLIENT_ID")
			ValidateProviderAttribute(resp, path.Root("client_secret"), "client secret", clientSecret, "POWER_PLATFORM_CLIENT_SECRET")
		}
	}

	switch cloud {
	case "public":
		p.Config.Urls.BapiUrl = constants.PUBLIC_BAPI_DOMAIN
		p.Config.Urls.PowerAppsUrl = constants.PUBLIC_POWERAPPS_API_DOMAIN
		p.Config.Urls.PowerAppsScope = constants.PUBLIC_POWERAPPS_SCOPE
		p.Config.Urls.PowerPlatformUrl = constants.PUBLIC_POWERPLATFORM_API_DOMAIN
		p.Config.Urls.PowerPlatformScope = constants.PUBLIC_POWERPLATFORM_API_SCOPE
		p.Config.Cloud = azcloud.AzurePublic
	case "gcc":
		p.Config.Urls.BapiUrl = constants.USGOV_BAPI_DOMAIN
		p.Config.Urls.PowerAppsUrl = constants.USGOV_POWERAPPS_API_DOMAIN
		p.Config.Urls.PowerAppsScope = constants.USGOV_POWERAPPS_SCOPE
		p.Config.Urls.PowerPlatformUrl = constants.USGOV_POWERPLATFORM_API_DOMAIN
		p.Config.Urls.PowerPlatformScope = constants.USGOV_POWERPLATFORM_API_SCOPE
		p.Config.Cloud = azcloud.AzurePublic //GCC uses public cloud for authentication
	case "gcchigh":
		p.Config.Urls.BapiUrl = constants.USGOVHIGH_BAPI_DOMAIN
		p.Config.Urls.PowerAppsUrl = constants.USGOVHIGH_POWERAPPS_API_DOMAIN
		p.Config.Urls.PowerAppsScope = constants.USGOVHIGH_POWERAPPS_SCOPE
		p.Config.Urls.PowerPlatformUrl = constants.USGOVHIGH_POWERPLATFORM_API_DOMAIN
		p.Config.Urls.PowerPlatformScope = constants.USGOVHIGH_POWERPLATFORM_API_SCOPE
		p.Config.Cloud = azcloud.AzureGovernment
	case "dod":
		p.Config.Urls.BapiUrl = constants.USDOD_BAPI_DOMAIN
		p.Config.Urls.PowerAppsUrl = constants.USDOD_POWERAPPS_API_DOMAIN
		p.Config.Urls.PowerAppsScope = constants.USDOD_POWERAPPS_SCOPE
		p.Config.Urls.PowerPlatformUrl = constants.USDOD_POWERPLATFORM_API_DOMAIN
		p.Config.Urls.PowerPlatformScope = constants.USDOD_POWERPLATFORM_API_SCOPE
		p.Config.Cloud = azcloud.AzureGovernment
	case "china":
		p.Config.Urls.BapiUrl = constants.CHINA_BAPI_DOMAIN
		p.Config.Urls.PowerAppsUrl = constants.CHINA_POWERAPPS_API_DOMAIN
		p.Config.Urls.PowerAppsScope = constants.CHINA_POWERAPPS_SCOPE
		p.Config.Urls.PowerPlatformUrl = constants.CHINA_POWERPLATFORM_API_DOMAIN
		p.Config.Urls.PowerPlatformScope = constants.CHINA_POWERPLATFORM_API_SCOPE
		p.Config.Cloud = azcloud.AzureChina
	case "ex":
		p.Config.Urls.BapiUrl = constants.EX_BAPI_DOMAIN
		p.Config.Urls.PowerAppsUrl = constants.EX_POWERAPPS_API_DOMAIN
		p.Config.Urls.PowerAppsScope = constants.EX_POWERAPPS_SCOPE
		p.Config.Urls.PowerPlatformUrl = constants.EX_POWERPLATFORM_API_DOMAIN
		p.Config.Urls.PowerPlatformScope = constants.EX_POWERPLATFORM_API_SCOPE
		p.Config.Cloud = azcloud.Configuration{
			ActiveDirectoryAuthorityHost: constants.EX_AUTHORITY_HOST,
			Services:                     map[azcloud.ServiceName]azcloud.ServiceConfiguration{},
		}
	case "rx":
		p.Config.Urls.BapiUrl = constants.RX_BAPI_DOMAIN
		p.Config.Urls.PowerAppsUrl = constants.RX_POWERAPPS_API_DOMAIN
		p.Config.Urls.PowerAppsScope = constants.RX_POWERAPPS_SCOPE
		p.Config.Urls.PowerPlatformUrl = constants.RX_POWERPLATFORM_API_DOMAIN
		p.Config.Urls.PowerPlatformScope = constants.RX_POWERPLATFORM_API_SCOPE
		p.Config.Cloud = azcloud.Configuration{
			ActiveDirectoryAuthorityHost: constants.RX_AUTHORITY_HOST,
			Services:                     map[azcloud.ServiceName]azcloud.ServiceConfiguration{},
		}
	default:
		resp.Diagnostics.AddAttributeError(
			path.Root("cloud"),
			"Unknown cloud",
			"The provider cannot create the API client as there is an unknown configuration value for `cloud`. "+
				"Either set the value in the provider configuration or use the POWER_PLATFORM_CLOUD environment variable.",
		)
	}

	p.Config.TelemetryOptout = config.TelemetryOptout.ValueBool()

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
		func() resource.Resource { return data_record.NewDataRecordResource() },
		func() resource.Resource { return env_settings.NewEnvironmentSettingsResource() },
		func() resource.Resource { return connection.NewConnectionResource() },
		func() resource.Resource { return rest.NewDataverseWebApiResource() },
		func() resource.Resource { return connections.NewConnectionShareResource() },
	}
}

func (p *PowerPlatformProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return connectors.NewConnectorsDataSource() },
		func() datasource.DataSource { return application.NewEnvironmentApplicationPackagesDataSource() },
		func() datasource.DataSource { return powerapps.NewEnvironmentPowerAppsDataSource() },
		func() datasource.DataSource { return environment.NewEnvironmentsDataSource() },
		func() datasource.DataSource { return environment_templates.NewEnvironmentTemplatesDataSource() },
		func() datasource.DataSource { return solution.NewSolutionsDataSource() },
		func() datasource.DataSource { return dlp_policy.NewDataLossPreventionPolicyDataSource() },
		func() datasource.DataSource { return tenant_settings.NewTenantSettingsDataSource() },
		func() datasource.DataSource { return licensing.NewBillingPoliciesDataSource() },
		func() datasource.DataSource { return licensing.NewBillingPoliciesEnvironmetsDataSource() },
		func() datasource.DataSource { return env_settings.NewEnvironmentSettingsDataSource() },
		func() datasource.DataSource { return locations.NewLocationsDataSource() },
		func() datasource.DataSource { return languages.NewLanguagesDataSource() },
		func() datasource.DataSource { return currencies.NewCurrenciesDataSource() },
		func() datasource.DataSource { return auth.NewSecurityRolesDataSource() },
		func() datasource.DataSource { return application.NewTenantApplicationPackagesDataSource() },
		func() datasource.DataSource { return data_record.NewDataRecordDataSource() },
		func() datasource.DataSource { return rest.NewDataverseWebApiDatasource() },
		func() datasource.DataSource { return connections.NewConnectionsDataSource() },
		func() datasource.DataSource { return connection.NewConnectionSharesDataSource() },
	}
}

func ValidateProviderAttribute(resp *provider.ConfigureResponse, path path.Path, name, value string, environmentVariableName string) {

	environmentVariableText := "Target apply the source of the value first, set the value statically in the configuration."
	if environmentVariableName != "" {
		environmentVariableText = fmt.Sprintf("Either target apply the source of the value first, set the value statically in the configuration, or use the %s environment variable.", environmentVariableName)
	}

	if value == "" {
		resp.Diagnostics.AddAttributeError(
			path,
			fmt.Sprintf("Unknown %s", name),
			fmt.Sprintf("The provider cannot create the API client as there is an unknown configuration value for %s. %s", name, environmentVariableText))
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
