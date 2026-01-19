// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package constants

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUnitCloudDomainsAndScopes(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		got  string
		want string
	}{
		{"ZERO_UUID", ZERO_UUID, "00000000-0000-0000-0000-000000000000"},
		{"PUBLIC_ADMIN_POWER_PLATFORM_URL", PUBLIC_ADMIN_POWER_PLATFORM_URL, "api.admin.powerplatform.microsoft.com"},
		{"PUBLIC_BAPI_DOMAIN", PUBLIC_BAPI_DOMAIN, "api.bap.microsoft.com"},
		{"PUBLIC_POWERAPPS_SCOPE", PUBLIC_POWERAPPS_SCOPE, "https://service.powerapps.com/.default"},
		{"PUBLIC_POWERPLATFORM_API_SCOPE", PUBLIC_POWERPLATFORM_API_SCOPE, "https://api.powerplatform.com/.default"},
		{"PUBLIC_LICENSING_API_DOMAIN", PUBLIC_LICENSING_API_DOMAIN, "licensing.powerplatform.microsoft.com"},
		{"USDOD_BAPI_DOMAIN", USDOD_BAPI_DOMAIN, "api.bap.appsplatform.us"},
		{"USDOD_POWERPLATFORM_API_SCOPE", USDOD_POWERPLATFORM_API_SCOPE, "https://api.appsplatform.us/.default"},
		{"USGOV_BAPI_DOMAIN", USGOV_BAPI_DOMAIN, "gov.api.bap.microsoft.us"},
		{"USGOV_POWERPLATFORM_API_SCOPE", USGOV_POWERPLATFORM_API_SCOPE, "https://api.gov.powerplatform.microsoft.us/.default"},
		{"USGOVHIGH_BAPI_DOMAIN", USGOVHIGH_BAPI_DOMAIN, "high.api.bap.microsoft.us"},
		{"USGOVHIGH_POWERAPPS_SCOPE", USGOVHIGH_POWERAPPS_SCOPE, "https://high.service.apps.appsplatform.us/.default"},
		{"CHINA_POWERPLATFORM_API_DOMAIN", CHINA_POWERPLATFORM_API_DOMAIN, "api.powerplatform.partner.microsoftonline.cn"},
		{"EX_POWERAPPS_SCOPE", EX_POWERAPPS_SCOPE, "https://service.powerapps.eaglex.ic.gov/.default"},
		{"RX_POWERPLATFORM_API_DOMAIN", RX_POWERPLATFORM_API_DOMAIN, "api.powerplatform.microsoft.scloud"},
		{"COPILOT_SCOPE", COPILOT_SCOPE, "96ff4394-9197-43aa-b393-6a41652e21f8"},
		{"PPAC_SCOPE", PPAC_SCOPE, "065d9450-1e87-434e-ac2f-69af271549ed"},
	}

	for _, tc := range cases {
		require.Equal(t, tc.want, tc.got, tc.name)
	}
}

func TestUnitTimeoutsAndRetryConstants(t *testing.T) {
	t.Parallel()

	require.Equal(t, 20*time.Minute, DEFAULT_RESOURCE_OPERATION_TIMEOUT_IN_MINUTES)
	require.Equal(t, 10, MAX_RETRY_COUNT)
}

func TestUnitHeaderAndErrorConstants(t *testing.T) {
	t.Parallel()

	require.Equal(t, "Odata-Entityid", HEADER_ODATA_ENTITY_ID)
	require.Equal(t, "Location", HEADER_LOCATION)
	require.Equal(t, "Operation-Location", HEADER_OPERATION_LOCATION)
	require.Equal(t, "Retry-After", HEADER_RETRY_AFTER)
	require.Equal(t, "authorization has been denied for this request. Make sure that your service principal is registered as an admin management application: https://learn.microsoft.com/en-us/power-platform/admin/powerplatform-api-create-service-principal#registering-an-admin-management-application", NO_MANAGEMENT_APPLICATION_ERROR_MSG)
	require.Equal(t, "claims=", CAE_CHALLENGE_CLAIMS_INDICATOR)
	require.Equal(t, "insufficient_claims", CAE_CHALLENGE_INSUFFICIENT_CLAIMS_INDICATOR)
	require.Equal(t, "OBJECT_NOT_FOUND", ERROR_OBJECT_NOT_FOUND)
	require.Equal(t, "ENVIRONMENT_URL_NOT_FOUND", ERROR_ENVIRONMENT_URL_NOT_FOUND)
	require.Equal(t, "ENVIRONMENTS_IN_ENV_GROUP", ERROR_ENVIRONMENTS_IN_ENV_GROUP)
	require.Equal(t, "POLICY_ASSIGNED_TO_ENV_GROUP", ERROR_POLICY_ASSIGNED_TO_ENV_GROUP)
	require.Equal(t, "ENVIRONMENT_SETTINGS_FAILED", ERROR_ENVIRONMENT_SETTINGS_FAILED)
	require.Equal(t, "ENVIRONMENT_CREATION", ERROR_ENVIRONMENT_CREATION)
}

func TestUnitEnvVarConstants(t *testing.T) {
	t.Parallel()

	require.Equal(t, "POWER_PLATFORM_CLOUD", ENV_VAR_POWER_PLATFORM_CLOUD)
	require.Equal(t, "POWER_PLATFORM_TENANT_ID", ENV_VAR_POWER_PLATFORM_TENANT_ID)
	require.Equal(t, "POWER_PLATFORM_CLIENT_ID", ENV_VAR_POWER_PLATFORM_CLIENT_ID)
	require.Equal(t, "POWER_PLATFORM_CLIENT_SECRET", ENV_VAR_POWER_PLATFORM_CLIENT_SECRET)
	require.Equal(t, "POWER_PLATFORM_USE_OIDC", ENV_VAR_POWER_PLATFORM_USE_OIDC)
	require.Equal(t, "POWER_PLATFORM_USE_CLI", ENV_VAR_POWER_PLATFORM_USE_CLI)
	require.Equal(t, "POWER_PLATFORM_USE_DEV_CLI", ENV_VAR_POWER_PLATFORM_USE_DEV_CLI)
	require.Equal(t, "POWER_PLATFORM_USE_MSI", ENV_VAR_POWER_PLATFORM_USE_MSI)
	require.Equal(t, "POWER_PLATFORM_CLIENT_CERTIFICATE", ENV_VAR_POWER_PLATFORM_CLIENT_CERTIFICATE)
	require.Equal(t, "POWER_PLATFORM_CLIENT_CERTIFICATE_FILE_PATH", ENV_VAR_POWER_PLATFORM_CLIENT_CERTIFICATE_FILE_PATH)
	require.Equal(t, "POWER_PLATFORM_CLIENT_CERTIFICATE_PASSWORD", ENV_VAR_POWER_PLATFORM_CLIENT_CERTIFICATE_PASSWORD)
	require.Equal(t, "POWER_PLATFORM_TELEMETRY_OPTOUT", ENV_VAR_POWER_PLATFORM_TELEMETRY_OPTOUT)
	require.Equal(t, "POWER_PLATFORM_AZDO_SERVICE_CONNECTION_ID", ENV_VAR_POWER_PLATFORM_AZDO_SERVICE_CONNECTION_ID)
	require.Equal(t, "POWER_PLATFORM_ENABLE_CAE", ENV_VAR_POWER_PLATFORM_ENABLE_CAE)
	require.Equal(t, "ARM_OIDC_REQUEST_URL", ENV_VAR_ARM_OIDC_REQUEST_URL)
	require.Equal(t, "ACTIONS_ID_TOKEN_REQUEST_URL", ENV_VAR_ACTIONS_ID_TOKEN_REQUEST_URL)
	require.Equal(t, "ARM_OIDC_REQUEST_TOKEN", ENV_VAR_ARM_OIDC_REQUEST_TOKEN)
	require.Equal(t, "ACTIONS_ID_TOKEN_REQUEST_TOKEN", ENV_VAR_ACTIONS_ID_TOKEN_REQUEST_TOKEN)
	require.Equal(t, "ARM_OIDC_TOKEN", ENV_VAR_ARM_OIDC_TOKEN)
	require.Equal(t, "ARM_OIDC_TOKEN_FILE_PATH", ENV_VAR_ARM_OIDC_TOKEN_FILE_PATH)
	require.Equal(t, "ARM_AUXILIARY_TENANT_IDS", ENV_VAR_ARM_AUXILIARY_TENANT_IDS)
	require.Equal(t, "POWER_PLATFORM_PARTNER_ID", ENV_VAR_POWER_PLATFORM_PARTNER_ID)
	require.Equal(t, "ARM_PARTNER_ID", ENV_VAR_ARM_PARTNER_ID)
	require.Equal(t, "POWER_PLATFORM_DISABLE_TERRAFORM_PARTNER_ID", ENV_VAR_POWER_PLATFORM_DISABLE_TERRAFORM_PARTNER_ID)
	require.Equal(t, "ARM_DISABLE_TERRAFORM_PARTNER_ID", ENV_VAR_ARM_DISABLE_TERRAFORM_PARTNER_ID)
}

func TestUnitVersionAndScopeConstants(t *testing.T) {
	t.Parallel()

	require.Equal(t, ">= 3.0.1", AZURE_AD_PROVIDER_VERSION_CONSTRAINT)
	require.Equal(t, ">= 3.6.3", RANDOM_PROVIDER_VERSION_CONSTRAINT)
	require.Equal(t, ">= 1.15.0", AZAPI_PROVIDER_VERSION_CONSTRAINT)
	require.Equal(t, "v9.2", DATAVERSE_API_VERSION)
	require.Equal(t, "2020-10-01", ADMIN_MANAGEMENT_APP_API_VERSION)
	require.Equal(t, "2019-10-01", ENTERPRISE_POLICY_API_VERSION)
	require.Equal(t, "2023-06-01", BAP_API_VERSION)
	require.Equal(t, "2021-04-01", BAP_2021_API_VERSION)
	require.Equal(t, "2022-05-01", BAP_2022_API_VERSION)
	require.Equal(t, "2022-03-01-preview", APPLICATION_API_VERSION)
	require.Equal(t, "2021-10-01-preview", ENVIRONMENT_GROUP_API_VERSION)
	require.Equal(t, "2019-05-01", CONNECTORS_API_VERSION)
	require.Equal(t, "2020-08-01", TENANT_SETTINGS_API_VERSION)
	require.Equal(t, "0ad12346-e108-40b8-a956-9a8f95ea18c9", SOLUTION_CHECKER_RULESET_ID)
	require.Equal(t, "222c6c49-1b0a-5959-a213-6608f9eb8820", DEFAULT_TERRAFORM_PARTNER_ID)
}
