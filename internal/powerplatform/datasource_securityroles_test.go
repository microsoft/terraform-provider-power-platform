// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccSecurityDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "env" {
					display_name      = "TestAccSecurityDataSource_Validate_Read"
					location          = "europe"
					language_code     = "1033"
					currency_code     = "USD"
					environment_type  = "Sandbox"
					security_group_id = "00000000-0000-0000-0000-000000000000"
				}

				data "powerplatform_security_roles" "all" {
					environment_id = powerplatform_environment.env.id
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "environment_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "security_roles.#", regexp.MustCompile(`^[1-9]\d*$`)),

					resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "security_roles.0.role_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "security_roles.0.name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "security_roles.0.is_managed", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "security_roles.0.business_unit_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
				),
			},
		},
	})
}

func TestUnitSecurityDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/authorization/tests/datasource/security_roles/Validate_Read/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/roles",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/authorization/tests/datasource/security_roles/Validate_Read/get_security_roles.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				data "powerplatform_security_roles" "all" {
					environment_id = "00000000-0000-0000-0000-000000000001"
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

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
