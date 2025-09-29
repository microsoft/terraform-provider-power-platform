// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package connectors_test

import (
	"net/http"
	"os"
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

const (
	SOLUTION_NAME          = "TerraformTestCustomConnector_1_0_0_1.zip"
	SOLUTION_RELATIVE_PATH = "tests/Test_Files/" + SOLUTION_NAME
)

func TestAccConnectorsDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_connectors" "all" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(

					// Verify the first power app to ensure all attributes are set.
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.description", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.display_name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.id", regexp.MustCompile(helpers.ApiIdRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.publisher", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.tier", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.type", regexp.MustCompile(helpers.ApiIdRegex)),
				),
			},
		},
	})
}

func TestUnitConnectorsDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/connectors/metadata/virtual`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Read/get_virtual.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/connectors/metadata/unblockable`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Read/get_unblockable.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.powerapps.com/providers/Microsoft.PowerApps/apis?%24filter=environment+eq+%27~Default%27&api-version=2019-05-01&hideDlpExemptApis=true&showAllDlpEnforceableApis=true&showApisWithToS=true`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Read/get_apis.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_connectors" "all" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify returned count.
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.#", "4"),

					// Verify the first power app to ensure all attributes are set
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.description", "SharePoint helps organizations share and collaborate with colleagues, partners, and customers. You can connect to SharePoint Online or to an on-premises SharePoint 2013 or 2016 farm using the On-Premises Data Gateway to manage documents and list items."),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.display_name", "SharePoint"),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.id", "/providers/Microsoft.PowerApps/apis/shared_sharepointonline"),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.name", "shared_sharepointonline"),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.publisher", "Microsoft"),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.tier", "Standard"),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.type", "Microsoft.PowerApps/apis"),
				),
			},
		},
	})
}

func TestAccTestUnitConnectorsDataSource_Validate_Read_With_Given_Environment(t *testing.T) {
	solutionFileBytes, err := os.ReadFile(SOLUTION_RELATIVE_PATH)
	if err != nil {
		t.Fatalf("Failed to read solution file: %s", err.Error())
	}

	err = os.WriteFile(SOLUTION_NAME, solutionFileBytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write solution file: %s", err.Error())
	}

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source: "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "environment" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                           = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}


				resource "time_sleep" "wait_90_seconds" {
					depends_on = [powerplatform_environment.environment]
					# Wait 90 seconds to allow the environment to be fully provisioned before importing the solution.
					# This is required as the environment may not be fully ready to accept solutions immediately after creation due to 
					# "Async operations are currently disabled for this organization." message.

					create_duration = "90s"
				}

				resource "powerplatform_solution" "solution" {
					depends_on      = [time_sleep.wait_90_seconds]

					environment_id = powerplatform_environment.environment.id
					solution_file    = "` + SOLUTION_NAME + `"
				}

				data "powerplatform_connectors" "all_connectors" {
					depends_on      = [powerplatform_solution.solution]

					environment_id = powerplatform_environment.environment.id
				}

				output "custom_connector" {
					value = {
						connectors = [for connector in data.powerplatform_connectors.all_connectors.connectors : connector if connector.display_name == "Custom Connector"]
					}
				}`,

				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValueAtPath(
						"custom_connector",
						// path -> connectors[0].display_name
						tfjsonpath.New("connectors").AtSliceIndex(0).AtMapKey("display_name"),
						knownvalue.StringExact("Custom Connector"),
					),
				},
			},
		},
	})
}

func TestUnitConnectorsDataSource_Validate_Read_With_Given_Environment(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/connectors/metadata/virtual`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Read/get_virtual.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/connectors/metadata/unblockable`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Read/get_unblockable.json").String()), nil
		})

	// Mock the API call with a specific environment ID
	httpmock.RegisterResponder("GET", `https://api.powerapps.com/providers/Microsoft.PowerApps/apis?%24filter=environment+eq+%2700000000-0000-0000-0000-000000000001%27&api-version=2019-05-01&hideDlpExemptApis=true&showAllDlpEnforceableApis=true&showApisWithToS=true`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Read/get_apis.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_connectors" "all_connectors" {
					environment_id = "00000000-0000-0000-0000-000000000001"
				}
					
				output "custom_connector" {
					value = {
						connectors = [for connector in data.powerplatform_connectors.all_connectors.connectors : connector if connector.display_name == "SharePoint"]
					}
				}
				`,

				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValueAtPath(
						"custom_connector",
						// path -> connectors[0].display_name
						tfjsonpath.New("connectors").AtSliceIndex(0).AtMapKey("display_name"),
						knownvalue.StringExact("SharePoint"),
					),
					statecheck.ExpectKnownOutputValueAtPath(
						"custom_connector",
						// path -> connectors[0].id
						tfjsonpath.New("connectors").AtSliceIndex(0).AtMapKey("id"),
						knownvalue.StringExact("/providers/Microsoft.PowerApps/apis/shared_sharepointonline"),
					),
				},
			},
		},
	})
}
