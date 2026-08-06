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

func TestUnitRoleBasedAccessAssignmentsDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/authorization/roleAssignments?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read_RoleAssignments/get_role_assignments.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_role_based_access_assignments" "test" {
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_role_based_access_assignments.test", "role_assignments.#", "2"),
					resource.TestCheckResourceAttr("data.powerplatform_role_based_access_assignments.test", "role_assignments.0.id", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckResourceAttr("data.powerplatform_role_based_access_assignments.test", "role_assignments.0.scope", "/tenants/dddddddd-dddd-dddd-dddd-dddddddddddd"),
					resource.TestCheckResourceAttr("data.powerplatform_role_based_access_assignments.test", "role_assignments.0.principal_type", "ApplicationUser"),
					resource.TestCheckResourceAttr("data.powerplatform_role_based_access_assignments.test", "role_assignments.0.enterprise_application_object_id", "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
					resource.TestCheckResourceAttr("data.powerplatform_role_based_access_assignments.test", "role_assignments.0.role_definition_id", "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
					resource.TestCheckResourceAttr("data.powerplatform_role_based_access_assignments.test", "role_assignments.0.created_by_principal_type", "User"),
					resource.TestCheckResourceAttr("data.powerplatform_role_based_access_assignments.test", "role_assignments.0.created_by_principal_object_id", "cccccccc-cccc-cccc-cccc-cccccccccccc"),
					resource.TestCheckResourceAttr("data.powerplatform_role_based_access_assignments.test", "role_assignments.0.created_on", "2026-06-22T17:00:00Z"),
					resource.TestCheckNoResourceAttr("data.powerplatform_role_based_access_assignments.test", "role_assignments.0.expires_on"),
					resource.TestCheckResourceAttr("data.powerplatform_role_based_access_assignments.test", "role_assignments.1.id", "22222222-2222-2222-2222-222222222222"),
					resource.TestCheckResourceAttr("data.powerplatform_role_based_access_assignments.test", "role_assignments.1.expires_on", "2027-06-23T17:00:00Z"),
				),
			},
		},
	})
}

func TestUnitRoleBasedAccessAssignmentsDataSource_Validate_Read_Empty(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/authorization/roleAssignments?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read_RoleAssignments_Empty/get_role_assignments.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_role_based_access_assignments" "test" {
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_role_based_access_assignments.test", "role_assignments.#", "0"),
				),
			},
		},
	})
}

func TestUnitRoleBasedAccessAssignmentsDataSource_Validate_Read_Error(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/authorization/roleAssignments?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusForbidden, `{"error":{"code":"Forbidden","message":"Access denied"}}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_role_based_access_assignments" "test" {
				}`,
				ExpectError: regexp.MustCompile(`Failed to list role assignments`),
			},
		},
	})
}
