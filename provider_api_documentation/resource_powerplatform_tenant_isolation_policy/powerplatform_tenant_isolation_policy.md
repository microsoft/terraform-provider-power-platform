# powerplatform_tenant_isolation_policy

This resource is used to configure the tenant isolation policy for a Power Platform tenant. The policy controls which external tenants can interact with resources in the current tenant and in which directions (inbound/outbound).

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud) |
| ------------------- | ----------- | -------------------------------- |
| Read                | `GET`       | `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/{tenant_id}/tenantIsolationPolicy` |
| Create / Update     | `PUT`       | `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/{tenant_id}/tenantIsolationPolicy` |

Where `{tenant_id}` is the Azure AD tenant ID whose isolation policy is being managed.

## Attribute Mapping

| Resource Attribute                             | API Request / Response JSON Field          |
| ---------------------------------------------- | ------------------------------------------ |
| `tenant_id`                                    | `properties.tenantId`                      |
| `is_disabled`                                  | `properties.isDisabled`                    |
| `allowed_tenants`                              | `properties.allowedTenants`                |
| `allowed_tenants[*].tenant_id`                 | `properties.allowedTenants[*].tenantId`    |
| `allowed_tenants[*].direction.inbound`         | `properties.allowedTenants[*].direction.inbound`  |
| `allowed_tenants[*].direction.outbound`        | `properties.allowedTenants[*].direction.outbound` |

### Example API Response

An example of the API response used by this resource (showing a tenant isolation policy with a single allowed tenant) can be found in the test fixture [`tenant_isolation_policy/tests/resource/Validate_Create/policy_response.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/tenant_isolation_policy/tests/resource/Validate_Create/policy_response.json).
