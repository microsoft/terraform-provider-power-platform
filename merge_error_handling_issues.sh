#!/bin/bash

# Get the list of all .md files in the error handling directory
files=(
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/01_missing_return_after_error/datasource_connectors_missing-return-on-error-high.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/01_missing_return_after_error/datasource_solutions_dataverse_exists_return_high.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/01_missing_return_after_error/datasource_tenant_application_packages_error_handling_high.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/03_nil_pointer_validation/api_application_Missing_Nil_Check_for_Header_Value_medium.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/03_nil_pointer_validation/models_null_unknown_checks_low.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/03_nil_pointer_validation/resource_user.go-nil_handling-high.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/04_diagnostic_usage/datasource_connection_shares_api-error-detail-medium.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/04_diagnostic_usage/helpers_improper-diagnostic-usage-high.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/04_diagnostic_usage/resource_data_record_Inconsistent_Error_Handling_in_convertColumnsToState_high.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/04_diagnostic_usage/resources_environment_settings_misleading_error_providerdata_type_high.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/05_json_marshal_unmarshal/api_languages_Potentially_Unhandled_JSON_Unmarshal_medium.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/05_json_marshal_unmarshal/datasource_connections.go-json_marshal-medium.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/05_json_marshal_unmarshal/datasource_languages_type_assertion_high.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/05_json_marshal_unmarshal/request_marshallto_decode_pointer_high.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/06_timeout_and_polling/api_environment_wave_polling_loop_without_timeout_high.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/07_context_and_defer/resource_billing_policy_context_leak_low.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/07_context_and_defer/resource_managed_environment_exit_context_always_deferred_medium.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/07_context_and_defer/resource_user.go-exit_context_defer-medium.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/08_data_structure_issues/datasource_analytics_data_exports_go_structure_medium.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/08_data_structure_issues/datasource_tenant_settings_schema_structure_medium.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/08_data_structure_issues/dto.go_conversion_duplication_medium.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/08_data_structure_issues/resource_environment_wave.go-feature_state-mapping-medium.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/09_validation_and_modifiers/dynamic_columns_validator-unchecked_getattr_error-medium.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/09_validation_and_modifiers/make_field_required_when_other_field_does_not_have_value_validator_Diagnostic_Message_Is_Not_Actionable_Low.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/09_validation_and_modifiers/other_field_required_when_value_of_validator_structure_low_2.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/09_validation_and_modifiers/set_bool_value_unknown_if_checksum_change_modifier_issue_wrong_hash_message_low.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/09_validation_and_modifiers/set_string_value_unknown_if_checksum_change_modifier-inadequate-diagnostic-context-for-attribute-in-checksum-calculation-medium.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/09_validation_and_modifiers/set_string_value_unknown_if_checksum_change_modifier_redundant-diagnostic-appending-medium.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/09_validation_and_modifiers/set_string_value_unknown_if_checksum_change_modifier_undescriptive-error-message-in-diagnostic-low.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/09_validation_and_modifiers/sync_attribute_plan_modifier_empty_checksum_medium.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/09_validation_and_modifiers/sync_attribute_plan_modifier_incorrect_checksum_message_low.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/10_api_response_handling/api_locations_structure_low.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/10_api_response_handling/api_managed_environment_low.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/11_generic_error_messages/datasource_environment_templates_error_message_configure_low.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/11_generic_error_messages/resource_tenant_isolation_policy_go_api_err_message_low.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/11_generic_error_messages/resource_tenant_settings.go-generic_error_titles-low.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/12_configuration_and_constants/client_go_magic_strings_low.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/12_configuration_and_constants/constants.go-incorrect_domain-medium.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/12_configuration_and_constants/provider_go_auxiliarytenantids_type_issue_medium.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/13_type_safety_and_formatting/resource_billing_policy_empty_string_medium.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/13_type_safety_and_formatting/url_format_error_interface_assertion_low.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/13_type_safety_and_formatting/uuid_value_string_formatting_low.md"
    "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/14_test_and_mock_issues/hash_test_structure_low.md"
)

# Output file
output_file="/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/organized_issues/Error handling/_merged_error_handling_issues.md"

# Create the header of the merged file
cat > "$output_file" << 'EOF'
# Terraform Provider Power Platform - Error Handling Issues

This document consolidates all identified error handling issues across the Terraform Provider for Power Platform. These issues have been categorized and prioritized to help improve the robustness and reliability of the provider.

---

EOF

# Counter for issue numbering
counter=1

# Loop through each file and append its content
for file in "${files[@]}"; do
    if [[ -f "$file" ]]; then
        echo "## ISSUE $counter" >> "$output_file"
        echo "" >> "$output_file"
        
        # Remove the markdown code block wrapper and add the content
        sed '1d;$d' "$file" >> "$output_file"
        
        echo "" >> "$output_file"
        echo "---" >> "$output_file"
        echo "" >> "$output_file"
        
        ((counter++))
    fi
done

# Add the final instructions
cat >> "$output_file" << 'EOF'
# To finish the task you have to 
1. Run linter and fix any issues 
2. Run UnitTest and fix any of failing ones
3. Generate docs 
4. Run Changie

# Changie Instructions
Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```
Where:
- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
EOF

echo "Merged file created: $output_file"
echo "Total issues merged: $((counter-1))"
