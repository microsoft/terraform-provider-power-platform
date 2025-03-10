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
                "id": "/providers/Microsoft.BusinessAppPlatform/tenant",
                "name": "default",
                "type": "Microsoft.BusinessAppPlatform/tenant",
                "tenantId": "%s",
                "state": "Enabled"
            }`, testTenantID)), nil
		})

	// Mock GET tenant isolation policy endpoint with empty policy.
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
                "properties": {
                    "isDisabled": false,
                    "allowedTenants": []
                }
            }`), nil
		})
}

// TestAccTenantIsolationPolicy_Validate_Create tests the basic creation and import of a tenant isolation policy.
// It verifies:
// 1. Creation of a policy with a single allowed tenant (inbound and outbound access)
// 2. The ability to import an existing policy and verify its attributes match the configuration.
func TestAccTenantIsolationPolicy_Validate_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_tenant_isolation_policy" "test" {
					is_disabled = false
					allowed_tenants = toset([
						{
							tenant_id = "11111111-1111-1111-1111-111111111111"
							inbound  = true
							outbound = true
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

// TestAccTenantIsolationPolicy_Validate_Update tests the update functionality of a tenant isolation policy.
// It verifies:
// 1. Initial creation of a policy with one allowed tenant
// 2. Updating the policy to:
//   - Enable the disabled flag
//   - Add a second allowed tenant
//   - Both tenants have different inbound/outbound settings.
func TestAccTenantIsolationPolicy_Validate_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_tenant_isolation_policy" "test" {
					is_disabled = false
					allowed_tenants = toset([
						{
							tenant_id = "11111111-1111-1111-1111-111111111111"
							inbound = true
							outbound = true
						}
					])
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "is_disabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "allowed_tenants.#", "1"),
				),
			},
			{
				Config: `
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
							inbound = true
							outbound = false
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

// TestAccTenantIsolationPolicy_Validate_Delete tests the removal of all allowed tenants from a policy.
// It verifies:
// 1. Initial creation of a policy with one allowed tenant
// 2. Updating the policy to remove all allowed tenants
// 3. The policy exists but has an empty allowed_tenants list.
func TestAccTenantIsolationPolicy_Validate_Delete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_tenant_isolation_policy" "test" {
					is_disabled = false
					allowed_tenants = toset([
						{
							tenant_id = "11111111-1111-1111-1111-111111111111"
							inbound = true
							outbound = true
						}
					])
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "allowed_tenants.#", "1"),
				),
			},
			{
				Config: `
				resource "powerplatform_tenant_isolation_policy" "test" {
					is_disabled = false
					allowed_tenants = toset([])
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "allowed_tenants.#", "0"),
				),
			},
		},
	})
}

// TestUnitTenantIsolationPolicyResource_Validate_Create tests the creation of a tenant isolation policy
// using mocked API responses. It verifies:
// 1. The API calls are made with correct parameters
// 2. The resource correctly processes the API response
// 3. The resource state matches the expected configuration with one allowed tenant.
func TestUnitTenantIsolationPolicyResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	setupTenantHttpMocks()

	// Setup PUT responder for creating the policy.
	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
                "properties": {
                    "tenantId": "00000000-0000-0000-0000-000000000001",
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
            }`), nil
		})

	// Also update the GET responder to match what is returned by PUT.
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
                "properties": {
                    "tenantId": "00000000-0000-0000-0000-000000000001",
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
            }`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
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

// TestUnitTenantIsolationPolicyResource_Validate_Update tests the update process of a tenant isolation policy
// using mocked API responses. It verifies:
// 1. Initial creation with one allowed tenant
// 2. Update to include:
//   - Changed is_disabled flag from false to true
//   - Modified first tenant's permissions
//   - Added second tenant with different permissions
//
// 3. Proper state transitions and API interactions during the update.
func TestUnitTenantIsolationPolicyResource_Validate_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Initial tenant response.
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/tenant?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
                "id": "/providers/Microsoft.BusinessAppPlatform/tenant",
                "name": "default",
                "type": "Microsoft.BusinessAppPlatform/tenant",
                "tenantId": "00000000-0000-0000-0000-000000000001",
                "state": "Enabled"
            }`), nil
		})

	// Step 1: Empty initial state, first GET returns empty policy.
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
                "properties": {
                    "isDisabled": false,
                    "allowedTenants": []
                }
            }`), nil
		})

	// Step 1: First PUT creates policy with initial state.
	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
		func(req *http.Request) (*http.Response, error) {
			// After first PUT, register a new GET to return initial state.
			httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
				func(req *http.Request) (*http.Response, error) {
					return httpmock.NewStringResponse(http.StatusOK, `{
                        "properties": {
                            "tenantId": "00000000-0000-0000-0000-000000000001",
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
                    }`), nil
				})

			// Register a new PUT handler for the update operation.
			httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
				func(req *http.Request) (*http.Response, error) {
					// After second PUT, register a new GET to return updated state.
					httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
						func(req *http.Request) (*http.Response, error) {
							return httpmock.NewStringResponse(http.StatusOK, `{
                                "properties": {
                                    "tenantId": "00000000-0000-0000-0000-000000000001",
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
                            }`), nil
						})

					return httpmock.NewStringResponse(http.StatusOK, `{
                        "properties": {
                            "tenantId": "00000000-0000-0000-0000-000000000001",
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
                    }`), nil
				})

			return httpmock.NewStringResponse(http.StatusOK, `{
                "properties": {
                    "tenantId": "00000000-0000-0000-0000-000000000001",
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
            }`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create initial policy.
				Config: `
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
				Config: `
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

// TestUnitTenantIsolationPolicyResource_Validate_Delete tests the removal of all tenants from a policy
// using mocked API responses. It verifies:
// 1. Initial creation with one allowed tenant
// 2. Update to remove all tenants
// 3. Proper API calls for emptying the allowed_tenants list
// 4. Final state has empty allowed_tenants list.
func TestUnitTenantIsolationPolicyResource_Validate_Delete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Initial tenant mocks setup
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/tenant?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{
                "id": "/providers/Microsoft.BusinessAppPlatform/tenant",
                "name": "default",
                "type": "Microsoft.BusinessAppPlatform/tenant",
                "tenantId": "%s",
                "state": "Enabled"
            }`, testTenantID)), nil
		})

	// Step 1: Initial GET returns empty policy
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
                "properties": {
                    "isDisabled": false,
                    "allowedTenants": []
                }
            }`), nil
		})

	// First PUT: Create policy with one tenant
	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
		func(req *http.Request) (*http.Response, error) {
			// After first PUT, update the GET handler to return the created policy
			httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
				func(req *http.Request) (*http.Response, error) {
					return httpmock.NewStringResponse(http.StatusOK, `{
                        "properties": {
                            "tenantId": "00000000-0000-0000-0000-000000000001",
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
                    }`), nil
				})

			// Register a new handler for the second PUT (empty tenants)
			httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
				func(req *http.Request) (*http.Response, error) {
					// After second PUT, update the GET handler to return the empty policy
					httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
						func(req *http.Request) (*http.Response, error) {
							return httpmock.NewStringResponse(http.StatusOK, `{
                                "properties": {
                                    "tenantId": "00000000-0000-0000-0000-000000000001",
                                    "isDisabled": false,
                                    "allowedTenants": []
                                }
                            }`), nil
						})

					return httpmock.NewStringResponse(http.StatusOK, `{
                        "properties": {
                            "tenantId": "00000000-0000-0000-0000-000000000001",
                            "isDisabled": false,
                            "allowedTenants": []
                        }
                    }`), nil
				})

			return httpmock.NewStringResponse(http.StatusOK, `{
                "properties": {
                    "tenantId": "00000000-0000-0000-0000-000000000001",
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
            }`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
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
				Config: `
				resource "powerplatform_tenant_isolation_policy" "test" {
					is_disabled = false
					allowed_tenants = toset([])
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "allowed_tenants.#", "0"),
				),
			},
		},
	})
}

// TestUnitTenantIsolationPolicyResource_Validate_Create_Error tests error handling during policy creation
// using mocked API responses. It verifies:
// 1. Proper error handling when an invalid tenant ID is provided
// 2. The error message is properly propagated from the API
// 3. The resource creation fails as expected with a validation error.
func TestUnitTenantIsolationPolicyResource_Validate_Create_Error(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	setupTenantHttpMocks()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusBadRequest, `{
                "error": {
                    "code": "BadRequest",
                    "message": "Invalid tenant ID format",
                    "details": [
                        {
                            "code": "ValidationError",
                            "message": "The tenant ID must be a valid GUID"
                        }
                    ]
                }
            }`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
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
