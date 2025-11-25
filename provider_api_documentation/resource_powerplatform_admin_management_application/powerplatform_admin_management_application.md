# `powerplatform_admin_management_application`

This resource is used to register a service principal as an administrator for Power Platform.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud)                                                                                             |
| ------------------- | ----------- | --------------------------------------------------------------------------------------------------------------------------- |
| Create              | `PUT`       | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/adminApplications/{clientId}?api-version=2022-03-01-preview` |
| Read                | `GET`       | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/adminApplications/{clientId}?api-version=2022-03-01-preview` |
| Delete              | `DELETE`    | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/adminApplications/{clientId}?api-version=2022-03-01-preview` |
| Update              | -           | Not Supported                                                                                                               |

## Attribute Mapping

| Resource Attribute | API Response JSON Field |
| ------------------ | ----------------------- |
| `id`               | `applicationId`         |
| `display_name`     | `displayName`           |
| `tenant_id`        | `tenantId`              |
