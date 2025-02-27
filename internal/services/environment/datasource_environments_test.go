// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccEnvironmentsDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "env" {
					display_name     = "` + mocks.TestName() + `"
					description      = "description"
					location         = "europe"
					azure_region     = "northeurope"
					environment_type = "Sandbox"
					cadence = "Moderate"
					dataverse = {
						language_code     = "1033"
						currency_code     = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				data "powerplatform_environments" "all" {
					depends_on = [powerplatform_environment.env]
				}
				
				output "test_environment"{
					value = one([for env in data.powerplatform_environments.all.environments : env if env.id == powerplatform_environment.env.id])
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test_environment", knownvalue.NotNull()),
					statecheck.ExpectKnownOutputValueAtPath("test_environment", tfjsonpath.New("id"), knownvalue.StringRegexp(regexp.MustCompile(helpers.GuidRegex))),
					statecheck.ExpectKnownOutputValueAtPath("test_environment", tfjsonpath.New("display_name"), knownvalue.StringExact(mocks.TestName())),
					statecheck.ExpectKnownOutputValueAtPath("test_environment", tfjsonpath.New("description"), knownvalue.StringExact("description")),
					statecheck.ExpectKnownOutputValueAtPath("test_environment", tfjsonpath.New("location"), knownvalue.StringExact("europe")),
					statecheck.ExpectKnownOutputValueAtPath("test_environment", tfjsonpath.New("azure_region"), knownvalue.StringExact("northeurope")),
					statecheck.ExpectKnownOutputValueAtPath("test_environment", tfjsonpath.New("environment_type"), knownvalue.StringExact("Sandbox")),
					statecheck.ExpectKnownOutputValueAtPath("test_environment", tfjsonpath.New("cadence"), knownvalue.StringExact("Moderate")),
					statecheck.ExpectKnownOutputValueAtPath("test_environment", tfjsonpath.New("dataverse").AtMapKey("language_code"), knownvalue.Int32Exact(1033)),
					statecheck.ExpectKnownOutputValueAtPath("test_environment", tfjsonpath.New("dataverse").AtMapKey("currency_code"), knownvalue.StringExact("USD")),
					statecheck.ExpectKnownOutputValueAtPath("test_environment", tfjsonpath.New("dataverse").AtMapKey("security_group_id"), knownvalue.StringRegexp(regexp.MustCompile(helpers.GuidOrEmptyValueRegex))),
					statecheck.ExpectKnownOutputValueAtPath("test_environment", tfjsonpath.New("dataverse").AtMapKey("organization_id"), knownvalue.StringRegexp(regexp.MustCompile(helpers.GuidRegex))),
					statecheck.ExpectKnownOutputValueAtPath("test_environment", tfjsonpath.New("dataverse").AtMapKey("url"), knownvalue.StringRegexp(regexp.MustCompile(helpers.UrlValidStringRegex))),
					statecheck.ExpectKnownOutputValueAtPath("test_environment", tfjsonpath.New("dataverse").AtMapKey("version"), knownvalue.StringRegexp(regexp.MustCompile(helpers.VersionRegex))),
					statecheck.ExpectKnownOutputValueAtPath("test_environment", tfjsonpath.New("dataverse").AtMapKey("unique_name"), knownvalue.StringRegexp(regexp.MustCompile(helpers.StringRegex))),
				},
			},
		},
	})
}

func TestUnitEnvironmentsDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24expand=properties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_environments.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000002?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_environment_00000000-0000-0000-0000-000000000002.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_environments" "all" {}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.#", "2"),

					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.cadence", "Moderate"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.description", "aaa"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.display_name", "Admin AdminOnMicrosoft's Environment"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.dataverse.domain", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.environment_type", "Developer"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.dataverse.language_code", "1033"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.dataverse.organization_id", "6450637c-f9a8-4988-8cf7-b03723d51ab7"),
					resource.TestCheckNoResourceAttr("data.powerplatform_environments.all", "environments.0.dataverse.security_group_id"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.dataverse.url", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.location", "europe"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.azure_region", "northeurope"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.dataverse.version", "9.2.23092.00206"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.dataverse.unique_name", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.dataverse.currency_code", "PLN"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.billing_policy_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.dataverse.linked_app_type", ""),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.dataverse.linked_app_id", ""),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.dataverse.linked_app_url", ""),
					resource.TestCheckNoResourceAttr("data.powerplatform_environments.all", "environments.0.dataverse.templates"),
					resource.TestCheckNoResourceAttr("data.powerplatform_environments.all", "environments.0.dataverse.template_metadata"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.environment_group_id", "00000000-0000-0000-0000-000000000001"),

					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.cadence", "Frequent"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.description", "bbb"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.display_name", "displayname"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.environment_type", "Sandbox"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.location", "europe"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.azure_region", "westeurope"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.billing_policy_id", ""),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.environment_group_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckNoResourceAttr("data.powerplatform_environments.all", "environments.1.dataverse.domain"),
					resource.TestCheckNoResourceAttr("data.powerplatform_environments.all", "environments.1.dataverse.language_code"),
					resource.TestCheckNoResourceAttr("data.powerplatform_environments.all", "environments.1.dataverse.organization_id"),
					resource.TestCheckNoResourceAttr("data.powerplatform_environments.all", "environments.1.dataverse.security_group_id"),
					resource.TestCheckNoResourceAttr("data.powerplatform_environments.all", "environments.1.dataverse.url"),
					resource.TestCheckNoResourceAttr("data.powerplatform_environments.all", "environments.1.dataverse.version"),
					resource.TestCheckNoResourceAttr("data.powerplatform_environments.all", "environments.1.dataverse.unique_name"),
					resource.TestCheckNoResourceAttr("data.powerplatform_environments.all", "environments.1.dataverse.currency_code"),
				),
			},
		},
	})
}
