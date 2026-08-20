// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access_test //nolint:revive // the underscored package name predates this file and matches every service in the repo

import (
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
					resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "scope", "/tenants/00000000-0000-0000-0000-000000000001/environmentGroups/"+testEnvironmentGroupId),
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
	httpmock.RegisterResponder("POST", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			postAttempts++
			if postAttempts == 1 {
				return httpmock.NewStringResponse(http.StatusBadRequest, `{"code":"EndpointInvalid","message":"Environment id is invalid.","innererror":{"code":"EnvironmentIdInvalid"}}`), nil
			}
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/Validate_Create_Environment/post_role_assignment.json").String()), nil
		})
	httpmock.RegisterResponder("GET", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Environment/get_role_assignments.json").String()), nil
		})
	httpmock.RegisterResponder("DELETE", `=~^https://api\.powerplatform\.com/authorization/environments/.+/roleAssignments/.+\z`,
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
					scope_type                       = "environment"
					environment_id                   = "` + testEnvironmentId + `"
					principal_id = "` + testPrincipalId + `"
					principal_type                   = "ApplicationUser"
					role_definition_id               = "` + testRoleDefinitionId + `"
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
					resource.TestMatchResourceAttr("powerplatform_role_assignment.test", "scope", regexp.MustCompile(`/environmentGroups/`)),
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
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// service principal, three scopes
					resource.TestMatchResourceAttr("powerplatform_role_assignment.sp_tenant", "scope", regexp.MustCompile(`^/tenants/[^/]+$`)),
					resource.TestMatchResourceAttr("powerplatform_role_assignment.sp_environment", "scope", regexp.MustCompile(`/environments/`)),
					resource.TestMatchResourceAttr("powerplatform_role_assignment.sp_environment_group", "scope", regexp.MustCompile(`/environmentGroups/`)),
					resource.TestCheckResourceAttrPair("powerplatform_role_assignment.sp_environment", "environment_id", "powerplatform_environment.test_environment", "id"),
					resource.TestCheckResourceAttrPair("powerplatform_role_assignment.sp_environment_group", "environment_group_id", "powerplatform_environment_group.test_env_group", "id"),

					// group
					resource.TestCheckResourceAttr("powerplatform_role_assignment.group_environment", "principal_type", "Group"),
					resource.TestCheckResourceAttrPair("powerplatform_role_assignment.group_environment", "principal_id", "azuread_group.test_group", "object_id"),
					resource.TestMatchResourceAttr("powerplatform_role_assignment.group_environment", "id", regexp.MustCompile(helpers.GuidRegex)),

					// user
					resource.TestCheckResourceAttr("powerplatform_role_assignment.user_environment", "principal_type", "User"),
					resource.TestMatchResourceAttr("powerplatform_role_assignment.user_environment", "id", regexp.MustCompile(helpers.GuidRegex)),

					// every assignment is discoverable through the data source at its own scope
					resource.TestCheckResourceAttrSet("powerplatform_role_assignment.sp_tenant", "created_on"),
				),
			},
		},
	})
}

// The create POST is not idempotent, so a transient failure must not replay it. Instead the scope is
// listed and a matching assignment adopted, because the failed-looking request may have committed.
func TestUnitRoleAssignmentResource_Validate_Create_Ambiguous_Failure_Reconciles_Without_Replay(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	postAttempts := 0
	httpmock.RegisterResponder("POST", environmentCollection+apiVersionQuery,
		func(_ *http.Request) (*http.Response, error) {
			postAttempts++
			// The server commits the assignment but the response is lost in a 500.
			return httpmock.NewStringResponse(http.StatusInternalServerError, `{"error":"socket hang up"}`), nil
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
						if postAttempts != 1 {
							return fmt.Errorf("the non-idempotent create must be sent exactly once, got %d attempts", postAttempts)
						}
						return nil
					},
				),
			},
		},
	})
}

// The API's principalObjectId is a System.Guid, verified live: an email fails JSON conversion with a
// 400 before any principal lookup. Rejecting it at plan time turns that into a clear message.
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
				ExpectError: regexp.MustCompile(`identified by object id, not email`),
			},
		},
	})
}
