// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_set_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccEnvironmentGroupRuleSetResource_Validate_Create(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_group_rule_set" "example_group_rule_set" {
  environment_group_id = "bd6b30f1-e31e-4cdd-b82b-689a4b674f2f"
  rules = [
    {
      type          = "Sharing",
      resource_type = "App",
    },
    {
      type = "AdminDigest"
    }
  ]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}
