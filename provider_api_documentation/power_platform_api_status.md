# Power Platform API Migration Status

This document provides a quick reference for the migration status of all provider resources and data sources to the Power Platform API.

## Status Definitions

- **Not in PPAPI**: The endpoint does not exist in the official Power Platform REST API
- **WebAPI**: Uses Dataverse Web API (out of scope for Power Platform API migration)
- **Requires changes**: Exists in PPAPI but requires modifications or has missing attributes
- **Ready for implementation**: Can be migrated to PPAPI without changes

## Resources

| Name | API | Status |
|------|-----|--------|
| powerplatform_admin_management_application | [API Documentation](./resource_powerplatform_admin_management_application/powerplatform_admin_management_application.md) | Not in PPAPI |
| powerplatform_billing_policy | [API Documentation](./resource_powerplatform_billing_policy/powerplatform_billing_policy.md) | Ready for implementation |
| powerplatform_billing_policy_environment | [API Documentation](./resource_powerplatform_billing_policy_environment/powerplatform_billing_policy_environment.md) | Ready for implementation |
| powerplatform_connection | [API Documentation](./resource_powerplatform_connection/powerplatform_connection.md) | Not in PPAPI |
| powerplatform_connection_share | [API Documentation](./resource_powerplatform_connection_share/powerplatform_connection_share.md) | Not in PPAPI |
| powerplatform_copilot_studio_application_insights | [API Documentation](./resource_powerplatform_copilot_studio_application_insights/powerplatform_copilot_studio_application_insights.md) | Not in PPAPI |
| powerplatform_data_loss_prevention_policy | [API Documentation](./resource_powerplatform_data_loss_prevention_policy/powerplatform_data_loss_prevention_policy.md) | Not in PPAPI |
| powerplatform_data_record | [API Documentation](./resource_powerplatform_data_record/powerplatform_data_record.md) | WebAPI |
| powerplatform_enterprise_policy | [API Documentation](./resource_powerplatform_enterprise_policy/powerplatform_enterprise_policy.md) | Not in PPAPI |
| powerplatform_environment | [API Documentation](./resource_powerplatform_environment/powerplatform_environment.md) | [Requires changes](./power_platform_api_diff.md#resource_powerplatform_environment) |
| powerplatform_environment_application_admin | [API Documentation](./resource_powerplatform_environment_application_admin/powerplatform_environment_application_admin.md) | Not in PPAPI |
| powerplatform_environment_application_package_install | [API Documentation](./resource_powerplatform_environment_application_package_install/powerplatform_environment_application_package_install.md) | Ready for implementation |
| powerplatform_environment_group | [API Documentation](./resource_powerplatform_environment_group/powerplatform_environment_group.md) | Ready for implementation |
| powerplatform_environment_group_rule_set | [API Documentation](./resource_powerplatform_environment_group_rule_set/powerplatform_environment_group_rule_set.md) | Ready for implementation |
| powerplatform_environment_settings | [API Documentation](./resource_powerplatform_environment_settings/powerplatform_environment_settings.md) | WebAPI |
| powerplatform_environment_wave | [API Documentation](./resource_powerplatform_environment_wave/powerplatform_environment_wave.md) | Not in PPAPI |
| powerplatform_managed_environment | [API Documentation](./resource_powerplatform_managed_environment/powerplatform_managed_environment.md) | [Requires changes](./power_platform_api_diff.md#resource_powerplatform_managed_environment) |
| powerplatform_rest | [API Documentation](./resource_powerplatform_rest/powerplatform_rest.md) | WebAPI |
| powerplatform_solution | [API Documentation](./resource_powerplatform_solution/powerplatform_solution.md) | WebAPI |
| powerplatform_tenant_isolation_policy | [API Documentation](./resource_powerplatform_tenant_isolation_policy/powerplatform_tenant_isolation_policy.md) | Not in PPAPI |
| powerplatform_tenant_settings | [API Documentation](./resource_powerplatform_tenant_settings/powerplatform_tenant_settings.md) | Not in PPAPI |
| powerplatform_user | [API Documentation](./resource_powerplatform_user/powerplatform_user.md) | [Requires changes](./power_platform_api_diff.md#resource_powerplatform_user) |

## Data Sources

| Name | API | Status |
|------|-----|--------|
| powerplatform_analytics_data_exports | [API Documentation](./datasource_powerplatform_analytics_data_exports/powerplatform_analytics_data_exports.md) | Not in PPAPI |
| powerplatform_billing_policies | [API Documentation](./datasource_powerplatform_billing_policies/powerplatform_billing_policies.md) | Ready for implementation |
| powerplatform_billing_policies_environments | [API Documentation](./datasource_powerplatform_billing_policies_environments/powerplatform_billing_policies_environments.md) | Ready for implementation |
| powerplatform_connections | [API Documentation](./datasource_powerplatform_connections/powerplatform_connections.md) | Not in PPAPI |
| powerplatform_connection_shares | [API Documentation](./datasource_powerplatform_connection_shares/powerplatform_connection_shares.md) | Not in PPAPI |
| powerplatform_connectors | [API Documentation](./datasource_powerplatform_connectors/powerplatform_connectors.md) | [Requires changes](./power_platform_api_diff.md#datasource_powerplatform_connectors) |
| powerplatform_currencies | [API Documentation](./datasource_powerplatform_currencies/powerplatform_currencies.md) | WebAPI |
| powerplatform_environment_application_packages | [API Documentation](./datasource_powerplatform_environment_application_packages/powerplatform_environment_application_packages.md) | Ready for implementation |
| powerplatform_environment_powerapps | [API Documentation](./datasource_powerplatform_environment_powerapps/powerplatform_environment_powerapps.md) | [Requires changes](./power_platform_api_diff.md#datasource_powerplatform_environment_powerapps) |
| powerplatform_environment_templates | [API Documentation](./datasource_powerplatform_environment_templates/powerplatform_environment_templates.md) | Not in PPAPI |
| powerplatform_languages | [API Documentation](./datasource_powerplatform_languages/powerplatform_languages.md) | Not in PPAPI |
| powerplatform_locations | [API Documentation](./datasource_powerplatform_locations/powerplatform_locations.md) | Not in PPAPI |
| powerplatform_rest_query | [API Documentation](./datasource_powerplatform_rest_query/powerplatform_rest_query.md) | WebAPI |
| powerplatform_solution_checker_rules | [API Documentation](./datasource_powerplatform_solution_checker_rules/powerplatform_solution_checker_rules.md) | Not in PPAPI |
| powerplatform_solutions | [API Documentation](./datasource_powerplatform_solutions/powerplatform_solutions.md) | WebAPI |
| powerplatform_tenant | [API Documentation](./datasource_powerplatform_tenant/powerplatform_tenant.md) | Not in PPAPI |
| powerplatform_tenant_application_packages | [API Documentation](./datasource_powerplatform_tenant_application_packages/powerplatform_tenant_application_packages.md) | Not in PPAPI |
| powerplatform_tenant_capacity | [API Documentation](./datasource_powerplatform_tenant_capacity/powerplatform_tenant_capacity.md) | Ready for implementation |
| powerplatform_tenant_settings_data_source | [API Documentation](./datasource_powerplatform_tenant_settings_data_source/powerplatform_tenant_settings_data_source.md) | Not in PPAPI |

## Summary

### Resources (23 total)
- **Ready for implementation**: 4 (17%)
- **Requires changes**: 3 (13%)
- **Not in PPAPI**: 12 (52%)
- **WebAPI**: 4 (17%)

### Data Sources (19 total)
- **Ready for implementation**: 4 (21%)
- **Requires changes**: 2 (11%)
- **Not in PPAPI**: 10 (53%)
- **WebAPI**: 3 (16%)

### Overall (42 total)
- **Ready for implementation**: 8 (19%)
- **Requires changes**: 5 (12%)
- **Not in PPAPI**: 22 (52%)
- **WebAPI**: 7 (17%)
