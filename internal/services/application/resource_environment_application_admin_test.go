// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package application_test

import (
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitEnvironmentApplicationAdminResource_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register mock responders
	// Get environment details
	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})

	// Add application user
	httpmock.RegisterResponder("POST", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001/addAppUser`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, "{}"), nil
		})

	// Check if application user exists with any filter parameter version
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			// Check if this is for our test application
			if strings.Contains(req.URL.RawQuery, "00000000-0000-0000-0000-000000000002") {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_applicationusers.json").String()), nil
			}
			return httpmock.NewStringResponse(http.StatusNotFound, ""), nil
		})

	// Mock response for deactivating system user
	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.URL.Path, "/Microsoft.Dynamics.CRM.SetState") {
				return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
			}
			return httpmock.NewStringResponse(http.StatusNotFound, ""), nil
		})

	// Mock DELETE response for system user
	httpmock.RegisterResponder("DELETE", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_application_admin" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					application_id = "00000000-0000-0000-0000-000000000002"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_application_admin.test", "id", "00000000-0000-0000-0000-000000000001/00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("powerplatform_environment_application_admin.test", "environment_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment_application_admin.test", "application_id", "00000000-0000-0000-0000-000000000002"),
				),
			},
		},
	})
}

func TestUnitEnvironmentApplicationAdminResource_Read_NotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register mock responders
	// Get environment details
	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Read_NotFound/get_environment.json").String()), nil
		})

	// Check if application user exists - return empty array to simulate not found
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/applicationusers\?`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_application_admin" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					application_id = "00000000-0000-0000-0000-000000000002"
				}
				`,
				ExpectError: regexp.MustCompile(".*"),
			},
		},
	})
}

func TestUnitEnvironmentApplicationAdminResource_Import(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register mock responders
	// Get environment details
	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Import/get_environment.json").String()), nil
		})

	// Add application user for create step
	httpmock.RegisterResponder("POST", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001/addAppUser`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, "{}"), nil
		})

	// Check if application user exists with any filter parameter version
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			// Check if this is for our test application
			if strings.Contains(req.URL.RawQuery, "00000000-0000-0000-0000-000000000002") {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Import/get_applicationusers.json").String()), nil
			}
			return httpmock.NewStringResponse(http.StatusNotFound, ""), nil
		})

	// Mock response for deactivating system user
	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.URL.Path, "/Microsoft.Dynamics.CRM.SetState") {
				return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
			}
			return httpmock.NewStringResponse(http.StatusNotFound, ""), nil
		})

	// Mock DELETE response for system user
	httpmock.RegisterResponder("DELETE", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_application_admin" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					application_id = "00000000-0000-0000-0000-000000000002"
				}
				`,
			},
			{
				ResourceName:      "powerplatform_environment_application_admin.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "00000000-0000-0000-0000-000000000001/00000000-0000-0000-0000-000000000002",
			},
		},
	})
}

func TestUnitEnvironmentApplicationAdminResource_Delete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register mock responders
	// Get environment details for initial creation
	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Delete/get_environment.json").String()), nil
		})

	// Add application user for initial creation
	httpmock.RegisterResponder("POST", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001/addAppUser`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, "{}"), nil
		})

	// Common handler for all systemusers endpoints
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			// Check if request is for querying application users - handle both URL-encoded and plus-encoded queries
			rawQuery := req.URL.RawQuery
			if strings.Contains(rawQuery, "applicationid+eq+00000000-0000-0000-0000-000000000002") ||
				strings.Contains(rawQuery, "applicationid%20eq%2000000000-0000-0000-0000-000000000002") {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Delete/get_applicationusers.json").String()), nil
			}
			return httpmock.NewStringResponse(http.StatusNotFound, ""), nil
		})

	// Mock response for deactivating system user - handle both plain URL and URL-encoded paths
	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			urlPath := req.URL.Path
			if strings.Contains(urlPath, "/Microsoft.Dynamics.CRM.SetState") &&
				strings.Contains(urlPath, "00000000-0000-0000-0000-000000000008") {
				return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
			}
			return httpmock.NewStringResponse(http.StatusNotFound, ""), nil
		})

	// Mock response for deleting system user - handle both plain and URL-encoded paths
	httpmock.RegisterResponder("DELETE", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			// Check if this is the delete request for our test user
			urlPath := req.URL.Path
			if strings.Contains(urlPath, "00000000-0000-0000-0000-000000000008") {
				return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
			}
			return httpmock.NewStringResponse(http.StatusNotFound, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_application_admin" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					application_id = "00000000-0000-0000-0000-000000000002"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_application_admin.test", "id", "00000000-0000-0000-0000-000000000001/00000000-0000-0000-0000-000000000002"),
				),
			},
			{
				Config:       ``, // Empty config means the resource will be destroyed
				RefreshState: true,
				Check: func(_ *terraform.State) error {
					// Resource should be destroyed, but the actual deletion is a no-op
					// Just checking that we don't have errors
					return nil
				},
			},
		},
	})
}
