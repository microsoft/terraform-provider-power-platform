// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package solution

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/config"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
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
    "SolutionDetails":{"SolutionUniqueName":"TerraformTestSolution"}
  }
}`), nil
		})
	httpmock.RegisterResponder("POST", "https://"+host+"/api/data/v9.2/ImportSolutionAsync",
		httpmock.NewStringResponder(http.StatusOK, `{"ImportJobKey":"22222222-2222-2222-2222-222222222222","AsyncOperationId":"11111111-1111-1111-1111-111111111111"}`))
	httpmock.RegisterResponder("GET", "https://"+host+"/api/data/v9.2/asyncoperations%2811111111-1111-1111-1111-111111111111%29",
		httpmock.NewStringResponder(http.StatusOK, `{"completedon":"2026-10-26T00:00:00Z"}`))
	httpmock.RegisterResponder("GET", "https://"+host+"/api/data/v9.0/RetrieveSolutionImportResult%28ImportJobId=22222222-2222-2222-2222-222222222222%29",
		httpmock.NewStringResponder(http.StatusOK, `{"SolutionOperationResult":{"Status":"Passed","ErrorMessages":[]}}`))
	httpmock.RegisterResponder("GET", "https://"+host+"/api/data/v9.2/solutions?%24expand=publisherid&%24filter=uniquename+eq+%27TerraformTestSolution%27",
		httpmock.NewStringResponder(http.StatusOK, `{"value":[{"solutionid":"86928ed8-df37-4ce2-add5-47030a833bff","uniquename":"TerraformTestSolution","friendlyname":"Terraform Test Solution","ismanaged":false,"version":"1.1.0.0"}]}`))
	return &stageCalls
}

func createCallerRolesSolution(environmentId string) (*SolutionDto, error) {
	cfg := &config.ProviderConfig{
		TestMode: true,
		Urls:     config.ProviderConfigUrls{BapiUrl: "api.bap.microsoft.com"},
	}
	client := NewSolutionClient(api.NewApiClientBase(cfg, api.NewAuthBase(cfg)))
	return client.CreateSolution(context.Background(), environmentId, []byte("solution"), nil)
}

func TestUnitCreateSolution_NewOrganization_WaitsForTheCallersRoles(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const host = "solution-new-org-roles.crm.dynamics.com"
	stageCalls := registerCallerRolesEnvironment(t, "00000000-0000-0000-0000-0000000000b1", host, 2)
	api.RecordNewDataverseOrg("https://" + host + "/")

	solution, err := createCallerRolesSolution("00000000-0000-0000-0000-0000000000b1")

	require.NoError(t, err, "a new organization that assigns the caller its roles within the time budget must not fail the import")
	require.Equal(t, "TerraformTestSolution", solution.Name)
	require.Equal(t, 3, *stageCalls)
}

func TestUnitCreateSolution_NewOrganization_StopsWaitingForTheCallersRolesAtTheTimeBudget(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const host = "solution-new-org-no-roles.crm.dynamics.com"
	stageCalls := registerCallerRolesEnvironment(t, "00000000-0000-0000-0000-0000000000b2", host, 1000)
	api.RecordNewDataverseOrg("https://" + host + "/")

	_, err := createCallerRolesSolution("00000000-0000-0000-0000-0000000000b2")

	replays := int(constants.NEW_DATAVERSE_ORG_CALLER_ROLES_RETRY_BUDGET / constants.NEW_DATAVERSE_ORG_CALLER_ROLES_RETRY_INTERVAL)
	require.ErrorContains(t, err, fmt.Sprintf("the caller still had no security roles after %d retries", replays))
	require.ErrorContains(t, err, "has not been assigned any roles")
	require.Equal(t, replays+1, *stageCalls)
}

func TestUnitCreateSolution_ExistingOrganization_FailsAtOnceWhenTheCallerHasNoRoles(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	stageCalls := registerCallerRolesEnvironment(t, "00000000-0000-0000-0000-0000000000b3", "solution-existing-org-roles.crm.dynamics.com", 2)

	_, err := createCallerRolesSolution("00000000-0000-0000-0000-0000000000b3")

	require.ErrorContains(t, err, "has not been assigned any roles")
	require.Equal(t, 1, *stageCalls, "a caller without roles in an organization the provider did not just create is not waited for")
}
