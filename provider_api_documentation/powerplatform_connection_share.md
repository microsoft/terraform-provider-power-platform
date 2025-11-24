# `powerplatform_connection_share`

This resource is used to manage connection shares in Power Platform.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud)                                                                                             |
| ------------------- | ----------- | --------------------------------------------------------------------------------------------------------------------------- |
| Create              | `POST`      | `https://{environment_host}/connectivity/connectors/{connectorName}/connections/{connectionId}/modifyPermissions?api-version=1&$filter=environment eq '{environmentId}'` |
| Read                | `GET`       | `https://{environment_host}/connectivity/connectors/{connectorName}/connections/{connectionId}/permissions?api-version=1&$filter=environment eq '{environmentId}'` |
| Update              | `POST`      | `https://{environment_host}/connectivity/connectors/{connectorName}/connections/{connectionId}/modifyPermissions?api-version=1&$filter=environment eq '{environmentId}'` |
| Delete              | `POST`      | `https://{environment_host}/connectivity/connectors/{connectorName}/connections/{connectionId}/modifyPermissions?api-version=1&$filter=environment eq '{environmentId}'` |

## Attribute Mapping

| Resource Attribute | API Response JSON Field |
| ------------------ | ----------------------- |
| `id`               | `id`                    |
| `environment_id`   | -                       |
| `connector_name`   | -                       |
| `connection_id`    | -                       |
| `role_name`        | `properties.roleName`   |
| `principal.entra_object_id` | `properties.principal.id` |
| `principal.display_name` | `properties.principal.displayName` |

### Example API Response

An example of the API response can be found [here](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/connection/tests/resource/connection_shares/Validate_Create/get_connection_shares.json).
