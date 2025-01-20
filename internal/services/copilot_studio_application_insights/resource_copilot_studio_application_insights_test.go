// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package copilot_studio_application_insights_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitCopilotStudioApplicationInsights_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
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

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
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
