// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package publisher_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitPublisherDataSource_Validate_Read_ByUniqueName(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()
	registerPublisherEnvironmentMock()

	httpmock.RegisterResponder("GET", "https://"+testPublisherHost+"/api/data/v9.2/publishers?%24filter=uniquename+eq+%27contoso%27",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[`+publisherCreateResponse()+`]}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "data.powerplatform_publisher.example",
				Config: `
data "powerplatform_publisher" "example" {
  environment_id = "` + testEnvironmentID + `"
  uniquename     = "contoso"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_publisher.example", "publisher_id", testPublisherID),
					resource.TestCheckResourceAttr("data.powerplatform_publisher.example", "friendly_name", "Contoso Publisher"),
					resource.TestCheckResourceAttr("data.powerplatform_publisher.example", "address.#", "2"),
				),
			},
		},
	})
}
