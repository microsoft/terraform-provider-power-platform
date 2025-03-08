# Schema Design for New Resources and Data Sources

This guide will help you effectively contribute new resources and data sources to the Microsoft Power Platform Terraform provider. We'll discuss important schema design considerations, best practices for choosing attribute types, naming attributes, writing Markdown descriptions, implementing custom types, using plan modifiers, validators, and structuring schemas with nesting—all illustrated with specific examples from the existing provider and the Terraform Plugin Framework documentation.

## Attribute Requirements

Decide carefully if an attribute should be:

- **Required:** For example, `environment_id` and `table_logical_name` in `powerplatform_data_record` are required.
- **Optional:** Attributes like display names or settings which have default values or can be omitted.
- **Computed:** Attributes like resource IDs (`id`) or system-generated values, determined by the Power Platform API post-creation.
- **Optional + Computed:** Used for attributes that can either be specified by the user or defaulted by the system. This pattern is common when the user may optionally provide a value that the system can otherwise determine. For example, certain default configurations or auto-generated fields can follow this approach.

## Naming Attributes

Attribute naming should closely mirror terminology used by the Power Platform API or the Power Platform Admin Center UI. When discrepancies exist due to feature renaming or updates, choose the modern, user-recognizable term. For instance, prefer Dataverse instead of older, deprecated terminology like CDS to ensure users familiar with current platform features easily understand the schema.

## Markdown Description

Use Markdown descriptions (`MarkdownDescription`) in schemas exclusively to document attributes clearly and effectively. These descriptions directly generate user-facing documentation, so make sure they are comprehensive and user-friendly. Avoid using the regular description field; Markdown alone is sufficient. Include relevant links to official Power Platform product documentation where appropriate, providing additional context and guidance to users. For instance, reference the official [Power Platform documentation](https://learn.microsoft.com/power-platform/) when describing complex configurations or referencing specific product features.

## Choosing Collection Types: Sets, Lists, and Tuples

Selecting the right collection type improves readability and reduces unnecessary plan diffs.

- Generally, prefer **Sets** for unique, unordered collections, commonly used in relationship fields. For example, relationships in `powerplatform_data_record` use sets (`toset([...])`) to avoid order-related changes in Terraform plans.
- Use **Lists** sparingly and only when order or duplication matters, which is rare for Power Platform resources.
- Avoid **Tuples**, as they are typically overly rigid and complex for practical use.

## When to Consider Custom Types

Custom attribute types can add complexity, so try to use built-in Terraform types with validators first.

Reserve custom types for complex scenarios requiring special parsing or logic. For instance, the provider uses a custom UUID type (`customtypes.UUID`) specifically for handling GUID/UUID values, ensuring proper parsing and validation of identifiers across resources and data sources.

## Effectively Using Plan Modifiers

Plan modifiers fine-tune Terraform’s planned changes before they're applied:

- **`UseStateForUnknown`:** Commonly used for computed fields like resource IDs. For example, the `id` field in `powerplatform_data_record` uses this modifier to preserve existing state values. Don't use this modifier when the attribute's value might genuinely need to change or when the API explicitly resets or regenerates the value, as this could mask necessary updates.
- **`RequiresReplace`:** Clearly indicates changes that necessitate recreating the resource, such as altering `environment_id` or `table_logical_name`.
- **Dynamic Defaults:** Explicitly set predictable defaults using plan modifiers to avoid confusion in Terraform plans.

Proper use ensures Terraform plans remain predictable and user-friendly.

## Implementing Validators

Validators ensure inputs meet constraints before Terraform applies changes. Clearly differentiate validator use based on your scenario:

- **Built-in Validators:** Use these for straightforward constraints like string length, regex patterns, or predefined enumerations. For instance, attributes like `environment_type` are suitable for built-in validation since new options usually necessitate broader schema updates.
- **Resource-level Validators:** Employ these when validation logic involves multiple interdependent attributes, such as mutually exclusive flags (`allow_bing_search` vs. `allow_moving_data_across_regions`).
- **Custom Validators:** Implement custom validators only for complex validation scenarios not covered by built-in options, ensuring validations remain efficient and maintainable.

Balance strictness and flexibility by validating only what you know to be invalid. Avoid overly constraining attributes to allow flexibility for backend changes that might introduce new valid values without requiring schema updates.

## Structuring Schemas with Nested Attributes (Blocks)

Nested blocks logically group related attributes, improving clarity:

- **SingleNested Blocks:** Used for singular configurations, such as Dataverse settings within the `powerplatform_environment` resource, grouping attributes like `language_code`, `currency_code`, and `security_group_id`.
- **SetNested or ListNested Blocks:** Used for repeatable structures, such as multiple connector references or multiple similar configurations.

Avoid overly deep nesting. Clearly represent resource or data source structure without unnecessary complexity.

## Loosely Typed vs. Strongly Typed Resources

The Power Platform Terraform provider uses two approaches to defining resources and data sources: loosely typed resources that accept flexible schemas, and strongly typed resources that have well-defined, explicit fields. This guide helps developers understand these approaches and decide which one is appropriate when implementing new resources or data sources.

| Decision Factor                    | Loosely Typed Resources                      | Strongly Typed Resources                              |
| ---------------------------------- | -------------------------------------------- | ----------------------------------------------------- |
| API Stability                      | Unstable, rapidly evolving APIs              | Stable, mature APIs                                   |
| Validation and Error Handling      | Limited validation; runtime errors possible  | Comprehensive validation; fewer runtime errors        |
| Development and Maintenance Effort | Low ongoing maintenance                      | Higher initial and ongoing maintenance                |
| Usability                          | Flexible but potentially error-prone         | Highly usable with clear guidance and auto-completion |
| Speed to Adopt New Features        | Quick adoption; no need for provider updates | Slower adoption; requires provider schema updates     |
| Debugging and Troubleshooting      | Harder; errors identified at runtime         | Easier; errors identified during plan                 |

### Loosely Typed Resources

Loosely typed resources accept flexible schemas, typically represented as raw JSON or maps. This flexibility accommodates rapidly changing APIs or extensive API surfaces without requiring explicit provider schema definitions.

In the Power Platform provider, common examples of loosely typed resources include:

- **`powerplatform_rest`**: Allows arbitrary Web API calls by specifying HTTP method, URL, and JSON body without fixed schema definitions ([Power Platform REST Resource](https://microsoft.github.io/terraform-provider-power-platform/resources/rest/)).
- **`powerplatform_connection`** (Connection Parameters): Parameters are defined as JSON strings or maps, allowing users to configure connectors flexibly without the provider explicitly knowing every connector's schema ([Connection Resource](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/connection/resource_connection.go)).
- **`powerplatform_data_record`** (Data Record): Uses dynamic column maps to handle arbitrary Dataverse entities and fields, accommodating diverse or evolving data schemas.

Loosely typed resources are best suited for rapidly changing APIs or highly customizable scenarios. They are appropriate for temporary or preview features lacking explicit provider support and custom or third-party integrations demanding flexibility. Avoid loosely typed resources in scenarios requiring high stability, rigorous configuration validation, clear guidance, user-friendly auto-completion, or detailed validation during planning stages.

### Strongly Typed Resources

Strongly typed resources feature explicitly defined schemas with predetermined attributes. Each attribute is validated at Terraform's planning stage, reducing runtime errors and enhancing usability.

Within the Power Platform provider, strongly typed resources include:

- **`powerplatform_environment`** (Environment Resource): Explicit fields like display name, location, and environment type provide clear and enforceable schemas ([Environment Resource](https://microsoft.github.io/terraform-provider-power-platform/resources/environment/)).
- **`powerplatform_dlp_policy`** (Data Loss Prevention Policies): Structured definitions of connectors, rules, and environment restrictions ensure precise policy configurations ([DLP Policy Resource](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/dlp_policy/resource_dlp_policy.go)).
- **`powerplatform_tenant_settings`** (Tenant Settings): Defined categories such as audit logs, email settings, and product features ensure consistent, validated configurations ([Tenant Settings Resource](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/tenant_settings/resource_tenant_settings.go)).
- **`powerplatform_managed_environment`** (Managed Environments): Defined configurations for security, administration, and environment controls provide clear guidance and enforcement ([Managed Environment Resource](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/managed_environment/resource_managed_environment.go)).

Strongly typed resources are recommended for established, stable resources where configuration consistency is essential. They are ideal for scenarios requiring detailed validation and early error detection in the Terraform lifecycle. Additionally, they benefit core infrastructure components that rely on clear documentation and enhanced usability. Avoid strongly typed resources for APIs that frequently change, requiring constant schema updates, or highly customized scenarios that cannot effectively be captured by fixed schemas.

## Additional Resources

Throughout your development, continuously reference the [Terraform Plugin Framework Documentation](https://developer.hashicorp.com/terraform/plugin/framework) and regularly review existing implementations in the provider repository to maintain consistency and alignment with best practices.
