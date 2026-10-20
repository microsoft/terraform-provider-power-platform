// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

import (
	"testing"

	"github.com/microsoft/terraform-provider-power-platform/internal/config"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/stretchr/testify/require"
)

const scopeTestId = "00000000-0000-0000-0000-000000000001"

type scopeTestCloud struct {
	name          string
	urls          config.ProviderConfigUrls
	dataverseHost string
}

func scopeTestClouds() []scopeTestCloud {
	return []scopeTestCloud{
		{
			name: "public",
			urls: config.ProviderConfigUrls{
				AdminPowerPlatformUrl: constants.PUBLIC_ADMIN_POWER_PLATFORM_URL,
				BapiUrl:               constants.PUBLIC_BAPI_DOMAIN,
				PowerAppsUrl:          constants.PUBLIC_POWERAPPS_API_DOMAIN,
				PowerAppsScope:        constants.PUBLIC_POWERAPPS_SCOPE,
				PowerPlatformUrl:      constants.PUBLIC_POWERPLATFORM_API_DOMAIN,
				PowerPlatformScope:    constants.PUBLIC_POWERPLATFORM_API_SCOPE,
				LicensingUrl:          constants.PUBLIC_LICENSING_API_DOMAIN,
				PowerAppsAdvisor:      constants.PUBLIC_POWERAPPS_ADVISOR_API_DOMAIN,
				PowerAppsAdvisorScope: constants.PUBLIC_POWERAPPS_ADVISOR_API_SCOPE,
				AnalyticsScope:        constants.PUBLIC_ANALYTICS_SCOPE,
			},
			dataverseHost: "org.crm4.dynamics.com",
		},
		{
			name: "gcc",
			urls: config.ProviderConfigUrls{
				AdminPowerPlatformUrl: constants.USGOV_ADMIN_POWER_PLATFORM_URL,
				BapiUrl:               constants.USGOV_BAPI_DOMAIN,
				PowerAppsUrl:          constants.USGOV_POWERAPPS_API_DOMAIN,
				PowerAppsScope:        constants.USGOV_POWERAPPS_SCOPE,
				PowerPlatformUrl:      constants.USGOV_POWERPLATFORM_API_DOMAIN,
				PowerPlatformScope:    constants.USGOV_POWERPLATFORM_API_SCOPE,
				LicensingUrl:          constants.USGOV_LICENSING_API_DOMAIN,
				PowerAppsAdvisor:      constants.USGOV_POWERAPPS_ADVISOR_API_DOMAIN,
				PowerAppsAdvisorScope: constants.USGOV_POWERAPPS_ADVISOR_API_SCOPE,
				AnalyticsScope:        constants.USGOV_ANALYTICS_SCOPE,
			},
			dataverseHost: "org.crm9.dynamics.com",
		},
		{
			name: "gcchigh",
			urls: config.ProviderConfigUrls{
				AdminPowerPlatformUrl: constants.USGOVHIGH_ADMIN_POWER_PLATFORM_URL,
				BapiUrl:               constants.USGOVHIGH_BAPI_DOMAIN,
				PowerAppsUrl:          constants.USGOVHIGH_POWERAPPS_API_DOMAIN,
				PowerAppsScope:        constants.USGOVHIGH_POWERAPPS_SCOPE,
				PowerPlatformUrl:      constants.USGOVHIGH_POWERPLATFORM_API_DOMAIN,
				PowerPlatformScope:    constants.USGOVHIGH_POWERPLATFORM_API_SCOPE,
				LicensingUrl:          constants.USGOVHIGH_LICENSING_API_DOMAIN,
				PowerAppsAdvisor:      constants.USGOVHIGH_POWERAPPS_ADVISOR_API_DOMAIN,
				PowerAppsAdvisorScope: constants.USGOVHIGH_POWERAPPS_ADVISOR_API_SCOPE,
				AnalyticsScope:        constants.USGOVHIGH_ANALYTICS_SCOPE,
			},
			dataverseHost: "org.crm.microsoftdynamics.us",
		},
		{
			name: "dod",
			urls: config.ProviderConfigUrls{
				AdminPowerPlatformUrl: constants.USDOD_ADMIN_POWER_PLATFORM_URL,
				BapiUrl:               constants.USDOD_BAPI_DOMAIN,
				PowerAppsUrl:          constants.USDOD_POWERAPPS_API_DOMAIN,
				PowerAppsScope:        constants.USDOD_POWERAPPS_SCOPE,
				PowerPlatformUrl:      constants.USDOD_POWERPLATFORM_API_DOMAIN,
				PowerPlatformScope:    constants.USDOD_POWERPLATFORM_API_SCOPE,
				LicensingUrl:          constants.USDOD_LICENSING_API_DOMAIN,
				PowerAppsAdvisor:      constants.USDOD_POWERAPPS_ADVISOR_API_DOMAIN,
				PowerAppsAdvisorScope: constants.USDOD_POWERAPPS_ADVISOR_API_SCOPE,
				AnalyticsScope:        constants.USDOD_ANALYTICS_SCOPE,
			},
			dataverseHost: "org.crm.appsplatform.us",
		},
		{
			name: "china",
			urls: config.ProviderConfigUrls{
				AdminPowerPlatformUrl: constants.CHINA_ADMIN_POWER_PLATFORM_URL,
				BapiUrl:               constants.CHINA_BAPI_DOMAIN,
				PowerAppsUrl:          constants.CHINA_POWERAPPS_API_DOMAIN,
				PowerAppsScope:        constants.CHINA_POWERAPPS_SCOPE,
				PowerPlatformUrl:      constants.CHINA_POWERPLATFORM_API_DOMAIN,
				PowerPlatformScope:    constants.CHINA_POWERPLATFORM_API_SCOPE,
				LicensingUrl:          constants.CHINA_LICENSING_API_DOMAIN,
				PowerAppsAdvisor:      constants.CHINA_POWERAPPS_ADVISOR_API_DOMAIN,
				PowerAppsAdvisorScope: constants.CHINA_POWERAPPS_ADVISOR_API_SCOPE,
				AnalyticsScope:        constants.CHINA_ANALYTICS_SCOPE,
			},
			dataverseHost: "org.crm.dynamics.cn",
		},
		{
			name: "ex",
			urls: config.ProviderConfigUrls{
				AdminPowerPlatformUrl: constants.EX_ADMIN_POWER_PLATFORM_URL,
				BapiUrl:               constants.EX_BAPI_DOMAIN,
				PowerAppsUrl:          constants.EX_POWERAPPS_API_DOMAIN,
				PowerAppsScope:        constants.EX_POWERAPPS_SCOPE,
				PowerPlatformUrl:      constants.EX_POWERPLATFORM_API_DOMAIN,
				PowerPlatformScope:    constants.EX_POWERPLATFORM_API_SCOPE,
				LicensingUrl:          constants.EX_LICENSING_API_DOMAIN,
				PowerAppsAdvisor:      constants.EX_POWERAPPS_ADVISOR_API_DOMAIN,
				PowerAppsAdvisorScope: constants.EX_POWERAPPS_ADVISOR_API_SCOPE,
				AnalyticsScope:        constants.EX_ANALYTICS_SCOPE,
			},
			dataverseHost: "org.crm.eaglex.ic.gov",
		},
		{
			name: "rx",
			urls: config.ProviderConfigUrls{
				AdminPowerPlatformUrl: constants.RX_ADMIN_POWER_PLATFORM_URL,
				BapiUrl:               constants.RX_BAPI_DOMAIN,
				PowerAppsUrl:          constants.RX_POWERAPPS_API_DOMAIN,
				PowerAppsScope:        constants.RX_POWERAPPS_SCOPE,
				PowerPlatformUrl:      constants.RX_POWERPLATFORM_API_DOMAIN,
				PowerPlatformScope:    constants.RX_POWERPLATFORM_API_SCOPE,
				LicensingUrl:          constants.RX_LICENSING_API_DOMAIN,
				PowerAppsAdvisor:      constants.RX_POWERAPPS_ADVISOR_API_DOMAIN,
				PowerAppsAdvisorScope: constants.RX_POWERAPPS_ADVISOR_API_SCOPE,
				AnalyticsScope:        constants.RX_ANALYTICS_SCOPE,
			},
			dataverseHost: "org.crm.microsoft.scloud",
		},
	}
}

func publicScopeTestUrls(t *testing.T) config.ProviderConfigUrls {
	t.Helper()
	clouds := scopeTestClouds()
	require.Equal(t, "public", clouds[0].name)
	return clouds[0].urls
}

func TestUnitTryGetScopeFromURL_KnownHosts(t *testing.T) {
	t.Parallel()

	for _, cloud := range scopeTestClouds() {
		urls := cloud.urls
		cases := []struct {
			name  string
			url   string
			scope string
		}{
			{"bapi", "https://" + urls.BapiUrl + "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?api-version=2023-06-01", urls.PowerAppsScope},
			{"regional bapi", "https://europe." + urls.BapiUrl + "/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/1", urls.PowerAppsScope},
			{"powerapps", "https://" + urls.PowerAppsUrl + "/providers/Microsoft.PowerApps/apps", urls.PowerAppsScope},
			{"power platform", "https://" + urls.PowerPlatformUrl + "/appmanagement/applicationPackages", urls.PowerPlatformScope},
			{"environment host", "https://" + helpers.BuildEnvironmentHostUri(scopeTestId, urls.PowerPlatformUrl) + "/connectivity/connectors", urls.PowerPlatformScope},
			{"tenant host", "https://" + helpers.BuildTenantHostUri(scopeTestId, urls.PowerPlatformUrl) + "/gateway/cluster?api-version=1", urls.PowerPlatformScope},
			{"regional advisor", "https://unitedstates." + urls.PowerAppsAdvisor + "/api/rule", urls.PowerAppsAdvisorScope},
			{"admin", "https://" + urls.AdminPowerPlatformUrl + "/api/tenant", constants.PPAC_SCOPE},
			{"licensing", "https://" + urls.LicensingUrl + "/v0.1-alpha/tenants", "https://" + urls.LicensingUrl + "/.default"},
			{"dataverse", "https://" + cloud.dataverseHost + "/api/data/v9.2/accounts", "https://" + cloud.dataverseHost + "/.default"},
		}
		for _, tc := range cases {
			t.Run(cloud.name+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				scope, err := tryGetScopeFromURL(tc.url, urls)
				require.NoError(t, err)
				require.Equal(t, tc.scope, scope)
			})
		}
	}
}

func TestUnitTryGetScopeFromURL_AnalyticsHosts(t *testing.T) {
	t.Parallel()

	urls := publicScopeTestUrls(t)
	for _, url := range []string{
		"https://na.csanalytics.powerplatform.microsoft.com/api/v2/connections",
		"https://gcc.csanalytics.powerplatform.microsoft.us/api/v2/connections",
		"https://high.csanalytics.powerplatform.microsoft.us/api/v2/connections",
		"https://dod.csanalytics.csanalytics.appsplatform.us/api/v2/connections",
	} {
		scope, err := tryGetScopeFromURL(url, urls)
		require.NoError(t, err, url)
		require.Equal(t, constants.PUBLIC_ANALYTICS_SCOPE, scope, url)
	}
}

func TestUnitTryGetScopeFromURL_HostMatchingIsCaseInsensitiveAndIgnoresPort(t *testing.T) {
	t.Parallel()

	urls := publicScopeTestUrls(t)

	scope, err := tryGetScopeFromURL("HTTPS://API.BAP.Microsoft.COM:443/providers", urls)
	require.NoError(t, err)
	require.Equal(t, constants.PUBLIC_POWERAPPS_SCOPE, scope)

	scope, err = tryGetScopeFromURL("https://Org.API.CRM4.Dynamics.com:443/api/data/v9.2/accounts", urls)
	require.NoError(t, err)
	require.Equal(t, "https://Org.API.CRM4.Dynamics.com:443/.default", scope)
}

// The URLs below each smuggle a known API domain into a part of the URL that is not its host, use plain http,
// or name a host that is not a known Power Platform or Dataverse host. None of them may get a scope, because
// the request would then carry a token to a host that should not receive it.
func TestUnitTryGetScopeFromURL_RejectsUntrustedUrls(t *testing.T) {
	t.Parallel()

	urls := publicScopeTestUrls(t)
	for _, url := range []string{
		// Known domain in the query, fragment, path, or user info of a foreign host.
		"https://evil.example/p?api.bap.microsoft.com",
		"https://evil.example/#api.bap.microsoft.com",
		"https://api.bap.microsoft.com@evil.example/",
		"https://evil.example/api.powerplatform.com/x",
		"https://evil.example/csanalytics",
		"https://evil.example/org.crm4.dynamics.com/api/data/v9.2/accounts",
		// Known domain as the prefix of a foreign host, or without a dot boundary.
		"https://api.bap.microsoft.com.evil.example/x",
		"https://evilapi.bap.microsoft.com/x",
		"https://evil-api.powerplatform.com/x",
		"https://evildynamics.com/api/data/v9.2/accounts",
		"https://csanalytics.evil.example/api/v2/connections",
		// Plain http, even to a known host.
		"http://api.bap.microsoft.com/insecure-scheme",
		"http://org.crm4.dynamics.com/api/data/v9.2/accounts",
		"http://127.0.0.1:8080/collect?x=api.bap.microsoft.com",
		// Unknown hosts no longer get a token minted for themselves.
		"https://evil.example/anything",
		"https://127.0.0.1/collect",
		// No host at all.
		"https:api.bap.microsoft.com",
		"https:///api.bap.microsoft.com",
	} {
		scope, err := tryGetScopeFromURL(url, urls)
		require.Error(t, err, url)
		require.Empty(t, scope, url)
	}
}

func TestUnitTryGetScopeFromURL_InvalidUrls(t *testing.T) {
	t.Parallel()

	urls := publicScopeTestUrls(t)
	for _, url := range []string{
		// An unknown host whose URL cannot be parsed used to panic with a nil pointer dereference.
		"https://evil.example/%zz",
		"https://api.bap.microsoft.com/%zz",
		"/relative/api.bap.microsoft.com",
		"",
	} {
		var scope string
		var err error
		require.NotPanics(t, func() { scope, err = tryGetScopeFromURL(url, urls) }, url)
		var urlErr *customerrors.UrlFormatError
		require.ErrorAs(t, err, &urlErr, url)
		require.Empty(t, scope, url)
	}
}

func TestUnitTryGetScopeFromURL_EmptyConfiguredDomainMatchesNothing(t *testing.T) {
	t.Parallel()

	scope, err := tryGetScopeFromURL("https://evil.example/anything", config.ProviderConfigUrls{PowerAppsScope: constants.PUBLIC_POWERAPPS_SCOPE})
	require.Error(t, err)
	require.Empty(t, scope)
}
