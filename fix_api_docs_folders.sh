#!/bin/bash

# Script to rename subfolders with correct "resource_" or "datasource_" prefix
# Based on actual analysis of file contents

cd /workspaces/terraform-provider-power-platform/provider_api_documentation || exit 1

# Define resources (should have resource_ prefix)
declare -a resources=(
    "powerplatform_admin_management_application"
    "powerplatform_analytics_data_exports"
    "powerplatform_connection"
    "powerplatform_connection_share"
    "powerplatform_copilot_studio_application_insights"
    "powerplatform_data_loss_prevention_policy"
    "powerplatform_data_record"
    "powerplatform_enterprise_policy"
    "powerplatform_environment"
    "powerplatform_environment_application_admin"
    "powerplatform_environment_application_package_install"
    "powerplatform_environment_group"
    "powerplatform_environment_group_rule_set"
    "powerplatform_managed_environment"
    "powerplatform_rest"
    "powerplatform_solution"
    "powerplatform_tenant_isolation_policy"
    "powerplatform_tenant_settings"
    "powerplatform_user"
)

# Define data sources (should have datasource_ prefix)
declare -a datasources=(
    "powerplatform_connection_shares"
    "powerplatform_connections"
    "powerplatform_connectors"
    "powerplatform_currencies"
    "powerplatform_data_loss_prevention_policies"
    "powerplatform_environment_application_packages"
    "powerplatform_environment_powerapps"
    "powerplatform_environment_settings"
    "powerplatform_environment_templates"
    "powerplatform_environment_wave"
    "powerplatform_languages"
    "powerplatform_locations"
    "powerplatform_rest_query"
    "powerplatform_solution_checker_rules"
    "powerplatform_solutions"
    "powerplatform_tenant"
    "powerplatform_tenant_application_packages"
    "powerplatform_tenant_capacity"
    "powerplatform_tenant_settings_data_source"
)

# Rename resources
echo "=== Renaming RESOURCES ==="
for name in "${resources[@]}"; do
    old_dir="datasource_${name}"
    new_dir="resource_${name}"
    
    if [ -d "$old_dir" ]; then
        mv "$old_dir" "$new_dir"
        echo "Renamed $old_dir to $new_dir"
    elif [ -d "$new_dir" ]; then
        echo "Skipped $name (already has resource_ prefix)"
    else
        echo "Warning: $old_dir not found"
    fi
done

# Keep data sources (already have correct prefix)
echo -e "\n=== Checking DATA SOURCES ==="
for name in "${datasources[@]}"; do
    dir="datasource_${name}"
    
    if [ -d "$dir" ]; then
        echo "OK: $dir (already correct)"
    else
        echo "Warning: $dir not found"
    fi
done

echo -e "\n=== Renaming complete! ==="
