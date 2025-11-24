# `powerplatform_data_loss_prevention_policy`

This resource is used to manage a Data Loss Prevention (DLP) policy in Power Platform.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud)                                        |
| ------------------- | ----------- | ---------------------------------------------------------------------- |
| Create              | `POST`      | `https://<geo>.bap.microsoft.com/providers/PowerPlatform.Governance/v2/policies` |
| Read                | `GET`       | `https://<geo>.bap.microsoft.com/providers/PowerPlatform.Governance/v2/policies/{name}` |
| Update              | `PATCH`     | `https://<geo>.bap.microsoft.com/providers/PowerPlatform.Governance/v2/policies/{name}` |
| Delete              | `DELETE`    | `https://<geo>.bap.microsoft.com/providers/PowerPlatform.Governance/v2/policies/{name}` |

## Attribute Mapping

| Resource Attribute                      | API Response JSON Field                                  |
| --------------------------------------- | -------------------------------------------------------- |
| `id`                                    | `policyDefinition.name`                                  |
| `display_name`                          | `policyDefinition.displayName`                           |
| `default_connectors_classification`     | `policyDefinition.defaultConnectorsClassification`       |
| `environment_type`                      | `policyDefinition.environmentType`                       |
| `created_by`                            | `policyDefinition.createdBy.displayName`                 |
| `created_time`                          | `policyDefinition.createdTime`                           |
| `last_modified_by`                      | `policyDefinition.lastModifiedBy.displayName`            |
| `last_modified_time`                    | `policyDefinition.lastModifiedTime`                      |
| `environments`                          | `policyDefinition.environments`                          |
| `non_business_connectors`               | `policyDefinition.connectorGroups` (where classification is `General`) |
| `business_connectors`                   | `policyDefinition.connectorGroups` (where classification is `Confidential`) |
| `blocked_connectors`                    | `policyDefinition.connectorGroups` (where classification is `Blocked`) |
| `custom_connectors_patterns`            | `customConnectorUrlPatternsDefinition.rules`             |

### Example API Response

An example of the API response can be found [here](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/dlp_policy/tests/resource/Validate_Create/get_policy_00000000-0000-0000-0000-000000000001.json).
