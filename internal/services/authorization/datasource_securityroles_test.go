// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package authorization_test

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccSecurityDataSource_Validate_Read(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "env" {
					display_name      = "` + mocks.TestName() + `"
					location          = "unitedstates"
					environment_type  = "Sandbox"
					dataverse = {
						language_code     = "1033"
						currency_code     = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				data "powerplatform_security_roles" "all" {
					environment_id = powerplatform_environment.env.id
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "environment_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "security_roles.#", regexp.MustCompile(`^[1-9]\d*$`)),

					resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "security_roles.0.role_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "security_roles.0.name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "security_roles.0.is_managed", regexp.MustCompile(helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "security_roles.0.business_unit_id", regexp.MustCompile(helpers.GuidRegex)),
				),
			},
		},
	})
}

func TestUnitSecurityDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/security_roles/Validate_Read/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/roles",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/security_roles/Validate_Read/get_security_roles.json").String()), nil
		})

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_security_roles" "all" {
					environment_id = "00000000-0000-0000-0000-000000000001"
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "id", regexp.MustCompile(helpers.GuidRegex)),

					resource.TestCheckResourceAttr("data.powerplatform_security_roles.all", "security_roles.#", "72"),
					resource.TestCheckResourceAttr("data.powerplatform_security_roles.all", "environment_id", "00000000-0000-0000-0000-000000000001"),

					resource.TestCheckResourceAttr("data.powerplatform_security_roles.all", "security_roles.0.role_id", "4931681d-8163-e811-a965-000d3a11fe32"),
					resource.TestCheckResourceAttr("data.powerplatform_security_roles.all", "security_roles.0.name", "Export Customizations (Solution Checker)"),
					resource.TestCheckResourceAttr("data.powerplatform_security_roles.all", "security_roles.0.is_managed", "true"),
					resource.TestCheckResourceAttr("data.powerplatform_security_roles.all", "security_roles.0.business_unit_id", "1360fdcb-b6e1-ee11-904c-002248dad9c1"),
				),
			},
		},
	})
}

func TestUnitSecurityDataSource_Validate_No_Dataverse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/security_roles/Validate_No_Dataverse/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/roles",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/security_roles/Validate_No_Dataverse/get_security_roles.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/security_roles/Validate_No_Dataverse/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/datasource/security_roles/Validate_No_Dataverse/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/security_roles/Validate_No_Dataverse/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	resource.ParallelTest(t, resource.TestCase{
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

				data "powerplatform_security_roles" "all" {
					environment_id = powerplatform_environment.env.id
				}`,
				ExpectError: regexp.MustCompile(`No Dataverse exists in environment`),
				Check:       resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestUnitSecurityDataSource_Validate_Read_Filter_BusinessUnit(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/security_roles/Validate_Read_Filter_BusinessUnit/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/roles?%24filter=_businessunitid_value+eq+00000000-0000-0000-0000-000000000002",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/security_roles/Validate_Read_Filter_BusinessUnit/get_security_roles.json").String()), nil
		})

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_security_roles" "all" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					business_unit_id = "00000000-0000-0000-0000-000000000002"
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "id", regexp.MustCompile(helpers.GuidRegex)),

					resource.TestCheckResourceAttr("data.powerplatform_security_roles.all", "security_roles.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_security_roles.all", "environment_id", "00000000-0000-0000-0000-000000000001"),

					resource.TestCheckResourceAttr("data.powerplatform_security_roles.all", "security_roles.0.role_id", "4931681d-8163-e811-a965-000d3a11fe32"),
					resource.TestCheckResourceAttr("data.powerplatform_security_roles.all", "security_roles.0.name", "Export Customizations (Solution Checker)"),
					resource.TestCheckResourceAttr("data.powerplatform_security_roles.all", "security_roles.0.is_managed", "true"),
					resource.TestCheckResourceAttr("data.powerplatform_security_roles.all", "security_roles.0.business_unit_id", "00000000-0000-0000-0000-000000000002"),
				),
			},
		},
	})
}
