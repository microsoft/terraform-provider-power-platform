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

func setupEnvironmentHttpMocks() {
	// Mock tenant endpoint that's called before CRUD operations
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

	// Mock GET tenant isolation policy endpoint
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"isDisabled": false,
				"allowedTenants": []
			}`), nil
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
				Config: testAccTenantIsolationPolicy_basic(),
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
				Config: testAccTenantIsolationPolicy_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "is_disabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "allowed_tenants.#", "1"),
				),
			},
			{
				Config: testAccTenantIsolationPolicy_update(),
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
				Config: testAccTenantIsolationPolicy_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_tenant_isolation_policy.test", "allowed_tenants.#", "1"),
				),
			},
			{
				Config: testAccTenantIsolationPolicy_empty(),
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

	setupEnvironmentHttpMocks()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"isDisabled": false,
				"allowedTenants": [
					{
						"tenantId": "11111111-1111-1111-1111-111111111111",
						"inbound": true,
						"outbound": false
					}
				]
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
					allowed_tenants = [
						{
							tenant_id = "11111111-1111-1111-1111-111111111111"
							inbound = true
							outbound = false
						}
					]
				}`,
			},
		},
	})
}

func TestUnitTenantIsolationPolicyResource_Validate_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	setupEnvironmentHttpMocks()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"isDisabled": true,
				"allowedTenants": [
					{
						"tenantId": "11111111-1111-1111-1111-111111111111",
						"inbound": true,
						"outbound": true
					},
					{
						"tenantId": "22222222-2222-2222-2222-222222222222",
						"inbound": false,
						"outbound": true
					}
				]
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
					allowed_tenants = [
						{
							tenant_id = "11111111-1111-1111-1111-111111111111"
							inbound = true
							outbound = false
						}
					]
				}`,
			},
			{
				Config: getProviderConfig() + `
				resource "powerplatform_tenant_isolation_policy" "test" {
					is_disabled = true
					allowed_tenants = [
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
					]
				}`,
			},
		},
	})
}

func TestUnitTenantIsolationPolicyResource_Validate_Delete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	setupEnvironmentHttpMocks()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfig() + `
				resource "powerplatform_tenant_isolation_policy" "test" {
					is_disabled = false
					allowed_tenants = [
						{
							tenant_id = "11111111-1111-1111-1111-111111111111"
							inbound = true
							outbound = false
						}
					]
				}`,
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

	setupEnvironmentHttpMocks()

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
					allowed_tenants = [
						{
							tenant_id = "invalid-tenant-id"
							inbound = true 
							outbound = false
						}
					]
				}`,
				ExpectError: regexp.MustCompile("Client error when creating tenant isolation policy: Invalid tenant ID format"),
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
                    allowed_tenants = [
                        {
                            tenant_id = ""  
                            inbound = true
                            outbound = false
                        }
                    ]
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
                    allowed_tenants = [
                        {
                            tenant_id = "11111111-1111-1111-1111-111111111111"
                            inbound = true
                            outbound = false
                        }
                    ]
                }`,
				ExpectError: regexp.MustCompile("The argument \"is_disabled\" is required"),
			},
		},
	})
}

func TestUnitTenantIsolationPolicyValidation_Invalid_TenantId(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	setupEnvironmentHttpMocks()

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
                    allowed_tenants = [
                        {
                            tenant_id = "invalid-guid-format"
                            inbound = true
                            outbound = false
                        }
                    ]
                }`,
				ExpectError: regexp.MustCompile("Client error when creating tenant isolation policy: Invalid tenant ID format"),
			},
		},
	})
}

// Helper functions for acceptance tests
func testAccTenantIsolationPolicy_basic() string {
	return getProviderConfig() + `
resource "powerplatform_tenant_isolation_policy" "test" {
  allowed_tenants {
    tenant_id = "11111111-1111-1111-1111-111111111111"
    inbound  = true
    outbound = true
  }
}
`
}

func testAccTenantIsolationPolicy_update() string {
	return getProviderConfig() + `
resource "powerplatform_tenant_isolation_policy" "test" {
  is_disabled = true
  allowed_tenants {
    tenant_id = "11111111-1111-1111-1111-111111111111"
    inbound  = true
    outbound = true
  }
  allowed_tenants {
    tenant_id = "22222222-2222-2222-2222-222222222222"
    inbound  = true
    outbound = false
  }
}
`
}

func testAccTenantIsolationPolicy_empty() string {
	return getProviderConfig() + `
resource "powerplatform_tenant_isolation_policy" "test" {
  allowed_tenants = []
}
`
}
