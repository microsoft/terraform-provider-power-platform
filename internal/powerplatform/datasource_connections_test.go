// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConnectionssDataSource_Validate_Read(t *testing.T) {

	t.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				data "powerplatform_connections" "all_connections" {
					environment_id = "0f555a0d-488a-ecd5-995c-47a85a167255"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}
