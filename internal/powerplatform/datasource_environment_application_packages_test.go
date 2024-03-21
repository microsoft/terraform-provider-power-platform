package powerplatform

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
)

func TestAccEnvironmentApplicationPackagesDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "env" {
					display_name      = "env_application_acceptance_test"
					location          = "europe"
					language_code     = "1033"
					currency_code     = "USD"
					environment_type  = "Sandbox"
					security_group_id = "00000000-0000-0000-0000-000000000000"
				}

				data "powerplatform_environment_application_packages" "all_applications" {
					environment_id = powerplatform_environment.env.id
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "id", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "environment_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.application_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.application_name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.unique_name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.version", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.description", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.publisher_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.publisher_name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.learn_more_url", regexp.MustCompile(powerplatform_helpers.UrlValidStringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.state", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.application_visibility", regexp.MustCompile(powerplatform_helpers.StringRegex)),
				),
			},
		},
	})
}

func TestUnitEnvironmentApplicationPackagesDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/appmanagement/environments/00000000-0000-0000-0000-000000000001/applicationPackages?api-version=2022-03-01-preview`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/application/tests/datasource/Validate_Read/get_applications.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				data "powerplatform_environment_application_packages" "all_applications" {
					environment_id = "00000000-0000-0000-0000-000000000001"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "id", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.powerplatform_environment_application_packages.all_applications", "environment_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.#", "2"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.application_id", "4bbd5362-21f6-47a8-bcd9-e2a75e8242ef"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.application_name", "Dynamics 365 Customer Voice"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.unique_name", "MicrosoftFormsPro"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.1.application_id", "f50a3059-435a-401b-a7ee-1bca67da5657"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.1.application_name", "Intelligent Order Management Portal"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.1.unique_name", "msdyn_IOMOrderReturnsPortalAnchor"),
				),
			},
		},
	})
}
