// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package publisher_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

const (
	testEnvironmentID = "00000000-0000-0000-0000-000000000001"
	testPublisherID   = "11111111-1111-1111-1111-111111111111"
	testPublisherHost = "00000000-0000-0000-0000-000000000001.crm4.dynamics.com"
)

func TestUnitPublisherResource_Validate_CRUD(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()
	registerPublisherEnvironmentMock()

	currentResponse := publisherCreateResponse()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://%s/api/data/v9.2/publishers", testPublisherHost),
		func(req *http.Request) (*http.Response, error) {
			body, _ := io.ReadAll(req.Body)
			if !strings.Contains(string(body), `"friendlyname":"Contoso Publisher"`) {
				return httpmock.NewStringResponse(http.StatusBadRequest, `{"error":"missing friendly name"}`), nil
			}
			if !strings.Contains(string(body), `"customizationoptionvalueprefix":77074`) {
				return httpmock.NewStringResponse(http.StatusBadRequest, `{"error":"missing derived customization option value prefix"}`), nil
			}

			resp := httpmock.NewStringResponse(http.StatusNoContent, "")
			resp.Header.Set("OData-EntityId", fmt.Sprintf("https://%s/api/data/v9.2/publishers(%s)", testPublisherHost, testPublisherID))
			return resp, nil
		})

	httpmock.RegisterResponder("GET", encodedPublisherURL(),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, currentResponse), nil
		})

	httpmock.RegisterResponder("PATCH", encodedPublisherURL(),
		func(req *http.Request) (*http.Response, error) {
			body, _ := io.ReadAll(req.Body)
			if !strings.Contains(string(body), `"friendlyname":"Updated Contoso Publisher"`) {
				return httpmock.NewStringResponse(http.StatusBadRequest, `{"error":"missing updated friendly name"}`), nil
			}
			if !strings.Contains(string(body), `"customizationoptionvalueprefix":72710`) {
				return httpmock.NewStringResponse(http.StatusBadRequest, `{"error":"missing explicit customization option value prefix override"}`), nil
			}
			if !strings.Contains(string(body), `"address2_city":null`) {
				return httpmock.NewStringResponse(http.StatusBadRequest, `{"error":"expected address slot 2 to be cleared"}`), nil
			}

			currentResponse = publisherUpdateResponse()
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("DELETE", encodedPublisherURL(),
		httpmock.NewStringResponder(http.StatusNoContent, ""))

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "powerplatform_publisher.example",
				Config: `
resource "powerplatform_publisher" "example" {
  environment_id                      = "` + testEnvironmentID + `"
  uniquename                          = "contoso"
  friendly_name                       = "Contoso Publisher"
  customization_prefix                = "cts"
  description                         = "Initial publisher"
  email_address                       = "publisher@contoso.example"
  supporting_website_url              = "https://contoso.example"

  address = [
    {
      slot         = 1
      line1        = "1 Collins Street"
      city         = "Melbourne"
      country      = "Australia"
      postal_code  = "3000"
      telephone1   = "+61-3-5555-0101"
    },
    {
      slot         = 2
      line1        = "100 Queen Street"
      city         = "Auckland"
      country      = "New Zealand"
      postal_code  = "1010"
      telephone1   = "+64-9-555-0102"
    }
  ]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_publisher.example", "id", testPublisherID),
					resource.TestCheckResourceAttr("powerplatform_publisher.example", "friendly_name", "Contoso Publisher"),
					resource.TestCheckResourceAttr("powerplatform_publisher.example", "customization_option_value_prefix", "77074"),
					resource.TestCheckResourceAttr("powerplatform_publisher.example", "address.#", "2"),
					resource.TestCheckResourceAttr("powerplatform_publisher.example", "address.0.slot", "1"),
					resource.TestCheckResourceAttr("powerplatform_publisher.example", "address.1.slot", "2"),
				),
			},
			{
				ResourceName: "powerplatform_publisher.example",
				Config: `
resource "powerplatform_publisher" "example" {
  environment_id                      = "` + testEnvironmentID + `"
  uniquename                          = "contoso"
  friendly_name                       = "Updated Contoso Publisher"
  customization_prefix                = "cts"
  customization_option_value_prefix   = 72710
  description                         = "Updated publisher"
  email_address                       = "updated@contoso.example"
  supporting_website_url              = "https://support.contoso.example"

  address = [
    {
      slot         = 1
      line1        = "200 Collins Street"
      city         = "Sydney"
      country      = "Australia"
      postal_code  = "2000"
      telephone1   = "+61-2-5555-0103"
    }
  ]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_publisher.example", "friendly_name", "Updated Contoso Publisher"),
					resource.TestCheckResourceAttr("powerplatform_publisher.example", "customization_option_value_prefix", "72710"),
					resource.TestCheckResourceAttr("powerplatform_publisher.example", "address.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_publisher.example", "address.0.city", "Sydney"),
				),
			},
			{
				ResourceName:      "powerplatform_publisher.example",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     testEnvironmentID + "_" + testPublisherID,
			},
		},
	})
}

func registerPublisherEnvironmentMock() {
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/"+testEnvironmentID+"?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
  "name": "`+testEnvironmentID+`",
  "properties": {
    "linkedEnvironmentMetadata": {
      "instanceUrl": "https://`+testPublisherHost+`/"
    }
  }
}`), nil
		})
}

func publisherCreateResponse() string {
	return `{
  "publisherid": "` + testPublisherID + `",
  "friendlyname": "Contoso Publisher",
  "uniquename": "contoso",
  "customizationprefix": "cts",
  "customizationoptionvalueprefix": 77074,
  "description": "Initial publisher",
  "emailaddress": "publisher@contoso.example",
  "supportingwebsiteurl": "https://contoso.example",
  "isreadonly": false,
  "address1_addressid": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
  "address1_city": "Melbourne",
  "address1_country": "Australia",
  "address1_line1": "1 Collins Street",
  "address1_postalcode": "3000",
  "address1_telephone1": "+61-3-5555-0101",
  "address2_addressid": "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
  "address2_city": "Auckland",
  "address2_country": "New Zealand",
  "address2_line1": "100 Queen Street",
  "address2_postalcode": "1010",
  "address2_telephone1": "+64-9-555-0102"
}`
}

func publisherUpdateResponse() string {
	return `{
  "publisherid": "` + testPublisherID + `",
  "friendlyname": "Updated Contoso Publisher",
  "uniquename": "contoso",
  "customizationprefix": "cts",
  "customizationoptionvalueprefix": 72710,
  "description": "Updated publisher",
  "emailaddress": "updated@contoso.example",
  "supportingwebsiteurl": "https://support.contoso.example",
  "isreadonly": false,
  "address1_addressid": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
  "address1_city": "Sydney",
  "address1_country": "Australia",
  "address1_line1": "200 Collins Street",
  "address1_postalcode": "2000",
  "address1_telephone1": "+61-2-5555-0103"
}`
}

func encodedPublisherURL() string {
	return fmt.Sprintf("https://%s/api/data/v9.2/publishers%%28%s%%29", testPublisherHost, testPublisherID)
}
