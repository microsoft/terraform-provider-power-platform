# powerplatform_locations (Data Source)

Fetches the list of available Dynamics 365 locations. For more information see [Power Platform Geos](https://learn.microsoft.com/power-platform/admin/regions-overview)

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud) |
| ------------------- | ----------- | -------------------------------- |
| Read                | `GET`       | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01` |

## Attribute Mapping

| Data Source Attribute                             | API Response JSON Field                    |
| ------------------------------------------------- | ------------------------------------------ |
| `locations`                                       | `value`                                    |
| `locations.id`                                    | `id`                                       |
| `locations.name`                                  | `name`                                     |
| `locations.display_name`                          | `properties.displayName`                   |
| `locations.code`                                  | `properties.code`                          |
| `locations.is_default`                            | `properties.isDefault`                     |
| `locations.is_disabled`                           | `properties.isDisabled`                    |
| `locations.can_provision_database`                | `properties.canProvisionDatabase`          |
| `locations.can_provision_customer_engagement_database` | `properties.canProvisionCustomerEngagementDatabase` |
| `locations.azure_regions`                         | `properties.azureRegions`                  |

### Example API Response

An example of the API response used by this data source can be found in the test fixture [`locations/tests/datasource/Validate_Read/get_locations.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/locations/tests/datasource/Validate_Read/get_locations.json).
