// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package git_integration_test

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitSolutionGitBranchResource_Validate_Create_And_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	createdBranch := false
	rootBranchCreated := false
	updatedBranch := false
	deletedBranch := false
	deleteReadCount := 0

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001\?%24expand=permissions%2Cproperties\.capacity%2Cproperties%2FbillingPolicy(%2Cproperties%2FcopilotPolicies)?&api-version=2023-06-01$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/shared/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.2/solutions\?.*solutionid\+eq\+33333333-3333-3333-3333-333333333333.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"solutionid":"33333333-3333-3333-3333-333333333333","uniquename":"sample_solution","friendlyname":"Sample Solution","ismanaged":false,"isvisible":true,"enabledforsourcecontrolintegration":false,"version":"1.0.0.0"}]}`), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/organizations(\?.*)?$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"organizationid":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","orgdborgsettings":"<OrgSettings><SourceControlIntegrationScope>SolutionScope</SourceControlIntegrationScope></OrgSettings>"}]}`), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/gitbranches\?%24filter=.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"organizationname":"example-org","projectname":"example-project","repositoryname":"example-repo","branchname":"main"},{"organizationname":"example-org","projectname":"example-project","repositoryname":"example-repo","branchname":"develop"}]}`), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/gitrepositories\?%24filter=.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"organizationname":"example-org","projectname":"example-project","repositoryname":"example-repo","defaultbranch":"refs/heads/main"}]}`), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/sourcecontrolconfigurations%2811111111-1111-1111-1111-111111111111%29$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_sourcecontrolconfiguration_1.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/sourcecontrolbranchconfigurations",
		func(req *http.Request) (*http.Response, error) {
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}

			body := map[string]any{}
			if err := json.Unmarshal(bodyBytes, &body); err != nil {
				return nil, err
			}

			switch body["partitionid"] {
			case "33333333-3333-3333-3333-333333333333":
				createdBranch = true
			case "00000000-0000-0000-0000-000000000000":
				rootBranchCreated = true
			default:
				return httpmock.NewStringResponse(http.StatusBadRequest, ""), nil
			}

			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/sourcecontrolbranchconfigurations\?partitionId=33333333-3333-3333-3333-333333333333$`),
		func(req *http.Request) (*http.Response, error) {
			if !createdBranch {
				return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
			}

			if deletedBranch {
				deleteReadCount++
				if deleteReadCount == 1 {
					return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"sourcecontrolbranchconfigurationid":"22222222-2222-2222-2222-222222222222","partitionid":"33333333-3333-3333-3333-333333333333","statuscode":1}]}`), nil
				}

				return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
			}

			if !updatedBranch {
				return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"sourcecontrolbranchconfigurationid":"22222222-2222-2222-2222-222222222222","name":"33333333-3333-3333-3333-333333333333","partitionid":"33333333-3333-3333-3333-333333333333","branchname":"main","upstreambranchname":"main","rootfolderpath":"solutions/sample-solution","branchsyncedcommitid":"abc123","upstreambranchsyncedcommitid":"abc123","statuscode":0,"_sourcecontrolconfigurationid_value":"11111111-1111-1111-1111-111111111111"}]}`), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"sourcecontrolbranchconfigurationid":"22222222-2222-2222-2222-222222222222","name":"33333333-3333-3333-3333-333333333333","partitionid":"33333333-3333-3333-3333-333333333333","branchname":"develop","upstreambranchname":"main","rootfolderpath":"solutions/sample-solution-updated","branchsyncedcommitid":"def456","upstreambranchsyncedcommitid":"def456","statuscode":0,"_sourcecontrolconfigurationid_value":"11111111-1111-1111-1111-111111111111"}]}`), nil
		},
	)

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/sourcecontrolbranchconfigurations?partitionId=00000000-0000-0000-0000-000000000000",
		func(req *http.Request) (*http.Response, error) {
			if !rootBranchCreated {
				return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"sourcecontrolbranchconfigurationid":"44444444-4444-4444-4444-444444444444","partitionid":"00000000-0000-0000-0000-000000000000","branchname":"main","upstreambranchname":"main","rootfolderpath":"dataverse","branchsyncedcommitid":"abc123","upstreambranchsyncedcommitid":"abc123","statuscode":0,"_sourcecontrolconfigurationid_value":"11111111-1111-1111-1111-111111111111"}]}`), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/PreValidateGitComponents",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"ValidationMessages":""}`), nil
		})

	httpmock.RegisterRegexpResponder("PATCH", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/sourcecontrolbranchconfigurations%28sourcecontrolbranchconfigurationid=22222222-2222-2222-2222-222222222222,partitionid=%2733333333-3333-3333-3333-333333333333%27%29$`),
		func(req *http.Request) (*http.Response, error) {
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}

			body := map[string]any{}
			if err := json.Unmarshal(bodyBytes, &body); err != nil {
				return nil, err
			}

			if body["statuscode"] == float64(1) {
				deletedBranch = true
				return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
			}

			updatedBranch = true
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "powerplatform_solution_git_branch" "sample" {
  environment_id     = "00000000-0000-0000-0000-000000000001"
  git_integration_id = "11111111-1111-1111-1111-111111111111"
  solution_id        = "00000000-0000-0000-0000-000000000001_33333333-3333-3333-3333-333333333333"
  branch_name        = "main"
  root_folder_path   = "solutions/sample-solution"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_solution_git_branch.sample", "git_integration_id", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckResourceAttr("powerplatform_solution_git_branch.sample", "solution_id", "00000000-0000-0000-0000-000000000001_33333333-3333-3333-3333-333333333333"),
					resource.TestCheckResourceAttr("powerplatform_solution_git_branch.sample", "branch_name", "main"),
					resource.TestCheckResourceAttr("powerplatform_solution_git_branch.sample", "upstream_branch_name", "main"),
					resource.TestCheckResourceAttr("powerplatform_solution_git_branch.sample", "root_folder_path", "solutions/sample-solution"),
				),
			},
			{
				Config: `
resource "powerplatform_solution_git_branch" "sample" {
  environment_id        = "00000000-0000-0000-0000-000000000001"
  git_integration_id    = "11111111-1111-1111-1111-111111111111"
  solution_id           = "00000000-0000-0000-0000-000000000001_33333333-3333-3333-3333-333333333333"
  branch_name           = "develop"
  upstream_branch_name  = "main"
  root_folder_path      = "solutions/sample-solution-updated"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_solution_git_branch.sample", "solution_id", "00000000-0000-0000-0000-000000000001_33333333-3333-3333-3333-333333333333"),
					resource.TestCheckResourceAttr("powerplatform_solution_git_branch.sample", "branch_name", "develop"),
					resource.TestCheckResourceAttr("powerplatform_solution_git_branch.sample", "upstream_branch_name", "main"),
					resource.TestCheckResourceAttr("powerplatform_solution_git_branch.sample", "root_folder_path", "solutions/sample-solution-updated"),
				),
			},
			{
				ResourceName:      "powerplatform_solution_git_branch.sample",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "00000000-0000-0000-0000-000000000001/11111111-1111-1111-1111-111111111111/33333333-3333-3333-3333-333333333333",
			},
		},
	})
}

func TestUnitSolutionGitBranchResource_Validate_ScopeMismatch(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001\?%24expand=permissions%2Cproperties\.capacity%2Cproperties%2FbillingPolicy(%2Cproperties%2FcopilotPolicies)?&api-version=2023-06-01$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/shared/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.2/solutions\?.*solutionid\+eq\+33333333-3333-3333-3333-333333333333.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"solutionid":"33333333-3333-3333-3333-333333333333","uniquename":"sample_solution","friendlyname":"Sample Solution","ismanaged":false,"isvisible":true,"enabledforsourcecontrolintegration":false,"version":"1.0.0.0"}]}`), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/organizations(\?.*)?$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"organizationid":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","orgdborgsettings":"<OrgSettings><SourceControlIntegrationScope>EnvironmentScope</SourceControlIntegrationScope></OrgSettings>"}]}`), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/sourcecontrolconfigurations%2811111111-1111-1111-1111-111111111111%29$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_sourcecontrolconfiguration_1.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "powerplatform_solution_git_branch" "sample" {
  environment_id     = "00000000-0000-0000-0000-000000000001"
  git_integration_id = "11111111-1111-1111-1111-111111111111"
  solution_id        = "00000000-0000-0000-0000-000000000001_33333333-3333-3333-3333-333333333333"
  branch_name        = "main"
  root_folder_path   = "solutions/sample-solution"
}
`,
				ExpectError: regexp.MustCompile(`Invalid git integration scope for solution binding`),
			},
		},
	})
}

func TestUnitSolutionGitBranchResource_Validate_DuplicateBinding(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001\?%24expand=permissions%2Cproperties\.capacity%2Cproperties%2FbillingPolicy(%2Cproperties%2FcopilotPolicies)?&api-version=2023-06-01$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/shared/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.2/solutions\?.*solutionid\+eq\+33333333-3333-3333-3333-333333333333.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"solutionid":"33333333-3333-3333-3333-333333333333","uniquename":"sample_solution","friendlyname":"Sample Solution","ismanaged":false,"isvisible":true,"enabledforsourcecontrolintegration":false,"version":"1.0.0.0"}]}`), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/organizations(\?.*)?$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"organizationid":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","orgdborgsettings":"<OrgSettings><SourceControlIntegrationScope>SolutionScope</SourceControlIntegrationScope></OrgSettings>"}]}`), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/sourcecontrolconfigurations%2811111111-1111-1111-1111-111111111111%29$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_sourcecontrolconfiguration_1.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/gitbranches\?%24filter=.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"organizationname":"example-org","projectname":"example-project","repositoryname":"example-repo","branchname":"main"}]}`), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/sourcecontrolbranchconfigurations\?partitionId=33333333-3333-3333-3333-333333333333$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"sourcecontrolbranchconfigurationid":"22222222-2222-2222-2222-222222222222","name":"33333333-3333-3333-3333-333333333333","partitionid":"33333333-3333-3333-3333-333333333333","branchname":"main","upstreambranchname":"main","rootfolderpath":"solutions/sample-solution","_sourcecontrolconfigurationid_value":"11111111-1111-1111-1111-111111111111"}]}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "powerplatform_solution_git_branch" "sample" {
  environment_id     = "00000000-0000-0000-0000-000000000001"
  git_integration_id = "11111111-1111-1111-1111-111111111111"
  solution_id        = "00000000-0000-0000-0000-000000000001_33333333-3333-3333-3333-333333333333"
  branch_name        = "main"
  root_folder_path   = "solutions/sample-solution"
}
`,
				ExpectError: regexp.MustCompile(`Duplicate solution git branch binding`),
			},
		},
	})
}
