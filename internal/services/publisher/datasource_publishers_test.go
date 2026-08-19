// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package publisher_test

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitPublishersDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()
	registerPublisherEnvironmentMock()

	httpmock.RegisterResponder("GET", "https://"+testPublisherHost+"/api/data/v9.2/publishers",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/PublishersDataSource_Validate_Read/get_publishers.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "data.powerplatform_publishers.example",
				Config: `
				data "powerplatform_publishers" "example" {
					environment_id = "00000000-0000-0000-0000-000000000001"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_publishers.example", "publishers.#", "2"),
					resource.TestCheckResourceAttr("data.powerplatform_publishers.example", "publishers.0.id", testPublisherID),
					resource.TestCheckResourceAttr("data.powerplatform_publishers.example", "publishers.0.friendly_name", "Contoso Publisher"),
					resource.TestCheckResourceAttr("data.powerplatform_publishers.example", "publishers.0.address.#", "2"),
					resource.TestCheckResourceAttr("data.powerplatform_publishers.example", "publishers.1.uniquename", "systempublisher"),
					resource.TestCheckResourceAttr("data.powerplatform_publishers.example", "publishers.1.address.#", "0"),
				),
			},
		},
	})
}

func TestAccPublishersDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source: "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
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
					uniquename           = "terraformpublisherds"
					friendly_name        = "Terraform Publisher DS"
					customization_prefix = "tpd"
				}

				data "powerplatform_publishers" "example" {
					environment_id = powerplatform_publisher.example.environment_id
				}
				`, mocks.TestName()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckPublisherAttrMatches("data.powerplatform_publishers.example", "terraformpublisherds", "id", regexp.MustCompile(helpers.GuidRegex)),
					testCheckPublisherAttr("data.powerplatform_publishers.example", "terraformpublisherds", "friendly_name", "Terraform Publisher DS"),
				),
			},
		},
	})
}

func findPublisherIndex(s *terraform.State, datasourceName, uniqueName string) (int, error) {
	rs := s.RootModule().Resources[datasourceName]
	if rs == nil {
		return -1, fmt.Errorf("datasource %q not found in state", datasourceName)
	}
	attrs := rs.Primary.Attributes
	var count int
	if _, err := fmt.Sscanf(attrs["publishers.#"], "%d", &count); err != nil || count == 0 {
		return -1, fmt.Errorf("no publishers found in %s", datasourceName)
	}
	for i := 0; i < count; i++ {
		if attrs[fmt.Sprintf("publishers.%d.uniquename", i)] == uniqueName {
			return i, nil
		}
	}
	return -1, fmt.Errorf("publisher with uniquename %q not found in %s (total publishers: %d)", uniqueName, datasourceName, count)
}

func testCheckPublisherAttr(datasourceName, uniqueName, attrSuffix, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		idx, err := findPublisherIndex(s, datasourceName, uniqueName)
		if err != nil {
			return err
		}
		return resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("publishers.%d.%s", idx, attrSuffix), expected)(s)
	}
}

func testCheckPublisherAttrMatches(datasourceName, uniqueName, attrSuffix string, re *regexp.Regexp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		idx, err := findPublisherIndex(s, datasourceName, uniqueName)
		if err != nil {
			return err
		}
		return resource.TestMatchResourceAttr(datasourceName, fmt.Sprintf("publishers.%d.%s", idx, attrSuffix), re)(s)
	}
}
