// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataRecordDatasource_Validate_Create(t *testing.T) {

	t.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				data "powerplatform_data_records" "example_data_records" {
					environment_id = "838f76c8-a192-e59c-a835-089ad8cfb047"
					entity_collection = "systemusers(1f70a364-5019-ef11-840b-002248ca35c3)"
					return_total_records_count = true
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestAccDataRecordDatasource_Validate_Create2(t *testing.T) {

	t.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				data "powerplatform_data_records" "example_data_records" {
					environment_id = "838f76c8-a192-e59c-a835-089ad8cfb047"
					entity_collection = "systemusers"
					select            = ["firstname", "lastname", "createdon"]
					//top               = 2
					return_total_records_count = true
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}
