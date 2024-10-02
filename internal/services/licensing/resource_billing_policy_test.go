// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package licensing_test

import (
	"fmt"
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
	println("rgName: " + rgName)
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

					body = jsonencode({
						properties = {}
					})
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

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/BillingPolicies?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("test/resource/policies/Validate_Create/post_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("test/resource/policies/Validate_Create/post_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://api.powerplatform.com/licensing/BillingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
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

					body = jsonencode({
						properties = {}
					})
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

					body = jsonencode({
						properties = {}
					})
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

					body = jsonencode({
						properties = {}
					})
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

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/BillingPolicies?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("test/resource/policies/Validate_Update/post_billing_policy_1.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			getResponseInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("test/resource/policies/Validate_Update/get_billing_policy_%d.json", getResponseInx)).String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://api.powerplatform.com/licensing/BillingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
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

					body = jsonencode({
						properties = {}
					})
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

					body = jsonencode({
						properties = {}
					})
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

					body = jsonencode({
						properties = {}
					})
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

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/BillingPolicies?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			postResponseInx++
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File(fmt.Sprintf("test/resource/policies/Validate_Update_ForceRecreate/post_billing_policy_%d.json", postResponseInx)).String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.powerplatform\.com/licensing/billingPolicies/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			getResponseInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("test/resource/policies/Validate_Update_ForceRecreate/get_billing_policy_%d.json", getResponseInx)).String()), nil
		})
	httpmock.RegisterResponder("DELETE", `=~^https://api\.powerplatform\.com/licensing/BillingPolicies/([\d-]+)\z`,
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

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/BillingPolicies?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("test/resource/policies/Validate_Create_WithoutFinalStatusInPostResponse/post_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("test/resource/policies/Validate_Create_WithoutFinalStatusInPostResponse/get_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://api.powerplatform.com/licensing/BillingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
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

// 	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/BillingPolicies?api-version=2022-03-01-preview",
// 		func(req *http.Request) (*http.Response, error) {
// 			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("test/resource/policies/Validate_Create_TimeoutWithoutFinalStatus/post_billing_policy.json").String()), nil
// 		})

// 	httpmock.RegisterResponder("GET", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
// 		func(req *http.Request) (*http.Response, error) {
// 			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("test/resource/policies/Validate_Create_TimeoutWithoutFinalStatus/get_billing_policy.json").String()), nil
// 		})

// 	httpmock.RegisterResponder("DELETE", "https://api.powerplatform.com/licensing/BillingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
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
