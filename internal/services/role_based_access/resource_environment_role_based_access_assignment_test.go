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
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitEnvironmentRoleBasedAccessAssignmentResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("POST", `https://api.powerplatform.com/authorization/environments/eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee/roleAssignments?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/Validate_Create_Environment/post_role_assignment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/authorization/environments/eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee/roleAssignments?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Environment/get_role_assignments.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", `https://api.powerplatform.com/authorization/environments/eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee/roleAssignments/33333333-3333-3333-3333-333333333333?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_role_based_access_assignment" "test" {
					environment_id                   = "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
					enterprise_application_object_id = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
					principal_type                   = "ApplicationUser"
					role_definition_id               = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_role_based_access_assignment.test", "id", "33333333-3333-3333-3333-333333333333"),
					resource.TestCheckResourceAttr("powerplatform_environment_role_based_access_assignment.test", "environment_id", "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"),
					resource.TestCheckResourceAttr("powerplatform_environment_role_based_access_assignment.test", "enterprise_application_object_id", "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
					resource.TestCheckResourceAttr("powerplatform_environment_role_based_access_assignment.test", "principal_type", "ApplicationUser"),
					resource.TestCheckResourceAttr("powerplatform_environment_role_based_access_assignment.test", "role_definition_id", "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
					resource.TestCheckResourceAttr("powerplatform_environment_role_based_access_assignment.test", "scope", "/tenants/00000000-0000-0000-0000-000000000001/environments/eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"),
					resource.TestCheckResourceAttr("powerplatform_environment_role_based_access_assignment.test", "created_on", "2026-06-22T17:00:00Z"),
				),
			},
			{
				ResourceName:      "powerplatform_environment_role_based_access_assignment.test",
				ImportState:       true,
				ImportStateId:     "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee/33333333-3333-3333-3333-333333333333",
				ImportStateVerify: true,
			},
		},
	})
}

// The RBAC API rejects a freshly created environment until its id propagates, so create is retried.
func TestUnitEnvironmentRoleBasedAccessAssignmentResource_Validate_Create_Retries_When_Environment_Not_Propagated(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	postAttempts := 0
	httpmock.RegisterResponder("POST", `https://api.powerplatform.com/authorization/environments/eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee/roleAssignments?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			postAttempts++
			if postAttempts == 1 {
				return httpmock.NewStringResponse(http.StatusBadRequest, `{"code":"EndpointInvalid","message":"Environment id eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee is invalid.","innererror":{"code":"EnvironmentIdInvalid"}}`), nil
			}
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/Validate_Create_Environment/post_role_assignment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/authorization/environments/eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee/roleAssignments?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Environment/get_role_assignments.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", `https://api.powerplatform.com/authorization/environments/eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee/roleAssignments/33333333-3333-3333-3333-333333333333?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_role_based_access_assignment" "test" {
					environment_id                   = "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
					enterprise_application_object_id = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
					principal_type                   = "ApplicationUser"
					role_definition_id               = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_role_based_access_assignment.test", "id", "33333333-3333-3333-3333-333333333333"),
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

func TestUnitEnvironmentRoleBasedAccessAssignmentResource_Validate_Create_Error(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("POST", `https://api.powerplatform.com/authorization/environments/eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee/roleAssignments?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusBadRequest, `{"code":"PrincipalDoesNotExist","message":"The service principal does not exist in tenant."}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_role_based_access_assignment" "test" {
					environment_id                   = "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
					enterprise_application_object_id = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
					principal_type                   = "ApplicationUser"
					role_definition_id               = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
				}`,
				ExpectError: regexp.MustCompile(`Failed to create environment role assignment`),
			},
		},
	})
}

func TestUnitEnvironmentRoleBasedAccessAssignmentResource_Validate_Import_InvalidId(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_role_based_access_assignment" "test" {
					environment_id                   = "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
					enterprise_application_object_id = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
					principal_type                   = "ApplicationUser"
					role_definition_id               = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
				}`,
				ResourceName:  "powerplatform_environment_role_based_access_assignment.test",
				ImportState:   true,
				ImportStateId: "33333333-3333-3333-3333-333333333333",
				ExpectError:   regexp.MustCompile(`Invalid import ID`),
			},
		},
	})
}
