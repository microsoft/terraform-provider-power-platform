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
// func TestAccBillingPolicyResourceEnvironment_Validate_Create(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { TestAccPreCheck(t) },
// 		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: AcceptanceTestsProviderConfig + `
// 				resource "powerplatform_billing_policy_environment" "pay_as_you_go_policy_envs" {
// 					billing_policy_id = "00000000-0000-0000-0000-000000000000"
// 					environments      = ["00000000-0000-0000-0000-000000000001"]
// 				}`,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
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

	mock_helpers.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000000/environments?api-version=2022-03-01-preview`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/licensing/test/resource/environments/Validate_Create/get_environments_for_policy.json").String()), nil
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
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_billing_policy_environment" "pay_as_you_go_policy_envs" {
					billing_policy_id = "00000000-0000-0000-0000-000000000000"
					environments      = ["00000000-0000-0000-0000-000000000001"]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
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
// 		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: AcceptanceTestsProviderConfig + `
// 				resource "powerplatform_billing_policy_environment" "pay_as_you_go_policy_envs" {
// 					billing_policy_id = "00000000-0000-0000-0000-000000000000"
// 					environments      = ["00000000-0000-0000-0000-000000000001"]
// 				}`,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.#", "1"),
// 					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.0", "00000000-0000-0000-0000-000000000001"),
// 				),
// 			},
// 			{
// 				Config: AcceptanceTestsProviderConfig + `
// 				resource "powerplatform_billing_policy_environment" "pay_as_you_go_policy_envs" {
// 					billing_policy_id = "00000000-0000-0000-0000-000000000000"
// 					environments      = ["00000000-0000-0000-0000-000000000001","00000000-0000-0000-0000-000000000002","00000000-0000-0000-0000-000000000003"]
// 				}`,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
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
// 					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
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

	mock_helpers.ActivateEnvironmentHttpMocks()

	getResponseInx := 0

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000000/environments?api-version=2022-03-01-preview`,
		func(req *http.Request) (*http.Response, error) {
			getResponseInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/licensing/test/resource/environments/Validate_Update/get_environments_for_policy_%d.json", getResponseInx)).String()), nil
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
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_billing_policy_environment" "pay_as_you_go_policy_envs" {
					billing_policy_id = "00000000-0000-0000-0000-000000000000"
					environments      = ["00000000-0000-0000-0000-000000000001"]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.0", "00000000-0000-0000-0000-000000000001"),
				),
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_billing_policy_environment" "pay_as_you_go_policy_envs" {
					billing_policy_id = "00000000-0000-0000-0000-000000000000"
					environments      = ["00000000-0000-0000-0000-000000000001","00000000-0000-0000-0000-000000000002","00000000-0000-0000-0000-000000000003"]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.#", "3"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.0", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.1", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.2", "00000000-0000-0000-0000-000000000003"),
				),
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_billing_policy_environment" "pay_as_you_go_policy_envs" {
					billing_policy_id = "00000000-0000-0000-0000-000000000000"
					environments      = ["00000000-0000-0000-0000-000000000001","00000000-0000-0000-0000-000000000003"]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.#", "2"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.0", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.1", "00000000-0000-0000-0000-000000000003"),
				),
			},
		},
	})
}