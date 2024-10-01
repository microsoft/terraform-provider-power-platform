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
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccBillingPoliciesDataSource_Validate_Read(t *testing.T) {
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
				ResourceName: "powerplatform_billing_policies.all",
				Config: `
				data "azapi_client_config" "current" {}

				resource "azapi_resource" "rg_example" {
					type      = "Microsoft.Resources/resourceGroups@2021-04-01"
					location  = "East US"
					name      = "power-platform-billing-` + mocks.TestName() + `"

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

				data "powerplatform_billing_policies" "all" {
					depends_on = [powerplatform_billing_policy.pay_as_you_go]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.location", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.status", regexp.MustCompile(`^(Enabled|Disabled|DisabledByLinkedResource)$`)),
					resource.TestMatchResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.billing_instrument.resource_group", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.billing_instrument.subscription_id", regexp.MustCompile(helpers.GuidRegex)),
				),
			},
		},
	})
}

func TestUnitTestBillingPoliciesDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/licensing/billingPolicies?api-version=2022-03-01-preview`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("test/datasource/policies/get_billing_policies.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_billing_policies" "all" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.#", "2"),

					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.name", "billingpolicy1"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.location", "europe"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.status", "Enabled"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.billing_instrument.resource_group", "rg-group1"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.billing_instrument.subscription_id", "00000000-0000-0000-0000-000000000001"),

					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.1.id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.1.name", "billingpolicy2"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.1.location", "europe"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.1.status", "Enabled"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.1.billing_instrument.resource_group", "rg-group2"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.1.billing_instrument.subscription_id", "00000000-0000-0000-0000-000000000002"),
				),
			},
		},
	})
}
