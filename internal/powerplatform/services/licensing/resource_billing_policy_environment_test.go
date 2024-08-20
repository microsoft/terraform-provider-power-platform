// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package licensing_test

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	mocks "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/provider"
)

// We can't test the create method as it requires a valid subscription id and resource group id
// func TestAccBillingPolicyResourceEnvironment_Validate_Create(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { TestAccPreCheck(t) },
// 		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: provider.TestsAcceptanceProviderConfig + `
// 				resource "powerplatform_billing_policy_environment" "pay_as_you_go_policy_envs" {
// 					billing_policy_id = "00000000-0000-0000-0000-000000000000"
// 					environments      = ["00000000-0000-0000-0000-000000000001"]
// 				}`,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(helpers.GuidRegex)),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.#", "1"),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.0", "00000000-0000-0000-0000-000000000001"),
// 				),
// 			},
// 		},
// 	})
// }

func TestUnitBillingPolicyResourceEnvironment_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000000/environments?api-version=2022-03-01-preview`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("test/resource/environments/Validate_Create/get_environments_for_policy.json").String()), nil
		})

	httpmock.RegisterResponder("POST", `https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000000/environments/remove?api-version=2022-03-01-preview`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	httpmock.RegisterResponder("POST", `https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000000/environments/add?api-version=2022-03-01-preview`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsUnitProviderConfig + `
				resource "powerplatform_billing_policy_environment" "pay_as_you_go_policy_envs" {
					billing_policy_id = "00000000-0000-0000-0000-000000000000"
					environments      = ["00000000-0000-0000-0000-000000000001"]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.0", "00000000-0000-0000-0000-000000000001"),
				),
			},
		},
	})
}

// We can't test the create method as it requires a valid subscription id and resource group id
// func TestAccBillingPolicyResourceEnvironment_Validate_Update(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { TestAccPreCheck(t) },
// 		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: provider.TestsAcceptanceProviderConfig + `
// 				resource "powerplatform_billing_policy_environment" "pay_as_you_go_policy_envs" {
// 					billing_policy_id = "00000000-0000-0000-0000-000000000000"
// 					environments      = ["00000000-0000-0000-0000-000000000001"]
// 				}`,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(helpers.GuidRegex)),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.#", "1"),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.0", "00000000-0000-0000-0000-000000000001"),
// 				),
// 			},
// 			{
// 				Config: provider.TestsAcceptanceProviderConfig + `
// 				resource "powerplatform_billing_policy_environment" "pay_as_you_go_policy_envs" {
// 					billing_policy_id = "00000000-0000-0000-0000-000000000000"
// 					environments      = ["00000000-0000-0000-0000-000000000001","00000000-0000-0000-0000-000000000002","00000000-0000-0000-0000-000000000003"]
// 				}`,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(helpers.GuidRegex)),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.#", "3"),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.0", "00000000-0000-0000-0000-000000000001"),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.1", "00000000-0000-0000-0000-000000000002"),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.2", "00000000-0000-0000-0000-000000000003"),
// 				),
// 			},
// 			{
// 				Config: AcceptanceTestsProviderConfig + `
// 				resource "powerplatform_billing_policy_environment" "pay_as_you_go_policy_envs" {
// 					billing_policy_id = "00000000-0000-0000-0000-000000000000"
// 					environments      = ["00000000-0000-0000-0000-000000000001","00000000-0000-0000-0000-000000000003"]
// 				}`,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(helpers.GuidRegex)),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.#", "2"),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.0", "00000000-0000-0000-0000-000000000001"),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.1", "00000000-0000-0000-0000-000000000003"),
// 				),
// 			},
// 		},
// 	})
// }

func TestUnitBillingPolicyResourceEnvironment_Validate_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	getResponseInx := 0

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000000/environments?api-version=2022-03-01-preview`,
		func(req *http.Request) (*http.Response, error) {
			getResponseInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("test/resource/environments/Validate_Update/get_environments_for_policy_%d.json", getResponseInx)).String()), nil
		})

	httpmock.RegisterResponder("POST", `https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000000/environments/add?api-version=2022-03-01-preview`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	httpmock.RegisterResponder("POST", `https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000000/environments/remove?api-version=2022-03-01-preview`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsUnitProviderConfig + `
				resource "powerplatform_billing_policy_environment" "pay_as_you_go_policy_envs" {
					billing_policy_id = "00000000-0000-0000-0000-000000000000"
					environments      = ["00000000-0000-0000-0000-000000000001"]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.0", "00000000-0000-0000-0000-000000000001"),
				),
			},
			{
				Config: provider.TestsUnitProviderConfig + `
				resource "powerplatform_billing_policy_environment" "pay_as_you_go_policy_envs" {
					billing_policy_id = "00000000-0000-0000-0000-000000000000"
					environments      = ["00000000-0000-0000-0000-000000000001","00000000-0000-0000-0000-000000000002","00000000-0000-0000-0000-000000000003"]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.#", "3"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.0", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.1", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.2", "00000000-0000-0000-0000-000000000003"),
				),
			},
			{
				Config: provider.TestsUnitProviderConfig + `
				resource "powerplatform_billing_policy_environment" "pay_as_you_go_policy_envs" {
					billing_policy_id = "00000000-0000-0000-0000-000000000000"
					environments      = ["00000000-0000-0000-0000-000000000001","00000000-0000-0000-0000-000000000003"]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.#", "2"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.0", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.1", "00000000-0000-0000-0000-000000000003"),
				),
			},
		},
	})
}
