// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package languages_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccLanguagesDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: constants.TestsAcceptanceProviderConfig + `
				data "powerplatform_languages" "all_languages_for_unitedstates" {
					location = "unitedstates"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "id", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.0.id", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.0.name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.0.display_name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.0.locale_id", regexp.MustCompile(helpers.StringRegex)),
				),
			},
		},
	})
}

func TestUnitLanguagesDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentLanguages?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_languages.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: constants.TestsUnitProviderConfig + `
				data "powerplatform_languages" "all_languages_for_unitedstates" {
					location = "unitedstates"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "id", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.#", "45"),
					resource.TestCheckResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.0.id", "/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentLanguages/1033"),
					resource.TestCheckResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.0.name", "1033"),
					resource.TestCheckResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.0.display_name", "English (United States)"),
					resource.TestCheckResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.0.locale_id", "1033"),
					resource.TestCheckResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.1.id", "/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentLanguages/1025"),
					resource.TestCheckResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.1.name", "1025"),
					resource.TestCheckResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.1.display_name", "العربية (المملكة العربية السعودية)"),
					resource.TestCheckResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.1.locale_id", "1025"),
				),
			},
		},
	})
}
