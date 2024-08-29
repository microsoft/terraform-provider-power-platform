// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestUnitTestTenantCapacityDataSource_Validate_Read(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://licensing.powerplatform.microsoft.com/v0.1-alpha/tenants/00000000-0000-0000-0000-000000000001/TenantCapacity`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/capacity/tests/datasource/Validate_Read/get_tenant_capacity.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsUnitProviderConfig + `
				data "powerplatform_tenant_capacity" "capacity" {
					tenant_id = "00000000-0000-0000-0000-000000000001"
					}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerplatform_tenant_capacity.capacity", "tenant_id"),
					resource.TestCheckResourceAttrSet("data.powerplatform_tenant_capacity.capacity", "license_model_type"),
					resource.TestCheckResourceAttrSet("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.capacity_type"),
					resource.TestCheckResourceAttrSet("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.capacity_units"),
					resource.TestCheckResourceAttrSet("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.total_capacity"),
					resource.TestCheckResourceAttrSet("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.max_capacity"),
					resource.TestCheckResourceAttrSet("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.consumption.actual"),
					resource.TestCheckResourceAttrSet("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.consumption.rated"),
					resource.TestCheckResourceAttrSet("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.consumption.actual_updated_on"),
					resource.TestCheckResourceAttrSet("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.consumption.rated_updated_on"),
					resource.TestCheckResourceAttrSet("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.status"),
				),
			},
		},
	})
}
