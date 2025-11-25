# `powerplatform_environment`

This resource is used to manage a Power Platform environment.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud) |
| ------------------- | ----------- | -------------------------------- |
| Create              | `POST`      | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01` |
| Read                | `GET`       | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/{environment_id}?api-version=2023-06-01` |
| Update              | `PATCH`     | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/{environment_id}?api-version=2023-06-01` |
| Delete              | `DELETE`    | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/{environment_id}?api-version=2023-06-01` |

## Attribute Mapping

| Resource Attribute | API Response JSON Field |
| ------------------ | ----------------------- |
| `display_name`     | `properties.displayName` |
| `location`         | `location`              |
| `environment_type` | `properties.environmentType` |
| `azure_region`     | `properties.azureRegion` |
| `templates`        | `properties.templates`  |
| `dataverse`        | `properties.linkedAppMetadata.instance` |
| `dataverse.url`    | `properties.linkedAppMetadata.instance.url` |
| `dataverse.domain` | `properties.linkedAppMetadata.instance.domain` |
| `dataverse.security_group_id` | `properties.linkedAppMetadata.instance.securityGroupId` |
| `dataverse.version` | `properties.linkedAppMetadata.instance.version` |
| `dataverse.language_code` | `properties.linkedAppMetadata.instance.languageCode` |
| `dataverse.currency_code` | `properties.linkedAppMetadata.instance.currencyCode` |
| `dataverse.currency_name` | `properties.linkedAppMetadata.instance.currencyName` |
| `dataverse.currency_symbol` | `properties.linkedAppMetadata.instance.currencySymbol` |
| `dataverse.base_language` | `properties.linkedAppMetadata.instance.baseLanguage` |
| `created_by`       | `properties.createdBy.userPrincipalName` |
| `created_time`     | `properties.createdTime` |
| `last_modified_by` | `properties.lastModifiedBy.userPrincipalName` |
| `last_modified_time` | `properties.lastModifiedTime` |

### Example API Response

An example of the API response used by this resource can be found in the test fixture [`environment/tests/resource/Validate_Create/get_environment_00000000-0000-0000-0000-000000000001.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/environment/tests/resource/Validate_Create/get_environment_00000000-0000-0000-0000-000000000001.json).
