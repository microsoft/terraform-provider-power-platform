# powerplatform_languages (Data Source)

Fetches the list of Dynamics 365 languages. For more information see [Power Platform Enable Languages](https://learn.microsoft.com/power-platform/admin/enable-languages)

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud) |
| ------------------- | ----------- | -------------------------------- |
| Read                | `GET`       | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/{location}/environmentLanguages?api-version=2023-06-01` |

## Attribute Mapping

| Data Source Attribute      | API Response JSON Field |
| -------------------------- | ----------------------- |
| `languages`                | `value`                 |
| `languages.id`             | `id`                    |
| `languages.name`           | `name`                  |
| `languages.display_name`   | `properties.displayName` |
| `languages.localized_name` | `properties.localizedName` |
| `languages.locale_id`      | `properties.localeId`   |
| `languages.is_tenant_default` | `properties.isTenantDefault` |

### Example API Response

An example of the API response used by this data source can be found in the test fixture [`languages/tests/datasource/Validate_Read/get_languages.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/languages/tests/datasource/Validate_Read/get_languages.json).
