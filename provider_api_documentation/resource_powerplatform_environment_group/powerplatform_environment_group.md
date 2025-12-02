# `powerplatform_environment_group`

This resource is used to manage environment groups in Power Platform.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud)                                                                                             |
| ------------------- | ----------- | --------------------------------------------------------------------------------------------------------------------------- |
| Create              | `POST`      | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environmentGroups?api-version=2021-10-01-preview`      |
| Read                | `GET`       | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environmentGroups/{environmentGroupId}?api-version=2021-10-01-preview` |
| Update              | `PUT`       | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environmentGroups/{environmentGroupId}?api-version=2021-10-01-preview` |
| Delete              | `DELETE`    | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environmentGroups/{environmentGroupId}?api-version=2021-10-01-preview` |

## Attribute Mapping

**403 Forbidden error when using SPN**

| Resource Attribute | API Response JSON Field |
| ------------------ | ----------------------- |
| `id`               | `id`                    |
| `display_name`     | `displayName`           |
| `description`      | `description`           |

### Example API Response

An example of the API response used by this resource can be found in the test fixture [`environment_groups/resources/get_environment_group.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/environment_groups/tests/resources/get_environment_group.json).
