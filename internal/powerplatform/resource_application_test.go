package powerplatform

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccApplicationResource_Validate_Install(t *testing.T) {
	envDomain := fmt.Sprintf("orgtest%d", rand.Intn(100000))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: AcceptanceTestsProviderConfig + `
				resource "powerplatform_environment" "environment" {
					display_name                              = "` + envDomain + `"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                           = "USD"
					environment_type                          = "Sandbox"
					domain = "` + envDomain + `"
					security_group_id = "00000000-0000-0000-0000-000000000000"
				}
				data "powerplatform_applications" "application_to_install" {
					environment_id = data.powerplatform_environments.all_environments.environments[0].id
					name           = "Power Automate Process Mining"
					publisher_name = "Microsoft Dynamics 365"
				}
				resource "powerplatform_application" "development" {
					environment_id = powerplatform_environment.environment.id
  				unique_name = data.powerplatform_applications.application_to_install.applications[0].unique_name
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_application.development", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_application.development", "environment_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_application.development", "unique_name", "ProcessMiningAnchor"),
				),
			},
		},
	})
}

func TestUnitApplicationResource_Validate_Install(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock_helpers.ActivateOAuthHttpMocks()

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/environment/tests/resource/Validate_Create/get_environment_%s.json", id)).String()), nil
		},
	)

	httpmock.RegisterResponder("GET", "https://api.powerplatform.com/appmanagement/environments/402c2b45-f5dc-e561-869f-368544f94a13/operations/475af49d-9bca-437f-8be1-9e467f44be8a?api-version=1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/application/tests/resource/Validate_Install/get_operation.json").String()), nil
		},
	)

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/appmanagement/environments/00000000-0000-0000-0000-000000000000/applicationPackages/ProcessMiningAnchor/install?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Operation-Location", "https://api.powerplatform.com/appmanagement/environments/402c2b45-f5dc-e561-869f-368544f94a13/operations/475af49d-9bca-437f-8be1-9e467f44be8a?api-version=1")
			return resp, nil
		},
	)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: UnitTestsProviderConfig + `
				resource "powerplatform_application" "development" {
					environment_id   = "00000000-0000-0000-0000-000000000000"
					unique_name      = "ProcessMiningAnchor"
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_application.development", "id", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("powerplatform_application.development", "environment_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_application.development", "unique_name", "ProcessMiningAnchor"),
				),
			},
		},
	})
}
