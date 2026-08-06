// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_based_policy_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitEnvironmentGroupRuleBasedPolicyResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	// List assignments returns empty (no existing policy)
	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/governance/ruleBasedPolicies/environmentGroups/00000000-0000-0000-0000-000000000000/assignments?api-version=2024-10-01&includeRuleSetCounts=true`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
		})

	httpmock.RegisterResponder("POST", `https://api.powerplatform.com/governance/ruleBasedPolicies?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/Validate_Create/post_policy.json").String()), nil
		})

	httpmock.RegisterResponder("POST", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000001/environmentGroups/00000000-0000-0000-0000-000000000000/assignments?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/Validate_Create/post_assignment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000001?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Create/get_policy.json").String()), nil
		})

	httpmock.RegisterResponder("PUT", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000001?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	httpmock.RegisterResponder("PATCH", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000001/removeRule?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_group_rule_based_policy" "test" {
					environment_group_id = "00000000-0000-0000-0000-000000000000"
					rule_sets = {
						advanced_connector_policies_only = {
							enabled = true
						}
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "environment_group_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.advanced_connector_policies_only.enabled", "true"),
				),
			},
		},
	})
}

func TestUnitEnvironmentGroupRuleBasedPolicyResource_Validate_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	get_policy_inx := -1

	// List assignments returns empty (no existing policy) for the create step
	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/governance/ruleBasedPolicies/environmentGroups/00000000-0000-0000-0000-000000000000/assignments?api-version=2024-10-01&includeRuleSetCounts=true`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
		})

	httpmock.RegisterResponder("POST", `https://api.powerplatform.com/governance/ruleBasedPolicies?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/Validate_Update/post_policy.json").String()), nil
		})

	httpmock.RegisterResponder("POST", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000001/environmentGroups/00000000-0000-0000-0000-000000000000/assignments?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/Validate_Update/post_assignment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000001?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			get_policy_inx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/Validate_Update/get_policy_%d.json", get_policy_inx)).String()), nil
		})

	httpmock.RegisterResponder("PUT", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000001?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Update/put_policy.json").String()), nil
		})

	httpmock.RegisterResponder("PATCH", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000001/removeRule?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_group_rule_based_policy" "test" {
					environment_group_id = "00000000-0000-0000-0000-000000000000"
					rule_sets = {
						advanced_connector_policies_only = {
							enabled = true
						}
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.advanced_connector_policies_only.enabled", "true"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment_group_rule_based_policy" "test" {
					environment_group_id = "00000000-0000-0000-0000-000000000000"
					rule_sets = {
						advanced_connector_policies_only = {
							enabled = false
						}
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.advanced_connector_policies_only.enabled", "false"),
				),
			},
		},
	})
}

func TestUnitEnvironmentGroupRuleBasedPolicyResource_Validate_Import(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("POST", `https://api.powerplatform.com/governance/ruleBasedPolicies?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/Validate_Import/get_policy.json").String()), nil
		})

	httpmock.RegisterResponder("POST", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000001/environmentGroups/00000000-0000-0000-0000-000000000000/assignments?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/Validate_Create/post_assignment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000001?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Import/get_policy.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/governance/ruleBasedPolicies/environmentGroups/00000000-0000-0000-0000-000000000000/assignments?api-version=2024-10-01&includeRuleSetCounts=true`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Import/get_assignments.json").String()), nil
		})

	httpmock.RegisterResponder("PUT", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000001?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	httpmock.RegisterResponder("PATCH", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000001/removeRule?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_group_rule_based_policy" "test" {
					environment_group_id = "00000000-0000-0000-0000-000000000000"
					rule_sets = {
						advanced_connector_policies_only = {
							enabled = true
						}
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "environment_group_id", "00000000-0000-0000-0000-000000000000"),
				),
			},
			{
				ResourceName:      "powerplatform_environment_group_rule_based_policy.test",
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "00000000-0000-0000-0000-000000000000",
			},
		},
	})
}

func TestAccEnvironmentGroupRuleBasedPolicyResource_Validate_Create(t *testing.T) {
	t.Skip("creating rule-based policies with SP is NOT yet supported")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_group" "example_group" {
					display_name = "` + mocks.TestName() + `"
					description  = "` + mocks.TestName() + `"
				}

				resource "powerplatform_environment_group_rule_based_policy" "test" {
					environment_group_id = powerplatform_environment_group.example_group.id
					rule_sets = {
						advanced_connector_policies_only = {
							enabled = true
						}
					}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.advanced_connector_policies_only.enabled", "true"),
				),
			},
		},
	})
}

func TestUnitEnvironmentGroupRuleBasedPolicyResource_Validate_Create_Conflict(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	// List assignments returns an existing policy for this environment group
	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/governance/ruleBasedPolicies/environmentGroups/00000000-0000-0000-0000-000000000000/assignments?api-version=2024-10-01&includeRuleSetCounts=true`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Create_Conflict/get_assignments.json").String()), nil
		})

	// Update the existing policy with desired state
	httpmock.RegisterResponder("PUT", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000002?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Create_Conflict/put_policy.json").String()), nil
		})

	// Read the policy after create (refresh)
	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000002?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Create_Conflict/get_policy.json").String()), nil
		})

	httpmock.RegisterResponder("PATCH", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000002/removeRule?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_group_rule_based_policy" "test" {
					environment_group_id = "00000000-0000-0000-0000-000000000000"
					rule_sets = {
						advanced_connector_policies_only = {
							enabled = true
						}
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "environment_group_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.advanced_connector_policies_only.enabled", "true"),
				),
			},
		},
	})
}

func TestUnitEnvironmentGroupRuleBasedPolicyResource_Validate_Create_CSP(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/governance/ruleBasedPolicies/environmentGroups/00000000-0000-0000-0000-000000000000/assignments?api-version=2024-10-01&includeRuleSetCounts=true`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
		})

	httpmock.RegisterResponder("POST", `https://api.powerplatform.com/governance/ruleBasedPolicies?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/Validate_Create_CSP/post_policy.json").String()), nil
		})

	httpmock.RegisterResponder("POST", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000003/environmentGroups/00000000-0000-0000-0000-000000000000/assignments?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/Validate_Create_CSP/post_assignment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000003?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Create_CSP/get_policy.json").String()), nil
		})

	httpmock.RegisterResponder("PUT", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000003?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	httpmock.RegisterResponder("PATCH", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000003/removeRule?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_group_rule_based_policy" "test" {
					environment_group_id = "00000000-0000-0000-0000-000000000000"
					rule_sets = {
						content_security_policy = {
							enabled               = true
							enabled_for_canvas    = true
							enabled_for_code_apps = true

							configuration = {
								img_src    = ["https://example.com"]
								strict_csp = true
							}

							configuration_for_canvas = {
								connect_src = ["https://canvas-api.example.com"]
								strict_csp  = true
							}

							configuration_for_code_apps = {
								script_src = ["https://cdn.example.com"]
							}
						}
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "id", "00000000-0000-0000-0000-000000000003"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.content_security_policy.enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.content_security_policy.enabled_for_canvas", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.content_security_policy.enabled_for_code_apps", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.content_security_policy.configuration.img_src.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.content_security_policy.configuration.img_src.0", "https://example.com"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.content_security_policy.configuration.strict_csp", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.content_security_policy.configuration_for_canvas.connect_src.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.content_security_policy.configuration_for_canvas.connect_src.0", "https://canvas-api.example.com"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.content_security_policy.configuration_for_canvas.strict_csp", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.content_security_policy.configuration_for_code_apps.script_src.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.content_security_policy.configuration_for_code_apps.script_src.0", "https://cdn.example.com"),
				),
			},
		},
	})
}

func TestUnitEnvironmentGroupRuleBasedPolicyResource_Validate_Create_ConnectorMgmt(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/governance/ruleBasedPolicies/environmentGroups/00000000-0000-0000-0000-000000000000/assignments?api-version=2024-10-01&includeRuleSetCounts=true`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
		})

	httpmock.RegisterResponder("POST", `https://api.powerplatform.com/governance/ruleBasedPolicies?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/Validate_Create_ConnectorMgmt/post_policy.json").String()), nil
		})

	httpmock.RegisterResponder("POST", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000004/environmentGroups/00000000-0000-0000-0000-000000000000/assignments?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/Validate_Create_ConnectorMgmt/post_assignment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000004?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Create_ConnectorMgmt/get_policy.json").String()), nil
		})

	httpmock.RegisterResponder("PUT", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000004?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	httpmock.RegisterResponder("PATCH", `https://api.powerplatform.com/governance/ruleBasedPolicies/00000000-0000-0000-0000-000000000004/removeRule?api-version=2024-10-01`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_group_rule_based_policy" "test" {
					environment_group_id = "00000000-0000-0000-0000-000000000000"
					rule_sets = {
						advanced_connector_policies = {
							allowed_connectors = [
								{
									connector_id = "shared_commondataservice"
									actions_mode = "all_allowed"
								},
								{
									connector_id    = "shared_office365"
									actions_mode    = "some_allowed"
									allowed_actions = ["SendEmail", "GetEvents"]
								}
							]
						}
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "id", "00000000-0000-0000-0000-000000000004"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.advanced_connector_policies.allowed_connectors.#", "2"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.advanced_connector_policies.allowed_connectors.0.connector_id", "shared_commondataservice"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.advanced_connector_policies.allowed_connectors.0.actions_mode", "all_allowed"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.advanced_connector_policies.allowed_connectors.1.connector_id", "shared_office365"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.advanced_connector_policies.allowed_connectors.1.actions_mode", "some_allowed"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.advanced_connector_policies.allowed_connectors.1.allowed_actions.#", "2"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.advanced_connector_policies.allowed_connectors.1.allowed_actions.0", "SendEmail"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_based_policy.test", "rule_sets.advanced_connector_policies.allowed_connectors.1.allowed_actions.1", "GetEvents"),
				),
			},
		},
	})
}
