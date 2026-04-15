// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package publisher_test

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
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
					resource.TestCheckResourceAttr("data.powerplatform_publisher.example", "id", testPublisherID),
					resource.TestCheckResourceAttr("data.powerplatform_publisher.example", "friendly_name", "Contoso Publisher"),
					resource.TestCheckResourceAttr("data.powerplatform_publisher.example", "address.#", "2"),
				),
			},
		},
	})
}

func TestAccPublisherDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source: "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: publisherAcceptanceDataSourceByUniqueNameConfig(mocks.TestName(), "terraformpublisherds", "Terraform Publisher DS", "tpd"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_publisher.example", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("data.powerplatform_publisher.example", "uniquename", "terraformpublisherds"),
					resource.TestCheckResourceAttr("data.powerplatform_publisher.example", "friendly_name", "Terraform Publisher DS"),
				),
			},
			{
				Config: publisherAcceptanceDataSourceByIDConfig(mocks.TestName(), "terraformpublisherds", "Terraform Publisher DS", "tpd"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_publisher.example", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttrPair("data.powerplatform_publisher.example", "id", "powerplatform_publisher.example", "id"),
				),
			},
		},
	})
}

func publisherAcceptanceDataSourceByUniqueNameConfig(environmentDisplayName, uniqueName, friendlyName, customizationPrefix string) string {
	return fmt.Sprintf(`
resource "powerplatform_environment" "environment" {
  display_name     = "%s"
  location         = "unitedstates"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

resource "time_sleep" "wait_120_seconds" {
  depends_on      = [powerplatform_environment.environment]
  create_duration = "120s"
}

resource "powerplatform_publisher" "example" {
  depends_on           = [time_sleep.wait_120_seconds]
  environment_id       = powerplatform_environment.environment.id
  uniquename           = "%s"
  friendly_name        = "%s"
  customization_prefix = "%s"
}

data "powerplatform_publisher" "example" {
  environment_id = powerplatform_environment.environment.id
  uniquename     = powerplatform_publisher.example.uniquename
}
`, environmentDisplayName, uniqueName, friendlyName, customizationPrefix)
}

func publisherAcceptanceDataSourceByIDConfig(environmentDisplayName, uniqueName, friendlyName, customizationPrefix string) string {
	return fmt.Sprintf(`
resource "powerplatform_environment" "environment" {
  display_name     = "%s"
  location         = "unitedstates"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

resource "time_sleep" "wait_120_seconds" {
  depends_on      = [powerplatform_environment.environment]
  create_duration = "120s"
}

resource "powerplatform_publisher" "example" {
  depends_on           = [time_sleep.wait_120_seconds]
  environment_id       = powerplatform_environment.environment.id
  uniquename           = "%s"
  friendly_name        = "%s"
  customization_prefix = "%s"
}

data "powerplatform_publisher" "example" {
  environment_id = powerplatform_environment.environment.id
  id             = powerplatform_publisher.example.id
}
`, environmentDisplayName, uniqueName, friendlyName, customizationPrefix)
}
