// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package licensing_test

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
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

func TestAccBillingPolicyResource_Validate_Create(t *testing.T) {
	rgName := "power-platform-billing-" + mocks.TestName() + strconv.Itoa(rand.IntN(9999))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azapi": {
				VersionConstraint: constants.AZAPI_PROVIDER_VERSION_CONSTRAINT,
				Source:            "azure/azapi",
			},
		},
		Steps: []resource.TestStep{
			{
				ResourceName: "powerplatform_billing_policy.pay_as_you_go",
				Config: `
				data "azapi_client_config" "current" {}

				resource "azapi_resource" "rg_example" {
					type      = "Microsoft.Resources/resourceGroups@2021-04-01"
					location  = "East US"
					name      = "` + rgName + `"
				}

				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "` + strings.ReplaceAll(mocks.TestName(), "_", "") + `"
					location = "unitedstates"
					status   = "Enabled"
					billing_instrument = {
						resource_group  = azapi_resource.rg_example.name
						subscription_id = data.azapi_client_config.current.subscription_id
					}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "name", strings.ReplaceAll(mocks.TestName(), "_", "")),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "location", "unitedstates"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "status", "Enabled"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "billing_instrument.resource_group", rgName),
					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "billing_instrument.subscription_id", regexp.MustCompile(helpers.GuidRegex)),
				),
			},
		},
	})
}

func TestUnitTestBillingPolicyResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/billingPolicies?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("test/resource/policies/Validate_Create/post_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("test/resource/policies/Validate_Create/post_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "payAsYouGoBillingPolicyExample"
					location = "europe"
					status   = "Enabled"
					billing_instrument = {
					  resource_group  = "resource_group_name"
					  subscription_id = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify placeholder id attribute.
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", "00000000-0000-0000-0000-000000000001"),

					// Verify the first power app to ensure all attributes are set.
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "name", "payAsYouGoBillingPolicyExample"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "status", "Enabled"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "billing_instrument.id", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource_group_name/providers/Microsoft.PowerPlatform/accounts/payAsYouGoBillingPolicyExample"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "billing_instrument.resource_group", "resource_group_name"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "billing_instrument.subscription_id", "00000000-0000-0000-0000-000000000000"),
				),
			},
		},
	})
}

func TestAccBillingPolicy_Validate_Update(t *testing.T) {
	rgName := "power-platform-billing-" + mocks.TestName() + strconv.Itoa(rand.IntN(9999))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azapi": {
				VersionConstraint: constants.AZAPI_PROVIDER_VERSION_CONSTRAINT,
				Source:            "azure/azapi",
			},
		},
		Steps: []resource.TestStep{
			{
				ResourceName: "powerplatform_billing_policy.pay_as_you_go",
				Config: `
				data "azapi_client_config" "current" {}

				resource "azapi_resource" "rg_example" {
					type      = "Microsoft.Resources/resourceGroups@2021-04-01"
					location  = "East US"
					name      = "` + rgName + `"
				}

				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "` + strings.ReplaceAll(mocks.TestName(), "_", "") + `"
					location = "unitedstates"
					status   = "Enabled"
					billing_instrument = {
						resource_group  = azapi_resource.rg_example.name
						subscription_id = data.azapi_client_config.current.subscription_id
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(helpers.GuidRegex)),
				),
			},
			{

				ResourceName: "powerplatform_billing_policy.pay_as_you_go",
				Config: `
				data "azapi_client_config" "current" {}

				resource "azapi_resource" "rg_example" {
					type      = "Microsoft.Resources/resourceGroups@2021-04-01"
					location  = "East US"
					name      = "` + rgName + `"
				}

				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "` + strings.ReplaceAll(mocks.TestName(), "_", "") + `1"
					location = "unitedstates"
					status   = "Enabled"
					billing_instrument = {
						resource_group  = azapi_resource.rg_example.name
						subscription_id = data.azapi_client_config.current.subscription_id
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "name", strings.ReplaceAll(mocks.TestName(), "_", "")+"1"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "status", "Enabled"),
				),
			},
			{

				ResourceName: "powerplatform_billing_policy.pay_as_you_go",
				Config: `
				data "azapi_client_config" "current" {}

				resource "azapi_resource" "rg_example" {
					type      = "Microsoft.Resources/resourceGroups@2021-04-01"
					location  = "East US"
					name      = "` + rgName + `"
				}

				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "` + strings.ReplaceAll(mocks.TestName(), "_", "") + `1"
					location = "unitedstates"
					status   = "Disabled"
					billing_instrument = {
						resource_group  = azapi_resource.rg_example.name
						subscription_id = data.azapi_client_config.current.subscription_id
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "name", strings.ReplaceAll(mocks.TestName(), "_", "")+"1"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "status", "Disabled"),
				),
			},
		},
	})
}

func TestUnitTestBillingPolicy_Validate_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	getResponseInx := 0

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/billingPolicies?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("test/resource/policies/Validate_Update/post_billing_policy_1.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			getResponseInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("test/resource/policies/Validate_Update/get_billing_policy_%d.json", getResponseInx)).String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("PUT", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("test/resource/policies/Validate_Update/put_billing_policy_1.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "payAsYouGoBillingPolicyExample"
					location = "europe"
					status   = "Enabled"
					billing_instrument = {
					  resource_group  = "resource_group_name"
					  subscription_id = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(helpers.GuidRegex)),
				),
			},
			{
				Config: `
				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "payAsYouGoBillingPolicyExample1"
					location = "europe"
					status   = "Enabled"
					billing_instrument = {
					  resource_group  = "resource_group_name"
					  subscription_id = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "name", "payAsYouGoBillingPolicyExample1"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "status", "Enabled"),
				),
			},
			{
				Config: `
				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "payAsYouGoBillingPolicyExample1"
					location = "europe"
					status   = "Disabled"
					billing_instrument = {
					  resource_group  = "resource_group_name"
					  subscription_id = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "name", "payAsYouGoBillingPolicyExample1"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "status", "Disabled"),
				),
			},
		},
	})
}

func TestAccBillingPolicy_Validate_Update_ForceRecreate(t *testing.T) {
	rgName := "power-platform-billing-" + mocks.TestName() + strconv.Itoa(rand.IntN(9999))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azapi": {
				VersionConstraint: constants.AZAPI_PROVIDER_VERSION_CONSTRAINT,
				Source:            "azure/azapi",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
				data "azapi_client_config" "current" {}

				resource "azapi_resource" "rg_example" {
					type      = "Microsoft.Resources/resourceGroups@2021-04-01"
					location  = "East US"
					name      = "` + rgName + `"
				}

				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "` + strings.ReplaceAll(mocks.TestName(), "_", "") + `"
					location = "europe"
					status   = "Enabled"
					billing_instrument = {
						resource_group  = azapi_resource.rg_example.name
						subscription_id = data.azapi_client_config.current.subscription_id
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(helpers.GuidRegex)),
				),
			},
			{

				Config: `
				data "azapi_client_config" "current" {}

				resource "azapi_resource" "rg_example" {
					type      = "Microsoft.Resources/resourceGroups@2021-04-01"
					location  = "East US"
					name      = "` + rgName + `"
				}

				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "` + strings.ReplaceAll(mocks.TestName(), "_", "") + `"
					location = "switzerland"
					status   = "Enabled"
					billing_instrument = {
						resource_group  = azapi_resource.rg_example.name
						subscription_id = data.azapi_client_config.current.subscription_id
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "location", "switzerland"),
				),
			},
			{
				Config: `
				data "azapi_client_config" "current" {}

				resource "azapi_resource" "rg_example" {
					type      = "Microsoft.Resources/resourceGroups@2021-04-01"
					location  = "East US"
					name      = "` + rgName + `1"
				}

				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "` + strings.ReplaceAll(mocks.TestName(), "_", "") + `"
					location = "switzerland"
					status   = "Enabled"
					billing_instrument = {
						resource_group  = azapi_resource.rg_example.name
						subscription_id = data.azapi_client_config.current.subscription_id
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "location", "switzerland"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "billing_instrument.resource_group", rgName+"1"),
				),
			},
		},
	})
}

func TestUnitTestBillingPolicy_Validate_Update_ForceRecreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	getResponseInx := 0
	postResponseInx := 0

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/billingPolicies?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			postResponseInx++
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File(fmt.Sprintf("test/resource/policies/Validate_Update_ForceRecreate/post_billing_policy_%d.json", postResponseInx)).String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.powerplatform\.com/licensing/billingPolicies/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			getResponseInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("test/resource/policies/Validate_Update_ForceRecreate/get_billing_policy_%d.json", getResponseInx)).String()), nil
		})
	httpmock.RegisterResponder("DELETE", `=~^https://api\.powerplatform\.com/licensing/billingPolicies/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "payAsYouGoBillingPolicyExample"
					location = "europe"
					status   = "Enabled"
					billing_instrument = {
						resource_group  = "resource_group_name"
						subscription_id = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", "00000000-0000-0000-0000-000000000001"),
				),
			},
			{
				Config: `
				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "payAsYouGoBillingPolicyExample"
					location = "switzerland"
					status   = "Enabled"
					billing_instrument = {
						resource_group  = "resource_group_name"
						subscription_id = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "location", "switzerland"),
				),
			},
			{
				Config: `
				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "payAsYouGoBillingPolicyExample"
					location = "switzerland"
					status   = "Enabled"
					billing_instrument = {
						resource_group  = "resource_group_name1"
						subscription_id = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", "00000000-0000-0000-0000-000000000003"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "location", "switzerland"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "billing_instrument.resource_group", "resource_group_name1"),
				),
			},
			{
				Config: `
				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "payAsYouGoBillingPolicyExample"
					location = "switzerland"
					status   = "Enabled"
					billing_instrument = {
						resource_group  = "resource_group_name1"
						subscription_id = "00000000-0000-0000-0000-000000000001"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", "00000000-0000-0000-0000-000000000004"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "location", "switzerland"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "billing_instrument.resource_group", "resource_group_name1"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "billing_instrument.subscription_id", "00000000-0000-0000-0000-000000000001"),
				),
			},
		},
	})
}

func TestUnitTestBillingPolicy_Validate_Create_WithoutFinalStatusInPostResponse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/billingPolicies?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("test/resource/policies/Validate_Create_WithoutFinalStatusInPostResponse/post_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("test/resource/policies/Validate_Create_WithoutFinalStatusInPostResponse/get_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "payAsYouGoBillingPolicyExample"
					location = "europe"
					status   = "Enabled"
					billing_instrument = {
					  resource_group  = "resource_group_name"
					  subscription_id = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "name", "payAsYouGoBillingPolicyExample"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "status", "Enabled"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "billing_instrument.id", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource_group_name/providers/Microsoft.PowerPlatform/accounts/payAsYouGoBillingPolicyExample"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "billing_instrument.resource_group", "resource_group_name"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "billing_instrument.subscription_id", "00000000-0000-0000-0000-000000000000"),
				),
			},
		},
	})
}

// commenting out until we can properly test timeouts
//
// func TestUnitTestBillingPolicy_Validate_Create_TimeoutWithoutFinalStatus(t *testing.T) {
// 	httpmock.Activate()
// 	defer httpmock.DeactivateAndReset()

// 	mocks.ActivateEnvironmentHttpMocks()

// 	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/billingPolicies?api-version=2022-03-01-preview",
// 		func(req *http.Request) (*http.Response, error) {
// 			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("test/resource/policies/Validate_Create_TimeoutWithoutFinalStatus/post_billing_policy.json").String()), nil
// 		})

// 	httpmock.RegisterResponder("GET", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
// 		func(req *http.Request) (*http.Response, error) {
// 			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("test/resource/policies/Validate_Create_TimeoutWithoutFinalStatus/get_billing_policy.json").String()), nil
// 		})

// 	httpmock.RegisterResponder("DELETE", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
// 		func(req *http.Request) (*http.Response, error) {
// 			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
// 		})

// 	resource.Test(t, resource.TestCase{
// 		IsUnitTest:               true,
// 		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: `
// 				resource "powerplatform_billing_policy" "pay_as_you_go" {
// 					name     = "payAsYouGoBillingPolicyExample"
// 					location = "europe"
// 					status   = "Enabled"
// 					billing_instrument = {
// 					  resource_group  = "resource_group_name"
// 					  subscription_id = "00000000-0000-0000-0000-000000000000"
// 					}
// 				}`,

// 				ExpectError: regexp.MustCompile("timeout reached while waiting for billing policy to reach a terminal state"),

// 				Check: resource.ComposeAggregateTestCheckFunc(),
// 			},
// 		},
// 	})
// }

const (
	billingPoliciesUrl = "https://api.powerplatform.com/licensing/billingPolicies?api-version=2022-03-01-preview"
	billingPolicyUrl   = "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview"
	billingPolicyName  = "powerplatform_billing_policy.pay_as_you_go"
)

// validateBillingPolicyRequest checks the status and meter states sent to the API.
// A nil expected map means that the request must not contain the payGoEntitlements property.
func validateBillingPolicyRequest(req *http.Request, expectedStatus string, expected map[string]bool) error {
	raw, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	body := map[string]json.RawMessage{}
	if err := json.Unmarshal(raw, &body); err != nil {
		return err
	}

	status := ""
	if err := json.Unmarshal(body["status"], &status); err != nil || status != expectedStatus {
		return fmt.Errorf("expected status '%s' in request body %s", expectedStatus, raw)
	}

	entitlementsJson, found := body["payGoEntitlements"]
	if expected == nil {
		if found {
			return fmt.Errorf("expected no payGoEntitlements in request body %s", raw)
		}
		return nil
	}

	entitlements := []struct {
		EntitlementId   string `json:"entitlementId"`
		PayAsYouGoState bool   `json:"payAsYouGoState"`
	}{}
	if err := json.Unmarshal(entitlementsJson, &entitlements); err != nil {
		return err
	}
	if len(entitlements) != len(expected) {
		return fmt.Errorf("expected %d payGoEntitlements in request body %s", len(expected), raw)
	}
	for _, entitlement := range entitlements {
		if state, ok := expected[entitlement.EntitlementId]; !ok || state != entitlement.PayAsYouGoState {
			return fmt.Errorf("unexpected payGoEntitlements in request body %s", raw)
		}
	}
	return nil
}

func allBillingPolicyEntitlements(enabled ...string) map[string]bool {
	entitlements := map[string]bool{
		"CloudFlowRuns": false, "PAUnattendedRPA": false, "PAAttendedRPA": false,
		"PowerPagesAuthenticated": false, "PowerPagesAnonymous": false,
		"File": false, "Log": false, "Database": false,
		"MCSMessages": false, "AppPass": false, "W365APAYGO": false,
	}
	for _, id := range enabled {
		entitlements[id] = true
	}
	return entitlements
}

func badRequestResponse(err error) (*http.Response, error) {
	return httpmock.NewStringResponse(http.StatusBadRequest, err.Error()), nil
}

func checkBillingPolicyMeters(dataverse, powerApps, powerAutomate, powerPages, copilotStudio, windows365ForAgents bool) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(billingPolicyName, "meters.dataverse", strconv.FormatBool(dataverse)),
		resource.TestCheckResourceAttr(billingPolicyName, "meters.power_apps", strconv.FormatBool(powerApps)),
		resource.TestCheckResourceAttr(billingPolicyName, "meters.power_automate", strconv.FormatBool(powerAutomate)),
		resource.TestCheckResourceAttr(billingPolicyName, "meters.power_pages", strconv.FormatBool(powerPages)),
		resource.TestCheckResourceAttr(billingPolicyName, "meters.copilot_studio", strconv.FormatBool(copilotStudio)),
		resource.TestCheckResourceAttr(billingPolicyName, "meters.windows_365_for_agents", strconv.FormatBool(windows365ForAgents)),
	)
}

func TestUnitBillingPolicyResource_Validate_Create_WithoutMeters(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("POST", billingPoliciesUrl,
		func(req *http.Request) (*http.Response, error) {
			// Without meters in the configuration the request must stay identical to earlier provider versions.
			if err := validateBillingPolicyRequest(req, "Enabled", nil); err != nil {
				return badRequestResponse(err)
			}
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("test/resource/policies/Validate_Create_WithoutMeters/post_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("GET", billingPolicyUrl,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("test/resource/policies/Validate_Create_WithoutMeters/post_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", billingPolicyUrl,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "payAsYouGoBillingPolicyExample"
					location = "europe"
					status   = "Enabled"
					billing_instrument = {
					  resource_group  = "resource_group_name"
					  subscription_id = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(billingPolicyName, "id", "00000000-0000-0000-0000-000000000001"),
					checkBillingPolicyMeters(false, false, false, false, false, false),
				),
			},
		},
	})
}

func TestUnitBillingPolicyResource_Validate_Create_WithMeters(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("POST", billingPoliciesUrl,
		func(req *http.Request) (*http.Response, error) {
			if err := validateBillingPolicyRequest(req, "Enabled", map[string]bool{"MCSMessages": true}); err != nil {
				return badRequestResponse(err)
			}
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("test/resource/policies/Validate_Create_WithMeters/post_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("GET", billingPolicyUrl,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("test/resource/policies/Validate_Create_WithMeters/post_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", billingPolicyUrl,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "payAsYouGoBillingPolicyExample"
					location = "europe"
					billing_instrument = {
					  resource_group  = "resource_group_name"
					  subscription_id = "00000000-0000-0000-0000-000000000000"
					}
					meters = {
					  copilot_studio = true
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(billingPolicyName, "status", "Enabled"),
					checkBillingPolicyMeters(false, false, false, false, true, false),
				),
			},
		},
	})
}

func TestUnitBillingPolicyResource_Validate_Update_Meters(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	phase := 1
	expectedPutRequests := []map[string]bool{
		allBillingPolicyEntitlements("AppPass"),
		allBillingPolicyEntitlements("AppPass"),
	}

	httpmock.RegisterResponder("POST", billingPoliciesUrl,
		func(req *http.Request) (*http.Response, error) {
			if err := validateBillingPolicyRequest(req, "Enabled", map[string]bool{"MCSMessages": true}); err != nil {
				return badRequestResponse(err)
			}
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("test/resource/policies/Validate_Update_Meters/post_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("GET", billingPolicyUrl,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("test/resource/policies/Validate_Update_Meters/get_billing_policy_%d.json", phase)).String()), nil
		})

	httpmock.RegisterResponder("PUT", billingPolicyUrl,
		func(req *http.Request) (*http.Response, error) {
			// Status isn't configured, so the current status must be sent instead of an empty value.
			if err := validateBillingPolicyRequest(req, "Enabled", expectedPutRequests[phase-1]); err != nil {
				return badRequestResponse(err)
			}
			response := httpmock.File(fmt.Sprintf("test/resource/policies/Validate_Update_Meters/put_billing_policy_%d.json", phase)).String()
			phase++
			return httpmock.NewStringResponse(http.StatusOK, response), nil
		})

	httpmock.RegisterResponder("DELETE", billingPolicyUrl,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "payAsYouGoBillingPolicyExample"
					location = "europe"
					billing_instrument = {
					  resource_group  = "resource_group_name"
					  subscription_id = "00000000-0000-0000-0000-000000000000"
					}
					meters = {
					  copilot_studio = true
					}
				}`,
				Check: checkBillingPolicyMeters(false, false, false, false, true, false),
			},
			{
				Config: `
				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "payAsYouGoBillingPolicyExample"
					location = "europe"
					billing_instrument = {
					  resource_group  = "resource_group_name"
					  subscription_id = "00000000-0000-0000-0000-000000000000"
					}
					meters = {
					  copilot_studio = false
					  power_apps     = true
					}
				}`,
				Check: checkBillingPolicyMeters(false, true, false, false, false, false),
			},
			{
				Config: `
				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "payAsYouGoBillingPolicyExample1"
					location = "europe"
					billing_instrument = {
					  resource_group  = "resource_group_name"
					  subscription_id = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(billingPolicyName, "name", "payAsYouGoBillingPolicyExample1"),
					checkBillingPolicyMeters(false, true, false, false, false, false),
				),
			},
		},
	})
}

func TestUnitBillingPolicyResource_Validate_Update_PreservesMetersWhenNotConfigured(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	phase := 1

	httpmock.RegisterResponder("POST", billingPoliciesUrl,
		func(req *http.Request) (*http.Response, error) {
			if err := validateBillingPolicyRequest(req, "Enabled", nil); err != nil {
				return badRequestResponse(err)
			}
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("test/resource/policies/Validate_Update_PreservesMeters/post_billing_policy.json").String()), nil
		})

	// The Copilot Studio meter is turned on outside of Terraform, for example in the Power Platform admin center.
	httpmock.RegisterResponder("GET", billingPolicyUrl,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("test/resource/policies/Validate_Update_PreservesMeters/get_billing_policy_%d.json", phase)).String()), nil
		})

	httpmock.RegisterResponder("PUT", billingPolicyUrl,
		func(req *http.Request) (*http.Response, error) {
			// Renaming the policy must not turn off meters that the configuration doesn't manage.
			if err := validateBillingPolicyRequest(req, "Enabled", allBillingPolicyEntitlements("MCSMessages")); err != nil {
				return badRequestResponse(err)
			}
			phase++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("test/resource/policies/Validate_Update_PreservesMeters/put_billing_policy_1.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", billingPolicyUrl,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "payAsYouGoBillingPolicyExample"
					location = "europe"
					status   = "Enabled"
					billing_instrument = {
					  resource_group  = "resource_group_name"
					  subscription_id = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: checkBillingPolicyMeters(false, false, false, false, false, false),
			},
			{
				Config: `
				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "payAsYouGoBillingPolicyExample1"
					location = "europe"
					status   = "Enabled"
					billing_instrument = {
					  resource_group  = "resource_group_name"
					  subscription_id = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(billingPolicyName, "name", "payAsYouGoBillingPolicyExample1"),
					checkBillingPolicyMeters(false, false, false, false, true, false),
				),
			},
		},
	})
}

func TestUnitBillingPolicyResource_Validate_Read_NumericProductCategory(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("POST", billingPoliciesUrl,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("test/resource/policies/Validate_Read_NumericProductCategory/post_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("GET", billingPolicyUrl,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("test/resource/policies/Validate_Read_NumericProductCategory/get_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", billingPolicyUrl,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	config := `
	resource "powerplatform_billing_policy" "pay_as_you_go" {
		name     = "payAsYouGoBillingPolicyExample"
		location = "europe"
		status   = "Enabled"
		billing_instrument = {
		  resource_group  = "resource_group_name"
		  subscription_id = "00000000-0000-0000-0000-000000000000"
		}
	}`

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				Config: config,
				Check:  checkBillingPolicyMeters(false, true, false, false, true, false),
			},
		},
	})
}

func TestUnitBillingPolicyResource_Validate_Update_Error(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("POST", billingPoliciesUrl,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("test/resource/policies/Validate_Update_Error/post_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("GET", billingPolicyUrl,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("test/resource/policies/Validate_Update_Error/get_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("PUT", billingPolicyUrl,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusBadRequest, httpmock.File("test/resource/policies/Validate_Update_Error/put_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", billingPolicyUrl,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "payAsYouGoBillingPolicyExample"
					location = "europe"
					status   = "Enabled"
					billing_instrument = {
					  resource_group  = "resource_group_name"
					  subscription_id = "00000000-0000-0000-0000-000000000000"
					}
				}`,
			},
			{
				Config: `
				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "payAsYouGoBillingPolicyExample"
					location = "europe"
					status   = "Enabled"
					billing_instrument = {
					  resource_group  = "resource_group_name"
					  subscription_id = "00000000-0000-0000-0000-000000000000"
					}
					meters = {
					  copilot_studio = true
					}
				}`,
				ExpectError: regexp.MustCompile("BillingPolicyInvalidEntitlementPayGoEntitlements"),
			},
		},
	})
}

func TestAccBillingPolicy_Validate_Meters(t *testing.T) {
	rgName := "power-platform-billing-" + mocks.TestName() + strconv.Itoa(rand.IntN(9999))
	policyName := strings.ReplaceAll(mocks.TestName(), "_", "")

	config := func(name string, meters string) string {
		return `
		data "azapi_client_config" "current" {}

		resource "azapi_resource" "rg_example" {
			type      = "Microsoft.Resources/resourceGroups@2021-04-01"
			location  = "East US"
			name      = "` + rgName + `"
		}

		resource "powerplatform_billing_policy" "pay_as_you_go" {
			name     = "` + name + `"
			location = "unitedstates"
			status   = "Enabled"
			billing_instrument = {
				resource_group  = azapi_resource.rg_example.name
				subscription_id = data.azapi_client_config.current.subscription_id
			}
			` + meters + `
		}`
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azapi": {
				VersionConstraint: constants.AZAPI_PROVIDER_VERSION_CONSTRAINT,
				Source:            "azure/azapi",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: config(policyName, `meters = {
					copilot_studio = true
				}`),
				Check: checkBillingPolicyMeters(false, false, false, false, true, false),
			},
			{
				Config: config(policyName, `meters = {
					copilot_studio = false
					power_apps     = true
					dataverse      = true
				}`),
				Check: checkBillingPolicyMeters(true, true, false, false, false, false),
			},
			{
				Config: config(policyName+"1", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(billingPolicyName, "name", policyName+"1"),
					checkBillingPolicyMeters(true, true, false, false, false, false),
				),
			},
		},
	})
}
