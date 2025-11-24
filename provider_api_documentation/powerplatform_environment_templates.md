## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud) |
| ------------------- | ----------- | -------------------------------- |
| Read                | `GET`       | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/{location}/environmentTemplates?api-version=2021-07-01` |

## Attribute Mapping

| Data Source Attribute                             | API Response JSON Field                    |
| ------------------------------------------------- | ------------------------------------------ |
| `environment_templates`                           | `value`                                    |
| `environment_templates.category`                  | derived from top-level array grouping      |
| `environment_templates.id`                        | `id`                                       |
| `environment_templates.name`                      | `name`                                     |
| `environment_templates.display_name`              | `properties.displayName`                   |
| `environment_templates.location`                  | `location`                                 |
| `environment_templates.is_disabled`               | `properties.isDisabled`                    |
| `environment_templates.disabled_reason_code`      | `properties.disabledReason.code`           |
| `environment_templates.disabled_reason_message`   | `properties.disabledReason.message`        |
| `environment_templates.is_customer_engagement`    | `properties.isCustomerEngagement`          |
| `environment_templates.is_supported_for_reset_operation` | `properties.isSupportedForResetOperation` |

### Example API Response

An example of the API response used by this data source can be found in the test fixture [`environment_templates/tests/datasource/Validate_Read/get_environment_templates.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/environment_templates/tests/datasource/Validate_Read/get_environment_templates.json).
