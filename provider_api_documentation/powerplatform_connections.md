# `powerplatform_connections` (data source)

This data source is used to read all connections in a Power Platform environment.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud)                                                                                             |
| ------------------- | ----------- | --------------------------------------------------------------------------------------------------------------------------- |
| Read                | `GET`       | `https://{environment_host}/connectivity/connections?api-version=1`                                                         |

## Attribute Mapping

| Resource Attribute          | API Response JSON Field |
| --------------------------- | ----------------------- |
| `id`                        | `id`                    |
| `name`                      | `name`                  |
| `display_name`              | `properties.displayName`|
| `status`                    | `properties.statuses`   |
| `connection_parameters`     | `properties.connectionParameters` |
| `connection_parameters_set` | `properties.connectionParametersSet` |

### Example API Response

An example of the API response can be found [here](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/connection/tests/datasource/connections/Validate_Read/get_connections.json).
