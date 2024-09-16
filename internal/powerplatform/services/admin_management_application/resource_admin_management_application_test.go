// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package admin_management_application_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccAdminManagementApplicationResource_Validate_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azuread": {
				VersionConstraint: ">= 2.53.1",
				Source:            "hashicorp/azuread",
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

	t.Setenv("client_id", "00000000-0000-0000-0000-000000000000")
	t.Setenv("client_secret", "00000000-0000-0000-0000-000000000000")
	t.Setenv("tenant_id", "00000000-0000-0000-0000-000000000000")

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client_id, _ := uuid.NewRandom()

	url := fmt.Sprintf("https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/adminApplications/%s?api-version=2020-10-01", client_id.String())
	body := fmt.Sprintf("{ \"applicationId\": \"%s\" }", client_id.String())

	httpmock.RegisterResponder(http.MethodPut, url, httpmock.NewStringResponder(http.StatusOK, body))
	httpmock.RegisterResponder(http.MethodGet, url, httpmock.NewStringResponder(http.StatusOK, body))
	httpmock.RegisterResponder(http.MethodDelete, url, httpmock.NewStringResponder(http.StatusNoContent, ""))

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{

				ConfigVariables: config.Variables{
					"client_id": config.StringVariable(client_id.String()),
				},
				Config: `
				variable "client_id" {
					type = string
				}
					
				resource "powerplatform_admin_management_application" "example_registration" {
					id = var.client_id
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_admin_management_application.example_registration", "id", client_id.String()),
				),
			},
		},
	})
}
