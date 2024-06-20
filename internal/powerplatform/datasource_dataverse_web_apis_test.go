// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package powerplatform

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestUnitDatasourceDataverse_Web_Apis_WhoAmI(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/dataverse_web_api/tests/datasource/Web_Apis_WhoAmI/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/whoami`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/dataverse_web_api/tests/datasource/Web_Apis_WhoAmI/get_whoami.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				data "powerplatform_dataverse_web_apis" "webapi_query" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					url            = "api/data/v9.2/whoami"
					method         = "GET"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}
