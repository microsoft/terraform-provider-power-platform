# powerplatform_tenant_application_packages (Data Source)

Fetches the list of Dynamics 365 tenant level applications.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud) |
| ------------------- | ----------- | -------------------------------- |
| Read                | `GET`       | `https://api.powerplatform.com/appmanagement/applicationPackages?api-version=2022-03-01-alpha` |

## Attribute Mapping

| Data Source Attribute                             | API Response JSON Field                    |
| ------------------------------------------------- | ------------------------------------------ |
| `applications`                                    | `value`                                    |
| `applications.unique_name`                        | `uniqueName`                               |
| `applications.localized_description`              | `localizedDescription`                     |
| `applications.localized_name`                     | `localizedName`                            |
| `applications.application_id`                     | `applicationId`                            |
| `applications.application_name`                   | `applicationName`                          |
| `applications.application_descprition`            | `applicationDescription`                   |
| `applications.publisher_name`                     | `publisherName`                            |
| `applications.publisher_id`                       | `publisherId`                              |
| `applications.learn_more_url`                     | `learnMoreUrl`                             |
| `applications.catalog_visibility`                 | `catalogVisibility`                        |
| `applications.application_visibility`             | `applicationVisibility`                    |
| `applications.last_error`                         | `errorDetails`                             |
| `applications.last_error[*].error_code`           | `errorDetails.errorCode`                   |
| `applications.last_error[*].error_name`           | `errorDetails.errorName`                   |
| `applications.last_error[*].message`              | `errorDetails.message`                     |
| `applications.last_error[*].source`               | `errorDetails.source`                      |
| `applications.last_error[*].status_code`          | `errorDetails.statusCode`                  |
| `applications.last_error[*].type`                 | `errorDetails.type`                        |

### Example API Response

An example of the API response used by this data source (showing a list of tenant-level applications) can be found in the test fixture [`application/tests/datasource/tenant_application_packages/Validate_Read/get_tenant_applications.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/application/tests/datasource/tenant_application_packages/Validate_Read/get_tenant_applications.json).
