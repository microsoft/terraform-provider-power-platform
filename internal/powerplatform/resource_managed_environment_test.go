package powerplatform

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccManagedEnvironmentsResource_Validate_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "example_managed_environment"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					security_group_id = "00000000-0000-0000-0000-000000000000"
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = true
					is_group_sharing_disabled  = true
					limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					maker_onboarding_markdown  = "this is test markdown"
					maker_onboarding_url       = "http://www.example.com"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),

					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "is_usage_insights_disabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "protection_level", "Standard"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "is_group_sharing_disabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "limit_sharing_mode", "ExcludeSharingToSecurityGroups"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "max_limit_user_sharing", "10"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_mode", "None"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "suppress_validation_emails", "true"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "maker_onboarding_markdown", "this is test markdown"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "maker_onboarding_url", "http://www.example.com"),
				),
			},
		},
	})
}

func TestAccManagedEnvironmentsResource_Validate_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "example_managed_environment"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					security_group_id = "00000000-0000-0000-0000-000000000000"
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = true
					is_group_sharing_disabled  = true
					limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					maker_onboarding_markdown  = "this is test markdown"
					maker_onboarding_url       = "http://www.example.com"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
				),
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "example_managed_environment"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					security_group_id = "00000000-0000-0000-0000-000000000000"
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = true
					limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					maker_onboarding_markdown  = "this is test markdown"
					maker_onboarding_url       = "http://www.example.com"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "is_usage_insights_disabled", "false"),
				),
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "example_managed_environment"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					security_group_id = "00000000-0000-0000-0000-000000000000"
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					maker_onboarding_markdown  = "this is test markdown"
					maker_onboarding_url       = "http://www.example.com"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "is_group_sharing_disabled", "false"),
				),
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "example_managed_environment"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					security_group_id = "00000000-0000-0000-0000-000000000000"
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					maker_onboarding_markdown  = "this is test markdown"
					maker_onboarding_url       = "http://www.example.com"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "limit_sharing_mode", "NoLimit"),
				),
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "example_managed_environment"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					security_group_id = "00000000-0000-0000-0000-000000000000"
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					maker_onboarding_markdown  = "this is test markdown"
					maker_onboarding_url       = "http://www.example.com"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "max_limit_user_sharing", "-1"),
				),
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "example_managed_environment"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					security_group_id = "00000000-0000-0000-0000-000000000000"
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "None"
					suppress_validation_emails = false
					maker_onboarding_markdown  = "this is test markdown"
					maker_onboarding_url       = "http://www.example.com"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "suppress_validation_emails", "false"),
				),
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "example_managed_environment"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					security_group_id = "00000000-0000-0000-0000-000000000000"
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "None"
					suppress_validation_emails = false
					maker_onboarding_markdown  = "this is test markdown changed"
					maker_onboarding_url       = "http://www.example.com"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "maker_onboarding_markdown", "this is test markdown changed"),
				),
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "example_managed_environment"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					security_group_id = "00000000-0000-0000-0000-000000000000"
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "None"
					suppress_validation_emails = false
					maker_onboarding_markdown  = "this is test markdown changed"
					maker_onboarding_url       = "http://www.example-changed.com"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "maker_onboarding_url", "http://www.example-changed.com"),
				),
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "example_managed_environment"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					security_group_id = "00000000-0000-0000-0000-000000000000"
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "Warn"
					suppress_validation_emails = false
					maker_onboarding_markdown  = "this is test markdown changed"
					maker_onboarding_url       = "http://www.example-changed.com"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_mode", "Warn"),
				),
			},
		},
	})
}

func TestUnitManagedEnvironmentsResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateEnvironmentHttpMocks()

	patchResponseInx := 0

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/environment/tests/resource/Validate_Create_And_Update/get_environments_%d.json", patchResponseInx)).String()), nil
		})
	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/managed_environment/tests/resource/Validate_Create_And_Update/get_lifecycle.json").String()), nil
		})
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/managed_environment/tests/resource/Validate_Create_And_Update/get_environment_create_response.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = "00000000-0000-0000-0000-000000000001"
					is_usage_insights_disabled = true
					is_group_sharing_disabled  = true
					limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					maker_onboarding_markdown  = "this is test markdown"
					maker_onboarding_url       = "http://www.example.com"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),

					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "is_usage_insights_disabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "protection_level", "Standard"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "is_group_sharing_disabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "limit_sharing_mode", "ExcludeSharingToSecurityGroups"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "max_limit_user_sharing", "10"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_mode", "None"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "suppress_validation_emails", "true"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "maker_onboarding_markdown", "this is test markdown"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "maker_onboarding_url", "http://www.example.com"),
				),
			},
		},
	})
}

func TestUnitManagedEnvironmentsResource_Validate_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateEnvironmentHttpMocks()

	patchResponseInx := -1

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Update/get_environments_0.json").String()), nil
		})
	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/managed_environment/tests/resource/Validate_Create_And_Update/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			patchResponseInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/managed_environment/tests/resource/Validate_Create_And_Update/get_environment_create_response_extended_%d.json", patchResponseInx)).String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
						resource "powerplatform_managed_environment" "managed_development" {
						environment_id             = "00000000-0000-0000-0000-000000000001"
						is_usage_insights_disabled = true
						is_group_sharing_disabled  = false
						limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
						max_limit_user_sharing     = 10
						solution_checker_mode      = "None"
						suppress_validation_emails = true
						maker_onboarding_markdown  = "this is test markdown"
						maker_onboarding_url       = "http://www.example.com"
					}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "is_group_sharing_disabled", "false"),
				),
			},
			{
				Config: TestsProviderConfig + `
					resource "powerplatform_managed_environment" "managed_development" {
						environment_id             = "00000000-0000-0000-0000-000000000001"
						is_usage_insights_disabled = true
						is_group_sharing_disabled  = false
						limit_sharing_mode         = "NoLimit"
						max_limit_user_sharing     = 10
						solution_checker_mode      = "None"
						suppress_validation_emails = true
						maker_onboarding_markdown  = "this is test markdown"
						maker_onboarding_url       = "http://www.example.com"
					}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "limit_sharing_mode", "NoLimit"),
				),
			},
			{
				Config: TestsProviderConfig + `
					resource "powerplatform_managed_environment" "managed_development" {
						environment_id             = "00000000-0000-0000-0000-000000000001"
						is_usage_insights_disabled = true
						is_group_sharing_disabled  = false
						limit_sharing_mode         = "NoLimit"
						max_limit_user_sharing     = -1
						solution_checker_mode      = "None"
						suppress_validation_emails = true
						maker_onboarding_markdown  = "this is test markdown"
						maker_onboarding_url       = "http://www.example.com"
					}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "max_limit_user_sharing", "-1"),
				),
			},
			{
				Config: TestsProviderConfig + `
					resource "powerplatform_managed_environment" "managed_development" {
						environment_id             = "00000000-0000-0000-0000-000000000001"
						is_usage_insights_disabled = true
						is_group_sharing_disabled  = false
						limit_sharing_mode         = "NoLimit"
						max_limit_user_sharing     = -1
						solution_checker_mode      = "Warn"
						suppress_validation_emails = true
						maker_onboarding_markdown  = "this is test markdown"
						maker_onboarding_url       = "http://www.example.com"
					}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_mode", "Warn"),
				),
			},
			{
				Config: TestsProviderConfig + `
					resource "powerplatform_managed_environment" "managed_development" {
						environment_id             = "00000000-0000-0000-0000-000000000001"
						is_usage_insights_disabled = true
						is_group_sharing_disabled  = false
						limit_sharing_mode         = "NoLimit"
						max_limit_user_sharing     = -1
						solution_checker_mode      = "Warn"
						suppress_validation_emails = false
						maker_onboarding_markdown  = "this is test markdown"
						maker_onboarding_url       = "http://www.example.com"
					}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "suppress_validation_emails", "false"),
				),
			},
			{
				Config: TestsProviderConfig + `
					resource "powerplatform_managed_environment" "managed_development" {
						environment_id             = "00000000-0000-0000-0000-000000000001"
						is_usage_insights_disabled = true
						is_group_sharing_disabled  = false
						limit_sharing_mode         = "NoLimit"
						max_limit_user_sharing     = -1
						solution_checker_mode      = "Warn"
						suppress_validation_emails = false
						maker_onboarding_markdown  = "this is test markdown 2"
						maker_onboarding_url       = "http://www.example.com"
					}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "maker_onboarding_markdown", "this is test markdown 2"),
				),
			},
			{
				Config: TestsProviderConfig + `
					resource "powerplatform_managed_environment" "managed_development" {
						environment_id             = "00000000-0000-0000-0000-000000000001"
						is_usage_insights_disabled = true
						is_group_sharing_disabled  = false
						limit_sharing_mode         = "NoLimit"
						max_limit_user_sharing     = -1
						solution_checker_mode      = "Warn"
						suppress_validation_emails = false
						maker_onboarding_markdown  = "this is test markdown 2"
						maker_onboarding_url       = "http://www.example2.com"
					}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "maker_onboarding_url", "http://www.example2.com"),
				),
			},
		},
	})
}
