
## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud) |
| ------------------- | ----------- | -------------------------------- |
| Read                | `GET`       | `https://api.powerapps.com/providers/Microsoft.PowerApps/apps?api-version=2016-11-01&$filter=environment eq '{environment_id}'` |

## Attribute Mapping

| Data Source Attribute | API Response JSON Field |
| --------------------- | ----------------------- |
| `powerapps`           | `value`                 |
| `powerapps.id`        | `id`                    |
| `powerapps.name`      | `name`                  |
| `powerapps.display_name` | `properties.displayName` |
| `powerapps.created_time` | `properties.createdTime` |

## Attribute Mapping

An example of the API response used by this data source can be found in the test fixture [`powerapps/tests/Validate_Read/get_apps_00000000-0000-0000-0000-000000000001.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/powerapps/tests/Validate_Read/get_apps_00000000-0000-0000-0000-000000000001.json).
