// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package git_integration_test

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitEnvironmentGitIntegrationResource_Validate_Create_And_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	updatedConfiguration := false
	rootBranchCreated := false
	configurationImplicitlyDeleted := false
	environmentScopeSolutionPatches := 0
	environmentScopeEnabled := map[string]bool{
		"33333333-3333-3333-3333-333333333333": false,
		"44444444-4444-4444-4444-444444444444": false,
	}

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001\?%24expand=permissions%2Cproperties\.capacity%2Cproperties%2FbillingPolicy(%2Cproperties%2FcopilotPolicies)?&api-version=2023-06-01$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/shared/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/gitorganizations",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_gitorganizations.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/organizations(\?.*)?$`),
		func(req *http.Request) (*http.Response, error) {
			if !updatedConfiguration {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/shared/get_organizations_solution_scope.json").String()), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/shared/get_organizations_environment_scope.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/gitprojects\?%24filter=%28organizationname\+eq\+%27example-org%27%29$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_gitprojects.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/gitrepositories\?%24filter=.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_gitrepositories_with_updated.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.2/solutions\?.*ismanaged\+eq\+false.*isvisible\+eq\+true.*enabledforsourcecontrolintegration.*$`),
		func(req *http.Request) (*http.Response, error) {
			// Solutions are enabled sequentially in list order, so the payload
			// progresses from none enabled, to solution-one enabled, to all enabled.
			switch {
			case environmentScopeEnabled["33333333-3333-3333-3333-333333333333"] && environmentScopeEnabled["44444444-4444-4444-4444-444444444444"]:
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_solutions_3.json").String()), nil
			case environmentScopeEnabled["33333333-3333-3333-3333-333333333333"]:
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_solutions_2.json").String()), nil
			default:
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_solutions_1.json").String()), nil
			}
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.2/solutions\?.*solutionid\+eq\+33333333-3333-3333-3333-333333333333.*$`),
		func(req *http.Request) (*http.Response, error) {
			if !environmentScopeEnabled["33333333-3333-3333-3333-333333333333"] {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_solution_33333333_1.json").String()), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_solution_33333333_2.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.2/solutions\?.*solutionid\+eq\+44444444-4444-4444-4444-444444444444.*$`),
		func(req *http.Request) (*http.Response, error) {
			if !environmentScopeEnabled["44444444-4444-4444-4444-444444444444"] {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_solution_44444444_1.json").String()), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_solution_44444444_2.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/sourcecontrolconfigurations",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/sourcecontrolconfigurations",
		func(req *http.Request) (*http.Response, error) {
			if configurationImplicitlyDeleted {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/shared/get_empty_value_list.json").String()), nil
			}

			if !updatedConfiguration {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_sourcecontrolconfigurations_1.json").String()), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_sourcecontrolconfigurations_2.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/sourcecontrolbranchconfigurations",
		func(req *http.Request) (*http.Response, error) {
			rootBranchCreated = true
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterRegexpResponder("PATCH", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/sourcecontrolconfigurations%2811111111-1111-1111-1111-111111111111%29$`),
		func(req *http.Request) (*http.Response, error) {
			updatedConfiguration = true
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterRegexpResponder("PATCH", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/organizations%28aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa%29$`),
		func(req *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			bodyText := string(body)
			if strings.Contains(bodyText, "organizationid") {
				return nil, fmt.Errorf("organization scope patch unexpectedly included organizationid: %s", bodyText)
			}
			if !strings.Contains(bodyText, "SourceControlIntegrationScope") {
				return nil, fmt.Errorf("organization scope patch missing SourceControlIntegrationScope: %s", bodyText)
			}
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterRegexpResponder("PATCH", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/solutions(?:%28|\()(33333333-3333-3333-3333-333333333333|44444444-4444-4444-4444-444444444444)(?:%29|\))$`),
		func(req *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			bodyText := string(body)
			if !strings.Contains(bodyText, `"enabledforsourcecontrolintegration":true`) {
				return nil, fmt.Errorf("solution enablement patch did not send boolean true: %s", bodyText)
			}
			matched := false
			for _, id := range []string{
				"33333333-3333-3333-3333-333333333333",
				"44444444-4444-4444-4444-444444444444",
			} {
				if strings.Contains(req.URL.String(), id) {
					environmentScopeEnabled[id] = true
					environmentScopeSolutionPatches++
					matched = true
					break
				}
			}
			if !matched {
				return nil, fmt.Errorf("unexpected solution enablement URL: %s", req.URL.String())
			}
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/sourcecontrolconfigurations%28[0-9a-f-]{36}%29$`),
		func(req *http.Request) (*http.Response, error) {
			if configurationImplicitlyDeleted {
				return httpmock.NewStringResponse(http.StatusNotFound, httpmock.File("tests/resource/environment_git_integration/error_sourcecontrolconfiguration_not_found.json").String()), nil
			}

			if !updatedConfiguration {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_sourcecontrolconfiguration_1.json").String()), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_sourcecontrolconfiguration_2.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/sourcecontrolbranchconfigurations?partitionId=00000000-0000-0000-0000-000000000000",
		func(req *http.Request) (*http.Response, error) {
			if !rootBranchCreated {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/shared/get_empty_value_list.json").String()), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_sourcecontrolbranchconfigurations_root.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/PreValidateGitComponents",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/shared/post_prevalidategitcomponents.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("PATCH", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/sourcecontrolbranchconfigurations%28sourcecontrolbranchconfigurationid=22222222-2222-2222-2222-222222222222,partitionid=%2700000000-0000-0000-0000-000000000000%27%29$`),
		func(req *http.Request) (*http.Response, error) {
			rootBranchCreated = false
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/sourcecontrolconfigurations%2811111111-1111-1111-1111-111111111111%29$`),
		func(req *http.Request) (*http.Response, error) {
			configurationImplicitlyDeleted = true
			return httpmock.NewStringResponse(http.StatusBadRequest, httpmock.File("tests/resource/environment_git_integration/error_sourcecontrolconfiguration_delete_conflict.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
variable "environment_id" {
  type    = string
  default = "00000000-0000-0000-0000-000000000001"
}

variable "organization_name" {
  type    = string
  default = "example-org"
}

variable "project_name" {
  type    = string
  default = "example-project"
}

variable "repository_name" {
  type    = string
  default = "example-repo"
}

resource "powerplatform_environment_git_integration" "test" {
  environment_id    = var.environment_id
  scope             = "Solution"
  organization_name = var.organization_name
  project_name      = var.project_name
  repository_name   = var.repository_name
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_git_integration.test", "scope", "Solution"),
					resource.TestCheckResourceAttr("powerplatform_environment_git_integration.test", "organization_name", "example-org"),
					resource.TestCheckResourceAttr("powerplatform_environment_git_integration.test", "project_name", "example-project"),
					resource.TestCheckResourceAttr("powerplatform_environment_git_integration.test", "repository_name", "example-repo"),
				),
			},
			{
				Config: `
resource "powerplatform_environment_git_integration" "test" {
  environment_id    = "00000000-0000-0000-0000-000000000001"
  scope             = "Environment"
  organization_name = "example-org"
  project_name      = "example-project"
  repository_name   = "example-repo-updated"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_git_integration.test", "scope", "Environment"),
					resource.TestCheckResourceAttr("powerplatform_environment_git_integration.test", "organization_name", "example-org"),
					resource.TestCheckResourceAttr("powerplatform_environment_git_integration.test", "project_name", "example-project"),
					resource.TestCheckResourceAttr("powerplatform_environment_git_integration.test", "repository_name", "example-repo-updated"),
					func(_ *terraform.State) error {
						if environmentScopeSolutionPatches != 2 {
							return fmt.Errorf("expected 2 environment-scope solution enablement patches, got %d", environmentScopeSolutionPatches)
						}
						return nil
					},
				),
			},
			{
				ResourceName:      "powerplatform_environment_git_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "00000000-0000-0000-0000-000000000001",
			},
		},
	})
}

func TestUnitEnvironmentGitIntegrationResource_Validate_Delete_When_Parent_Environment_Is_Deleted(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	environmentDeleted := false

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001\?%24expand=permissions%2Cproperties\.capacity%2Cproperties%2FbillingPolicy(%2Cproperties%2FcopilotPolicies)?&api-version=2023-06-01$`),
		func(req *http.Request) (*http.Response, error) {
			if environmentDeleted {
				return httpmock.NewStringResponse(http.StatusNotFound, httpmock.File("tests/shared/error_environment_not_found.json").String()), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/shared/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/gitorganizations",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_gitorganizations.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/organizations(\?.*)?$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/shared/get_organizations_solution_scope.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/gitprojects\?%24filter=%28organizationname\+eq\+%27example-org%27%29$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_gitprojects.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/gitrepositories\?%24filter=.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/shared/get_gitrepositories.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/sourcecontrolconfigurations",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/sourcecontrolconfigurations",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_sourcecontrolconfigurations_1.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/sourcecontrolconfigurations%28[0-9a-f-]{36}%29$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_sourcecontrolconfiguration_1.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("PATCH", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/organizations%28aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa%29$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "powerplatform_environment_git_integration" "test" {
  environment_id    = "00000000-0000-0000-0000-000000000001"
  scope             = "Solution"
  organization_name = "example-org"
  project_name      = "example-project"
  repository_name   = "example-repo"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_git_integration.test", "repository_name", "example-repo"),
				),
			},
			{
				// Simulate the parent environment being deleted out-of-band. The refresh plan
				// removes the resource from the in-memory state (Read tolerates the missing
				// environment), and the automatic post-test destroy (which runs with refresh
				// disabled) exercises Delete against the deleted environment, which must
				// succeed silently instead of erroring.
				PreConfig: func() {
					environmentDeleted = true
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				Config: `
resource "powerplatform_environment_git_integration" "test" {
  environment_id    = "00000000-0000-0000-0000-000000000001"
  scope             = "Solution"
  organization_name = "example-org"
  project_name      = "example-project"
  repository_name   = "example-repo"
}
`,
			},
		},
	})
}
