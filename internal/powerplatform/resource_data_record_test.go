// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataRecordResource_Validate_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "test_env" {
					display_name     = "powerplatform_data_record_example"
					location         = "europe"
					environment_type = "Sandbox"
					dataverse = {
					  language_code     = "1033"
					  currency_code     = "USD"
					  security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}
				
				resource "powerplatform_data_record" "data_record_sample_contact1" {
					environment_id     = powerplatform_environment.data_record_example_env.id
					table_logical_name = "contact"
					columns = {
					  firstname          = "John"
					  lastname           = "Doe"
					  telephone1         = "555-555-5555"
					  emailaddress1      = "johndoe@contoso.com"
					  address1_composite = "123 Main St\nRedmond\nWA\n98052\nUS"
					  anniversary        = "2024-04-10"
					  annualincome       = 1234.56
					  birthdate          = "2024-04-10"
					  description        = "This is the description of the the terraform \n\nsample contact"
					}
				}
				
				
				
				`,

				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}
