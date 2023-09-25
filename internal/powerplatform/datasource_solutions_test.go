package powerplatform

import (
	"context"
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

func TestUnitSolutionsDataSource_Validate_Read(t *testing.T) {
	clientMock := mocks.NewUnitTestMockDataverseClientInterface(t)

	envId := "00000000-0000-0000-0000-000000000001"
	sol := models.SolutionDto{
		Id:              "00000000-0000-0000-0000-000000000002",
		EnvironmentName: envId,
		DisplayName:     "Solution",
		Name:            "solution",
		CreatedTime:     "2020-01-01T00:00:00Z",
		ModifiedTime:    "2020-01-01T00:00:00Z",
		InstallTime:     "2020-01-01T00:00:00Z",
		Version:         "1.2.3.4",
		IsManaged:       true,
	}
	solutions := make([]models.SolutionDto, 0)
	solutions = append(solutions, sol)

	clientMock.EXPECT().GetSolutions(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, id string) ([]models.SolutionDto, error) {
		return solutions, nil
	}).AnyTimes()

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"powerplatform": powerPlatformProviderServerApiMock(nil, clientMock, nil),
		},
		Steps: []resource.TestStep{
			{
				Config: uniTestsProviderConfig + `
				data "powerplatform_solutions" "all" {
					environment_name = "00000000-0000-0000-0000-000000000001"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.name", sol.Name),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.environment_name", sol.EnvironmentName),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.display_name", sol.DisplayName),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.created_time", sol.CreatedTime),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.modified_time", sol.ModifiedTime),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.install_time", sol.InstallTime),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.is_managed", strconv.FormatBool(sol.IsManaged)),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.version", sol.Version),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.id", sol.Id),
				),
			},
		},
	})
}

func TestAccSolutionsDataSource_Validate_Read(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "testaccsolutionsdatasource"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					domain           = "testaccsolutionsdatasource"
					security_group_id = "00000000-0000-0000-0000-000000000000"
				}

				data "powerplatform_solutions" "all" {
					environment_name = powerplatform_environment.development.environment_name
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					// Verify the first power app to ensure all attributes are set
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.environment_name", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.display_name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.created_time", regexp.MustCompile(powerplatform_helpers.TimeRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.modified_time", regexp.MustCompile(powerplatform_helpers.TimeRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.install_time", regexp.MustCompile(powerplatform_helpers.TimeRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.is_managed", regexp.MustCompile(`^(true|false)$`)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.version", regexp.MustCompile(powerplatform_helpers.VersionRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
				),
			},
		},
	})
}
