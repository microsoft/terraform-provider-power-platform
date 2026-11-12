// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package managedsolution

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/config"
	"github.com/stretchr/testify/require"
)

// The StageSolution answer of a Dataverse organization created minutes earlier that has not yet assigned
// the caller its roles, as seen in the nightly acceptance run of 2026-10-26.
const stageCallerWithoutRoles = `{
  "StageSolutionResults":{
    "StageSolutionUploadId":"upload-id",
    "StageSolutionStatus":"Failed",
    "SolutionValidationResults":[{
      "SolutionValidationResultType":"Error",
      "Message":"TryGetMaxPrivilegeDepthForUserAcrossBusinessUnits: The user with id a068a392-04d1-f111-a059-000d3a374fad has not been assigned any roles. They need a role with the prvImportCustomization privilege."
    }],
    "MissingDependencies":[]
  }
}`

// registerCallerRolesEnvironment serves environment environmentId, whose organization lives at host, and a
// StageSolution that refuses the caller for having no roles the first refusals times. It returns the count
// of StageSolution calls.
func registerCallerRolesEnvironment(t *testing.T, environmentId, host string, refusals int) *int {
	t.Helper()
	httpmock.RegisterNoResponder(func(req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("no responder found for %s %s", req.Method, req.URL)
	})
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/"+environmentId+"?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		httpmock.NewStringResponder(http.StatusOK, `{"id":"`+environmentId+`","name":"env","properties":{"linkedEnvironmentMetadata":{"instanceURL":"https://`+host+`/"}}}`))

	stageCalls := 0
	httpmock.RegisterResponder("POST", "https://"+host+"/api/data/v9.2/StageSolution",
		func(req *http.Request) (*http.Response, error) {
			stageCalls++
			if stageCalls <= refusals {
				return httpmock.NewStringResponse(http.StatusOK, stageCallerWithoutRoles), nil
			}
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
	httpmock.RegisterResponder("POST", "https://"+host+"/api/data/v9.2/ImportSolutionAsync",
		httpmock.NewStringResponder(http.StatusOK, `{"ImportJobKey":"22222222-2222-2222-2222-222222222222","AsyncOperationId":"11111111-1111-1111-1111-111111111111"}`))
	httpmock.RegisterResponder("GET", "https://"+host+"/api/data/v9.2/asyncoperations%2811111111-1111-1111-1111-111111111111%29",
		httpmock.NewStringResponder(http.StatusOK, `{"completedon":"2026-10-26T00:00:00Z"}`))
	httpmock.RegisterResponder("GET", "https://"+host+"/api/data/v9.0/RetrieveSolutionImportResult%28ImportJobId=22222222-2222-2222-2222-222222222222%29",
		httpmock.NewStringResponder(http.StatusOK, `{"SolutionOperationResult":{"Status":"Passed","ErrorMessages":[]}}`))
	httpmock.RegisterResponder("GET", "https://"+host+"/api/data/v9.2/solutions?%24expand=publisherid&%24filter=uniquename+eq+%27MetaForm%27",
		httpmock.NewStringResponder(http.StatusOK, `{"value":[{"solutionid":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","uniquename":"MetaForm","friendlyname":"Meta Form","ismanaged":true,"version":"2.0.246"}]}`))
	return &stageCalls
}

func applyCallerRolesManagedSolution(environmentId string) (string, error) {
	cfg := &config.ProviderConfig{
		TestMode: true,
		Urls:     config.ProviderConfigUrls{BapiUrl: "api.bap.microsoft.com"},
	}
	client := NewManagedSolutionClient(api.NewApiClientBase(cfg, api.NewAuthBase(cfg)))
	solution, err := client.ApplyManagedSolution(context.Background(), environmentId, []byte("managed-package"), nil, importOperationInstall, true)
	if err != nil {
		return "", err
	}
	return solution.Name, nil
}

func TestUnitApplyManagedSolution_NewOrganization_WaitsForTheCallersRoles(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const host = "managed-new-org-roles.crm.dynamics.com"
	stageCalls := registerCallerRolesEnvironment(t, "00000000-0000-0000-0000-0000000000a1", host, 2)
	api.RecordNewDataverseOrg("https://" + host + "/")

	name, err := applyCallerRolesManagedSolution("00000000-0000-0000-0000-0000000000a1")

	require.NoError(t, err, "a new organization that assigns the caller its roles within the time budget must not fail the import")
	require.Equal(t, "MetaForm", name)
	require.Equal(t, 3, *stageCalls)
}

func TestUnitApplyManagedSolution_ExistingOrganization_FailsAtOnceWhenTheCallerHasNoRoles(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	stageCalls := registerCallerRolesEnvironment(t, "00000000-0000-0000-0000-0000000000a2", "managed-existing-org-roles.crm.dynamics.com", 2)

	_, err := applyCallerRolesManagedSolution("00000000-0000-0000-0000-0000000000a2")

	require.ErrorContains(t, err, "has not been assigned any roles")
	require.Equal(t, 1, *stageCalls, "a caller without roles in an organization the provider did not just create is not waited for")
}
