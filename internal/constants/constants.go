// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package constants

import "time"

// Cloud	BAPI	Power Apps API	Power Platform API	OAuth Authority
// public	api.bap.microsoft.com	api.powerapps.com	api.powerplatform.com	login.microsoftonline.com
// gcc	gov.api.bap.microsoft.us	gov.api.powerapps.us	api.gov.powerplatform.microsoft.us	login.microsoftonline.com
// gcchigh	high.api.bap.microsoft.us	high.api.powerapps.us	api.high.powerplatform.microsoft.us	login.microsoftonline.us
// dod	api.bap.appsplatform.us	api.apps.appsplatform.us	api.appsplatform.us	login.microsoftonline.us
// ex	api.bap.eaglex.ic.gov	api.powerapps.eaglex.ic.gov	api.powerplatform.eaglex.ic.gov	login.microsoftonline.eaglex.ic.gov
// rx	api.bap.microsoft.scloud	api.powerapps.microsoft.scloud	api.powerplatform.microsoft.scloud	login.microsoftonline.microsoft.scloud
// china	api.bap.partner.microsoftonline.cn	api.powerapps.cn	api.powerplatform.partner.microsoftonline.cn	login.chinacloudapi.cn

const ZERO_UUID = "00000000-0000-0000-0000-000000000000"

const (
	PUBLIC_OAUTH_AUTHORITY_URL          = "https://login.microsoftonline.com/"
	PUBLIC_BAPI_DOMAIN                  = "api.bap.microsoft.com"
	PUBLIC_POWERAPPS_API_DOMAIN         = "api.powerapps.com"
	PUBLIC_POWERAPPS_SCOPE              = "https://service.powerapps.com/.default"
	PUBLIC_POWERPLATFORM_API_DOMAIN     = "api.powerplatform.com"
	PUBLIC_POWERPLATFORM_API_SCOPE      = "https://api.powerplatform.com/.default"
	PUBLIC_LICENSING_API_DOMAIN         = "licensing.powerplatform.microsoft.com"
	PUBLIC_POWERAPPS_ADVISOR_API_DOMAIN = "api.advisor.powerapps.com"
	PUBLIC_POWERAPPS_ADVISOR_API_SCOPE  = "https://advisor.powerapps.com/.default"
)

const (
	USDOD_OAUTH_AUTHORITY_URL          = "https://login.microsoftonline.us/"
	USDOD_BAPI_DOMAIN                  = "api.bap.appsplatform.us"
	USDOD_POWERAPPS_API_DOMAIN         = "api.apps.appsplatform.us"
	USDOD_POWERAPPS_SCOPE              = "https://service.apps.appsplatform.us/.default"
	USDOD_POWERPLATFORM_API_DOMAIN     = "api.appsplatform.us"
	USDOD_POWERPLATFORM_API_SCOPE      = "https://api.appsplatform.us/.default"
	USDOD_LICENSING_API_DOMAIN         = "licensing.appsplatform.us"
	USDOD_POWERAPPS_ADVISOR_API_DOMAIN = "api.advisor.powerapps.us"
	USDOD_POWERAPPS_ADVISOR_API_SCOPE  = "https://advisor.powerapps.us/.default"
)

const (
	USGOV_OAUTH_AUTHORITY_URL          = "https://login.microsoftonline.com/"
	USGOV_BAPI_DOMAIN                  = "gov.api.bap.microsoft.us"
	USGOV_POWERAPPS_API_DOMAIN         = "gov.api.powerapps.us"
	USGOV_POWERAPPS_SCOPE              = "https://service.powerapps.us/.default"
	USGOV_POWERPLATFORM_API_DOMAIN     = "api.gov.powerplatform.microsoft.us"
	USGOV_POWERPLATFORM_API_SCOPE      = "https://api.gov.powerplatform.microsoft.us/.default"
	USGOV_LICENSING_API_DOMAIN         = "gov.licensing.powerplatform.microsoft.us"
	USGOV_POWERAPPS_ADVISOR_API_DOMAIN = "gov.api.advisor.powerapps.us"
	USGOV_POWERAPPS_ADVISOR_API_SCOPE  = "https://gov.advisor.powerapps.us/.default"
)

const (
	USGOVHIGH_OAUTH_AUTHORITY_URL          = "https://login.microsoftonline.us/"
	USGOVHIGH_BAPI_DOMAIN                  = "high.api.bap.microsoft.us"
	USGOVHIGH_POWERAPPS_API_DOMAIN         = "high.api.powerapps.us"
	USGOVHIGH_POWERAPPS_SCOPE              = "https://high.service.apps.appsplatform.us/.default"
	USGOVHIGH_POWERPLATFORM_API_DOMAIN     = "api.appsplatform.us"
	USGOVHIGH_POWERPLATFORM_API_SCOPE      = "https://api.appsplatform.us/.default"
	USGOVHIGH_LICENSING_API_DOMAIN         = "high.licensing.powerplatform.microsoft.us"
	USGOVHIGH_POWERAPPS_ADVISOR_API_DOMAIN = "high.api.advisor.powerapps.us"
	USGOVHIGH_POWERAPPS_ADVISOR_API_SCOPE  = "https://high.advisor.powerapps.us/.default"
)

const (
	CHINA_OAUTH_AUTHORITY_URL          = "https://login.chinacloudapi.cn/"
	CHINA_BAPI_DOMAIN                  = "api.bap.partner.microsoftonline.cn"
	CHINA_POWERAPPS_API_DOMAIN         = "api.powerapps.cn"
	CHINA_POWERAPPS_SCOPE              = "https://service.powerapps.cn/.default"
	CHINA_POWERPLATFORM_API_DOMAIN     = "api.powerplatform.partner.microsoftonline.cn"
	CHINA_POWERPLATFORM_API_SCOPE      = "https://api.powerplatform.partner.microsoftonline.cn/.default"
	CHINA_LICENSING_API_DOMAIN         = "licensing.partner.microsoftonline.cn"
	CHINA_POWERAPPS_ADVISOR_API_DOMAIN = "api.advisor.powerapps.cn"
	CHINA_POWERAPPS_ADVISOR_API_SCOPE  = "https://advisor.powerapps.cn/.default"
)

const (
	EX_OAUTH_AUTHORITY_URL          = "https://login.microsoftonline.eaglex.ic.gov/"
	EX_BAPI_DOMAIN                  = "api.bap.eaglex.ic.gov"
	EX_POWERAPPS_API_DOMAIN         = "api.powerapps.eaglex.ic.gov"
	EX_POWERAPPS_SCOPE              = "https://service.powerapps.eaglex.ic.gov/.default"
	EX_POWERPLATFORM_API_DOMAIN     = "api.powerplatform.eaglex.ic.gov"
	EX_POWERPLATFORM_API_SCOPE      = "https://api.powerplatform.eaglex.ic.gov/.default"
	EX_AUTHORITY_HOST               = "https://login.microsoftonline.eaglex.ic.gov/"
	EX_LICENSING_API_DOMAIN         = "licensing.eaglex.ic.gov"
	EX_POWERAPPS_ADVISOR_API_DOMAIN = "api.advisor.powerapps.eaglex.ic.gov"
	EX_POWERAPPS_ADVISOR_API_SCOPE  = "https://advisor.powerapps.eaglex.ic.gov/.default"
)

const (
	RX_OAUTH_AUTHORITY_URL          = "https://login.microsoftonline.microsoft.scloud/"
	RX_BAPI_DOMAIN                  = "api.bap.microsoft.scloud"
	RX_POWERAPPS_API_DOMAIN         = "api.powerapps.microsoft.scloud"
	RX_POWERAPPS_SCOPE              = "https://service.powerapps.microsoft.scloud/.default"
	RX_POWERPLATFORM_API_DOMAIN     = "api.powerplatform.microsoft.scloud"
	RX_POWERPLATFORM_API_SCOPE      = "https://api.powerplatform.microsoft.scloud/.default"
	RX_AUTHORITY_HOST               = "https://login.microsoftonline.microsoft.scloud/"
	RX_LICENSING_API_DOMAIN         = "licensing.microsoft.scloud"
	RX_POWERAPPS_ADVISOR_API_DOMAIN = "api.advisor.powerapps.eaglex.ic.gov"
	RX_POWERAPPS_ADVISOR_API_SCOPE  = "https://advisor.powerapps.eaglex.ic.gov/.default"
)

const (
	COPILOT_SCOPE = "96ff4394-9197-43aa-b393-6a41652e21f8"
)

const (
	DATAVERSE_API_VERSION     = "v9.2"
	HEADER_ODATA_ENTITY_ID    = "Odata-Entityid"
	HEADER_LOCATION           = "Location"
	HEADER_OPERATION_LOCATION = "Operation-Location"
	HEADER_RETRY_AFTER        = "Retry-After"
	HTTPS                     = "https"
	API_VERSION_PARAM         = "api-version"
)

const (
	DEFAULT_RESOURCE_OPERATION_TIMEOUT_IN_MINUTES = 20 * time.Minute
)

const (
	ENV_VAR_POWER_PLATFORM_CLOUD                        = "POWER_PLATFORM_CLOUD"
	ENV_VAR_POWER_PLATFORM_TENANT_ID                    = "POWER_PLATFORM_TENANT_ID"
	ENV_VAR_POWER_PLATFORM_CLIENT_ID                    = "POWER_PLATFORM_CLIENT_ID"
	ENV_VAR_POWER_PLATFORM_CLIENT_SECRET                = "POWER_PLATFORM_CLIENT_SECRET"
	ENV_VAR_POWER_PLATFORM_USE_OIDC                     = "POWER_PLATFORM_USE_OIDC"
	ENV_VAR_POWER_PLATFORM_USE_CLI                      = "POWER_PLATFORM_USE_CLI"
	ENV_VAR_POWER_PLATFORM_USE_MSI                      = "POWER_PLATFORM_USE_MSI"
	ENV_VAR_POWER_PLATFORM_CLIENT_CERTIFICATE           = "POWER_PLATFORM_CLIENT_CERTIFICATE"
	ENV_VAR_POWER_PLATFORM_CLIENT_CERTIFICATE_FILE_PATH = "POWER_PLATFORM_CLIENT_CERTIFICATE_FILE_PATH"
	ENV_VAR_POWER_PLATFORM_CLIENT_CERTIFICATE_PASSWORD  = "POWER_PLATFORM_CLIENT_CERTIFICATE_PASSWORD"
	ENV_VAR_POWER_PLATFORM_TELEMETRY_OPTOUT             = "POWER_PLATFORM_TELEMETRY_OPTOUT"
	ENV_VAR_POWER_PLATFORM_AZDO_SERVICE_CONNECTION_ID   = "POWER_PLATFORM_AZDO_SERVICE_CONNECTION_ID"

	ENV_VAR_ARM_OIDC_REQUEST_URL           = "ARM_OIDC_REQUEST_URL"
	ENV_VAR_ACTIONS_ID_TOKEN_REQUEST_URL   = "ACTIONS_ID_TOKEN_REQUEST_URL"
	ENV_VAR_ARM_OIDC_REQUEST_TOKEN         = "ARM_OIDC_REQUEST_TOKEN"
	ENV_VAR_ACTIONS_ID_TOKEN_REQUEST_TOKEN = "ACTIONS_ID_TOKEN_REQUEST_TOKEN"
	ENV_VAR_ARM_OIDC_TOKEN                 = "ARM_OIDC_TOKEN"
	ENV_VAR_ARM_OIDC_TOKEN_FILE_PATH       = "ARM_OIDC_TOKEN_FILE_PATH"
)

const (
	AZURE_AD_PROVIDER_VERSION_CONSTRAINT = ">= 3.0.1"
	RANDOM_PROVIDER_VERSION_CONSTRAINT   = ">= 3.6.3"
	AZAPI_PROVIDER_VERSION_CONSTRAINT    = ">= 1.15.0"
)

const (
	SOLUTION_CHECKER_RULESET_ID = "0ad12346-e108-40b8-a956-9a8f95ea18c9"
)
