# powerplatform_tenant (Data Source)

This data source is used to read information about the current Power Platform tenant, including location, geos, and compliance-related flags.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud) |
| ------------------- | ----------- | -------------------------------- |
| Read                | `GET`       | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/tenant?api-version=2021-04-01` |

## Attribute Mapping

| Data Source Attribute                      | API Response JSON Field |
| ------------------------------------------ | ----------------------- |
| `tenant_id`                                | `tenantId`              |
| `state`                                    | `state`                 |
| `location`                                 | `location`              |
| `aad_country_geo`                          | `aadCountryGeo`         |
| `data_storage_geo`                         | `dataStorageGeo`        |
| `default_environment_geo`                  | `defaultEnvironmentGeo` |
| `aad_data_boundary`                        | `aadDataBoundary`       |
| `fedramp_high_certification_required`      | `fedRAMPHighCertificationRequired` |

### Example API Response

An example of the API response used by this data source can be found in the test fixture [`tenant/tests/datasource/Validate_Read/get_tenant.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/tenant/tests/datasource/Validate_Read/get_tenant.json).
