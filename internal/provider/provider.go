// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/common"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/config"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/admin_management_application"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/application"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/authorization"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/capacity"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/connection"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/connectors"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/currencies"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/data_record"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/dlp_policy"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/enterprise_policy"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment_group_rule_set"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment_groups"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment_settings"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment_templates"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/languages"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/licensing"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/locations"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/managed_environment"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/powerapps"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/rest"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/solution"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/tenant"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/tenant_settings"
)

var _ provider.Provider = &PowerPlatformProvider{}

type PowerPlatformProvider struct {
	Config *config.ProviderConfig
	Api    *api.Client
}

func NewPowerPlatformProvider(ctx context.Context, testModeEnabled ...bool) func() provider.Provider {
	cloudUrls, cloudConfig := getCloudPublicUrls()
	providerConfig := config.ProviderConfig{
		Urls:             *cloudUrls,
		Cloud:            *cloudConfig,
		TerraformVersion: "unknown",
		TelemetryOptout:  false,
	}

	if len(testModeEnabled) > 0 && testModeEnabled[0] {
		tflog.Warn(ctx, "Test mode enabled. Authentication requests will not be sent to the backend APIs.")
		providerConfig.TestMode = true
	}

	return func() provider.Provider {
		p := &PowerPlatformProvider{
			Config: &providerConfig,
			Api:    api.NewApiClientBase(&providerConfig, api.NewAuthBase(&providerConfig)),
		}
		return p
	}
}

func (p *PowerPlatformProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "powerplatform"
	resp.Version = common.ProviderVersion

	tflog.Debug(ctx, "Provider Metadata request received", map[string]any{
		"version":  resp.Version,
		"typeName": resp.TypeName,
		"branch":   common.Branch,
		"commit":   common.Commit,
	})
}

func (p *PowerPlatformProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	ctx, exitContext := helpers.EnterProviderContext(ctx, req)
	defer exitContext()

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
			"use_msi": schema.BoolAttribute{
				Description:         "Flag to indicate whether to use managed identity for authentication",
				MarkdownDescription: "Flag to indicate whether to use managed identity for authentication",
				Optional:            true,
			},
			"azdo_service_connection_id": schema.StringAttribute{
				Description:         "The service connection id of the Azure DevOps service connection. For use in workload identity federation.",
				MarkdownDescription: "The service connection id of the Azure DevOps service connection. For use in workload identity federation.",
				Optional:            true,
			},
		},
	}
}

func (p *PowerPlatformProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	ctx, exitContext := helpers.EnterProviderContext(ctx, req)
	defer exitContext()

	// Get Provider Configuration from the provider block in the configuration.
	var configValue config.ProviderConfigModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &configValue)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Provider Configuration from the configuration, environment variables, or defaults.
	cloudType := helpers.GetConfigString(ctx, configValue.Cloud, constants.ENV_VAR_POWER_PLATFORM_CLOUD, "public")
	tenantId := helpers.GetConfigString(ctx, configValue.TenantId, constants.ENV_VAR_POWER_PLATFORM_TENANT_ID, "")
	clientId := helpers.GetConfigString(ctx, configValue.ClientId, constants.ENV_VAR_POWER_PLATFORM_CLIENT_ID, "")
	clientSecret := helpers.GetConfigString(ctx, configValue.ClientSecret, constants.ENV_VAR_POWER_PLATFORM_CLIENT_SECRET, "")
	useOidc := helpers.GetConfigBool(ctx, configValue.UseOidc, constants.ENV_VAR_POWER_PLATFORM_USE_OIDC, false)
	useCli := helpers.GetConfigBool(ctx, configValue.UseCli, constants.ENV_VAR_POWER_PLATFORM_USE_CLI, false)
	clientCertificate := helpers.GetConfigString(ctx, configValue.ClientCertificate, constants.ENV_VAR_POWER_PLATFORM_CLIENT_CERTIFICATE, "")
	clientCertificateFilePath := helpers.GetConfigString(ctx, configValue.ClientCertificateFilePath, constants.ENV_VAR_POWER_PLATFORM_CLIENT_CERTIFICATE_FILE_PATH, "")
	clientCertificatePassword := helpers.GetConfigString(ctx, configValue.ClientCertificatePassword, constants.ENV_VAR_POWER_PLATFORM_CLIENT_CERTIFICATE_PASSWORD, "")
	useMsi := helpers.GetConfigBool(ctx, configValue.UseMsi, constants.ENV_VAR_POWER_PLATFORM_USE_MSI, false)
	azdoServiceConnectionId := helpers.GetConfigString(ctx, configValue.AzDOServiceConnectionID, constants.ENV_VAR_POWER_PLATFORM_AZDO_SERVICE_CONNECTION_ID, "")

	// Check for AzDO and GitHub environment variables
	oidcRequestUrl := helpers.GetConfigMultiString(ctx, configValue.OidcRequestUrl, []string{constants.ENV_VAR_ARM_OIDC_REQUEST_URL, constants.ENV_VAR_ACTIONS_ID_TOKEN_REQUEST_URL}, "")
	oidcRequestToken := helpers.GetConfigMultiString(ctx, configValue.OidcRequestToken, []string{constants.ENV_VAR_ARM_OIDC_REQUEST_TOKEN, constants.ENV_VAR_ACTIONS_ID_TOKEN_REQUEST_TOKEN}, "")
	oidcToken := helpers.GetConfigString(ctx, configValue.OidcToken, constants.ENV_VAR_ARM_OIDC_TOKEN, "")
	oidcTokenFilePath := helpers.GetConfigString(ctx, configValue.OidcTokenFilePath, constants.ENV_VAR_ARM_OIDC_TOKEN_FILE_PATH, "")

	// Check for telemetry opt out
	telemetryOptOut := helpers.GetConfigBool(ctx, configValue.TelemetryOptout, constants.ENV_VAR_POWER_PLATFORM_TELEMETRY_OPTOUT, false)

	// Set the configuration values

	if p.Config.TestMode {
		tflog.Info(ctx, "Test mode enabled. Authentication requests will not be sent to the backend APIs.")
	} else if useCli {
		tflog.Info(ctx, "Using CLI for authentication")
		p.Config.UseCli = true
	} else if useOidc && azdoServiceConnectionId != "" {
		tflog.Info(ctx, "Using Workload Identity Federation for Azure Pipelines")
		p.Config.UseOidc = true
		p.Config.AzDOServiceConnectionID = azdoServiceConnectionId
	} else if useOidc {
		tflog.Info(ctx, "Using OpenID Connect for authentication")
		validateProviderAttribute(resp, path.Root("tenant_id"), "tenant id", tenantId, constants.ENV_VAR_POWER_PLATFORM_TENANT_ID)
		validateProviderAttribute(resp, path.Root("client_id"), "client id", clientId, constants.ENV_VAR_POWER_PLATFORM_CLIENT_ID)

		p.Config.UseOidc = true
		p.Config.TenantId = tenantId
		p.Config.ClientId = clientId
		p.Config.OidcRequestToken = oidcRequestToken
		p.Config.OidcRequestUrl = oidcRequestUrl
		p.Config.OidcToken = oidcToken
		p.Config.OidcTokenFilePath = oidcTokenFilePath
	} else if useMsi {
		tflog.Info(ctx, "Using Managed Identity for authentication")
		p.Config.ClientId = clientId // No client ID as it could be blank for system-managed or populated for user-managed.
		p.Config.UseMsi = true
	} else if clientCertificatePassword != "" && (clientCertificate != "" || clientCertificateFilePath != "") {
		tflog.Info(ctx, "Using client certificate for authentication")
		validateProviderAttribute(resp, path.Root("tenant_id"), "tenant id", tenantId, constants.ENV_VAR_POWER_PLATFORM_TENANT_ID)
		validateProviderAttribute(resp, path.Root("client_id"), "client id", clientId, constants.ENV_VAR_POWER_PLATFORM_CLIENT_ID)

		cert, err := helpers.GetCertificateRawFromCertOrFilePath(clientCertificate, clientCertificateFilePath)
		if err != nil {
			resp.Diagnostics.AddAttributeError(path.Root("client_certificate"), "Error getting certificate", err.Error())
		}
		p.Config.ClientCertificateRaw = cert
		p.Config.ClientCertificatePassword = clientCertificatePassword
		p.Config.TenantId = tenantId
		p.Config.ClientId = clientId
	} else {
		tflog.Info(ctx, "Using client id and secret for authentication")
		if tenantId != "" && clientId != "" && clientSecret != "" {
			p.Config.TenantId = tenantId
			p.Config.ClientId = clientId
			p.Config.ClientSecret = clientSecret
		} else {
			validateProviderAttribute(resp, path.Root("tenant_id"), "tenant id", tenantId, constants.ENV_VAR_POWER_PLATFORM_TENANT_ID)
			validateProviderAttribute(resp, path.Root("client_id"), "client id", clientId, constants.ENV_VAR_POWER_PLATFORM_CLIENT_ID)
			validateProviderAttribute(resp, path.Root("client_secret"), "client secret", clientSecret, constants.ENV_VAR_POWER_PLATFORM_CLIENT_SECRET)
		}
	}

	var providerConfigUrls *config.ProviderConfigUrls
	var cloudConfiguration *cloud.Configuration
	switch cloudType {
	case "public":
		providerConfigUrls, cloudConfiguration = getCloudPublicUrls()
	case "gcc":
		providerConfigUrls, cloudConfiguration = getGccUrls()
	case "gcchigh":
		providerConfigUrls, cloudConfiguration = getGccHighUrls()
	case "dod":
		providerConfigUrls, cloudConfiguration = getDodUrls()
	case "china":
		providerConfigUrls, cloudConfiguration = getChinaUrls()
	case "ex":
		providerConfigUrls, cloudConfiguration = getExUrls()
	case "rx":
		providerConfigUrls, cloudConfiguration = getRxUrls()
	default:
		resp.Diagnostics.AddAttributeError(
			path.Root("cloud"),
			"Unknown cloud",
			fmt.Sprintf("The provider cannot create the API client as there is an unknown configuration value for `cloud`. Either set the value in the provider configuration or use the '%s' environment variable.", constants.ENV_VAR_POWER_PLATFORM_CLOUD),
		)
	}

	p.Config.Urls = *providerConfigUrls
	p.Config.Cloud = *cloudConfiguration
	p.Config.TelemetryOptout = telemetryOptOut
	p.Config.TerraformVersion = req.TerraformVersion

	providerClient := api.ProviderClient{
		Config: p.Config,
		Api:    p.Api,
	}
	resp.DataSourceData = &providerClient
	resp.ResourceData = &providerClient
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
		func() resource.Resource { return authorization.NewUserResource() },
		func() resource.Resource { return data_record.NewDataRecordResource() },
		func() resource.Resource { return environment_settings.NewEnvironmentSettingsResource() },
		func() resource.Resource { return connection.NewConnectionResource() },
		func() resource.Resource { return rest.NewDataverseWebApiResource() },
		func() resource.Resource { return connection.NewConnectionShareResource() },
		func() resource.Resource { return environment_groups.NewEnvironmentGroupResource() },
		func() resource.Resource { return admin_management_application.NewAdminManagementApplicationResource() },
		func() resource.Resource { return environment_group_rule_set.NewEnvironmentGroupRuleSetResource() },
		func() resource.Resource { return enterprise_policy.NewEnterpisePolicyResource() },
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
		func() datasource.DataSource { return environment_settings.NewEnvironmentSettingsDataSource() },
		func() datasource.DataSource { return locations.NewLocationsDataSource() },
		func() datasource.DataSource { return languages.NewLanguagesDataSource() },
		func() datasource.DataSource { return currencies.NewCurrenciesDataSource() },
		func() datasource.DataSource { return authorization.NewSecurityRolesDataSource() },
		func() datasource.DataSource { return application.NewTenantApplicationPackagesDataSource() },
		func() datasource.DataSource { return data_record.NewDataRecordDataSource() },
		func() datasource.DataSource { return rest.NewDataverseWebApiDatasource() },
		func() datasource.DataSource { return connection.NewConnectionsDataSource() },
		func() datasource.DataSource { return connection.NewConnectionSharesDataSource() },
		func() datasource.DataSource { return capacity.NewTenantCapcityDataSource() },
		func() datasource.DataSource { return tenant.NewTenantDataSource() },
	}
}

func validateProviderAttribute(resp *provider.ConfigureResponse, attrPath path.Path, name, value string, environmentVariableName string) {
	environmentVariableText := "Target apply the source of the value first, set the value statically in the configuration."
	if environmentVariableName != "" {
		environmentVariableText = fmt.Sprintf("Either target apply the source of the value first, set the value statically in the configuration, or use the %s environment variable.", environmentVariableName)
	}

	if value == "" {
		resp.Diagnostics.AddAttributeError(
			attrPath,
			fmt.Sprintf("Unknown %s", name),
			fmt.Sprintf("The provider cannot create the API client as there is an unknown configuration value for %s. %s", name, environmentVariableText))
	}
}

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

func getCloudPublicUrls() (*config.ProviderConfigUrls, *cloud.Configuration) {
	return &config.ProviderConfigUrls{
		BapiUrl:            constants.PUBLIC_BAPI_DOMAIN,
		PowerAppsUrl:       constants.PUBLIC_POWERAPPS_API_DOMAIN,
		PowerAppsScope:     constants.PUBLIC_POWERAPPS_SCOPE,
		PowerPlatformUrl:   constants.PUBLIC_POWERPLATFORM_API_DOMAIN,
		PowerPlatformScope: constants.PUBLIC_POWERPLATFORM_API_SCOPE,
		LicensingUrl:       constants.PUBLIC_LICENSING_API_DOMAIN,
	}, &cloud.AzurePublic
}

func getGccUrls() (*config.ProviderConfigUrls, *cloud.Configuration) {
	return &config.ProviderConfigUrls{
		BapiUrl:            constants.USGOV_BAPI_DOMAIN,
		PowerAppsUrl:       constants.USGOV_POWERAPPS_API_DOMAIN,
		PowerAppsScope:     constants.USGOV_POWERAPPS_SCOPE,
		PowerPlatformUrl:   constants.USGOV_POWERPLATFORM_API_DOMAIN,
		PowerPlatformScope: constants.USGOV_POWERPLATFORM_API_SCOPE,
		LicensingUrl:       constants.USGOV_LICENSING_API_DOMAIN,
	}, &cloud.AzurePublic // GCC uses public cloud for authentication.
}

func getGccHighUrls() (*config.ProviderConfigUrls, *cloud.Configuration) {
	return &config.ProviderConfigUrls{
		BapiUrl:            constants.USGOVHIGH_BAPI_DOMAIN,
		PowerAppsUrl:       constants.USGOVHIGH_POWERAPPS_API_DOMAIN,
		PowerAppsScope:     constants.USGOVHIGH_POWERAPPS_SCOPE,
		PowerPlatformUrl:   constants.USGOVHIGH_POWERPLATFORM_API_DOMAIN,
		PowerPlatformScope: constants.USGOVHIGH_POWERPLATFORM_API_SCOPE,
		LicensingUrl:       constants.USGOVHIGH_LICENSING_API_DOMAIN,
	}, &cloud.AzureGovernment
}

func getDodUrls() (*config.ProviderConfigUrls, *cloud.Configuration) {
	return &config.ProviderConfigUrls{
		BapiUrl:            constants.USDOD_BAPI_DOMAIN,
		PowerAppsUrl:       constants.USDOD_POWERAPPS_API_DOMAIN,
		PowerAppsScope:     constants.USDOD_POWERAPPS_SCOPE,
		PowerPlatformUrl:   constants.USDOD_POWERPLATFORM_API_DOMAIN,
		PowerPlatformScope: constants.USDOD_POWERPLATFORM_API_SCOPE,
		LicensingUrl:       constants.USDOD_LICENSING_API_DOMAIN,
	}, &cloud.AzureGovernment
}

func getChinaUrls() (*config.ProviderConfigUrls, *cloud.Configuration) {
	return &config.ProviderConfigUrls{
		BapiUrl:            constants.CHINA_BAPI_DOMAIN,
		PowerAppsUrl:       constants.CHINA_POWERAPPS_API_DOMAIN,
		PowerAppsScope:     constants.CHINA_POWERAPPS_SCOPE,
		PowerPlatformUrl:   constants.CHINA_POWERPLATFORM_API_DOMAIN,
		PowerPlatformScope: constants.CHINA_POWERPLATFORM_API_SCOPE,
		LicensingUrl:       constants.CHINA_LICENSING_API_DOMAIN,
	}, &cloud.AzureChina
}

func getExUrls() (*config.ProviderConfigUrls, *cloud.Configuration) {
	return &config.ProviderConfigUrls{
			BapiUrl:            constants.EX_BAPI_DOMAIN,
			PowerAppsUrl:       constants.EX_POWERAPPS_API_DOMAIN,
			PowerAppsScope:     constants.EX_POWERAPPS_SCOPE,
			PowerPlatformUrl:   constants.EX_POWERPLATFORM_API_DOMAIN,
			PowerPlatformScope: constants.EX_POWERPLATFORM_API_SCOPE,
			LicensingUrl:       constants.EX_LICENSING_API_DOMAIN,
		}, &cloud.Configuration{
			ActiveDirectoryAuthorityHost: constants.EX_AUTHORITY_HOST,
			Services:                     map[cloud.ServiceName]cloud.ServiceConfiguration{},
		}
}

func getRxUrls() (*config.ProviderConfigUrls, *cloud.Configuration) {
	return &config.ProviderConfigUrls{
			BapiUrl:            constants.RX_BAPI_DOMAIN,
			PowerAppsUrl:       constants.RX_POWERAPPS_API_DOMAIN,
			PowerAppsScope:     constants.RX_POWERAPPS_SCOPE,
			PowerPlatformUrl:   constants.RX_POWERPLATFORM_API_DOMAIN,
			PowerPlatformScope: constants.RX_POWERPLATFORM_API_SCOPE,
			LicensingUrl:       constants.RX_LICENSING_API_DOMAIN,
		}, &cloud.Configuration{
			ActiveDirectoryAuthorityHost: constants.RX_AUTHORITY_HOST,
			Services:                     map[cloud.ServiceName]cloud.ServiceConfiguration{},
		}
}
