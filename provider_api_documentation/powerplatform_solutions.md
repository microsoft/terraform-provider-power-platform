# `powerplatform_solutions` (data source)

This data source is used to read all solutions in a Power Platform environment.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud)                                                                                             |
| ------------------- | ----------- | --------------------------------------------------------------------------------------------------------------------------- |
| Read                | `GET`       | `https://{environment_host}/api/data/v9.2/solutions?$expand=publisherid&$filter=(isvisible eq true)&$orderby=createdon desc` |

## Attribute Mapping

| Data Source Attribute | API Response JSON Field |
| ---------------------- | ----------------------- |
| `solutions`            | `value`                 |
| `solutions.id`         | `solutionid`            |
| `solutions.display_name` | `friendlyname`        |
| `solutions.name`       | `uniquename`            |
| `solutions.created_time` | `createdon`           |
| `solutions.modified_time` | `modifiedon`         |
| `solutions.install_time` | `installedon`         |
| `solutions.version`    | `version`               |
| `solutions.is_managed` | `ismanaged`             |
| `solutions.environment_id` | - (from provider/environment context) |

### Example API Response

An example of the API response used by this data source (showing a list of solutions returned from Dataverse) can be found in the test fixture [`solution/tests/datasource/Validate_Read/get_solution.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/solution/tests/datasource/Validate_Read/get_solution.json).
