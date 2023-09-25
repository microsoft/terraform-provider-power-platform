package powerplatform

import (
	"context"
	"regexp"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	mocks "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/models"
)

func TestAccEnvironmentsResource_Validate_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "Example2"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "USD"
					environment_type                          = "Sandbox"
					security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					domain									  = "terraformtest2"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),

					// Verify the first power app to ensure all attributes are set
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "Example2"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", "terraformtest2"),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "environment_name", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Sandbox"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "language_code", "1033"),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "organization_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "security_group_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "url", "https://terraformtest2.crm4.dynamics.com/"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
				),
			},
			{
				Config: providerConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "Example3"
					domain									  = "terraformtest3"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "USD"
					environment_type                          = "Sandbox"
					security_group_id 						  = "00000000-0000-0000-0000-000000000000"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "Example3"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", "terraformtest3"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "url", "https://terraformtest3.crm4.dynamics.com/"),
				),
			},
		},
	})
}

func TestAccEnvironmentsResource_Validate_Create(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "Example1"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "USD"
					environment_type                          = "Sandbox"
					security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					domain									  = "terraformtest1"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),

					// Verify the first power app to ensure all attributes are set
					resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", "terraformtest1"),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "environment_name", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "environment_type", regexp.MustCompile(`^(Default|Sandbox|Developer)$`)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "language_code", "1033"),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "organization_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "security_group_id", regexp.MustCompile(powerplatform_helpers.GuidOrEmptyValueRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "url", "https://terraformtest1.crm4.dynamics.com/"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "Example1"),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "version", regexp.MustCompile(powerplatform_helpers.VersionRegex)),
				),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Create_And_Force_Recreate(t *testing.T) {
	clientMock := mocks.NewUnitTestsMockBapiClientInterface(t)

	envIdBeforeChanges := "00000000-0000-0000-0000-000000000001"
	envIdAfterLocationChanges := "00000000-0000-0000-0000-000000000002"
	envIdAfterCurrencyChanges := "00000000-0000-0000-0000-000000000003"
	envIdAfterLanguageChanges := "00000000-0000-0000-0000-000000000004"
	envIdAfterEnvironmentTypeChanges := "00000000-0000-0000-0000-000000000005"

	env := models.EnvironmentDto{
		Name: envIdBeforeChanges,
		Properties: models.EnvironmentPropertiesDto{
			LinkedEnvironmentMetadata: models.LinkedEnvironmentMetadataDto{
				ResourceId:      "org1",
				SecurityGroupId: "security1",
				DomainName:      "domain",
				InstanceURL:     "url",
				Version:         "version",
			},
		},
	}

	steps := []resource.TestStep{
		{
			Config: uniTestsProviderConfig + `
			resource "powerplatform_environment" "development" {
				display_name                              = "Example1"
				location                                  = "europe"
				language_code                             = "1033"
				currency_code                             = "USD"
				environment_type                          = "Sandbox"
				domain									  = "domain"
				security_group_id 						  = "security1"

			}`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_name", envIdBeforeChanges),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
			),
		},
		{
			Config: uniTestsProviderConfig + `
			resource "powerplatform_environment" "development" {
				display_name                              = "Example1"
				location                                  = "unitedstates"
				language_code                             = "1033"
				currency_code                             = "USD"
				environment_type                          = "Sandbox"
				domain									  = "domain"
				security_group_id 						  = "security1"
			}`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_name", envIdAfterLocationChanges),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "unitedstates"),
			),
		},
		{
			Config: uniTestsProviderConfig + `
			resource "powerplatform_environment" "development" {
				display_name                              = "Example1"
				location                                  = "unitedstates"
				language_code                             = "1033"
				currency_code                             = "EUR"
				environment_type                          = "Sandbox"
				domain									  = "domain"
				security_group_id 						  = "security1"
			}`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_name", envIdAfterCurrencyChanges),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "currency_code", "EUR"),
			),
		},
		{
			Config: uniTestsProviderConfig + `
			resource "powerplatform_environment" "development" {
				display_name                              = "Example1"
				location                                  = "unitedstates"
				language_code                             = "1033"
				currency_code                             = "EUR"
				environment_type                          = "Trial"
				domain									  = "domain"
				security_group_id 						  = "security1"
			}`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_name", envIdAfterEnvironmentTypeChanges),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Trial"),
			),
		},
		{
			Config: uniTestsProviderConfig + `
			resource "powerplatform_environment" "development" {
				display_name                              = "Example1"
				location                                  = "europe"
				language_code                             = "1031"
				currency_code                             = "EUR"
				environment_type                          = "Sandbox"
				domain									  = "domain"
				security_group_id 						  = "security1"
			}`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_name", envIdAfterLanguageChanges),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "language_code", "1031"),
			),
		},
	}

	clientMock.EXPECT().GetEnvironment(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, id string) (*models.EnvironmentDto, error) {
		return &env, nil
	}).AnyTimes()

	clientMock.EXPECT().CreateEnvironment(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, envToCreate models.EnvironmentCreateDto) (*models.EnvironmentDto, error) {

		env = models.EnvironmentDto{
			Name:     envIdBeforeChanges,
			Id:       envIdBeforeChanges,
			Location: envToCreate.Location,
			Properties: models.EnvironmentPropertiesDto{
				DisplayName:    envToCreate.Properties.DisplayName,
				EnvironmentSku: envToCreate.Properties.EnvironmentSku,
				LinkedEnvironmentMetadata: models.LinkedEnvironmentMetadataDto{
					DomainName:      envToCreate.Properties.LinkedEnvironmentMetadata.DomainName,
					BaseLanguage:    envToCreate.Properties.LinkedEnvironmentMetadata.BaseLanguage,
					SecurityGroupId: envToCreate.Properties.LinkedEnvironmentMetadata.SecurityGroupId,
				},
			},
		}

		if envToCreate.Location == "unitedstates" {
			env.Name = envIdAfterLocationChanges
			env.Id = envIdAfterLocationChanges
		}
		if envToCreate.Properties.LinkedEnvironmentMetadata.Currency.Code == "EUR" {
			env.Name = envIdAfterCurrencyChanges
			env.Id = envIdAfterCurrencyChanges
		}
		if envToCreate.Properties.EnvironmentSku == "Trial" {
			env.Name = envIdAfterEnvironmentTypeChanges
			env.Id = envIdAfterEnvironmentTypeChanges
		}
		if envToCreate.Properties.LinkedEnvironmentMetadata.BaseLanguage == 1031 {
			env.Name = envIdAfterLanguageChanges
			env.Id = envIdAfterLanguageChanges
		}

		return &env, nil
		//we expect create to be called twice because we are forcing a recreate
	}).Times(len(steps))

	clientMock.EXPECT().UpdateEnvironment(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, id string, envToUpdate models.EnvironmentDto) error {
		return nil
		//we expect update to be never called as we are forcing a recreate
	}).Times(0)

	clientMock.EXPECT().DeleteEnvironment(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, id string) error {
		return nil
	}).AnyTimes()

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"powerplatform": powerPlatformProviderServerApiMock(clientMock, nil, nil),
		},
		Steps: steps,
	})

}

func TestUnitEnvironmentsResource_Validate_Create_And_Update(t *testing.T) {
	clientMock := mocks.NewUnitTestsMockBapiClientInterface(t)

	envId := "00000000-0000-0000-0000-000000000001"
	env := models.EnvironmentDto{
		Name: envId,
		Properties: models.EnvironmentPropertiesDto{
			EnvironmentSku: "Sandbox",
			LinkedEnvironmentMetadata: models.LinkedEnvironmentMetadataDto{
				ResourceId:      "org1",
				SecurityGroupId: "security1",
				DomainName:      "domain",
				InstanceURL:     "url",
				Version:         "version",
			},
		},
	}

	steps := []resource.TestStep{
		{
			Config: uniTestsProviderConfig + `
			resource "powerplatform_environment" "development" {
				display_name                              = "Example1"
				location                                  = "europe"
				language_code                             = "1033"
				currency_code                             = "USD"
				environment_type                          = "Sandbox"
				domain									  = "domain"
				security_group_id 						  = "security1"
			}`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_name", envId),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "Example1"),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", "domain"),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "security_group_id", "security1"),
			),
		},
		{
			Config: uniTestsProviderConfig + `
			resource "powerplatform_environment" "development" {
				display_name                              = "Example123"
				location                                  = "europe"
				language_code                             = "1033"
				currency_code                             = "USD"
				environment_type                          = "Sandbox"
				domain									  = "domain"
				security_group_id 						  = "security1"
			}`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_name", envId),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "Example123"),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", "domain"),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "security_group_id", "security1"),
			),
		},
		{
			Config: uniTestsProviderConfig + `
			resource "powerplatform_environment" "development" {
				display_name                              = "Example123"
				location                                  = "europe"
				language_code                             = "1033"
				currency_code                             = "USD"
				environment_type                          = "Sandbox"
				domain									  = "domain123"
				security_group_id 						  = "security1"
			}`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_name", envId),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "Example123"),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", "domain123"),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "security_group_id", "security1"),
			),
		},
		{
			Config: uniTestsProviderConfig + `
			resource "powerplatform_environment" "development" {
				display_name                              = "Example123"
				location                                  = "europe"
				language_code                             = "1033"
				currency_code                             = "USD"
				environment_type                          = "Sandbox"
				domain									  = "domain123"
				security_group_id 						  = "security123"
			}`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_name", envId),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "Example123"),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", "domain123"),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "security_group_id", "security123"),
			),
		},
	}

	clientMock.EXPECT().GetEnvironment(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, id string) (*models.EnvironmentDto, error) {
		return &env, nil
	}).AnyTimes()

	clientMock.EXPECT().CreateEnvironment(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, envToCreate models.EnvironmentCreateDto) (*models.EnvironmentDto, error) {
		env = models.EnvironmentDto{
			Id:       envId,
			Location: envToCreate.Location,
			Name:     envId,
			Properties: models.EnvironmentPropertiesDto{
				DisplayName:    envToCreate.Properties.DisplayName,
				EnvironmentSku: env.Properties.EnvironmentSku,
				LinkedEnvironmentMetadata: models.LinkedEnvironmentMetadataDto{
					DomainName:      "domain",
					InstanceURL:     "url",
					BaseLanguage:    envToCreate.Properties.LinkedEnvironmentMetadata.BaseLanguage,
					SecurityGroupId: envToCreate.Properties.LinkedEnvironmentMetadata.SecurityGroupId,
					Version:         "version",
					ResourceId:      "org1",
				},
			},
		}
		return &env, nil
	}).Times(1)

	clientMock.EXPECT().UpdateEnvironment(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, environmentId string, environment models.EnvironmentDto) (*models.EnvironmentDto, error) {
		env.Name = environment.Name
		env.Id = environment.Id
		env.Properties.DisplayName = environment.Properties.DisplayName
		env.Properties.LinkedEnvironmentMetadata.DomainName = environment.Properties.LinkedEnvironmentMetadata.DomainName
		env.Properties.LinkedEnvironmentMetadata.SecurityGroupId = environment.Properties.LinkedEnvironmentMetadata.SecurityGroupId
		return &env, nil
	}).Times(len(steps) - 1)

	clientMock.EXPECT().DeleteEnvironment(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, id string) error {
		return nil
	}).AnyTimes()

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"powerplatform": powerPlatformProviderServerApiMock(clientMock, nil, nil),
		},
		Steps: steps,
	})

}

func TestUnitEnvironmentsResource_Validate_Create(t *testing.T) {
	clientMock := mocks.NewUnitTestsMockBapiClientInterface(t)

	env := models.EnvironmentDto{
		Name: "00000000-0000-0000-0000-000000000001",
		Properties: models.EnvironmentPropertiesDto{
			EnvironmentSku: "Sandbox",
			LinkedEnvironmentMetadata: models.LinkedEnvironmentMetadataDto{
				BaseLanguage:    1033,
				ResourceId:      "org1",
				SecurityGroupId: "security1",
				DomainName:      "domain",
				InstanceURL:     "url",
				Version:         "version",
			},
		},
	}

	clientMock.EXPECT().GetEnvironment(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, id string) (*models.EnvironmentDto, error) {
		return &env, nil
	}).Times(1)

	clientMock.EXPECT().CreateEnvironment(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, envToCreate models.EnvironmentCreateDto) (*models.EnvironmentDto, error) {
		env = models.EnvironmentDto{
			Name:     "00000000-0000-0000-0000-000000000001",
			Location: envToCreate.Location,
			Properties: models.EnvironmentPropertiesDto{
				DisplayName:    envToCreate.Properties.DisplayName,
				EnvironmentSku: "Sandbox",
				LinkedEnvironmentMetadata: models.LinkedEnvironmentMetadataDto{
					BaseLanguage:    env.Properties.LinkedEnvironmentMetadata.BaseLanguage,
					ResourceId:      "org1",
					SecurityGroupId: "security1",
					DomainName:      "domain",
					InstanceURL:     "url",
					Version:         "version",
				},
			},
		}
		return &env, nil
	}).Times(1)

	clientMock.EXPECT().DeleteEnvironment(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, id string) error {
		return nil
	}).Times(1)

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"powerplatform": powerPlatformProviderServerApiMock(clientMock, nil, nil),
		},
		Steps: []resource.TestStep{
			{
				Config: uniTestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "Example1"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "USD"
					environment_type                          = "Sandbox"
					domain                                    = "domain"
					security_group_id                         = "security1"
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "environment_name", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_name", env.Name),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "Example1"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "url", env.Properties.LinkedEnvironmentMetadata.InstanceURL),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", env.Properties.LinkedEnvironmentMetadata.DomainName),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Sandbox"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "language_code", "1033"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "currency_code", "USD"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "organization_id", env.Properties.LinkedEnvironmentMetadata.ResourceId),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "security_group_id", env.Properties.LinkedEnvironmentMetadata.SecurityGroupId),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "version", env.Properties.LinkedEnvironmentMetadata.Version),
				),
			},
		},
	})

}
