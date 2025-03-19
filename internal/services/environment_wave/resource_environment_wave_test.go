// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_wave_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitEnvironmentWaveResource_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	// Register enable endpoint
	httpmock.RegisterResponder("POST", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features/October2024Update/enable$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	// Register mock for first GET call - returns Upgrading state
	httpmock.RegisterResponder("GET", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, map[string]interface{}{
				"values": []map[string]interface{}{
					{
						"FeatureName":      "October2024Update",
						"DisplayName":      "2024 release wave 2",
						"CanBeReset":       false,
						"Enabled":          true,
						"IsAllowed":        true,
						"NotBefore":        "2024-06-30T00:00:00+00:00",
						"NotAfter":         "2030-01-01T00:00:00+00:00",
						"MinVersion":       "9.0",
						"MaxVersion":       "9.3",
						"State":            "Upgrading",
						"AppsUpgradeState": "Upgrading",
					},
				},
			})
		})

	// Register mock for subsequent GET calls - returns ON state
	httpmock.RegisterResponder("GET", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, map[string]interface{}{
				"values": []map[string]interface{}{
					{
						"FeatureName":      "October2024Update",
						"DisplayName":      "2024 release wave 2",
						"CanBeReset":       false,
						"Enabled":          true,
						"IsAllowed":        true,
						"NotBefore":        "2024-06-30T00:00:00+00:00",
						"NotAfter":         "2030-01-01T00:00:00+00:00",
						"MinVersion":       "9.0",
						"MaxVersion":       "9.3",
						"State":            "ON",
						"AppsUpgradeState": "ON",
					},
				},
			})
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

	// Register enable endpoint
	httpmock.RegisterResponder("POST", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features/October2024Update/enable$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	// Register mock for GET calls - returns Failed state
	httpmock.RegisterResponder("GET", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, map[string]interface{}{
				"values": []map[string]interface{}{
					{
						"FeatureName":      "October2024Update",
						"DisplayName":      "2024 release wave 2",
						"CanBeReset":       false,
						"Enabled":          true,
						"IsAllowed":        true,
						"NotBefore":        "2024-06-30T00:00:00+00:00",
						"NotAfter":         "2030-01-01T00:00:00+00:00",
						"MinVersion":       "9.0",
						"MaxVersion":       "9.3",
						"State":            "Failed",
						"AppsUpgradeState": "Failed",
					},
				},
			})
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

	// Register enable endpoint
	httpmock.RegisterResponder("POST", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features/October2024Update/enable$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	// Register mock for first GET call - returns Upgrading state
	httpmock.RegisterResponder("GET", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, map[string]interface{}{
				"values": []map[string]interface{}{
					{
						"FeatureName":      "October2024Update",
						"DisplayName":      "2024 release wave 2",
						"CanBeReset":       false,
						"Enabled":          true,
						"IsAllowed":        true,
						"AppsUpgradeState": "Upgrading",
					},
				},
			})
		})

	// Register mock for subsequent GET calls - returns Failed state
	httpmock.RegisterResponder("GET", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, map[string]interface{}{
				"values": []map[string]interface{}{
					{
						"FeatureName":      "October2024Update",
						"DisplayName":      "2024 release wave 2",
						"CanBeReset":       false,
						"Enabled":          true,
						"IsAllowed":        true,
						"AppsUpgradeState": "Failed",
					},
				},
			})
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

	// Register enable endpoint
	httpmock.RegisterResponder("POST", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features/October2024Update/enable$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	// Register mock for GET calls - returns unknown state
	httpmock.RegisterResponder("GET", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, map[string]interface{}{
				"values": []map[string]interface{}{
					{
						"FeatureName":      "October2024Update",
						"DisplayName":      "2024 release wave 2",
						"CanBeReset":       false,
						"Enabled":          true,
						"IsAllowed":        true,
						"AppsUpgradeState": "Unknown",
					},
				},
			})
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

func TestUnitEnvironmentWaveResource_Import(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	// Register enable endpoint for initial resource creation
	httpmock.RegisterResponder("POST", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features/October2024Update/enable$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	// Register mock for GET calls - returns ON state for both creation and import
	httpmock.RegisterResponder("GET", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, map[string]interface{}{
				"values": []map[string]interface{}{
					{
						"FeatureName":      "October2024Update",
						"DisplayName":      "2024 release wave 2",
						"CanBeReset":       false,
						"Enabled":          true,
						"IsAllowed":        true,
						"AppsUpgradeState": "ON",
					},
				},
			})
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
					resource.TestCheckResourceAttr("powerplatform_environment_wave.test", "id", "00000000-0000-0000-0000-000000000001/October2024Update"),
				),
			},
			{
				ResourceName:      "powerplatform_environment_wave.test",
				ImportState:       true,
				ImportStateId:     "00000000-0000-0000-0000-000000000001/October2024Update",
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
			{
				ResourceName:  "powerplatform_environment_wave.test",
				ImportState:   true,
				ImportStateId: "invalid-format",
				ExpectError:   regexp.MustCompile("Invalid import ID"),
			},
		},
	})
}
