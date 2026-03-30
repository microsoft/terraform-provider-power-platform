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
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitUnmanagedSolutionDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/"+unmanagedSolutionEnvironmentID+"?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Unmanaged_Create/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=uniquename+eq+%27TerraformTestSolution%27",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Unmanaged_Create/get_solution.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_unmanaged_solution" "solution" {
					environment_id = "` + unmanagedSolutionEnvironmentID + `"
					uniquename     = "TerraformTestSolution"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_unmanaged_solution.solution", "id", unmanagedSolutionEnvironmentID+"_"+unmanagedSolutionID),
					resource.TestCheckResourceAttr("data.powerplatform_unmanaged_solution.solution", "environment_id", unmanagedSolutionEnvironmentID),
					resource.TestCheckResourceAttr("data.powerplatform_unmanaged_solution.solution", "uniquename", "TerraformTestSolution"),
					resource.TestCheckResourceAttr("data.powerplatform_unmanaged_solution.solution", "display_name", "Terraform Test Solution"),
					resource.TestCheckResourceAttr("data.powerplatform_unmanaged_solution.solution", "publisher_id", unmanagedPublisherID),
					resource.TestCheckResourceAttr("data.powerplatform_unmanaged_solution.solution", "description", "Created by Terraform"),
				),
			},
		},
	})
}

func TestUnitUnmanagedSolutionDataSource_Validate_No_Dataverse(t *testing.T) {
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
				data "powerplatform_unmanaged_solution" "solution" {
					environment_id = "` + unmanagedSolutionEnvironmentID + `"
					uniquename     = "TerraformTestSolution"
				}`,
				ExpectError: regexp.MustCompile(fmt.Sprintf("No Dataverse exists in environment '%s'", unmanagedSolutionEnvironmentID)),
			},
		},
	})
}

func TestUnitUnmanagedSolutionDataSource_Validate_Not_Found(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/"+unmanagedSolutionEnvironmentID+"?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Unmanaged_Create/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=uniquename+eq+%27MissingSolution%27",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_unmanaged_solution" "solution" {
					environment_id = "` + unmanagedSolutionEnvironmentID + `"
					uniquename     = "MissingSolution"
				}`,
				ExpectError: regexp.MustCompile("Unmanaged solution 'MissingSolution' not found"),
			},
		},
	})
}
