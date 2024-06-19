// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package powerplatform

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestUnitTestDataverse_Web_Api_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/dataverse_web_api/tests/resource/Web_Api_Validate_Create/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("POST", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts?$select=name,accountid`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("services/dataverse_web_api/tests/resource/Web_Api_Validate_Create/post_account.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_dataverse_web_api" "query" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					create = {
						url    = "api/data/v9.2/accounts?$select=name,accountid"
						method = "POST"
						body   = jsonencode({
							"accountid": "00000000-0000-0000-0000-000000000031",
							"name": "powerplatform_dataverse_web_api",
							"creditonhold": false,
							"address1_latitude": 47.639583,
							"description": "This is the description of the sample account",
							"revenue": 5000000,
							"accountcategorycode": 1
						})
						headers = [
							{
								name  = "Content-Type"
								value = "application/json; charset=utf-8"
							},
							{
								name  = "OData-MaxVersion"
								value = "4.0"
							},
							{
								name  = "OData-Version"
								value = "4.0"
							}
						]
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}
