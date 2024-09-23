// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package licensing_test

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccBillingPolicyResourceEnvironment_Validate_Create(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azapi": {
				VersionConstraint: ">= 1.15.0",
				Source:            "azure/azapi",
			},
		},
		Steps: []resource.TestStep{
			{
				ResourceName: "powerplatform_billing_policy_environment.pay_as_you_go_policy_envs",
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

				resource "powerplatform_environment" "env1" {
					display_name     = "billing_policy_example_environment_1_` + mocks.TestName() + `"	
					location         = "unitedstates"
					environment_type = "Sandbox"
				}

				resource "powerplatform_environment" "env2" {
					display_name     = "billing_policy_example_environment_2_` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
				}

				resource "powerplatform_billing_policy_environment" "pay_as_you_go_policy_envs" {
					billing_policy_id = powerplatform_billing_policy.pay_as_you_go.id
					environments      = [powerplatform_environment.env1.id, powerplatform_environment.env2.id]
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.#", "2"),
					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.0", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.1", regexp.MustCompile(helpers.GuidRegex)),
				),
			},
		},
	})
}

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
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
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

func TestAccBillingPolicyResourceEnvironment_Validate_Update(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azapi": {
				VersionConstraint: ">= 1.15.0",
				Source:            "azure/azapi",
			},
		},
		Steps: []resource.TestStep{
			{
				ResourceName: "powerplatform_billing_policy_environment.pay_as_you_go_policy_envs",
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

				resource "powerplatform_environment" "env1" {
					display_name     = "billing_policy_example_environment_1_` + mocks.TestName() + `"	
					location         = "unitedstates"
					environment_type = "Sandbox"
				}

				resource "powerplatform_environment" "env2" {
					display_name     = "billing_policy_example_environment_2_` + mocks.TestName() + `"	
					location         = "unitedstates"
					environment_type = "Sandbox"
				}

				resource "powerplatform_environment" "env3" {
					display_name     = "billing_policy_example_environment_3_` + mocks.TestName() + `"	
					location         = "unitedstates"
					environment_type = "Sandbox"
				}

				resource "powerplatform_billing_policy_environment" "pay_as_you_go_policy_envs" {
					billing_policy_id = powerplatform_billing_policy.pay_as_you_go.id
					environments      = [powerplatform_environment.env1.id]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.#", "1"),
				),
			},
			{
				ResourceName: "powerplatform_billing_policy_environment.pay_as_you_go_policy_envs",
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

				resource "powerplatform_environment" "env1" {
					display_name     = "billing_policy_example_environment_1_` + mocks.TestName() + `"	
					location         = "unitedstates"
					environment_type = "Sandbox"
				}

				resource "powerplatform_environment" "env2" {
					display_name     = "billing_policy_example_environment_2_` + mocks.TestName() + `"	
					location         = "unitedstates"
					environment_type = "Sandbox"
				}

				resource "powerplatform_environment" "env3" {
					display_name     = "billing_policy_example_environment_3_` + mocks.TestName() + `"	
					location         = "unitedstates"
					environment_type = "Sandbox"
				}

				resource "powerplatform_billing_policy_environment" "pay_as_you_go_policy_envs" {
					billing_policy_id = powerplatform_billing_policy.pay_as_you_go.id
					environments      = [powerplatform_environment.env1.id, powerplatform_environment.env2.id, powerplatform_environment.env3.id]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.#", "3"),
				),
			},
			{
				ResourceName: "powerplatform_billing_policy_environment.pay_as_you_go_policy_envs",
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

				resource "powerplatform_environment" "env1" {
					display_name     = "billing_policy_example_environment_1_` + mocks.TestName() + `"	
					location         = "unitedstates"
					environment_type = "Sandbox"
				}

				resource "powerplatform_environment" "env2" {
					display_name     = "billing_policy_example_environment_2_` + mocks.TestName() + `"	
					location         = "unitedstates"
					environment_type = "Sandbox"
				}

				resource "powerplatform_environment" "env3" {
					display_name     = "billing_policy_example_environment_3_` + mocks.TestName() + `"	
					location         = "unitedstates"
					environment_type = "Sandbox"
				}

				resource "powerplatform_billing_policy_environment" "pay_as_you_go_policy_envs" {
					billing_policy_id = powerplatform_billing_policy.pay_as_you_go.id
					environments      = [powerplatform_environment.env1.id, powerplatform_environment.env3.id]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.#", "2"),
				),
			},
		},
	})
}

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
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
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
				Config: `
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
				Config: `
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
