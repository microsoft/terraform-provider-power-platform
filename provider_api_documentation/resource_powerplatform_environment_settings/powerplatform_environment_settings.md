## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud) |
| ------------------- | ----------- | -------------------------------- |
| Read                | `GET`       | `https://{environment_host}/api/data/v9.2/retrievesettinglist` |
| Update              | `POST`      | `https://{environment_host}/api/data/v9.2/RetrieveCurrentOrganization2` |

## Attribute Mapping

| Resource Attribute                                        | API Response JSON Field |
| --------------------------------------------------------- | ----------------------- |
| `environment_id`                                          | from request / Terraform config |
| `audit_and_logs.plugin_trace_log_setting`                 | `plugintracelogsetting` |
| `audit_and_logs.audit_settings.is_audit_enabled`          | `isauditenabled` |
| `audit_and_logs.audit_settings.is_user_access_audit_enabled` | `isuseraccessauditenabled` |
| `audit_and_logs.audit_settings.is_read_audit_enabled`     | `isreadauditenabled` |
| `audit_and_logs.audit_settings.log_retention_period_in_days` | `auditretentionperiodv2` |
| `email.email_settings.max_upload_file_size_in_bytes`      | `maxuploadfilesize` |
| `product.behavior_settings.show_dashboard_cards_in_expanded_state` | `bounddashboarddefaultcardexpanded` |
| `product.features.power_apps_component_framework_for_canvas_apps` | `iscustomcontrolsincanvasappsenabled` |
| `product.features.enable_powerapps_maker_bot`             | `powerappsmakerbotenabled` |
| `product.features.enable_access_to_session_transcripts_for_copilot_studio` | `blockaccesstosessiontranscriptsforcopilotstudio` (inverted) |
| `product.features.enable_transcript_recording_for_copilot_studio` | `blocktranscriptrecordingforcopilotstudio` (inverted) |
| `product.features.enable_ai_prompts`                      | `aipromptsenabled` |
| `product.features.enable_copilot_studio_share_data_with_viva_insights` | `enablecopilotstudiosharedatawithvivainsights` |
| `product.features.enable_copilot_studio_cross_geo_share_data_with_viva_insights` | `enablecopilotstudiocrossgeosharedatawithvivainsights` |
| `product.features.enable_preview_and_experimental_ai_models` | `paipreviewscenarioenabled` |
| `product.features.enable_copilot_answer_control`          | `SettingDetailCollection[EnableCopilotAnswerControl]` |
| `product.features.enable_ai_powered_chat`                 | `SettingDetailCollection[appcopilotenabled]` |
| `product.features.ai_form_fill_automatic_suggestions`     | `SettingDetailCollection[FormPredictEnabled]` |
| `product.features.ai_form_fill_smart_paste_and_file_suggestions` | `SettingDetailCollection[FormPredictSmartPasteEnabledOnByDefault]` |
| `product.features.ai_form_fill_toolbar`                   | `SettingDetailCollection[FormFillBarUXEnabled]` |
| `product.features.natural_language_grid_and_view_search`  | `SettingDetailCollection[NLGridSearchSetting]` |
| `product.features.allow_ai_to_generate_charts`            | `SettingDetailCollection[NLChartDataVisualizationSetting]` |
| `product.security.allow_application_user_access`          | `allowapplicationuseraccess` |
| `product.security.allow_microsoft_trusted_service_tags`   | `allowmicrosofttrustedservicetags` |
| `product.security.allowed_ip_range_for_firewall`          | `allowediprangeforfirewall` |
| `product.security.allowed_service_tags_for_firewall`      | `allowedservicetagsforfirewall` |
| `product.security.enable_ip_based_cookie_binding`         | `enableipbasedcookiebinding` |
| `product.security.enable_ip_based_firewall_rule`          | `enableipbasedfirewallrule` |
| `product.security.enable_ip_based_firewall_rule_in_audit_mode` | `enableipbasedfirewallruleinauditmode` |
| `product.security.reverse_proxy_ip_addresses`             | `reverseproxyipaddresses` |

### Example API Response

Examples of API responses used by this resource can be found in the test fixtures:

- [`environment_settings/tests/resources/Validate_Read/get_retrievesettinglist.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/environment_settings/tests/resources/Validate_Read/get_retrievesettinglist.json)
- [`environment_settings/tests/resources/Validate_Read/get_organisations_1.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/environment_settings/tests/resources/Validate_Read/get_organisations_1.json)
- [`environment_settings/tests/resources/Validate_Read/get_organisations_2.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/environment_settings/tests/resources/Validate_Read/get_organisations_2.json)
- [`environment_settings/tests/resources/Validate_Read/get_organisations_3.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/environment_settings/tests/resources/Validate_Read/get_organisations_3.json)
- Update sequence examples:
  - [`environment_settings/tests/resources/Validate_Update/get_retrievesettinglist_step1.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/environment_settings/tests/resources/Validate_Update/get_retrievesettinglist_step1.json)
  - [`environment_settings/tests/resources/Validate_Update/get_retrievesettinglist_step2.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/environment_settings/tests/resources/Validate_Update/get_retrievesettinglist_step2.json)
