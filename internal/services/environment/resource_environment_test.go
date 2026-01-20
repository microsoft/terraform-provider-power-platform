// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitEnvironmentsResource_Validate_Attribute_Validators(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Attribute_Validators/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Attribute_Validators/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `

				variable "power_platform_environment" {
					type = object({
						id = string
						display_name = string
					    location = string
						type = string
						language_code = string
						currency_code = string
						domain = string
						security_group_id = string
					})
					default = {
						id = ""
						display_name = "displayname"
						location = "europe"
						type = "Sandbox"
						language_code = "1033"
						currency_code = "PLN"
						domain = "00000000-0000-0000-0000-000000000001"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "powerplatform_environment" "development" {
					display_name                              = var.power_platform_environment.display_name
					location                                  = var.power_platform_environment.location
					environment_type                          = var.power_platform_environment.type
					dataverse = {
						language_code                             = var.power_platform_environment.language_code
						currency_code                             = var.power_platform_environment.currency_code
						domain                                    = var.power_platform_environment.domain
						security_group_id                         = var.power_platform_environment.security_group_id
					}
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "displayname"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.url", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.domain", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Sandbox"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.language_code", "1033"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.currency_code", "PLN"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.organization_id", "orgid"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.security_group_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.version", "9.2.23092.00206"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.unique_name", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "release_cycle", "Standard"),
				),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Do_Not_Retry_On_NoCapacity(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusConflict, httpmock.File("tests/resource/Validate_Do_Not_Retry_On_NoCapacity/post_environment.json").String())
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile(".*InsufficientCapacity_StorageDriven.*"),
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						domain                                    = "00000000-0000-0000-0000-000000000001"
						security_group_id                         = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Retry_On_Running_LifecycleOperation(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	deleteRetryCount := 0
	lifecycleOperationInProgressCount := 5

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			deleteRetryCount++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Retry_On_Running_LifecycleOperation/get_lifecycle_delete_%d.json", deleteRetryCount)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Retry_On_Running_LifecycleOperation/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			if lifecycleOperationInProgressCount > 0 {
				lifecycleOperationInProgressCount--
				resp := httpmock.NewStringResponse(http.StatusConflict, httpmock.File("tests/resource/Validate_Retry_On_Running_LifecycleOperation/post_environment_operation_in_progress.json").String())
				return resp, nil
			}
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Retry_On_Running_LifecycleOperation/get_environment_%s.json", id)).String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						domain                                    = "00000000-0000-0000-0000-000000000001"
						security_group_id                         = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Retry_LifecycleOperation(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	deleteRetryCount := 0

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			deleteRetryCount++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Retry_LidecycleOperation/get_lifecycle_delete_%d.json", deleteRetryCount)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Retry_LidecycleOperation/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Retry_LidecycleOperation/get_environment_%s.json", id)).String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						domain                                    = "00000000-0000-0000-0000-000000000001"
						security_group_id                         = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccEnvironmentsResource_Validate_Update_Name_Field(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
  				resource "powerplatform_environment" "development" {
					display_name                              = "aaa"
					location                                  = "unitedstates"
					environment_type                       	  = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "USD"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "aaa1"
					location                                  = "unitedstates"
					environment_type                       	  = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "USD"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "aaa1"),
				),
			},
		},
	})
}

func TestAccEnvironmentsResource_Validate_CreateGenerativeAiFeatures_Non_US_Region_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "powerplatform_environment" "development" {
						display_name                              = "` + fmt.Sprintf("%s_%d", t.Name(), rand.Intn(100000)) + `"
						location                                  = "europe"
						environment_type                       	  = "Sandbox"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "allow_bing_search", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "allow_moving_data_across_regions", "true"),
				),
			},
			{
				Config: `
					resource "powerplatform_environment" "development" {
						display_name                              = "` + fmt.Sprintf("%s_%d", t.Name(), rand.Intn(100000)) + `"
						location                                  = "europe"
						environment_type                       	  = "Sandbox"

						allow_bing_search                = false
						allow_moving_data_across_regions = true
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "allow_bing_search", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "allow_moving_data_across_regions", "true"),
				),
			},
			{
				Config: `
					resource "powerplatform_environment" "development" {
						display_name                              = "` + fmt.Sprintf("%s_%d", t.Name(), rand.Intn(100000)) + `"
						location                                  = "europe"
						environment_type                       	  = "Sandbox"

						allow_bing_search                = false
						allow_moving_data_across_regions = false
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "allow_bing_search", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "allow_moving_data_across_regions", "false"),
				),
			},
		},
	})
}

func TestAccEnvironmentsResource_Validate_CreateGenerativeAiFeatures_US_Region_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "powerplatform_environment" "development" {
						display_name                              = "` + fmt.Sprintf("%s_%d", t.Name(), rand.Intn(100000)) + `"
						location                                  = "unitedstates"
						environment_type                       	  = "Sandbox"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "allow_bing_search", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "allow_moving_data_across_regions", "false"),
				),
			},
			{
				Config: `
					resource "powerplatform_environment" "development" {
						display_name                              = "` + fmt.Sprintf("%s_%d", t.Name(), rand.Intn(100000)) + `"
						location                                  = "unitedstates"
						environment_type                       	  = "Sandbox"

						allow_bing_search                = true
						//on usa region, moving data across regions is not allowed and always false
						allow_moving_data_across_regions = false
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "allow_bing_search", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "allow_moving_data_across_regions", "false"),
				),
			},
			{
				Config: `
					resource "powerplatform_environment" "development" {
						display_name                              = "` + fmt.Sprintf("%s_%d", t.Name(), rand.Intn(100000)) + `"
						location                                  = "unitedstates"
						environment_type                       	  = "Sandbox"

						allow_bing_search                = false
						allow_moving_data_across_regions = false
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "allow_bing_search", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "allow_moving_data_across_regions", "false"),
				),
			},
		},
	})
}

func TestAccEnvironmentsResource_Validate_CreateGenerativeAiFeatures_US_Region_Expect_Fail(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile(".*moving data across regions is not supported in the unitedstates location.*"),
				Config: `
					resource "powerplatform_environment" "development" {
						display_name                              = "` + fmt.Sprintf("%s_%d", t.Name(), rand.Intn(100000)) + `"
						location                                  = "unitedstates"
						environment_type                       	  = "Sandbox"

						allow_bing_search                = false
						allow_moving_data_across_regions = true
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestAccEnvironmentsResource_Validate_CreateGenerativeAiFeatures_Non_US_Region_Expect_Fail(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile(".*to enable ai generative features, moving data across regions must be enabled.*"),
				Config: `
					resource "powerplatform_environment" "development" {
						display_name                              = "` + fmt.Sprintf("%s_%d", t.Name(), rand.Intn(100000)) + `"
						location                                  = "europe"
						environment_type                       	  = "Sandbox"

						allow_bing_search                = true
						allow_moving_data_across_regions = false
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestAccEnvironmentsResource_Validate_CreateDeveloperEnvironment(t *testing.T) {
	t.Skip("creating dev environments with SP is NOT yet supported")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azuread": {
				VersionConstraint: constants.AZURE_AD_PROVIDER_VERSION_CONSTRAINT,
				Source:            "hashicorp/azuread",
			},
			"random": {
				VersionConstraint: constants.RANDOM_PROVIDER_VERSION_CONSTRAINT,
				Source:            "hashicorp/random",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
				data "azuread_domains" "aad_domains" {
					only_initial = true
				}

				locals {
					domain_name = data.azuread_domains.aad_domains.domains[0].domain_name
				}

				resource "random_password" "passwords" {
				    min_lower = 1
					min_upper        = 1
					min_numeric      = 1
					min_special      = 1
					length           = 16
					special          = true
					override_special = "_%@"
				}

				resource "azuread_user" "test_user" {
					user_principal_name = "` + mocks.TestName() + `@${local.domain_name}"
					display_name        = "` + mocks.TestName() + `"
					mail_nickname       = "` + mocks.TestName() + `"
					password            = random_password.passwords.result
					usage_location      = "US"
				}

				resource "powerplatform_environment" "development" {
					display_name         = "` + mocks.TestName() + `"
					location             = "europe"
					environment_type     = "Developer"
					owner_id = azuread_user.test_user.id
					dataverse = {
						language_code     = "1033"
						currency_code     = "USD"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Create_Error_Check_Environment_Group(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile(".*Attribute \"dataverse\" must be specified when \"environment_group_id\".*"),
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "displayname"
					location         = "europe"
					environment_type = "Sandbox"
					environment_group_id = "00000000-0000-0000-0000-000000000001"
				}`,
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_CreateDevelopmentEnvironment_Error_Check_Security_Group(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile(".*Attribute \"dataverse.security_group_id\" cannot be specified when \"owner_id\".*"),
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "displayname"
					location         = "europe"
					environment_type = "Developer"
					owner_id = "00000000-0000-0000-0000-000000000001"
					dataverse = {
						language_code = "1033"
						currency_code = "PLN"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}`,
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_CreateDevelopmentEnvironment_Error_Check_No_Dataverse(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile(".*Attribute \"dataverse\" must be specified when \"owner_id\" is specified.*"),
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "displayname"
					location         = "europe"
					environment_type = "Developer"
					owner_id = "00000000-0000-0000-0000-000000000001"
				}`,
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_CreateDevelopmentEnvironment_Error_Check_No_Developer_Env(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile(".*owner_id can be used only when environment_type is `Developer`.*"),
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "displayname"
					location         = "europe"
					environment_type = "Sandbox"
					owner_id = "00000000-0000-0000-0000-000000000001"
					dataverse = {
						language_code = "1033"
						currency_code = "PLN"
					}
				}`,
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_CreateDevelopmentEnvironment_Error_Check_No_OwnerId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile(".*owner_id must be set when environment_type is `Developer`.*"),
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "displayname"
					location         = "europe"
					environment_type = "Developer"
					dataverse = {
						language_code = "1033"
						currency_code = "PLN"
					}
				}`,
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_CreateDevelopmentEnvironment_Error_Check_Empty_OwnerId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile(".*owner_id must be set when environment_type is `Developer`.*"),
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "displayname"
					location         = "europe"
					environment_type = "Developer"
					owner_id = ""
					dataverse = {
						language_code = "1033"
						currency_code = "PLN"
					}
				}`,
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_CreateDeveloperEnvironment(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Dev_Env/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Create_Dev_Env/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Dev_Env/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "displayname"
					location         = "europe"
					azure_region     = "westeurope"
					environment_type = "Developer"
					cadence          = "Frequent"
					owner_id = "00000000-0000-0000-0000-000000000001"
					dataverse = {
						language_code = "1033"
						currency_code = "PLN"
					}
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "displayname"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.url", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.domain", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Developer"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "owner_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.language_code", "1033"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.currency_code", "PLN"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.organization_id", "orgid"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.security_group_id"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.version", "9.2.23092.00206"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.unique_name", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "release_cycle", "Standard"),
				),
			},
		},
	})
}

func TestAccEnvironmentsResource_Validate_Create_Early_Release_Cycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					release_cycle                             = "Early"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "release_cycle", "Early"),
				),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Create_Early_Release_Cycle(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Early_Release_Cycle/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Create_Early_Release_Cycle/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Early_Release_Cycle/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "powerplatform_environment" "development" {
						display_name                              = "displayname"
						location                                  = "europe"
						environment_type                          = "Sandbox"
						release_cycle                             = "Early"
						dataverse = {
							language_code                             = "1033"
							currency_code                             = "PLN"
							domain                                    = "00000000-0000-0000-0000-000000000001"
							security_group_id                         = "00000000-0000-0000-0000-000000000000"
						}
					}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "release_cycle", "Early"),
				),
			},
		},
	})
}

func TestAccEnvironmentsResource_Validate_Update(t *testing.T) {
	domainName := fmt.Sprintf("terraformprovidertest%d", rand.Intn(100000))
	newDomainName := fmt.Sprintf("terraformprovidertest%d", rand.Intn(100000))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					description                               = "aaaa"
					cadence								   	  = "Moderate"
					location                                  = "unitedstates"
					environment_type                       	  = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "USD"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
						domain									  =  "` + domainName + `"
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify placeholder id attribute.
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),

					// Verify the first power app to ensure all attributes are set.
					resource.TestCheckResourceAttr("powerplatform_environment.development", "description", "aaaa"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "cadence", "Moderate"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", mocks.TestName()),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.domain", domainName),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Sandbox"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.language_code", "1033"),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "dataverse.organization_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.security_group_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.url", "https://"+domainName+".crm.dynamics.com/"),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "dataverse.unique_name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "unitedstates"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000000"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					description                               = "bbbb"
					cadence								      = "Frequent"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					dataverse = {
						domain									  =  "` + newDomainName + `"
						language_code                             = "1033"
						currency_code                             = "USD"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "description", "bbbb"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "cadence", "Frequent"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.domain", newDomainName),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.url", "https://"+newDomainName+".crm.dynamics.com/"),
				),
			},
		},
	})
}

func TestAccEnvironmentsResource_Validate_Domain_Uniqueness_On_Update(t *testing.T) {
	domainName := fmt.Sprintf("terraformprovidertest%d", rand.Intn(100000))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					description                               = "desc"
					cadence								   	  = "Moderate"
					location                                  = "unitedstates"
					environment_type                       	  = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "USD"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
						domain									  =  "` + domainName + `"
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "description", "desc"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.domain", domainName),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.security_group_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.url", "https://"+domainName+".crm.dynamics.com/"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					description                               = "desc test"
					cadence								   	  = "Moderate"
					location                                  = "unitedstates"
					environment_type                       	  = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "USD"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
						domain									  =  "` + domainName + `"
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "description", "desc test"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.domain", domainName),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.security_group_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.url", "https://"+domainName+".crm.dynamics.com/"),
				),
			},
		},
	})
}

func TestAccEnvironmentsResource_Validate_Create(t *testing.T) {
	domainName := fmt.Sprintf("orgtest%d", rand.Intn(100000))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "USD"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
						domain									  =  "` + domainName + `"
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),

					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "environment_type", regexp.MustCompile(`^(Default|Sandbox|Developer)$`)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.language_code", "1033"),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "dataverse.organization_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "dataverse.security_group_id", regexp.MustCompile(helpers.GuidOrEmptyValueRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.domain", domainName),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.url", "https://"+domainName+".crm.dynamics.com/"),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "dataverse.unique_name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "unitedstates"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", mocks.TestName()),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "dataverse.version", regexp.MustCompile(helpers.VersionRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "release_cycle", "Standard"),
					// resource.TestMatchResourceAttr("powerplatform_environment.development", "templates", regexp.MustCompile(`D365_FinOps_Finance$`)),
					// resource.TestMatchResourceAttr("powerplatform_environment.development", "template_metadata", regexp.MustCompile(`{"PostProvisioningPackages": [{ "applicationUniqueName": "msdyn_FinanceAndOperationsProvisioningAppAnchor",\n "parameters": "DevToolsEnabled=true\|DemoDataEnabled=true"\n }\n ]\n }`)),
					// resource.TestMatchResourceAttr("powerplatform_environment.development", "linked_app_url", regexp.MustCompile(`\.operations\.dynamics\.com$`)),
				),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Create_And_Force_Recreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	envIdResponseInx := -1
	envIdResponseArray := []string{"00000000-0000-0000-0000-000000000001",
		"00000000-0000-0000-0000-000000000002",
		"00000000-0000-0000-0000-000000000003",
		"00000000-0000-0000-0000-000000000004"}

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Force_Recreate/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://([\d-]+)\.crm4\.dynamics\.com/api/data/v9\.2/transactioncurrencies\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Create_And_Force_Recreate/get_transactioncurrencies_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://([\d-]+)\.crm4\.dynamics\.com/api/data/v9\.2/organizations\z`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"value": [
					{
						"_basecurrencyid_value": "xyz"
					}]}`), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			envIdResponseInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Create_And_Force_Recreate/get_lifecycle_%s.json", envIdResponseArray[envIdResponseInx])).String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Create_And_Force_Recreate/get_environment_%s.json", id)).String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						domain									  = "00000000-0000-0000-0000-000000000001"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.currency_code", "PLN"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.unique_name", "00000000-0000-0000-0000-000000000001"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						domain									  = "00000000-0000-0000-0000-000000000002"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "unitedstates"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.currency_code", "PLN"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "Example1"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "EUR"
						domain									  = "00000000-0000-0000-0000-000000000003"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000003"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.currency_code", "EUR"),
				),
			},
			{
				Config: `
					resource "powerplatform_environment" "development" {
					display_name                              = "Example1"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1031"
						currency_code                             = "EUR"
						domain									  = "00000000-0000-0000-0000-000000000004"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000004"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.language_code", "1031"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.currency_code", "EUR"),
				),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Create_And_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	getLifecycleResponseInx := 0
	patchResponseInx := 0

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Create_And_Update/get_lifecycle_%d.json", getLifecycleResponseInx)).String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			getLifecycleResponseInx++
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Create_And_Update/get_environment_%d.json", patchResponseInx)).String()), nil
		})

	httpmock.RegisterResponder("PATCH", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			patchResponseInx++
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24expand=properties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Create_And_Update/get_environments_%d.json", patchResponseInx)).String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "Example1"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						domain									  = "00000000-0000-0000-0000-000000000001"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "Example1"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.domain", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.security_group_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.unique_name", "00000000-0000-0000-0000-000000000001"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "Example123"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						domain									  = "00000000-0000-0000-0000-000000000001"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "Example123"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.domain", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.security_group_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.unique_name", "00000000-0000-0000-0000-000000000001"),
				),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Create/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						domain                                    = "00000000-0000-0000-0000-000000000001"
						security_group_id                         = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "displayname"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.url", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.domain", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Sandbox"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.language_code", "1033"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.currency_code", "PLN"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.organization_id", "orgid"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.security_group_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.version", "9.2.23092.00206"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.unique_name", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "release_cycle", "Standard"),
				),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Create_With_Billing_Policy(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_With_Billing_Policy/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Create_With_Billing_Policy/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_With_Billing_Policy/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					billing_policy_id                         = "00000000-0000-0000-0000-000000000002"
					environment_type                          = "Sandbox"
					dataverse = {
						currency_code                             = "PLN"
						language_code                             = "1033"
						domain                                    = "00000000-0000-0000-0000-000000000001"
						security_group_id                         = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000002"),
				),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Update_With_Billing_Policy(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	getResponseInx := 0
	patchResponseInx := 0

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Update_With_Billing_Policy/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			getResponseInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Update_With_Billing_Policy/get_environment_%d.json", getResponseInx)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Update_With_Billing_Policy/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001/environments/add?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000002/environments/add?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001/environments/remove?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000002/environments/remove?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	httpmock.RegisterResponder("PATCH", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("PATCH", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24expand=properties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			patchResponseInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Update_With_Billing_Policy/get_environments_%d.json", patchResponseInx)).String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						domain                                    = "00000000-0000-0000-0000-000000000001"
						security_group_id                         = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000000"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					billing_policy_id                         = "00000000-0000-0000-0000-000000000001"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						domain                                    = "00000000-0000-0000-0000-000000000001"
						security_group_id                         = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000001"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					billing_policy_id                         = "00000000-0000-0000-0000-000000000002"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						domain                                    = "00000000-0000-0000-0000-000000000001"
						security_group_id                         = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000002"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					billing_policy_id                         = "00000000-0000-0000-0000-000000000000"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						domain                                    = "00000000-0000-0000-0000-000000000001"
						security_group_id                         = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000000"),
				),
			},
		},
	})
}
func TestUnitEnvironmentsResource_Validate_Create_With_D365_Template(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_With_D365_Template/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Create_With_D365_Template/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_With_D365_Template/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("PATCH", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusAccepted, ""), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24expand=properties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_With_D365_Template/get_environments.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						domain                                    = "00000000-0000-0000-0000-000000000001"
						security_group_id                         = "00000000-0000-0000-0000-000000000000"
						templates = ["D365_FinOps_Finance"]
						template_metadata = "{\"PostProvisioningPackages\":[{\"applicationUniqueName\":\"msdyn_FinanceAndOperationsProvisioningAppAnchor\",\"parameters\":\"DevToolsEnabled=true|DemoDataEnabled=true\"}]}"
					}
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttr("powerplatform_environment.development", "dataverse.templates.*", "D365_FinOps_Finance"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.template_metadata", "{\"PostProvisioningPackages\":[{\"applicationUniqueName\":\"msdyn_FinanceAndOperationsProvisioningAppAnchor\",\"parameters\":\"DevToolsEnabled=true|DemoDataEnabled=true\"}]}"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						domain                                    = "00000000-0000-0000-0000-000000000001"
						security_group_id                         = "00000000-0000-0000-0000-000000000000"
						templates = ["D365_FinOps_Finance"]
						template_metadata = "{\"PostProvisioningPackages\":[{\"applicationUniqueName\":\"msdyn_FinanceAndOperationsProvisioningAppAnchor\",\"parameters\":\"DevToolsEnabled=true|DemoDataEnabled=true\"}]}"
					}
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttr("powerplatform_environment.development", "dataverse.templates.*", "D365_FinOps_Finance"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.template_metadata", "{\"PostProvisioningPackages\":[{\"applicationUniqueName\":\"msdyn_FinanceAndOperationsProvisioningAppAnchor\",\"parameters\":\"DevToolsEnabled=true|DemoDataEnabled=true\"}]}"),
				),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Taken_Domain_Name(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/validateEnvironmentDetails?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusBadRequest, `{
				"error": {
					"code": "InvalidDomainName",
					"message": "The specified domain name with a value of 'wrong domain name' is invalid. A domain name must start with a letter and contain only characters, A-Z, a-z, 0-9 and '-'."
				}
			}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile(".*InvalidDomainName.*"),
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						domain                                    = "wrong domain name"
						security_group_id                         = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Domain_Format_Valid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PlanOnly: true,
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "displayname"
					location         = "europe"
					environment_type = "Sandbox"
					dataverse = {
						language_code     = "1033"
						currency_code     = "PLN"
						domain            = "example-env-2026"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}`,
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Domain_Format_Invalid_Characters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile("domain must start with a lowercase letter"),
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "displayname"
					location         = "europe"
					environment_type = "Sandbox"
					dataverse = {
						language_code     = "1033"
						currency_code     = "PLN"
						domain            = "example_env"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}`,
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Create_No_Dataverse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Dataverse/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Create_No_Dataverse/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Dataverse/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					description                               = "description"
					cadence								      = "Moderate"
					location                                  = "europe"
					environment_type                          = "Sandbox"
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "displayname"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "description", "description"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "cadence", "Moderate"),

					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Sandbox"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000000"),

					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.url"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.domain"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.language_code"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.currency_code"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.organization_id"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.security_group_id"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.version"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.unique_name"),
				),
			},
		},
	})
}

func TestAccEnvironmentsResource_Validate_Create_No_Dataverse(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "europe"
					environment_type                          = "Sandbox"
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", mocks.TestName()),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Sandbox"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000000"),

					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.language_code"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.currency_code"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.organization_id"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.security_group_id"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.version"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.unique_name"),
				),
			},
		},
	})
}

func TestAccEnvironmentsResource_Validate_Create_Them_Try_Remove_Dataverse(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						security_group_id                         = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "unitedstates"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Sandbox"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.language_code", "1033"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.currency_code", "PLN"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.security_group_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000000"),
				),
			},
			{
				Config: `
					resource "powerplatform_environment" "development" {
						display_name                              = "` + mocks.TestName() + `"
						location                                  = "unitedstates"
						environment_type                          = "Sandbox"
					}`,
				Check: resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Create_Them_Try_Remove_Dataverse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Them_Try_Remove_Dataverse/get_lifecycle_delete.json").String()), nil
		})

	envQueryInx := 0
	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			envQueryInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Create_Them_Try_Remove_Dataverse/get_environment_%d.json", envQueryInx)).String()), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Them_Try_Remove_Dataverse/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						domain                                    = "00000000-0000-0000-0000-000000000001"
						security_group_id                         = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "displayname"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.url", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.domain", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Sandbox"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.language_code", "1033"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.currency_code", "PLN"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.organization_id", "orgid"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.security_group_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.version", "9.2.23092.00206"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000000"),
				),
			},
			{
				Config: `
					resource "powerplatform_environment" "development" {
						display_name                              = "displayname"
						location                                  = "europe"
						environment_type                          = "Sandbox"
					}`,
				Check: resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Create_Environment_And_Dataverse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Environment_And_Dataverse/get_lifecycle_delete.json").String()), nil
		})

	envQueryInx := 0
	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			envQueryInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Create_Environment_And_Dataverse/get_environment_%d.json", envQueryInx)).String()), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Environment_And_Dataverse/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/provisionInstance?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000002?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000002?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Environment_And_Dataverse/get_lifecycle_new_dataverse.json").String()), nil
		})

	httpmock.RegisterResponder("PATCH", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2021-04-01`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("PATCH", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2021-04-01`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24expand=properties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Environment_And_Dataverse/get_environments.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "displayname"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Sandbox"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000000"),

					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.url"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.domain"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.language_code"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.currency_code"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.organization_id"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.security_group_id"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.version"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.unique_name"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						domain                                    = "00000000-0000-0000-0000-000000000001"
						security_group_id                         = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "displayname"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.url", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.domain", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Sandbox"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.language_code", "1033"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.currency_code", "PLN"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.organization_id", "orgid"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.security_group_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.version", "9.2.23092.00206"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.unique_name", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000000"),
				),
			},
		},
	})
}

func TestAccEnvironmentsResource_Validate_Create_Environment_And_Dataverse(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", mocks.TestName()),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "unitedstates"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Sandbox"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000000"),

					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.url"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.domain"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.language_code"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.currency_code"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.organization_id"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.security_group_id"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.version"),
					resource.TestCheckNoResourceAttr("powerplatform_environment.development", "dataverse.unique_name"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						security_group_id                         = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", mocks.TestName()),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "unitedstates"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Sandbox"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.language_code", "1033"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.currency_code", "PLN"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.security_group_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000000"),
				),
			},
		},
	})
}

func TestAccEnvironmentsResource_Validate_Locations_And_Azure_Regions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               false,
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "TestAccEnvironmentsResource_Validate_Locations_And_Azure_Regions"
					location                                  = "foo"
					environment_type                          = "Sandbox"
				}`,
				ExpectError: regexp.MustCompile(".*location 'foo' is not valid.*"),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "TestAccEnvironmentsResource_Validate_Locations_And_Azure_Regions"
					location                                  = "europe"
					azure_region 							  = "bar"
					environment_type                          = "Sandbox"
				}`,
				ExpectError: regexp.MustCompile(".*region 'bar' is not valid for location europe.*"),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "TestAccEnvironmentsResource_Validate_Locations_And_Azure_Regions"
					location                                  = "europe"
					azure_region 							  = "westeurope"
					environment_type                          = "Sandbox"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "azure_region", "westeurope"),
				),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Locations_And_Azure_Regions(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Locations_And_Azure_Regions/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Locations_And_Azure_Regions/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Locations_And_Azure_Regions/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "foo"
					environment_type                          = "Sandbox"
				}`,
				ExpectError: regexp.MustCompile(".*location 'foo' is not valid.*"),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					azure_region 							  = "bar"
					environment_type                          = "Sandbox"
				}`,
				ExpectError: regexp.MustCompile(".*region 'bar' is not valid for location europe.*"),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					azure_region 							  = "westeurope"
					environment_type                          = "Sandbox"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "azure_region", "westeurope"),
				),
			},
		},
	})
}

func TestAccEnvironmentsResource_Validate_Enable_Admin_Mode(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "USD"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),

					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "environment_type", regexp.MustCompile(`^(Default|Sandbox|Developer)$`)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.language_code", "1033"),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "dataverse.organization_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "dataverse.security_group_id", regexp.MustCompile(helpers.GuidOrEmptyValueRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "unitedstates"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", mocks.TestName()),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "dataverse.version", regexp.MustCompile(helpers.VersionRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000000"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "USD"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
						administration_mode_enabled 			 = true
						background_operation_enabled		     = false
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),

					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.administration_mode_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.background_operation_enabled", "false"),
				),
			},
			{
				Config: `
					resource "powerplatform_environment" "development" {
						display_name                              = "` + mocks.TestName() + `"
						location                                  = "unitedstates"
						environment_type                          = "Sandbox"
						dataverse = {
							language_code                             = "1033"
							currency_code                             = "USD"
							security_group_id 						  = "00000000-0000-0000-0000-000000000000"
							administration_mode_enabled 			 = true
							background_operation_enabled		     = true
						}
					}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),

					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.administration_mode_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.background_operation_enabled", "true"),
				),
			},
			{
				Config: `
					resource "powerplatform_environment" "development" {
						display_name                              = "` + mocks.TestName() + `"
						location                                  = "unitedstates"
						environment_type                          = "Sandbox"
						dataverse = {
							language_code                             = "1033"
							currency_code                             = "USD"
							security_group_id 						  = "00000000-0000-0000-0000-000000000000"
							administration_mode_enabled 			 = false
							background_operation_enabled		     = true
						}
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),

					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.administration_mode_enabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.background_operation_enabled", "true"),
				),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Create_Environment_With_Env_Group(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24filter=properties%2FparentEnvironmentGroup%2Fid+eq+00000000-0000-0000-0000-000000000001&api-version=2021-04-01",
		httpmock.NewStringResponder(http.StatusOK, `{"value":[]}`),
	)

	httpmock.RegisterResponder("GET", "https://000000000000000000000000000000.01.tenant.api.powerplatform.com/governance/environmentGroups/00000000-0000-0000-0000-000000000001/ruleSets?api-version=2021-10-01-preview",
		httpmock.NewStringResponder(http.StatusOK, `{"value": [{"parameters": [],"id": "00000000-0000-0000-0000-000000000001","environmentFilter": {"type": "Include","values": []}}]}`),
	)

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Create_Environment_With_Env_Group/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Create_Environment_With_Env_Group/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Create_Environment_With_Env_Group/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environmentGroups?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/Create_Environment_With_Env_Group/get_environment_group.json").String())
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environmentGroups/00000000-0000-0000-0000-000000000001?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Create_Environment_With_Env_Group/get_environment_group.json").String())
			return resp, nil
		},
	)

	httpmock.RegisterResponder("DELETE", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environmentGroups/00000000-0000-0000-0000-000000000001?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			return resp, nil
		},
	)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_group" "env_group" {
					display_name                              = "test_env_group"
					description                               = "test env group"
				}

				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					environment_group_id					  = powerplatform_environment_group.env_group.id
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_group_id", "00000000-0000-0000-0000-000000000001"),
				),
			},
		},
	})
}

func TestAccEnvironmentsResource_Create_Environment_With_Env_Group(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_group" "env_group" {
					display_name                              = "` + mocks.TestName() + `"
					description                               = "` + mocks.TestName() + `"
				}

				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					environment_group_id					  = powerplatform_environment_group.env_group.id
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "USD"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "environment_group_id", regexp.MustCompile(helpers.GuidRegex)),
				),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Create_Environment_And_Add_Env_Group(t *testing.T) {
	postEnvGroupRequestInx := 0
	getEnvironmentRequestInx := 0

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24filter=properties%2FparentEnvironmentGroup%2Fid+eq+00000000-0000-0000-0000-000000000001&api-version=2021-04-01",
		httpmock.NewStringResponder(http.StatusOK, `{"value":[]}`),
	)

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24filter=properties%2FparentEnvironmentGroup%2Fid+eq+00000000-0000-0000-0000-000000000002&api-version=2021-04-01",
		httpmock.NewStringResponder(http.StatusOK, `{"value":[]}`),
	)

	httpmock.RegisterResponder("GET", "https://000000000000000000000000000000.01.tenant.api.powerplatform.com/governance/environmentGroups/00000000-0000-0000-0000-000000000001/ruleSets?api-version=2021-10-01-preview",
		httpmock.NewStringResponder(http.StatusOK, `{"value": [{"parameters": [],"id": "00000000-0000-0000-0000-000000000001","environmentFilter": {"type": "Include","values": []}}]}`),
	)

	httpmock.RegisterResponder("GET", "https://000000000000000000000000000000.01.tenant.api.powerplatform.com/governance/environmentGroups/00000000-0000-0000-0000-000000000002/ruleSets?api-version=2021-10-01-preview",
		httpmock.NewStringResponder(http.StatusOK, `{"value": [{"parameters": [],"id": "00000000-0000-0000-0000-000000000001","environmentFilter": {"type": "Include","values": []}}]}`),
	)

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Create_Environment_And_Add_Env_Group/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			getEnvironmentRequestInx++
			if getEnvironmentRequestInx < 14 {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Create_Environment_And_Add_Env_Group/get_environment_%s_1.json", id)).String()), nil
			} else if getEnvironmentRequestInx < 21 {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Create_Environment_And_Add_Env_Group/get_environment_%s_2.json", id)).String()), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Create_Environment_And_Add_Env_Group/get_environment_%s_3.json", id)).String()), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Create_Environment_And_Add_Env_Group/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environmentGroups?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			postEnvGroupRequestInx++
			resp := httpmock.NewStringResponse(http.StatusCreated, httpmock.File(fmt.Sprintf("tests/resource/Create_Environment_And_Add_Env_Group/post_environment_group_%d.json", postEnvGroupRequestInx)).String())
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environmentGroups/00000000-0000-0000-0000-000000000001?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Create_Environment_And_Add_Env_Group/get_environment_group_1.json").String())
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environmentGroups/00000000-0000-0000-0000-000000000002?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Create_Environment_And_Add_Env_Group/get_environment_group_2.json").String())
			return resp, nil
		},
	)

	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/environmentGroups/00000000-0000-0000-0000-00000000000(1|2)\?api-version=2021-04-01$`),
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			return resp, nil
		},
	)

	httpmock.RegisterResponder("PATCH", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("PATCH", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24expand=properties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Create_Environment_And_Add_Env_Group/get_environments_1.json").String())
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
				),
			},
			{
				Config: `
				resource "powerplatform_environment_group" "env_group" {
					display_name                              = "test_env_group"
					description                               = "test env group"
				}

				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					environment_group_id					  = powerplatform_environment_group.env_group.id
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_group_id", "00000000-0000-0000-0000-000000000001"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment_group" "env_group" {
					display_name                              = "test_env_group"
					description                               = "test env group"
				}

				resource "powerplatform_environment_group" "env_group_new" {
					display_name                              = "test_env_group_new"
					description                               = "test env group new"
				}

				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					environment_group_id					  = powerplatform_environment_group.env_group_new.id
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_group_id", "00000000-0000-0000-0000-000000000002"),
				),
			},
			{
				Config: `
					resource "powerplatform_environment_group" "env_group" {
					display_name                              = "test_env_group"
					description                               = "test env group"
				}

				resource "powerplatform_environment_group" "env_group_new" {
					display_name                              = "test_env_group_new"
					description                               = "test env group new"
				}

				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					environment_group_id					  = "00000000-0000-0000-0000-000000000000"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_group_id", "00000000-0000-0000-0000-000000000000"),
				),
			},
		},
	})
}

func TestAccEnvironmentsResource_Create_Environment_And_Add_Env_Group(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_group" "env_group" {
					display_name                              = "` + mocks.TestName() + `"
					description                               = "` + mocks.TestName() + `"
				}

				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					environment_group_id					  = powerplatform_environment_group.env_group.id
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "USD"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "environment_group_id", regexp.MustCompile(helpers.GuidRegex)),
				),
			},
			{
				Config: `
				resource "powerplatform_environment_group" "env_group" {
					display_name                              = "` + mocks.TestName() + `"
					description                               = "` + mocks.TestName() + `"
				}

				resource "powerplatform_environment_group" "env_group_new" {
					display_name                              = "` + mocks.TestName() + `_new"
					description                               = "` + mocks.TestName() + `_new"
				}

				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					environment_group_id					  = powerplatform_environment_group.env_group_new.id
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "USD"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "environment_group_id", regexp.MustCompile(helpers.GuidRegex)),
				),
			},
			{
				Config: `
				resource "powerplatform_environment_group" "env_group" {
					display_name                              = "` + mocks.TestName() + `"
					description                               = "` + mocks.TestName() + `"
				}

				resource "powerplatform_environment_group" "env_group_new" {
					display_name                              = "` + mocks.TestName() + `_new"
					description                               = "` + mocks.TestName() + `_new"
				}

				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					environment_group_id					  = "00000000-0000-0000-0000-000000000000"
					allow_bing_search 						  = false
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "USD"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_group_id", "00000000-0000-0000-0000-000000000000"),
				),
			},
		},
	})
}

func TestAccEnvironmentsResource_Create_Environment_No_Dataverse_Add_Dataverse_And_Add_Env_Group(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "USD"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
				),
			},
			{
				Config: `
				resource "powerplatform_environment_group" "env_group" {
					display_name                              = "` + mocks.TestName() + `"
					description                               = "` + mocks.TestName() + `"
				}
					
				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					environment_group_id					  = powerplatform_environment_group.env_group.id
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "USD"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
						
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "environment_group_id", regexp.MustCompile(helpers.GuidRegex)),
				),
			},
			{
				Config: `
				resource "powerplatform_environment_group" "env_group" {
					display_name                              = "` + mocks.TestName() + `"
					description                               = "` + mocks.TestName() + `"
				}
					
				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					environment_group_id					  = "00000000-0000-0000-0000-000000000000"
					allow_bing_search 						  = false
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "USD"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_group_id", "00000000-0000-0000-0000-000000000000"),
				),
			},
		},
	})
}

func TestAccEnvironmentsResource_Create_Environment_No_Dataverse_Add_Env_Group(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_group" "env_group" {
					display_name                              = "` + mocks.TestName() + `"
					description                               = "` + mocks.TestName() + `"
				}

				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					environment_group_id					  = powerplatform_environment_group.env_group.id
				}`,
				ExpectError: regexp.MustCompile(".*Attribute \"dataverse\" must be specified when \"environment_group_id\" is.*"),
				Check:       resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Update_Environment_Type(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	getLifecycleResponseInx := 0
	patchResponseInx := 0
	getResponseInx := 0
	deleteResponseInx := 0

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001/modifySku?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Operation-Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			deleteResponseInx++
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			if deleteResponseInx > 1 {
				panic("Environment was recreated unexpectedly. Check the test case.")
			}
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Update_Environment_Type/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Update_Environment_Type/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			getLifecycleResponseInx++
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			getResponseInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Update_Environment_Type/get_environment_%d.json", getResponseInx)).String()), nil
		})

	httpmock.RegisterResponder("PATCH", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			patchResponseInx++
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24expand=properties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Update_Environment_Type/get_environments_%d.json", patchResponseInx)).String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "Example"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Sandbox"),
				),
			},
			{
				Config: `
					resource "powerplatform_environment" "development" {
					display_name                              = "Example"
					location                                  = "europe"
					environment_type                          = "Production"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Production"),
				),
			},
		},
	})
}

func TestAccEnvironmentsResource_Validate_Update_Environment_Type(t *testing.T) {
	var environmentIdStep1 = &mocks.StateValue{}
	var environmentIdStep2 = &mocks.StateValue{}

	resource.Test(t, resource.TestCase{
		IsUnitTest:               false,
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Sandbox"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("powerplatform_environment.development", tfjsonpath.New("id"), mocks.GetStateValue(environmentIdStep1)),
				},
			},
			{
				Config: `
					resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `1"
					location                                  = "unitedstates"
					environment_type                          = "Production"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Production"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("powerplatform_environment.development", tfjsonpath.New("id"), mocks.GetStateValue(environmentIdStep2)),
				},
			},
			{
				Config: `
					resource "powerplatform_environment" "development" {
					display_name                              = "` + mocks.TestName() + `1"
					location                                  = "unitedstates"
					environment_type                          = "Production"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}`,
				Check: resource.ComposeTestCheckFunc(
					mocks.TestStateValueMatch(environmentIdStep1, environmentIdStep2, func(a, b *mocks.StateValue) error {
						if a.Value != b.Value {
							return fmt.Errorf("Environment ID from before environment_type change are not equal, recreate was triggred. '%s' != '%s'", a.Value, b.Value)
						}
						return nil
					}),
				),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Create_No_Dataverse_Region_Not_Available(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Dataverse_Region_Not_Available/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Create_No_Dataverse_Region_Not_Available/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Dataverse_Region_Not_Available/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile(".*Provisioning environment in azure region .* failed"),
				Config: `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					description                               = "description"
					cadence								      = "Frequent"
					location                                  = "europe"
					azure_region 							  = "northeurope"
					environment_type                          = "Sandbox"
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
				),
			},
		},
	})
}
