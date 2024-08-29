// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestUnitTestTenantDataSource_Validate_Read(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/tenant?api-version=2021-04-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/tenant/tests/datasource/Validate_Read/get_tenant.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsUnitProviderConfig + `
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
