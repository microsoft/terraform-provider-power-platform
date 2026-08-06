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

func TestUnitRoleBasedAccessAssignmentResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("POST", `https://api.powerplatform.com/authorization/roleAssignments?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/Validate_Create_Tenant/post_role_assignment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/authorization/roleAssignments?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Tenant/get_role_assignments.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", `https://api.powerplatform.com/authorization/roleAssignments/11111111-1111-1111-1111-111111111111?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_based_access_assignment" "test" {
					principal_object_id = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
					principal_type      = "ApplicationUser"
					role_definition_id  = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_role_based_access_assignment.test", "id", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckResourceAttr("powerplatform_role_based_access_assignment.test", "principal_object_id", "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
					resource.TestCheckResourceAttr("powerplatform_role_based_access_assignment.test", "principal_type", "ApplicationUser"),
					resource.TestCheckResourceAttr("powerplatform_role_based_access_assignment.test", "role_definition_id", "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
					resource.TestCheckResourceAttr("powerplatform_role_based_access_assignment.test", "scope", "/tenants/99999999-9999-9999-9999-999999999999"),
					resource.TestCheckResourceAttr("powerplatform_role_based_access_assignment.test", "created_on", "2026-06-22T15:09:35Z"),
				),
			},
		},
	})
}
