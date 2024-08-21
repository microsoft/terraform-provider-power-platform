// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package powerapps_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	mocks "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/provider"
)

func TestAccEnvironmentPowerAppsDataSource_Basic(t *testing.T) {

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsAcceptanceProviderConfig + `
				data "powerplatform_environment_powerapps" "all" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("data.powerplatform_environment_powerapps.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					// Verify the first power app to ensure all attributes are set
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

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24expand=properties%2FbillingPolicy&api-version=2023-06-01`,
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
		ProtoV6ProviderFactories: provider.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsUnitProviderConfig + `
				data "powerplatform_environment_powerapps" "all" {}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_environment_powerapps.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

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
