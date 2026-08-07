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
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/api/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/sourcecontrolbranchconfigurations?partitionId=33333333-3333-3333-3333-333333333333",
		func(req *http.Request) (*http.Response, error) {
			if deleted {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/shared/get_empty_value_list.json").String()), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/api/get_sourcecontrolbranchconfigurations_33333333.json").String()), nil
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
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/api/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/organizations?%24select=organizationid%2Corgdborgsettings",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/api/get_organizations_broken_scope.json").String()), nil
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
