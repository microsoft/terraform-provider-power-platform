// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package powerplatform

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccTestRest_Validate_Create(t *testing.T) {
	beforeUpdateRegex := `^\{"@odata\.context":"https:\/\/org[0-9a-fA-F]{8}\.crm4\.dynamics\.com\/api\/data\/v9\.2\/\$metadata#accounts\(name,accountid\)\/\$entity","@odata\.etag":"W\/\\"[0-9]{7}\\"","name":"powerplatform_rest","accountid":"00000000-0000-0000-0000-000000000001"\}$`
	afterUpdateRegex := `^\{"@odata\.context":"https:\/\/org[0-9a-fA-F]{8}\.crm4\.dynamics\.com\/api\/data\/v9\.2\/\$metadata#accounts\(name,accountid\)\/\$entity","@odata\.etag":"W\/\\"[0-9]{7}\\"","name":"powerplatform_rest_change","accountid":"00000000-0000-0000-0000-000000000001"\}$`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `

				resource "powerplatform_environment" "env" {
					display_name     = "` + mock_helpers.TestName() + `"
					location         = "europe"
					environment_type = "Sandbox"
					dataverse = {
						language_code     = "1033"
						currency_code     = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}


				locals {
					body = jsonencode({
						"accountid" : "00000000-0000-0000-0000-000000000001",
						"name" : "powerplatform_rest",
						"creditonhold" : true,
						"address1_latitude" : 47.6396,
						"description" : "This is the updated description of the sample account",
						"revenue" : 6000000,
						"accountcategorycode" : 2
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
						},
						{
							name  = "Prefer"
							value = "return=representation"
						}
					]
				}

				resource "powerplatform_rest" "query" {
					environment_id = powerplatform_environment.env.id
					create = {
						url     = "/api/data/v9.2/accounts?$select=name,accountid"
						method  = "POST"
						body    = local.body
						headers = local.headers
					}
					read = {
						url    = "api/data/v9.2/accounts(00000000-0000-0000-0000-000000000001)?$select=name,accountid"
						method = "GET"
					}
					update = {
						url     = "api/data/v9.2/accounts(00000000-0000-0000-0000-000000000001)?$select=name,accountid"
						method  = "PATCH"
						body    = local.body
						headers = local.headers
					}
					destroy = {
						url    = "api/data/v9.2/accounts(00000000-0000-0000-0000-000000000001)"
						method = "DELETE"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_rest.query", "output.body", regexp.MustCompile(beforeUpdateRegex)),
				),
			},
			{
				Config: TestsProviderConfig + `

				locals {
					body = jsonencode({
						"accountid" : "00000000-0000-0000-0000-000000000001",
						"name" : "powerplatform_rest_change",
						"creditonhold" : true,
						"address1_latitude" : 47.6396,
						"description" : "This is the updated description of the sample account",
						"revenue" : 6000000,
						"accountcategorycode" : 2
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
						},
						{
							name  = "Prefer"
							value = "return=representation"
						}
					]
				}

				resource "powerplatform_rest" "query" {
					environment_id = powerplatform_environment.env.id
					create = {
						url     = "/api/data/v9.2/accounts?$select=name,accountid"
						method  = "POST"
						body    = local.body
						headers = local.headers
					}
					read = {
						url    = "api/data/v9.2/accounts(00000000-0000-0000-0000-000000000001)?$select=name,accountid"
						method = "GET"
					}
					update = {
						url     = "api/data/v9.2/accounts(00000000-0000-0000-0000-000000000001)?$select=name,accountid"
						method  = "PATCH"
						body    = local.body
						headers = local.headers
					}
					destroy = {
						url    = "api/data/v9.2/accounts(00000000-0000-0000-0000-000000000001)"
						method = "DELETE"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_rest.query", "output.body", regexp.MustCompile(afterUpdateRegex)),
				),
			},
		},
	})
}

func TestUnitTestRest_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/rest/tests/resource/Web_Api_Validate_Create/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("POST", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts?$select=name,accountid`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("services/rest/tests/resource/Web_Api_Validate_Create/post_account.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `

				locals {
					body = jsonencode({
						"accountid" : "00000000-0000-0000-0000-000000000001",
						"name" : "powerplatform_rest",
						"creditonhold" : true,
						"address1_latitude" : 47.6396,
						"description" : "This is the updated description of the sample account",
						"revenue" : 6000000,
						"accountcategorycode" : 2
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
						},
						{
							name  = "Prefer"
							value = "return=representation"
						}
					]
				}

				resource "powerplatform_rest" "query" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					create = {
						url     = "api/data/v9.2/accounts?$select=name,accountid"
						method  = "POST"
						body    = local.body
						headers = local.headers
					}
					read = {
						url    = "api/data/v9.2/accounts(00000000-0000-0000-0000-000000000001)?$select=name,accountid"
						method = "GET"
					}
					update = {
						url     = "api/data/v9.2/accounts(00000000-0000-0000-0000-000000000001)?$select=name,accountid"
						method  = "PATCH"
						body    = local.body
						headers = local.headers
					}
					destroy = {
						url    = "api/data/v9.2/accounts(00000000-0000-0000-0000-000000000001)"
						method = "DELETE"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_rest.query", "output.body", httpmock.File("services/rest/tests/resource/Web_Api_Validate_Create/post_account.json").String()),
				),
			},
		},
	})
}
