// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package git_integration

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/config"
	"github.com/stretchr/testify/require"
)

func TestUnitDeleteSolutionGitBranch_UsesLookedUpBranchID(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterNoResponder(func(req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("no responder found for %s %s", req.Method, req.URL)
	})

	deleted := false

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001\?%24expand=permissions%2Cproperties\.capacity%2Cproperties%2FbillingPolicy(%2Cproperties%2FcopilotPolicies)?&api-version=2023-06-01$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"id": "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001",
				"name": "00000000-0000-0000-0000-000000000001",
				"type": "Microsoft.BusinessAppPlatform/scopes/admin/environments",
				"location": "europe",
				"properties": {
					"displayName": "Test",
					"linkedEnvironmentMetadata": {
						"instanceUrl": "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/"
					}
				}
			}`), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/sourcecontrolbranchconfigurations?partitionId=33333333-3333-3333-3333-333333333333",
		func(req *http.Request) (*http.Response, error) {
			if deleted {
				return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"sourcecontrolbranchconfigurationid":"22222222-2222-2222-2222-222222222222","partitionid":"33333333-3333-3333-3333-333333333333","statuscode":0,"_sourcecontrolconfigurationid_value":"11111111-1111-1111-1111-111111111111"}]}`), nil
		})

	httpmock.RegisterRegexpResponder("PATCH", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/sourcecontrolbranchconfigurations%28sourcecontrolbranchconfigurationid=22222222-2222-2222-2222-222222222222,partitionid=%2733333333-3333-3333-3333-333333333333%27%29$`),
		func(req *http.Request) (*http.Response, error) {
			deleted = true
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	cfg := &config.ProviderConfig{
		TestMode: true,
		Urls: config.ProviderConfigUrls{
			BapiUrl: "api.bap.microsoft.com",
		},
	}
	apiClient := api.NewApiClientBase(cfg, api.NewAuthBase(cfg))
	client := newGitIntegrationClient(apiClient)

	err := client.DeleteSolutionGitBranch(context.Background(), "00000000-0000-0000-0000-000000000001", "11111111-1111-1111-1111-111111111111", "33333333-3333-3333-3333-333333333333")
	require.NoError(t, err)
	require.True(t, deleted)
}

func TestUnitGetSourceControlIntegrationScope_RejectsUnknownScope(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterNoResponder(func(req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("no responder found for %s %s", req.Method, req.URL)
	})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001\?%24expand=permissions%2Cproperties\.capacity%2Cproperties%2FbillingPolicy(%2Cproperties%2FcopilotPolicies)?&api-version=2023-06-01$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"id": "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001",
				"name": "00000000-0000-0000-0000-000000000001",
				"type": "Microsoft.BusinessAppPlatform/scopes/admin/environments",
				"location": "europe",
				"properties": {
					"displayName": "Test",
					"linkedEnvironmentMetadata": {
						"instanceUrl": "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/"
					}
				}
			}`), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/organizations?%24select=organizationid%2Corgdborgsettings",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"organizationid":"44444444-4444-4444-4444-444444444444","orgdborgsettings":"<OrgSettings><SourceControlIntegrationScope>BrokenScope</SourceControlIntegrationScope></OrgSettings>"}]}`), nil
		})

	cfg := &config.ProviderConfig{
		TestMode: true,
		Urls: config.ProviderConfigUrls{
			BapiUrl: "api.bap.microsoft.com",
		},
	}
	apiClient := api.NewApiClientBase(cfg, api.NewAuthBase(cfg))
	client := newGitIntegrationClient(apiClient)

	_, err := client.GetSourceControlIntegrationScope(context.Background(), "00000000-0000-0000-0000-000000000001")
	require.ErrorContains(t, err, "could not be determined")
}
