// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package languages_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccLanguagesDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_languages" "all_languages_for_unitedstates" {
					location = "unitedstates"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
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
