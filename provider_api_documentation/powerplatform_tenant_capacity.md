# powerplatform_tenant_capacity (Data Source)

Fetches the tenant-level capacity information for Dataverse and related storage types.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud) |
| ------------------- | ----------- | -------------------------------- |
| Read                | `GET`       | `https://{licensing_host}/v0.1-alpha/tenants/{tenant_id}/TenantCapacity` |

Where:

- `{licensing_host}` is the Licensing API host configured in the provider.
- `{tenant_id}` is the Azure AD tenant ID whose capacity is being queried.

## Attribute Mapping

| Data Source Attribute                             | API Response JSON Field                    |
| ------------------------------------------------- | ------------------------------------------ |
| `tenant_id`                                       | `tenantId`                                 |
| `license_model_type`                              | `licenseModelType`                         |
| `capacities`                                      | `tenantCapacities`                         |
| `capacities[*].capacity_type`                     | `capacityType`                             |
| `capacities[*].capacity_units`                    | `capacityUnits`                            |
| `capacities[*].total_capacity`                    | `totalCapacity`                            |
| `capacities[*].max_capacity`                      | `maxCapacity`                              |
| `capacities[*].status`                            | `status`                                   |
| `capacities[*].consumption.actual`                | `consumption.actual`                       |
| `capacities[*].consumption.rated`                 | `consumption.rated`                        |
| `capacities[*].consumption.actual_updated_on`     | `consumption.actualUpdatedOn`              |
| `capacities[*].consumption.rated_updated_on`      | `consumption.ratedUpdatedOn`               |

### Example API Response

An example of the API response used by this data source (showing tenant capacity across database, file, and log storage) can be found in the test fixture [`capacity/tests/datasource/Validate_Read/get_tenant_capacity.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/capacity/tests/datasource/Validate_Read/get_tenant_capacity.json).
