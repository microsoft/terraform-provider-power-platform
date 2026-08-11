// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environmentvariable_test

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccEnvironmentVariableResource_Validate_Create_HappyPath(t *testing.T) {
	solutionPath := filepath.Join("../solution/tests/resource/Test_Files", "TerraformTestSolution_Complex_1_1_0_0.zip")
	solutionBytes, err := os.ReadFile(solutionPath)
	if err != nil {
		t.Fatalf("failed to read solution file: %s", err.Error())
	}

	const solutionFileName = "TerraformTestSolution_Complex_1_1_0_0.zip"
	err = os.WriteFile(solutionFileName, solutionBytes, 0644)
	if err != nil {
		t.Fatalf("failed to write solution file: %s", err.Error())
	}
	t.Cleanup(func() {
		_ = os.Remove(solutionFileName)
	})

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
  depends_on      = [powerplatform_environment.environment]
  create_duration = "120s"
}

resource "powerplatform_solution" "solution" {
  depends_on     = [time_sleep.wait_120_seconds]
  environment_id = powerplatform_environment.environment.id
  solution_file  = "` + solutionFileName + `"
}

resource "powerplatform_environment_variable" "text" {
  depends_on     = [powerplatform_solution.solution]
  environment_id = powerplatform_environment.environment.id
  schema_name    = "cra6e_SolutionVariableText"
  value          = "https://prod.contoso.example"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_variable.text", "schema_name", "cra6e_SolutionVariableText"),
					resource.TestCheckResourceAttr("powerplatform_environment_variable.text", "display_name", "SolutionVariableText"),
					resource.TestCheckResourceAttr("powerplatform_environment_variable.text", "type", "String"),
					resource.TestCheckResourceAttr("powerplatform_environment_variable.text", "secret_store", "Microsoft Dataverse"),
					resource.TestMatchResourceAttr("powerplatform_environment_variable.text", "environment_variable_definition_id", regexp.MustCompile(`^[0-9a-fA-F-]{36}$`)),
					resource.TestMatchResourceAttr("powerplatform_environment_variable.text", "environment_variable_value_id", regexp.MustCompile(`^[0-9a-fA-F-]{36}$`)),
				),
			},
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
  depends_on      = [powerplatform_environment.environment]
  create_duration = "120s"
}

resource "powerplatform_solution" "solution" {
  depends_on     = [time_sleep.wait_120_seconds]
  environment_id = powerplatform_environment.environment.id
  solution_file  = "` + solutionFileName + `"
}

resource "powerplatform_environment_variable" "text" {
  depends_on     = [powerplatform_solution.solution]
  environment_id = powerplatform_environment.environment.id
  schema_name    = "cra6e_SolutionVariableText"
  value          = "https://test.contoso.example"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_variable.text", "schema_name", "cra6e_SolutionVariableText"),
				),
			},
		},
	})
}

func TestUnitEnvironmentVariableResource_Validate_Create_HappyPath(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerEnvironmentVariableHappyPathResponders()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "powerplatform_environment_variable" "text" {
  environment_id = "00000000-0000-0000-0000-000000000001"
  schema_name    = "contoso_ApiBaseUrl"
  value          = "https://api.contoso.example"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_variable.text", "id", "00000000-0000-0000-0000-000000000001/contoso_ApiBaseUrl"),
					resource.TestCheckResourceAttr("powerplatform_environment_variable.text", "schema_name", "contoso_ApiBaseUrl"),
					resource.TestCheckResourceAttr("powerplatform_environment_variable.text", "display_name", "API Base URL"),
					resource.TestCheckResourceAttr("powerplatform_environment_variable.text", "description", "Base URL for downstream API"),
					resource.TestCheckResourceAttr("powerplatform_environment_variable.text", "type", "String"),
					resource.TestCheckResourceAttr("powerplatform_environment_variable.text", "secret_store", "Microsoft Dataverse"),
					resource.TestCheckResourceAttr("powerplatform_environment_variable.text", "value_schema", "{\"type\":\"string\"}"),
					resource.TestCheckResourceAttr("powerplatform_environment_variable.text", "environment_variable_definition_id", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckResourceAttr("powerplatform_environment_variable.text", "environment_variable_value_id", "22222222-2222-2222-2222-222222222222"),
				),
			},
		},
	})
}

func TestUnitEnvironmentVariableResource_Validate_Update_HappyPath(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerEnvironmentVariableUpdateResponders()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "powerplatform_environment_variable" "text" {
  environment_id = "00000000-0000-0000-0000-000000000001"
  schema_name    = "contoso_ApiBaseUrl"
  value          = "https://api.contoso.example"
}
`,
			},
			{
				Config: `
resource "powerplatform_environment_variable" "text" {
  environment_id = "00000000-0000-0000-0000-000000000001"
  schema_name    = "contoso_ApiBaseUrl"
  value          = "https://api2.contoso.example"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_variable.text", "environment_variable_value_id", "22222222-2222-2222-2222-222222222222"),
				),
			},
		},
	})
}

func TestUnitEnvironmentVariableResource_Validate_Create_Fails_When_Definition_Is_Missing(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerEnvironmentResponder()
	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/environmentvariabledefinitions?%24filter=schemaname+eq+%27contoso_Missing%27&%24select=environmentvariabledefinitionid%2Cschemaname%2Cdisplayname%2Cdescription%2Cdefaultvalue%2Ctype%2Cvalueschema%2Csecretstore",
		httpmock.NewStringResponder(http.StatusOK, `{"value":[]}`))

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "powerplatform_environment_variable" "missing" {
  environment_id = "00000000-0000-0000-0000-000000000001"
  schema_name    = "contoso_Missing"
  value          = "x"
}
`,
				ExpectError: regexp.MustCompile(`contoso_Missing'.*not found`),
			},
		},
	})
}

func TestUnitEnvironmentVariableResource_Validate_Create_Fails_When_No_Dataverse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		httpmock.NewStringResponder(http.StatusOK, `{
  "name":"00000000-0000-0000-0000-000000000001",
  "id":"/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001",
  "type":"Microsoft.BusinessAppPlatform/scopes/admin/environments",
  "location":"unitedstates",
  "properties":{
    "displayName":"No Dataverse",
    "azureRegion":"unitedstates",
    "createdTime":"2024-01-01T00:00:00Z",
    "environmentSku":"Sandbox",
    "linkedEnvironmentMetadata":{
      "instanceUrl":""
    }
  }
}`))

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "powerplatform_environment_variable" "missing" {
  environment_id = "00000000-0000-0000-0000-000000000001"
  schema_name    = "contoso_ApiBaseUrl"
  value          = "x"
}
`,
				ExpectError: regexp.MustCompile(`does not have Dataverse\s+linked`),
			},
		},
	})
}

func TestUnitEnvironmentVariableResource_Validate_Create_Fails_When_Multiple_Values_Exist(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerEnvironmentResponder()
	registerDefinitionResponder()
	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/environmentvariablevalues?%24filter=_environmentvariabledefinitionid_value+eq+11111111-1111-1111-1111-111111111111&%24select=environmentvariablevalueid%2Cschemaname%2Cvalue",
		httpmock.NewStringResponder(http.StatusOK, `{"value":[
  {"environmentvariablevalueid":"22222222-2222-2222-2222-222222222222","schemaname":"contoso_ApiBaseUrl","value":"a"},
  {"environmentvariablevalueid":"33333333-3333-3333-3333-333333333333","schemaname":"contoso_ApiBaseUrl","value":"b"}
]}`))

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "powerplatform_environment_variable" "broken" {
  environment_id = "00000000-0000-0000-0000-000000000001"
  schema_name    = "contoso_ApiBaseUrl"
  value          = "x"
}
`,
				ExpectError: regexp.MustCompile(`multiple current values found`),
			},
		},
	})
}

func registerEnvironmentVariableHappyPathResponders() {
	registerEnvironmentResponder()
	registerDefinitionResponder()

	getValuesCallCount := 0
	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/environmentvariablevalues?%24filter=_environmentvariabledefinitionid_value+eq+11111111-1111-1111-1111-111111111111&%24select=environmentvariablevalueid%2Cschemaname%2Cvalue",
		func(req *http.Request) (*http.Response, error) {
			getValuesCallCount++
			if getValuesCallCount == 1 {
				return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"environmentvariablevalueid":"22222222-2222-2222-2222-222222222222","schemaname":"contoso_ApiBaseUrl","value":"https://api.contoso.example"}]}`), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/environmentvariablevalues",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusNoContent, "")
			resp.Header.Set("Odata-Entityid", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/environmentvariablevalues(22222222-2222-2222-2222-222222222222)")
			return resp, nil
		})

	httpmock.RegisterResponder("DELETE", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/environmentvariablevalues%2822222222-2222-2222-2222-222222222222%29",
		httpmock.NewStringResponder(http.StatusNoContent, ""))

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/environmentvariablevalues%2822222222-2222-2222-2222-222222222222%29?%24select=environmentvariablevalueid%2Cschemaname%2Cvalue",
		httpmock.NewStringResponder(http.StatusOK, `{"environmentvariablevalueid":"22222222-2222-2222-2222-222222222222","schemaname":"contoso_ApiBaseUrl","value":"https://api.contoso.example"}`))
}

func registerEnvironmentVariableUpdateResponders() {
	registerEnvironmentResponder()
	registerDefinitionResponder()

	getValuesCallCount := 0
	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/environmentvariablevalues?%24filter=_environmentvariabledefinitionid_value+eq+11111111-1111-1111-1111-111111111111&%24select=environmentvariablevalueid%2Cschemaname%2Cvalue",
		func(req *http.Request) (*http.Response, error) {
			getValuesCallCount++
			if getValuesCallCount == 1 {
				return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[{"environmentvariablevalueid":"22222222-2222-2222-2222-222222222222","schemaname":"contoso_ApiBaseUrl","value":"https://api.contoso.example"}]}`), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/environmentvariablevalues",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusNoContent, "")
			resp.Header.Set("Odata-Entityid", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/environmentvariablevalues(22222222-2222-2222-2222-222222222222)")
			return resp, nil
		})

	httpmock.RegisterResponder("PATCH", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/environmentvariablevalues%2822222222-2222-2222-2222-222222222222%29",
		httpmock.NewStringResponder(http.StatusNoContent, ""))

	httpmock.RegisterResponder("DELETE", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/environmentvariablevalues%2822222222-2222-2222-2222-222222222222%29",
		httpmock.NewStringResponder(http.StatusNoContent, ""))

	readValueCallCount := 0
	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/environmentvariablevalues%2822222222-2222-2222-2222-222222222222%29?%24select=environmentvariablevalueid%2Cschemaname%2Cvalue",
		func(req *http.Request) (*http.Response, error) {
			readValueCallCount++
			if readValueCallCount == 1 {
				return httpmock.NewStringResponse(http.StatusOK, `{"environmentvariablevalueid":"22222222-2222-2222-2222-222222222222","schemaname":"contoso_ApiBaseUrl","value":"https://api.contoso.example"}`), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, `{"environmentvariablevalueid":"22222222-2222-2222-2222-222222222222","schemaname":"contoso_ApiBaseUrl","value":"https://api2.contoso.example"}`), nil
		})
}

func registerEnvironmentResponder() {
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		httpmock.NewStringResponder(http.StatusOK, `{
  "name":"00000000-0000-0000-0000-000000000001",
  "id":"/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001",
  "type":"Microsoft.BusinessAppPlatform/scopes/admin/environments",
  "location":"unitedstates",
  "properties":{
    "displayName":"Test",
    "azureRegion":"unitedstates",
    "createdTime":"2024-01-01T00:00:00Z",
    "environmentSku":"Sandbox",
    "linkedEnvironmentMetadata":{
      "instanceUrl":"https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/"
    }
  }
}`))
}

func registerDefinitionResponder() {
	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/environmentvariabledefinitions?%24filter=schemaname+eq+%27contoso_ApiBaseUrl%27&%24select=environmentvariabledefinitionid%2Cschemaname%2Cdisplayname%2Cdescription%2Cdefaultvalue%2Ctype%2Cvalueschema%2Csecretstore",
		httpmock.NewStringResponder(http.StatusOK, `{"value":[{"environmentvariabledefinitionid":"11111111-1111-1111-1111-111111111111","schemaname":"contoso_ApiBaseUrl","displayname":"API Base URL","description":"Base URL for downstream API","defaultvalue":"https://default.contoso.example","type":100000000,"valueschema":"{\"type\":\"string\"}","secretstore":1}]}`))
}
