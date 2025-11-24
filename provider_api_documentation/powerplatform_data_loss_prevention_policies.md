# `powerplatform_data_loss_prevention_policies` (data source)

This data source is used to read all Data Loss Prevention (DLP) policies in Power Platform.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud)                                                                                             |
| ------------------- | ----------- | --------------------------------------------------------------------------------------------------------------------------- |
| Read                | `GET`       | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/policies?api-version=2020-10-01`                         |

## Attribute Mapping

| Resource Attribute                                       | API Response JSON Field                                  |
| -------------------------------------------------------- | -------------------------------------------------------- |
| `policies`                                               | `value`                                                  |
| `policies.name`                                          | `name`                                                   |
| `policies.display_name`                                  | `properties.displayName`                                 |
| `policies.created_by`                                    | `properties.createdBy.displayName`                       |
| `policies.created_time`                                  | `properties.createdTime`                                 |
| `policies.last_modified_by`                              | `properties.lastModifiedBy.displayName`                  |
| `policies.last_modified_time`                            | `properties.lastModifiedTime`                            |
| `policies.etag`                                          | `properties.eTag`                                        |
| `policies.default_environment_connector_group_name`      | `properties.defaultEnvironmentConnectorGroup.name`       |
| `policies.default_environment_connector_group_classification` | `properties.defaultEnvironmentConnectorGroup.classification` |
| `policies.environment_connector_groups`                  | `properties.environmentConnectorGroups`                  |
| `policies.environment_connector_groups.environment_name` | `properties.environmentConnectorGroups.environment.name` |
| `policies.environment_connector_groups.connector_group_name` | `properties.environmentConnectorGroups.connectorGroup.name` |
| `policies.environment_connector_groups.connector_group_classification` | `properties.environmentConnectorGroups.connectorGroup.classification` |
| `policies.connector_groups`                              | `properties.connectorGroups`                             |
| `policies.connector_groups.name`                         | `properties.connectorGroups.name`                        |
| `policies.connector_groups.classification`               | `properties.connectorGroups.classification`              |
| `policies.connector_groups.connectors`                   | `properties.connectorGroups.connectors`                  |
| `policies.connector_groups.connectors.id`                | `properties.connectorGroups.connectors.id`               |
| `policies.connector_groups.connectors.name`              | `properties.connectorGroups.connectors.name`             |
| `policies.connector_groups.connectors.type`              | `properties.connectorGroups.connectors.type`             |

### Example API Response

An example of the API response can be found [here](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/dlp_policy/tests/datasource/Validate_Read/get_policies.json).
