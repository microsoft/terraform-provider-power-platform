package powerplatform

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
)

func TestAccCurrenciesDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				data "powerplatform_currencies" "all_currencies_for_unitedstates" {
					location = "unitedstates"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_currencies.all_currencies_for_unitedstates", "id", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.powerplatform_currencies.all_currencies_for_unitedstates", "currencies.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.powerplatform_currencies.all_currencies_for_unitedstates", "currencies.0.id", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_currencies.all_currencies_for_unitedstates", "currencies.0.name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_currencies.all_currencies_for_unitedstates", "currencies.0.display_name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_currencies.all_currencies_for_unitedstates", "currencies.0.locale_id", regexp.MustCompile(powerplatform_helpers.StringRegex)),
				),
			},
		},
	})
}

func TestUnitCurrenciesDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentCurrencies?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/currencies/tests/datasource/Validate_Read/get_currencies.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				data "powerplatform_currencies" "all_currencies_for_unitedstates" {
					location = "unitedstates"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_currencies.all_currencies_for_unitedstates", "id", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.powerplatform_currencies.all_currencies_for_unitedstates", "currencies.#", "112"),
					resource.TestCheckResourceAttr("data.powerplatform_currencies.all_currencies_for_unitedstates", "currencies.0.id", "/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentCurrencies/DJF"),
					resource.TestCheckResourceAttr("data.powerplatform_currencies.all_currencies_for_unitedstates", "currencies.0.name", "DJF"),
					resource.TestCheckResourceAttr("data.powerplatform_currencies.all_currencies_for_unitedstates", "currencies.0.type", "Microsoft.BusinessAppPlatform/locations/environmentCurrencies"),
					resource.TestCheckResourceAttr("data.powerplatform_currencies.all_currencies_for_unitedstates", "currencies.0.symbol", "Fdj"),
					resource.TestCheckResourceAttr("data.powerplatform_currencies.all_currencies_for_unitedstates", "currencies.1.id", "/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentCurrencies/ZAR"),
					resource.TestCheckResourceAttr("data.powerplatform_currencies.all_currencies_for_unitedstates", "currencies.1.name", "ZAR"),
					resource.TestCheckResourceAttr("data.powerplatform_currencies.all_currencies_for_unitedstates", "currencies.1.type", "Microsoft.BusinessAppPlatform/locations/environmentCurrencies"),
					resource.TestCheckResourceAttr("data.powerplatform_currencies.all_currencies_for_unitedstates", "currencies.1.symbol", "R"),
				),
			},
		},
	})
}
