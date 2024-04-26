// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package settings

// Cloud	BAPI	Power Apps API	Power Platform API	OAuth Authority
// public	api.bap.microsoft.com	api.powerapps.com	api.powerplatform.com	login.microsoftonline.com
// gcc	gov.api.bap.microsoft.us	gov.api.powerapps.us	api.gov.powerplatform.microsoft.us	login.microsoftonline.com
// gcchigh	high.api.bap.microsoft.us	high.api.powerapps.us	api.high.powerplatform.microsoft.us	login.microsoftonline.us
// dod	api.bap.appsplatform.us	api.apps.appsplatform.us	api.appsplatform.us	login.microsoftonline.us
// ex	api.bap.eaglex.ic.gov	api.powerapps.eaglex.ic.gov	api.powerplatform.eaglex.ic.gov	login.microsoftonline.eaglex.ic.gov
// rx	api.bap.microsoft.scloud	api.powerapps.microsoft.scloud	api.powerplatform.microsoft.scloud	login.microsoftonline.microsoft.scloud
// china	api.bap.partner.microsoftonline.cn	api.powerapps.cn	api.powerplatform.partner.microsoftonline.cn	login.chinacloudapi.cn

const (
	PUBLIC_OAUTH_AUTHORITY_URL      = "https://login.microsoftonline.com/"
	PUBLIC_BAPI_DOMAIN              = "api.bap.microsoft.com"
	PUBLIC_POWERAPPS_API_DOMAIN     = "api.powerapps.com"
	PUBLIC_POWERAPPS_SCOPE          = "https://service.powerapps.com/.default"
	PUBLIC_POWERPLATFORM_API_DOMAIN = "api.powerplatform.com"
	PUBLIC_POWERPLATFORM_API_SCOPE  = "https://api.powerplatform.com/.default"
)

const (
	USDOD_OAUTH_AUTHORITY_URL      = "https://login.microsoftonline.us/"
	USDOD_BAPI_DOMAIN              = "api.bap.appsplatform.us"
	USDOD_POWERAPPS_API_DOMAIN     = "api.apps.appsplatform.us"
	USDOD_POWERAPPS_SCOPE          = "https://service.apps.appsplatform.us/.default"
	USDOD_POWERPLATFORM_API_DOMAIN = "api.appsplatform.us"
	USDOD_POWERPLATFORM_API_SCOPE  = "https://api.appsplatform.us/.default"
)

const (
	USGOV_OAUTH_AUTHORITY_URL      = "https://login.microsoftonline.com/"
	USGOV_BAPI_DOMAIN              = "gov.api.bap.microsoft.us"
	USGOV_POWERAPPS_API_DOMAIN     = "gov.api.powerapps.us"
	USGOV_POWERAPPS_SCOPE          = "https://service.powerapps.us/.default"
	USGOV_POWERPLATFORM_API_DOMAIN = "api.gov.powerplatform.microsoft.us"
	USGOV_POWERPLATFORM_API_SCOPE  = "https://api.gov.powerplatform.microsoft.us/.default"
)

const (
	USGOVHIGH_OAUTH_AUTHORITY_URL      = "https://login.microsoftonline.us/"
	USGOVHIGH_BAPI_DOMAIN              = "high.api.bap.microsoft.us"
	USGOVHIGH_POWERAPPS_API_DOMAIN     = "high.api.powerapps.us"
	USGOVHIGH_POWERAPPS_SCOPE          = "https://high.service.apps.appsplatform.us/.default"
	USGOVHIGH_POWERPLATFORM_API_DOMAIN = "api.appsplatform.us"
	USGOVHIGH_POWERPLATFORM_API_SCOPE  = "https://api.appsplatform.us/.default"
)

const (
	CHINA_OAUTH_AUTHORITY_URL      = "https://login.chinacloudapi.cn/"
	CHINA_BAPI_DOMAIN              = "api.bap.partner.microsoftonline.cn"
	CHINA_POWERAPPS_API_DOMAIN     = "api.powerapps.cn"
	CHINA_POWERAPPS_SCOPE          = "https://service.powerapps.cn/.default"
	CHINA_POWERPLATFORM_API_DOMAIN = "api.powerplatform.partner.microsoftonline.cn"
	CHINA_POWERPLATFORM_API_SCOPE  = "https://api.powerplatform.partner.microsoftonline.cn/.default"
)

const (
	EX_OAUTH_AUTHORITY_URL      = "https://login.microsoftonline.eaglex.ic.gov/"
	EX_BAPI_DOMAIN              = "api.bap.eaglex.ic.gov"
	EX_POWERAPPS_API_DOMAIN     = "api.powerapps.eaglex.ic.gov"
	EX_POWERAPPS_SCOPE          = "https://service.powerapps.eaglex.ic.gov/.default"
	EX_POWERPLATFORM_API_DOMAIN = "api.powerplatform.eaglex.ic.gov"
	EX_POWERPLATFORM_API_SCOPE  = "https://api.powerplatform.eaglex.ic.gov/.default"
)

const (
	RX_OAUTH_AUTHORITY_URL      = "https://login.microsoftonline.microsoft.scloud/"
	RX_BAPI_DOMAIN              = "api.bap.microsoft.scloud"
	RX_POWERAPPS_API_DOMAIN     = "api.powerapps.microsoft.scloud"
	RX_POWERAPPS_SCOPE          = "https://service.powerapps.microsoft.scloud/.default"
	RX_POWERPLATFORM_API_DOMAIN = "api.powerplatform.microsoft.scloud"
	RX_POWERPLATFORM_API_SCOPE  = "https://api.powerplatform.microsoft.scloud/.default"
)
