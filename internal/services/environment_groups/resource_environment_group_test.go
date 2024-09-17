// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package environment_groups_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccEnvironmentGroupResource_Validate_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_group" "test_env_group" {
					display_name = "test_env_group"
					description = "test env group"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group.test_env_group", "display_name", "test_env_group"),
					resource.TestCheckResourceAttr("powerplatform_environment_group.test_env_group", "description", "test env group"),
					resource.TestMatchResourceAttr("powerplatform_environment_group.test_env_group", "id", regexp.MustCompile(helpers.GuidRegex)),
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
			resp := httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resources/get_environment_group.json").String())
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environmentGroups/00000000-0000-0000-0000-000000000001?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resources/get_environment_group.json").String())
			return resp, nil
		},
	)

	httpmock.RegisterResponder("DELETE", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environmentGroups/00000000-0000-0000-0000-000000000001?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			return resp, nil
		},
	)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_group" "test_env_group" {
					display_name = "test_env_group"
					description = "test env group"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group.test_env_group", "display_name", "test_env_group"),
					resource.TestCheckResourceAttr("powerplatform_environment_group.test_env_group", "description", "test env group"),
					resource.TestCheckResourceAttrSet("powerplatform_environment_group.test_env_group", "id"),
					resource.TestCheckResourceAttr("powerplatform_environment_group.test_env_group", "id", "00000000-0000-0000-0000-000000000001"),
				),
			},
		},
	})
}
