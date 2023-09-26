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

func TestAccPowerAppsDataSource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				data "powerplatform_powerapps" "all" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("data.powerplatform_powerapps.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					// Verify the first power app to ensure all attributes are set
					resource.TestMatchResourceAttr("data.powerplatform_powerapps.all", "powerapps.0.name", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_powerapps.all", "powerapps.0.environment_name", regexp.MustCompile(powerplatform_helpers.UrlValidStringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_powerapps.all", "powerapps.0.display_name", regexp.MustCompile(powerplatform_helpers.UrlValidStringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_powerapps.all", "powerapps.0.created_time", regexp.MustCompile(`^\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+([+-][0-2]\d:[0-5]\d|Z)$`)),
				),
			},
		},
	})
}

func TestUnitPowerAppsDataSource_Validate_Read(t *testing.T) {
	clientMock := mocks.NewUnitTestsMockBapiClientInterface(t)

	apps := make([]models.PowerAppBapi, 0)
	apps = append(apps, models.PowerAppBapi{
		Name: "name1",
		Properties: models.PowerAppPropertiesBapi{
			DisplayName: "display_name1",
			CreatedTime: "created_time1",
			Environment: models.PowerAppEnvironmentDto{
				Name: "environment",
			},
		},
	})
	apps = append(apps, models.PowerAppBapi{
		Name: "name2",
		Properties: models.PowerAppPropertiesBapi{
			DisplayName: "display_name2",
			CreatedTime: "created_time2",
			Environment: models.PowerAppEnvironmentDto{
				Name: "environment",
			},
		},
	})
	clientMock.EXPECT().GetPowerApps(gomock.Any(), gomock.Any()).Return(apps, nil).AnyTimes()

	app1 := ConvertFromPowerAppDto(apps[0])
	app2 := ConvertFromPowerAppDto(apps[1])

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"powerplatform": powerPlatformProviderServerApiMock(clientMock, nil, nil),
		},
		Steps: []resource.TestStep{
			{
				Config: uniTestsProviderConfig + `
				data "powerplatform_powerapps" "all" {}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_powerapps.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					resource.TestCheckResourceAttr("data.powerplatform_powerapps.all", "powerapps.#", strconv.Itoa(len(apps))),

					resource.TestCheckResourceAttr("data.powerplatform_powerapps.all", "powerapps.0.name", app1.Name.ValueString()),
					resource.TestCheckResourceAttr("data.powerplatform_powerapps.all", "powerapps.0.environment_name", app1.EnvironmentName.ValueString()),
					resource.TestCheckResourceAttr("data.powerplatform_powerapps.all", "powerapps.0.display_name", app1.DisplayName.ValueString()),
					resource.TestCheckResourceAttr("data.powerplatform_powerapps.all", "powerapps.0.created_time", app1.CreatedTime.ValueString()),

					resource.TestCheckResourceAttr("data.powerplatform_powerapps.all", "powerapps.1.name", app2.Name.ValueString()),
					resource.TestCheckResourceAttr("data.powerplatform_powerapps.all", "powerapps.1.environment_name", app2.EnvironmentName.ValueString()),
					resource.TestCheckResourceAttr("data.powerplatform_powerapps.all", "powerapps.1.display_name", app2.DisplayName.ValueString()),
					resource.TestCheckResourceAttr("data.powerplatform_powerapps.all", "powerapps.1.created_time", app2.CreatedTime.ValueString()),
				),
			},
		},
	})
}
