// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package connectors_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccConnectorsDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_connectors" "all" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(

					// Verify the first power app to ensure all attributes are set.
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.description", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.display_name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.id", regexp.MustCompile(helpers.ApiIdRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.publisher", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.tier", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.type", regexp.MustCompile(helpers.ApiIdRegex)),
				),
			},
		},
	})
}

func TestUnitConnectorsDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/connectors/metadata/virtual`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Read/get_virtual.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/connectors/metadata/unblockable`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Read/get_unblockable.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.powerapps.com/providers/Microsoft.PowerApps/apis?%24filter=environment+eq+%27~Default%27&api-version=2019-05-01&hideDlpExemptApis=true&showAllDlpEnforceableApis=true&showApisWithToS=true`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Read/get_apis.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_connectors" "all" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify returned count.
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.#", "4"),

					// Verify the first power app to ensure all attributes are set
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.description", "SharePoint helps organizations share and collaborate with colleagues, partners, and customers. You can connect to SharePoint Online or to an on-premises SharePoint 2013 or 2016 farm using the On-Premises Data Gateway to manage documents and list items."),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.display_name", "SharePoint"),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.id", "/providers/Microsoft.PowerApps/apis/shared_sharepointonline"),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.name", "shared_sharepointonline"),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.publisher", "Microsoft"),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.tier", "Standard"),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.type", "Microsoft.PowerApps/apis"),
				),
			},
		},
	})
}
