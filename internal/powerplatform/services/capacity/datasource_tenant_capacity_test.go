// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package capacity_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestUnitTenantCapacityDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://licensing.powerplatform.microsoft.com/v0.1-alpha/tenants/00000000-0000-0000-0000-000000000001/TenantCapacity`,
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
					tenant_id = "00000000-0000-0000-0000-000000000001"
					}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "license_model_type", "StorageDriven"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.capacity_type", "Database"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.capacity_units", "MB"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.total_capacity", "11264"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.max_capacity", "0"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.consumption.actual", "2101.093994140625"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.consumption.rated", "0"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.status", "Available"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.consumption.actual_updated_on", "2024-08-28T18:55:26.0217309+00:00"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_capacities.0.consumption.rated_updated_on", "2024-08-28T18:55:26.0217309+00:00"),
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
				Config: constants.TestsAcceptanceProviderConfig + `
				data "powerplatform_tenant" "tenant" {}

				data "powerplatform_tenant_capacity" "capacity" {
					tenant_id = data.powerplatform_tenant.tenant.tenant_id
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
