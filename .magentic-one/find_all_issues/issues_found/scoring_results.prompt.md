# PROMPT

The following directory has a list of issue found in the repository, read all of them and give them a significance number between 1-10 where the is the most valuable for fixing and 1 is not important or false positive.

- Return result as a table with 4 colums: Order Number, Score, Markdown File Name, Description. Table should be in markdown string so it is easy to copy.
- Be sure to score all the files

the files are in the following folder: /workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/critical

## Result

### Critical Issues

| Order Number | Score | Markdown File Name | Description |
|-------------|-------|-------------------|-------------|
| 1 | 10 | api_rest_panic_issue.md | Potential panic in REST API client due to unchecked error conditions which could crash the provider without proper recovery mechanisms |
| 2 | 10 | string_to_set_improper_usage_of_panic.md | Direct usage of panic in code path instead of proper error handling, causing unrecoverable application termination |
| 3 | 9 | api_licensing_infinite_retry.md | Infinite retry logic in API client could cause application hang and resource exhaustion, with no retry limit implemented |
| 4 | 9 | resource_environment_missing_error_check_create.md | Missing error checks in environment creation flow could lead to silent failures and inconsistent state |
| 5 | 9 | resource_environment_missing_error_check_read.md | Missing validation in environment read operations could cause inconsistent state |
| 6 | 9 | uninitialized_connections_client.md | Use of uninitialized client could cause runtime crashes |
| 7 | 9 | cert_potential_panic_due_to_unchecked_return_values.md | Unchecked certificate operations may cause panic in authentication flow |
| 8 | 8 | provider_configure_missing_error_handling.md | Insufficient error handling in provider configuration could prevent proper initialization |
| 9 | 8 | Configure-Lack-Error-Handling.md | Missing error handling for client initialization which could lead to nil pointer exceptions |
| 10 | 8 | api_dlp_policy_Critical_Issue_1.md | Improper error handling in DLP policy API calls leading to potential state corruption |
| 11 | 8 | resource_data_record_columns_error.md | Improper handling of column data could lead to data corruption or loss |
| 12 | 8 | datasource_data_record_critical_issue_1.md | Potential for nil pointer dereference when handling data records |
| 13 | 8 | datasource_currencies_read_null_pointer.md | Null pointer dereference risk in currencies data source read operations |
| 14 | 8 | provider_configure_null_handling_fix.md | Improper null check in provider configuration could lead to runtime errors |
| 15 | 8 | critical_Type_Assertion_Validation_Configure_Method.md | Unsafe type assertion without validation could cause panic |
| 16 | 8 | unexpected_providerdata_type.md | Incorrect type assertion could lead to runtime errors |
| 17 | 8 | resource_environment_application_package_install_1.md | Unhandled errors in application package installation flow |
| 18 | 7 | other_field_required_when_value_of_validator_regex_null_check.md | Validator missing null checks could lead to validation bypass |
| 19 | 7 | api_solution_issue_missing_import_completion_mechanism.md | Incomplete solution import mechanism could lead to partial deployments |
| 20 | 7 | resource_tenant_settings_delete_critical.md | Improper deletion flow in tenant settings could leave orphaned resources |
| 21 | 7 | api_environment_DeleteEnvironment_Uninitialized_Type.md | Uninitialized type in environment deletion could cause crashes |
| 22 | 7 | unbounded-for-loop-installApplication.md | Unbounded loop in application installation could cause resource exhaustion |
| 23 | 7 | api_environment_settings_issue_1.md | Improper error handling in environment settings API |
| 24 | 7 | concurrent-goroutine-management-for-timeout-cancellation.md | Improper context handling could lead to goroutine leaks |
| 25 | 7 | api_environment_settings_issue_2.md | Missing validation in environment settings API |
| 26 | 7 | improper_error_handling_execute_method.md | Improper error handling in API execute method could mask failures |
| 27 | 7 | environment_settings_resource_create_function_dataverse_check.md | Missing validation for Dataverse environment requirements |
| 28 | 7 | lifecycle_missing_context_cancellation.md | Missing context cancellation could lead to resource leaks |
| 29 | 7 | error_handling_in_read_function.md | Insufficient error checking in REST execute operations |
| 30 | 6 | set_string_value_unknown_if_checksum_change_modifier_slice_access.md | Unsafe slice access in modifiers could cause out-of-bounds panics |
| 31 | 6 | require_replace_object_to_empty_modifier_missing_error_handling.md | Missing error handling in resource replacement logic |
| 32 | 6 | resource_environment_missing_validation_update.md | Inadequate validation in environment update operations |
| 33 | 6 | resource_connection_test_missing_comma.md | Syntax error in test file preventing proper test execution |
| 34 | 6 | datasource_connectors_issue_nil_connectors.md | Missing nil check for connectors could cause crashes |
| 35 | 6 | resource_dlp_policy_test_incomplete_error_handling.md | Incomplete error handling in DLP policy tests |
| 36 | 6 | models_missing_initialization_error_handling.md | Missing error handling during model initialization |
| 37 | 6 | force_string_value_unknown_modifier_critical_error_handling.md | Missing error handling in string value modifier |
| 38 | 5 | unexpected_http_status_code_error_issue_exposing_http_body.md | Potential for leaking sensitive information in error messages |
| 39 | 5 | models-missing-validations.md | Missing input validations across multiple model structures |
| 40 | 5 | httpmock_file_issue_in_mocks.md | Issues with HTTP mock setup affecting test reliability |
| 41 | 5 | api_dlp_policy_Critical_Issue_3.md | Improper API call construction in DLP policy management |
| 42 | 5 | api_dlp_policy_Critical_Issue_2.md | Missing error checks in DLP policy API client |
| 43 | 5 | request_resource_leak_defer.md | Resource leak due to improper defer usage |
| 44 | 5 | auth_getAssertion_HTTP_client_handling.md | Improper HTTP client handling in authentication |
| 45 | 5 | models.go_inconsistent_error_checking.md | Inconsistent error checking patterns across model code |
| 46 | 5 | models.go_invalid_error_propagation.md | Improper error propagation in model conversions |
| 47 | 5 | models-undefined-client.md | Use of undefined client references |
| 48 | 5 | metadata_context_initialization_check.md | Missing initialization checks for metadata context |
| 49 | 4 | datasource_locations_input_validation.md | Missing input validation for location data sources |
| 50 | 4 | api_languages_improper_error_handling.md | Suboptimal error handling in languages API client |
| 51 | 4 | resource_billing_policy_environment_context_misuse.md | Improper context usage in billing policy resource |
| 52 | 4 | resource_environment_application_package_install_2.md | Incomplete validation in application package installation |
| 53 | 4 | resource_copilot_studio_application_insights_issue_4.md | Inconsistent error handling in Copilot Studio insights |
| 54 | 4 | datasource_tenant_read_function_unverified_tenant_variable.md | Missing verification of tenant variables |
| 55 | 4 | models-no-client-implementation.md | Missing client implementation in model interfaces |
| 56 | 4 | issue1_context_async_usage.md | Improper async context usage |
| 57 | 4 | excessive_dependencies_on_external_resources.md | Over-reliance on external dependencies |
| 58 | 3 | api_connection_invalid_query_param_format.md | Incorrect query parameter formatting in connection API |
| 59 | 3 | dynamic_columns_issue_2.md | Inconsistent handling of dynamic columns |
| 60 | 3 | constants_UUID.md | Non-compliance with UUID format standards |
| 61 | 3 | resource_solution_issue_1.md | Inconsistent state handling in solution resource |
| 62 | 3 | dto.md | Inconsistent DTO structure definitions |
| 63 | 3 | datasource_connections_test.go_issue_invalid_json.md | Invalid JSON in test fixtures |
| 64 | 3 | datasource_billing_policies_environments_test_issue1.md | Inconsistent test verification for billing policies |
| 65 | 3 | dto_go_issue_6.md | Minor inconsistencies in DTO models |
| 66 | 3 | Lack_of_Validation_for_Struct_Fields.md | General lack of validation for struct fields |
| 67 | 3 | issue_1.md | Generic issue with minor impact |
| 68 | 3 | hardcoded_ids_in_resource_configuration.md | Hardcoded IDs in resource configurations |
| 69 | 2 | TestUnitEnvironmentApplicationPackagesDataSource_Validate_Read_Invalid_URLs.md | Imprecise error validation in application package tests |
| 70 | 2 | datasource_connectors_issue_exit_context.md | Missing context cleanup in connectors data source |
| 71 | 2 | api_currencies_missing_unit_tests.md | Lack of unit tests for currencies API client |
| 72 | 2 | provider_test_issue_2.md | Minor issues in provider test setup |
| 73 | 2 | datasource_solutions_test_issue1.md | Inconsistent test assertions in solutions tests |
| 74 | 2 | api_dlp_policy_Critical_Issue_4.md | Minor inconsistencies in DLP policy API implementation |
| 75 | 2 | other_field_required_when_value_of.md | Minor issue with field validation implementation |
| 76 | 2 | Lack_of_Unit_Tests.md | General lack of unit tests across codebase |
| 77 | 2 | models-go-null-pointer-dereference.md | Minor null pointer risk with low probability |
| 78 | 1 | resource_data_record.md | Documentation inconsistencies in data record resource |
| 79 | 1 | infinite_retry_logic.md | Low-impact issue with retry logic |
| 80 | 1 | models_validation.md | Minor model validation inconsistencies |
| 81 | 1 | models_missing_validation.md | Inconsistencies in model validations with minimal impact |
| 82 | 1 | critical.md | Improper error handling in application packages unit test |
| 83 | 1 | resource_environment_issue_unexpected_providerdata.md | Minor type assertion issue in environment resource |
| 84 | 1 | lifecycle_missing_context_cancellation.md | Low-impact context cleanup issue |
| 85 | 1 | TestUnitEnvironmentApplicationPackagesDataSource_Validate_Read_Invalid_URLs.md | Minor test validation improvement opportunity |

### High Issues

| Order Number | Score | Markdown File Name | Description |
|-------------|-------|-------------------|-------------|
| 1 | 8 | Update-Null-Dereference.md | Potential null dereference in REST resource Update function could cause runtime crash during API operations |
| 2 | 8 | potential_nil_pointer_dereference_in_read.md | Risk of nil pointer dereference in Read operation that could crash the provider |
| 3 | 8 | potential_null_pointer_dereference.md | Unchecked null pointer access in core functionality that could lead to runtime panic |
| 4 | 8 | api_environment_Infinite_Loop_AddDataverseToEnvironment.md | Possible infinite loop in AddDataverseToEnvironment function causing resource exhaustion |
| 5 | 8 | api_rest_null_pointer_scope.md | REST API null pointer risk due to improper scope management |
| 6 | 7 | config_null_pointer_dereference.md | Configuration component has potential null pointer dereference that could crash provider initialization |
| 7 | 7 | datasource_connectors_issue_undefined_listdatasourcemodel.md | Undefined model reference in connectors data source could cause unexpected behavior |
| 8 | 7 | models.go_incorrect_type_assertions.md | Unsafe type assertions without verification in models.go that could lead to panics |
| 9 | 7 | datasource_tenant_capacity-invalid_api_client.md | Invalid API client usage in tenant capacity data source could cause operations to fail |
| 10 | 7 | resource_managed_environment_issue_unexpected_provider_data.md | Incorrect type handling in managed environment resource |
| 11 | 7 | models.go_unhandled_null_conditions.md | Critical null value conditions not properly handled in model conversions |
| 12 | 7 | api_solution_issue_missing_nil_checks.md | Missing nil checks in solution API client could lead to runtime crashes |
| 13 | 7 | Create-Improper-Error-Handling.md | Inadequate error handling in resource creation that could mask critical failures |
| 14 | 7 | resource_environment_wave_nil_value_check.md | Missing nil value check in environment wave resource could cause panic |
| 15 | 7 | api_tenant_settings_nil_pointer_filterdto.md | Nil pointer risk in tenant settings filter DTO handling |
| 16 | 7 | datasource_billing_policies_environments_unverified_cast.md | Unverified type casting in billing policies environments data source |
| 17 | 7 | resource_admin_management_application_create_error.md | Improper error handling in admin application creation could lead to orphaned resources |
| 18 | 6 | api_dlp_policy_High_Issue_1.md | Improper error propagation in DLP policy API client |
| 19 | 6 | api_dlp_policy_High_Issue_2.md | Inconsistent error handling pattern in DLP policy client methods |
| 20 | 6 | api_licensing_invalid_error_handling.md | Invalid error handling approach in licensing API that could mask failures |
| 21 | 6 | api_licensing_missing_pagination.md | Missing pagination handling could lead to incomplete results in licensing operations |
| 22 | 6 | api_data_record_misuse_of_fmtsprintf_with_unescaped_characters_url_construction.md | Potential for URL injection due to unescaped characters in URL construction |
| 23 | 6 | api_currencies_improper_defer_handling.md | Resource leaks due to improper defer handling in currencies API |
| 24 | 6 | api_languages_json_unmarshalling_error.md | Missing error handling for JSON unmarshalling in languages API |
| 25 | 6 | resource_billing_policy_improper_error_handling.md | Improper error handling in billing policy resource that could lead to inconsistent state |
| 26 | 6 | Find_generic_handling.md | Generic error handling without proper context in core functionality |
| 27 | 6 | request_improper_error_handling.md | Insufficient error handling in request processing could lead to misleading error messages |
| 28 | 6 | lack_error_type_handling.md | Missing error type differentiation that could lead to incorrect error responses |
| 29 | 6 | unhandled_json_marshal.md | Unhandled JSON marshaling errors that could lead to incomplete operations |
| 30 | 6 | unhandled_json_marshall.md | Missing error handling for JSON marshalling operations |
| 31 | 6 | resource_tenant_settings_create_high.md | Improper state handling in tenant settings creation |
| 32 | 6 | resource_tenant_settings_configure_high.md | Configuration issues in tenant settings resource |
| 33 | 6 | convertToDto_error_handling_issue.md | Insufficient error handling in DTO conversion functions |
| 34 | 6 | api_locations_missing_error_wrapping.md | Missing context in error messages from locations API client |
| 35 | 6 | HardcodingURLs.md | Hardcoded URLs throughout the codebase make environment configuration difficult |
| 36 | 5 | datasource_environment_powerapps_test_missing_error_handling.md | Incomplete error handling in Power Apps test data source |
| 37 | 5 | lack_of_validation_for_url_input.md | Missing URL input validation could allow malformed requests |
| 38 | 5 | resource_enterprise_policy_create_retry.md | Inconsistent retry mechanism in enterprise policy creation |
| 39 | 5 | api_connectors_error_handling.md | Suboptimal error handling in connectors API client |
| 40 | 5 | resource_solution_issue_2.md | State inconsistency issues in solution resource |
| 41 | 5 | models_error_handling.md | Inconsistent error handling patterns in model operations |
| 42 | 5 | other_field_required_when_value_of_validator_missing_error_handling.md | Missing error handling in validator that could lead to validation bypass |
| 43 | 5 | improper-error-handling-in-enterTimeoutContext.md | Timeout context error handling issues could lead to resource leaks |
| 44 | 5 | hash_error_propagation.md | Inadequate error propagation in hash functions |
| 45 | 5 | datasource_analytics_data_exports_issue_null_check.md | Missing null check in analytics data exports could cause panic |
| 46 | 5 | datasource_data_record_high_issue_2.md | Improper data handling in data record data source |
| 47 | 5 | datasource_connectors_issue_type_assertion.md | Unsafe type assertion in connectors data source |
| 48 | 5 | api_environment_settings_issue_3.md | Incomplete validation in environment settings API |
| 49 | 5 | api_tenant_isolation_policy_Lack_of_Input_Validation_for_tenantId.md | Missing tenant ID validation in isolation policy API |
| 50 | 5 | resource_data_record_dynamic_validator.md | Insufficient validation in dynamic data record fields |
| 51 | 5 | resource_environment_group_rule_set_issue_1.md | Inconsistent error handling in environment group rule set |
| 52 | 5 | resource_environment_group_issue_redundant_error_calls.md | Redundant error handling causing confusing diagnostics |
| 53 | 5 | resource_data_record_test-missing-error-handling.md | Missing error handling in data record tests |
| 54 | 5 | resource_environment_group_rule_set_issue_3.md | State inconsistency in environment group rule sets |
| 55 | 5 | resource_copilot_studio_application_insights_issue_1.md | Error handling inconsistency in Copilot Studio insights |
| 56 | 5 | resource_billing_policy_environment_error_handling.md | Suboptimal error handling in billing policy environment resource |
| 57 | 5 | models_error_context.md | Missing error context in model operations |
| 58 | 5 | force_string_value_unknown_modifier_exit_condition_issue.md | Potential infinite loop in string modifier |
| 59 | 5 | datasource_environments_incorrect_error_handling.md | Improper error handling in environments data source |
| 60 | 5 | Overly_Permissive_UrlValidation.md | URL validation too permissive, allowing potentially malformed URLs |
| 61 | 5 | Misleading_Error_Message_DeleteEnvironmentGroupRuleSet.md | Error messages that don't accurately describe the actual issue |
| 62 | 5 | api_tenant_settings_lack_of_tests_filterdto.md | Insufficient test coverage for tenant settings filter |
| 63 | 5 | api_environment_group_improper_error_handling.md | Inadequate error handling in environment group operations |
| 64 | 5 | models_high_missing_field_validation.md | Critical fields lacking validation in model structures |
| 65 | 5 | api_enterprise_policy_retry_mechanism.md | Inefficient retry mechanism in enterprise policy API |
| 66 | 5 | issue_3_resource_billing_policy_test.md | Inconsistent test behavior in billing policy tests |
| 67 | 5 | datasource_environments_incorrect_providerdata_handling.md | Type assertion issues in environments data source |
| 68 | 5 | unhandled_errors_in_GetSolutionCheckerRules_Function_Call.md | Unchecked errors in solution checker rules client |
| 69 | 5 | resource_solution_test_improper_error_handling.md | Inadequate error assertions in solution resource tests |
| 70 | 5 | api_copilot_studio_application_insights_http_status_handling.md | Missing HTTP status code validation |
| 71 | 4 | datasource_environments_test__error_handling_mock.md | Incomplete mock response handling in environments tests |
| 72 | 4 | datasource_tenant_test_missing_error_handling.md | Missing error validation in tenant data source tests |
| 73 | 4 | Missing_Validation_of_Input_Parameters.md | Input parameters in multiple components lack validation |
| 74 | 4 | datasource_tenant_application_packages_client_initialization.md | Potential nil client in tenant application packages data source |
| 75 | 4 | datasource_tenant_application_packages_read_error_handling.md | Inconsistent error handling in application packages read operations |
| 76 | 4 | resource_billing_policy_environment_test_hardcoded_ids_issue.md | Hardcoded test IDs that could make tests brittle |
| 77 | 4 | resource_connection_share_test_missing_error_handling.md | Incomplete error handling in connection share tests |
| 78 | 4 | UserResource_ErrorHandling.md | Suboptimal error handling pattern in user resource |
| 79 | 4 | map_initialization_issue_in_mocks.md | Potential nil map access in mock implementations |
| 80 | 4 | missing_environment_id_validation_resource_user.md | Missing validation for environment ID in user resource |
| 81 | 4 | models_missing_error_handling_for_client.md | Insufficient error handling for client operations in models |
| 82 | 4 | require_replace_int_attribute_modifier.md | Potential issue with int attribute replacement modifier |
| 83 | 4 | datasource_locations_test_missing_error_handling_for_http_responses.md | Missing HTTP response validation in locations tests |
| 84 | 4 | auth_error_wrapping_getAssertion.md | Insufficient error context in authentication assertions |
| 85 | 4 | datasource_environment_application_packages_missing_environment_id_validation.md | Missing validation for environment ID in application packages |
| 86 | 4 | missing_validation_dto_fields.md | Insufficient validation in DTO field definitions |
| 87 | 4 | requires_replace_string_from_non_empty_modifier_Limited_error_handling_in_PlanModifyString.md | Limited error handling in string replacement modifier |
| 88 | 4 | models_function_error_handling.md | Inconsistent error handling in model conversion functions |
| 89 | 4 | resource_tenant_isolation_policy_test_irregular_error_handling.md | Irregular error handling in isolation policy tests |
| 90 | 4 | resource_rest_test_01_hardcoded_metadata_url.md | Hardcoded URLs in REST resource tests |
| 91 | 4 | Missing_JSON_Tags_in_linkedEnvironmentIdMetadataDto_Struct.md | Missing JSON tags in DTO struct could cause serialization issues |
| 92 | 4 | resource_tenant_isolation_policy_modularization.md | Insufficient modularization in isolation policy resource |
| 93 | 4 | resource_billing_policy_environment_test_randIntn_issue.md | Non-deterministic test behavior due to random values |
| 94 | 4 | resource_environment_settings_test_01.md | Inconsistencies in environment settings test |
| 95 | 4 | dto_validation_missing.md | DTOs lacking proper validation could allow invalid data |
| 96 | 4 | dto_Type_ambiguity.md | Type ambiguity in DTO definitions could lead to conversion errors |
| 97 | 4 | models_fields_without_validation.md | Model fields lacking validation constraints |
| 98 | 4 | models-struct-tag-errors.md | Incorrect struct tags in model definitions |
| 99 | 4 | datasource_rest_query_expected_http_status_default.md | Hard-coded HTTP status in REST query data source |
| 100 | 4 | api_managed_environment_missing_api_client_validation.md | Missing client validation in managed environment API |
| 101 | 3 | dto_improper_error_handling_ai_generative.md | Suboptimal error handling in AI generative DTO |
| 102 | 3 | resource_connection_test_inadequate_status_assertions.md | Insufficient status code validation in connection tests |
| 103 | 3 | datasource_tenant_configure_method_improper_error_handling.md | Inconsistent error handling in tenant data source configuration |
| 104 | 3 | datasource_billing_policies_environments_test_issue2.md | Minor issues in billing policies environments test |
| 105 | 3 | resource_environment_application_package_install_3.md | Minor issues in application package installation resource |
| 106 | 3 | datasource_analytics_data_exports_test.md | Inconsistent test patterns in analytics exports tests |
| 107 | 3 | api_connectors_response_validation.md | Limited validation of API responses in connectors client |
| 108 | 3 | lack-of-validation-environmentId-issue.md | Missing environment ID validation in multiple components |
| 109 | 3 | resource_environment_settings_test_04.md | Minor inconsistencies in environment settings tests |
| 110 | 3 | location_dto_unexported.md | Unexported location DTO fields limiting extensibility |
| 111 | 3 | location_properties_unexported.md | Unexported location properties reducing API usability |
| 112 | 3 | locations_array_dto_unexported.md | Unexported array fields in location DTOs |
| 113 | 3 | parameter_conversion.md | Inconsistent parameter conversion in multiple components |
| 114 | 3 | warning_error_validation_missing.md | Insufficient validation in warning/error handling |
| 115 | 3 | api_environment_wave_polling_efficiency.md | Inefficient polling mechanism in environment wave API |
| 116 | 3 | dto-time-field-misuse.md | Improper time field handling in DTOs |
| 117 | 3 | models_use_string_for_capacity_enum.md | Inconsistent enum representation in models |
| 118 | 3 | datasource_tenant_settings_test_HardcodedAPIEndpoint.md | Hardcoded API endpoints in tenant settings tests |
| 119 | 3 | restore_original_value_modifier_issue_1.md | Minor issue in value restoration modifier |
| 120 | 3 | restore_original_value_modifier_issue_2.md | Secondary issue in value restoration logic |
| 121 | 3 | resource_tenant_settings_test_missing_assertions.md | Insufficient test assertions in tenant settings tests |
| 122 | 3 | datasource_currencies_read_state_handling.md | Suboptimal state handling in currencies data source |
| 123 | 3 | datasource_solutions_providerdata_type_handling.md | Type handling issue in solutions data source |
| 124 | 3 | resource_managed_environment_test_Invalid_Resource_Attribute_Settings.md | Invalid attribute settings in managed environment tests |
| 125 | 3 | resource_managed_environment_test_Lack_of_Edge_Case_Validation.md | Missing edge case testing in managed environment resource |
| 126 | 3 | datasource_tenant_capacity-missing_attribute_handling.md | Incomplete attribute handling in tenant capacity data source |
| 127 | 3 | resource_dlp_policy_test_httpmock_issue.md | HTTP mock setup issues in DLP policy tests |
| 128 | 3 | datasource_tenant_capacity_test_issue_mock_file_handling.md | Inconsistent mock file handling in tenant capacity tests |
| 129 | 3 | datasource_tenant_application_packages_test_incorrect_test_data_path.md | Incorrect test data paths in application packages tests |
| 130 | 3 | config_test_environment_persistence_bool.md | Boolean persistence issues in environment test config |
| 131 | 3 | config_test_environment_variable_persistence.md | Variable persistence problems in test configuration |
| 132 | 3 | api_environment_group_hardcoded_api_versions.md | Hardcoded API versions could break with API changes |
| 133 | 3 | resource_environment_group_test_defer_cleanup.md | Missing defer cleanup in environment group tests |
| 134 | 3 | resource_billing_policy_unnecessary_type_assertion.md | Unnecessary type assertions in billing policy resource |
| 135 | 3 | datasource_billing_policies_test_error_handling.md | Minor error handling issues in billing policies tests |
| 136 | 3 | api_admin_management_application_missing_clientId_validation.md | Missing client ID validation in admin management |
| 137 | 3 | inconsistent_field_validations_issue.md | Inconsistent pattern in field validations |
| 138 | 3 | environment_settings_resource_metadata_error_handling.md | Limited error context in environment settings resource |
| 139 | 3 | datasource_connectors_test_issue2.md | Minor inconsistencies in connectors data source tests |
| 140 | 3 | dto_lack_error_handling_solution_checker.md | Limited error handling in solution checker DTOs |
| 141 | 3 | uuid_function_NewUUIDValueMust_flaw.md | Minor flaw in UUID value function implementation |
| 142 | 3 | config_global_map_concurrency.md | Global map usage without proper concurrency controls |
| 143 | 3 | other_field_required_when_value_of_validator_improper_handling_unknown_values.md | Inadequate handling of unknown values in validators |
| 144 | 3 | invalid_http_mock_responses_in_unit_test.md | Improper HTTP mock setup in unit tests |
| 145 | 3 | dto_instance_url_missing_json_tag.md | Missing JSON tag for instance URL in DTO |
| 146 | 3 | dto_missing_json_tag_instanceurl.md | Missing JSON tag for instance URL field |
| 147 | 2 | DTO_conversion_inconsistent_error_handling.md | Minor inconsistency in DTO conversion error handling |
| 148 | 2 | enironmentDeleteDto_struct_typo_fix.md | Typo in struct name could cause confusion |
| 149 | 2 | BuildODataQueryFromModel.md | Inefficiencies in OData query building |
| 150 | 2 | ClusterDto_field_misspelling.md | Field name misspelling in cluster DTO |
| 151 | 2 | ErrorHandlingInGetPowerApps.md | Minor error handling improvement needed in Power Apps retrieval |
| 152 | 2 | dto-missing-validation.md | Non-critical validation missing in DTOs |
| 153 | 2 | environmentTemplateDto_efficiency.md | Efficiency improvements possible in environment template DTO |
| 154 | 2 | resource_data_record_test-redundant-code.md | Redundant code in data record tests |
| 155 | 2 | dto_unexported_structs.md | Non-critical unexported structs in DTOs |
| 156 | 2 | dto_unexported_structs_usability.md | Usability impact of unexported structs |
| 157 | 2 | dto_go_issue_3.md | Minor issue in DTO implementation |
| 158 | 2 | config_test_multi_string.md | Multi-string handling in test configuration |
| 159 | 2 | main_issue_2.md | Secondary issue in main implementation |
| 160 | 2 | set_bool_value_unknown_if_checksum_change_modifier_high.md | Minor issue in boolean value modifier |
| 161 | 2 | set_bool_value_unknown_if_checksum_change_modifier_index_validation_high.md | Index validation in boolean modifier |
| 162 | 2 | dto_float_precision_issues.md | Potential float precision issues in DTOs |
| 163 | 2 | Unexported_Struct_Issue.md | Non-critical unexported struct issue |
| 164 | 2 | httpmock_initializer_issues.md | Minor issues in HTTP mock initialization |
| 165 | 2 | dto_Unblockable_missing_tag.md | Missing tag in unblockable DTO |
| 166 | 2 | api_analytics_data_exports_empty_url_path_check.md | Inefficient empty URL path checking |
| 167 | 6 | unhandled_json_marshal_new.md | Unhandled JSON marshal errors in newer code paths could cause silent failures |
| 168 | 6 | issue_uncontrolled_nil_pointer_wrap_into_provider_error.md | Uncontrolled nil pointer when wrapping errors into provider errors |
| 169 | 6 | sync_attribute_plan_modifier_improper_null_assignment.md | Improper null handling in attribute plan modifier could cause data corruption |
| 170 | 5 | provider_missing_error_handling_in_validation.md | Missing error handling in provider validation could allow invalid configurations |
| 171 | 5 | provider_schema_missing_attribute_validations.md | Missing attribute validations in provider schema |
| 172 | 5 | responding_json_processing.md | Issues with response JSON processing causing potential data inconsistency |
| 173 | 5 | datasource_environment_templates_test_high_file_handling.md | Improper file handling in environment templates test |
| 174 | 5 | resource_admin_management_application_read_error_handling.md | Inadequate error handling in admin application read operation |
| 175 | 5 | resource_copilot_studio_application_insights_test_incorrect_rand.md | Non-deterministic testing due to incorrect random value usage |
| 176 | 5 | read_function_empty_list_handling.md | Improper handling of empty lists in read operations |
| 177 | 5 | lack_of_input_validation_for_environment_id.md | Insufficient validation for environment IDs across multiple components |
| 178 | 4 | helpers-go-issue-loop-flow.md | Control flow issues in helper function loops |
| 179 | 4 | datasource_languages_test_edgecases.md | Edge cases not properly tested in languages data source |
| 180 | 4 | datasource_environment_settings_test.md | Incomplete test coverage for environment settings data source |
| 181 | 4 | datasource_environment_templates.md | Implementation issues in environment templates data source |
| 182 | 4 | itemDto_field_validations.md | Missing validations for item DTO fields |
| 183 | 4 | hardcoded_identifications_test_cases.md | Hardcoded IDs in test cases limiting test flexibility |
| 184 | 4 | http_mock_error_handling.md | Improper error handling in HTTP mock implementations |
| 185 | 3 | constants_TIMEOUT.md | Timeout constants that could be better configured for variable network conditions |
| 186 | 3 | unclear_context_for_constants.md | Ambiguous contextual usage of constants |
| 187 | 3 | issue_1_resource_billing_policy_test.md | Implementation issues in billing policy test |
| 188 | 3 | test_file.md | Generic test implementation issues |
| 189 | 3 | models_unused_import.md | Unnecessary imports in models causing potential confusion |
| 190 | 3 | dto_unnecessary_import_assert.md | Redundant import assertions in DTO code |
| 191 | 3 | dto_unused_tags.md | Unused JSON tags in DTO definitions |
| 192 | 3 | uuid_value_valueuuid_error_handling.md | Insufficient error handling in UUID value conversion |
| 193 | 2 | environment_settings_resource_metadata_error_handling.md | Minor issues in environment settings resource metadata handling |

### Medium Issues

| Order Number | Score | Markdown File Name | Description |
|-------------|-------|-------------------|-------------|
| 1 | 8 | Update-Null-Dereference.md | Potential null dereference in REST resource Update function could cause runtime crash during API operations |
| 2 | 8 | potential_nil_pointer_dereference_in_read.md | Risk of nil pointer dereference in Read operation that could crash the provider |
| 3 | 8 | potential_null_pointer_dereference.md | Unchecked null pointer access in core functionality that could lead to runtime panic |
| 4 | 8 | api_environment_Infinite_Loop_AddDataverseToEnvironment.md | Possible infinite loop in AddDataverseToEnvironment function causing resource exhaustion |
| 5 | 8 | api_rest_null_pointer_scope.md | REST API null pointer risk due to improper scope management |
| 6 | 7 | config_null_pointer_dereference.md | Configuration component has potential null pointer dereference that could crash provider initialization |
| 7 | 7 | datasource_connectors_issue_undefined_listdatasourcemodel.md | Undefined model reference in connectors data source could cause unexpected behavior |
| 8 | 7 | models.go_incorrect_type_assertions.md | Unsafe type assertions without verification in models.go that could lead to panics |
| 9 | 7 | datasource_tenant_capacity-invalid_api_client.md | Invalid API client usage in tenant capacity data source could cause operations to fail |
| 10 | 7 | resource_managed_environment_issue_unexpected_provider_data.md | Incorrect type handling in managed environment resource |
| 11 | 7 | models.go_unhandled_null_conditions.md | Critical null value conditions not properly handled in model conversions |
| 12 | 7 | api_solution_issue_missing_nil_checks.md | Missing nil checks in solution API client could lead to runtime crashes |
| 13 | 7 | Create-Improper-Error-Handling.md | Inadequate error handling in resource creation that could mask critical failures |
| 14 | 7 | resource_environment_wave_nil_value_check.md | Missing nil value check in environment wave resource could cause panic |
| 15 | 7 | api_tenant_settings_nil_pointer_filterdto.md | Nil pointer risk in tenant settings filter DTO handling |
| 16 | 7 | datasource_billing_policies_environments_unverified_cast.md | Unverified type casting in billing policies environments data source |
| 17 | 7 | resource_admin_management_application_create_error.md | Improper error handling in admin application creation could lead to orphaned resources |
| 18 | 6 | api_dlp_policy_High_Issue_1.md | Improper error propagation in DLP policy API client |
| 19 | 6 | api_dlp_policy_High_Issue_2.md | Inconsistent error handling pattern in DLP policy client methods |
| 20 | 6 | api_licensing_invalid_error_handling.md | Invalid error handling approach in licensing API that could mask failures |
| 21 | 6 | api_licensing_missing_pagination.md | Missing pagination handling could lead to incomplete results in licensing operations |
| 22 | 6 | api_data_record_misuse_of_fmtsprintf_with_unescaped_characters_url_construction.md | Potential for URL injection due to unescaped characters in URL construction |
| 23 | 6 | api_currencies_improper_defer_handling.md | Resource leaks due to improper defer handling in currencies API |
| 24 | 6 | api_languages_json_unmarshalling_error.md | Missing error handling for JSON unmarshalling in languages API |
| 25 | 6 | resource_billing_policy_improper_error_handling.md | Improper error handling in billing policy resource that could lead to inconsistent state |
| 26 | 6 | Find_generic_handling.md | Generic error handling without proper context in core functionality |
| 27 | 6 | request_improper_error_handling.md | Insufficient error handling in request processing could lead to misleading error messages |
| 28 | 6 | lack_error_type_handling.md | Missing error type differentiation that could lead to incorrect error responses |
| 29 | 6 | unhandled_json_marshal.md | Unhandled JSON marshaling errors that could lead to incomplete operations |
| 30 | 6 | unhandled_json_marshall.md | Missing error handling for JSON marshalling operations |
| 31 | 6 | resource_tenant_settings_create_high.md | Improper state handling in tenant settings creation |
| 32 | 6 | resource_tenant_settings_configure_high.md | Configuration issues in tenant settings resource |
| 33 | 6 | convertToDto_error_handling_issue.md | Insufficient error handling in DTO conversion functions |
| 34 | 6 | api_locations_missing_error_wrapping.md | Missing context in error messages from locations API client |
| 35 | 6 | HardcodingURLs.md | Hardcoded URLs throughout the codebase make environment configuration difficult |
| 36 | 5 | datasource_environment_powerapps_test_missing_error_handling.md | Incomplete error handling in Power Apps test data source |
| 37 | 5 | lack_of_validation_for_url_input.md | Missing URL input validation could allow malformed requests |
| 38 | 5 | resource_enterprise_policy_create_retry.md | Inconsistent retry mechanism in enterprise policy creation |
| 39 | 5 | api_connectors_error_handling.md | Suboptimal error handling in connectors API client |
| 40 | 5 | resource_solution_issue_2.md | State inconsistency issues in solution resource |
| 41 | 5 | models_error_handling.md | Inconsistent error handling patterns in model operations |
| 42 | 5 | other_field_required_when_value_of_validator_missing_error_handling.md | Missing error handling in validator that could lead to validation bypass |
| 43 | 5 | improper-error-handling-in-enterTimeoutContext.md | Timeout context error handling issues could lead to resource leaks |
| 44 | 5 | hash_error_propagation.md | Inadequate error propagation in hash functions |
| 45 | 5 | datasource_analytics_data_exports_issue_null_check.md | Missing null check in analytics data exports could cause panic |
| 46 | 5 | datasource_data_record_high_issue_2.md | Improper data handling in data record data source |
| 47 | 5 | datasource_connectors_issue_type_assertion.md | Unsafe type assertion in connectors data source |
| 48 | 5 | api_environment_settings_issue_3.md | Incomplete validation in environment settings API |
| 49 | 5 | api_tenant_isolation_policy_Lack_of_Input_Validation_for_tenantId.md | Missing tenant ID validation in isolation policy API |
| 50 | 5 | resource_data_record_dynamic_validator.md | Insufficient validation in dynamic data record fields |
| 51 | 5 | resource_environment_group_rule_set_issue_1.md | Inconsistent error handling in environment group rule set |
| 52 | 5 | resource_environment_group_issue_redundant_error_calls.md | Redundant error handling causing confusing diagnostics |
| 53 | 5 | resource_data_record_test-missing-error-handling.md | Missing error handling in data record tests |
| 54 | 5 | resource_environment_group_rule_set_issue_3.md | State inconsistency in environment group rule sets |
| 55 | 5 | resource_copilot_studio_application_insights_issue_1.md | Error handling inconsistency in Copilot Studio insights |
| 56 | 5 | resource_billing_policy_environment_error_handling.md | Suboptimal error handling in billing policy environment resource |
| 57 | 5 | models_error_context.md | Missing error context in model operations |
| 58 | 5 | force_string_value_unknown_modifier_exit_condition_issue.md | Potential infinite loop in string modifier |
| 59 | 5 | datasource_environments_incorrect_error_handling.md | Improper error handling in environments data source |
| 60 | 5 | Overly_Permissive_UrlValidation.md | URL validation too permissive, allowing potentially malformed URLs |
| 61 | 5 | Misleading_Error_Message_DeleteEnvironmentGroupRuleSet.md | Error messages that don't accurately describe the actual issue |
| 62 | 5 | api_tenant_settings_lack_of_tests_filterdto.md | Insufficient test coverage for tenant settings filter |
| 63 | 5 | api_environment_group_improper_error_handling.md | Inadequate error handling in environment group operations |
| 64 | 5 | models_high_missing_field_validation.md | Critical fields lacking validation in model structures |
| 65 | 5 | api_enterprise_policy_retry_mechanism.md | Inefficient retry mechanism in enterprise policy API |
| 66 | 5 | issue_3_resource_billing_policy_test.md | Inconsistent test behavior in billing policy tests |
| 67 | 5 | datasource_environments_incorrect_providerdata_handling.md | Type assertion issues in environments data source |
| 68 | 5 | unhandled_errors_in_GetSolutionCheckerRules_Function_Call.md | Unchecked errors in solution checker rules client |
| 69 | 5 | resource_solution_test_improper_error_handling.md | Inadequate error assertions in solution resource tests |
| 70 | 5 | api_copilot_studio_application_insights_http_status_handling.md | Missing HTTP status code validation |
| 71 | 4 | datasource_environments_test__error_handling_mock.md | Incomplete mock response handling in environments tests |
| 72 | 4 | datasource_tenant_test_missing_error_handling.md | Missing error validation in tenant data source tests |
| 73 | 4 | Missing_Validation_of_Input_Parameters.md | Input parameters in multiple components lack validation |
| 74 | 4 | datasource_tenant_application_packages_client_initialization.md | Potential nil client in tenant application packages data source |
| 75 | 4 | datasource_tenant_application_packages_read_error_handling.md | Inconsistent error handling in application packages read operations |
| 76 | 4 | resource_billing_policy_environment_test_hardcoded_ids_issue.md | Hardcoded test IDs that could make tests brittle |
| 77 | 4 | resource_connection_share_test_missing_error_handling.md | Incomplete error handling in connection share tests |
| 78 | 4 | UserResource_ErrorHandling.md | Suboptimal error handling pattern in user resource |
| 79 | 4 | map_initialization_issue_in_mocks.md | Potential nil map access in mock implementations |
| 80 | 4 | missing_environment_id_validation_resource_user.md | Missing validation for environment ID in user resource |
| 81 | 4 | models_missing_error_handling_for_client.md | Insufficient error handling for client operations in models |
| 82 | 4 | require_replace_int_attribute_modifier.md | Potential issue with int attribute replacement modifier |
| 83 | 4 | datasource_locations_test_missing_error_handling_for_http_responses.md | Missing HTTP response validation in locations tests |
| 84 | 4 | auth_error_wrapping_getAssertion.md | Insufficient error context in authentication assertions |
| 85 | 4 | datasource_environment_application_packages_missing_environment_id_validation.md | Missing validation for environment ID in application packages |
| 86 | 4 | missing_validation_dto_fields.md | Insufficient validation in DTO field definitions |
| 87 | 4 | requires_replace_string_from_non_empty_modifier_Limited_error_handling_in_PlanModifyString.md | Limited error handling in string replacement modifier |
| 88 | 4 | models_function_error_handling.md | Inconsistent error handling in model conversion functions |
| 89 | 4 | resource_tenant_isolation_policy_test_irregular_error_handling.md | Irregular error handling in isolation policy tests |
| 90 | 4 | resource_rest_test_01_hardcoded_metadata_url.md | Hardcoded URLs in REST resource tests |
| 91 | 4 | Missing_JSON_Tags_in_linkedEnvironmentIdMetadataDto_Struct.md | Missing JSON tags in DTO struct could cause serialization issues |
| 92 | 4 | resource_tenant_isolation_policy_modularization.md | Insufficient modularization in isolation policy resource |
| 93 | 4 | resource_billing_policy_environment_test_randIntn_issue.md | Non-deterministic test behavior due to random values |
| 94 | 4 | resource_environment_settings_test_01.md | Inconsistencies in environment settings test |
| 95 | 4 | dto_validation_missing.md | DTOs lacking proper validation could allow invalid data |
| 96 | 4 | dto_Type_ambiguity.md | Type ambiguity in DTO definitions could lead to conversion errors |
| 97 | 4 | models_fields_without_validation.md | Model fields lacking validation constraints |
| 98 | 4 | models-struct-tag-errors.md | Incorrect struct tags in model definitions |
| 99 | 4 | datasource_rest_query_expected_http_status_default.md | Hard-coded HTTP status in REST query data source |
| 100 | 4 | api_managed_environment_missing_api_client_validation.md | Missing client validation in managed environment API |
| 101 | 3 | dto_improper_error_handling_ai_generative.md | Suboptimal error handling in AI generative DTO |
| 102 | 3 | resource_connection_test_inadequate_status_assertions.md | Insufficient status code validation in connection tests |
| 103 | 3 | datasource_tenant_configure_method_improper_error_handling.md | Inconsistent error handling in tenant data source configuration |
| 104 | 3 | datasource_billing_policies_environments_test_issue2.md | Minor issues in billing policies environments test |
| 105 | 3 | resource_environment_application_package_install_3.md | Minor issues in application package installation resource |
| 106 | 3 | datasource_analytics_data_exports_test.md | Inconsistent test patterns in analytics exports tests |
| 107 | 3 | api_connectors_response_validation.md | Limited validation of API responses in connectors client |
| 108 | 3 | lack-of-validation-environmentId-issue.md | Missing environment ID validation in multiple components |
| 109 | 3 | resource_environment_settings_test_04.md | Minor inconsistencies in environment settings tests |
| 110 | 3 | location_dto_unexported.md | Unexported location DTO fields limiting extensibility |
| 111 | 3 | location_properties_unexported.md | Unexported location properties reducing API usability |
| 112 | 3 | locations_array_dto_unexported.md | Unexported array fields in location DTOs |
| 113 | 3 | parameter_conversion.md | Inconsistent parameter conversion in multiple components |
| 114 | 3 | warning_error_validation_missing.md | Insufficient validation in warning/error handling |
| 115 | 3 | api_environment_wave_polling_efficiency.md | Inefficient polling mechanism in environment wave API |
| 116 | 3 | dto-time-field-misuse.md | Improper time field handling in DTOs |
| 117 | 3 | models_use_string_for_capacity_enum.md | Inconsistent enum representation in models |
| 118 | 3 | datasource_tenant_settings_test_HardcodedAPIEndpoint.md | Hardcoded API endpoints in tenant settings tests |
| 119 | 3 | restore_original_value_modifier_issue_1.md | Minor issue in value restoration modifier |
| 120 | 3 | restore_original_value_modifier_issue_2.md | Secondary issue in value restoration logic |
| 121 | 3 | resource_tenant_settings_test_missing_assertions.md | Insufficient test assertions in tenant settings tests |
| 122 | 3 | datasource_currencies_read_state_handling.md | Suboptimal state handling in currencies data source |
| 123 | 3 | datasource_solutions_providerdata_type_handling.md | Type handling issue in solutions data source |
| 124 | 3 | resource_managed_environment_test_Invalid_Resource_Attribute_Settings.md | Invalid attribute settings in managed environment tests |
| 125 | 3 | resource_managed_environment_test_Lack_of_Edge_Case_Validation.md | Missing edge case testing in managed environment resource |
| 126 | 3 | datasource_tenant_capacity-missing_attribute_handling.md | Incomplete attribute handling in tenant capacity data source |
| 127 | 3 | resource_dlp_policy_test_httpmock_issue.md | HTTP mock setup issues in DLP policy tests |
| 128 | 3 | datasource_tenant_capacity_test_issue_mock_file_handling.md | Inconsistent mock file handling in tenant capacity tests |
| 129 | 3 | datasource_tenant_application_packages_test_incorrect_test_data_path.md | Incorrect test data paths in application packages tests |
| 130 | 3 | config_test_environment_persistence_bool.md | Boolean persistence issues in environment test config |
| 131 | 3 | config_test_environment_variable_persistence.md | Variable persistence problems in test configuration |
| 132 | 3 | api_environment_group_hardcoded_api_versions.md | Hardcoded API versions could break with API changes |
| 133 | 3 | resource_environment_group_test_defer_cleanup.md | Missing defer cleanup in environment group tests |
| 134 | 3 | resource_billing_policy_unnecessary_type_assertion.md | Unnecessary type assertions in billing policy resource |
| 135 | 3 | datasource_billing_policies_test_error_handling.md | Minor error handling issues in billing policies tests |
| 136 | 3 | api_admin_management_application_missing_clientId_validation.md | Missing client ID validation in admin management |
| 137 | 3 | inconsistent_field_validations_issue.md | Inconsistent pattern in field validations |
| 138 | 3 | environment_settings_resource_metadata_error_handling.md | Limited error context in environment settings resource |
| 139 | 3 | datasource_connectors_test_issue2.md | Minor inconsistencies in connectors data source tests |
| 140 | 3 | dto_lack_error_handling_solution_checker.md | Limited error handling in solution checker DTOs |
| 141 | 3 | uuid_function_NewUUIDValueMust_flaw.md | Minor flaw in UUID value function implementation |
| 142 | 3 | config_global_map_concurrency.md | Global map usage without proper concurrency controls |
| 143 | 3 | other_field_required_when_value_of_validator_improper_handling_unknown_values.md | Inadequate handling of unknown values in validators |
| 144 | 3 | invalid_http_mock_responses_in_unit_test.md | Improper HTTP mock setup in unit tests |
| 145 | 3 | dto_instance_url_missing_json_tag.md | Missing JSON tag for instance URL in DTO |
| 146 | 3 | dto_missing_json_tag_instanceurl.md | Missing JSON tag for instance URL field |
| 147 | 2 | DTO_conversion_inconsistent_error_handling.md | Minor inconsistency in DTO conversion error handling |
| 148 | 2 | enironmentDeleteDto_struct_typo_fix.md | Typo in struct name could cause confusion |
| 149 | 2 | BuildODataQueryFromModel.md | Inefficiencies in OData query building |
| 150 | 2 | ClusterDto_field_misspelling.md | Field name misspelling in cluster DTO |
| 151 | 2 | ErrorHandlingInGetPowerApps.md | Minor error handling improvement needed in Power Apps retrieval |
| 152 | 2 | dto-missing-validation.md | Non-critical validation missing in DTOs |
| 153 | 2 | environmentTemplateDto_efficiency.md | Efficiency improvements possible in environment template DTO |
| 154 | 2 | resource_data_record_test-redundant-code.md | Redundant code in data record tests |
| 155 | 2 | dto_unexported_structs.md | Non-critical unexported structs in DTOs |
| 156 | 2 | dto_unexported_structs_usability.md | Usability impact of unexported structs |
| 157 | 2 | dto_go_issue_3.md | Minor issue in DTO implementation |
| 158 | 2 | config_test_multi_string.md | Multi-string handling in test configuration |
| 159 | 2 | main_issue_2.md | Secondary issue in main implementation |
| 160 | 2 | set_bool_value_unknown_if_checksum_change_modifier_high.md | Minor issue in boolean value modifier |
| 161 | 2 | set_bool_value_unknown_if_checksum_change_modifier_index_validation_high.md | Index validation in boolean modifier |
| 162 | 2 | dto_float_precision_issues.md | Potential float precision issues in DTOs |
| 163 | 2 | Unexported_Struct_Issue.md | Non-critical unexported struct issue |
| 164 | 2 | httpmock_initializer_issues.md | Minor issues in HTTP mock initialization |
| 165 | 2 | dto_Unblockable_missing_tag.md | Missing tag in unblockable DTO |
| 166 | 2 | api_analytics_data_exports_empty_url_path_check.md | Inefficient empty URL path checking |
| 167 | 3 | DocumentationMissingFromStructMembers.md | Missing documentation for struct members reduces code maintainability |
| 168 | 3 | ImproperUseGlobalFunctions_MockedTesting.md | Use of global functions in mocked tests makes tests brittle |
| 169 | 3 | Diff_map_initialization.md | Inefficient map initialization in diff functions |
| 170 | 3 | Diff_slice_initialization.md | Inefficient slice initialization in diff operations |
| 171 | 3 | UnnecessaryField_DisableDelete.md | Unnecessary disable delete field increases schema complexity |
| 172 | 3 | UnusedImportsTightCoupling.md | Unused imports causing tight coupling between components |
| 173 | 3 | Unused_Variables_and_Functions.md | Dead code in the form of unused variables and functions |
| 174 | 3 | uuid_type_valuefromterraform_error_handling.md | Suboptimal error handling in UUID terraform value conversion |
| 175 | 3 | Incorrect_Error_Handling.md | Error handling inconsistencies across multiple components |
| 176 | 3 | ambiguous_defaults_in_convertFromDto_issue.md | Ambiguous default values in DTO to model conversion |
| 177 | 3 | api_admin_management_application_hardcoded_api_version.md | Hardcoded API versions in admin management client |
| 178 | 3 | api_copilot_studio_application_insights_error_details.md | Incomplete error details in Copilot Studio insights |
| 179 | 3 | api_copilot_studio_application_insights_parse_import_id.md | Import ID parsing issues in Copilot Studio component |
| 180 | 3 | api_dlp_policy_Medium_Issue_1.md | Medium severity issues in DLP policy API implementation |
| 181 | 3 | api_dlp_policy_Medium_Issue_2.md | Secondary issues in DLP policy API error handling |
| 182 | 3 | cert_handle_file_read_errors.md | Inadequate handling of file read errors in certificate operations |
| 183 | 3 | cert_hardcoded_error_messages_limiting_localization.md | Hardcoded error messages limit localization options |
| 184 | 3 | helpers-go-issue-error-custom-patterns.md | Non-standard error handling patterns in helpers |
| 185 | 3 | hardcoded_network_isolation.md | Hardcoded network isolation settings reducing flexibility |
| 186 | 3 | improper_configvalue_check.md | Improper configuration value validation |
| 187 | 3 | known_state_value_test_function_error_logging.md | Suboptimal error logging in state value test functions |
| 188 | 3 | auth_logging_sanitization_AuthenticateOIDC.md | Incomplete logging sanitization in OIDC authentication |
| 189 | 3 | make_field_required_when_other_field_does_not_have_value_missing_docs.md | Missing documentation for conditional field requirements |
| 190 | 3 | model_fields_without_validation.md | Multiple model fields lacking validation logic |
| 191 | 3 | issue_improper_use_provider_error.md | Non-standard usage of provider error types |
| 192 | 3 | unit_test_error_handling.md | Inconsistent error handling approach in unit tests |
| 193 | 3 | unnecessary_composed_test_check_functions.md | Excessive test check function composition |
| 194 | 3 | typeinfo.md | Issues with type information implementation and usage |
| 195 | 3 | resource_solution_issue_3.md | Solution resource implementation concerns |
| 196 | 3 | resource_connection.md | Connection resource implementation issues |
| 197 | 3 | resource_data_record_test-unused-imports.md | Unused imports in data record test files |
| 198 | 3 | resource_enterprise_policy_test_mock_overuse.md | Excessive use of mocks in enterprise policy tests |
| 199 | 3 | resource_environment_application_package_install_4.md | Implementation issues in application package installation |
| 200 | 3 | struct_field_validation_missing.md | Missing validations for multiple struct fields |
| 201 | 2 | Delete-Diagnostic-Error-Handling.md | Minor issues with diagnostic error handling in delete operations |
| 202 | 2 | MissingUnitTestsForClientInitialization.md | Lack of unit tests for client initialization functions |
| 203 | 2 | constants.go-inappropriate-naming.md | Non-standard naming conventions in constants |
| 204 | 2 | dto_go_issue_1.md | Minor issues in DTO implementation |
| 205 | 2 | release.md | Minor issues in release management code |
| 206 | 2 | main_issue_1.md | Minor implementation issues in main package |
| 207 | 2 | provider_logging_verbosity.md | Suboptimal logging verbosity settings in provider |
| 208 | 2 | uuid_function_NewUUIDPointerValueMust_flaw.md | Minor flaw in UUID pointer value function |
| 209 | 2 | models_missing_docs.md | Missing documentation for model functions |
| 210 | 2 | MissingHTTPStatusValidation.md | HTTP status validation missing in some API calls |
| 211 | 2 | HardcodedAPIVersionUsage.md | API version hardcoding in various components |
| 212 | 2 | constants_ANALYTICS_SCOPE.md | Potential issues with analytics scope constant definition |
| 213 | 2 | constants_naming_convention_issue.md | Naming convention inconsistencies in constants |
| 214 | 2 | function_naming_issue_in_mocks.md | Non-standard function naming in mock implementations |
| 215 | 2 | hash_file_close_error.md | Minor file handle closing issues in hash functions |
| 216 | 2 | models-go-empty-environments-handling.md | Minor issues handling empty environments in models |
| 217 | 2 | models-missing-validation.md | Non-critical validations missing in some models |
| 218 | 2 | models.md | Minor documentation and implementation issues in models |
| 219 | 2 | itemDisabledReasonDto_field_constraints.md | Missing constraints for disabled reason fields |
| 220 | 2 | skipping_dto_conversion_with_empty_ids_issue.md | Minor issue with empty ID handling in DTO conversion |
| 221 | 1 | Magic_Strings_in_API_Version_Usage.md | Low-impact use of magic strings for API versions |
| 222 | 1 | Hardcoding_of_API_Version_Query_Parameter.md | Low-severity API version parameter hardcoding |
| 223 | 1 | Inconsistent_JSON_Tags.md | Low-impact inconsistencies in JSON tag formatting |

### Low Issues

| Order Number | Score | Markdown File Name | Description |
|-------------|-------|-------------------|-------------|
| 1 | 3 | unexpected_http_status_code_error_issue_error_interface_validation.md | HTTP status code error handling could be improved for better error reporting |
| 2 | 3 | unexpected_provider_data_type_error.md | Type handling for provider data could be more robust |
| 3 | 3 | require_replace_int_attribute_modifier.md | Integer attribute replacement modifier has minor implementation issues |
| 4 | 3 | missing_diagnostic_handling_for_configuration_validation.md | Configuration validation could provide more detailed diagnostic information |
| 5 | 3 | api_analytics_data_exports_error_handling_context.md | Error handling in analytics exports API could include more context |
| 6 | 3 | api_solution_issue_inconsistent_error_handling.md | Inconsistent error handling patterns in solution API client |
| 7 | 3 | dynamic_columns_issue_1.md | Minor implementation issue in dynamic columns handling |
| 8 | 3 | Metadata-Context-Management.md | Context management in metadata operations could be improved |
| 9 | 2 | api_dlp_policy_Low_Issue_1.md | Minor issue in DLP policy API implementation |
| 10 | 2 | api_dlp_policy_Low_Issue_2.md | Secondary minor issue in DLP policy API error handling |
| 11 | 2 | resource_copilot_studio_application_insights_issue_3.md | Minor issue in Copilot Studio application insights resource |
| 12 | 2 | resource_copilot_studio_application_insights_test_missing_error_handling.md | Non-critical error handling improvements possible in tests |
| 13 | 2 | environment_settings_resource_import_state_function.md | Minor refinements possible in import state function |
| 14 | 2 | resource_environment_application_package_install_5.md | Minor issue in application package installation resource |
| 15 | 2 | api_tenant_isolation_policy_Hardcoded_Default_Retry_Duration_for_Retry_After_Header_Parsing.md | Hardcoded retry duration value could be configurable |
| 16 | 2 | api_environment_settings_issue_5.md | Minor issue in environment settings API client |
| 17 | 2 | api_environment_wave_missing_documentation.md | Documentation could be improved for environment wave API |
| 18 | 2 | api_environment_group_api_call_inefficiency.md | Minor inefficiency in environment group API calls |
| 19 | 2 | raw_string_concatenation_issue.md | String concatenation could be more efficient |
| 20 | 2 | datasource_tenant_capacity-repeated_initialization_pattern.md | Repeated code pattern in tenant capacity data source |
| 21 | 2 | auth_hardcoded_bearer_getAssertion.md | Bearer token prefix is hardcoded |
| 22 | 2 | api_rest_variable_initialization.md | Inefficient variable initialization in REST API client |
| 23 | 2 | datasource_securityroles_typo.md | Typographical error in security roles data source |
| 24 | 2 | datasource_solutions_context_redundant.md | Redundant context handling in solutions data source |
| 25 | 2 | datasource_tenant_application_packages_test_unnecessary_httpmock_activation.md | Unnecessary mock activation in test |
| 26 | 2 | models_issue_1.md | Minor issue in model implementation |
| 27 | 2 | models-go-unoptimized-slice-allocation.md | Slice allocation could be optimized |
| 28 | 2 | models_redundant_checks.md | Redundant validation checks in models |
| 29 | 2 | resource_tenant_settings_create_low.md | Minor issue in tenant settings resource creation |
| 30 | 2 | resource_billing_policy_environment_redundant_api_calls.md | Redundant API calls in billing policy environment resource |
| 31 | 2 | resource_environment_application_package_install_test.md | Test implementation could be improved |
| 32 | 2 | resource_environment_group_rule_set_issue_4.md | Minor issue in environment group rule set resource |
| 33 | 2 | resource_environment_group_test_issue_2.md | Test implementation issue in environment group tests |
| 34 | 2 | resource_solution_issue_4.md | Minor issue in solution resource |
| 35 | 2 | Unnecessarily_High_Integer_Type_for_statusCode_Field.md | Using int64 for HTTP status code when int would suffice |
| 36 | 2 | UnusedStructDefinitions.md | Unused struct definitions increase code size |
| 37 | 2 | Unused_Constants.md | Unused constants in the codebase |
| 38 | 2 | datasource_rest_query_markdown_description_issue.md | Markdown description improvements needed |
| 39 | 2 | resource_data_record_test-variable-naming-inconsistent.md | Inconsistent variable naming in tests |
| 40 | 2 | inefficient_loop_operation_in_GetListStringValues.md | Loop operation could be more efficient |
| 41 | 2 | models-inconsistent-field-naming.md | Inconsistent naming patterns for model fields |
| 42 | 2 | unnecessary_code_duplication_prepare_request_body.md | Duplicated code in request body preparation |
| 43 | 2 | url_format_error_variable_naming.md | Variable naming issue in URL formatting code |
| 44 | 2 | models_type_conversion.md | Type conversion could be implemented more clearly |
| 45 | 2 | require_replace_object_to_empty_modifier_misleading_description.md | Misleading description in object replacement modifier |
| 46 | 2 | redundancy_in_fallback_mechanism_for_GetConfigString.md | Redundant fallback handling in configuration |
| 47 | 2 | restore_original_value_modifier_issue_3.md | Minor issue in value restoration modifier |
| 48 | 2 | api_data_record_potential_memory_allocation_inefficiency_applydatarecord_function.md | Memory allocation could be more efficient |
| 49 | 2 | resource_data_record_metadata_logging.md | Metadata logging could be improved |
| 50 | 2 | buildODataSavedQuery.md | OData query building could be more efficient |
| 51 | 2 | constants_Organization.md | Organization constants could be better organized |
| 52 | 1 | datasource_tenant_settings_metadata_logging.md | Low impact metadata logging issue |
| 53 | 1 | datasource_tenant_settings_schema_url.md | URL in schema could be improved |
| 54 | 1 | datasource_environments_test__unused_imports.md | Unused imports in test file |
| 55 | 1 | datasource_environments_unused_policyattributeschema.md | Unused policy attribute schema |
| 56 | 1 | datasource_languages_test_documentation.md | Documentation improvements in test file |
| 57 | 1 | datasource_locations_test_missing_comments_on_test_cases.md | Missing comments in test cases |
| 58 | 1 | datasource_solutions_test_issue3.md | Minor test implementation issue |
| 59 | 1 | datasource_currencies_metadata_logging.md | Metadata logging improvements possible |
| 60 | 1 | datasource_currencies_test_Issue2.md | Minor test implementation issue |
| 61 | 1 | datasource_dlp_policy_test_duplicate_attr_check.md | Duplicate attribute check in test |
| 62 | 1 | datasource_environment_powerapps_test_hardcoded_paths.md | Hardcoded paths in test could be variables |
| 63 | 1 | datasource_environment_settings.md | Minor implementation issue |
| 64 | 1 | datasource_environment_settings_test.md | Minor test implementation issue |
| 65 | 1 | datasource_analytics_data_exports_test_mock_responder_issues.md | Mock responder could be improved |
| 66 | 1 | datasource_billing_policies_test_mismatch_rand.md | Mismatched random value usage in test |
| 67 | 1 | datasource_connections_test.go_hardcoded_mock_url.md | Hardcoded URL in test mock |
| 68 | 1 | variable_naming_issue_datasource_environment_settings.md | Variable naming could be improved |
| 69 | 1 | unused_import.md | Unused import in file |
| 70 | 1 | unused_mock_cleanup_command.md | Unused mock cleanup command |
| 71 | 1 | api_analytics_data_exports_redundant_uppercase_conversion.md | Redundant conversion to uppercase |
| 72 | 1 | api_capacity_unnecessary_url_allocation.md | Unnecessary URL allocation |
| 73 | 1 | api_enterprise_policy_typo_variable_name.md | Typo in variable name |
| 74 | 1 | api_licensing_incorrect_method_name.md | Method name could be more descriptive |
| 75 | 1 | api_managed_environment_duplicated_url_construction.md | Duplicated URL construction logic |
| 76 | 1 | global_variable_issue_in_mocks.md | Global variable usage in mocks |
| 77 | 1 | client_test_Inconsistent_Use_of_cancel_HappyPath.md | Inconsistent cancel usage in tests |
| 78 | 1 | dto.md | Minor DTO implementation issues |
| 79 | 1 | dto_go_issue_2.md | Minor issue in DTO file |
| 80 | 1 | dto_go_issue_4.md | Additional minor issue in DTO file |
| 81 | 1 | dto_go_issue_5.md | Further minor issue in DTO file |
| 82 | 1 | dto_improper_comment_clarity.md | Comment clarity could be improved |
| 83 | 1 | dto_lack_of_documentation.md | Documentation could be added or improved |
| 84 | 1 | dto_low_documentation_comments.md | Low-quality documentation comments |
| 85 | 1 | PowerPagesSettingsDto_unused_type.md | Unused DTO type in codebase |
| 86 | 1 | permissionPropertiesDto-incorrect-json-tag.md | Incorrect JSON tag in DTO |
| 87 | 1 | shareConnectionRequestPutPropertiesDto-incorrect-json-tag.md | Incorrect JSON tag in connection DTO |
| 88 | 1 | resource_connection_share.go_markdown_description_validation.md | Markdown description needs validation |
| 89 | 1 | resource_data_record_test-redundant-comments.md | Redundant comments in test file |
| 90 | 1 | resource_dlp_policy_test_issue_unused_variables.md | Unused variables in test |
| 91 | 1 | resource_enterprise_policy_test_inconsistent_test_name.md | Inconsistent test naming pattern |
| 92 | 1 | resource_enterprise_policy_test_unused_dependency.md | Unused dependency in test |
| 93 | 1 | resource_enterprise_policy_update.md | Minor update function issue |
| 94 | 1 | resource_tenant_isolation_policy_test_comment_improvement.md | Test comments could be improved |
| 95 | 1 | spelling_error_log_message.md | Spelling error in log message |
| 96 | 1 | test_permission_issue.md | Minor permission issue in test |
| 97 | 1 | issue_5_resource_billing_policy_test.md | Minor test implementation issue |
| 98 | 1 | logging_in_metadata_function.md | Logging could be improved in metadata function |
| 99 | 1 | make_field_required_when_other_field_does_not_have_value_unused_import.md | Unused import in validator file |
| 100 | 1 | md_solution_checker_rules.md | Documentation issue |
| 101 | 1 | models-unused-imports.md | Unused imports in models file |
| 102 | 1 | models.go_import_grouping.md | Import grouping could be improved |
| 103 | 1 | models.go_nonstandard_naming.md | Non-standard naming pattern used |
| 104 | 1 | models_documentation_missing.md | Missing documentation for models |
| 105 | 1 | models_generic_typeinfo_issue.md | Generic type info could be improved |
| 106 | 1 | models_go_basic_validation.md | Basic validation could be enhanced |
| 107 | 1 | models_go_import_not_used.md | Unused import in models file |
| 108 | 1 | models_lack_struct_comments.md | Missing comments for struct definitions |
| 109 | 1 | models_low_unnecessary_json_tags.md | Unnecessary JSON tags in models |
| 110 | 1 | models_missing_documentation.md | Documentation missing for models |
| 111 | 1 | models_unused_import.md | Unused import in models file |
| 112 | 1 | models_unused_imports.md | Multiple unused imports |
| 113 | 1 | uuid_function_NewUUIDNull_flaw.md | Minor implementation flaw in UUID function |
| 114 | 1 | uuid_function_NewUUIDPointerValue_flaw.md | Implementation flaw in UUID pointer value function |
| 115 | 1 | uuid_function_NewUUIDUnknown_flaw.md | Implementation flaw in UUID unknown function |
| 116 | 1 | uuid_function_NewUUIDValue_flaw.md | Implementation flaw in UUID value function |
| 117 | 1 | uuid_value_validateattribute_error.md | Attribute validation error in UUID value handling |
| 118 | 1 | CommentedCopyright.md | Copyright notice formatted as comment |
| 119 | 1 | helpers-go-issue-function-name.md | Function naming issue in helpers |
| 120 | 1 | no_issues_found.md | File indicating no issues found in a component |
| 121 | 2 | Hardcoded_Regex_No_Comments.md | Regular expression patterns are hardcoded without explanatory comments |
| 122 | 2 | RedundantTestCaseConfigurations.md | Test configurations contain duplicate or redundant setup code |
| 123 | 2 | Redundant_HTTP_Status_Handling_GetEnvironmentGroupRuleSet.md | Repetitive HTTP status handling in environment group rule set code |
| 124 | 2 | Redundant_StringRegex_Definition.md | Multiple similar string regex definitions that could be consolidated |
| 125 | 2 | readability_config_issues.md | Configuration code could be more readable and better structured |
| 126 | 2 | request_unnecessary_session_initialization.md | Session is initialized unnecessarily in request handling |
| 127 | 2 | requires_replace_string_from_non_empty_modifier_Unclear_struct_documentation.md | Unclear documentation in string replacement modifier struct |
| 128 | 2 | sync_attribute_plan_modifier_redundant_null_checks.md | Redundant null checking in attribute plan modifier |
| 129 | 1 | provider_unnecessary_verbose_comments.md | Overly verbose comments in provider code |
| 130 | 1 | resource_environment_application_package_install_test.md | Minor improvements needed in application package install tests |
| 131 | 1 | resource_tenant_isolation_policy_test_comment_improvement.md | Comments in tenant isolation policy tests could be clearer |
| 132 | 1 | resource_tenant_settings_create_low.md | Minor issue in tenant settings creation code |
