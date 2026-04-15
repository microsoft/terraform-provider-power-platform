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
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"organizationname":"example-org"}]}`), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/organizations(\?.*)?$`),
		func(req *http.Request) (*http.Response, error) {
			if !updatedConfiguration {
				return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"organizationid":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","orgdborgsettings":"<OrgSettings><SourceControlIntegrationScope>SolutionScope</SourceControlIntegrationScope></OrgSettings>"}]}`), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"organizationid":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","orgdborgsettings":"<OrgSettings><SourceControlIntegrationScope>EnvironmentScope</SourceControlIntegrationScope></OrgSettings>"}]}`), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/gitprojects\?%24filter=%28organizationname\+eq\+%27example-org%27%29$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"organizationname":"example-org","projectname":"example-project"}]}`), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/gitrepositories\?%24filter=.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"organizationname":"example-org","projectname":"example-project","repositoryname":"example-repo","defaultbranch":"refs/heads/main"},{"organizationname":"example-org","projectname":"example-project","repositoryname":"example-repo-updated","defaultbranch":"refs/heads/main"}]}`), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.2/solutions\?.*ismanaged\+eq\+false.*isvisible\+eq\+true.*enabledforsourcecontrolintegration.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{"value":[
{"solutionid":"00000001-0000-0000-0001-00000000009b","uniquename":"common-default","friendlyname":"Common Data Services Default Solution","ismanaged":false,"isvisible":true,"enabledforsourcecontrolintegration":false},
{"solutionid":"fd140aaf-4df4-11dd-bd17-0019b9312238","uniquename":"Default","friendlyname":"Default Solution","ismanaged":false,"isvisible":true,"enabledforsourcecontrolintegration":false},
{"solutionid":"33333333-3333-3333-3333-333333333333","uniquename":"solution-one","friendlyname":"solution-one","ismanaged":false,"isvisible":true,"enabledforsourcecontrolintegration":%t},
{"solutionid":"44444444-4444-4444-4444-444444444444","uniquename":"solution-two","friendlyname":"solution-two","ismanaged":false,"isvisible":true,"enabledforsourcecontrolintegration":%t}
]}`, environmentScopeEnabled["33333333-3333-3333-3333-333333333333"], environmentScopeEnabled["44444444-4444-4444-4444-444444444444"])), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.2/solutions\?.*solutionid\+eq\+33333333-3333-3333-3333-333333333333.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{"value":[{"solutionid":"33333333-3333-3333-3333-333333333333","uniquename":"solution-one","friendlyname":"solution-one","ismanaged":false,"isvisible":true,"enabledforsourcecontrolintegration":%t,"version":"1.0.0.0"}]}`, environmentScopeEnabled["33333333-3333-3333-3333-333333333333"])), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.2/solutions\?.*solutionid\+eq\+44444444-4444-4444-4444-444444444444.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{"value":[{"solutionid":"44444444-4444-4444-4444-444444444444","uniquename":"solution-two","friendlyname":"solution-two","ismanaged":false,"isvisible":true,"enabledforsourcecontrolintegration":%t,"version":"1.0.0.0"}]}`, environmentScopeEnabled["44444444-4444-4444-4444-444444444444"])), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/sourcecontrolconfigurations",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/sourcecontrolconfigurations",
		func(req *http.Request) (*http.Response, error) {
			if configurationImplicitlyDeleted {
				return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
			}

			if !updatedConfiguration {
				return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"sourcecontrolconfigurationid":"11111111-1111-1111-1111-111111111111","organizationname":"example-org","projectname":"example-project","repositoryname":"example-repo","gitprovider":0}]}`), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"sourcecontrolconfigurationid":"11111111-1111-1111-1111-111111111111","organizationname":"example-org","projectname":"example-project","repositoryname":"example-repo-updated","gitprovider":0}]}`), nil
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
				return httpmock.NewStringResponse(http.StatusNotFound, `{"error":{"code":"0x80040217","message":"source control configuration not found"}}`), nil
			}

			if !updatedConfiguration {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/environment_git_integration/get_sourcecontrolconfiguration_1.json").String()), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, `{
				"sourcecontrolconfigurationid":"11111111-1111-1111-1111-111111111111",
				"name":"",
				"organizationname":"example-org",
				"projectname":"example-project",
				"repositoryname":"example-repo-updated",
				"gitprovider":0
			}`), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/sourcecontrolbranchconfigurations?partitionId=00000000-0000-0000-0000-000000000000",
		func(req *http.Request) (*http.Response, error) {
			if !rootBranchCreated {
				return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"sourcecontrolbranchconfigurationid":"22222222-2222-2222-2222-222222222222","partitionid":"00000000-0000-0000-0000-000000000000","branchname":"main","upstreambranchname":"main","rootfolderpath":"dataverse","branchsyncedcommitid":"abc123","upstreambranchsyncedcommitid":"abc123","statuscode":0,"_sourcecontrolconfigurationid_value":"11111111-1111-1111-1111-111111111111"}]}`), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/PreValidateGitComponents",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"ValidationMessages":""}`), nil
		})

	httpmock.RegisterRegexpResponder("PATCH", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/sourcecontrolbranchconfigurations%28sourcecontrolbranchconfigurationid=22222222-2222-2222-2222-222222222222,partitionid=%2700000000-0000-0000-0000-000000000000%27%29$`),
		func(req *http.Request) (*http.Response, error) {
			rootBranchCreated = false
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/sourcecontrolconfigurations%2811111111-1111-1111-1111-111111111111%29$`),
		func(req *http.Request) (*http.Response, error) {
			configurationImplicitlyDeleted = true
			return httpmock.NewStringResponse(http.StatusBadRequest, `{"error":{"code":"0x80040265","message":"Existing source control configurations can't be deleted."}}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "powerplatform_environment_git_integration" "test" {
  environment_id    = "00000000-0000-0000-0000-000000000001"
  git_provider      = "AzureDevOps"
  scope             = "Solution"
  organization_name = "example-org"
  project_name      = "example-project"
  repository_name   = "example-repo"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_git_integration.test", "git_provider", "AzureDevOps"),
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
  git_provider      = "AzureDevOps"
  scope             = "Environment"
  organization_name = "example-org"
  project_name      = "example-project"
  repository_name   = "example-repo-updated"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_git_integration.test", "git_provider", "AzureDevOps"),
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
