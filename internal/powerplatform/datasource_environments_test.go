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

func TestAccEnvironmentsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				data "powerplatform_environments" "all" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					// Verify the first power app to ensure all attributes are set
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.display_name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.domain", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.environment_name", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.environment_type", regexp.MustCompile(`^(Default|Sandbox|Developer)$`)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.language_code", regexp.MustCompile(`^(1033|1031)$`)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.organization_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.security_group_id", regexp.MustCompile(powerplatform_helpers.GuidOrEmptyValueRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.url", regexp.MustCompile(powerplatform_helpers.UrlValidStringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.location", regexp.MustCompile(`^(unitedstates|europe)$`)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.version", regexp.MustCompile(powerplatform_helpers.VersionRegex)),
				),
			},
		},
	})
}

func TestUnitEnvironmentsDataSource_Validate_Read(t *testing.T) {
	clientMock := mocks.NewUnitTestsMockBapiClientInterface(t)

	envs := make([]models.EnvironmentDto, 0)
	envs = append(envs, models.EnvironmentDto{
		Location: "europe",
		Name:     "test_environment",
		Properties: models.EnvironmentPropertiesDto{
			DisplayName:    "Test Environment",
			EnvironmentSku: "Sandbox",

			LinkedEnvironmentMetadata: models.LinkedEnvironmentMetadataDto{
				DomainName:      "org1",
				InstanceURL:     "https://org1.crm.dynamics.com",
				SecurityGroupId: "00000000-0000-0000-0000-000000000000",
				Version:         "9.2.21044.00148",
				BaseLanguage:    1033,
			},
		},
	})
	envs = append(envs, models.EnvironmentDto{
		Location: "unitedstates",
		Name:     "test_environment_2",
		Properties: models.EnvironmentPropertiesDto{
			DisplayName:    "Test Environment 2",
			EnvironmentSku: "Sandbox",

			LinkedEnvironmentMetadata: models.LinkedEnvironmentMetadataDto{
				DomainName:      "org2",
				InstanceURL:     "https://org2.crm.dynamics.com",
				SecurityGroupId: "00000000-0000-0000-0000-000000000000",
				Version:         "9.2.21044.00100",
				BaseLanguage:    1031,
			},
		},
	})

	clientMock.EXPECT().GetEnvironments(gomock.Any()).Return(envs, nil).AnyTimes()

	env := ConvertFromEnvironmentDto(envs[0])

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"powerplatform": powerPlatformProviderServerApiMock(clientMock, nil, nil),
		},
		Steps: []resource.TestStep{
			{
				Config: uniTestsProviderConfig + `
				data "powerplatform_environments" "all" {}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.#", strconv.Itoa(len(envs))),

					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.display_name", env.DisplayName.ValueString()),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.domain", env.Domain.ValueString()),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.environment_name", env.EnvironmentName.ValueString()),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.environment_type", env.EnvironmentType.ValueString()),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.language_code", strconv.Itoa(int(env.LanguageName.ValueInt64()))),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.organization_id", env.OrganizationId.ValueString()),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.security_group_id", env.SecurityGroupId.ValueString()),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.url", env.Url.ValueString()),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.location", env.Location.ValueString()),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.0.version", env.Version.ValueString()),
				),
			},
		},
	})
}
