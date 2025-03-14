// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package solution_checker_rules_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccSolutionCheckerRulesDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "test" {
					display_name     = "` + mocks.TestName() + `"
					location        = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
						language_code     = "1033"
						currency_code     = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				data "powerplatform_solution_checker_rules" "test" {
					environment_id = powerplatform_environment.test.id
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_solution_checker_rules.test", "rules.0.code", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solution_checker_rules.test", "rules.0.description", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solution_checker_rules.test", "rules.0.summary", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solution_checker_rules.test", "rules.0.guidance_url", regexp.MustCompile(helpers.UrlValidStringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solution_checker_rules.test", "rules.0.primary_category_description", regexp.MustCompile(helpers.StringRegex)),
				),
			},
		},
	})
}

func TestUnitSolutionCheckerRulesDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register mock responses
	httpmock.RegisterResponder(
		"GET",
		"https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001",
		httpmock.NewStringResponder(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_environment.json").String()))

	// Register responder for the rules API with exact URL match
	httpmock.RegisterResponder(
		"GET",
		"https://unitedstates.api.advisor.powerapps.com/api/rule?api-version=2.0&ruleset=0ad12346-e108-40b8-a956-9a8f95ea18c9",
		httpmock.NewStringResponder(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_rules.json").String()))

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_solution_checker_rules" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_solution_checker_rules.test", "rules.0.code", "meta-remove-dup-reg"),
					resource.TestCheckResourceAttr("data.powerplatform_solution_checker_rules.test", "rules.0.description", "Checks for duplicate Dataverse plug-in registrations"),
					resource.TestCheckResourceAttr("data.powerplatform_solution_checker_rules.test", "rules.0.summary", "Duplicate plug-in registration"),
					resource.TestCheckResourceAttr("data.powerplatform_solution_checker_rules.test", "rules.0.how_to_fix", ""),
					resource.TestCheckResourceAttr("data.powerplatform_solution_checker_rules.test", "rules.0.guidance_url", "https://learn.microsoft.com/powerapps/developer/data-platform/best-practices/business-logic/do-not-duplicate-plugin-step-registration"),
					resource.TestCheckResourceAttr("data.powerplatform_solution_checker_rules.test", "rules.0.component_type", "0"),
					resource.TestCheckResourceAttr("data.powerplatform_solution_checker_rules.test", "rules.0.primary_category", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_solution_checker_rules.test", "rules.0.primary_category_description", "Performance"),
					resource.TestCheckResourceAttr("data.powerplatform_solution_checker_rules.test", "rules.0.include", "true"),
					resource.TestCheckResourceAttr("data.powerplatform_solution_checker_rules.test", "rules.0.severity", "5"),
				),
			},
		},
	})
}
