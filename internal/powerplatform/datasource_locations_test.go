// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
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
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				data "powerplatform_locations" "all_locations" {
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_locations.all_locations", "id", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.powerplatform_locations.all_locations", "locations.#", "18"),
					resource.TestCheckResourceAttr("data.powerplatform_locations.all_locations", "locations.0.id", "/providers/Microsoft.BusinessAppPlatform/locations/unitedstates"),
					resource.TestCheckResourceAttr("data.powerplatform_locations.all_locations", "locations.0.name", "unitedstates"),
					resource.TestCheckResourceAttr("data.powerplatform_locations.all_locations", "locations.0.display_name", "United States"),
					resource.TestCheckResourceAttr("data.powerplatform_locations.all_locations", "locations.0.code", "NA"),
					resource.TestCheckResourceAttr("data.powerplatform_locations.all_locations", "locations.1.id", "/providers/Microsoft.BusinessAppPlatform/locations/unitedstatesfirstrelease"),
					resource.TestCheckResourceAttr("data.powerplatform_locations.all_locations", "locations.1.name", "unitedstatesfirstrelease"),
					resource.TestCheckResourceAttr("data.powerplatform_locations.all_locations", "locations.1.display_name", "Preview (United States)"),
					resource.TestCheckResourceAttr("data.powerplatform_locations.all_locations", "locations.1.code", "NA"),
				),
			},
		},
	})
}
