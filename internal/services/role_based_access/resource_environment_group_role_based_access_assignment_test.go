// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access_test

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

func TestAccEnvironmentGroupRoleBasedAccessAssignmentResource_Validate_Create(t *testing.T) {
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
					description  = "Environment group for role assignment acceptance test"
				}

				data "powerplatform_role_definitions" "all" {
				}

				locals {
					role_definition_id = [
						for role in data.powerplatform_role_definitions.all.role_definitions :
						role.role_definition_id if role.role_definition_name == "` + roleBasedAccessAdministratorRoleName + `"
					][0]
				}

				resource "powerplatform_environment_group_role_based_access_assignment" "test" {
					environment_group_id             = powerplatform_environment_group.test_env_group.id
					enterprise_application_object_id = azuread_service_principal.test_sp.object_id
					principal_type                   = "ApplicationUser"
					role_definition_id               = local.role_definition_id

					depends_on = [time_sleep.wait_for_service_principal]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment_group_role_based_access_assignment.test", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttrPair("powerplatform_environment_group_role_based_access_assignment.test", "environment_group_id", "powerplatform_environment_group.test_env_group", "id"),
					resource.TestCheckResourceAttrPair("powerplatform_environment_group_role_based_access_assignment.test", "enterprise_application_object_id", "azuread_service_principal.test_sp", "object_id"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_role_based_access_assignment.test", "principal_type", "ApplicationUser"),
					resource.TestMatchResourceAttr("powerplatform_environment_group_role_based_access_assignment.test", "scope", regexp.MustCompile(`/environmentGroups/`)),
					resource.TestCheckResourceAttrSet("powerplatform_environment_group_role_based_access_assignment.test", "created_on"),
				),
			},
		},
	})
}

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
					environment_group_id             = "dddddddd-dddd-dddd-dddd-dddddddddddd"
					enterprise_application_object_id = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
					principal_type                   = "ApplicationUser"
					role_definition_id               = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group_role_based_access_assignment.test", "id", "22222222-2222-2222-2222-222222222222"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_role_based_access_assignment.test", "environment_group_id", "dddddddd-dddd-dddd-dddd-dddddddddddd"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_role_based_access_assignment.test", "enterprise_application_object_id", "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_role_based_access_assignment.test", "principal_type", "ApplicationUser"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_role_based_access_assignment.test", "role_definition_id", "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_role_based_access_assignment.test", "scope", "/tenants/00000000-0000-0000-0000-000000000001/environmentGroups/dddddddd-dddd-dddd-dddd-dddddddddddd"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_role_based_access_assignment.test", "created_on", "2026-06-22T16:00:00Z"),
				),
			},
			{
				ResourceName:      "powerplatform_environment_group_role_based_access_assignment.test",
				ImportState:       true,
				ImportStateId:     "dddddddd-dddd-dddd-dddd-dddddddddddd/22222222-2222-2222-2222-222222222222",
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitEnvironmentGroupRoleBasedAccessAssignmentResource_Validate_Import_InvalidId(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_group_role_based_access_assignment" "test" {
					environment_group_id             = "dddddddd-dddd-dddd-dddd-dddddddddddd"
					enterprise_application_object_id = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
					principal_type                   = "ApplicationUser"
					role_definition_id               = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
				}`,
				ResourceName:  "powerplatform_environment_group_role_based_access_assignment.test",
				ImportState:   true,
				ImportStateId: "22222222-2222-2222-2222-222222222222",
				ExpectError:   regexp.MustCompile(`Invalid import ID`),
			},
		},
	})
}
