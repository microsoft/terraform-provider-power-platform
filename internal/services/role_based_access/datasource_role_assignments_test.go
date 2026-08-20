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

const testDataSourceEnvironmentGroupId = "ffffffff-ffff-ffff-ffff-ffffffffffff"

func registerListMock(collectionPath, fixture string) {
	httpmock.RegisterResponder("GET", collectionPath+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fixture).String()), nil
		})
}

// With neither identifier set the data source reads the tenant scoped assignments.
func TestUnitRoleAssignmentsDataSource_Validate_Read_Tenant_Scope(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()
	registerListMock(tenantCollection, "tests/datasource/Validate_Read_RoleAssignments/get_role_assignments.json")

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "powerplatform_role_assignments" "test" { scope_type = "tenant" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_role_assignments.test", "role_assignments.#", "2"),
					resource.TestCheckResourceAttr("data.powerplatform_role_assignments.test", "role_assignments.0.id", tenantAssignmentId),
					resource.TestCheckResourceAttr("data.powerplatform_role_assignments.test", "role_assignments.0.principal_type", "ApplicationUser"),
					resource.TestCheckResourceAttr("data.powerplatform_role_assignments.test", "role_assignments.0.role_definition_id", testRoleDefinitionId),
				),
			},
		},
	})
}

func TestUnitRoleAssignmentsDataSource_Validate_Read_Environment_Scope(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()
	registerListMock(environmentCollection, "tests/datasource/Validate_Read_EnvironmentRoleAssignments/get_role_assignments.json")

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_role_assignments" "test" {
					scope_type     = "environment"
					environment_id = "` + testEnvironmentId + `"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_role_assignments.test", "role_assignments.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_role_assignments.test", "role_assignments.0.id", environmentAssignmentId),
					resource.TestCheckResourceAttr("data.powerplatform_role_assignments.test", "role_assignments.0.scope", "/environments/"+testEnvironmentId),
				),
			},
		},
	})
}

func TestUnitRoleAssignmentsDataSource_Validate_Read_EnvironmentGroup_Scope(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()
	registerListMock("https://api.powerplatform.com/authorization/environmentGroups/"+testDataSourceEnvironmentGroupId+"/roleAssignments",
		"tests/datasource/Validate_Read_EnvironmentGroupRoleAssignments/get_role_assignments.json")

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_role_assignments" "test" {
					scope_type           = "environment_group"
					environment_group_id = "` + testDataSourceEnvironmentGroupId + `"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_role_assignments.test", "role_assignments.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_role_assignments.test", "role_assignments.0.id", "44444444-4444-4444-4444-444444444444"),
					resource.TestCheckResourceAttr("data.powerplatform_role_assignments.test", "role_assignments.0.scope", "/environmentGroups/"+testDataSourceEnvironmentGroupId),
				),
			},
		},
	})
}

func TestUnitRoleAssignmentsDataSource_Validate_Read_Empty(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()
	registerListMock(tenantCollection, "tests/datasource/Validate_Read_RoleAssignments_Empty/get_role_assignments.json")

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "powerplatform_role_assignments" "test" { scope_type = "tenant" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_role_assignments.test", "role_assignments.#", "0"),
				),
			},
		},
	})
}
