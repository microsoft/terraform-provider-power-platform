# `powerplatform_environment_application_admin`

This resource is used to manage application administrators in a specific Power Platform environment.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud)                                                                |
| ------------------- | ----------- | ---------------------------------------------------------------------------------------------- |
| Create              | `POST`      | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/{environmentId}/addAppUser` |
| Read                | `GET`       | `https://{environment_host}/api/data/v9.2/systemusers`                                         |
| Delete (deactivate) | `PATCH`     | `https://{environment_host}/api/data/v9.0/systemusers({systemUserId})`                         |
| Delete              | `DELETE`    | `https://{environment_host}/api/data/v9.2/systemusers({systemUserId})`                         |

## Attribute Mapping

| Resource Attribute | JSON Path in API Response                       |
| ------------------ | ----------------------------------------------- |
| `application_id`   | `value[*].applicationid`                        |
| `environment_id`   | not present in body; provided by `environment_id` |
| `id`               | composed as `{environment_id}/{application_id}` |

### Example API Response

An example of the `systemusers` API response used by this resource can be found in the test fixture [`application_admin/Create/get_applicationusers.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/application/tests/resource/application_admin/Create/get_applicationusers.json).
