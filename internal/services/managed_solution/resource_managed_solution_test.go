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
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccManagedSolutionResource_Validate_Create_HappyPath(t *testing.T) {
	solutionPath, err := filepath.Abs("tests/resource/TerraformSolutionExample_1_0_0_1_managed.zip")
	if err != nil {
		t.Fatalf("failed to resolve solution path: %v", err)
	}

	guid := strings.Trim(helpers.GuidRegex, "^$")
	idRegex := regexp.MustCompile(fmt.Sprintf(`^%s/%s$`, guid, guid))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source: "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "powerplatform_environment" "environment" {
  display_name     = "%s"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

# A freshly provisioned Dataverse organization rejects solution imports with
# "Async operations are currently disabled for this organization" for a few minutes.
resource "time_sleep" "wait_for_dataverse" {
  create_duration = "240s"

  depends_on = [powerplatform_environment.environment]
}

resource "powerplatform_connection" "dataverse_connection" {
  environment_id = powerplatform_environment.environment.id
  name           = "shared_commondataserviceforapps"
  display_name   = "Dataverse Connection"

  connection_parameters = jsonencode({
  })

  lifecycle {
    ignore_changes = [
      connection_parameters
    ]
  }

  depends_on = [time_sleep.wait_for_dataverse]
}

resource "powerplatform_managed_solution" "solution" {
  environment_id = powerplatform_environment.environment.id
  unique_name    = "TerraformSolutionExample"
  version        = "1.0.0.1"

  source = {
    path = %q
  }

  connection_references = {
    terr_SolutionConnectionReference = powerplatform_connection.dataverse_connection.id
  }

  depends_on = [powerplatform_connection.dataverse_connection]
}
`, mocks.TestName(), solutionPath),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "unique_name", "TerraformSolutionExample"),
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "version", "1.0.0.1"),
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "display_name", "Terraform Solution Example"),
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "source.path", solutionPath),
					resource.TestCheckNoResourceAttr("powerplatform_managed_solution.solution", "source.url"),
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "skip_product_update_dependencies", "true"),
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "publish_all_customizations", "false"),
					resource.TestMatchResourceAttr("powerplatform_managed_solution.solution", "solution_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_managed_solution.solution", "environment_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_managed_solution.solution", "id", idRegex),
					resource.TestCheckResourceAttrPair("powerplatform_managed_solution.solution", "environment_id", "powerplatform_environment.environment", "id"),
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "connection_references.%", "1"),
					resource.TestCheckResourceAttrPair("powerplatform_managed_solution.solution", "connection_references.terr_SolutionConnectionReference", "powerplatform_connection.dataverse_connection", "id"),
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
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_HappyPath/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://000000000000000000000000000000.01.environment.api.powerplatform.com/connectivity/connections?api-version=1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_HappyPath/get_connections.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24orderby=createdon+desc",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_HappyPath/get_installed_solutions.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/StageSolution",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_HappyPath/post_stage_solution.json").String()), nil
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

			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_HappyPath/post_import_solution_async.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/asyncoperations%28async-id%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_HappyPath/get_async_operations.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/RetrieveSolutionImportResult%28ImportJobId=job-id%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_HappyPath/get_solution_import_result.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=uniquename+eq+%27TerraformSimpleTestSolution%27",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_HappyPath/get_solution.json").String()), nil
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
	registerManagedSolutionEnvironmentResponder("AdoptsExactInstalledManagedSolution_AndImportsState")

	const installedSolutionFile = "tests/resource/AdoptsExactInstalledManagedSolution_AndImportsState/get_solution.json"
	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24orderby=createdon+desc",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(installedSolutionFile).String()), nil
		})
	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=uniquename+eq+%27MetaForm%27",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(installedSolutionFile).String()), nil
		})
	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=solutionid+eq+86928ed8-df37-4ce2-add5-47030a833bff",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(installedSolutionFile).String()), nil
		})
	httpmock.RegisterResponder("GET", "https://000000000000000000000000000000.01.environment.api.powerplatform.com/connectivity/connections?api-version=1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/AdoptsExactInstalledManagedSolution_AndImportsState/get_connections.json").String()), nil
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
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "id", "00000000-0000-0000-0000-000000000001/86928ed8-df37-4ce2-add5-47030a833bff"),
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "solution_id", "86928ed8-df37-4ce2-add5-47030a833bff"),
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "version", "2.0.246"),
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "skip_product_update_dependencies", "true"),
					resource.TestCheckResourceAttr("powerplatform_managed_solution.solution", "publish_all_customizations", "false"),
				),
			},
			{
				ResourceName:      "powerplatform_managed_solution.solution",
				ImportState:       true,
				ImportStateId:     "00000000-0000-0000-0000-000000000001/86928ed8-df37-4ce2-add5-47030a833bff",
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

	registerManagedSolutionEnvironmentResponder("Validate_Create_Fails_When_ConnectionReferenceMissing")

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

func TestUnitManagedSolutionResource_Validate_Create_Fails_When_UnmanagedDefinitionWouldBeCaptured(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	solutionPath := createTestSolutionZip(t, map[string]string{
		"solution.xml":       `<ImportExportXml><SolutionManifest><UniqueName>CodeEditor</UniqueName><Version>1.0.0.0</Version><Managed>1</Managed><LocalizedNames><LocalizedName description="Code Editor" /></LocalizedNames></SolutionManifest></ImportExportXml>`,
		"customizations.xml": `<ImportExportXml></ImportExportXml>`,
		"environmentvariabledefinitions/codeeditor_secret/environmentvariabledefinition.xml": `<environmentvariabledefinition schemaname="codeeditor_secret"></environmentvariabledefinition>`,
	})

	registerManagedSolutionEnvironmentResponder("Validate_Create_Fails_When_UnmanagedDefinitionCapture")
	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/environmentvariabledefinitions?%24filter=schemaname+eq+%27codeeditor_secret%27&%24select=schemaname%2Cismanaged",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Fails_When_UnmanagedDefinitionCapture/get_environment_variable_definitions.json").String()), nil
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
				ExpectError: regexp.MustCompile(`(?s)environment variable packaging validation failed: codeeditor_secret already\s+exists as an unmanaged definition`),
			},
		},
	})
}

func TestUnitManagedSolutionResource_Validate_Create_Retries_When_Dataverse_Denies_Caller(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	solutionPath := createTestSolutionZip(t, map[string]string{
		"solution.xml":       `<ImportExportXml><SolutionManifest><UniqueName>CodeEditor</UniqueName><Version>1.0.0.0</Version><Managed>1</Managed><LocalizedNames><LocalizedName description="Code Editor" /></LocalizedNames></SolutionManifest></ImportExportXml>`,
		"customizations.xml": `<ImportExportXml></ImportExportXml>`,
		"environmentvariabledefinitions/codeeditor_secret/environmentvariabledefinition.xml": `<environmentvariabledefinition schemaname="codeeditor_secret"></environmentvariabledefinition>`,
	})

	registerManagedSolutionEnvironmentResponder("Validate_Create_Retries_When_Dataverse_Denies_Caller")
	definitionAttempts := 0
	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/environmentvariabledefinitions?%24filter=schemaname+eq+%27codeeditor_secret%27&%24select=schemaname%2Cismanaged",
		func(req *http.Request) (*http.Response, error) {
			definitionAttempts++
			if definitionAttempts == 1 {
				return httpmock.NewStringResponse(http.StatusForbidden, `{"error":{"code":"0x80072560","message":"The user is not a member of the organization."}}`), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Retries_When_Dataverse_Denies_Caller/get_environment_variable_definitions.json").String()), nil
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
				// The packaging error proves the retried read reached Dataverse after the 403.
				ExpectError: regexp.MustCompile(`(?s)environment variable packaging validation failed: codeeditor_secret already\s+exists as an unmanaged definition`),
			},
		},
	})

	if definitionAttempts < 2 {
		t.Fatalf("expected the environment variable definition read to be retried after a 403; observed %d attempt(s)", definitionAttempts)
	}
}

func TestUnitManagedSolutionResource_Validate_Create_Fails_When_DependencyVersionIsTooLow(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	solutionPath := createTestSolutionZip(t, map[string]string{
		"solution.xml":       `<ImportExportXml><SolutionManifest><UniqueName>CodeEditor</UniqueName><Version>1.0.0.0</Version><Managed>1</Managed><LocalizedNames><LocalizedName description="Code Editor" /></LocalizedNames><MissingDependencies><MissingDependency><Required type="66" schemaName="base" solution="BaseLib (2.0.0.0)" /></MissingDependency></MissingDependencies></SolutionManifest></ImportExportXml>`,
		"customizations.xml": `<ImportExportXml></ImportExportXml>`,
	})

	registerManagedSolutionEnvironmentResponder("Validate_Create_Fails_When_DependencyVersionIsTooLow")
	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24orderby=createdon+desc",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Fails_When_DependencyVersionIsTooLow/get_installed_solutions.json").String()), nil
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

func registerManagedSolutionEnvironmentResponder(testFolder string) {
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/%s/get_environment_00000000-0000-0000-0000-000000000001.json", testFolder)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24orderby=createdon+desc",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/%s/get_installed_solutions.json", testFolder)).String()), nil
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
