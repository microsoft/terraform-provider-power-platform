// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package locations_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccMacroRegionsDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_macro_regions" "all_macro_regions" {
				}`,

				// Tenants that provision by datacenter location return an empty list, so only the
				// presence of the attribute can be asserted without knowing the tenant mode.
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerplatform_macro_regions.all_macro_regions", "macro_regions.#"),
				),
			},
		},
	})
}

func TestUnitMacroRegionsDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read_Macro_Regions/get_macro_regions.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_macro_regions" "all_macro_regions" {
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_macro_regions.all_macro_regions", "macro_regions.#", "6"),
					resource.TestCheckResourceAttr("data.powerplatform_macro_regions.all_macro_regions", "macro_regions.0.macro_region_id", "asia-pacific"),
					resource.TestCheckResourceAttr("data.powerplatform_macro_regions.all_macro_regions", "macro_regions.0.display_name", "Asia Pacific"),
					resource.TestCheckResourceAttr("data.powerplatform_macro_regions.all_macro_regions", "macro_regions.0.data_residency_note", "Your data will reside within the Asia-Pacific region."),
					resource.TestCheckResourceAttr("data.powerplatform_macro_regions.all_macro_regions", "macro_regions.2.macro_region_id", "eu-efta"),
					resource.TestCheckResourceAttr("data.powerplatform_macro_regions.all_macro_regions", "macro_regions.2.display_name", "European Union & EFTA"),
					resource.TestCheckResourceAttr("data.powerplatform_macro_regions.all_macro_regions", "macro_regions.3.macro_region_id", "europe-uk"),
					resource.TestCheckResourceAttr("data.powerplatform_macro_regions.all_macro_regions", "macro_regions.3.display_name", "Europe & UK"),
					// The live API omits dataResidencyNote for europe-uk, which must surface as null.
					resource.TestCheckNoResourceAttr("data.powerplatform_macro_regions.all_macro_regions", "macro_regions.3.data_residency_note"),
					resource.TestCheckResourceAttr("data.powerplatform_macro_regions.all_macro_regions", "macro_regions.4.macro_region_id", "north-america"),
					resource.TestCheckResourceAttr("data.powerplatform_macro_regions.all_macro_regions", "macro_regions.5.macro_region_id", "the-americas"),
				),
			},
		},
	})
}

func TestUnitMacroRegionsDataSource_Validate_Read_Classic_Tenant(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read_Macro_Regions_Classic_Tenant/get_macro_regions.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_macro_regions" "all_macro_regions" {
				}`,

				// A tenant that provisions by location returns "macroRegions": [], which is an empty
				// list rather than an error.
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_macro_regions.all_macro_regions", "macro_regions.#", "0"),
				),
			},
		},
	})
}
