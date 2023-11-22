package powerplatform

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
)

func TestAccApplicationsDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: AcceptanceTestsProviderConfig + `
				resource "powerplatform_environment" "env" {
					display_name      = "env_application_acceptance_test"
					location          = "europe"
					language_code     = "1033"
					currency_code     = "USD"
					environment_type  = "Sandbox"
					domain            = "applicationacceptancetest"
					security_group_id = "00000000-0000-0000-0000-000000000000"
				}

				data "powerplatform_applications" "all_applications" {
					environment_id = powerplatform_environment.env.id
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_applications.all_applications", "id", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.powerplatform_applications.all_applications", "environment_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_applications.all_applications", "applications.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.powerplatform_applications.all_applications", "applications.0.application_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_applications.all_applications", "applications.0.application_name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_applications.all_applications", "applications.0.unique_name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
				),
			},
		},
	})
}
