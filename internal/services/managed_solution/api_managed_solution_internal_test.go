// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package managedsolution

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/config"
	"github.com/stretchr/testify/require"
)

func TestUnitManagedSolutionImportDecisionSeparatesPackageUpgradeFromTargetRebinding(t *testing.T) {
	tests := []struct {
		name                string
		installed           string
		configured          string
		allowAdoption       bool
		hasTargetBindings   bool
		wantAdoption        bool
		wantStageAndUpgrade bool
		wantError           bool
	}{
		{
			name:          "create adopts an exact installed package",
			installed:     "2.0.273.0",
			configured:    "2.0.273",
			allowAdoption: true,
			wantAdoption:  true,
		},
		{
			name:              "create reimports exact package when target bindings must be applied",
			installed:         "2.0.273.0",
			configured:        "2.0.273",
			allowAdoption:     true,
			hasTargetBindings: true,
		},
		{
			name:       "same package version reimports target bindings without upgrade",
			installed:  "2.0.273",
			configured: "2.0.273.0",
		},
		{
			name:                "higher package version stages and upgrades",
			installed:           "2.0.273",
			configured:          "2.0.274",
			wantStageAndUpgrade: true,
		},
		{
			name:       "lower package version fails before import",
			installed:  "2.0.273",
			configured: "2.0.1",
			wantError:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adopt, operation, err := managedSolutionImportDecision(
				test.installed,
				test.configured,
				test.allowAdoption,
				test.hasTargetBindings)
			if test.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.wantAdoption, adopt)
			wantOperation := importOperationInstall
			if test.wantStageAndUpgrade {
				wantOperation = importOperationStageAndUpgrade
			}
			require.Equal(t, wantOperation, operation)
		})
	}
}

func TestUnitApplyManagedSolution_UsesStageAndUpgradeForUpdates(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterNoResponder(func(req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("no responder found for %s %s", req.Method, req.URL)
	})
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
  "id":"00000000-0000-0000-0000-000000000001",
  "name":"env",
  "properties":{"linkedEnvironmentMetadata":{"instanceURL":"https://example.crm.dynamics.com/"}}
}`), nil
		})
	httpmock.RegisterResponder("POST", "https://example.crm.dynamics.com/api/data/v9.2/StageSolution",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
  "StageSolutionResults":{
    "StageSolutionUploadId":"upload-id",
    "StageSolutionStatus":"Passed",
    "SolutionValidationResults":[],
    "MissingDependencies":[],
    "SolutionDetails":{"SolutionUniqueName":"MetaForm"}
  }
}`), nil
		})
	stageAndUpgradeCalled := false
	var importRequest importSolutionDto
	httpmock.RegisterResponder("POST", "https://example.crm.dynamics.com/api/data/v9.2/StageAndUpgradeAsync",
		func(req *http.Request) (*http.Response, error) {
			stageAndUpgradeCalled = true
			body, err := io.ReadAll(req.Body)
			require.NoError(t, err)
			require.NoError(t, json.Unmarshal(body, &importRequest))
			return httpmock.NewStringResponse(http.StatusOK, `{"ImportJobKey":"22222222-2222-2222-2222-222222222222","AsyncOperationId":"11111111-1111-1111-1111-111111111111"}`), nil
		})
	httpmock.RegisterResponder("GET", "https://example.crm.dynamics.com/api/data/v9.2/asyncoperations%2811111111-1111-1111-1111-111111111111%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"completedon":"2026-07-13T00:00:00Z"}`), nil
		})
	httpmock.RegisterResponder("GET", "https://example.crm.dynamics.com/api/data/v9.0/RetrieveSolutionImportResult%28ImportJobId=22222222-2222-2222-2222-222222222222%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"SolutionOperationResult":{"Status":"Passed","ErrorMessages":[]}}`), nil
		})
	httpmock.RegisterResponder("GET", "https://example.crm.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=uniquename+eq+%27MetaForm%27",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"solutionid":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","uniquename":"MetaForm","friendlyname":"Meta Form","ismanaged":true,"version":"2.0.246"}]}`), nil
		})

	cfg := &config.ProviderConfig{
		TestMode: true,
		Urls:     config.ProviderConfigUrls{BapiUrl: "api.bap.microsoft.com"},
	}
	client := NewManagedSolutionClient(api.NewApiClientBase(cfg, api.NewAuthBase(cfg)))
	solution, err := client.ApplyManagedSolution(
		context.Background(),
		"00000000-0000-0000-0000-000000000001",
		[]byte("managed-package"),
		nil,
		importOperationStageAndUpgrade,
		true)

	require.NoError(t, err)
	require.True(t, stageAndUpgradeCalled)
	require.True(t, importRequest.PublishWorkflows)
	require.True(t, importRequest.OverwriteUnmanagedCustomizations)
	require.True(t, importRequest.SkipProductUpdateDependencies)
	require.Equal(t, "MetaForm", solution.Name)
}

func TestUnitApplyManagedSolution_DoesNotReplayAmbiguousImportStart(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerManagedSolutionApiEnvironment(t)
	registerManagedSolutionStage(t)
	attempts := 0
	httpmock.RegisterResponder("POST", "https://example.crm.dynamics.com/api/data/v9.2/ImportSolutionAsync",
		func(req *http.Request) (*http.Response, error) {
			attempts++
			return httpmock.NewStringResponse(http.StatusServiceUnavailable, `{"error":"ambiguous"}`), nil
		})

	cfg := &config.ProviderConfig{
		TestMode: true,
		Urls:     config.ProviderConfigUrls{BapiUrl: "api.bap.microsoft.com"},
	}
	client := NewManagedSolutionClient(api.NewApiClientBase(cfg, api.NewAuthBase(cfg)))
	_, err := client.ApplyManagedSolution(
		context.Background(),
		"00000000-0000-0000-0000-000000000001",
		[]byte("managed-package"),
		nil,
		importOperationInstall,
		true)

	require.Error(t, err)
	require.Equal(t, 1, attempts, "an ambiguous ImportSolutionAsync response must not replay the import")
}

func TestUnitPublishAllCustomizations_UsesDataversePublishAllXml(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerManagedSolutionApiEnvironment(t)
	called := 0
	httpmock.RegisterResponder("POST", "https://example.crm.dynamics.com/api/data/v9.2/PublishAllXml",
		func(req *http.Request) (*http.Response, error) {
			called++
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	cfg := &config.ProviderConfig{
		TestMode: true,
		Urls:     config.ProviderConfigUrls{BapiUrl: "api.bap.microsoft.com"},
	}
	client := NewManagedSolutionClient(api.NewApiClientBase(cfg, api.NewAuthBase(cfg)))
	err := client.PublishAllCustomizations(context.Background(), "00000000-0000-0000-0000-000000000001")

	require.NoError(t, err)
	require.Equal(t, 1, called)
}

func registerManagedSolutionApiEnvironment(t *testing.T) {
	t.Helper()
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
  "id":"00000000-0000-0000-0000-000000000001",
  "name":"env",
  "properties":{"linkedEnvironmentMetadata":{"instanceURL":"https://example.crm.dynamics.com/"}}
}`), nil
		})
}

func registerManagedSolutionStage(t *testing.T) {
	t.Helper()
	httpmock.RegisterResponder("POST", "https://example.crm.dynamics.com/api/data/v9.2/StageSolution",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
  "StageSolutionResults":{
    "StageSolutionUploadId":"upload-id",
    "StageSolutionStatus":"Passed",
    "SolutionValidationResults":[],
    "MissingDependencies":[],
    "SolutionDetails":{"SolutionUniqueName":"MetaForm"}
  }
}`), nil
		})
}
