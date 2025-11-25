## powerplatform_user

This resource associates a user to a Power Platform environment and manages their Dataverse or environment-level security roles.

## API Endpoints

| Operation | Method | URL |
| --------- | ------ | --- |
| Read user | GET | `https://{dataverse_org_url}/api/data/v9.2/systemusers({system_user_id})` |
| List users (lookup) | GET | `https://{dataverse_org_url}/api/data/v9.2/systemusers?$filter=azureactivedirectoryobjectid eq {aad_object_id}` |
| List role assignments (environment) | GET | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/{environment_id}/roleAssignments?api-version=2020-10-01` |
| Modify role assignments (environment) | POST | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/{environment_id}/modifyRoleAssignments?api-version=2020-10-01` |

## Attribute Mapping

| Data Source Attribute   | API Response JSON Field              |
| ----------------------- | ------------------------------------ |
| `id`                    | `systemuserid`                       |
| `environment_id`        | *(path parameter – not in body)*     |
| `aad_id`                | `azureactivedirectoryobjectid`       |
| `business_unit_id`      | `_businessunitid_value`              |
| `user_principal_name`   | `domainname`                         |
| `first_name`            | `firstname`                          |
| `last_name`             | `lastname`                           |
| `security_roles`        | `systemuserroles_association[*].roleid` |
| `disable_delete`        | *(Terraform-only flag, not in body)* |

### Example API Response

The following example shows a Dataverse `systemuser` record used by this resource:

- [get_systemuser_00000000-0000-0000-0000-000000000002.json](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/authorization/tests/resource/user/Validate_Create/get_systemuser_00000000-0000-0000-0000-000000000002.json)
