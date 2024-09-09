// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_templates_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccEnvironmentTemplatesDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: constants.TestsAcceptanceProviderConfig + `
				data "powerplatform_environment_templates" "all_environment_templates_for_unitedstates" {
					location = "unitedstates"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "id", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.0.id", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.0.name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.0.display_name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.0.category", regexp.MustCompile(helpers.StringRegex)),
				),
			},
		},
	})
}

func TestUnitEnvironmentTemplatesDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/templates?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_environment_templates.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: constants.TestsUnitProviderConfig + `
				data "powerplatform_environment_templates" "all_environment_templates_for_unitedstates" {
					location = "unitedstates"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "id", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.#", "53"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.0.id", "/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentTemplates/D365_CDSSampleApp"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.0.name", "D365_CDSSampleApp"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.0.display_name", "Sample App"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.0.category", "developer"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.1.id", "/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentTemplates/D365_Sales"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.1.name", "D365_Sales"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.1.display_name", "Sales Enterprise"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.1.category", "production"),
				),
			},
		},
	})
}
