// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package solution_test

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

const (
	unmanagedSolutionEnvironmentID = "00000000-0000-0000-0000-000000000001"
	unmanagedSolutionID            = "86928ed8-df37-4ce2-add5-47030a833bff"
	unmanagedPublisherID           = "aa47dc6c-bf13-490b-a007-1da95a0d1e3f"
)

func TestUnitUnmanagedSolutionResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/"+unmanagedSolutionEnvironmentID+"?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Unmanaged_Create/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/solutions",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusCreated, "")
			resp.Header.Set("OData-EntityId", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/solutions("+unmanagedSolutionID+")")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=solutionid+eq+"+unmanagedSolutionID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Unmanaged_Create/get_solution.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.2/solutions%28`+regexp.QuoteMeta(unmanagedSolutionID)+`%29$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `

				resource "powerplatform_unmanaged_solution" "solution" {
					environment_id = "` + unmanagedSolutionEnvironmentID + `"
					uniquename     = "TerraformTestSolution"
					display_name   = "Terraform Test Solution"
					publisher_id   = "` + unmanagedPublisherID + `"
					description    = "Created by Terraform"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_unmanaged_solution.solution", "id", unmanagedSolutionID),
					resource.TestCheckResourceAttr("powerplatform_unmanaged_solution.solution", "environment_id", unmanagedSolutionEnvironmentID),
					resource.TestCheckResourceAttr("powerplatform_unmanaged_solution.solution", "uniquename", "TerraformTestSolution"),
					resource.TestCheckResourceAttr("powerplatform_unmanaged_solution.solution", "display_name", "Terraform Test Solution"),
					resource.TestCheckResourceAttr("powerplatform_unmanaged_solution.solution", "publisher_id", unmanagedPublisherID),
					resource.TestCheckResourceAttr("powerplatform_unmanaged_solution.solution", "description", "Created by Terraform"),
				),
			},
			{
				ResourceName:      "powerplatform_unmanaged_solution.solution",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     unmanagedSolutionEnvironmentID + "_" + unmanagedSolutionID,
			},
		},
	})
}

func TestUnitUnmanagedSolutionResource_Validate_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	updated := false

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/"+unmanagedSolutionEnvironmentID+"?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Unmanaged_Update/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/solutions",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusCreated, "")
			resp.Header.Set("OData-EntityId", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/solutions("+unmanagedSolutionID+")")
			return resp, nil
		})

	httpmock.RegisterRegexpResponder("PATCH", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.0/solutions%28`+regexp.QuoteMeta(unmanagedSolutionID)+`%29$`),
		func(req *http.Request) (*http.Response, error) {
			updated = true
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=solutionid+eq+"+unmanagedSolutionID,
		func(req *http.Request) (*http.Response, error) {
			fileName := "tests/resource/Validate_Unmanaged_Update/get_solution_before.json"
			if updated {
				fileName = "tests/resource/Validate_Unmanaged_Update/get_solution_after.json"
			}
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fileName).String()), nil
		})

	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.2/solutions%28`+regexp.QuoteMeta(unmanagedSolutionID)+`%29$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `

				resource "powerplatform_unmanaged_solution" "solution" {
					environment_id = "` + unmanagedSolutionEnvironmentID + `"
					uniquename     = "TerraformTestSolution"
					display_name   = "Terraform Test Solution"
					publisher_id   = "` + unmanagedPublisherID + `"
					description    = "Created by Terraform"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_unmanaged_solution.solution", "display_name", "Terraform Test Solution"),
					resource.TestCheckResourceAttr("powerplatform_unmanaged_solution.solution", "description", "Created by Terraform"),
				),
			},
			{
				Config: `

				resource "powerplatform_unmanaged_solution" "solution" {
					environment_id = "` + unmanagedSolutionEnvironmentID + `"
					uniquename     = "TerraformTestSolution"
					display_name   = "Terraform Test Solution Updated"
					publisher_id   = "` + unmanagedPublisherID + `"
					description    = "Updated by Terraform"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_unmanaged_solution.solution", "display_name", "Terraform Test Solution Updated"),
					resource.TestCheckResourceAttr("powerplatform_unmanaged_solution.solution", "description", "Updated by Terraform"),
				),
			},
		},
	})
}

func TestUnitUnmanagedSolutionResource_Validate_Create_No_Dataverse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/"+unmanagedSolutionEnvironmentID+"?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Unmanaged_No_Dataverse/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `

				resource "powerplatform_unmanaged_solution" "solution" {
					environment_id = "` + unmanagedSolutionEnvironmentID + `"
					uniquename     = "TerraformTestSolution"
					display_name   = "Terraform Test Solution"
					publisher_id   = "` + unmanagedPublisherID + `"
				}`,
				ExpectError: regexp.MustCompile(fmt.Sprintf("No Dataverse exists in environment '%s'", unmanagedSolutionEnvironmentID)),
			},
		},
	})
}

func TestUnitUnmanagedSolutionResource_Validate_Create_Eventual_Consistency(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	readAttempts := 0

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/"+unmanagedSolutionEnvironmentID+"?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Unmanaged_Create/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/solutions",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=uniquename+eq+%27TerraformTestSolution%27",
		func(req *http.Request) (*http.Response, error) {
			readAttempts++
			if readAttempts == 1 {
				return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Unmanaged_Create/get_solution.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=solutionid+eq+"+unmanagedSolutionID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Unmanaged_Create/get_solution.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.2/solutions%28`+regexp.QuoteMeta(unmanagedSolutionID)+`%29$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `

				resource "powerplatform_unmanaged_solution" "solution" {
					environment_id = "` + unmanagedSolutionEnvironmentID + `"
					uniquename     = "TerraformTestSolution"
					display_name   = "Terraform Test Solution"
					publisher_id   = "` + unmanagedPublisherID + `"
					description    = "Created by Terraform"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_unmanaged_solution.solution", "id", unmanagedSolutionID),
					resource.TestCheckResourceAttr("powerplatform_unmanaged_solution.solution", "uniquename", "TerraformTestSolution"),
				),
			},
		},
	})

	if readAttempts < 2 {
		t.Fatalf("expected create flow to poll for visibility, got %d read attempt(s)", readAttempts)
	}
}

func TestUnitUnmanagedSolutionResource_Validate_Create_Uses_Response_ID(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/"+unmanagedSolutionEnvironmentID+"?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Unmanaged_Create/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/solutions",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, `{
				"solutionid":"`+unmanagedSolutionID+`",
				"uniquename":"TerraformTestSolution",
				"friendlyname":"Terraform Test Solution",
				"description":"Created by Terraform",
				"_publisherid_value":"`+unmanagedPublisherID+`"
			}`), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=solutionid+eq+"+unmanagedSolutionID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Unmanaged_Create/get_solution.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.2/solutions%28`+regexp.QuoteMeta(unmanagedSolutionID)+`%29$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `

				resource "powerplatform_unmanaged_solution" "solution" {
					environment_id = "` + unmanagedSolutionEnvironmentID + `"
					uniquename     = "TerraformTestSolution"
					display_name   = "Terraform Test Solution"
					publisher_id   = "` + unmanagedPublisherID + `"
					description    = "Created by Terraform"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_unmanaged_solution.solution", "id", unmanagedSolutionID),
					resource.TestCheckResourceAttr("powerplatform_unmanaged_solution.solution", "publisher_id", unmanagedPublisherID),
				),
			},
		},
	})
}

func TestUnitUnmanagedSolutionResource_Validate_Managed_Solution(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/"+unmanagedSolutionEnvironmentID+"?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Unmanaged_Create/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/solutions",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusNoContent, "")
			resp.Header.Set("OData-EntityId", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/solutions("+unmanagedSolutionID+")")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=solutionid+eq+"+unmanagedSolutionID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"value": [
					{
						"solutionid": "`+unmanagedSolutionID+`",
						"uniquename": "TerraformTestSolution",
						"friendlyname": "Terraform Test Solution",
						"_publisherid_value": "`+unmanagedPublisherID+`",
						"ismanaged": true
					}
				]
			}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `

				resource "powerplatform_unmanaged_solution" "solution" {
					environment_id = "` + unmanagedSolutionEnvironmentID + `"
					uniquename     = "TerraformTestSolution"
					display_name   = "Terraform Test Solution"
					publisher_id   = "` + unmanagedPublisherID + `"
					description    = "Created by Terraform"
				}`,
				ExpectError: regexp.MustCompile(`solution 'TerraformTestSolution' is managed and cannot be used with\s+powerplatform_unmanaged_solution`),
			},
		},
	})
}

func TestAccUnmanagedSolutionResource_Validate_Create_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source: "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccUnmanagedSolutionConfig(
					mocks.TestName(),
					"TerraformUnmanagedSolutionAcc",
					"Terraform Unmanaged Solution",
					"Created by Terraform acceptance test",
					true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_unmanaged_solution.solution", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_unmanaged_solution.solution", "environment_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_unmanaged_solution.solution", "uniquename", "TerraformUnmanagedSolutionAcc"),
					resource.TestCheckResourceAttr("powerplatform_unmanaged_solution.solution", "display_name", "Terraform Unmanaged Solution"),
					resource.TestCheckResourceAttr("powerplatform_unmanaged_solution.solution", "description", "Created by Terraform acceptance test"),
					resource.TestMatchResourceAttr("powerplatform_unmanaged_solution.solution", "publisher_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_unmanaged_solution.lookup", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("data.powerplatform_unmanaged_solution.lookup", "uniquename", "TerraformUnmanagedSolutionAcc"),
					resource.TestCheckResourceAttr("data.powerplatform_unmanaged_solution.lookup", "display_name", "Terraform Unmanaged Solution"),
					resource.TestCheckResourceAttr("data.powerplatform_unmanaged_solution.lookup", "description", "Created by Terraform acceptance test"),
					resource.TestMatchResourceAttr("data.powerplatform_unmanaged_solution.lookup", "publisher_id", regexp.MustCompile(helpers.GuidRegex)),
				),
			},
			{
				Config: testAccUnmanagedSolutionConfig(
					mocks.TestName(),
					"TerraformUnmanagedSolutionAcc",
					"Terraform Unmanaged Solution Updated",
					"Updated by Terraform acceptance test",
					true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_unmanaged_solution.solution", "display_name", "Terraform Unmanaged Solution Updated"),
					resource.TestCheckResourceAttr("powerplatform_unmanaged_solution.solution", "description", "Updated by Terraform acceptance test"),
					resource.TestCheckResourceAttr("data.powerplatform_unmanaged_solution.lookup", "display_name", "Terraform Unmanaged Solution Updated"),
					resource.TestCheckResourceAttr("data.powerplatform_unmanaged_solution.lookup", "description", "Updated by Terraform acceptance test"),
				),
			},
		},
	})
}

func testAccUnmanagedSolutionConfig(environmentDisplayName, uniqueName, displayName, description string, includeLookup bool) string {
	lookupBlock := ""
	if includeLookup {
		lookupBlock = `

				data "powerplatform_unmanaged_solution" "lookup" {
					depends_on     = [powerplatform_unmanaged_solution.solution]
					environment_id = powerplatform_environment.environment.id
					uniquename     = powerplatform_unmanaged_solution.solution.uniquename
				}`
	}

	return fmt.Sprintf(`

				resource "powerplatform_environment" "environment" {
					display_name     = "%s"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
						language_code     = "1033"
						currency_code     = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "time_sleep" "wait_120_seconds" {
					depends_on      = [powerplatform_environment.environment]
					create_duration = "120s"
				}

				data "powerplatform_rest_query" "default_publisher" {
					depends_on            = [time_sleep.wait_120_seconds]
					scope                 = "${powerplatform_environment.environment.dataverse.url}/.default"
					url                   = "${powerplatform_environment.environment.dataverse.url}api/data/v9.2/publishers?$select=publisherid,uniquename&$filter=startswith(uniquename,'DefaultPublisher')&$top=1"
					method                = "GET"
					expected_http_status  = [200]
				}

				locals {
					default_publisher = jsondecode(data.powerplatform_rest_query.default_publisher.output.body).value[0]
				}

				resource "powerplatform_unmanaged_solution" "solution" {
					environment_id = powerplatform_environment.environment.id
					uniquename     = "%s"
					display_name   = "%s"
					publisher_id   = local.default_publisher.publisherid
					description    = "%s"
				}%s
	`, environmentDisplayName, uniqueName, displayName, description, lookupBlock)
}
