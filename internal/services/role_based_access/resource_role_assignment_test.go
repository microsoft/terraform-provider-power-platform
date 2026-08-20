// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access_test //nolint:revive // the underscored package name predates this file and matches every service in the repo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

// Microsoft has recased this role's display name, so match on the stable id instead.
const roleBasedAccessAdministratorRoleId = "95e94555-018c-447b-8691-bdac8e12211e"

const (
	tenantAssignmentId      = "11111111-1111-1111-1111-111111111111"
	envGroupAssignmentId    = "22222222-2222-2222-2222-222222222222"
	environmentAssignmentId = "33333333-3333-3333-3333-333333333333"
	testEnvironmentId       = "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
	testEnvironmentGroupId  = "dddddddd-dddd-dddd-dddd-dddddddddddd"
	testPrincipalId         = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	testRoleDefinitionId    = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
)

// registerScopeMocks wires create, list and delete for one scope of the RBAC API. collectionPath is
// the roleAssignments collection for that scope, and fixtureDir holds its recorded responses.
func registerScopeMocks(collectionPath, fixtureDir, assignmentId string) {
	httpmock.RegisterResponder("POST", collectionPath+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File(fixtureDir+"/post_role_assignment.json").String()), nil
		})
	httpmock.RegisterResponder("GET", collectionPath+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fixtureDir+"/get_role_assignments.json").String()), nil
		})
	httpmock.RegisterResponder("DELETE", collectionPath+"/"+assignmentId+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})
}

const apiVersionQuery = "?api-version=2024-10-01"

var (
	tenantCollection      = "https://api.powerplatform.com/authorization/roleAssignments"
	environmentCollection = "https://api.powerplatform.com/authorization/environments/" + testEnvironmentId + "/roleAssignments"
	envGroupCollection    = "https://api.powerplatform.com/authorization/environmentGroups/" + testEnvironmentGroupId + "/roleAssignments"
)

// A role assignment with neither identifier set applies to the tenant.
func TestUnitRoleAssignmentResource_Validate_Create_Tenant_Scope(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()
	registerScopeMocks(tenantCollection, "tests/resource/Validate_Create_Tenant", tenantAssignmentId)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_assignment" "test" {
					scope_type         = "tenant"
					principal_id       = "` + testPrincipalId + `"
					principal_type                   = "ApplicationUser"
					role_definition_id               = "` + testRoleDefinitionId + `"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "id", tenantAssignmentId),
					resource.TestCheckNoResourceAttr("powerplatform_role_assignment.test", "environment_id"),
					resource.TestCheckNoResourceAttr("powerplatform_role_assignment.test", "environment_group_id"),
					resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "scope", "/tenants/00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "created_on", "2026-06-22T15:09:35Z"),
				),
			},
			{
				ResourceName:      "powerplatform_role_assignment.test",
				ImportState:       true,
				ImportStateId:     tenantAssignmentId,
				ImportStateVerify: true,
			},
		},
	})
}

// Setting environment_id routes the same resource to the environment collection.
func TestUnitRoleAssignmentResource_Validate_Create_Environment_Scope(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()
	registerScopeMocks(environmentCollection, "tests/resource/Validate_Create_Environment", environmentAssignmentId)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_assignment" "test" {
					scope_type                       = "environment"
					environment_id                   = "` + testEnvironmentId + `"
					principal_id = "` + testPrincipalId + `"
					principal_type                   = "ApplicationUser"
					role_definition_id               = "` + testRoleDefinitionId + `"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "id", environmentAssignmentId),
					resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "environment_id", testEnvironmentId),
					resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "scope", "/tenants/00000000-0000-0000-0000-000000000001/environments/"+testEnvironmentId),
				),
			},
			{
				ResourceName:      "powerplatform_role_assignment.test",
				ImportState:       true,
				ImportStateId:     "environments/" + testEnvironmentId + "/" + environmentAssignmentId,
				ImportStateVerify: true,
			},
		},
	})
}

// Setting environment_group_id routes the same resource to the environment group collection.
func TestUnitRoleAssignmentResource_Validate_Create_EnvironmentGroup_Scope(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()
	registerScopeMocks(envGroupCollection, "tests/resource/Validate_Create_EnvGroup", envGroupAssignmentId)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_assignment" "test" {
					scope_type                       = "environment_group"
					environment_group_id             = "` + testEnvironmentGroupId + `"
					principal_id = "` + testPrincipalId + `"
					principal_type                   = "ApplicationUser"
					role_definition_id               = "` + testRoleDefinitionId + `"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "id", envGroupAssignmentId),
					resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "environment_group_id", testEnvironmentGroupId),
					resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "scope", "/tenants/00000000-0000-0000-0000-000000000001/environmentgroups/"+testEnvironmentGroupId),
				),
			},
			{
				ResourceName:      "powerplatform_role_assignment.test",
				ImportState:       true,
				ImportStateId:     "environmentGroups/" + testEnvironmentGroupId + "/" + envGroupAssignmentId,
				ImportStateVerify: true,
			},
		},
	})
}

// The two scope identifiers are mutually exclusive.
func TestUnitRoleAssignmentResource_Validate_Scopes_Are_Mutually_Exclusive(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_assignment" "test" {
					scope_type                       = "environment"
					environment_id                   = "` + testEnvironmentId + `"
					environment_group_id             = "` + testEnvironmentGroupId + `"
					principal_id = "` + testPrincipalId + `"
					principal_type                   = "ApplicationUser"
					role_definition_id               = "` + testRoleDefinitionId + `"
				}`,
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

// The RBAC API rejects a freshly created environment until its id propagates, so create is retried.
func TestUnitRoleAssignmentResource_Validate_Create_Retries_When_Scope_Not_Propagated(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	postAttempts := 0
	created := false
	httpmock.RegisterResponder("POST", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			postAttempts++
			if postAttempts == 1 {
				return httpmock.NewStringResponse(http.StatusBadRequest, `{"code":"EndpointInvalid","message":"Environment id is invalid.","innererror":{"code":"EnvironmentIdInvalid"}}`), nil
			}
			created = true
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/Validate_Create_Environment/post_role_assignment.json").String()), nil
		})
	httpmock.RegisterResponder("GET", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			if !created {
				return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Environment/get_role_assignments.json").String()), nil
		})
	httpmock.RegisterResponder("DELETE", environmentCollection+"/"+environmentAssignmentId+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_assignment" "test" {
					scope_type         = "environment"
					environment_id     = "` + testEnvironmentId + `"
					principal_id       = "` + testPrincipalId + `"
					principal_type     = "ApplicationUser"
					role_definition_id = "` + testRoleDefinitionId + `"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "id", environmentAssignmentId),
					func(_ *terraform.State) error {
						if postAttempts < 2 {
							return fmt.Errorf("expected the create request to be retried, got %d attempts", postAttempts)
						}
						return nil
					},
				),
			},
		},
	})
}

func TestUnitRoleAssignmentResource_Validate_Create_Error(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("POST", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusBadRequest, `{"code":"PrincipalDoesNotExist","message":"The service principal does not exist in tenant."}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_assignment" "test" {
					scope_type                       = "environment"
					environment_id                   = "` + testEnvironmentId + `"
					principal_id = "` + testPrincipalId + `"
					principal_type                   = "ApplicationUser"
					role_definition_id               = "` + testRoleDefinitionId + `"
				}`,
				ExpectError: regexp.MustCompile(`Failed to create role assignment at environment .* scope`),
			},
		},
	})
}

func TestUnitRoleAssignmentResource_Validate_Import_InvalidId(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()
	registerScopeMocks(environmentCollection, "tests/resource/Validate_Create_Environment", environmentAssignmentId)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_assignment" "test" {
					scope_type                       = "environment"
					environment_id                   = "` + testEnvironmentId + `"
					principal_id = "` + testPrincipalId + `"
					principal_type                   = "ApplicationUser"
					role_definition_id               = "` + testRoleDefinitionId + `"
				}`,
				ResourceName:  "powerplatform_role_assignment.test",
				ImportState:   true,
				ImportStateId: "environments/" + testEnvironmentId,
				ExpectError:   regexp.MustCompile(`Invalid import ID`),
			},
		},
	})
}

// acceptanceProviders are the external providers every acceptance case below needs.
func acceptanceProviders() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"azuread": {
			VersionConstraint: constants.AZURE_AD_PROVIDER_VERSION_CONSTRAINT,
			Source:            "hashicorp/azuread",
		},
		"time": {
			Source: "hashicorp/time",
		},
	}
}

// acceptancePreamble creates the service principal to assign to and resolves a role definition by
// name. scopeConfig adds whatever the scope under test needs, and scopeAttribute is spliced into the
// role assignment itself.
func acceptancePreamble(scopeConfig, scopeAttribute string) string {
	return `
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

		data "powerplatform_role_definitions" "all" {
		}

		locals {
			role_definition_id = "` + roleBasedAccessAdministratorRoleId + `"
		}
		` + scopeConfig + `

		resource "powerplatform_role_assignment" "test" {
			` + scopeAttribute + `
			principal_id       = azuread_service_principal.test_sp.object_id
			principal_type     = "ApplicationUser"
			role_definition_id = local.role_definition_id

			depends_on = [time_sleep.wait_for_service_principal]
		}`
}

// commonAcceptanceChecks are true of an assignment at any scope.
func commonAcceptanceChecks() []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		resource.TestMatchResourceAttr("powerplatform_role_assignment.test", "id", regexp.MustCompile(helpers.GuidRegex)),
		resource.TestCheckResourceAttrPair("powerplatform_role_assignment.test", "principal_id", "azuread_service_principal.test_sp", "object_id"),
		resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "principal_type", "ApplicationUser"),
		resource.TestCheckResourceAttrSet("powerplatform_role_assignment.test", "created_on"),
	}
}

// With neither identifier set the assignment lands on the tenant.
func TestAccRoleAssignmentResource_Validate_Create_Tenant_Scope(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        acceptanceProviders(),
		Steps: []resource.TestStep{
			{
				Config: acceptancePreamble("", `scope_type = "tenant"`),
				Check: resource.ComposeAggregateTestCheckFunc(append(commonAcceptanceChecks(),
					resource.TestMatchResourceAttr("powerplatform_role_assignment.test", "scope", regexp.MustCompile(`^/tenants/`)),
					resource.TestCheckNoResourceAttr("powerplatform_role_assignment.test", "environment_id"),
					resource.TestCheckNoResourceAttr("powerplatform_role_assignment.test", "environment_group_id"),
				)...),
			},
		},
	})
}

func TestAccRoleAssignmentResource_Validate_Create_Environment_Scope(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        acceptanceProviders(),
		Steps: []resource.TestStep{
			{
				Config: acceptancePreamble(`
		resource "powerplatform_environment" "test_environment" {
			display_name     = "`+mocks.TestName()+`"
			location         = "unitedstates"
			environment_type = "Sandbox"
		}`, "scope_type     = \"environment\"\n\t\t\tenvironment_id = powerplatform_environment.test_environment.id"),
				Check: resource.ComposeAggregateTestCheckFunc(append(commonAcceptanceChecks(),
					resource.TestCheckResourceAttrPair("powerplatform_role_assignment.test", "environment_id", "powerplatform_environment.test_environment", "id"),
					resource.TestMatchResourceAttr("powerplatform_role_assignment.test", "scope", regexp.MustCompile(`/environments/`)),
				)...),
			},
		},
	})
}

func TestAccRoleAssignmentResource_Validate_Create_EnvironmentGroup_Scope(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        acceptanceProviders(),
		Steps: []resource.TestStep{
			{
				Config: acceptancePreamble(`
		resource "powerplatform_environment_group" "test_env_group" {
			display_name = "`+mocks.TestName()+`"
			description  = "Environment group for role assignment acceptance test"
		}`, "scope_type           = \"environment_group\"\n\t\t\tenvironment_group_id = powerplatform_environment_group.test_env_group.id"),
				Check: resource.ComposeAggregateTestCheckFunc(append(commonAcceptanceChecks(),
					resource.TestCheckResourceAttrPair("powerplatform_role_assignment.test", "environment_group_id", "powerplatform_environment_group.test_env_group", "id"),
					resource.TestMatchResourceAttr("powerplatform_role_assignment.test", "scope", regexp.MustCompile(`(?i)/environmentgroups/`)),
				)...),
			},
		},
	})
}

// TestAccRoleAssignmentResource_Validate_All_Principal_Types_And_Scopes covers the full matrix in one
// apply: a user, a group and a service principal, assigned at tenant, environment and environment
// group scope.
//
// The identifier differs by principal type. Microsoft documents users as identified by email address,
// groups by group id, and service principals by their enterprise object id. See
// https://learn.microsoft.com/en-us/power-platform/admin/security/role-based-access-control.
// Unit tests only prove those type strings deserialise, so this is the test that proves the
// identifier contract is right.
func TestAccRoleAssignmentResource_Validate_All_Principal_Types_And_Scopes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azuread": {
				VersionConstraint: constants.AZURE_AD_PROVIDER_VERSION_CONSTRAINT,
				Source:            "hashicorp/azuread",
			},
			"random": {
				VersionConstraint: constants.RANDOM_PROVIDER_VERSION_CONSTRAINT,
				Source:            "hashicorp/random",
			},
			"time": {
				Source: "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
				data "azuread_domains" "aad_domains" {
					only_initial = true
				}

				locals {
					domain_name = data.azuread_domains.aad_domains.domains[0].domain_name
				}

				# --- the three principal kinds -------------------------------------------------

				resource "random_password" "user" {
					length           = 16
					min_lower        = 1
					min_upper        = 1
					min_numeric      = 1
					min_special      = 1
					special          = true
					override_special = "_%@"
				}

				resource "azuread_user" "test_user" {
					user_principal_name = "` + mocks.TestName() + `@${local.domain_name}"
					display_name        = "` + mocks.TestName() + `"
					mail_nickname       = "` + mocks.TestName() + `"
					password            = random_password.user.result
					usage_location      = "US"
				}

				resource "azuread_group" "test_group" {
					display_name     = "` + mocks.TestName() + `"
					security_enabled = true
				}

				resource "azuread_application_registration" "test_app" {
					display_name = "` + mocks.TestName() + `"
				}

				resource "azuread_service_principal" "test_sp" {
					client_id = azuread_application_registration.test_app.client_id
				}

				resource "time_sleep" "wait_for_principals" {
					create_duration = "60s"

					depends_on = [
						azuread_user.test_user,
						azuread_group.test_group,
						azuread_service_principal.test_sp,
					]
				}

				# --- the two non tenant scopes -------------------------------------------------

				resource "powerplatform_environment" "test_environment" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
				}

				resource "powerplatform_environment_group" "test_env_group" {
					display_name = "` + mocks.TestName() + `"
					description  = "Environment group for the role assignment matrix"
				}

				data "powerplatform_role_definitions" "all" {
				}

				locals {
					role_definition_id = "` + roleBasedAccessAdministratorRoleId + `"
				}

				# --- service principal, at all three scopes ------------------------------------

				resource "powerplatform_role_assignment" "sp_tenant" {
					scope_type         = "tenant"
					principal_id       = azuread_service_principal.test_sp.object_id
					principal_type     = "ApplicationUser"
					role_definition_id = local.role_definition_id

					depends_on = [time_sleep.wait_for_principals]
				}

				resource "powerplatform_role_assignment" "sp_environment" {
					scope_type         = "environment"
					environment_id     = powerplatform_environment.test_environment.id
					principal_id       = azuread_service_principal.test_sp.object_id
					principal_type     = "ApplicationUser"
					role_definition_id = local.role_definition_id

					depends_on = [time_sleep.wait_for_principals]
				}

				resource "powerplatform_role_assignment" "sp_environment_group" {
					scope_type           = "environment_group"
					environment_group_id = powerplatform_environment_group.test_env_group.id
					principal_id         = azuread_service_principal.test_sp.object_id
					principal_type       = "ApplicationUser"
					role_definition_id   = local.role_definition_id

					depends_on = [time_sleep.wait_for_principals]
				}

				# --- group, at environment scope ------------------------------------------------

				resource "powerplatform_role_assignment" "group_environment" {
					scope_type         = "environment"
					environment_id     = powerplatform_environment.test_environment.id
					principal_id       = azuread_group.test_group.object_id
					principal_type     = "Group"
					role_definition_id = local.role_definition_id

					depends_on = [time_sleep.wait_for_principals]
				}

				# --- user, at environment scope -------------------------------------------------
				#
				# Verified live 2026-08-20: the API's principalObjectId is a System.Guid, so a user
				# is identified by object id like every other principal. Microsoft's RBAC doc says
				# email, but that describes the portal; a UPN fails JSON conversion with a 400
				# before any lookup happens.

				resource "powerplatform_role_assignment" "user_environment" {
					scope_type         = "environment"
					environment_id     = powerplatform_environment.test_environment.id
					principal_id       = azuread_user.test_user.object_id
					principal_type     = "User"
					role_definition_id = local.role_definition_id

					depends_on = [time_sleep.wait_for_principals]
				}

				# --- the assignments are discoverable through the data source at their scopes -----

				data "powerplatform_role_assignments" "environment" {
					scope_type     = "environment"
					environment_id = powerplatform_environment.test_environment.id

					depends_on = [
						powerplatform_role_assignment.sp_environment,
						powerplatform_role_assignment.group_environment,
						powerplatform_role_assignment.user_environment,
					]
				}

				data "powerplatform_role_assignments" "environment_group" {
					scope_type           = "environment_group"
					environment_group_id = powerplatform_environment_group.test_env_group.id

					depends_on = [powerplatform_role_assignment.sp_environment_group]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// service principal, three scopes
					resource.TestMatchResourceAttr("powerplatform_role_assignment.sp_tenant", "scope", regexp.MustCompile(`^/tenants/[^/]+$`)),
					resource.TestMatchResourceAttr("powerplatform_role_assignment.sp_environment", "scope", regexp.MustCompile(`/environments/`)),
					resource.TestMatchResourceAttr("powerplatform_role_assignment.sp_environment_group", "scope", regexp.MustCompile(`(?i)/environmentgroups/`)),
					resource.TestCheckResourceAttrPair("powerplatform_role_assignment.sp_environment", "environment_id", "powerplatform_environment.test_environment", "id"),
					resource.TestCheckResourceAttrPair("powerplatform_role_assignment.sp_environment_group", "environment_group_id", "powerplatform_environment_group.test_env_group", "id"),

					// group
					resource.TestCheckResourceAttr("powerplatform_role_assignment.group_environment", "principal_type", "Group"),
					resource.TestCheckResourceAttrPair("powerplatform_role_assignment.group_environment", "principal_id", "azuread_group.test_group", "object_id"),
					resource.TestMatchResourceAttr("powerplatform_role_assignment.group_environment", "id", regexp.MustCompile(helpers.GuidRegex)),

					// user
					resource.TestCheckResourceAttr("powerplatform_role_assignment.user_environment", "principal_type", "User"),
					resource.TestMatchResourceAttr("powerplatform_role_assignment.user_environment", "id", regexp.MustCompile(helpers.GuidRegex)),

					// the data source finds the assignments at their scopes: three on the
					// environment (service principal, group, user) and one on the group
					resource.TestCheckResourceAttr("data.powerplatform_role_assignments.environment", "role_assignments.#", "3"),
					resource.TestCheckResourceAttr("data.powerplatform_role_assignments.environment_group", "role_assignments.#", "1"),

					resource.TestCheckResourceAttrSet("powerplatform_role_assignment.sp_tenant", "created_on"),
				),
			},
		},
	})
}

// The create POST is not idempotent and the API issues no correlation token, so when a transport
// failure hides the outcome the provider must not replay the POST, must not go looking for an
// assignment to adopt, and must not write state. The mock commits the assignment to prove the
// deliberate gap: the grant exists remotely while Terraform has claimed nothing.
func TestUnitRoleAssignmentResource_Validate_Create_Ambiguous_Commit_Is_Unknown_Outcome(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	envSuffix := "/environments/" + testEnvironmentId
	postAttempts := 0
	gets := 0
	committed := false
	httpmock.RegisterResponder("POST", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			postAttempts++
			// The server commits the assignment but the response never arrives.
			committed = true
			return nil, errors.New("connection reset by peer")
		})
	httpmock.RegisterResponder("GET", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			gets++
			if committed {
				return httpmock.NewStringResponse(http.StatusOK, `{"value":[`+assignmentJSON(environmentAssignmentId, testPrincipalId, testRoleDefinitionId, envSuffix, "")+`]}`), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_assignment" "test" {
					scope_type         = "environment"
					environment_id     = "` + testEnvironmentId + `"
					principal_id       = "` + testPrincipalId + `"
					principal_type     = "ApplicationUser"
					role_definition_id = "` + testRoleDefinitionId + `"
				}`,
				ExpectError: regexp.MustCompile(`(?s)the create outcome is\s+unknown`),
			},
			{
				Config: `# empty`,
				Check: resource.ComposeAggregateTestCheckFunc(
					func(state *terraform.State) error {
						if len(state.RootModule().Resources) != 0 {
							return fmt.Errorf("no assignment may reach state after an unknown outcome, found %d resources", len(state.RootModule().Resources))
						}
						if postAttempts != 1 {
							return fmt.Errorf("the non-idempotent create must be sent exactly once, got %d attempts", postAttempts)
						}
						if gets != 1 {
							return fmt.Errorf("nothing may be listed after the failed POST, got %d list calls", gets)
						}
						if !committed {
							return errors.New("the mock must hold the committed assignment for this test to prove the deliberate gap")
						}
						return nil
					},
				),
			},
		},
	})
}

// A relationship that already exists uniquely is adopted without a POST, since the configuration
// identifies the relationship and the assignment id is computed.
func TestUnitRoleAssignmentResource_Validate_Create_Adopts_Unique_Existing_Assignment(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	postAttempts := 0
	httpmock.RegisterResponder("POST", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			postAttempts++
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/Validate_Create_Environment/post_role_assignment.json").String()), nil
		})
	httpmock.RegisterResponder("GET", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Environment/get_role_assignments.json").String()), nil
		})
	httpmock.RegisterResponder("DELETE", environmentCollection+"/"+environmentAssignmentId+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_assignment" "test" {
					scope_type         = "environment"
					environment_id     = "` + testEnvironmentId + `"
					principal_id       = "` + testPrincipalId + `"
					principal_type     = "ApplicationUser"
					role_definition_id = "` + testRoleDefinitionId + `"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "id", environmentAssignmentId),
					func(_ *terraform.State) error {
						if postAttempts != 0 {
							return fmt.Errorf("an existing unique relationship must be adopted without a POST, got %d attempts", postAttempts)
						}
						return nil
					},
				),
			},
		},
	})
}

// Duplicate existing relationships cannot be adopted safely, so create refuses them.
func TestUnitRoleAssignmentResource_Validate_Create_Refuses_Duplicate_Existing_Assignments(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[
				{"roleAssignmentId":"33333333-3333-3333-3333-333333333333","scope":"/tenants/00000000-0000-0000-0000-000000000001/environments/`+testEnvironmentId+`","principalType":"ApplicationUser","principalObjectId":"`+testPrincipalId+`","roleDefinitionId":"`+testRoleDefinitionId+`","createdByPrincipalType":"User","createdByPrincipalObjectId":"cccccccc-cccc-cccc-cccc-cccccccccccc","createdOn":"2026-06-22T17:00:00Z","expiresOn":null},
				{"roleAssignmentId":"55555555-5555-5555-5555-555555555555","scope":"/tenants/00000000-0000-0000-0000-000000000001/environments/`+testEnvironmentId+`","principalType":"ApplicationUser","principalObjectId":"`+testPrincipalId+`","roleDefinitionId":"`+testRoleDefinitionId+`","createdByPrincipalType":"User","createdByPrincipalObjectId":"cccccccc-cccc-cccc-cccc-cccccccccccc","createdOn":"2026-06-22T17:01:00Z","expiresOn":null}
			]}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_assignment" "test" {
					scope_type         = "environment"
					environment_id     = "` + testEnvironmentId + `"
					principal_id       = "` + testPrincipalId + `"
					principal_type     = "ApplicationUser"
					role_definition_id = "` + testRoleDefinitionId + `"
				}`,
				ExpectError: regexp.MustCompile(`(?s)found 2 existing role assignments.*deduplicate`),
			},
		},
	})
}

// A definitive rejection never committed anything, so there is nothing to reconcile and no polling.
func TestUnitRoleAssignmentResource_Validate_Create_Definitive_Failure_Does_Not_Poll(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	gets := 0
	httpmock.RegisterResponder("GET", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			gets++
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
		})
	httpmock.RegisterResponder("POST", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusBadRequest, `{"code":"PrincipalDoesNotExist","message":"The service principal does not exist in tenant."}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_assignment" "test" {
					scope_type         = "environment"
					environment_id     = "` + testEnvironmentId + `"
					principal_id       = "` + testPrincipalId + `"
					principal_type     = "ApplicationUser"
					role_definition_id = "` + testRoleDefinitionId + `"
				}`,
				ExpectError: regexp.MustCompile(`PrincipalDoesNotExist`),
			},
			{
				Config: `# empty`,
				Check: resource.ComposeAggregateTestCheckFunc(
					func(_ *terraform.State) error {
						if gets != 1 {
							return fmt.Errorf("a definitive 400 must not trigger reconcile polling, got %d list calls", gets)
						}
						return nil
					},
				),
			},
		},
	})
}

// scope_type must arrive with its matching id: the per-attribute validators cannot catch an id
// that is absent entirely, so ValidateConfig covers that side.
func TestUnitRoleAssignmentResource_Validate_ScopeType_Requires_Its_Id(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_assignment" "missing_env_id" {
					scope_type         = "environment"
					principal_id       = "` + testPrincipalId + `"
					principal_type     = "ApplicationUser"
					role_definition_id = "` + testRoleDefinitionId + `"
				}`,
				ExpectError: regexp.MustCompile(`environment_id is required when scope_type`),
			},
			{
				Config: `
				resource "powerplatform_role_assignment" "missing_group_id" {
					scope_type         = "environment_group"
					principal_id       = "` + testPrincipalId + `"
					principal_type     = "ApplicationUser"
					role_definition_id = "` + testRoleDefinitionId + `"
				}`,
				ExpectError: regexp.MustCompile(`environment_group_id is required when scope_type`),
			},
			{
				Config: `
				resource "powerplatform_role_assignment" "tenant_with_id" {
					scope_type         = "tenant"
					environment_id     = "` + testEnvironmentId + `"
					principal_id       = "` + testPrincipalId + `"
					principal_type     = "ApplicationUser"
					role_definition_id = "` + testRoleDefinitionId + `"
				}`,
				ExpectError: regexp.MustCompile(`must not be set when scope_type is .tenant.`),
			},
		},
	})
}

// A 404 from the tenant collection cannot mean a deleted scope, so it must stay an error rather
// than silently untracking an active tenant-wide grant. A 404 on a deleted environment is the
// opposite: the scope took its assignments with it, so the resource leaves state.
func TestUnitRoleAssignmentResource_Validate_Read_404_Tenant_Errors_Environment_Removes(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	// tenant scope: create fine, then the collection starts returning 404
	tenantGets := 0
	httpmock.RegisterResponder("POST", tenantCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/Validate_Create_Tenant/post_role_assignment.json").String()), nil
		})
	httpmock.RegisterResponder("GET", tenantCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			tenantGets++
			// preflight and the post-apply refresh see the assignment; the 404 starts afterwards
			if tenantGets <= 2 {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Tenant/get_role_assignments.json").String()), nil
			}
			return httpmock.NewStringResponse(http.StatusNotFound, ""), nil
		})
	httpmock.RegisterResponder("DELETE", tenantCollection+"/"+tenantAssignmentId+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_assignment" "test" {
					scope_type         = "tenant"
					principal_id       = "` + testPrincipalId + `"
					principal_type     = "ApplicationUser"
					role_definition_id = "` + testRoleDefinitionId + `"
				}`,
			},
			{
				RefreshState: true,
				ExpectError:  regexp.MustCompile(`(?s)Failed to list role assignments at tenant scope.*cannot\s+mean\s+a\s+deleted\s+scope`),
			},
		},
	})
}

func TestUnitRoleAssignmentResource_Validate_Read_Removes_State_When_Environment_Deleted(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	envGets := 0
	httpmock.RegisterResponder("POST", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/Validate_Create_Environment/post_role_assignment.json").String()), nil
		})
	httpmock.RegisterResponder("GET", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			envGets++
			// preflight and the post-apply refresh see the assignment; then the environment is
			// deleted out of band
			if envGets <= 2 {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Environment/get_role_assignments.json").String()), nil
			}
			return httpmock.NewStringResponse(http.StatusNotFound, ""), nil
		})
	httpmock.RegisterResponder("DELETE", environmentCollection+"/"+environmentAssignmentId+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_assignment" "test" {
					scope_type         = "environment"
					environment_id     = "` + testEnvironmentId + `"
					principal_id       = "` + testPrincipalId + `"
					principal_type     = "ApplicationUser"
					role_definition_id = "` + testRoleDefinitionId + `"
				}`,
			},
			{
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// The API's principalObjectId is a System.Guid, verified live: an email fails JSON conversion with
// a 400 before any principal lookup. The UUID custom type rejects it at plan time instead.
func TestUnitRoleAssignmentResource_Validate_PrincipalId_Must_Be_A_Guid(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_assignment" "test" {
					scope_type         = "environment"
					environment_id     = "` + testEnvironmentId + `"
					principal_id       = "someone@contoso.com"
					principal_type     = "User"
					role_definition_id = "` + testRoleDefinitionId + `"
				}`,
				ExpectError: regexp.MustCompile(`Invalid UUID String Value`),
			},
		},
	})
}

// Adopting or refusing duplicates is only sound against a real listing, so a failed preflight
// stops the create before any POST.
func TestUnitRoleAssignmentResource_Validate_Create_Fails_When_Preflight_List_Fails(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	postAttempts := 0
	httpmock.RegisterResponder("GET", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusForbidden, `{"error":"forbidden"}`), nil
		})
	httpmock.RegisterResponder("POST", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			postAttempts++
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/Validate_Create_Environment/post_role_assignment.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_assignment" "test" {
					scope_type         = "environment"
					environment_id     = "` + testEnvironmentId + `"
					principal_id       = "` + testPrincipalId + `"
					principal_type     = "ApplicationUser"
					role_definition_id = "` + testRoleDefinitionId + `"
				}`,
				ExpectError: regexp.MustCompile(`(?s)could not list the existing role assignments before\s+creating`),
			},
			{
				Config: `# empty`,
				Check: resource.ComposeAggregateTestCheckFunc(
					func(_ *terraform.State) error {
						if postAttempts != 0 {
							return fmt.Errorf("create must not POST when the preflight list fails, got %d attempts", postAttempts)
						}
						return nil
					},
				),
			},
		},
	})
}

// assignmentJSON renders one assignment for list fixtures built inline.
func assignmentJSON(id, principal, role, scopeSuffix, expiresOn string) string {
	expires := "null"
	if expiresOn != "" {
		expires = `"` + expiresOn + `"`
	}
	return `{"roleAssignmentId":"` + id + `","scope":"/tenants/00000000-0000-0000-0000-000000000001` + scopeSuffix + `","principalType":"ApplicationUser","principalObjectId":"` + principal + `","roleDefinitionId":"` + role + `","createdByPrincipalType":"User","createdByPrincipalObjectId":"cccccccc-cccc-cccc-cccc-cccccccccccc","createdOn":"2026-06-22T17:00:00Z","expiresOn":` + expires + `}`
}

// Every status that hides the create outcome (timeout, throttle, server error) yields the explicit
// unknown-outcome error: exactly one POST, no listing afterwards, nothing adopted.
func TestUnitRoleAssignmentResource_Validate_Create_Retryable_Statuses_Are_Unknown_Outcome(t *testing.T) {
	for _, status := range []int{http.StatusRequestTimeout, http.StatusTooManyRequests, http.StatusBadGateway} {
		t.Run(fmt.Sprintf("status_%d", status), func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			mocks.ActivateEnvironmentHttpMocks()

			postAttempts := 0
			gets := 0
			httpmock.RegisterResponder("GET", environmentCollection+apiVersionQuery,
				func(_ *http.Request) (*http.Response, error) {
					gets++
					return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
				})
			httpmock.RegisterResponder("POST", environmentCollection+apiVersionQuery,
				func(_ *http.Request) (*http.Response, error) {
					postAttempts++
					return httpmock.NewStringResponse(status, `{"error":"outcome hidden"}`), nil
				})

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: `
						resource "powerplatform_role_assignment" "test" {
							scope_type         = "environment"
							environment_id     = "` + testEnvironmentId + `"
							principal_id       = "` + testPrincipalId + `"
							principal_type     = "ApplicationUser"
							role_definition_id = "` + testRoleDefinitionId + `"
						}`,
						ExpectError: regexp.MustCompile(`(?s)the create outcome is\s+unknown`),
					},
					{
						Config: `# empty`,
						Check: func(_ *terraform.State) error {
							if postAttempts != 1 {
								return fmt.Errorf("the non-idempotent create must be sent exactly once, got %d attempts", postAttempts)
							}
							if gets != 1 {
								return fmt.Errorf("nothing may be listed after the failed POST, got %d list calls", gets)
							}
							return nil
						},
					},
				},
			})
		})
	}
}

// The documented recovery from an unknown outcome: the operator inspects the scope, finds the
// committed assignment, and imports it with the scope-shaped id. From there the ordinary lifecycle
// owns it, and destroy removes exactly that assignment.
func TestUnitRoleAssignmentResource_Validate_Import_After_Ambiguous_Commit_Succeeds(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	envSuffix := "/environments/" + testEnvironmentId
	committed := false
	httpmock.RegisterResponder("GET", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			if committed {
				return httpmock.NewStringResponse(http.StatusOK, `{"value":[`+assignmentJSON(environmentAssignmentId, testPrincipalId, testRoleDefinitionId, envSuffix, "")+`]}`), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
		})
	httpmock.RegisterResponder("POST", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			// The server commits the assignment but answers with a server error.
			committed = true
			return httpmock.NewStringResponse(http.StatusInternalServerError, `{"error":"boom"}`), nil
		})
	deleted := ""
	httpmock.RegisterResponder("DELETE", environmentCollection+"/"+environmentAssignmentId+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			deleted = environmentAssignmentId
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	config := `
	resource "powerplatform_role_assignment" "test" {
		scope_type         = "environment"
		environment_id     = "` + testEnvironmentId + `"
		principal_id       = "` + testPrincipalId + `"
		principal_type     = "ApplicationUser"
		role_definition_id = "` + testRoleDefinitionId + `"
	}`

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`(?s)the create outcome is\s+unknown`),
			},
			{
				Config:             config,
				ResourceName:       "powerplatform_role_assignment.test",
				ImportState:        true,
				ImportStateId:      "environments/" + testEnvironmentId + "/" + environmentAssignmentId,
				ImportStatePersist: true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected exactly one imported instance, got %d", len(states))
					}
					if states[0].Attributes["id"] != environmentAssignmentId {
						return fmt.Errorf("expected the committed assignment %s to be imported, got %s", environmentAssignmentId, states[0].Attributes["id"])
					}
					return nil
				},
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if deleted != environmentAssignmentId {
				return fmt.Errorf("destroy must remove exactly the imported assignment, deleted '%s'", deleted)
			}
			return nil
		},
	})
}

// An expiring assignment is a relationship this resource cannot represent, so it is never adopted:
// the create proceeds and makes the permanent assignment the configuration declares.
func TestUnitRoleAssignmentResource_Validate_Create_Does_Not_Adopt_Expiring_Assignment(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	envSuffix := "/environments/" + testEnvironmentId
	postAttempts := 0
	created := false
	httpmock.RegisterResponder("GET", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			rows := assignmentJSON("44444444-4444-4444-4444-444444444444", testPrincipalId, testRoleDefinitionId, envSuffix, "2027-01-01T00:00:00Z")
			if created {
				rows += "," + assignmentJSON(environmentAssignmentId, testPrincipalId, testRoleDefinitionId, envSuffix, "")
			}
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[`+rows+`]}`), nil
		})
	httpmock.RegisterResponder("POST", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			postAttempts++
			created = true
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/Validate_Create_Environment/post_role_assignment.json").String()), nil
		})
	httpmock.RegisterResponder("DELETE", environmentCollection+"/"+environmentAssignmentId+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_assignment" "test" {
					scope_type         = "environment"
					environment_id     = "` + testEnvironmentId + `"
					principal_id       = "` + testPrincipalId + `"
					principal_type     = "ApplicationUser"
					role_definition_id = "` + testRoleDefinitionId + `"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "id", environmentAssignmentId),
					func(_ *terraform.State) error {
						if postAttempts != 1 {
							return fmt.Errorf("the expiring assignment must not be adopted; expected one POST, got %d", postAttempts)
						}
						return nil
					},
				),
			},
		},
	})
}

// Import ids must be guids in every guid position.
func TestUnitRoleAssignmentResource_Validate_Import_Rejects_Malformed_Guids(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()
	registerScopeMocks(environmentCollection, "tests/resource/Validate_Create_Environment", environmentAssignmentId)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_assignment" "test" {
					scope_type         = "environment"
					environment_id     = "` + testEnvironmentId + `"
					principal_id       = "` + testPrincipalId + `"
					principal_type     = "ApplicationUser"
					role_definition_id = "` + testRoleDefinitionId + `"
				}`,
				ResourceName:  "powerplatform_role_assignment.test",
				ImportState:   true,
				ImportStateId: "environments/not-a-guid/" + environmentAssignmentId,
				ExpectError:   regexp.MustCompile(`(?s)Invalid import ID.*Every\s+id\s+segment\s+must\s+be\s+a\s+guid`),
			},
		},
	})
}

// A cancellation during the POST hides the outcome exactly like a transport failure: the result is
// the explicit unknown-outcome error, never a success and never an adoption.
func TestUnitRoleAssignmentResource_Validate_Create_Cancellation_During_Post_Is_Unknown_Outcome(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	gets := 0
	httpmock.RegisterResponder("GET", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			gets++
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
		})
	httpmock.RegisterResponder("POST", environmentCollection+apiVersionQuery,
		httpmock.NewErrorResponder(context.Canceled))

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_assignment" "test" {
					scope_type         = "environment"
					environment_id     = "` + testEnvironmentId + `"
					principal_id       = "` + testPrincipalId + `"
					principal_type     = "ApplicationUser"
					role_definition_id = "` + testRoleDefinitionId + `"
				}`,
				ExpectError: regexp.MustCompile(`(?s)the create outcome is\s+unknown`),
			},
			{
				Config: `# empty`,
				Check: func(_ *terraform.State) error {
					if gets != 1 {
						return fmt.Errorf("nothing may be listed after the cancelled POST, got %d list calls", gets)
					}
					return nil
				},
			},
		},
	})
}

// Reserved path words are only reserved in their own position; used anywhere else they are just
// malformed ids and the import is rejected.
func TestUnitRoleAssignmentResource_Validate_Import_Rejects_Reserved_Word_Bypasses(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()
	registerScopeMocks(environmentCollection, "tests/resource/Validate_Create_Environment", environmentAssignmentId)

	config := `
	resource "powerplatform_role_assignment" "test" {
		scope_type         = "environment"
		environment_id     = "` + testEnvironmentId + `"
		principal_id       = "` + testPrincipalId + `"
		principal_type     = "ApplicationUser"
		role_definition_id = "` + testRoleDefinitionId + `"
	}`

	for _, importId := range []string{
		"environments",
		"environments/" + testEnvironmentId + "/environmentGroups",
		"environmentGroups/environments/" + environmentAssignmentId,
	} {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config:        config,
					ResourceName:  "powerplatform_role_assignment.test",
					ImportState:   true,
					ImportStateId: importId,
					ExpectError:   regexp.MustCompile(`Invalid import ID`),
				},
			},
		})
	}
}
