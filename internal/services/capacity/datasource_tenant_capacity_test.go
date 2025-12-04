// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package capacity_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitTenantCapacityDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/licensing/tenantCapacity?api-version=2022-03-01-preview`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_tenant_capacity.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_tenant_capacity" "capacity" {
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "license_model_type", "StorageDriven"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.capacity_type", "Database"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.capacity_units", "MB"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.total_capacity", "20480"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.max_capacity", "0"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.consumption.actual", "3150.945068359375"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.consumption.rated", "2379.89794921875"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.status", "Available"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.consumption.actual_updated_on", "2025-11-30T08:52:06.3330393+00:00"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.consumption.rated_updated_on", "2025-11-30T08:52:06.3330393+00:00"),
				),
			},
		},
	})
}

func TestAccTenantCapacityDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_tenant_capacity" "capacity" {
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
