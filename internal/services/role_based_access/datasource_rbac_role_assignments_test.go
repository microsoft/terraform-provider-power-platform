// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access_test //nolint:revive // the underscored package name predates this file and matches every service in the repo

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
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

// A tenant scoped read lists the assignment just made at the tenant. The tenant also holds grants
// this test did not create, so only the shape of the list is asserted here; the scoped tests below
// read a freshly created scope and can pin the exact assignment.
func TestAccRoleAssignmentsDataSource_Validate_Read_Tenant_Scope(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azuread": {
				VersionConstraint: constants.AZURE_AD_PROVIDER_VERSION_CONSTRAINT,
				Source:            "hashicorp/azuread",
			},
			"time": {
				Source: "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
				resource "azuread_application_registration" "test_app" {
					display_name = "` + mocks.TestName() + `"
				}

				resource "azuread_service_principal" "test_sp" {
					client_id = azuread_application_registration.test_app.client_id
				}

				resource "time_sleep" "wait_for_service_principal" {
					create_duration = "60s"

					depends_on = [azuread_service_principal.test_sp]
				}

				resource "powerplatform_rbac_role_assignment" "test" {
					scope_type         = "tenant"
					principal_id       = azuread_service_principal.test_sp.object_id
					principal_type     = "ApplicationUser"
					role_definition_id = "` + roleBasedAccessAdministratorRoleId + `"

					depends_on = [time_sleep.wait_for_service_principal]
				}

				data "powerplatform_rbac_role_assignments" "test" {
					scope_type = "tenant"

					depends_on = [powerplatform_rbac_role_assignment.test]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_rbac_role_assignments.test", "scope_type", "tenant"),
					resource.TestCheckNoResourceAttr("data.powerplatform_rbac_role_assignments.test", "environment_id"),
					resource.TestCheckNoResourceAttr("data.powerplatform_rbac_role_assignments.test", "environment_group_id"),
					resource.TestMatchResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.0.id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.0.scope", regexp.MustCompile(`^/tenants/`)),
				),
			},
		},
	})
}

// The environment is created by this test, so the read returns exactly the one assignment made
// against it and every element can be pinned by index.
func TestAccRoleAssignmentsDataSource_Validate_Read_Environment_Scope(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azuread": {
				VersionConstraint: constants.AZURE_AD_PROVIDER_VERSION_CONSTRAINT,
				Source:            "hashicorp/azuread",
			},
			"time": {
				Source: "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
				resource "azuread_application_registration" "test_app" {
					display_name = "` + mocks.TestName() + `"
				}

				resource "azuread_service_principal" "test_sp" {
					client_id = azuread_application_registration.test_app.client_id
				}

				resource "time_sleep" "wait_for_service_principal" {
					create_duration = "60s"

					depends_on = [azuread_service_principal.test_sp]
				}

				resource "powerplatform_environment" "test_environment" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
				}

				resource "powerplatform_rbac_role_assignment" "test" {
					scope_type         = "environment"
					environment_id     = powerplatform_environment.test_environment.id
					principal_id       = azuread_service_principal.test_sp.object_id
					principal_type     = "ApplicationUser"
					role_definition_id = "` + roleBasedAccessAdministratorRoleId + `"

					depends_on = [time_sleep.wait_for_service_principal]
				}

				data "powerplatform_rbac_role_assignments" "test" {
					scope_type     = "environment"
					environment_id = powerplatform_environment.test_environment.id

					depends_on = [powerplatform_rbac_role_assignment.test]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.powerplatform_rbac_role_assignments.test", "environment_id", "powerplatform_environment.test_environment", "id"),
					resource.TestCheckResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.#", "1"),
					resource.TestCheckResourceAttrPair("data.powerplatform_rbac_role_assignments.test", "role_assignments.0.id", "powerplatform_rbac_role_assignment.test", "id"),
					resource.TestCheckResourceAttrPair("data.powerplatform_rbac_role_assignments.test", "role_assignments.0.principal_id", "azuread_service_principal.test_sp", "object_id"),
					resource.TestCheckResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.0.principal_type", "ApplicationUser"),
					resource.TestCheckResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.0.role_definition_id", roleBasedAccessAdministratorRoleId),
					resource.TestMatchResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.0.scope", regexp.MustCompile(`/environments/`)),
					resource.TestCheckResourceAttrSet("data.powerplatform_rbac_role_assignments.test", "role_assignments.0.created_on"),
				),
			},
		},
	})
}

// The environment group is created by this test, so the read returns exactly the one assignment made
// against it and every element can be pinned by index.
func TestAccRoleAssignmentsDataSource_Validate_Read_EnvironmentGroup_Scope(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azuread": {
				VersionConstraint: constants.AZURE_AD_PROVIDER_VERSION_CONSTRAINT,
				Source:            "hashicorp/azuread",
			},
			"time": {
				Source: "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
				resource "azuread_application_registration" "test_app" {
					display_name = "` + mocks.TestName() + `"
				}

				resource "azuread_service_principal" "test_sp" {
					client_id = azuread_application_registration.test_app.client_id
				}

				resource "time_sleep" "wait_for_service_principal" {
					create_duration = "60s"

					depends_on = [azuread_service_principal.test_sp]
				}

				resource "powerplatform_environment_group" "test_env_group" {
					display_name = "` + mocks.TestName() + `"
					description  = "Environment group for the role assignments data source acceptance test"
				}

				resource "powerplatform_rbac_role_assignment" "test" {
					scope_type           = "environment_group"
					environment_group_id = powerplatform_environment_group.test_env_group.id
					principal_id         = azuread_service_principal.test_sp.object_id
					principal_type       = "ApplicationUser"
					role_definition_id   = "` + roleBasedAccessAdministratorRoleId + `"

					depends_on = [time_sleep.wait_for_service_principal]
				}

				data "powerplatform_rbac_role_assignments" "test" {
					scope_type           = "environment_group"
					environment_group_id = powerplatform_environment_group.test_env_group.id

					depends_on = [powerplatform_rbac_role_assignment.test]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.powerplatform_rbac_role_assignments.test", "environment_group_id", "powerplatform_environment_group.test_env_group", "id"),
					resource.TestCheckResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.#", "1"),
					resource.TestCheckResourceAttrPair("data.powerplatform_rbac_role_assignments.test", "role_assignments.0.id", "powerplatform_rbac_role_assignment.test", "id"),
					resource.TestCheckResourceAttrPair("data.powerplatform_rbac_role_assignments.test", "role_assignments.0.principal_id", "azuread_service_principal.test_sp", "object_id"),
					resource.TestCheckResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.0.principal_type", "ApplicationUser"),
					resource.TestCheckResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.0.role_definition_id", roleBasedAccessAdministratorRoleId),
					resource.TestMatchResourceAttr("data.powerplatform_rbac_role_assignments.test", "role_assignments.0.scope", regexp.MustCompile(`(?i)/environmentgroups/`)),
					resource.TestCheckResourceAttrSet("data.powerplatform_rbac_role_assignments.test", "role_assignments.0.created_on"),
				),
			},
		},
	})
}
