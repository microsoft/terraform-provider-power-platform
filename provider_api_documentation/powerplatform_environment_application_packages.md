# `powerplatform_environment_application_packages` (data source)

This data source is used to read all application packages in a Power Platform environment.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud) |
| ------------------- | ----------- | -------------------------------- |
| Read                | `GET`       | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/{environment_id}/applicationPackages?api-version=2022-03-01-alpha` |

## Attribute Mapping

| Resource Attribute | API Response JSON Field |
| ------------------ | ----------------------- |
| `packages`         | `value`                 |
| `packages.id`      | `id`                    |
| `packages.name`    | `name`                  |
| `packages.properties.displayName` | `properties.displayName` |
| `packages.properties.description` | `properties.description` |
| `packages.properties.publisherName` | `properties.publisherName` |
| `packages.properties.learnMoreUrl` | `properties.learnMoreUrl` |
| `packages.properties.applicationId` | `properties.applicationId` |
| `packages.properties.uniqueName` | `properties.uniqueName` |
| `packages.properties.version` | `properties.version` |
| `packages.properties.status` | `properties.status` |
| `packages.properties.error` | `properties.error` |
| `packages.properties.lastModifiedTime` | `properties.lastModifiedTime` |

### Example API Response

An example of the API response used by this data source can be found in the test fixture [`environment_application_packages/Validate_Read/get_applications.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/application/tests/datasource/environment_application_packages/Validate_Read/get_applications.json).
