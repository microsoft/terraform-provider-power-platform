// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package solution_test

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

const (
	SOLUTION_1_NAME          = "TerraformTestSolution_Complex_1_1_0_0.zip"
	SOLUTION_1_RELATIVE_PATH = "tests/resource/Test_Files/" + SOLUTION_1_NAME

	SOLUTION_2_NAME          = "TerraformSimpleTestSolution_1_0_0_1_managed.zip"
	SOLUTION_2_RELATIVE_PATH = "tests/resource/Test_Files/" + SOLUTION_2_NAME
)

func TestAccSolutionResource_Uninstall_Multiple_Solutions(t *testing.T) {
	solutionFileBytes1, err := os.ReadFile(SOLUTION_1_RELATIVE_PATH)
	if err != nil {
		t.Fatalf("Failed to read solution file: %v", err)
	}

	err = os.WriteFile(SOLUTION_1_NAME, solutionFileBytes1, 0644)
	if err != nil {
		t.Fatalf("Failed to write solution file: %v", err)
	}

	solutionFileBytes2, err := os.ReadFile(filepath.Join(SOLUTION_2_RELATIVE_PATH))
	if err != nil {
		t.Fatalf("Failed to read solution file: %v", err)
	}

	err = os.WriteFile(SOLUTION_2_NAME, solutionFileBytes2, 0644)
	if err != nil {
		t.Fatalf("Failed to write solution file: %v", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
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

				resource "powerplatform_solution" "solution1" {
					environment_id = powerplatform_environment.environment.id
					solution_file    = "` + SOLUTION_1_NAME + `"
				}
					
				resource "powerplatform_solution" "solution2" {
					environment_id = powerplatform_environment.environment.id
					solution_file    = "` + SOLUTION_2_NAME + `"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestAccSolutionResource_Validate_Create_No_Settings_File(t *testing.T) {
	solutionFileBytes, err := os.ReadFile(SOLUTION_1_RELATIVE_PATH)
	if err != nil {
		t.Fatalf("Failed to read solution file: %v", err)
	}

	err = os.WriteFile(SOLUTION_1_NAME, solutionFileBytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write solution file: %v", err)
	}

	solutionFileChecksum, _ := helpers.CalculateSHA256(SOLUTION_1_NAME)
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
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

				resource "powerplatform_solution" "solution" {
					environment_id = powerplatform_environment.environment.id
					solution_file    = "` + SOLUTION_1_NAME + `"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("powerplatform_solution.solution", "settings_file_checksum"),
					resource.TestCheckNoResourceAttr("powerplatform_solution.solution", "settings_file"),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_file_checksum", solutionFileChecksum),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_file", SOLUTION_1_NAME),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "display_name", "Terraform Test Solution"),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "is_managed", "false"),
					resource.TestMatchResourceAttr("powerplatform_solution.solution", "environment_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_version", "1.1.0.0"),
				),
			},
		},
	})
}

func TestUnitSolutionResource_Validate_Create_With_Settings_File(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	solution_checksum := createFile("test_solution.zip", "test_solution")
	settings_checksum := createFile("test_solution_settings.json", "")

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_With_Settings_File/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/StageSolution",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_With_Settings_File/post_stage_solution.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/ImportSolutionAsync",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_With_Settings_File/post_import_solution_async.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/asyncoperations%28310799b8-dc6c-ee11-9ae7-000d3aaae21d%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_With_Settings_File/get_async_operations.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/RetrieveSolutionImportResult%28ImportJobId=1b1fa80d-aa0f-4291-b60c-b0745304ce24%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_With_Settings_File/get_solution_import_result.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=%28isvisible+eq+true%29&%24orderby=createdon+desc",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_With_Settings_File/get_solution.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=uniquename+eq+%27TerraformTestSolution%27",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_With_Settings_File/get_solution.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=solutionid+eq+86928ed8-df37-4ce2-add5-47030a833bff",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_With_Settings_File/get_solution.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions%2886928ed8-df37-4ce2-add5-47030a833bff%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, httpmock.File("tests/resource/Validate_Create_With_Settings_File/get_solution.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `

				resource "powerplatform_solution" "solution" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					solution_file    = "test_solution.zip"
					settings_file 	 = "test_solution_settings.json"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_version", "1.1.0.0"),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_file_checksum", solution_checksum),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "settings_file_checksum", settings_checksum),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "display_name", "Terraform Test Solution"),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "is_managed", strconv.FormatBool(false)),
				),
			},
		},
	})
}

func TestAccSolutionResource_Validate_Create_With_Settings_File(t *testing.T) {
	solutionSettingsFileName := "test_solution_settings.json"

	solutionFileBytes, err := os.ReadFile(SOLUTION_1_RELATIVE_PATH)
	if err != nil {
		t.Fatalf("Failed to read solution file: %v", err)
	}

	err = os.WriteFile(SOLUTION_1_NAME, solutionFileBytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write solution file: %v", err)
	}

	solutionSettingsContent := []byte(`{
		"EnvironmentVariables": [
		  {
			"SchemaName": "cra6e_SolutionVariableDataSource",
			"Value": "/sites/Shared%20Documents"
		  },
		  {
			"SchemaName": "cra6e_SolutionVariableJson",
			"Value": "{ \"value\": 1234, \"text\": \"abc\" }"
		  },
		  {
			"SchemaName": "cra6e_SolutionVariableText",
			"Value": "cd930b48-4bcc-e444-92e9-547b85c2fd4"
		  }
		],
		"ConnectionReferences": [
		  {
			"LogicalName": "cra6e_ConnectionReferenceSharePoint",
			"ConnectionId": "",
			"ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_sharepointonline"
		  }
		]
	  }`)

	err = os.WriteFile(solutionSettingsFileName, solutionSettingsContent, 0644)
	if err != nil {
		t.Fatalf("Failed to write settings file: %v", err)
	}

	solutionFileChecksum, _ := helpers.CalculateSHA256(SOLUTION_1_NAME)
	settingsFileChecksum, _ := helpers.CalculateSHA256(solutionSettingsFileName)
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `

				resource "powerplatform_environment" "environment" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                         = "1033"
						currency_code                         = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}
				
				resource "powerplatform_solution" "solution" {
					environment_id = powerplatform_environment.environment.id
					solution_file    = "` + SOLUTION_1_NAME + `"
					settings_file 	 = "` + solutionSettingsFileName + `"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_file_checksum", solutionFileChecksum),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "settings_file_checksum", settingsFileChecksum),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "settings_file", solutionSettingsFileName),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_file", SOLUTION_1_NAME),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "display_name", "Terraform Test Solution"),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "is_managed", "false"),
					resource.TestMatchResourceAttr("powerplatform_solution.solution", "environment_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_version", "1.1.0.0"),
				),
			},
		},
	})
}

func TestUnitSolutionResource_Validate_Create_No_Settings_File(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	solution_checksum := createFile("test_solution.zip", "test_solution")

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Settings_File/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/StageSolution",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Settings_File/post_stage_solution.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/ImportSolutionAsync",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Settings_File/post_import_solution_async.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/asyncoperations%28310799b8-dc6c-ee11-9ae7-000d3aaae21d%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Settings_File/get_async_operations.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/RetrieveSolutionImportResult%28ImportJobId=1b1fa80d-aa0f-4291-b60c-b0745304ce24%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Settings_File/get_solution_import_result.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=%28isvisible+eq+true%29&%24orderby=createdon+desc",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Settings_File/get_solution.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=uniquename+eq+%27TerraformTestSolution%27",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Settings_File/get_solution.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=solutionid+eq+86928ed8-df37-4ce2-add5-47030a833bff",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Settings_File/get_solution.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions%2886928ed8-df37-4ce2-add5-47030a833bff%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, httpmock.File("tests/resource/Validate_Create_No_Settings_File/get_solution.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `

				resource "powerplatform_solution" "solution" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					solution_file    = "test_solution.zip"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("powerplatform_solution.solution", "settings_file_checksum"),
					resource.TestCheckNoResourceAttr("powerplatform_solution.solution", "settings_file"),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_file_checksum", solution_checksum),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_version", "1.1.0.0"),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "display_name", "Terraform Test Solution"),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "is_managed", strconv.FormatBool(false)),
				),
			},
		},
	})
}
func TestUnitSolutionResource_Validate_Create_And_Force_Recreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	solution_before_checksum := createFile("test_solution_before.zip", "test_solution_before")
	settings_before_checksum := createFile("test_settings_before.json", "")
	solution_after_checksum := createFile("test_solution_after.zip", "test_solution_after")
	settings_after_checksum := createFile("test_settings_after.json", "")

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Force_Recreate/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/StageSolution",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Force_Recreate/post_stage_solution.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/ImportSolutionAsync",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Force_Recreate/post_import_solution_async.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/asyncoperations%28310799b8-dc6c-ee11-9ae7-000d3aaae21d%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Force_Recreate/get_async_operations.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/RetrieveSolutionImportResult%28ImportJobId=1b1fa80d-aa0f-4291-b60c-b0745304ce24%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Force_Recreate/get_solution_import_result.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=%28isvisible+eq+true%29&%24orderby=createdon+desc",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Force_Recreate/get_solution.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=uniquename+eq+%27TerraformTestSolution%27",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Force_Recreate/get_solution.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=solutionid+eq+86928ed8-df37-4ce2-add5-47030a833bff",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Force_Recreate/get_solution.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions%2886928ed8-df37-4ce2-add5-47030a833bff%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, httpmock.File("tests/resource/Validate_Create_And_Force_Recreate/get_solution.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
		
				resource "powerplatform_solution" "solution" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					solution_file    = "test_solution_before.zip"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_file_checksum", solution_before_checksum),
					resource.TestCheckNoResourceAttr("powerplatform_solution.solution", "settings_file_checksum"),
				),
			},
			{
				Config: `

				resource "powerplatform_solution" "solution" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					solution_file    = "test_solution_before.zip"
					settings_file 	 = "test_settings_before.json"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_file_checksum", solution_before_checksum),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "settings_file_checksum", settings_before_checksum),
				),
			},
			{
				Config: `

				resource "powerplatform_solution" "solution" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					solution_file    = "test_solution_after.zip"
					settings_file 	 = "test_settings_before.json"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_file_checksum", solution_after_checksum),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "settings_file_checksum", settings_before_checksum),
				),
			},
			{
				Config: `

				resource "powerplatform_solution" "solution" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					solution_file    = "test_solution_after.zip"
					settings_file 	 = "test_settings_after.json"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_file_checksum", solution_after_checksum),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "settings_file_checksum", settings_after_checksum),
				),
			},
		},
	})
}

func TestAccSolutionResource_Validate_Create_No_Dataverse(t *testing.T) {
	solutionFileBytes, err := os.ReadFile(SOLUTION_1_RELATIVE_PATH)
	if err != nil {
		t.Fatalf("Failed to read solution file: %v", err)
	}

	err = os.WriteFile(SOLUTION_1_NAME, solutionFileBytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write solution file: %v", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `

				resource "powerplatform_environment" "environment" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
				}
				
				resource "powerplatform_solution" "solution" {
					environment_id = powerplatform_environment.environment.id
					solution_file    = "` + SOLUTION_1_NAME + `"
				}`,
				ExpectError: regexp.MustCompile("No Dataverse exists in environment"),
				Check:       resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestUnitSolutionResource_Validate_Create_No_Dataverse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Dataverse/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/StageSolution",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Dataverse/post_stage_solution.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/ImportSolutionAsync",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Dataverse/post_import_solution_async.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/asyncoperations%28310799b8-dc6c-ee11-9ae7-000d3aaae21d%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Dataverse/get_async_operations.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/RetrieveSolutionImportResult%28ImportJobId=1b1fa80d-aa0f-4291-b60c-b0745304ce24%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Dataverse/get_solution_import_result.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=%28isvisible+eq+true%29&%24orderby=createdon+desc",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Dataverse/get_solution.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions%2886928ed8-df37-4ce2-add5-47030a833bff%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, httpmock.File("tests/resource/Validate_Create_No_Dataverse/get_solution.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Dataverse/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Create_No_Dataverse/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Dataverse/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `

				resource "powerplatform_environment" "env" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
				}

				resource "powerplatform_solution" "solution" {
					environment_id = powerplatform_environment.env.id
					solution_file    = "test_solution.zip"
				}`,
				ExpectError: regexp.MustCompile("No Dataverse exists in environment '00000000-0000-0000-0000-000000000001'"),
				Check:       resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func createFile(fileName string, content string) string {
	file, err := os.Create(fileName)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = file.Write([]byte(content))
	if err != nil {
		panic(err)
	}

	fileChecksum, _ := helpers.CalculateSHA256(fileName)
	return fileChecksum
}
