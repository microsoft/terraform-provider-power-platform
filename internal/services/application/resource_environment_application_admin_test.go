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

	// Check if application user exists with v9.0 API for the new method
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.0/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			// Check if this is for our test system user ID lookup
			if strings.Contains(req.URL.RawQuery, "systemuserid+eq") {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_applicationusers.json").String()), nil
			}
			return httpmock.NewStringResponse(http.StatusNotFound, ""), nil
		})

	// Mock response for deactivating system user (now using PATCH)
	httpmock.RegisterResponder("PATCH", `=~^https://test-env.crm.dynamics.com/api/data/v9.0/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

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

	// Check if application user exists with v9.0 API for the new method
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.0/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			// Check if this is for our test system user ID lookup
			if strings.Contains(req.URL.RawQuery, "systemuserid+eq") {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Import/get_applicationusers.json").String()), nil
			}
			return httpmock.NewStringResponse(http.StatusNotFound, ""), nil
		})

	// Mock response for deactivating system user (now using PATCH)
	httpmock.RegisterResponder("PATCH", `=~^https://test-env.crm.dynamics.com/api/data/v9.0/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	// Mock DELETE response for system user (updated to v9.2)
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

	// Check if application user exists with v9.0 API for the new method
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.0/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			// Check if this is for our test system user ID lookup
			if strings.Contains(req.URL.RawQuery, "systemuserid+eq") {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Delete/get_applicationusers.json").String()), nil
			}
			return httpmock.NewStringResponse(http.StatusNotFound, ""), nil
		})

	// Mock response for deactivating system user (now using PATCH)
	httpmock.RegisterResponder("PATCH", `=~^https://test-env.crm.dynamics.com/api/data/v9.0/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			// Check if this is the deactivate request for our test user
			urlPath := req.URL.Path
			if strings.Contains(urlPath, "00000000-0000-0000-0000-000000000008") {
				return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
			}
			return httpmock.NewStringResponse(http.StatusNotFound, ""), nil
		})

	// Mock response for deleting system user (updated to v9.2)
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
					// Resource should be destroyed - application user deactivated and deleted
					// Just checking that we don't have errors
					return nil
				},
			},
		},
	})
}

func TestUnitEnvironmentApplicationAdminResource_Read_ParentDeleted(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var getEnvironmentCallCount = 0
	// The application client's getEnvironment() (exact URL with api-version only).
	// Step 1 (Create + Reads) makes 5 calls that succeed; call 6 (step 2 RefreshState)
	// returns 404 to simulate the environment being deleted out-of-band.
	// Calls 7+ succeed again to allow orderly post-test cleanup.
	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			getEnvironmentCallCount++
			if getEnvironmentCallCount == 6 {
				return httpmock.NewStringResponse(http.StatusNotFound, httpmock.File("tests/resource/application_admin/Read_ParentDeleted/get_environment.json").String()), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})

	// The EnvironmentClient's GetEnvironment() returns 404 (URL includes $expand query params).
	// This is only called during tier-2 parent check when the application user read fails.
	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001\?.+expand`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNotFound, httpmock.File("tests/resource/application_admin/Read_ParentDeleted/get_environment.json").String()), nil
		})

	// Add application user
	httpmock.RegisterResponder("POST", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001/addAppUser`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, "{}"), nil
		})

	// Check if application user exists
	httpmock.RegisterResponder("GET", `=~^https://test-env\.crm\.dynamics\.com/api/data/v9\.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.URL.RawQuery, "00000000-0000-0000-0000-000000000002") {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_applicationusers.json").String()), nil
			}
			return httpmock.NewStringResponse(http.StatusNotFound, ""), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env\.crm\.dynamics\.com/api/data/v9\.0/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.URL.RawQuery, "systemuserid+eq") {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_applicationusers.json").String()), nil
			}
			return httpmock.NewStringResponse(http.StatusNotFound, ""), nil
		})

	httpmock.RegisterResponder("PATCH", `=~^https://test-env\.crm\.dynamics\.com/api/data/v9\.0/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://test-env\.crm\.dynamics\.com/api/data/v9\.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: create the application admin successfully.
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
				// Step 2: refresh when the parent environment has been deleted out-of-band.
				// The ApplicationUserExists call fails because getEnvironment returns 404;
				// the provider detects the parent is gone and removes the resource from state.
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
