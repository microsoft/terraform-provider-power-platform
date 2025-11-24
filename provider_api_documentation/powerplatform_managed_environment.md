# `powerplatform_managed_environment`

This resource is used to enable or disable Managed Environments governance features for an existing Power Platform environment.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud) |
| ------------------- | ----------- | -------------------------------- |
| Create (Enable)     | `POST`      | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/{environment_id}/governanceConfiguration?api-version=2021-04-01` |
| Update              | `POST`      | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/{environment_id}/governanceConfiguration?api-version=2021-04-01` |
| Delete (Disable)    | `POST`      | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/{environment_id}/governanceConfiguration?api-version=2021-04-01` |

## Attribute Mapping

| Resource Attribute                      | API Request/Response JSON Field |
| --------------------------------------- | -------------------------------- |
| `id`                                    | - (Terraform-only identifier, typically `{environment_id}`) |
| `environment_id`                        | Path parameter in `.../environments/{environment_id}/governanceConfiguration` |
| `protection_level`                      | `protectionLevel` (request body) |
| `is_usage_insights_disabled`           | `excludeEnvironmentFromAnalysis` (request body; `true` means usage insights disabled) |
| `is_group_sharing_disabled`            | `isGroupSharingDisabled` (request body) |
| `max_limit_user_sharing`               | `maxLimitUserSharing` (request body) |
| `limit_sharing_mode`                   | `limitSharingMode` (request body) |
| `solution_checker_mode`                | `solutionCheckerMode` (request body) |
| `suppress_validation_emails`           | `suppressValidationEmails` (request body) |
| `maker_onboarding_url`                 | `makerOnboardingUrl` (request body) |
| `maker_onboarding_markdown`            | `makerOnboardingMarkdown` (request body) |
| `solution_checker_rule_overrides`      | `solutionCheckerRuleOverrides` (request body) |

> Note: The Managed Environment API does not return the full governance configuration in the lifecycle response; these fields are primarily sent in the request body. The resource uses the standard Environment APIs to fetch additional environment details when needed.

### Example API Response

An example of the environment response used together with this resource (showing an environment after Managed Environment has been enabled) can be found in the test fixture [`managed_environment/tests/resource/Validate_Create_And_Update/get_environment_create_response_extended_0.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/managed_environment/tests/resource/Validate_Create_And_Update/get_environment_create_response_extended_0.json).
