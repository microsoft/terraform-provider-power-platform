// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package rest_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

const whoAmIResponseRegex = `^{"@odata.context":"https:\/\/[^"]+","BusinessUnitId":"[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}","UserId":"[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}","OrganizationId":"[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"}$`

func TestUnitDatasourceRestQuery_WhoAmI_Using_Scope(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Web_Apis_WhoAmI/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/whoami`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Web_Apis_WhoAmI/get_whoami.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_rest_query" "webapi_query" {
					scope 		   = "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/.default"
					url            = "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/whoami"
					method         = "GET"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_rest_query.webapi_query", "output.body", httpmock.File("tests/datasource/Web_Apis_WhoAmI/get_whoami.json").String()),
				),
			},
		},
	})
}

func TestAccDatasourceRestQuery_WhoAmI_Using_Scope(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "env" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
						language_code     = "1033"
						currency_code     = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}
				

				data "powerplatform_rest_query" "webapi_query" {
					scope                = "${powerplatform_environment.env.dataverse.url}/.default"
					url                  = "${powerplatform_environment.env.dataverse.url}api/data/v9.2/WhoAmI"
					method               = "GET"
					expected_http_status = [200]
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_rest_query.webapi_query", "output.body", regexp.MustCompile(whoAmIResponseRegex)),
				),
			},
		},
	})
}
