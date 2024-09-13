// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package admin_management_application_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccAdminManagementApplicationResource_Validate_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azuread": {
				VersionConstraint: ">= 2.53.1",
				Source: 		 "hashicorp/azuread",
			},
		},
		Steps: []resource.TestStep{
			{
				ResourceName: "powerplatform_admin_management_application.example_registration",
				Config: `
				resource "azuread_application_registration" "example_app" {
					display_name = "TestAccAdminManagementApplicationResource Application"
				}

				resource "azuread_service_principal" "example_sp" {
					client_id = azuread_application_registration.example_app.client_id
				}

				resource "powerplatform_admin_management_application" "example_registration" {
					id = azuread_application_registration.example_app.client_id
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("powerplatform_admin_management_application.example_registration", "id"),
				),
			},
		},
	})
}

func TestUnitAdminManagementApplicationResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	
	// zerod guid
	httpmock.RegisterResponder(http.MethodPut, "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/adminApplications/00000000-0000-0000-0000-000000000001?api-version=2020-10-01", httpmock.NewStringResponder(http.StatusOK, "{ \"applicationId\": \"00000000-0000-0000-0000-000000000001\" }"))
	httpmock.RegisterResponder(http.MethodGet, "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/adminApplications/00000000-0000-0000-0000-000000000001?api-version=2020-10-01", httpmock.NewStringResponder(http.StatusOK, "{ \"applicationId\": \"00000000-0000-0000-0000-000000000001\" }"))
	httpmock.RegisterResponder(http.MethodDelete, "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/adminApplications/00000000-0000-0000-0000-000000000001?api-version=2020-10-01", httpmock.NewStringResponder(http.StatusNoContent, ""))

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_admin_management_application" "example_registration" {
					id = "00000000-0000-0000-0000-000000000001"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_admin_management_application.example_registration", "id", constants.TEST_UUID),
				),
			},
		},
	})
}
