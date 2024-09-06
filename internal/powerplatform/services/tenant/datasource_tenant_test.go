// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestUnitTenantDataSource_Validate_Read(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/tenant?api-version=2021-04-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_tenant.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: constants.TestsUnitProviderConfig + `
				data "powerplatform_tenant" "tenant" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_tenant.tenant", "tenant_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant.tenant", "state", "Enabled"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant.tenant", "location", "unitedstates"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant.tenant", "aad_country_geo", "unitedstates"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant.tenant", "data_storage_geo", "unitedstates"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant.tenant", "default_environment_geo", "unitedstates"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant.tenant", "aad_data_boundary", "none"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant.tenant", "fed_ramp_high_certification_required", "false"),
				),
			},
		},
	})
}

func TestAccTenantDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsAcceptanceProviderConfig + `
				data "powerplatform_tenant" "tenant" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerplatform_tenant.tenant", "tenant_id"),
					resource.TestCheckResourceAttrSet("data.powerplatform_tenant.tenant", "state"),
					resource.TestCheckResourceAttrSet("data.powerplatform_tenant.tenant", "location"),
					resource.TestCheckResourceAttrSet("data.powerplatform_tenant.tenant", "aad_country_geo"),
					resource.TestCheckResourceAttrSet("data.powerplatform_tenant.tenant", "data_storage_geo"),
					resource.TestCheckResourceAttrSet("data.powerplatform_tenant.tenant", "default_environment_geo"),
					resource.TestCheckResourceAttrSet("data.powerplatform_tenant.tenant", "aad_data_boundary"),
					resource.TestCheckResourceAttrSet("data.powerplatform_tenant.tenant", "fed_ramp_high_certification_required"),
				),
			},
		},
	})
}
