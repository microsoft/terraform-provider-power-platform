package powerplatform

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func aTestUnitDlpPolicyDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `=~^https://login\.microsoftonline\.com/*./v2.0/\.well-known/openid-configuration`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusTeapot, ""), nil
		})

	httpmock.RegisterNoResponder(func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusTeapot, ""), nil
	})

	// connectors := make([]models.ConnectorDto, 0)
	// connectors = append(connectors, models.ConnectorDto{
	// 	Properties: models.ConnectorPropertiesDto{
	// 		DisplayName: "Test Display Name 1",
	// 		Description: "Test Description 1",
	// 		Tier:        "Test Tier 1",
	// 		Publisher:   "Test Publisher 1",
	// 	},
	// 	Id:   "Test Id 1",
	// 	Name: "Test Name 1",
	// 	Type: "Test Type 1",
	// })
	// connectors = append(connectors, models.ConnectorDto{
	// 	Properties: models.ConnectorPropertiesDto{
	// 		DisplayName: "Test Display Name 2",
	// 		Description: "Test Description 2",
	// 		Tier:        "Test Tier 2",
	// 		Publisher:   "Test Publisher 2",
	// 	},
	// 	Id:   "Test Id 2",
	// 	Name: "Test Name 2",
	// 	Type: "Test Type 2",
	// })

	//conn := ConvertFromConnectorDto(connectors[0])

	//clientMock.EXPECT().GetConnectors(gomock.Any()).Return(connectors, nil).AnyTimes()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: UniTestsProviderConfig + `
				data "powerplatform_data_loss_prevention_policies" "all" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					//Verify returned count
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "connectors.#", "2"),

					// Verify the first power app to ensure all attributes are set
					// resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.description", conn.Description.ValueString()),
					// resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.display_name", conn.DisplayName.ValueString()),
					// resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.id", conn.Id.ValueString()),
					// resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.name", conn.Name.ValueString()),
					// resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.publisher", conn.Publisher.ValueString()),
					// resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.tier", conn.Tier.ValueString()),
					// resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.type", conn.Type.ValueString()),
				),
			},
		},
	})
}
