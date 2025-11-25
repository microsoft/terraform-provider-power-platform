# `powerplatform_enterprise_policy`

This resource is used to manage enterprise policy environment assignments in Power Platform.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud)                                                                                             |
| ------------------- | ----------- | --------------------------------------------------------------------------------------------------------------------------- |
| Create              | `POST`      | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/{environmentId}/enterprisePolicies/{policyType}/link?api-version=2022-06-01-preview` |
| Read                | `GET`       | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/{environmentId}?$expand=properties/enterprisePolicies&api-version=2023-06-01` |
| Delete              | `POST`      | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/{environmentId}/enterprisePolicies/{policyType}/unlink?api-version=2022-06-01-preview` |

## Attribute Mapping

| Resource Attribute | API Response JSON Field |
| ------------------ | ----------------------- |
| `id`               | -                       |
| `environment_id`   | -                       |
| `system_id`        | `properties.enterprisePolicies.vnets.systemId` or `properties.enterprisePolicies.customerManagedKeys.systemId` |
| `policy_type`      | -                       |

### Example API Response

An example of the API response can be found [here](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/enterprise_policy/tests/Validate_Create/get_environment_00000000-0000-0000-0000-000000000001.json).
