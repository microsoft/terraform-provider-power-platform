# `powerplatform_solution_checker_rules` (Data Source)

This data source is used to fetch the list of Solution Checker rules for a specific Power Platform environment. Solution Checker analyzes solutions against a ruleset of best practices and returns detailed rule definitions that can be used to configure Managed Environments.

## API Endpoints

The data source uses the environment's PowerApps Advisor endpoint to retrieve rules. It first resolves the environment details, then calls the advisor API:

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud) |
| ------------------- | ----------- | -------------------------------- |
| Read                | `GET`       | `https://{powerapp_advisor_host}/api/rule?ruleset={ruleset_id}&api-version=2.0` |

Where:

- `{powerapp_advisor_host}` is derived from the environment's `properties.runtimeEndpoints.microsoft.PowerAppsAdvisor` URL.
- `{ruleset_id}` is a fixed identifier defined by the provider (Solution Checker ruleset ID).

## Attribute Mapping

| Data Source Attribute                  | API Response JSON Field |
| -------------------------------------- | ----------------------- |
| `environment_id`                       | Used to look up the environment via the Environments API; not part of the rules API payload. |
| `rules`                                | Entire rules array returned by the API response. |
| `rules[*].code`                        | `code` |
| `rules[*].description`                 | `description` |
| `rules[*].summary`                     | `summary` |
| `rules[*].how_to_fix`                  | `howToFix` |
| `rules[*].guidance_url`                | `guidanceUrl` |
| `rules[*].component_type`              | `componentType` |
| `rules[*].primary_category`            | `primaryCategory` |
| `rules[*].primary_category_description`| Derived from `primaryCategory` using a local mapping (for example, `0` → `Error`, `1` → `Performance`, etc.). |
| `rules[*].include`                     | `include` |
| `rules[*].severity`                    | `severity` |

### Example API Response

An example of the API response used by this data source (showing two sample Solution Checker rules) can be found in the test fixture [`solution_checker_rules/tests/datasource/Validate_Read/get_rules.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/solution_checker_rules/tests/datasource/Validate_Read/get_rules.json).
