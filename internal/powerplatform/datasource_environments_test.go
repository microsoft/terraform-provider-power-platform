package powerplatform

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccEnvironmentsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: AcceptanceTestsProviderConfig + `
				data "powerplatform_environments" "all" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					// Verify the first power app to ensure all attributes are set
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.display_name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.domain", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.environment_type", regexp.MustCompile(`^(Default|Sandbox|Developer)$`)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.language_code", regexp.MustCompile(`^(1033|1031)$`)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.organization_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.security_group_id", regexp.MustCompile(powerplatform_helpers.GuidOrEmptyValueRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.url", regexp.MustCompile(powerplatform_helpers.UrlValidStringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.location", regexp.MustCompile(`^(unitedstates|europe)$`)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.version", regexp.MustCompile(powerplatform_helpers.VersionRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.currency_code", regexp.MustCompile(powerplatform_helpers.StringRegex)),
				),
			},
		},
	})
}

func TestUnitEnvironmentsDataSource_Validate_Read(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock_helpers.ActivateOAuthHttpMocks()
	mock_helpers.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/datasource/Validate_Read/get_environments.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/datasource/Validate_Read/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000002?%24expand=permissions%2Cproperties.capacity&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/datasource/Validate_Read/get_environment_00000000-0000-0000-0000-000000000002.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: UnitTestsProviderConfig + `
				data "powerplatform_environments" "all" {}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.#", "2"),

					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.display_name", "Admin AdminOnMicrosoft's Environment"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.domain", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.environment_type", "Developer"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.language_code", "1033"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.organization_id", "6450637c-f9a8-4988-8cf7-b03723d51ab7"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.security_group_id", ""),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.url", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.location", "europe"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.version", "9.2.23092.00206"),
					//resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.linked_app_type", ""),
					//resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.linked_app_id", ""),
					//resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.linked_app_url", ""),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.currency_code", "PLN"),
				),
			},
		},
	})
}
