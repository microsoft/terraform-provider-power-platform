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

func TestAccEnvironmentGroupResource_Validate_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsAcceptanceProviderConfig + `
				resource "powerplatform_environment_group" "test_env_group" {
					display_name = "test_env_group"
					description = "test env group"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group.test_env_group", "display_name", "test_env_group"),
					resource.TestCheckResourceAttr("powerplatform_environment_group.test_env_group", "description", "test env group"),
					resource.TestMatchResourceAttr("powerplatform_environment_group.test_env_group", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
				),
			},
		},
	})
}

func TestUnitEnvirionmentGroupResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environmentGroups?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment_groups/test/resources/get_environment_group.json").String())
			return resp, nil
		},
	)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsUnitProviderConfig + `
				resource "powerplatform_environment_group" "test_env_group" {
					display_name = "test_env_group"
					description = "test env group"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group.test_env_group", "display_name", "test_env_group"),
					resource.TestCheckResourceAttr("powerplatform_environment_group.test_env_group", "description", "test env group"),
					resource.TestCheckResourceAttrSet("powerplatform_environment_group.test_env_group", "id"),
				),
			},
		},
	})
}
