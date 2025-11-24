# `powerplatform_solution`

This resource is used to import and export solutions in Power Platform environments.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud)                                                                                             |
| ------------------- | ----------- | --------------------------------------------------------------------------------------------------------------------------- |
| Create              | `POST`      | `https://{environment_host}/api/data/v9.2/StageSolution`                                                                    |
| Create              | `POST`      | `https://{environment_host}/api/data/v9.2/ImportSolutionAsync`                                                              |
| Read                | `GET`       | `https://{environment_host}/api/data/v9.2/solutions?$expand=publisherid&$filter=solutionid eq {solutionId}`                  |
| Delete              | `DELETE`    | `https://{environment_host}/api/data/v9.2/solutions({solutionId})`                                                          |
| Update              | `POST`      | `https://{environment_host}/api/data/v9.2/StageSolution`                                                                    |
| Update              | `POST`      | `https://{environment_host}/api/data/v9.2/ImportSolutionAsync`                                                              |

## Attribute Mapping

| Resource Attribute     | API Response JSON Field |
| ---------------------- | ----------------------- |
| `id`                   | `solutionid`            |
| `display_name`         | `friendlyname`          |
| `is_managed`           | `ismanaged`             |
| `solution_version`     | `version`               |
| `solution_file`        | -                       |
| `settings_file`        | -                       |
| `solution_file_checksum` | -                     |
| `settings_file_checksum` | -                     |
| `environment_id`       | -                       |

### Example API Response

An example of the API response used by this resource (showing a solution record returned from Dataverse after import) can be found in the test fixture [`solution/tests/resource/Validate_Create_With_Settings_File/get_solution.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/solution/tests/resource/Validate_Create_With_Settings_File/get_solution.json).
