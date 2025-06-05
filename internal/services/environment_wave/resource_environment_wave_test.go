// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_wave_test

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func loadTestResponse(t *testing.T, testFolder string, filename string) string {
	path := filepath.Join("test", "resource", testFolder, filename)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test response file %s: %v", filename, err)
	}
	return string(content)
}

func registerOrganizationsMock(t *testing.T, testFolder string) {
	httpmock.RegisterResponder("GET", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/tenants/mytenant/organizations$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, loadTestResponse(t, testFolder, "get_organizations.json")), nil
		})
}

func registerEnvironmentMock(t *testing.T, testFolder string) {
	// Register a mock for the environment API call
	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001(\?.+)?$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, loadTestResponse(t, testFolder, "get_environment_00000000-0000-0000-0000-000000000001.json")), nil
		})
}

func TestUnitEnvironmentWaveResource_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()
	registerEnvironmentMock(t, "EnvironmentWaveResource_Create")

	// Register organizations mock
	registerOrganizationsMock(t, "EnvironmentWaveResource_Create")

	// Register enable endpoint
	httpmock.RegisterResponder("POST", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features/October2024Update/enable$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	// Register mock for first GET call - returns Upgrading state
	httpmock.RegisterResponder("GET", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, loadTestResponse(t, "EnvironmentWaveResource_Create", "get_features_upgrading.json")), nil
		})

	// Register mock for subsequent GET calls - returns ON state
	httpmock.RegisterResponder("GET", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, loadTestResponse(t, "EnvironmentWaveResource_Create", "get_features_enabled.json")), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_wave" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					feature_name  = "October2024Update"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_wave.test", "environment_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment_wave.test", "feature_name", "October2024Update"),
					resource.TestCheckResourceAttr("powerplatform_environment_wave.test", "state", "enabled"),
				),
			},
		},
	})
}

func TestUnitEnvironmentWaveResource_Error(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()
	registerEnvironmentMock(t, "EnvironmentWaveResource_Error")

	// Register organizations mock
	registerOrganizationsMock(t, "EnvironmentWaveResource_Error")

	// Register enable endpoint
	httpmock.RegisterResponder("POST", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features/October2024Update/enable$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	// Register mock for GET calls - returns Failed state
	httpmock.RegisterResponder("GET", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, loadTestResponse(t, "EnvironmentWaveResource_Error", "get_features_failed.json")), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_wave" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					feature_name  = "October2024Update"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_wave.test", "environment_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment_wave.test", "feature_name", "October2024Update"),
					resource.TestCheckResourceAttr("powerplatform_environment_wave.test", "state", "error"),
				),
			},
		},
	})
}

func TestUnitEnvironmentWaveResource_NotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()
	registerEnvironmentMock(t, "EnvironmentWaveResource_NotFound")

	// Register organizations mock
	registerOrganizationsMock(t, "EnvironmentWaveResource_NotFound")

	// Register enable endpoint
	httpmock.RegisterResponder("POST", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features/October2024Update/enable$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	// Register mock for GET calls - returns 404
	httpmock.RegisterResponder("GET", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(404, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_wave" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					feature_name  = "October2024Update"
				}`,
				ExpectError: regexp.MustCompile(`.*404.*`),
			},
		},
	})
}

func TestUnitEnvironmentWaveResource_FailedDuringUpgrade(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()
	registerEnvironmentMock(t, "EnvironmentWaveResource_FailedDuringUpgrade")

	// Register organizations mock
	registerOrganizationsMock(t, "EnvironmentWaveResource_FailedDuringUpgrade")

	// Register enable endpoint
	httpmock.RegisterResponder("POST", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features/October2024Update/enable$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	// Register mock for first GET call - returns Upgrading state
	httpmock.RegisterResponder("GET", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, loadTestResponse(t, "EnvironmentWaveResource_FailedDuringUpgrade", "get_features_upgrading.json")), nil
		})

	// Register mock for subsequent GET calls - returns Failed state
	httpmock.RegisterResponder("GET", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, loadTestResponse(t, "EnvironmentWaveResource_FailedDuringUpgrade", "get_features_failed.json")), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_wave" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					feature_name  = "October2024Update"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_wave.test", "environment_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment_wave.test", "feature_name", "October2024Update"),
					resource.TestCheckResourceAttr("powerplatform_environment_wave.test", "state", "error"),
				),
			},
		},
	})
}

func TestUnitEnvironmentWaveResource_UnsupportedState(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()
	registerEnvironmentMock(t, "EnvironmentWaveResource_UnsupportedState")

	// Register organizations mock
	registerOrganizationsMock(t, "EnvironmentWaveResource_UnsupportedState")

	// Register enable endpoint
	httpmock.RegisterResponder("POST", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features/October2024Update/enable$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	// Register mock for GET calls - returns unknown state
	httpmock.RegisterResponder("GET", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, loadTestResponse(t, "EnvironmentWaveResource_UnsupportedState", "get_features_unknown.json")), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_wave" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					feature_name  = "October2024Update"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_wave.test", "environment_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment_wave.test", "feature_name", "October2024Update"),
					resource.TestCheckResourceAttr("powerplatform_environment_wave.test", "state", "error"),
				),
			},
		},
	})
}

func TestUnitEnvironmentWaveResource_Timeout(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()
	registerEnvironmentMock(t, "EnvironmentWaveResource_Timeout")

	// Register organizations mock
	registerOrganizationsMock(t, "EnvironmentWaveResource_Timeout")

	// Register enable endpoint
	httpmock.RegisterResponder("POST", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features/October2024Update/enable$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	// Register mock for GET calls - always returns Upgrading state (never completes)
	httpmock.RegisterResponder("GET", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, loadTestResponse(t, "EnvironmentWaveResource_Timeout", "get_features_upgrading.json")), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_wave" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					feature_name  = "October2024Update"
					timeouts = {
						create = "1s"
					}
				}`,
				ExpectError: regexp.MustCompile(`context deadline exceeded`),
			},
		},
	})
}
