package powerplatform

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
)

func TestAccManagedEnvironmentsResource_Validate_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "example_managed_environment"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					domain           = "mydomainmanagedenvironment"
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
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "example_managed_environment"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					domain           = "mydomainmanagedenvironment"
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
				Config: ProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "example_managed_environment"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					domain           = "mydomainmanagedenvironment"
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
				Config: ProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "example_managed_environment"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					domain           = "mydomainmanagedenvironment"
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
				Config: ProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "example_managed_environment"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					domain           = "mydomainmanagedenvironment"
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
				Config: ProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "example_managed_environment"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					domain           = "mydomainmanagedenvironment"
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
				Config: ProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "example_managed_environment"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					domain           = "mydomainmanagedenvironment"
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
				Config: ProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "example_managed_environment"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					domain           = "mydomainmanagedenvironment"
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
				Config: ProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "example_managed_environment"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					domain           = "mydomainmanagedenvironment"
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
				Config: ProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "example_managed_environment"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					domain           = "mydomainmanagedenvironment"
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
