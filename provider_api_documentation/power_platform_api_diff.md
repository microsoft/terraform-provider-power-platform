# Power Platform API Documentation Comparison

This document tracks the differences between the API endpoints documented in this provider and the official [Microsoft Power Platform REST API documentation](https://learn.microsoft.com/en-us/rest/api/power-platform/).

## Summary

Total Resources: 20  
Total Data Sources: 18

## Resources

### resource_powerplatform_admin_management_application

**Provider API Endpoints:**

[See resource_powerplatform_admin_management_application/powerplatform_admin_management_application.md](./resource_powerplatform_admin_management_application/powerplatform_admin_management_application.md)

**Official Microsoft Docs:**

**DOES NOT EXIST IS POWER PLATFORM API**

---

### resource_powerplatform_connection

**Provider API Endpoints:**

[See resource_powerplatform_connection/powerplatform_connection.md](./resource_powerplatform_connection/powerplatform_connection.md)

**Official Microsoft Docs:**

**DOES NOT EXIST IS POWER PLATFORM API**

---

### resource_powerplatform_connection_share

**Provider API Endpoints:**

[See resource_powerplatform_connection_share/powerplatform_connection_share.md](./resource_powerplatform_connection_share/powerplatform_connection_share.md)

**Official Microsoft Docs:**

**DOES NOT EXIST IS POWER PLATFORM API**

---

### resource_powerplatform_copilot_studio_application_insights

**Provider API Endpoints:**

[See resource_powerplatform_copilot_studio_application_insights/powerplatform_copilot_studio_application_insights.md](./resource_powerplatform_copilot_studio_application_insights/powerplatform_copilot_studio_application_insights.md)

**Official Microsoft Docs:**

**DOES NOT EXIST IS POWER PLATFORM API**

---

### resource_powerplatform_data_loss_prevention_policy

**Provider API Endpoints:**

[See resource_powerplatform_data_loss_prevention_policy/powerplatform_data_loss_prevention_policy.md](./resource_powerplatform_data_loss_prevention_policy/powerplatform_data_loss_prevention_policy.md)

**Official Microsoft Docs:**

**DOES NOT EXIST IS POWER PLATFORM API**

---

### resource_powerplatform_data_record

**Provider API Endpoints:**

[See resource_powerplatform_data_record/powerplatform_data_record.md](./resource_powerplatform_data_record/powerplatform_data_record.md)

**Official Microsoft Docs:**

**WEBAPI OUT OF SCOPE**

---

### resource_powerplatform_enterprise_policy

**Provider API Endpoints:**

[See resource_powerplatform_enterprise_policy/powerplatform_enterprise_policy.md](./resource_powerplatform_enterprise_policy/powerplatform_enterprise_policy.md)

**Official Microsoft Docs:**

**DOES NOT EXIST IS POWER PLATFORM API**

---

### resource_powerplatform_environment

**Provider API Endpoints:**

[See resource_powerplatform_environment/powerplatform_environment.md](./resource_powerplatform_environment/powerplatform_environment.md)

**Official Microsoft Docs:**

- [Environments](https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environments)

**Differences:**

1. Create operation Missing
2. Read does not work with SPN accounts:

    ```JSON
    {"code":"PrincipalNotExists","message":"The request cannot be completed because the principal in tenant 1dbbeae5-8fa6-462e-a5a1-9932a520a1dc with object id 8cc077e9-3c82-4550-ba63-1416a149f7fe does not exist."}
    ```

    Missing Read attributes:
    - allow_bing_search (Boolean) Allow Bing search in the environment
allow_moving_data_across_regions (Boolean) Allow moving data across regions
    - billing_policy_id
    - cadence (frequent, moderate)
    - description
    - owner_id (entra id of the developer type environment)
    - release_cycle
    - security_group_id
    - temaplate_metadata (D365 template: FieldService, Sales, etc)
    - templates (List of String) The selected instance provisioning template (if any). ERP based templates.
    - (Read only) enterprise_policies (id, location, status, system_id, type)
    - (Read only) linked_app_id (String) The GUID of the linked D365 application
    - (Read only) linked_app_type (String) The type of the linked D365 application
    - (Read only) linked_app_url (String) The URL of the linked D365 application
    - (Read only) unique_name (String) Unique name of the Dataverse environment

---

### resource_powerplatform_environment_application_admin

**Provider API Endpoints:**

[See resource_powerplatform_environment_application_admin/powerplatform_environment_application_admin.md](./resource_powerplatform_environment_application_admin/powerplatform_environment_application_admin.md)

**Official Microsoft Docs:**

**DOES NOT EXIST IS POWER PLATFORM API**

---

### resource_powerplatform_environment_application_package_install

**Provider API Endpoints:**

[See resource_powerplatform_environment_application_package_install/powerplatform_environment_application_package_install.md](./resource_powerplatform_environment_application_package_install/powerplatform_environment_application_package_install.md)

**Official Microsoft Docs:**

- [Applications - Install Application Package](https://learn.microsoft.com/en-us/rest/api/power-platform/appmanagement/applications/install-application-package)

**Differences:**

**Nothing is missing can be migrated to PPAPI asap**

---

### resource_powerplatform_environment_group

**Provider API Endpoints:**

[See resource_powerplatform_environment_group/powerplatform_environment_group.md](./resource_powerplatform_environment_group/powerplatform_environment_group.md)

**Official Microsoft Docs:**

- [Environment Groups - Create Environment Group](https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-groups/create-environment-group)

**Differences:**

**Nothing is missing can be migrated to PPAPI asap**

---

### resource_powerplatform_environment_group_rule_set

**Provider API Endpoints:**

[See resource_powerplatform_environment_group_rule_set/powerplatform_environment_group_rule_set.md](./resource_powerplatform_environment_group_rule_set/powerplatform_environment_group_rule_set.md)

**Official Microsoft Docs:**

- [Rule Based Policies - Create Rule Based Policy](https://learn.microsoft.com/en-us/rest/api/power-platform/governance/rule-based-policies/create-rule-based-policy#ruleset)

**Nothing is missing can be migrated to PPAPI asap**

---

### resource_powerplatform_environment_settings

**Provider API Endpoints:**

[See resource_powerplatform_environment_settings/powerplatform_environment_settings.md](./resource_powerplatform_environment_settings/powerplatform_environment_settings.md)

**Official Microsoft Docs:**

**WEBAPI OUT OF SCOPE**

---

### resource_powerplatform_environment_wave

**Provider API Endpoints:**

[See resource_powerplatform_environment_wave/powerplatform_environment_wave.md](./resource_powerplatform_environment_wave/powerplatform_environment_wave.md)

**Official Microsoft Docs:**

**DOES NOT EXIST IS POWER PLATFORM API**

---

### resource_powerplatform_managed_environment

**Provider API Endpoints:**

[See resource_powerplatform_managed_environment/powerplatform_managed_environment.md](./resource_powerplatform_managed_environment/powerplatform_managed_environment.md)

**Official Microsoft Docs:**

- [Environment Managed Governance - Enable Managed Environment](https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-managed-governance/enable-managed-environment)
- [Environment Management Settings - Create Environment Management Settings](https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-management-settings/create-environment-management-settings)

**Differences:**

```
https://api.powerplatform.com/environmentmanagement/environments/3105cd9a-ab72-f011-8589-000d3a5b5d2b/settings?api-version=2022-03-01-preview

{"code":"AuthorizationHeaderInvalid","message":"The provided Authorization header is invalid: S2S17001: SAL was able to validate the protocol, but validation failed as none of the inbound policies were satisfied. Validation failures: \u0027Unified: InvalidSignature\u0027."}
```

Question: Does this API allow to set managed environment settings? (<https://registry.terraform.io/providers/microsoft/power-platform/latest/docs/resources/managed_environment>)

---

### resource_powerplatform_rest

**Provider API Endpoints:**

[See resource_powerplatform_rest/powerplatform_rest.md](./resource_powerplatform_rest/powerplatform_rest.md)

**Official Microsoft Docs:**

**WEBAPI OUT OF SCOPE**

---

### resource_powerplatform_solution

**Provider API Endpoints:**

[See resource_powerplatform_solution/powerplatform_solution.md](./resource_powerplatform_solution/powerplatform_solution.md)

**Official Microsoft Docs:**

**WEBAPI OUT OF SCOPE**

---

### resource_powerplatform_tenant_isolation_policy

**Provider API Endpoints:**

[See resource_powerplatform_tenant_isolation_policy/powerplatform_tenant_isolation_policy.md](./resource_powerplatform_tenant_isolation_policy/powerplatform_tenant_isolation_policy.md)

**Official Microsoft Docs:**

**DOES NOT EXIST IS POWER PLATFORM API**

---

### resource_powerplatform_tenant_settings

**Provider API Endpoints:**

[See resource_powerplatform_tenant_settings/powerplatform_tenant_settings.md](./resource_powerplatform_tenant_settings/powerplatform_tenant_settings.md)

**Official Microsoft Docs:**

**DOES NOT EXIST IS POWER PLATFORM API**

---

### resource_powerplatform_user

**Provider API Endpoints:**

[See resource_powerplatform_user/powerplatform_user.md](./resource_powerplatform_user/powerplatform_user.md)

**Official Microsoft Docs:**

- [Users - Apply Admin Role](https://learn.microsoft.com/en-us/rest/api/power-platform/usermanagement/users/apply-admin-role)

**Differences:**

**PPAPI only alllows to assign admn roles**

---

## Data Sources

### datasource_powerplatform_analytics_data_exports

**Provider API Endpoints:**

[See datasource_powerplatform_analytics_data_exports/powerplatform_analytics_data_exports.md](./datasource_powerplatform_analytics_data_exports/powerplatform_analytics_data_exports.md)

**Official Microsoft Docs:**

- TBD

**Differences:**

- TBD

---

### datasource_powerplatform_connections

**Provider API Endpoints:**

[See datasource_powerplatform_connections/powerplatform_connections.md](./datasource_powerplatform_connections/powerplatform_connections.md)

**Official Microsoft Docs:**

- TBD

**Differences:**

- TBD

---

### datasource_powerplatform_connection_shares

**Provider API Endpoints:**

[See datasource_powerplatform_connection_shares/powerplatform_connection_shares.md](./datasource_powerplatform_connection_shares/powerplatform_connection_shares.md)

**Official Microsoft Docs:**

- TBD

**Differences:**

- TBD

---

### datasource_powerplatform_connectors

**Provider API Endpoints:**

[See datasource_powerplatform_connectors/powerplatform_connectors.md](./datasource_powerplatform_connectors/powerplatform_connectors.md)

**Official Microsoft Docs:**

- TBD

**Differences:**

- TBD

---

### datasource_powerplatform_currencies

**Provider API Endpoints:**

[See datasource_powerplatform_currencies/powerplatform_currencies.md](./datasource_powerplatform_currencies/powerplatform_currencies.md)

**Official Microsoft Docs:**

- TBD

**Differences:**

- TBD

---

### datasource_powerplatform_data_loss_prevention_policies

**Provider API Endpoints:**

[See datasource_powerplatform_data_loss_prevention_policies/powerplatform_data_loss_prevention_policies.md](./datasource_powerplatform_data_loss_prevention_policies/powerplatform_data_loss_prevention_policies.md)

**Official Microsoft Docs:**

- TBD

**Differences:**

- TBD

---

### datasource_powerplatform_environment_application_packages

**Provider API Endpoints:**

[See datasource_powerplatform_environment_application_packages/powerplatform_environment_application_packages.md](./datasource_powerplatform_environment_application_packages/powerplatform_environment_application_packages.md)

**Official Microsoft Docs:**

- TBD

**Differences:**

- TBD

---

### datasource_powerplatform_environment_powerapps

**Provider API Endpoints:**

[See datasource_powerplatform_environment_powerapps/powerplatform_environment_powerapps.md](./datasource_powerplatform_environment_powerapps/powerplatform_environment_powerapps.md)

**Official Microsoft Docs:**

- TBD

**Differences:**

- TBD

---

### datasource_powerplatform_environment_templates

**Provider API Endpoints:**

[See datasource_powerplatform_environment_templates/powerplatform_environment_templates.md](./datasource_powerplatform_environment_templates/powerplatform_environment_templates.md)

**Official Microsoft Docs:**

- TBD

**Differences:**

- TBD

---

### datasource_powerplatform_languages

**Provider API Endpoints:**

[See datasource_powerplatform_languages/powerplatform_languages.md](./datasource_powerplatform_languages/powerplatform_languages.md)

**Official Microsoft Docs:**

- TBD

**Differences:**

- TBD

---

### datasource_powerplatform_locations

**Provider API Endpoints:**

[See datasource_powerplatform_locations/powerplatform_locations.md](./datasource_powerplatform_locations/powerplatform_locations.md)

**Official Microsoft Docs:**

- TBD

**Differences:**

- TBD

---

### datasource_powerplatform_rest_query

**Provider API Endpoints:**

[See datasource_powerplatform_rest_query/powerplatform_rest_query.md](./datasource_powerplatform_rest_query/powerplatform_rest_query.md)

**Official Microsoft Docs:**

- TBD

**Differences:**

- TBD

---

### datasource_powerplatform_solution_checker_rules

**Provider API Endpoints:**

[See datasource_powerplatform_solution_checker_rules/powerplatform_solution_checker_rules.md](./datasource_powerplatform_solution_checker_rules/powerplatform_solution_checker_rules.md)

**Official Microsoft Docs:**

- TBD

**Differences:**

- TBD

---

### datasource_powerplatform_solutions

**Provider API Endpoints:**

[See datasource_powerplatform_solutions/powerplatform_solutions.md](./datasource_powerplatform_solutions/powerplatform_solutions.md)

**Official Microsoft Docs:**

- TBD

**Differences:**

- TBD

---

### datasource_powerplatform_tenant

**Provider API Endpoints:**

[See datasource_powerplatform_tenant/powerplatform_tenant.md](./datasource_powerplatform_tenant/powerplatform_tenant.md)

**Official Microsoft Docs:**

- TBD

**Differences:**

- TBD

---

### datasource_powerplatform_tenant_application_packages

**Provider API Endpoints:**

[See datasource_powerplatform_tenant_application_packages/powerplatform_tenant_application_packages.md](./datasource_powerplatform_tenant_application_packages/powerplatform_tenant_application_packages.md)

**Official Microsoft Docs:**

- TBD

**Differences:**

- TBD

---

### datasource_powerplatform_tenant_capacity

**Provider API Endpoints:**

[See datasource_powerplatform_tenant_capacity/powerplatform_tenant_capacity.md](./datasource_powerplatform_tenant_capacity/powerplatform_tenant_capacity.md)

**Official Microsoft Docs:**

- TBD

**Differences:**

- TBD

---

### datasource_powerplatform_tenant_settings_data_source

**Provider API Endpoints:**

[See datasource_powerplatform_tenant_settings_data_source/powerplatform_tenant_settings_data_source.md](./datasource_powerplatform_tenant_settings_data_source/powerplatform_tenant_settings_data_source.md)

**Official Microsoft Docs:**

- TBD

**Differences:**

- TBD

---
