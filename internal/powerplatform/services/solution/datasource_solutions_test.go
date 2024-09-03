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
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/provider"
)

func TestUnitSolutionsDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_environments.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=%28isvisible+eq+true%29&%24orderby=createdon+desc`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_solution.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsUnitProviderConfig + `
				data "powerplatform_solutions" "all" {
					environment_id = "00000000-0000-0000-0000-000000000001"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.name", "ProductivityToolsAnchor"),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.environment_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.display_name", "ProductivityTools"),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.created_time", "2023-10-10T08:09:56Z"),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.modified_time", "2023-10-10T08:09:58Z"),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.install_time", "2023-10-10T08:09:56Z"),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.is_managed", "true"),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.version", "9.2.1.1020"),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.id", "70edca66-e4c2-4384-92e0-4300465c1894"),
				),
			},
		},
	})
}

func TestUnitSolutionsDataSource_Validate_No_Dataverse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_No_Dataverse(/get_environments.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_No_Dataverse/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=%28isvisible+eq+true%29&%24orderby=createdon+desc`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_No_Dataverse/get_solution.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../../services/application/tests/resource/Validate_No_Dataverse/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("../../services/application/tests/resource/Validate_No_Dataverse/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../../services/application/tests/resource/Validate_No_Dataverse/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsUnitProviderConfig + `
				resource "powerplatform_environment" "env" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
				}

				data "powerplatform_solutions" "all" {
					environment_id = powerplatform_environment.env.id
				}`,
				ExpectError: regexp.MustCompile(`No Dataverse exists in environment`),
				Check:       resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestAccSolutionsDataSource_Validate_Read(t *testing.T) {

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsAcceptanceProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
						language_code    = "1033"
						currency_code    = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				data "powerplatform_solutions" "all" {
					environment_id = powerplatform_environment.development.id
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					// Verify the first power app to ensure all attributes are set
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.environment_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.display_name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.created_time", regexp.MustCompile(helpers.TimeRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.modified_time", regexp.MustCompile(helpers.TimeRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.install_time", regexp.MustCompile(helpers.TimeRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.is_managed", regexp.MustCompile(`^(true|false)$`)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.version", regexp.MustCompile(helpers.VersionRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.id", regexp.MustCompile(helpers.GuidRegex)),
				),
			},
		},
	})
}
