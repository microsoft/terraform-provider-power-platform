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

const roleBasedAccessAdministratorRoleName = "Power Platform Role Based Access Control Administrator"

func TestAccRoleBasedAccessAssignmentResource_Validate_Create(t *testing.T) {
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

				data "powerplatform_role_definitions" "all" {
				}

				locals {
					role_definition_id = [
						for role in data.powerplatform_role_definitions.all.role_definitions :
						role.role_definition_id if role.role_definition_name == "` + roleBasedAccessAdministratorRoleName + `"
					][0]
				}

				resource "powerplatform_role_based_access_assignment" "test" {
					enterprise_application_object_id = azuread_service_principal.test_sp.object_id
					principal_type                   = "ApplicationUser"
					role_definition_id               = local.role_definition_id

					depends_on = [time_sleep.wait_for_service_principal]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_role_based_access_assignment.test", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttrPair("powerplatform_role_based_access_assignment.test", "enterprise_application_object_id", "azuread_service_principal.test_sp", "object_id"),
					resource.TestCheckResourceAttr("powerplatform_role_based_access_assignment.test", "principal_type", "ApplicationUser"),
					resource.TestMatchResourceAttr("powerplatform_role_based_access_assignment.test", "scope", regexp.MustCompile(`^/tenants/`)),
					resource.TestCheckResourceAttrSet("powerplatform_role_based_access_assignment.test", "created_on"),
				),
			},
		},
	})
}

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
					enterprise_application_object_id = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
					principal_type                   = "ApplicationUser"
					role_definition_id               = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_role_based_access_assignment.test", "id", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckResourceAttr("powerplatform_role_based_access_assignment.test", "enterprise_application_object_id", "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
					resource.TestCheckResourceAttr("powerplatform_role_based_access_assignment.test", "principal_type", "ApplicationUser"),
					resource.TestCheckResourceAttr("powerplatform_role_based_access_assignment.test", "role_definition_id", "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
					resource.TestCheckResourceAttr("powerplatform_role_based_access_assignment.test", "scope", "/tenants/00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_role_based_access_assignment.test", "created_on", "2026-06-22T15:09:35Z"),
				),
			},
			{
				ResourceName:      "powerplatform_role_based_access_assignment.test",
				ImportState:       true,
				ImportStateId:     "11111111-1111-1111-1111-111111111111",
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitRoleBasedAccessAssignmentResource_Validate_Create_Error(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("POST", `https://api.powerplatform.com/authorization/roleAssignments?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusForbidden, `{"error":{"code":"Forbidden","message":"Access denied"}}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_based_access_assignment" "test" {
					enterprise_application_object_id = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
					principal_type                   = "ApplicationUser"
					role_definition_id               = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
				}`,
				ExpectError: regexp.MustCompile(`Failed to create role assignment`),
			},
		},
	})
}
