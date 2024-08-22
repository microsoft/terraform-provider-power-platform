// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package licensing_test

import (
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	mocks "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/provider"
)

// We can't test this until we don't have tenant with billing policies
func TestAccBillingPoliciesEnvironmentsDataSource_Validate_Read(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsAcceptanceProviderConfig + `
					provider "azurerm" {
					features {}
				}

				data "azurerm_client_config" "current" {
				}

				resource "azurerm_resource_group" "rg_example" {
					name     = "power-platform-billing-` + mocks.TestName() + `"
					location = "westeurope"
				}

				resource "powerplatform_billing_policy" "pay_as_you_go" {
					name     = "` + strings.ReplaceAll(mocks.TestName(), "_", "") + `"
					location = "unitedstates"
					status   = "Enabled"
					billing_instrument = {
						resource_group  = azurerm_resource_group.rg_example.name
						subscription_id = data.azurerm_client_config.current.subscription_id
					}
				}

				resource "powerplatform_billing_policy_environment" "pay_as_you_go_policy_envs" {
					billing_policy_id = powerplatform_billing_policy.pay_as_you_go.id
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_billing_policies_environments.all_pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies_environments.all_pay_as_you_go_policy_envs", "environments.#", "0"),
				),
			},
		},
	})
}

func TestUnitTestBillingPoliciesEnvironmentsDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000000/environments?api-version=2022-03-01-preview`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("test/datasource/environments/get_environments_for_policy.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsUnitProviderConfig + `
				data "powerplatform_billing_policies_environments" "all_pay_as_you_go_policy_envs" {
					billing_policy_id = "00000000-0000-0000-0000-000000000000"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies_environments.all_pay_as_you_go_policy_envs", "billing_policy_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies_environments.all_pay_as_you_go_policy_envs", "environments.#", "3"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies_environments.all_pay_as_you_go_policy_envs", "environments.0", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies_environments.all_pay_as_you_go_policy_envs", "environments.1", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies_environments.all_pay_as_you_go_policy_envs", "environments.2", "00000000-0000-0000-0000-000000000003"),
				),
			},
		},
	})
}
