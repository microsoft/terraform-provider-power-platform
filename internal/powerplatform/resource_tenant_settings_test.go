package powerplatform

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccTenantSettingsResource_Validate_Create(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: AcceptanceTestsProviderConfig + `
				resource "powerplatform_tenant_settings" "settings" {
					walk_me_opt_out                                       = false
					disable_support_tickets_visible_by_all_users          = false
					disable_trial_environment_creation_by_non_admin_users = false
					disable_capacity_allocation_by_environment_admins     = false
					disable_environment_creation_by_non_admin_users       = false
					disable_portals_creation_by_non_admin_users           = false
					disable_newsletter_sendout                            = false
					disable_nps_comments_reachout                         = false
					disable_survey_feedback                               = false
				  
					power_platform = {
					  search = {
						disable_docs_search       = false
						disable_community_search  = false
						disable_bing_video_search = false
					  }
					  teams_integration = {
						share_with_colleagues_user_limit = 10001
					  }
					  power_apps = {
						disable_share_with_everyone              = false
						enable_guests_to_make                    = false
						disable_members_indicator                = false
						disable_maker_match                      = false
						disable_unused_license_assignment        = false
						disable_create_from_image                = false
						disable_create_from_figma                = false
						disable_connection_sharing_with_everyone = false
					  }
					  power_automate = {
						disable_copilot           = false
						disable_copilot_with_bing = false
					  }
					  environments = {
						disable_preferred_data_location_for_teams_environment = false
					  }
					  governance = {
						disable_admin_digest                                      = false
						disable_developer_environment_creation_by_non_admin_users = false
						enable_default_environment_routing                        = false
						policy = {
						  enable_desktop_flow_data_policy_management = false
						}
					  }
					  licensing = {
						disable_billing_policy_creation_by_non_admin_users    = false
						enable_tenant_capacity_report_for_environment_admins  = false
						storage_capacity_consumption_warning_threshold        = 85
						enable_tenant_licensing_report_for_environment_admins = false
						disable_use_of_unassigned_ai_builder_credits          = false
						apply_auto_claim_to_only_managed_environments         = true
					  }
					  power_pages = {}
					  champions = {
						disable_champions_invitation_reachout    = false
						disable_skills_match_invitation_reachout = false
					  }
					  intelligence = {
						disable_copilot                   = false
						enable_open_ai_bot_publishing     = false
						disable_copilot_feedback          = false
						disable_copilot_feedback_metadata = false
					  }
					  model_experimentation = {
						enable_model_data_sharing = false
						disable_data_logging      = false
					  }
					  catalog_settings = {
						power_catalog_audience_setting = "All"
					  }
					  user_management_settings = {
						enable_delete_disabled_user_in_all_environments = false
					  }
					}
				  }`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_capacity_allocation_by_environment_admins", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_environment_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_newsletter_sendout", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_nps_comments_reachout", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_portals_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_support_tickets_visible_by_all_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_survey_feedback", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_trial_environment_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.catalog_settings.power_catalog_audience_setting", "All"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.champions.disable_champions_invitation_reachout", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.champions.disable_skills_match_invitation_reachout", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.environments.disable_preferred_data_location_for_teams_environment", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.disable_admin_digest", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.disable_developer_environment_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.enable_default_environment_routing", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.policy.enable_desktop_flow_data_policy_management", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.intelligence.disable_copilot", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.intelligence.enable_open_ai_bot_publishing", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.disable_billing_policy_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.disable_use_of_unassigned_ai_builder_credits", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.enable_tenant_capacity_report_for_environment_admins", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.enable_tenant_licensing_report_for_environment_admins", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.storage_capacity_consumption_warning_threshold", "85"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.model_experimentation.disable_data_logging", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.model_experimentation.enable_model_data_sharing", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_connection_sharing_with_everyone", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_create_from_figma", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_create_from_image", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_maker_match", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_members_indicator", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_share_with_everyone", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_unused_license_assignment", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.enable_guests_to_make", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_automate.disable_copilot", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.search.disable_bing_video_search", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.search.disable_community_search", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.search.disable_docs_search", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.teams_integration.share_with_colleagues_user_limit", "10001"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.user_management_settings.enable_delete_disabled_user_in_all_environments", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "walk_me_opt_out", "false"),
				),
			},
		},
	})
}

func TestUnitTestTenantSettingsResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock_helpers.ActivateOAuthHttpMocks()

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/updateTenantSettings?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/tenant_settings/tests/resources/Validate_Create/post_update_tenant_settings.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/listTenantSettings?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/tenant_settings/tests/resources/Validate_Create/post_update_tenant_settings.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: UnitTestsProviderConfig + `
				resource "powerplatform_tenant_settings" "settings" {
					walk_me_opt_out                                       = false
					disable_support_tickets_visible_by_all_users          = false
					disable_trial_environment_creation_by_non_admin_users = false
					disable_capacity_allocation_by_environment_admins     = false
					disable_environment_creation_by_non_admin_users       = false
					disable_portals_creation_by_non_admin_users           = false
					disable_newsletter_sendout                            = false
					disable_nps_comments_reachout                         = false
					disable_survey_feedback                               = false
				  
					power_platform = {
					  search = {
						disable_docs_search       = false
						disable_community_search  = false
						disable_bing_video_search = false
					  }
					  teams_integration = {
						share_with_colleagues_user_limit = 10001
					  }
					  power_apps = {
						disable_share_with_everyone              = false
						enable_guests_to_make                    = false
						disable_members_indicator                = false
						disable_maker_match                      = false
						disable_unused_license_assignment        = false
						disable_create_from_image                = false
						disable_create_from_figma                = false
						disable_connection_sharing_with_everyone = false
					  }
					  power_automate = {
						disable_copilot           = false
						disable_copilot_with_bing = false
					  }
					  environments = {
						disable_preferred_data_location_for_teams_environment = false
					  }
					  governance = {
						disable_admin_digest                                      = false
						disable_developer_environment_creation_by_non_admin_users = false
						enable_default_environment_routing                        = false
						policy = {
						  enable_desktop_flow_data_policy_management = false
						}
					  }
					  licensing = {
						disable_billing_policy_creation_by_non_admin_users    = false
						enable_tenant_capacity_report_for_environment_admins  = false
						storage_capacity_consumption_warning_threshold        = 85
						enable_tenant_licensing_report_for_environment_admins = false
						disable_use_of_unassigned_ai_builder_credits          = false
						apply_auto_claim_to_only_managed_environments         = true
					  }
					  power_pages = {}
					  champions = {
						disable_champions_invitation_reachout    = false
						disable_skills_match_invitation_reachout = false
					  }
					  intelligence = {
						disable_copilot                   = false
						enable_open_ai_bot_publishing     = false
						disable_copilot_feedback          = false
						disable_copilot_feedback_metadata = false
					  }
					  model_experimentation = {
						enable_model_data_sharing = false
						disable_data_logging      = false
					  }
					  catalog_settings = {
						power_catalog_audience_setting = "All"
					  }
					  user_management_settings = {
						enable_delete_disabled_user_in_all_environments = false
					  }
					}
				  }`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_capacity_allocation_by_environment_admins", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_environment_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_newsletter_sendout", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_nps_comments_reachout", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_portals_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_support_tickets_visible_by_all_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_survey_feedback", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_trial_environment_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.catalog_settings.power_catalog_audience_setting", "All"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.champions.disable_champions_invitation_reachout", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.champions.disable_skills_match_invitation_reachout", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.environments.disable_preferred_data_location_for_teams_environment", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.disable_admin_digest", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.disable_developer_environment_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.enable_default_environment_routing", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.policy.enable_desktop_flow_data_policy_management", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.intelligence.disable_copilot", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.intelligence.enable_open_ai_bot_publishing", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.disable_billing_policy_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.disable_use_of_unassigned_ai_builder_credits", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.enable_tenant_capacity_report_for_environment_admins", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.enable_tenant_licensing_report_for_environment_admins", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.storage_capacity_consumption_warning_threshold", "85"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.model_experimentation.disable_data_logging", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.model_experimentation.enable_model_data_sharing", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_connection_sharing_with_everyone", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_create_from_figma", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_create_from_image", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_maker_match", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_members_indicator", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_share_with_everyone", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_unused_license_assignment", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.enable_guests_to_make", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_automate.disable_copilot", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.search.disable_bing_video_search", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.search.disable_community_search", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.search.disable_docs_search", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.teams_integration.share_with_colleagues_user_limit", "10001"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.user_management_settings.enable_delete_disabled_user_in_all_environments", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "walk_me_opt_out", "false"),
				),
			},
		},
	})
}

func TestAccTenantSettingsResource_Validate_Update(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: AcceptanceTestsProviderConfig + `
				resource "powerplatform_tenant_settings" "settings" {
					walk_me_opt_out                                       = true
					disable_support_tickets_visible_by_all_users          = true
					disable_trial_environment_creation_by_non_admin_users = true
					disable_capacity_allocation_by_environment_admins     = true
					disable_environment_creation_by_non_admin_users       = true
					disable_portals_creation_by_non_admin_users           = true
					disable_newsletter_sendout                            = true
					disable_nps_comments_reachout                         = true
					disable_survey_feedback                               = true
				  
					power_platform = {
					  search = {
						disable_docs_search       = true
						disable_community_search  = true
						disable_bing_video_search = true
					  }
					  teams_integration = {
						share_with_colleagues_user_limit = 999
					  }
					  power_apps = {
						disable_share_with_everyone              = true
						enable_guests_to_make                    = true
						disable_members_indicator                = true
						disable_maker_match                      = true
						disable_unused_license_assignment        = true
						disable_create_from_image                = true
						disable_create_from_figma                = true
						disable_connection_sharing_with_everyone = true
					  }
					  power_automate = {
						disable_copilot           = true
						disable_copilot_with_bing = true
					  }
					  environments = {
						disable_preferred_data_location_for_teams_environment = true
					  }
					  governance = {
						disable_admin_digest                                      = true
						disable_developer_environment_creation_by_non_admin_users = false
						enable_default_environment_routing                        = true
						policy = {
						  enable_desktop_flow_data_policy_management = true
						}
					  }
					  licensing = {
						disable_billing_policy_creation_by_non_admin_users    = true
						enable_tenant_capacity_report_for_environment_admins  = true
						storage_capacity_consumption_warning_threshold        = 95
						enable_tenant_licensing_report_for_environment_admins = true
						disable_use_of_unassigned_ai_builder_credits          = true
						apply_auto_claim_to_only_managed_environments         = false
					  }
					  power_pages = {}
					  champions = {
						disable_champions_invitation_reachout    = true
						disable_skills_match_invitation_reachout = true
					  }
					  intelligence = {
						disable_copilot                   = true
						enable_open_ai_bot_publishing     = true
						disable_copilot_feedback          = true
						disable_copilot_feedback_metadata = true
					  }
					  model_experimentation = {
						enable_model_data_sharing = true
						disable_data_logging      = true
					  }
					  catalog_settings = {
						power_catalog_audience_setting = "All"
					  }
					  user_management_settings = {
						enable_delete_disabled_user_in_all_environments = true
					  }
					}
				  }`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_capacity_allocation_by_environment_admins", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_environment_creation_by_non_admin_users", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_newsletter_sendout", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_nps_comments_reachout", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_portals_creation_by_non_admin_users", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_support_tickets_visible_by_all_users", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_survey_feedback", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_trial_environment_creation_by_non_admin_users", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.catalog_settings.power_catalog_audience_setting", "All"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.champions.disable_champions_invitation_reachout", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.champions.disable_skills_match_invitation_reachout", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.environments.disable_preferred_data_location_for_teams_environment", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.disable_admin_digest", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.disable_developer_environment_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.enable_default_environment_routing", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.policy.enable_desktop_flow_data_policy_management", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.intelligence.disable_copilot", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.intelligence.enable_open_ai_bot_publishing", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.disable_billing_policy_creation_by_non_admin_users", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.disable_use_of_unassigned_ai_builder_credits", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.enable_tenant_capacity_report_for_environment_admins", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.enable_tenant_licensing_report_for_environment_admins", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.storage_capacity_consumption_warning_threshold", "95"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.model_experimentation.disable_data_logging", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.model_experimentation.enable_model_data_sharing", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_connection_sharing_with_everyone", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_create_from_figma", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_create_from_image", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_maker_match", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_members_indicator", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_share_with_everyone", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_unused_license_assignment", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.enable_guests_to_make", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_automate.disable_copilot", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.search.disable_bing_video_search", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.search.disable_community_search", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.search.disable_docs_search", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.teams_integration.share_with_colleagues_user_limit", "999"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.user_management_settings.enable_delete_disabled_user_in_all_environments", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "walk_me_opt_out", "true"),
				),
			},
			{
				Config: AcceptanceTestsProviderConfig + `
				resource "powerplatform_tenant_settings" "settings" {
					walk_me_opt_out                                       = false
					disable_support_tickets_visible_by_all_users          = false
					disable_trial_environment_creation_by_non_admin_users = false
					disable_capacity_allocation_by_environment_admins     = false
					disable_environment_creation_by_non_admin_users       = false
					disable_portals_creation_by_non_admin_users           = false
					disable_newsletter_sendout                            = false
					disable_nps_comments_reachout                         = false
					disable_survey_feedback                               = false
				  
					power_platform = {
					  search = {
						disable_docs_search       = false
						disable_community_search  = false
						disable_bing_video_search = false
					  }
					  teams_integration = {
						share_with_colleagues_user_limit = 10001
					  }
					  power_apps = {
						disable_share_with_everyone              = false
						enable_guests_to_make                    = false
						disable_members_indicator                = false
						disable_maker_match                      = false
						disable_unused_license_assignment        = false
						disable_create_from_image                = false
						disable_create_from_figma                = false
						disable_connection_sharing_with_everyone = false
					  }
					  power_automate = {
						disable_copilot           = false
						disable_copilot_with_bing = false
					  }
					  environments = {
						disable_preferred_data_location_for_teams_environment = false
					  }
					  governance = {
						disable_admin_digest                                      = false
						disable_developer_environment_creation_by_non_admin_users = false
						enable_default_environment_routing                        = false
						policy = {
						  enable_desktop_flow_data_policy_management = false
						}
					  }
					  licensing = {
						disable_billing_policy_creation_by_non_admin_users    = false
						enable_tenant_capacity_report_for_environment_admins  = false
						storage_capacity_consumption_warning_threshold        = 85
						enable_tenant_licensing_report_for_environment_admins = false
						disable_use_of_unassigned_ai_builder_credits          = false
						apply_auto_claim_to_only_managed_environments         = true
					  }
					  power_pages = {}
					  champions = {
						disable_champions_invitation_reachout    = false
						disable_skills_match_invitation_reachout = false
					  }
					  intelligence = {
						disable_copilot                   = false
						enable_open_ai_bot_publishing     = false
						disable_copilot_feedback          = false
						disable_copilot_feedback_metadata = false
					  }
					  model_experimentation = {
						enable_model_data_sharing = false
						disable_data_logging      = false
					  }
					  catalog_settings = {
						power_catalog_audience_setting = "All"
					  }
					  user_management_settings = {
						enable_delete_disabled_user_in_all_environments = false
					  }
					}
				  }`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_capacity_allocation_by_environment_admins", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_environment_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_newsletter_sendout", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_nps_comments_reachout", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_portals_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_support_tickets_visible_by_all_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_survey_feedback", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_trial_environment_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.catalog_settings.power_catalog_audience_setting", "All"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.champions.disable_champions_invitation_reachout", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.champions.disable_skills_match_invitation_reachout", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.environments.disable_preferred_data_location_for_teams_environment", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.disable_admin_digest", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.disable_developer_environment_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.enable_default_environment_routing", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.policy.enable_desktop_flow_data_policy_management", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.intelligence.disable_copilot", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.intelligence.enable_open_ai_bot_publishing", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.disable_billing_policy_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.disable_use_of_unassigned_ai_builder_credits", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.enable_tenant_capacity_report_for_environment_admins", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.enable_tenant_licensing_report_for_environment_admins", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.storage_capacity_consumption_warning_threshold", "85"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.model_experimentation.disable_data_logging", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.model_experimentation.enable_model_data_sharing", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_connection_sharing_with_everyone", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_create_from_figma", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_create_from_image", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_maker_match", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_members_indicator", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_share_with_everyone", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_unused_license_assignment", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.enable_guests_to_make", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_automate.disable_copilot", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.search.disable_bing_video_search", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.search.disable_community_search", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.search.disable_docs_search", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.teams_integration.share_with_colleagues_user_limit", "10001"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.user_management_settings.enable_delete_disabled_user_in_all_environments", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "walk_me_opt_out", "false"),
				),
			},
		},
	})
}

func TestUnitTestTenantSettingsResource_Validate_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock_helpers.ActivateOAuthHttpMocks()

	updateSettingsResponseCount := 0
	listSettingsResponseCount := 0

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/updateTenantSettings?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			updateSettingsResponseCount++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/tenant_settings/tests/resources/Validate_Update/post_update_tenant_settings_%d.json", updateSettingsResponseCount)).String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/listTenantSettings?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			listSettingsResponseCount++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/tenant_settings/tests/resources/Validate_Update/post_update_tenant_settings_%d.json", listSettingsResponseCount)).String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: UnitTestsProviderConfig + `
				resource "powerplatform_tenant_settings" "settings" {
					walk_me_opt_out                                       = true
					disable_support_tickets_visible_by_all_users          = true
					disable_trial_environment_creation_by_non_admin_users = true
					disable_capacity_allocation_by_environment_admins     = true
					disable_environment_creation_by_non_admin_users       = true
					disable_portals_creation_by_non_admin_users           = true
					disable_newsletter_sendout                            = true
					disable_nps_comments_reachout                         = true
					disable_survey_feedback                               = true
				  
					power_platform = {
					  search = {
						disable_docs_search       = true
						disable_community_search  = true
						disable_bing_video_search = true
					  }
					  teams_integration = {
						share_with_colleagues_user_limit = 999
					  }
					  power_apps = {
						disable_share_with_everyone              = true
						enable_guests_to_make                    = true
						disable_members_indicator                = true
						disable_maker_match                      = true
						disable_unused_license_assignment        = true
						disable_create_from_image                = true
						disable_create_from_figma                = true
						disable_connection_sharing_with_everyone = true
					  }
					  power_automate = {
						disable_copilot           = true
						disable_copilot_with_bing = true
					  }
					  environments = {
						disable_preferred_data_location_for_teams_environment = true
					  }
					  governance = {
						disable_admin_digest                                      = true
						disable_developer_environment_creation_by_non_admin_users = false
						enable_default_environment_routing                        = true
						policy = {
						  enable_desktop_flow_data_policy_management = true
						}
					  }
					  licensing = {
						disable_billing_policy_creation_by_non_admin_users    = true
						enable_tenant_capacity_report_for_environment_admins  = true
						storage_capacity_consumption_warning_threshold        = 95
						enable_tenant_licensing_report_for_environment_admins = true
						disable_use_of_unassigned_ai_builder_credits          = true
						apply_auto_claim_to_only_managed_environments         = false
					  }
					  power_pages = {}
					  champions = {
						disable_champions_invitation_reachout    = true
						disable_skills_match_invitation_reachout = true
					  }
					  intelligence = {
						disable_copilot                   = true
						enable_open_ai_bot_publishing     = true
						disable_copilot_feedback          = true
						disable_copilot_feedback_metadata = true
					  }
					  model_experimentation = {
						enable_model_data_sharing = true
						disable_data_logging      = true
					  }
					  catalog_settings = {
						power_catalog_audience_setting = "All"
					  }
					  user_management_settings = {
						enable_delete_disabled_user_in_all_environments = true
					  }
					}
				  }`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_capacity_allocation_by_environment_admins", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_environment_creation_by_non_admin_users", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_newsletter_sendout", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_nps_comments_reachout", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_portals_creation_by_non_admin_users", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_support_tickets_visible_by_all_users", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_survey_feedback", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_trial_environment_creation_by_non_admin_users", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.catalog_settings.power_catalog_audience_setting", "All"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.champions.disable_champions_invitation_reachout", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.champions.disable_skills_match_invitation_reachout", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.environments.disable_preferred_data_location_for_teams_environment", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.disable_admin_digest", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.disable_developer_environment_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.enable_default_environment_routing", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.policy.enable_desktop_flow_data_policy_management", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.intelligence.disable_copilot", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.intelligence.enable_open_ai_bot_publishing", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.disable_billing_policy_creation_by_non_admin_users", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.disable_use_of_unassigned_ai_builder_credits", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.enable_tenant_capacity_report_for_environment_admins", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.enable_tenant_licensing_report_for_environment_admins", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.storage_capacity_consumption_warning_threshold", "95"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.model_experimentation.disable_data_logging", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.model_experimentation.enable_model_data_sharing", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_connection_sharing_with_everyone", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_create_from_figma", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_create_from_image", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_maker_match", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_members_indicator", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_share_with_everyone", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_unused_license_assignment", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.enable_guests_to_make", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_automate.disable_copilot", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.search.disable_bing_video_search", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.search.disable_community_search", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.search.disable_docs_search", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.teams_integration.share_with_colleagues_user_limit", "999"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.user_management_settings.enable_delete_disabled_user_in_all_environments", "true"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "walk_me_opt_out", "true"),
				),
			},
			{
				Config: UnitTestsProviderConfig + `
				resource "powerplatform_tenant_settings" "settings" {
					walk_me_opt_out                                       = false
					disable_support_tickets_visible_by_all_users          = false
					disable_trial_environment_creation_by_non_admin_users = false
					disable_capacity_allocation_by_environment_admins     = false
					disable_environment_creation_by_non_admin_users       = false
					disable_portals_creation_by_non_admin_users           = false
					disable_newsletter_sendout                            = false
					disable_nps_comments_reachout                         = false
					disable_survey_feedback                               = false

					power_platform = {
					  search = {
						disable_docs_search       = false
						disable_community_search  = false
						disable_bing_video_search = false
					  }
					  teams_integration = {
						share_with_colleagues_user_limit = 10001
					  }
					  power_apps = {
						disable_share_with_everyone              = false
						enable_guests_to_make                    = false
						disable_members_indicator                = false
						disable_maker_match                      = false
						disable_unused_license_assignment        = false
						disable_create_from_image                = false
						disable_create_from_figma                = false
						disable_connection_sharing_with_everyone = false
					  }
					  power_automate = {
						disable_copilot           = false
						disable_copilot_with_bing = false
					  }
					  environments = {
						disable_preferred_data_location_for_teams_environment = false
					  }
					  governance = {
						disable_admin_digest                                      = false
						disable_developer_environment_creation_by_non_admin_users = false
						enable_default_environment_routing                        = false
						policy = {
						  enable_desktop_flow_data_policy_management = false
						}
					  }
					  licensing = {
						disable_billing_policy_creation_by_non_admin_users    = false
						enable_tenant_capacity_report_for_environment_admins  = false
						storage_capacity_consumption_warning_threshold        = 85
						enable_tenant_licensing_report_for_environment_admins = false
						disable_use_of_unassigned_ai_builder_credits          = false
						apply_auto_claim_to_only_managed_environments         = true
					  }
					  power_pages = {}
					  champions = {
						disable_champions_invitation_reachout    = false
						disable_skills_match_invitation_reachout = false
					  }
					  intelligence = {
						disable_copilot                   = false
						enable_open_ai_bot_publishing     = false
						disable_copilot_feedback          = false
						disable_copilot_feedback_metadata = false
					  }
					  model_experimentation = {
						enable_model_data_sharing = false
						disable_data_logging      = false
					  }
					  catalog_settings = {
						power_catalog_audience_setting = "All"
					  }
					  user_management_settings = {
						enable_delete_disabled_user_in_all_environments = false
					  }
					}
				  }`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_capacity_allocation_by_environment_admins", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_environment_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_newsletter_sendout", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_nps_comments_reachout", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_portals_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_support_tickets_visible_by_all_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_survey_feedback", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_trial_environment_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.catalog_settings.power_catalog_audience_setting", "All"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.champions.disable_champions_invitation_reachout", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.champions.disable_skills_match_invitation_reachout", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.environments.disable_preferred_data_location_for_teams_environment", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.disable_admin_digest", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.disable_developer_environment_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.enable_default_environment_routing", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.governance.policy.enable_desktop_flow_data_policy_management", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.intelligence.disable_copilot", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.intelligence.enable_open_ai_bot_publishing", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.disable_billing_policy_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.disable_use_of_unassigned_ai_builder_credits", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.enable_tenant_capacity_report_for_environment_admins", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.enable_tenant_licensing_report_for_environment_admins", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.storage_capacity_consumption_warning_threshold", "85"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.model_experimentation.disable_data_logging", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.model_experimentation.enable_model_data_sharing", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_connection_sharing_with_everyone", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_create_from_figma", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_create_from_image", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_maker_match", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_members_indicator", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_share_with_everyone", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_unused_license_assignment", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_apps.enable_guests_to_make", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.power_automate.disable_copilot", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.search.disable_bing_video_search", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.search.disable_community_search", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.search.disable_docs_search", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.teams_integration.share_with_colleagues_user_limit", "10001"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.user_management_settings.enable_delete_disabled_user_in_all_environments", "false"),
					resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "walk_me_opt_out", "false"),
				),
			},
		},
	})
}
