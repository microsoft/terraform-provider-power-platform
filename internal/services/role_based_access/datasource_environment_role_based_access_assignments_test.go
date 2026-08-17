// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitEnvironmentRoleBasedAccessAssignmentsDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/authorization/environments/eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee/roleAssignments?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read_EnvironmentRoleAssignments/get_role_assignments.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_environment_role_based_access_assignments" "test" {
					environment_id = "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_environment_role_based_access_assignments.test", "role_assignments.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_role_based_access_assignments.test", "role_assignments.0.id", "33333333-3333-3333-3333-333333333333"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_role_based_access_assignments.test", "role_assignments.0.scope", "/environments/eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_role_based_access_assignments.test", "role_assignments.0.enterprise_application_object_id", "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_role_based_access_assignments.test", "role_assignments.0.role_definition_id", "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
				),
			},
		},
	})
}

func TestUnitEnvironmentRoleBasedAccessAssignmentsDataSource_Validate_Read_Error(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/authorization/environments/eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee/roleAssignments?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusForbidden, `{"error":{"code":"Forbidden","message":"Access denied"}}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_environment_role_based_access_assignments" "test" {
					environment_id = "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
				}`,
				ExpectError: regexp.MustCompile(`Failed to list environment role assignments`),
			},
		},
	})
}
