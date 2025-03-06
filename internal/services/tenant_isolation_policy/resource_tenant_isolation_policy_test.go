// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant_isolation_policy_test

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

const testTenantID = "00000000-0000-0000-0000-000000000001"

func setupTenantHttpMocks() {
	// Mock tenant endpoint that's called before CRUD operations.
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/tenant?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{
				"tenantId": "%s",
				"state": "Enabled",
				"location": "unitedstates",
				"aadCountryGeo": "unitedstates",
				"dataStorageGeo": "unitedstates",
				"defaultEnvironmentGeo": "unitedstates",
				"aadDataBoundary": "none",
				"fedRAMPHighCertificationRequired": false
			}`, testTenantID)), nil
		})

	// Mock GET tenant isolation policy endpoint with empty policy.
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{
				"properties": {
					"tenantId": "%s",
					"isDisabled": false,
					"allowedTenants": []
				}
			}`, testTenantID)), nil
		})
}

func getProviderConfig() string {
	return fmt.Sprintf(`
provider "powerplatform" {
	tenant_id = "%s"
	use_cli = false
	client_id = "test-client-id"
	client_secret = "test-client-secret"
}
`, testTenantID)
}

// Test functions with updated API paths...
func TestAccTenantIsolationPolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: _testAccTenantIsolationPolicy_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "is_disabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "allowed_tenants.#", "1"),
				),
			},
			{
				ResourceName:      "powerplatform_tenant_isolation_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTenantIsolationPolicy_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: _testAccTenantIsolationPolicy_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "is_disabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "allowed_tenants.#", "1"),
				),
			},
			{
				Config: _testAccTenantIsolationPolicy_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "is_disabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "allowed_tenants.#", "2"),
				),
			},
		},
	})
}

func TestAccTenantIsolationPolicy_remove(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: _testAccTenantIsolationPolicy_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "allowed_tenants.#", "1"),
				),
			},
			{
				Config: _testAccTenantIsolationPolicy_empty(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "allowed_tenants.#", "0"),
				),
			},
		},
	})
}

func TestUnitTenantIsolationPolicyResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	setupTenantHttpMocks()

	// Initial state for first response.
	firstResponseJson := fmt.Sprintf(`{
		"properties": {
			"tenantId": "%s",
			"isDisabled": false,
			"allowedTenants": [
				{
					"tenantId": "11111111-1111-1111-1111-111111111111",
					"direction": {
						"inbound": true,
						"outbound": false
					}
				}
			]
		}
	}`, testTenantID)

	// Setup PUT responder for creating the policy.
	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, firstResponseJson), nil
		})

	// Also update the GET responder to match what is returned by PUT.
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, firstResponseJson), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfig() + `
				resource "powerplatform_tenant_isolation_policy" "test" {
					is_disabled = false
					allowed_tenants = toset([
						{
							tenant_id = "11111111-1111-1111-1111-111111111111"
							inbound = true
							outbound = false
						}
					])
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "is_disabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "allowed_tenants.#", "1"),
				),
			},
		},
	})
}

func TestUnitTenantIsolationPolicyResource_Validate_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Initial tenant response.
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/tenant?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{
				"tenantId": "%s",
				"state": "Enabled"
			}`, testTenantID)), nil
		})

	// Initial state for the first test step.
	initialState := fmt.Sprintf(`{
		"properties": {
			"tenantId": "%s",
			"isDisabled": false,
			"allowedTenants": [
				{
					"tenantId": "11111111-1111-1111-1111-111111111111",
					"direction": {
						"inbound": true,
						"outbound": false
					}
				}
			]
		}
	}`, testTenantID)

	// Updated state for the second test step.
	updatedState := fmt.Sprintf(`{
		"properties": {
			"tenantId": "%s",
			"isDisabled": true,
			"allowedTenants": [
				{
					"tenantId": "11111111-1111-1111-1111-111111111111", 
					"direction": {
						"inbound": true,
						"outbound": true
					}
				},
				{
					"tenantId": "22222222-2222-2222-2222-222222222222",
					"direction": {
						"inbound": false,
						"outbound": true
					}
				}
			]
		}
	}`, testTenantID)

	// Step 1: Empty initial state, first GET returns empty policy.
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{
				"properties": {
					"tenantId": "%s",
					"isDisabled": false,
					"allowedTenants": []
				}
			}`, testTenantID)), nil
		})

	// Step 1: First PUT creates policy with initial state.
	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
		func(req *http.Request) (*http.Response, error) {
			// After first PUT, register a new GET to return initial state.
			httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
				func(req *http.Request) (*http.Response, error) {
					return httpmock.NewStringResponse(http.StatusOK, initialState), nil
				})

			// Register a new PUT handler for the update operation.
			httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
				func(req *http.Request) (*http.Response, error) {
					// After second PUT, register a new GET to return updated state.
					httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
						func(req *http.Request) (*http.Response, error) {
							return httpmock.NewStringResponse(http.StatusOK, updatedState), nil
						})

					return httpmock.NewStringResponse(http.StatusOK, updatedState), nil
				})

			return httpmock.NewStringResponse(http.StatusOK, initialState), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create initial policy.
				Config: getProviderConfig() + `
				resource "powerplatform_tenant_isolation_policy" "test" {
					is_disabled = false
					allowed_tenants = toset([
						{
							tenant_id = "11111111-1111-1111-1111-111111111111"
							inbound = true
							outbound = false
						}
					])
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "is_disabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "allowed_tenants.#", "1"),
				),
			},
			{
				// Step 2: Update the policy.
				Config: getProviderConfig() + `
				resource "powerplatform_tenant_isolation_policy" "test" {
					is_disabled = true
					allowed_tenants = toset([
						{
							tenant_id = "11111111-1111-1111-1111-111111111111"
							inbound = true
							outbound = true
						},
						{
							tenant_id = "22222222-2222-2222-2222-222222222222"
							inbound = false
							outbound = true
						}
					])
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "is_disabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "allowed_tenants.#", "2"),
				),
			},
		},
	})
}

func TestUnitTenantIsolationPolicyResource_Validate_Delete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	setupTenantHttpMocks()

	// Create a consistent response to use in both PUT and GET.
	policyJson := fmt.Sprintf(`{
		"properties": {
			"tenantId": "%s",
			"isDisabled": false,
			"allowedTenants": [
				{
					"tenantId": "11111111-1111-1111-1111-111111111111",
					"direction": {
						"inbound": true,
						"outbound": false
					}
				}
			]
		}
	}`, testTenantID)

	// Register PUT responder for initial creation.
	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, policyJson), nil
		})

	// Register GET responder to return the created policy.
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, policyJson), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfig() + `
				resource "powerplatform_tenant_isolation_policy" "test" {
					is_disabled = false
					allowed_tenants = toset([
						{
							tenant_id = "11111111-1111-1111-1111-111111111111"
							inbound = true
							outbound = false
						}
					])
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "is_disabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "allowed_tenants.#", "1"),
				),
			},
			{
				ResourceName:      "powerplatform_tenant_isolation_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitTenantIsolationPolicyResource_Validate_Create_Error(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	setupTenantHttpMocks()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusBadRequest, `{
				"error": {
					"code": "BadRequest",
					"message": "Invalid tenant ID format"
				}
			}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfig() + `
				resource "powerplatform_tenant_isolation_policy" "test" {
					is_disabled = false
					allowed_tenants = toset([
						{
							tenant_id = "invalid-tenant-id"
							inbound = true 
							outbound = false
						}
					])
				}`,
				ExpectError: regexp.MustCompile("Invalid tenant ID format"),
			},
		},
	})
}

func TestUnitTenantIsolationPolicyResource_Validate_Empty_TenantId(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfig() + `
                resource "powerplatform_tenant_isolation_policy" "test" {
                    is_disabled = false
                    allowed_tenants = toset([
                        {
                            tenant_id = ""  
                            inbound = true
                            outbound = false
                        }
                    ])
                }`,
				ExpectError: regexp.MustCompile("string length must be at least 1"),
			},
		},
	})
}

func TestUnitTenantIsolationPolicyResource_Validate_Missing_IsDisabled(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfig() + `
                resource "powerplatform_tenant_isolation_policy" "test" {
                    allowed_tenants = toset([
                        {
                            tenant_id = "11111111-1111-1111-1111-111111111111"
                            inbound = true
                            outbound = false
                        }
                    ])
                }`,
				ExpectError: regexp.MustCompile("The argument \"is_disabled\" is required"),
			},
		},
	})
}

func TestUnitTenantIsolationPolicyValidation_Invalid_TenantId(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	setupTenantHttpMocks()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusBadRequest, `{
				"error": {
					"code": "BadRequest",
					"message": "Invalid tenant ID format"
				}
			}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfig() + `
                resource "powerplatform_tenant_isolation_policy" "test" {
                    is_disabled = false
                    allowed_tenants = toset([
                        {
                            tenant_id = "invalid-guid-format"
                            inbound = true
                            outbound = false
                        }
                    ])
                }`,
				ExpectError: regexp.MustCompile("Invalid tenant ID format"),
			},
		},
	})
}

// func TestUnitTenantIsolationPolicyResource_Validate_Import(t *testing.T) {
// 	httpmock.Activate()
// 	defer httpmock.DeactivateAndReset()

// 	setupTenantHttpMocks()

// 	// Register GET responder with the expected state
// 	policyJson := fmt.Sprintf(`{
// 		"properties": {
// 			"tenantId": "%s",
// 			"isDisabled": false,
// 			"allowedTenants": [
// 				{
// 					"tenantId": "11111111-1111-1111-1111-111111111111",
// 					"direction": {
// 						"inbound": true,
// 						"outbound": false
// 					}
// 				}
// 			]
// 		}
// 	}`, testTenantID)

// 	httpmock.RegisterResponder("GET",
// 		fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
// 		func(req *http.Request) (*http.Response, error) {
// 			return httpmock.NewStringResponse(http.StatusOK, policyJson), nil
// 		})

// 	resource.Test(t, resource.TestCase{
// 		IsUnitTest:               true,
// 		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: getProviderConfig() + `
// 				import {
// 					to: powerplatform_tenant_isolation_policy.test
// 					id: "00000000-0000-0000-0000-000000000001"
// 				}

// 				resource "powerplatform_tenant_isolation_policy" "test" {
// 					is_disabled = false
// 					allowed_tenants = toset([
// 						{
// 							tenant_id = "11111111-1111-1111-1111-111111111111"
// 							inbound = true
// 							outbound = false
// 						}
// 					])
// 				}`,
// 				ResourceName:      "powerplatform_tenant_isolation_policy.test",
// 				ImportState:       true,
// 				ImportStateId:     testTenantID,
// 				ImportStateVerify: true,
// 			},
// 		},
// 	})
// }

// Helper functions for acceptance tests.
func _testAccTenantIsolationPolicy_basic() string {
	return getProviderConfig() + `
resource "powerplatform_tenant_isolation_policy" "test" {
  is_disabled = false
  allowed_tenants = toset([
    {
      tenant_id = "11111111-1111-1111-1111-111111111111"
      inbound  = true
      outbound = true
    }
  ])
}
`
}

func _testAccTenantIsolationPolicy_update() string {
	return getProviderConfig() + `
resource "powerplatform_tenant_isolation_policy" "test" {
  is_disabled = true
  allowed_tenants = toset([
    {
      tenant_id = "11111111-1111-1111-1111-111111111111"
      inbound  = true
      outbound = true
    },
    {
      tenant_id = "22222222-2222-2222-2222-222222222222"
      inbound  = true
      outbound = false
    }
  ])
}
`
}

func _testAccTenantIsolationPolicy_empty() string {
	return getProviderConfig() + `
resource "powerplatform_tenant_isolation_policy" "test" {
  allowed_tenants = toset([])
}
`
}
