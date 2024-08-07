// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

// We can't test the create method as it requires a valid subscription id and resource group id
// func TestAccBillingPolicyResource_Validate_Create(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: TestsAcceptanceProviderConfig + `
// 				resource "powerplatform_billing_policy" "pay_as_you_go" {
// 					name     = "` + mock_helpers.TestName() + `"
// 					location = "europe"
// 					status   = "Enabled"
// 					billing_instrument = {
// 					  resource_group  = "resource_group_name"
// 					  subscription_id = "00000000-0000-0000-0000-000000000000"
// 					}
// 				}`,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					//Verify placeholder id attribute
// 					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
// 					// Verify the first power app to ensure all attributes are set
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "name", mock_helpers.TestName()),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "location", "europe"),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "status", "Enabled"),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "billing_instrument.resource_group", "resource_group_name"),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "billing_instrument.subscription_id", "00000000-0000-0000-0000-000000000000"),
// 				),
// 			},
// 		},
// 	})
// }

func TestUnitTestBillingPolicyResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/BillingPolicies?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("services/licensing/test/resource/policies/Validate_Create/post_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/licensing/test/resource/policies/Validate_Create/post_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://api.powerplatform.com/licensing/BillingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsUnitProviderConfig + `
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
					//Verify placeholder id attribute
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", "00000000-0000-0000-0000-000000000001"),

					// Verify the first power app to ensure all attributes are set
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

// We can't test the create method as it requires a valid subscription id and resource group id
// func TestAccBillingPolicy_Validate_Update(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: TestsAcceptanceProviderConfig + `
// 				resource "powerplatform_billing_policy" "pay_as_you_go" {
// 					name     = "` + mock_helpers.TestName() + `"
// 					location = "europe"
// 					status   = "Enabled"
// 					billing_instrument = {
// 						resource_group  = "resource_group_name"
// 						subscription_id = "00000000-0000-0000-0000-000000000000"
// 					}
// 				}`,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
// 				),
// 			},
// 			{
// 				Config: TestsAcceptanceProviderConfig + `
// 				resource "powerplatform_billing_policy" "pay_as_you_go" {
// 					name     = "` + mock_helpers.TestName() + `1"
// 					location = "europe"
// 					status   = "Enabled"
// 					billing_instrument = {
// 						resource_group  = "resource_group_name"
// 						subscription_id = "00000000-0000-0000-0000-000000000000"
// 					}
// 				}`,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "name", mock_helpers.TestName()+"1"),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "status", "Enabled"),
// 				),
// 			},
// 			{
// 				Config: TestsAcceptanceProviderConfig + `
// 				resource "powerplatform_billing_policy" "pay_as_you_go" {
// 					name     = "` + mock_helpers.TestName() + `1"
// 					location = "europe"
// 					status   = "Disabled"
// 					billing_instrument = {
// 						resource_group  = "resource_group_name"
// 						subscription_id = "00000000-0000-0000-0000-000000000000"
// 					}
// 				}`,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "name", mock_helpers.TestName()+"1"),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "status", "Disabled"),
// 				),
// 			},
// 		},
// 	})
// }

func TestUnitTestBillingPolicy_Validate_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateEnvironmentHttpMocks()

	getResponseInx := 0

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/BillingPolicies?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("services/licensing/test/resource/policies/Validate_Update/post_billing_policy_1.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			getResponseInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/licensing/test/resource/policies/Validate_Update/get_billing_policy_%d.json", getResponseInx)).String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://api.powerplatform.com/licensing/BillingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("PUT", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/licensing/test/resource/policies/Validate_Update/put_billing_policy_1.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsUnitProviderConfig + `
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
					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
				),
			},
			{
				Config: TestsUnitProviderConfig + `
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
					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "name", "payAsYouGoBillingPolicyExample1"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "status", "Enabled"),
				),
			},
			{
				Config: TestsUnitProviderConfig + `
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
					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "name", "payAsYouGoBillingPolicyExample1"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "status", "Disabled"),
				),
			},
		},
	})
}

// We can't test the create method as it requires a valid subscription id and resource group id
// func TestAccBillingPolicy_Validate_Update_ForceRecreate(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: TestsAcceptanceProviderConfig + `
// 				resource "powerplatform_billing_policy" "pay_as_you_go" {
// 					name     = "` + mock_helpers.TestName() + `"
// 					location = "europe"
// 					status   = "Enabled"
// 					billing_instrument = {
// 						resource_group  = "resource_group_name"
// 						subscription_id = "00000000-0000-0000-0000-000000000000"
// 					}
// 				}`,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
// 				),
// 			},
// 			{
// 				Config: TestsAcceptanceProviderConfig + `
// 				resource "powerplatform_billing_policy" "pay_as_you_go" {
// 					name     = "` + mock_helpers.TestName() + `"
// 					location = "switzerland"
// 					status   = "Enabled"
// 					billing_instrument = {
// 						resource_group  = "resource_group_name"
// 						subscription_id = "00000000-0000-0000-0000-000000000000"
// 					}
// 				}`,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "location", "switzerland"),
// 				),
// 			},
// 			{
// 				Config: TestsAcceptanceProviderConfig + `
// 				resource "powerplatform_billing_policy" "pay_as_you_go" {
// 					name     = "` + mock_helpers.TestName() + `"
// 					location = "switzerland"
// 					status   = "Enabled"
// 					billing_instrument = {
// 						resource_group  = "resource_group_name1"
// 						subscription_id = "00000000-0000-0000-0000-000000000000"
// 					}
// 				}`,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "location", "switzerland"),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "billing_instrument.resource_group", "resource_group_name1"),
// 				),
// 			},
// 			{
// 				Config: TestsAcceptanceProviderConfig + `
// 				resource "powerplatform_billing_policy" "pay_as_you_go" {
// 					name     = "` + mock_helpers.TestName() + `"
// 					location = "switzerland"
// 					status   = "Enabled"
// 					billing_instrument = {
// 						resource_group  = "resource_group_name1"
// 						subscription_id = "00000000-0000-0000-0000-000000000001"
// 					}
// 				}`,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "location", "switzerland"),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "billing_instrument.resource_group", "resource_group_name1"),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "billing_instrument.subscription_id", "00000000-0000-0000-0000-000000000001"),
// 				),
// 			},
// 		},
// 	})
// }

func TestUnitTestBillingPolicy_Validate_Update_ForceRecreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateEnvironmentHttpMocks()

	getResponseInx := 0
	postResponseInx := 0

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/BillingPolicies?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			postResponseInx++
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File(fmt.Sprintf("services/licensing/test/resource/policies/Validate_Update_ForceRecreate/post_billing_policy_%d.json", postResponseInx)).String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.powerplatform\.com/licensing/billingPolicies/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			getResponseInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/licensing/test/resource/policies/Validate_Update_ForceRecreate/get_billing_policy_%d.json", getResponseInx)).String()), nil
		})
	httpmock.RegisterResponder("DELETE", `=~^https://api\.powerplatform\.com/licensing/BillingPolicies/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsUnitProviderConfig + `
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
				Config: TestsUnitProviderConfig + `
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
				Config: TestsUnitProviderConfig + `
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
				Config: TestsUnitProviderConfig + `
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

	mock_helpers.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/BillingPolicies?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("services/licensing/test/resource/policies/Validate_Create_WithoutFinalStatusInPostResponse/post_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/licensing/test/resource/policies/Validate_Create_WithoutFinalStatusInPostResponse/get_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://api.powerplatform.com/licensing/BillingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsUnitProviderConfig + `
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

func TestUnitTestBillingPolicy_Validate_Create_TimeoutWithoutFinalStatus(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/BillingPolicies?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("services/licensing/test/resource/policies/Validate_Create_TimeoutWithoutFinalStatus/post_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/licensing/test/resource/policies/Validate_Create_TimeoutWithoutFinalStatus/get_billing_policy.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://api.powerplatform.com/licensing/BillingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsUnitProviderConfig + `
				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "payAsYouGoBillingPolicyExample"
					location = "europe"
					status   = "Enabled"
					billing_instrument = {
					  resource_group  = "resource_group_name"
					  subscription_id = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				ExpectError: regexp.MustCompile("timeout reached while waiting for billing policy to reach a terminal state"),

				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}
