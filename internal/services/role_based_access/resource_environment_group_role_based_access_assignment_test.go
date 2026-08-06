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

func TestUnitEnvironmentGroupRoleBasedAccessAssignmentResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("POST", `https://api.powerplatform.com/authorization/environmentGroups/dddddddd-dddd-dddd-dddd-dddddddddddd/roleAssignments?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/Validate_Create_EnvGroup/post_role_assignment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/authorization/environmentGroups/dddddddd-dddd-dddd-dddd-dddddddddddd/roleAssignments?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_EnvGroup/get_role_assignments.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", `https://api.powerplatform.com/authorization/environmentGroups/dddddddd-dddd-dddd-dddd-dddddddddddd/roleAssignments/22222222-2222-2222-2222-222222222222?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_group_role_based_access_assignment" "test" {
					environment_group_id = "dddddddd-dddd-dddd-dddd-dddddddddddd"
					principal_object_id  = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
					principal_type       = "ApplicationUser"
					role_definition_id   = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group_role_based_access_assignment.test", "id", "22222222-2222-2222-2222-222222222222"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_role_based_access_assignment.test", "environment_group_id", "dddddddd-dddd-dddd-dddd-dddddddddddd"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_role_based_access_assignment.test", "principal_object_id", "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_role_based_access_assignment.test", "principal_type", "ApplicationUser"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_role_based_access_assignment.test", "role_definition_id", "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_role_based_access_assignment.test", "scope", "/environmentGroups/dddddddd-dddd-dddd-dddd-dddddddddddd"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_role_based_access_assignment.test", "created_on", "2026-06-22T16:00:00Z"),
				),
			},
		},
	})
}
