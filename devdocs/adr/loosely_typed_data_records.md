# ADR: Loosely Typed Resources (`powerplatform_data_record` and `powerplatform_data_records`)

## Status

Accepted

## Summary

This document outlines the rationale for implementing loosely typed Terraform resources (such as `powerplatform_rest`, `powerplatform_data_record`, `powerplatform_data_records`) within the Terraform Provider for Power Platform. The design decision was inspired by the AzAPI Terraform provider and was driven by the need to support Dataverse configuration settings (e.g., application users, roles, business units, teams) in a flexible manner.

## Context

Terraform providers traditionally define strongly typed resources, where each resource corresponds to a specific API entity with a fixed schema. This approach provides strict validation, clear documentation, and improved state management. However, it can become restrictive when APIs evolve frequently or are highly customizable, as is the case with Dataverse in Microsoft Power Platform.

Dataverse administrative settings (such as application users, roles, business units, and teams) are stored in flexible, schema-driven tables. Maintaining strongly typed Terraform resources for each Dataverse entity would require frequent provider updates and significant Go development effort, slowing down the adoption of new features.

## Decision

We decided to implement loosely typed Terraform resources, specifically:

- `powerplatform_data_record`
- `powerplatform_data_records`

These resources use a flexible schema (key-value maps) to support arbitrary Dataverse entities and attributes. This design allows Terraform users to configure Dataverse settings dynamically without requiring explicit provider updates or Go-based contributions.

## Rationale

- **Reduced Maintenance Overhead:** Inspired by the Azure AzAPI Terraform provider, this approach significantly reduces the maintenance burden. It eliminates the need for frequent provider updates whenever Dataverse APIs evolve or new configuration options become available.
- **Faster Time-to-Market:** Contributors can immediately leverage new Dataverse features without waiting for provider updates.
- **Contributor Accessibility:** Terraform module developers familiar with HCL but not Go can easily define complex Dataverse configurations, broadening the contributor base.
- **Limited Scope:** This decision specifically targets Dataverse configuration and generic REST interactions. Strongly typed resources remain the preferred approach for stable, well-defined Power Platform entities.

## Comparison: Strongly Typed vs. Loosely Typed Resources

- **Strongly Typed Approach:**
    - Each resource has predefined attributes in Terraform’s schema, aligning closely with the API contract.
    - **Benefits:** Strong validation, better documentation, improved state management.
    - **Drawbacks:** Requires provider updates for new API features; less flexible for evolving platforms like Dataverse.
- **Loosely Typed Approach (used in `powerplatform_rest`, `powerplatform_data_record`, `powerplatform_data_records`):**
    - Uses a flexible schema, often relying on key-value maps (`map[string]interface{}`) to support arbitrary API attributes.
    - **Benefits:** Immediate support for all Dataverse settings, faster iteration, and reduced provider maintenance overhead.
    - **Drawbacks:** Less strict validation; requires users to refer to external API documentation for field definitions; potential for misuse.

## Consequences

### Positive

- **Scalability:** Provider maintenance is simplified as Microsoft continues to evolve Dataverse APIs.
- **Flexibility:** Users gain immediate access to new Dataverse features and configurations.
- **Community Contribution:** Easier for community members to contribute Terraform modules without deep Go expertise.

### Negative

- **Validation Challenges:** Terraform cannot perform deep schema validation on dynamic attributes, potentially leading to runtime errors if incorrect attributes are provided.
- **Documentation Dependency:** Users must rely on external Dataverse API documentation to understand available fields and configurations clearly.

## Future Enhancements

To mitigate validation and usability challenges, future improvements may include:

- Auto-generating schemas from Dataverse API metadata.
- Developing validation plugins or tooling to assist users in creating correct configurations.

## Conclusion

Adopting loosely typed resources (`powerplatform_data_record` and `powerplatform_data_records`) balances flexibility, ease of contribution, and maintainability with acceptable trade-offs in validation and documentation. This decision aligns with established best practices from the AzAPI provider and supports a broader Terraform user base.
