// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package disaster_recovery_test

import (
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitDisasterRecoveryResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("PATCH", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000099?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000099?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_environment_enabled.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_disaster_recovery" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					enabled        = true
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_disaster_recovery.test", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment_disaster_recovery.test", "environment_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment_disaster_recovery.test", "enabled", "true"),
				),
			},
		},
	})
}

func TestUnitDisasterRecoveryResource_Validate_Create_And_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	getEnvironmentResponseInx := 0

	httpmock.RegisterResponder("PATCH", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000099?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000099?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_lifecycle.json").String()), nil
		})

	environmentFiles := []string{
		"tests/resource/Validate_Create_And_Update/get_environment_enabled.json",  // Step 1: Create's post-create Read
		"tests/resource/Validate_Create_And_Update/get_environment_enabled.json",  // Step 1: drift-check Read
		"tests/resource/Validate_Create_And_Update/get_environment_enabled.json",  // Step 2: refresh Read (must still be enabled so TF detects diff)
		"tests/resource/Validate_Create_And_Update/get_environment_disabled.json", // Step 2: Update's post-update Read
		"tests/resource/Validate_Create_And_Update/get_environment_disabled.json", // Step 2: drift-check Read
		"tests/resource/Validate_Create_And_Update/get_environment_disabled.json", // Destroy fallback
	}

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			filePath := environmentFiles[getEnvironmentResponseInx]
			if getEnvironmentResponseInx < len(environmentFiles)-1 {
				getEnvironmentResponseInx++
			}
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(filePath).String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with DR enabled
				Config: `
				resource "powerplatform_environment_disaster_recovery" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					enabled        = true
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_disaster_recovery.test", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment_disaster_recovery.test", "enabled", "true"),
				),
			},
			{
				// Step 2: Update to disable DR
				Config: `
				resource "powerplatform_environment_disaster_recovery" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					enabled        = false
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_disaster_recovery.test", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment_disaster_recovery.test", "enabled", "false"),
				),
			},
		},
	})
}

func TestUnitDisasterRecoveryResource_Validate_Force_Recreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("PATCH", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/.*\?api-version=2021-04-01$`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000099?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000099?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_environment_enabled.json").String()), nil
		})

	env2Response := httpmock.File("tests/resource/Validate_Create_And_Update/get_environment_enabled.json").String()
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000002?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, env2Response), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_disaster_recovery" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					enabled        = true
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_disaster_recovery.test", "environment_id", "00000000-0000-0000-0000-000000000001"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment_disaster_recovery" "test" {
					environment_id = "00000000-0000-0000-0000-000000000002"
					enabled        = true
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_disaster_recovery.test", "environment_id", "00000000-0000-0000-0000-000000000002"),
				),
			},
		},
	})
}

func TestUnitDisasterRecoveryResource_Validate_Import(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("PATCH", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000099?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000099?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_environment_enabled.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_disaster_recovery" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					enabled        = true
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_disaster_recovery.test", "id", "00000000-0000-0000-0000-000000000001"),
				),
			},
			{
				ResourceName:      "powerplatform_environment_disaster_recovery.test",
				ImportState:       true,
				ImportStateId:     "00000000-0000-0000-0000-000000000001",
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitDisasterRecoveryResource_Validate_Create_Error(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("PATCH", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusForbidden, `{"error":{"code":"Forbidden","message":"Insufficient permissions to enable disaster recovery"}}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_disaster_recovery" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					enabled        = true
				}`,
				ExpectError: regexp.MustCompile(".*error.*"),
			},
		},
	})
}

func TestUnitDisasterRecoveryResource_Validate_Default_Enabled(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("PATCH", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000099?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000099?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_environment_enabled.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create without specifying enabled - should default to true
				Config: `
				resource "powerplatform_environment_disaster_recovery" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_disaster_recovery.test", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment_disaster_recovery.test", "enabled", "true"),
				),
			},
		},
	})
}

func TestUnitDisasterRecoveryResource_Validate_Create_Disabled(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("PATCH", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000099?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000099?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_environment_disabled.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_disaster_recovery" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					enabled        = false
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_disaster_recovery.test", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment_disaster_recovery.test", "enabled", "false"),
				),
			},
		},
	})
}

func TestUnitDisasterRecoveryResource_Validate_Update_Enable(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("PATCH", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000099?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000099?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_lifecycle.json").String()), nil
		})

	getEnvironmentResponseInx := 0
	environmentFiles := []string{
		"tests/resource/Validate_Create_And_Update/get_environment_disabled.json", // Step 1: Create's post-create Read (disabled)
		"tests/resource/Validate_Create_And_Update/get_environment_disabled.json", // Step 1: drift-check Read (disabled)
		"tests/resource/Validate_Create_And_Update/get_environment_disabled.json", // Step 2: refresh Read (disabled — so TF sees diff to enabled=true)
		"tests/resource/Validate_Create_And_Update/get_environment_enabled.json",  // Step 2: Update's post-update Read (enabled)
		"tests/resource/Validate_Create_And_Update/get_environment_enabled.json",  // Step 2: drift-check Read (enabled)
		"tests/resource/Validate_Create_And_Update/get_environment_enabled.json",  // Destroy fallback
	}

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			filePath := environmentFiles[getEnvironmentResponseInx]
			if getEnvironmentResponseInx < len(environmentFiles)-1 {
				getEnvironmentResponseInx++
			}
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(filePath).String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with DR disabled
				Config: `
				resource "powerplatform_environment_disaster_recovery" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					enabled        = false
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_disaster_recovery.test", "enabled", "false"),
				),
			},
			{
				// Step 2: Update to enable DR
				Config: `
				resource "powerplatform_environment_disaster_recovery" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					enabled        = true
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_disaster_recovery.test", "enabled", "true"),
				),
			},
		},
	})
}

func TestAccDisasterRecoveryResource_Validate_Create(t *testing.T) {
	t.Skip("Skipping live test for Create validation until we can reliably set up and tear down disaster recovery")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azapi": {
				VersionConstraint: constants.AZAPI_PROVIDER_VERSION_CONSTRAINT,
				Source:            "azure/azapi",
			},
			"time": {
				Source: "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
				data "azapi_client_config" "current" {}

				resource "azapi_resource" "rg_dr_test" {
					type     = "Microsoft.Resources/resourceGroups@2021-04-01"
					location = "East US"
					name     = "power-platform-billing-` + mocks.TestName() + strconv.Itoa(rand.Intn(9999)) + `"
				}

				resource "powerplatform_billing_policy" "dr_test" {
					name     = "` + strings.ReplaceAll(mocks.TestName(), "_", "") + strconv.Itoa(rand.Intn(9999)) + `"
					location = "unitedstates"
					status   = "Enabled"
					billing_instrument = {
						resource_group  = "powerplatform_billing"
    					subscription_id = "2bc1f261-7e26-490c-9fd5-b7ca72032ad3"
					}
				}

				resource "powerplatform_environment" "dr_test" {
					display_name      = "` + mocks.TestName() + `"
					location          = "unitedstates"
					environment_type  = "Production"
					billing_policy_id = powerplatform_billing_policy.dr_test.id
					dataverse = {
						language_code     = "1033"
						currency_code     = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "time_sleep" "wait_for_dataverse" {
					create_duration = "120s"

					depends_on = [powerplatform_environment.dr_test]
				}

				resource "powerplatform_managed_environment" "dr_test" {
					environment_id             = powerplatform_environment.dr_test.id
					is_usage_insights_disabled = true
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "None"
					suppress_validation_emails = true

					depends_on = [time_sleep.wait_for_dataverse]
				}

				resource "powerplatform_environment_disaster_recovery" "dr_test" {
					environment_id = powerplatform_environment.dr_test.id
					enabled        = true

					timeouts = {
						create = "30m"
						update = "30m"
						delete = "30m"
						read   = "30m"
					}

					depends_on = [powerplatform_managed_environment.dr_test]
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.dr_test", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment_disaster_recovery.dr_test", "enabled", "true"),
				),
			},
		},
	})
}

func TestAccDisasterRecoveryResource_Validate_Update(t *testing.T) {
	t.Skip("Skipping live test for Create validation until we can reliably set up and tear down disaster recovery")
	envName := mocks.TestName()
	rgName := "power-platform-billing-" + mocks.TestName() + strconv.Itoa(rand.Intn(9999))
	bpName := strings.ReplaceAll(mocks.TestName(), "_", "") + strconv.Itoa(rand.Intn(9999))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azapi": {
				VersionConstraint: constants.AZAPI_PROVIDER_VERSION_CONSTRAINT,
				Source:            "azure/azapi",
			},
			"time": {
				Source: "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
				data "azapi_client_config" "current" {}

				resource "azapi_resource" "rg_dr_test" {
					type     = "Microsoft.Resources/resourceGroups@2021-04-01"
					location = "East US"
					name     = "` + rgName + `"
				}

				resource "powerplatform_billing_policy" "dr_test" {
					name     = "` + bpName + `"
					location = "unitedstates"
					status   = "Enabled"
					billing_instrument = {
						resource_group  = azapi_resource.rg_dr_test.name
						subscription_id = data.azapi_client_config.current.subscription_id
					}
				}

				resource "powerplatform_environment" "dr_test" {
					display_name      = "` + envName + `"
					location          = "unitedstates"
					environment_type  = "Production"
					billing_policy_id = powerplatform_billing_policy.dr_test.id
					dataverse = {
						language_code     = "1033"
						currency_code     = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "time_sleep" "wait_for_dataverse" {
					create_duration = "120s"

					depends_on = [powerplatform_environment.dr_test]
				}

				resource "powerplatform_managed_environment" "dr_test" {
					environment_id             = powerplatform_environment.dr_test.id
					is_usage_insights_disabled = true
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "None"
					suppress_validation_emails = true

					depends_on = [time_sleep.wait_for_dataverse]
				}

				resource "powerplatform_environment_disaster_recovery" "dr_test" {
					environment_id = powerplatform_environment.dr_test.id
					enabled        = true

					timeouts = {
						create = "30m"
						update = "30m"
						delete = "30m"
						read   = "30m"
					}

					depends_on = [powerplatform_managed_environment.dr_test]
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_disaster_recovery.dr_test", "enabled", "true"),
				),
			},
			{
				Config: `
				data "azapi_client_config" "current" {}

				resource "azapi_resource" "rg_dr_test" {
					type     = "Microsoft.Resources/resourceGroups@2021-04-01"
					location = "East US"
					name     = "` + rgName + `"
				}

				resource "powerplatform_billing_policy" "dr_test" {
					name     = "` + bpName + `"
					location = "unitedstates"
					status   = "Enabled"
					billing_instrument = {
						resource_group  = azapi_resource.rg_dr_test.name
						subscription_id = data.azapi_client_config.current.subscription_id
					}
				}

				resource "powerplatform_environment" "dr_test" {
					display_name      = "` + envName + `"
					location          = "unitedstates"
					environment_type  = "Production"
					billing_policy_id = powerplatform_billing_policy.dr_test.id
					dataverse = {
						language_code     = "1033"
						currency_code     = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "time_sleep" "wait_for_dataverse" {
					create_duration = "120s"

					depends_on = [powerplatform_environment.dr_test]
				}

				resource "powerplatform_managed_environment" "dr_test" {
					environment_id             = powerplatform_environment.dr_test.id
					is_usage_insights_disabled = true
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "None"
					suppress_validation_emails = true

					depends_on = [time_sleep.wait_for_dataverse]
				}

				resource "powerplatform_environment_disaster_recovery" "dr_test" {
					environment_id = powerplatform_environment.dr_test.id
					enabled        = false

					timeouts = {
						create = "30m"
						update = "30m"
						delete = "30m"
						read   = "30m"
					}

					depends_on = [powerplatform_managed_environment.dr_test]
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_disaster_recovery.dr_test", "enabled", "false"),
				),
			},
		},
	})
}
