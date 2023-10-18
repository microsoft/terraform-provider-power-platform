package powerplatform

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccSolutionResource_Validate_Create_No_Settings_File(t *testing.T) {

	solutionName := "TerraformTestSolution"
	solutionFileName := solutionName + "_Complex_1_1_0_0.zip"
	rand.Seed(time.Now().UnixNano())
	envDomain := fmt.Sprintf("orgtest%d", rand.Intn(100000))

	solutionFileBytes, err := os.ReadFile(filepath.Join("../../examples/resources/powerplatform_solution", solutionFileName))
	if err != nil {
		t.Fatalf("Failed to read solution file: %v", err)
	}

	err = os.WriteFile(solutionFileName, solutionFileBytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write solution file: %v", err)
	}

	solutionFileChecksum, _ := powerplatform_helpers.CalculateMd5(solutionFileName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + `

				resource "powerplatform_environment" "environment" {
					display_name                              = "` + envDomain + `"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                           = "USD"
					environment_type                          = "Sandbox"
					domain = "` + envDomain + `"
					security_group_id = "00000000-0000-0000-0000-000000000000"
					environment_id = "00000000-0000-0000-0000-000000000000"
				}

				resource "powerplatform_solution" "solution" {
					environment_id = powerplatform_environment.environment.id
					solution_name    = "` + solutionName + `"
					solution_file    = "` + solutionFileName + `"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("powerplatform_solution.solution", "settings_file_checksum"),
					resource.TestCheckNoResourceAttr("powerplatform_solution.solution", "settings_file"),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_name", solutionName),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_file_checksum", solutionFileChecksum),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_file", solutionFileName),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "display_name", "Terraform Test Solution"),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "is_managed", "false"),
					resource.TestMatchResourceAttr("powerplatform_solution.solution", "environment_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_version", "1.1.0.0"),
				),
			},
		},
	})
}

func TestUnitSolutionResource_Validate_Create_With_Settings_File(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock_helpers.ActivateOAuthHttpMocks()
	mock_helpers.ActivateEnvironmentHttpMocks()

	solution_checksum := createFile("test_solution.zip", "test_solution")
	settings_checksum := createFile("test_solution_settings.json", "")

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource_solution_test/Validate_Create_With_Settings_File/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/StageSolution",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource_solution_test/Validate_Create_With_Settings_File/post_stage_solution.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/ImportSolutionAsync",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource_solution_test/Validate_Create_With_Settings_File/post_import_solution_async.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/asyncoperations%28310799b8-dc6c-ee11-9ae7-000d3aaae21d%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource_solution_test/Validate_Create_With_Settings_File/get_async_operations.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/RetrieveSolutionImportResult%28ImportJobId=1b1fa80d-aa0f-4291-b60c-b0745304ce24%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource_solution_test/Validate_Create_With_Settings_File/get_solution_import_result.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=%28isvisible+eq+true%29&%24orderby=createdon+desc",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource_solution_test/Validate_Create_With_Settings_File/get_solution.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions%2886928ed8-df37-4ce2-add5-47030a833bff%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, httpmock.File("tests/resource_solution_test/Validate_Create_With_Settings_File/get_solution.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: UnitTestsProviderConfig + `

				resource "powerplatform_solution" "solution" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					solution_name    = "TerraformTestSolution"
					solution_file    = "test_solution.zip"
					settings_file 	 = "test_solution_settings.json"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_name", "TerraformTestSolution"),
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

	solutionName := "TerraformTestSolution"
	solutionFileName := solutionName + "_Complex_1_1_0_0.zip"
	solutionSettingsFileName := "test_solution_settings.json"
	rand.Seed(time.Now().UnixNano())
	envDomain := fmt.Sprintf("orgtest%d", rand.Intn(100000))

	solutionFileBytes, err := os.ReadFile(filepath.Join("../../examples/resources/powerplatform_solution", solutionFileName))
	if err != nil {
		t.Fatalf("Failed to read solution file: %v", err)
	}

	err = os.WriteFile(solutionFileName, solutionFileBytes, 0644)
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

	solutionFileChecksum, _ := powerplatform_helpers.CalculateMd5(solutionFileName)
	settingsFileChecksum, _ := powerplatform_helpers.CalculateMd5(solutionSettingsFileName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + `

				resource "powerplatform_environment" "environment" {
					display_name                              = "TestAccSolutionResource_Settings_File"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                         = "USD"
					environment_type                          = "Sandbox"
					domain 									  = "` + envDomain + `"
					security_group_id = "00000000-0000-0000-0000-000000000000"
					environment_id = "00000000-0000-0000-0000-000000000000"
				}

				resource "powerplatform_solution" "solution" {
					environment_id = powerplatform_environment.environment.id
					solution_name    = "TerraformTestSolution"
					solution_file    = "` + solutionFileName + `"
					settings_file 	 = "` + solutionSettingsFileName + `"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_name", solutionName),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_file_checksum", solutionFileChecksum),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "settings_file_checksum", settingsFileChecksum),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "settings_file", solutionSettingsFileName),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_file", solutionFileName),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "display_name", "Terraform Test Solution"),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "is_managed", "false"),
					resource.TestMatchResourceAttr("powerplatform_solution.solution", "environment_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_version", "1.1.0.0"),
				),
			},
		},
	})
}

func TestUnitSolutionResource_Validate_Create_With_Settings_File(t *testing.T) {
	clientDataverseMock := mocks.NewUnitTestMockDataverseClientInterface(t)
	clientBapiMock := mocks.NewUnitTestsMockBapiClientInterface(t)

	solutionFileName := "test_solution.zip"
	solutionSettingsFileName := "test_solution_settings.json"

	solutionFileChecksum := createFile(solutionFileName, "")
	solutionSettingsFileChecksum := createFile(solutionSettingsFileName, "")

	environmentStub := models.EnvironmentDto{
		Id: "00000000-0000-0000-0000-000000000000",
		Properties: models.EnvironmentPropertiesDto{
			LinkedEnvironmentMetadata: models.LinkedEnvironmentMetadataDto{
				ResourceId:      "org1",
				SecurityGroupId: "security1",
				DomainName:      "domain",
				InstanceURL:     "url",
				Version:         "version",
			},
			LinkedAppMetadata: models.LinkedAppMetadataDto{
				Type: "Internal",
				Id:   "00000000-0000-0000-0000-000000000000",
				Url:  "https://url.operations.dynamics.com/",
			},
		},
	}

	solutionStub := models.SolutionDto{
		Id:            "00000000-0000-0000-0000-000000000002",
		EnvironmentId: environmentStub.Id,
		DisplayName:   "Solution",
		Name:          "solution",
		CreatedTime:   "2020-01-01T00:00:00Z",
		ModifiedTime:  "2020-01-01T00:00:00Z",
		InstallTime:   "2020-01-01T00:00:00Z",
		Version:       "1.2.3.4",
		IsManaged:     true,
	}

	clientDataverseMock.EXPECT().GetDefaultCurrencyForEnvironment(gomock.Any(), gomock.Any()).Return(&models.TransactionCurrencyDto{IsoCurrencyCode: "USD"}, nil).AnyTimes()

	clientBapiMock.EXPECT().CreateEnvironment(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, envToCreate models.EnvironmentCreateDto) (*models.EnvironmentDto, error) {
		environmentStub = models.EnvironmentDto{
			Name:     environmentStub.Name,
			Id:       environmentStub.Id,
			Location: envToCreate.Location,
			Properties: models.EnvironmentPropertiesDto{
				DisplayName:    envToCreate.Properties.DisplayName,
				EnvironmentSku: envToCreate.Properties.EnvironmentSku,
				LinkedEnvironmentMetadata: models.LinkedEnvironmentMetadataDto{
					DomainName:      envToCreate.Properties.LinkedEnvironmentMetadata.DomainName,
					BaseLanguage:    envToCreate.Properties.LinkedEnvironmentMetadata.BaseLanguage,
					SecurityGroupId: envToCreate.Properties.LinkedEnvironmentMetadata.SecurityGroupId,
					InstanceURL:     environmentStub.Properties.LinkedEnvironmentMetadata.InstanceURL,
					Version:         environmentStub.Properties.LinkedEnvironmentMetadata.Version,
					ResourceId:      environmentStub.Properties.LinkedEnvironmentMetadata.ResourceId,
				},
			},
		}
		return &environmentStub, nil
	}).Times(1)

	clientBapiMock.EXPECT().GetEnvironment(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, id string) (*models.EnvironmentDto, error) {
		return &environmentStub, nil
	}).AnyTimes()

	clientBapiMock.EXPECT().DeleteEnvironment(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	clientDataverseMock.EXPECT().GetSolutions(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, environmentId string) ([]models.SolutionDto, error) {
		solutions := []models.SolutionDto{}
		solutions = append(solutions, solutionStub)

		return solutions, nil
	}).AnyTimes()

	clientDataverseMock.EXPECT().CreateSolution(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, EnvironmentId string, solutionToCreate models.ImportSolutionDto, content []byte, settings []byte) (*models.SolutionDto, error) {
		return &solutionStub, nil
	}).Times(1)

	clientDataverseMock.EXPECT().DeleteSolution(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"powerplatform": powerPlatformProviderServerApiMock(clientBapiMock, clientDataverseMock, nil),
		},
		Steps: []resource.TestStep{
			{
				Config: UnitTestsProviderConfig + `

				resource "powerplatform_environment" "environment" {
					display_name                              = "Solution Import Acceptance Test"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "USD"
					environment_type                          = "Sandbox"
					security_group_id                         = "00000000-0000-0000-0000-000000000000"
					domain 								  	  = "domain"
				}

				resource "powerplatform_solution" "solution" {
					environment_id = "` + environmentStub.Id + `"
					solution_name    = "` + solutionStub.Name + `"
					solution_file    = "` + solutionFileName + `"
					settings_file 	 = "` + solutionSettingsFileName + `"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.environment", "id", environmentStub.Id),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_name", solutionStub.Name),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_version", solutionStub.Version),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_file_checksum", solutionFileChecksum),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "settings_file_checksum", solutionSettingsFileChecksum),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "display_name", solutionStub.DisplayName),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "is_managed", strconv.FormatBool(solutionStub.IsManaged)),
				),
			},
		},
	})
}

func TestUnitSolutionResource_Validate_Create_No_Settings_File(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock_helpers.ActivateOAuthHttpMocks()
	mock_helpers.ActivateEnvironmentHttpMocks()

	solution_checksum := createFile("test_solution.zip", "test_solution")

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource_solution_test/Validate_Create_No_Settings_File/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/StageSolution",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource_solution_test/Validate_Create_No_Settings_File/post_stage_solution.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/ImportSolutionAsync",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource_solution_test/Validate_Create_No_Settings_File/post_import_solution_async.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/asyncoperations%28310799b8-dc6c-ee11-9ae7-000d3aaae21d%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource_solution_test/Validate_Create_No_Settings_File/get_async_operations.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/RetrieveSolutionImportResult%28ImportJobId=1b1fa80d-aa0f-4291-b60c-b0745304ce24%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource_solution_test/Validate_Create_No_Settings_File/get_solution_import_result.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=%28isvisible+eq+true%29&%24orderby=createdon+desc",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource_solution_test/Validate_Create_No_Settings_File/get_solution.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions%2886928ed8-df37-4ce2-add5-47030a833bff%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, httpmock.File("tests/resource_solution_test/Validate_Create_No_Settings_File/get_solution.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: UnitTestsProviderConfig + `

				resource "powerplatform_solution" "solution" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					solution_name    = "TerraformTestSolution"
					solution_file    = "test_solution.zip"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("powerplatform_solution.solution", "settings_file_checksum"),
					resource.TestCheckNoResourceAttr("powerplatform_solution.solution", "settings_file"),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_file_checksum", solution_checksum),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_name", "TerraformTestSolution"),
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
	mock_helpers.ActivateOAuthHttpMocks()
	mock_helpers.ActivateEnvironmentHttpMocks()

	solution_before_checksum := createFile("test_solution_before.zip", "test_solution_before")
	settings_before_checksum := createFile("test_settings_before.json", "")
	solution_after_checksum := createFile("test_solution_after.zip", "test_solution_after")
	settings_after_checksum := createFile("test_settings_after.json", "")

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock_helpers.ActivateOAuthHttpMocks()
	mock_helpers.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource_solution_test/Validate_Create_And_Force_Recreate/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/StageSolution",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource_solution_test/Validate_Create_And_Force_Recreate/post_stage_solution.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/ImportSolutionAsync",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource_solution_test/Validate_Create_And_Force_Recreate/post_import_solution_async.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/asyncoperations%28310799b8-dc6c-ee11-9ae7-000d3aaae21d%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource_solution_test/Validate_Create_And_Force_Recreate/get_async_operations.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/RetrieveSolutionImportResult%28ImportJobId=1b1fa80d-aa0f-4291-b60c-b0745304ce24%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource_solution_test/Validate_Create_And_Force_Recreate/get_solution_import_result.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=%28isvisible+eq+true%29&%24orderby=createdon+desc",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource_solution_test/Validate_Create_And_Force_Recreate/get_solution.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions%2886928ed8-df37-4ce2-add5-47030a833bff%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, httpmock.File("tests/resource_solution_test/Validate_Create_And_Force_Recreate/get_solution.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: UnitTestsProviderConfig + `
		
				resource "powerplatform_solution" "solution" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					solution_name    = "TerraformTestSolution"
					solution_file    = "test_solution_before.zip"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_name", "TerraformTestSolution"),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_file_checksum", solution_before_checksum),
					resource.TestCheckNoResourceAttr("powerplatform_solution.solution", "settings_file_checksum"),
				),
			},
			{
				Config: UnitTestsProviderConfig + `

				resource "powerplatform_solution" "solution" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					solution_name    = "TerraformTestSolution"
					solution_file    = "test_solution_before.zip"
					settings_file 	 = "test_settings_before.json"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_name", "TerraformTestSolution"),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_file_checksum", solution_before_checksum),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "settings_file_checksum", settings_before_checksum),
				),
			},
			{
				Config: UnitTestsProviderConfig + `

				resource "powerplatform_solution" "solution" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					solution_name    = "TerraformTestSolution"
					solution_file    = "test_solution_after.zip"
					settings_file 	 = "test_settings_before.json"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_name", "TerraformTestSolution"),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_file_checksum", solution_after_checksum),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "settings_file_checksum", settings_before_checksum),
				),
			},
			{
				Config: UnitTestsProviderConfig + `

				resource "powerplatform_solution" "solution" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					solution_name    = "TerraformTestSolution"
					solution_file    = "test_solution_after.zip"
					settings_file 	 = "test_settings_after.json"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_name", "TerraformTestSolution"),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "solution_file_checksum", solution_after_checksum),
					resource.TestCheckResourceAttr("powerplatform_solution.solution", "settings_file_checksum", settings_after_checksum),
				),
			},
		},
	})
}

func createFile(fileName string, content string) string {
	file, _ := os.Create(fileName)
	file.Write([]byte(content))
	file.Close()
	fileChecksum, _ := powerplatform_helpers.CalculateMd5(fileName)
	return fileChecksum
}
