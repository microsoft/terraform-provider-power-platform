# powerplatform_tenant_settings

This resource is used to manage tenant-wide settings for Power Platform, including governance, licensing, and user experience controls.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud) |
| ------------------- | ----------- | -------------------------------- |
| Read                | `POST`      | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/listTenantSettings?api-version=2020-10-01` |
| Update              | `POST`      | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/updateTenantSettings?api-version=2023-06-01` |

## Attribute Mapping

| Resource Attribute                                       | API Request / Response JSON Field                       |
| -------------------------------------------------------- | -------------------------------------------------------- |
| `walk_me_opt_out`                                        | `walkMeOptOut`                                           |
| `disable_nps_comments_reachout`                          | `disableNPSCommentsReachout`                             |
| `disable_newsletter_sendout`                             | `disableNewsletterSendout`                               |
| `disable_environment_creation_by_non_admin_users`        | `disableEnvironmentCreationByNonAdminUsers`              |
| `disable_portals_creation_by_non_admin_users`            | `disablePortalsCreationByNonAdminUsers`                  |
| `disable_survey_feedback`                                | `disableSurveyFeedback`                                  |
| `disable_trial_environment_creation_by_non_admin_users`  | `disableTrialEnvironmentCreationByNonAdminUsers`         |
| `disable_capacity_allocation_by_environment_admins`      | `disableCapacityAllocationByEnvironmentAdmins`           |
| `disable_support_tickets_visible_by_all_users`           | `disableSupportTicketsVisibleByAllUsers`                 |
| `power_platform.search.disable_docs_search`              | `powerPlatform.search.disableDocsSearch`                 |
| `power_platform.search.disable_community_search`         | `powerPlatform.search.disableCommunitySearch`            |
| `power_platform.search.disable_bing_video_search`        | `powerPlatform.search.disableBingVideoSearch`            |
| `power_platform.teams_integration.share_with_colleagues_user_limit` | `powerPlatform.teamsIntegration.shareWithColleaguesUserLimit` |
| `power_platform.power_apps.disable_share_with_everyone`  | `powerPlatform.powerApps.disableShareWithEveryone`       |
| `power_platform.power_apps.enable_guests_to_make`        | `powerPlatform.powerApps.enableGuestsToMake`             |
| `power_platform.power_apps.disable_maker_match`          | `powerPlatform.powerApps.disableMakerMatch`              |
| `power_platform.power_apps.disable_unused_license_assignment` | `powerPlatform.powerApps.disableUnusedLicenseAssignment` |
| `power_platform.power_apps.disable_create_from_image`    | `powerPlatform.powerApps.disableCreateFromImage`         |
| `power_platform.power_apps.disable_create_from_figma`    | `powerPlatform.powerApps.disableCreateFromFigma`         |
| `power_platform.power_apps.disable_connection_sharing_with_everyone` | `powerPlatform.powerApps.disableConnectionSharingWithEveryone` |
| `power_platform.power_automate.disable_copilot`          | `powerPlatform.powerAutomate.disableCopilot`             |
| `power_platform.environments.disable_preferred_data_location_for_teams_environment` | `powerPlatform.environments.disablePreferredDataLocationForTeamsEnvironment` |
| `power_platform.governance.disable_admin_digest`         | `powerPlatform.governance.disableAdminDigest`            |
| `power_platform.governance.disable_developer_environment_creation_by_non_admin_users` | `powerPlatform.governance.disableDeveloperEnvironmentCreationByNonAdminUsers` |
| `power_platform.governance.enable_default_environment_routing` | `powerPlatform.governance.enableDefaultEnvironmentRouting` |
| `power_platform.governance.environment_routing_all_makers` | `powerPlatform.governance.environmentRoutingAllMakers`   |
| `power_platform.governance.environment_routing_target_environment_group_id` | `powerPlatform.governance.environmentRoutingTargetEnvironmentGroupId` |
| `power_platform.governance.environment_routing_target_security_group_id` | `powerPlatform.governance.environmentRoutingTargetSecurityGroupId` |
| `power_platform.governance.policy.enable_desktop_flow_data_policy_management` | `powerPlatform.governance.policy.enableDesktopFlowDataPolicyManagement` |
| `power_platform.licensing.disable_billing_policy_creation_by_non_admin_users` | `powerPlatform.licensing.disableBillingPolicyCreationByNonAdminUsers` |
| `power_platform.licensing.enable_tenant_capacity_report_for_environment_admins` | `powerPlatform.licensing.enableTenantCapacityReportForEnvironmentAdmins` |
| `power_platform.licensing.storage_capacity_consumption_warning_threshold` | `powerPlatform.licensing.storageCapacityConsumptionWarningThreshold` |
| `power_platform.licensing.enable_tenant_licensing_report_for_environment_admins` | `powerPlatform.licensing.enableTenantLicensingReportForEnvironmentAdmins` |
| `power_platform.licensing.disable_use_of_unassigned_ai_builder_credits` | `powerPlatform.licensing.disableUseOfUnassignedAIBuilderCredits` |
| `power_platform.champions.disable_champions_invitation_reachout` | `powerPlatform.champions.disableChampionsInvitationReachout` |
| `power_platform.champions.disable_skills_match_invitation_reachout` | `powerPlatform.champions.disableSkillsMatchInvitationReachout` |
| `power_platform.intelligence.disable_copilot`            | `powerPlatform.intelligence.disableCopilot`              |
| `power_platform.intelligence.enable_open_ai_bot_publishing` | `powerPlatform.intelligence.enableOpenAiBotPublishing` |
| `power_platform.model_experimentation.enable_model_data_sharing` | `powerPlatform.modelExperimentation.enableModelDataSharing` |
| `power_platform.model_experimentation.disable_data_logging` | `powerPlatform.modelExperimentation.disableDataLogging` |
| `power_platform.catalog_settings.power_catalog_audience_setting` | `powerPlatform.catalogSettings.powerCatalogAudienceSetting` |
| `power_platform.user_management_settings.enable_delete_disabled_user_in_all_environments` | `powerPlatform.userManagementSettings.enableDeleteDisabledUserinAllEnvironments` |

### Example API Response

An example of the API response used by this resource (showing updated tenant settings) can be found in the test fixture [`tenant_settings/tests/resources/Validate_Create/post_update_tenant_settings.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/tenant_settings/tests/resources/Validate_Create/post_update_tenant_settings.json).
