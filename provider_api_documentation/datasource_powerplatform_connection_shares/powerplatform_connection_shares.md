# `powerplatform_connection_shares` (data source)

This data source is used to read all connection shares in a Power Platform environment.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud)                                                                                             |
| ------------------- | ----------- | --------------------------------------------------------------------------------------------------------------------------- |
| Read                | `GET`       | `https://{environment_host}/connectivity/connectors/{connectorName}/connections/{connectionId}/permissions?api-version=1&$filter=environment eq '{environmentId}'` |

## Attribute Mapping

| Resource Attribute | API Response JSON Field |
| ------------------ | ----------------------- |
| `id`               | `id`                    |
| `role_name`        | `properties.roleName`   |
| `principal.entra_object_id` | `properties.principal.id` |
| `principal.display_name` | `properties.principal.displayName` |

### Example API Response

An example of the API response can be found [here](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/connection/tests/datasource/connection_shares/Validate_Read/get_connection_shares.json).
