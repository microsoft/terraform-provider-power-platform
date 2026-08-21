// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access_test //nolint:revive // the underscored package name predates this file and matches every service in the repo

import (
	"net/http"
	"regexp"
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

// scope_type "tenant" reads the tenant scoped assignments, with no identifier to accompany it.
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
				Config: `data "powerplatform_rbac_role_assignments" "test" { scope_type = "tenant" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.#", "2"),
					resource.TestCheckResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.0.id", tenantAssignmentId),
					resource.TestCheckResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.0.principal_type", "ApplicationUser"),
					resource.TestCheckResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.0.role_definition_id", testRoleDefinitionId),
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
				data "powerplatform_rbac_role_assignments" "test" {
					scope_type     = "environment"
					environment_id = "` + testEnvironmentId + `"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.0.id", environmentAssignmentId),
					resource.TestCheckResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.0.scope", "/environments/"+testEnvironmentId),
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
				data "powerplatform_rbac_role_assignments" "test" {
					scope_type           = "environment_group"
					environment_group_id = "` + testDataSourceEnvironmentGroupId + `"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.0.id", "44444444-4444-4444-4444-444444444444"),
					resource.TestCheckResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.0.scope", "/environmentgroups/"+testDataSourceEnvironmentGroupId),
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
				Config: `data "powerplatform_rbac_role_assignments" "test" { scope_type = "tenant" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.#", "0"),
				),
			},
		},
	})
}

// A failing list surfaces the scoped diagnostic instead of an empty result, at every scope.
func TestUnitRoleAssignmentsDataSource_Validate_Read_Error(t *testing.T) {
	cases := []struct {
		name       string
		collection string
		config     string
	}{
		{"tenant", tenantCollection, `
			scope_type = "tenant"`},
		{"environment", environmentCollection, `
			scope_type     = "environment"
			environment_id = "` + testEnvironmentId + `"`},
		{"environment_group", envGroupCollection, `
			scope_type           = "environment_group"
			environment_group_id = "` + testEnvironmentGroupId + `"`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			mocks.ActivateEnvironmentHttpMocks()

			httpmock.RegisterResponder("GET", tc.collection+apiVersionQuery,
				func(_ *http.Request) (*http.Response, error) {
					return httpmock.NewStringResponse(http.StatusForbidden, `{"error":"forbidden"}`), nil
				})

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: `
						data "powerplatform_rbac_role_assignments" "test" {` + tc.config + `
						}`,
						ExpectError: regexp.MustCompile(`(?s)Failed to list role assignments at.*scope`),
					},
				},
			})
		})
	}
}
