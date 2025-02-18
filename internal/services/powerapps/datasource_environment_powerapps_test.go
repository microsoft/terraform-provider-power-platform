// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package powerapps_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccEnvironmentPowerAppsDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "env" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
				}

				data "powerplatform_environment_powerapps" "all" {
					depends_on = [powerplatform_environment.env]
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the first power app to ensure all attributes are set.
					resource.TestMatchResourceAttr("data.powerplatform_environment_powerapps.all", "powerapps.0.name", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_powerapps.all", "powerapps.0.id", regexp.MustCompile(helpers.UrlValidStringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_powerapps.all", "powerapps.0.display_name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_powerapps.all", "powerapps.0.created_time", regexp.MustCompile(`^\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+([+-][0-2]\d:[0-5]\d|Z)$`)),
				),
			},
		},
	})
}

func TestUnitEnvironmentPowerAppsDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24expand=properties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Read/get_environments.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.powerapps\.com/providers/Microsoft\.PowerApps/scopes/admin/environments/([\d-]+)/apps`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Read/get_apps_"+id+".json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_environment_powerapps" "all" {}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_environment_powerapps.all", "powerapps.#", "4"),

					resource.TestCheckResourceAttr("data.powerplatform_environment_powerapps.all", "powerapps.0.name", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_powerapps.all", "powerapps.0.id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_powerapps.all", "powerapps.0.display_name", "Overview"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_powerapps.all", "powerapps.0.created_time", "2023-09-27T07:08:47.1964785Z"),

					resource.TestCheckResourceAttr("data.powerplatform_environment_powerapps.all", "powerapps.2.name", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_powerapps.all", "powerapps.2.id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_powerapps.all", "powerapps.2.display_name", "Overview"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_powerapps.all", "powerapps.2.created_time", "2023-09-27T07:08:47.1964785Z"),
				),
			},
		},
	})
}
