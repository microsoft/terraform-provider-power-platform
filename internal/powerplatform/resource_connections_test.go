// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConnectionssResource_Validate_Create(t *testing.T) {

	t.Setenv("TF_ACC", "1")
	t.Setenv("TF_LOG", "WARN")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_connection" "new_connection" {
					environment_id = "00000000-0000-0000-0000-000000000000"
					name = "shared_azureopenai"
					display_name = "OpenAI Connection"
					connection_parameters = jsonencode({
						"azureOpenAIResourceName":"aaa",
						"azureOpenAIApiKey":"bbb",
						"azureSearchEndpointUrl":"ccc",
						"azureSearchApiKey":"dddd"
					})
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}
