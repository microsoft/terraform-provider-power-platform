// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package copilot_studio_application_insights_test

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccCopilotStudioApplicationInsights_Validate_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azapi": {
				VersionConstraint: constants.AZAPI_PROVIDER_VERSION_CONSTRAINT,
				Source:            "azure/azapi",
			},
			"time": {
				Source: "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
					resource "azapi_resource" "rg_example" {
						type     = "Microsoft.Resources/resourceGroups@2021-04-01"
						location = "East US"
						name     = "power-platform-app-insights-rg-` + mocks.TestName() + strconv.Itoa(rand.IntN(9999)) + `"
					}

					resource "azapi_resource" "app_insights" {
						schema_validation_enabled = false

						type = "Microsoft.Insights/components@2020-02-02"
						location = azapi_resource.rg_example.location
						name = "power-platform-app-insights-` + mocks.TestName() + strconv.Itoa(rand.IntN(9999)) + `"
						parent_id = azapi_resource.rg_example.id

						body = {
							properties = {
								Application_Type = "web"
								Flow_Type = "Bluefield"
							}
						}
					}

					resource "powerplatform_environment" "environment" {
						display_name     = "` + mocks.TestName() + `"
						location         = "unitedstates"
						environment_type = "Sandbox"
						dataverse = {
							language_code     = "1033"
							currency_code     = "USD"
							security_group_id = "00000000-0000-0000-0000-000000000000"
						}
					}

					resource "powerplatform_solution" "solution" {
						environment_id = powerplatform_environment.environment.id
						solution_file  = "tests/Test_Files/testagent_1_0_0_1_managed.zip"
					}

					resource "time_sleep" "wait_60_seconds" {
						depends_on = [powerplatform_solution.solution]
						create_duration = "60s"
					}

					data "powerplatform_data_records" "bot_data_query" {
						environment_id    = powerplatform_environment.environment.id
						entity_collection = "bots"
						filter            = "name eq 'Test Agent'"
						select            = ["botid", "name"]
						top               = 1

						depends_on = [powerplatform_solution.solution, time_sleep.wait_60_seconds]
					}

					resource "powerplatform_copilot_studio_application_insights" "cps_app_insights_config" {
						environment_id                         = powerplatform_environment.environment.id
						bot_id                                 = data.powerplatform_data_records.bot_data_query.rows[0].botid
						application_insights_connection_string = azapi_resource.app_insights.output.properties.ConnectionString
						include_activities                     = true
						include_sensitive_information          = true
						include_actions                        = true

						depends_on = [azapi_resource.app_insights, time_sleep.wait_60_seconds]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_copilot_studio_application_insights.cps_app_insights_config", "application_insights_connection_string", regexp.MustCompile(`^InstrumentationKey=.*$`)),
					resource.TestCheckResourceAttr("powerplatform_copilot_studio_application_insights.cps_app_insights_config", "include_sensitive_information", "true"),
					resource.TestCheckResourceAttr("powerplatform_copilot_studio_application_insights.cps_app_insights_config", "include_activities", "true"),
					resource.TestCheckResourceAttr("powerplatform_copilot_studio_application_insights.cps_app_insights_config", "include_actions", "true"),
				),
			},
		},
	})
}

func TestUnitCopilotStudioApplicationInsights_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Create/get_environments.json").String()), nil
		})

	httpmock.RegisterResponder("PUT", "https://powervamg.eu-il107.gateway.prod.island.powerapps.com/api/botmanagement/2022-01-15/environments/00000000-0000-0000-0000-000000000001/bots/00000000-0000-0000-0000-000000000002/applicationinsightsconfiguration",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Create/app_insights.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://powervamg.eu-il107.gateway.prod.island.powerapps.com/api/botmanagement/2022-01-15/environments/00000000-0000-0000-0000-000000000001/bots/00000000-0000-0000-0000-000000000002/applicationinsightsconfiguration",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Create/app_insights.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "powerplatform_copilot_studio_application_insights" "cps_app_insights_config" {
						environment_id                         = "00000000-0000-0000-0000-000000000001"
						bot_id                                 = "00000000-0000-0000-0000-000000000002"
						application_insights_connection_string = "<<connection_string>>"
						include_sensitive_information          = false
						include_activities                     = false
						include_actions                        = false
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_copilot_studio_application_insights.cps_app_insights_config", "environment_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_copilot_studio_application_insights.cps_app_insights_config", "bot_id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("powerplatform_copilot_studio_application_insights.cps_app_insights_config", "application_insights_connection_string", "<<connection_string>>"),
					resource.TestCheckResourceAttr("powerplatform_copilot_studio_application_insights.cps_app_insights_config", "include_sensitive_information", "false"),
					resource.TestCheckResourceAttr("powerplatform_copilot_studio_application_insights.cps_app_insights_config", "include_activities", "false"),
					resource.TestCheckResourceAttr("powerplatform_copilot_studio_application_insights.cps_app_insights_config", "include_actions", "false"),
				),
			},
		},
	})
}

func TestUnitCopilotStudioApplicationInsights_Validate_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	put_app_insights_inx := 0
	get_app_insights_inx := 0

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Update/get_environments.json").String()), nil
		})

	httpmock.RegisterResponder("PUT", "https://powervamg.eu-il107.gateway.prod.island.powerapps.com/api/botmanagement/2022-01-15/environments/00000000-0000-0000-0000-000000000001/bots/00000000-0000-0000-0000-000000000002/applicationinsightsconfiguration",
		func(req *http.Request) (*http.Response, error) {
			put_app_insights_inx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/Validate_Update/put_app_insights_%d.json", put_app_insights_inx)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://powervamg.eu-il107.gateway.prod.island.powerapps.com/api/botmanagement/2022-01-15/environments/00000000-0000-0000-0000-000000000001/bots/00000000-0000-0000-0000-000000000002/applicationinsightsconfiguration",
		func(req *http.Request) (*http.Response, error) {
			get_app_insights_inx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/Validate_Update/get_app_insights_%d.json", get_app_insights_inx)).String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "powerplatform_copilot_studio_application_insights" "cps_app_insights_config" {
						environment_id                         = "00000000-0000-0000-0000-000000000001"
						bot_id                                 = "00000000-0000-0000-0000-000000000002"
						application_insights_connection_string = "<<connection_string>>"
						include_sensitive_information          = false
						include_activities                     = false
						include_actions                        = false
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_copilot_studio_application_insights.cps_app_insights_config", "environment_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_copilot_studio_application_insights.cps_app_insights_config", "bot_id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("powerplatform_copilot_studio_application_insights.cps_app_insights_config", "application_insights_connection_string", "<<connection_string>>"),
					resource.TestCheckResourceAttr("powerplatform_copilot_studio_application_insights.cps_app_insights_config", "include_sensitive_information", "false"),
					resource.TestCheckResourceAttr("powerplatform_copilot_studio_application_insights.cps_app_insights_config", "include_activities", "false"),
					resource.TestCheckResourceAttr("powerplatform_copilot_studio_application_insights.cps_app_insights_config", "include_actions", "false"),
				),
			},
			{
				Config: `
					resource "powerplatform_copilot_studio_application_insights" "cps_app_insights_config" {
						environment_id                         = "00000000-0000-0000-0000-000000000001"
						bot_id                                 = "00000000-0000-0000-0000-000000000002"
						application_insights_connection_string = "<<connection_string_new>>"
						include_sensitive_information          = true
						include_activities                     = true
						include_actions                        = true
					}
				`,
			},
		},
	})
}
