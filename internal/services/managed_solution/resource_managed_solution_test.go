// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package managedsolution_test

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccManagedSolutionResource_Validate_Create_HappyPath(t *testing.T) {
	solutionFileBytes, err := os.ReadFile("../solution/tests/resource/Test_Files/TerraformSimpleTestSolution_1_0_0_1_managed.zip")
	if err != nil {
		t.Fatalf("Failed to read solution file: %s", err.Error())
	}

	err = os.WriteFile("TerraformSimpleTestSolution_1_0_0_1_managed.zip", solutionFileBytes, 0644)
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
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
						language_code     = "1033"
						currency_code     = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "time_sleep" "wait_120_seconds" {
					depends_on = [powerplatform_environment.environment]
					create_duration = "120s"
				}

				resource "powerplatform_managed_solution" "solution" {
					depends_on     = [time_sleep.wait_120_seconds]
					environment_id = powerplatform_environment.environment.id
					unique_name    = "TerraformSimpleTestSolution"
					version        = "1.0.0.1"

					source = {
						path = "TerraformSimpleTestSolution_1_0_0_1_managed.zip"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "unique_name", "TerraformSimpleTestSolution"),
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "version", "1.0.0.1"),
				),
			},
		},
	})
}

func TestUnitManagedSolutionResource_Validate_Create_HappyPath(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	solutionPath, err := filepath.Abs("../solution/tests/resource/Test_Files/TerraformSimpleTestSolution_1_0_0_1_managed.zip")
	if err != nil {
		t.Fatalf("failed to resolve solution path: %v", err)
	}

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
  "id": "00000000-0000-0000-0000-000000000001",
  "name": "env",
  "properties": {
    "linkedEnvironmentMetadata": {
      "instanceURL": "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/"
    }
  }
}`), nil
		})

	httpmock.RegisterResponder("GET", "https://000000000000000000000000000000.01.environment.api.powerplatform.com/connectivity/connections?api-version=1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24orderby=createdon+desc",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/StageSolution",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
  "StageSolutionResults": {
    "StageSolutionUploadId": "upload-id",
    "StageSolutionStatus": "Passed",
    "SolutionValidationResults": [],
    "MissingDependencies": [],
    "SolutionDetails": {
      "SolutionUniqueName": "TerraformSimpleTestSolution"
    }
  }
}`), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/ImportSolutionAsync",
		func(req *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("failed to read import request body: %v", err)
			}

			requestBody := map[string]any{}
			if err := json.Unmarshal(body, &requestBody); err != nil {
				t.Fatalf("failed to unmarshal import request body: %v", err)
			}

			if _, exists := requestBody["ComponentParameters"]; exists {
				t.Fatalf("ComponentParameters should be omitted when the package declares no connection references: %s", string(body))
			}

			return httpmock.NewStringResponse(http.StatusOK, `{
  "ImportJobKey": "job-id",
  "AsyncOperationId": "async-id"
}`), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/asyncoperations%28async-id%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"completedon":"2024-01-01T00:00:00Z"}`), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/RetrieveSolutionImportResult%28ImportJobId=job-id%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
  "SolutionOperationResult": {
    "Status": "Passed",
    "ErrorMessages": []
  }
}`), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=uniquename+eq+%27TerraformSimpleTestSolution%27",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
  "value": [
    {
      "solutionid": "86928ed8-df37-4ce2-add5-47030a833bff",
      "environment_id": "00000000-0000-0000-0000-000000000001",
      "uniquename": "TerraformSimpleTestSolution",
      "friendlyname": "Terraform Simple Test Solution",
      "ismanaged": true,
      "createdon": "2024-01-01T00:00:00Z",
      "version": "1.0.0.1",
      "modifiedon": "2024-01-01T00:00:00Z",
      "installedon": "2024-01-01T00:00:00Z"
    }
  ]
}`), nil
		})

	httpmock.RegisterResponder("DELETE", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions%2886928ed8-df37-4ce2-add5-47030a833bff%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ``), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "powerplatform_managed_solution" "solution" {
  environment_id = "00000000-0000-0000-0000-000000000001"
  unique_name    = "TerraformSimpleTestSolution"
  version        = "1.0.0.1"

  source = {
    path = %q
  }
}
`, solutionPath),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "unique_name", "TerraformSimpleTestSolution"),
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "version", "1.0.0.1"),
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "display_name", "Terraform Simple Test Solution"),
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "solution_id", "86928ed8-df37-4ce2-add5-47030a833bff"),
				),
			},
		},
	})
}

func TestUnitManagedSolutionResource_AdoptsExactInstalledManagedSolution_AndImportsState(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()
	solutionPath := createTestSolutionZip(t, map[string]string{
		"solution.xml":       `<ImportExportXml><SolutionManifest><UniqueName>MetaForm</UniqueName><Version>2.0.246</Version><Managed>1</Managed><LocalizedNames><LocalizedName description="Meta Form" /></LocalizedNames></SolutionManifest></ImportExportXml>`,
		"customizations.xml": `<ImportExportXml></ImportExportXml>`,
	})
	registerManagedSolutionEnvironmentResponder()

	const installed = `{"value":[{"solutionid":"86928ed8-df37-4ce2-add5-47030a833bff","uniquename":"MetaForm","friendlyname":"Meta Form","ismanaged":true,"version":"2.0.246.0"}]}`
	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24orderby=createdon+desc",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, installed), nil
		})
	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=uniquename+eq+%27MetaForm%27",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, installed), nil
		})
	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=solutionid+eq+86928ed8-df37-4ce2-add5-47030a833bff",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, installed), nil
		})
	httpmock.RegisterResponder("GET", "https://000000000000000000000000000000.01.environment.api.powerplatform.com/connectivity/connections?api-version=1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
		})
	importStarts := 0
	httpmock.RegisterResponder("POST", `=~^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.2/(StageSolution|ImportSolutionAsync|StageAndUpgradeAsync)$`,
		func(req *http.Request) (*http.Response, error) {
			importStarts++
			return httpmock.NewStringResponse(http.StatusInternalServerError, "existing exact solution must be adopted"), nil
		})
	httpmock.RegisterResponder("DELETE", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions%2886928ed8-df37-4ce2-add5-47030a833bff%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	config := fmt.Sprintf(`
resource "powerplatform_managed_solution" "solution" {
  environment_id = "00000000-0000-0000-0000-000000000001"
  unique_name    = "MetaForm"
  version        = "2.0.246"

  source = {
    path = %q
  }
}
`, solutionPath)
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "id", "00000000-0000-0000-0000-000000000001_86928ed8-df37-4ce2-add5-47030a833bff"),
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "solution_id", "86928ed8-df37-4ce2-add5-47030a833bff"),
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "version", "2.0.246"),
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "skip_product_update_dependencies", "true"),
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "publish_all_customizations", "false"),
				),
			},
			{
				ResourceName:      "powerplatform_managed_solution.solution",
				ImportState:       true,
				ImportStateId:     "00000000-0000-0000-0000-000000000001_86928ed8-df37-4ce2-add5-47030a833bff",
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"source",
					"version",
					"skip_product_update_dependencies",
					"publish_all_customizations",
				},
			},
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "version", "2.0.246"),
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "source.path", solutionPath),
				),
			},
		},
	})
	if importStarts != 0 {
		t.Fatalf("exact installed managed solution should be adopted without an import; observed %d import request(s)", importStarts)
	}
}

func TestUnitManagedSolutionResource_Validate_Create_Fails_When_ConnectionReferenceMissing(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	solutionPath := createTestSolutionZip(t, map[string]string{
		"solution.xml":       `<ImportExportXml><SolutionManifest><UniqueName>CodeEditor</UniqueName><Version>1.0.0.0</Version><Managed>1</Managed><LocalizedNames><LocalizedName description="Code Editor" /></LocalizedNames></SolutionManifest></ImportExportXml>`,
		"customizations.xml": `<ImportExportXml><connectionreferences><connectionreference connectionreferencelogicalname="codeeditor_sharepoint"><connectorid>/providers/Microsoft.PowerApps/apis/shared_sharepointonline</connectorid></connectionreference></connectionreferences></ImportExportXml>`,
	})

	registerManagedSolutionEnvironmentResponder()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "powerplatform_managed_solution" "solution" {
  environment_id = "00000000-0000-0000-0000-000000000001"
  unique_name    = "CodeEditor"
  version        = "1.0.0.0"

  source = {
    path = %q
  }
}
`, solutionPath),
				ExpectError: regexp.MustCompile(`(?s)missing required connection\s+references: codeeditor_sharepoint`),
			},
		},
	})
}

func TestUnitManagedSolutionResource_Validate_Create_Fails_When_EnvironmentVariableIsUnsatisfied(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	solutionPath := createTestSolutionZip(t, map[string]string{
		"solution.xml":       `<ImportExportXml><SolutionManifest><UniqueName>CodeEditor</UniqueName><Version>1.0.0.0</Version><Managed>1</Managed><LocalizedNames><LocalizedName description="Code Editor" /></LocalizedNames></SolutionManifest></ImportExportXml>`,
		"customizations.xml": `<ImportExportXml></ImportExportXml>`,
		"environmentvariabledefinitions/codeeditor_secret/environmentvariabledefinition.xml": `<environmentvariabledefinition schemaname="codeeditor_secret"></environmentvariabledefinition>`,
	})

	registerManagedSolutionEnvironmentResponder()
	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/environmentvariablevalues?%24filter=schemaname+eq+%27codeeditor_secret%27&%24select=environmentvariablevalueid%2Cschemaname%2Cvalue",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "powerplatform_managed_solution" "solution" {
  environment_id = "00000000-0000-0000-0000-000000000001"
  unique_name    = "CodeEditor"
  version        = "1.0.0.0"

  source = {
    path = %q
  }
}
`, solutionPath),
				ExpectError: regexp.MustCompile(`(?s)environment variable validation failed: codeeditor_secret has no packaged\s+default and no existing environment value`),
			},
		},
	})
}

func TestUnitManagedSolutionResource_Validate_Create_Fails_When_DependencyVersionIsTooLow(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	solutionPath := createTestSolutionZip(t, map[string]string{
		"solution.xml":       `<ImportExportXml><SolutionManifest><UniqueName>CodeEditor</UniqueName><Version>1.0.0.0</Version><Managed>1</Managed><LocalizedNames><LocalizedName description="Code Editor" /></LocalizedNames><MissingDependencies><MissingDependency><Required type="66" schemaName="base" solution="BaseLib (2.0.0.0)" /></MissingDependency></MissingDependencies></SolutionManifest></ImportExportXml>`,
		"customizations.xml": `<ImportExportXml></ImportExportXml>`,
	})

	registerManagedSolutionEnvironmentResponder()
	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24orderby=createdon+desc",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
  "value": [
    {
      "solutionid": "base-id",
      "uniquename": "BaseLib",
      "friendlyname": "Base Lib",
      "ismanaged": true,
      "createdon": "2024-01-01T00:00:00Z",
      "version": "1.0.0.0",
      "modifiedon": "2024-01-01T00:00:00Z",
      "installedon": "2024-01-01T00:00:00Z"
    }
  ]
}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "powerplatform_managed_solution" "solution" {
  environment_id = "00000000-0000-0000-0000-000000000001"
  unique_name    = "CodeEditor"
  version        = "1.0.0.0"

  source = {
    path = %q
  }
}
`, solutionPath),
				ExpectError: regexp.MustCompile(`(?s)dependency validation failed: required dependency "BaseLib" >= 2.0.0.0 but\s+1.0.0.0 is installed`),
			},
		},
	})
}

func registerManagedSolutionEnvironmentResponder() {
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
  "id": "00000000-0000-0000-0000-000000000001",
  "name": "env",
  "properties": {
    "linkedEnvironmentMetadata": {
      "instanceURL": "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/"
    }
  }
}`), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24orderby=createdon+desc",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
		})
}

func createTestSolutionZip(t *testing.T, files map[string]string) string {
	t.Helper()

	file, err := os.CreateTemp("", "managed-solution-resource-test-*.zip")
	if err != nil {
		t.Fatalf("failed to create temp zip: %v", err)
	}
	defer file.Close()

	writer := zip.NewWriter(file)
	for name, content := range files {
		entry, err := writer.Create(name)
		if err != nil {
			t.Fatalf("failed to create zip entry %s: %v", name, err)
		}
		if _, err := entry.Write([]byte(content)); err != nil {
			t.Fatalf("failed to write zip entry %s: %v", name, err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close zip writer: %v", err)
	}

	return file.Name()
}
