// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

import (
	"strings"
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
	// dataverseIsWholeNamespace is set for the clouds whose Dataverse domain is their whole cloud namespace, so that a
	// look-alike of one of their API hosts falls back to a token for that look-alike host itself.
	dataverseIsWholeNamespace bool
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
			dataverseHost:             "org.crm.eaglex.ic.gov",
			dataverseIsWholeNamespace: true,
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
			dataverseHost:             "org.crm.microsoft.scloud",
			dataverseIsWholeNamespace: true,
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

	for _, cloud := range scopeTestClouds() {
		for _, rawUrl := range []string{
			"https://na.csanalytics.powerplatform.microsoft.com/api/v2/connections",
			"https://gcc.csanalytics.powerplatform.microsoft.us/api/v2/connections",
			"https://high.csanalytics.powerplatform.microsoft.us/api/v2/connections",
			"https://dod.csanalytics.csanalytics.appsplatform.us/api/v2/connections",
		} {
			scope, err := tryGetScopeFromURL(rawUrl, cloud.urls)
			require.NoError(t, err, cloud.name+" "+rawUrl)
			require.Equal(t, cloud.urls.AnalyticsScope, scope, cloud.name+" "+rawUrl)
		}
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

// untrustedUrls returns URLs built from the cloud's own domains that each smuggle a known domain into a part of the URL
// that is not its host, use plain http, or name a host that is not a known Power Platform or Dataverse host. None of
// them may get the scope of the domain they mention, because the request would then carry that domain's token to a
// host that should not receive it. lookAlike marks the hosts that imitate an API host without a dot boundary.
func untrustedUrls(cloud scopeTestCloud) []struct {
	url       string
	lookAlike bool
} {
	urls := cloud.urls
	return []struct {
		url       string
		lookAlike bool
	}{
		// Known domain in the query, fragment, path, or user info of a foreign host.
		{url: "https://evil.example/p?" + urls.BapiUrl},
		{url: "https://evil.example/#" + urls.BapiUrl},
		{url: "https://" + urls.BapiUrl + "@evil.example/"},
		{url: "https://evil.example/?" + urls.PowerAppsUrl},
		{url: "https://evil.example/" + urls.PowerPlatformUrl + "/x"},
		{url: "https://evil.example/" + urls.PowerAppsAdvisor + "/api/rule"},
		{url: "https://evil.example/" + urls.AdminPowerPlatformUrl + "/api/tenant"},
		{url: "https://evil.example/csanalytics"},
		{url: "https://evil.example/" + cloud.dataverseHost + "/api/data/v9.2/accounts"},
		// Known domain as the prefix of a foreign host, or without a dot boundary.
		{url: "https://" + urls.BapiUrl + ".evil.example/x"},
		{url: "https://" + urls.PowerPlatformUrl + ".evil.example/x"},
		{url: "https://" + cloud.dataverseHost + ".evil.example/api/data/v9.2/accounts"},
		{url: "https://evil" + urls.BapiUrl + "/x", lookAlike: true},
		{url: "https://evil-" + urls.PowerPlatformUrl + "/x", lookAlike: true},
		{url: "https://evildynamics.com/api/data/v9.2/accounts"},
		{url: "https://csanalytics.evil.example/api/v2/connections"},
		// Plain http, even to a known host.
		{url: "http://" + urls.BapiUrl + "/insecure-scheme"},
		{url: "http://" + urls.PowerPlatformUrl + "/insecure-scheme"},
		{url: "http://" + cloud.dataverseHost + "/api/data/v9.2/accounts"},
		{url: "http://127.0.0.1:8080/collect?x=" + urls.BapiUrl},
		// Unknown hosts no longer get a token minted for themselves.
		{url: "https://evil.example/anything"},
		{url: "https://127.0.0.1/collect"},
		// No host at all.
		{url: "https:" + urls.BapiUrl},
		{url: "https:///" + urls.BapiUrl},
	}
}

func TestUnitTryGetScopeFromURL_RejectsUntrustedUrls(t *testing.T) {
	t.Parallel()

	for _, cloud := range scopeTestClouds() {
		t.Run(cloud.name, func(t *testing.T) {
			t.Parallel()
			for _, tc := range untrustedUrls(cloud) {
				scope, err := tryGetScopeFromURL(tc.url, cloud.urls)
				if tc.lookAlike && cloud.dataverseIsWholeNamespace {
					// The look-alike host is inside this cloud's namespace, so it falls back to a token for itself, never
					// to the token of the API host it imitates.
					require.NoError(t, err, tc.url)
					require.Equal(t, strings.TrimSuffix(tc.url, "/x")+"/.default", scope, tc.url)
					continue
				}
				require.Error(t, err, tc.url)
				require.Empty(t, scope, tc.url)
			}
		})
	}
}

func TestUnitTryGetScopeFromURL_InvalidUrls(t *testing.T) {
	t.Parallel()

	for _, cloud := range scopeTestClouds() {
		t.Run(cloud.name, func(t *testing.T) {
			t.Parallel()
			for _, rawUrl := range []string{
				// An unknown host whose URL cannot be parsed used to panic with a nil pointer dereference.
				"https://evil.example/%zz",
				"https://" + cloud.urls.BapiUrl + "/%zz",
				"https://" + cloud.dataverseHost + "/%zz",
				"/relative/" + cloud.urls.BapiUrl,
				"",
			} {
				var scope string
				var err error
				require.NotPanics(t, func() { scope, err = tryGetScopeFromURL(rawUrl, cloud.urls) }, rawUrl)
				var urlErr *customerrors.UrlFormatError
				require.ErrorAs(t, err, &urlErr, rawUrl)
				require.Empty(t, scope, rawUrl)
			}
		})
	}
}

func TestUnitTryGetScopeFromURL_EmptyConfiguredDomainMatchesNothing(t *testing.T) {
	t.Parallel()

	scope, err := tryGetScopeFromURL("https://evil.example/anything", config.ProviderConfigUrls{PowerAppsScope: constants.PUBLIC_POWERAPPS_SCOPE})
	require.Error(t, err)
	require.Empty(t, scope)
}
