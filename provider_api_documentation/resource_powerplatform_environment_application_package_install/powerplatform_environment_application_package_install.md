# `powerplatform_environment_application_package_install`

This resource is used to install an application package in a Power Platform environment.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud) |
| ------------------- | ----------- | -------------------------------- |
| Create              | `POST`      | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/{environment_id}/applicationPackages/{application_id}/install?api-version=2022-03-01-alpha` |
| Read                | `GET`       | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/{environment_id}/applicationPackages/{application_id}?api-version=2022-03-01-alpha` |
| Update              | `POST`      | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/{environment_id}/applicationPackages/{application_id}/install?api-version=2022-03-01-alpha` |
| Delete              | `DELETE`    | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/{environment_id}/applicationPackages/{application_id}?api-version=2022-03-01-alpha` |

## Attribute Mapping

| Resource Attribute   | API Response JSON Field        |
| -------------------- | ------------------------------ |
| `environment_id`     | `environmentId`                |
| `package_id`         | `id`                           |
| `status`             | `properties.status`            |
| `error`              | `properties.error`             |
| `last_modified_time` | `properties.lastModifiedTime`  |


