package powerplatform

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	mocks "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/models"
)

func TestAccConnectorsDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				data "powerplatform_connectors" "all" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					// Verify the first power app to ensure all attributes are set
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.description", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.display_name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.id", regexp.MustCompile(powerplatform_helpers.ApiIdRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.publisher", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.tier", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.type", regexp.MustCompile(powerplatform_helpers.ApiIdRegex)),
				),
			},
		},
	})
}

func TestUnitConnectorsDataSource_Validate_Read(t *testing.T) {
	clientMock := mocks.NewUnitTestsMockBapiClientInterface(t)

	connectors := make([]models.ConnectorDto, 0)
	connectors = append(connectors, models.ConnectorDto{
		Properties: models.ConnectorPropertiesDto{
			DisplayName: "Test Display Name 1",
			Description: "Test Description 1",
			Tier:        "Test Tier 1",
			Publisher:   "Test Publisher 1",
		},
		Id:   "Test Id 1",
		Name: "Test Name 1",
		Type: "Test Type 1",
	})
	connectors = append(connectors, models.ConnectorDto{
		Properties: models.ConnectorPropertiesDto{
			DisplayName: "Test Display Name 2",
			Description: "Test Description 2",
			Tier:        "Test Tier 2",
			Publisher:   "Test Publisher 2",
		},
		Id:   "Test Id 2",
		Name: "Test Name 2",
		Type: "Test Type 2",
	})

	conn := ConvertFromConnectorDto(connectors[0])

	clientMock.EXPECT().GetConnectors(gomock.Any()).Return(connectors, nil).AnyTimes()

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"powerplatform": powerPlatformProviderServerApiMock(clientMock, nil, nil),
		},
		Steps: []resource.TestStep{
			{
				Config: uniTestsProviderConfig + `
				data "powerplatform_connectors" "all" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					//Verify returned count
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.#", strconv.Itoa(len(connectors))),

					// Verify the first power app to ensure all attributes are set
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.description", conn.Description.ValueString()),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.display_name", conn.DisplayName.ValueString()),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.id", conn.Id.ValueString()),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.name", conn.Name.ValueString()),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.publisher", conn.Publisher.ValueString()),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.tier", conn.Tier.ValueString()),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.type", conn.Type.ValueString()),
				),
			},
		},
	})
}
