// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitRoleDefinitionsDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/authorization/roleDefinitions?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read_RoleDefinitions/get_role_definitions.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_role_definitions" "test" {
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_role_definitions.test", "role_definitions.#", "2"),
					resource.TestCheckResourceAttr("data.powerplatform_role_definitions.test", "role_definitions.0.role_definition_id", "ff954d61-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_role_definitions.test", "role_definitions.0.role_definition_name", "Power Platform Administrator"),
					resource.TestCheckResourceAttr("data.powerplatform_role_definitions.test", "role_definitions.0.permissions.#", "2"),
					resource.TestCheckResourceAttr("data.powerplatform_role_definitions.test", "role_definitions.1.role_definition_id", "ff954d61-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("data.powerplatform_role_definitions.test", "role_definitions.1.role_definition_name", "Environment Administrator"),
				),
			},
		},
	})
}
