// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access_test

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

const roleBasedAccessAdministratorRoleName = "Power Platform Role Based Access Control Administrator"

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
					principal_id = "` + testPrincipalId + `"
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

func TestAccRoleAssignmentResource_Validate_Create_Environment_Scope(t *testing.T) {
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

				data "powerplatform_role_definitions" "all" {
				}

				locals {
					role_definition_id = [
						for role in data.powerplatform_role_definitions.all.role_definitions :
						role.role_definition_id if role.role_definition_name == "` + roleBasedAccessAdministratorRoleName + `"
					][0]
				}

				resource "powerplatform_role_assignment" "test" {
					environment_id                   = powerplatform_environment.test_environment.id
					principal_id = azuread_service_principal.test_sp.object_id
					principal_type                   = "ApplicationUser"
					role_definition_id               = local.role_definition_id

					depends_on = [time_sleep.wait_for_service_principal]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_role_assignment.test", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttrPair("powerplatform_role_assignment.test", "environment_id", "powerplatform_environment.test_environment", "id"),
					resource.TestMatchResourceAttr("powerplatform_role_assignment.test", "scope", regexp.MustCompile(`/environments/`)),
					resource.TestCheckResourceAttrSet("powerplatform_role_assignment.test", "created_on"),
				),
			},
		},
	})
}

// An unrecognised principal type is rejected at plan time rather than by the API at apply time.
func TestUnitRoleAssignmentResource_Validate_PrincipalType_Is_Enumerated(t *testing.T) {
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
					environment_id     = "` + testEnvironmentId + `"
					principal_id       = "` + testPrincipalId + `"
					principal_type     = "ServicePrincipal"
					role_definition_id = "` + testRoleDefinitionId + `"
				}`,
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}
