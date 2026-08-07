// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environmentvariable_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitEnvironmentVariableResource_Validate_Read_ParentDeleted(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerEnvironmentVariableHappyPathResponders()

	envGetCallCount := 0
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			envGetCallCount++
			// Step 1 (create + post-apply plan) makes 10 environment lookups; from the step 2
			// refresh onward the parent environment has been deleted out of band.
			if envGetCallCount > 10 {
				return httpmock.NewStringResponse(http.StatusNotFound, `{"error":{"code":"EnvironmentNotFound","message":"The environment could not be found."}}`), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, `{
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
}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: create the environment variable value successfully.
				Config: `
resource "powerplatform_environment_variable" "text" {
  environment_id = "00000000-0000-0000-0000-000000000001"
  schema_name    = "contoso_ApiBaseUrl"
  value          = "https://api.contoso.example"
}
`,
				Check: resource.TestCheckResourceAttr("powerplatform_environment_variable.text", "schema_name", "contoso_ApiBaseUrl"),
			},
			{
				// Step 2: refresh after the parent environment was deleted out of band.
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
