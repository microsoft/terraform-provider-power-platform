package powerplatform

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccTenantSettingsDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: AcceptanceTestsProviderConfig + `
				data "powerplatform_tenant_settings" "settings" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),

					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "disable_capacity_allocation_by_environment_admins", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "disable_environment_creation_by_non_admin_users", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "disable_newsletter_sendout", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "disable_nps_comments_reachout", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "disable_portals_creation_by_non_admin_users", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "disable_support_tickets_visible_by_all_users", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "disable_survey_feedback", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "disable_trial_environment_creation_by_non_admin_users", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.catalog_settings.power_catalog_audience_setting", regexp.MustCompile("^(All|None)$")),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.champions.disable_champions_invitation_reachout", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.champions.disable_skills_match_invitation_reachout", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.environments.disable_preferred_data_location_for_teams_environment", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.governance.disable_admin_digest", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.governance.disable_developer_environment_creation_by_non_admin_users", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.governance.enable_default_environment_routing", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.governance.policy.enable_desktop_flow_data_policy_management", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.intelligence.disable_copilot", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.intelligence.enable_open_ai_bot_publishing", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.licensing.disable_billing_policy_creation_by_non_admin_users", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.licensing.disable_use_of_unassigned_ai_builder_credits", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.licensing.enable_tenant_capacity_report_for_environment_admins", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.licensing.enable_tenant_licensing_report_for_environment_admins", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.licensing.storage_capacity_consumption_warning_threshold", regexp.MustCompile("^[0-9]+$")),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.model_experimentation.disable_data_logging", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.model_experimentation.enable_model_data_sharing", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_connection_sharing_with_everyone", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_create_from_figma", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_create_from_image", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_maker_match", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_members_indicator", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_share_with_everyone", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_unused_license_assignment", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.power_apps.enable_guests_to_make", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.power_automate.disable_copilot", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.search.disable_bing_video_search", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.search.disable_community_search", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.search.disable_docs_search", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.teams_integration.share_with_colleagues_user_limit", regexp.MustCompile("^[0-9]+$")),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.user_management_settings.enable_delete_disabled_user_in_all_environments", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "walk_me_opt_out", regexp.MustCompile(powerplatform_helpers.BooleanRegex)),
				),
			},
		},
	})
}

func TestUnitTestTenantSettingsDataSource_Validate_Read(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock_helpers.ActivateOAuthHttpMocks()

	httpmock.RegisterResponder("POST", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/listTenantSettings?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/tenant_settings/tests/datasource/post_list_tenant_settings.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: UnitTestsProviderConfig + `
				data "powerplatform_tenant_settings" "settings" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "disable_capacity_allocation_by_environment_admins", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "disable_environment_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "disable_newsletter_sendout", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "disable_nps_comments_reachout", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "disable_portals_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "disable_support_tickets_visible_by_all_users", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "disable_survey_feedback", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "disable_trial_environment_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.catalog_settings.power_catalog_audience_setting", "All"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.champions.disable_champions_invitation_reachout", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.champions.disable_skills_match_invitation_reachout", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.environments.disable_preferred_data_location_for_teams_environment", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.governance.disable_admin_digest", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.governance.disable_developer_environment_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.governance.enable_default_environment_routing", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.governance.policy.enable_desktop_flow_data_policy_management", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.intelligence.disable_copilot", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.intelligence.enable_open_ai_bot_publishing", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.licensing.disable_billing_policy_creation_by_non_admin_users", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.licensing.disable_use_of_unassigned_ai_builder_credits", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.licensing.enable_tenant_capacity_report_for_environment_admins", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.licensing.enable_tenant_licensing_report_for_environment_admins", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.licensing.storage_capacity_consumption_warning_threshold", "38"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.model_experimentation.disable_data_logging", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.model_experimentation.enable_model_data_sharing", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_connection_sharing_with_everyone", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_create_from_figma", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_create_from_image", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_maker_match", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_members_indicator", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_share_with_everyone", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.power_apps.disable_unused_license_assignment", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.power_apps.enable_guests_to_make", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.power_automate.disable_copilot", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.search.disable_bing_video_search", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.search.disable_community_search", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.search.disable_docs_search", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.teams_integration.share_with_colleagues_user_limit", "123"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.user_management_settings.enable_delete_disabled_user_in_all_environments", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", "walk_me_opt_out", "false"),
				),
			},
		},
	})
}
