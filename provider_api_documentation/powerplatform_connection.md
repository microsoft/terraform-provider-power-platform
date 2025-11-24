# `powerplatform_connection`

This resource is used to manage connections in Power Platform.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud)                                                                                             |
| ------------------- | ----------- | --------------------------------------------------------------------------------------------------------------------------- |
| Create              | `PUT`       | `https://{environment_host}/connectivity/connectors/{connectorName}/connections/{connectionId}?api-version=1&$filter=environment eq '{environmentId}'` |
| Read                | `GET`       | `https://{environment_host}/connectivity/connectors/{connectorName}/connections/{connectionId}?api-version=1&$filter=environment eq '{environmentId}'` |
| Update              | `PUT`       | `https://{environment_host}/connectivity/connectors/{connectorName}/connections/{connectionId}?api-version=1&$filter=environment eq '{environmentId}'` |
| Delete              | `DELETE`    | `https://{environment_host}/connectivity/connectors/{connectorName}/connections/{connectionId}?api-version=1&$filter=environment eq '{environmentId}'` |

## Attribute Mapping

| Resource Attribute          | API Response JSON Field |
| --------------------------- | ----------------------- |
| `id`                        | `id`                    |
| `environment_id`            | -                       |
| `name`                      | `name`                  |
| `display_name`              | `properties.displayName`|
| `status`                    | `properties.statuses`   |
| `connection_parameters`     | `properties.connectionParameters` |
| `connection_parameters_set` | `properties.connectionParametersSet` |

### Example API Response

An example of the API response can be found [here](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/connection/tests/resource/connections/Validate_Create/put_connection.json).
