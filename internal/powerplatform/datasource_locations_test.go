package powerplatform

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
)

func TestAccLocationsDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				data "powerplatform_locations" "all_locations" {
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_locations.all_locations", "id", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.powerplatform_locations.all_locations", "locations.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.powerplatform_locations.all_locations", "locations.0.id", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_locations.all_locations", "locations.0.name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_locations.all_locations", "locations.0.display_name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_locations.all_locations", "locations.0.code", regexp.MustCompile(powerplatform_helpers.StringRegex)),
				),
			},
		},
	})
}

func TestUnitLocationsDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/locations/tests/datasource/Validate_Read/get_locations.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				data "powerplatform_locations" "all_locations" {
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_locations.all_locations", "id", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.powerplatform_locations.all_locations", "applications.#", "2"),
					resource.TestCheckResourceAttr("data.powerplatform_locations.all_locations", "applications.0.id", "4bbd5362-21f6-47a8-bcd9-e2a75e8242ef"),
					resource.TestCheckResourceAttr("data.powerplatform_locations.all_locations", "applications.0.name", "Dynamics 365 Customer Voice"),
					resource.TestCheckResourceAttr("data.powerplatform_locations.all_locations", "applications.0.display_name", "MicrosoftFormsPro"),
					resource.TestCheckResourceAttr("data.powerplatform_locations.all_locations", "applications.0.code", "MicrosoftFormsPro"),
					resource.TestCheckResourceAttr("data.powerplatform_locations.all_locations", "applications.1.id", "f50a3059-435a-401b-a7ee-1bca67da5657"),
					resource.TestCheckResourceAttr("data.powerplatform_locations.all_locations", "applications.1.name", "Intelligent Order Management Portal"),
					resource.TestCheckResourceAttr("data.powerplatform_locations.all_locations", "applications.1.display_name", "msdyn_IOMOrderReturnsPortalAnchor"),
					resource.TestCheckResourceAttr("data.powerplatform_locations.all_locations", "applications.1.code", "msdyn_IOMOrderReturnsPortalAnchor"),
				),
			},
		},
	})
}
