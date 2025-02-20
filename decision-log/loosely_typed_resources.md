# Decision Log: Loosely Typed Resources in Terraform Power Platform Provider

## Summary

This document outlines the rationale for implementing loosely typed Terraform resources (such as `powerplatform_rest`, `powerplatform_data_record`, `powerplatform_data_records`) within the Terraform Provider for Power Platform. The design decision was inspired by the AzAPI Terraform provider and was driven by the need to support Dataverse configuration settings (e.g., application users, roles, business units, teams) in a flexible manner.

## Background

Terraform providers traditionally define strongly typed resources, where each resource corresponds to a specific API entity with a fixed schema. This approach offers strict validation and type safety, but it can be limiting when APIs are frequently changing or highly customizable.

In contrast, loosely typed resources provide a more generic interface for managing platform entities. This reduces the need for frequent provider updates, while allowing contributors who are familiar with HCL but not Go to extend functionality through Terraform modules.

## Comparison: Strongly Typed vs. Loosely Typed Resources

- **Strongly Typed Approach:**
    - Each resource has predefined attributes in Terraform’s schema, aligning closely with the API contract.
    - **Benefits:** Strong validation, better documentation, improved state management.
    - **Drawbacks:** Requires provider updates for new API features; less flexible for evolving platforms like Dataverse.
- **Loosely Typed Approach (used in `powerplatform_rest`, `powerplatform_data_record`, `powerplatform_data_records`):**
    - Uses a flexible schema, often relying on key-value maps (`map[string]interface{}`) to support arbitrary API attributes.
    - **Benefits:** Immediate support for all Dataverse settings, faster iteration, and reduced provider maintenance overhead.
    - **Drawbacks:** Less strict validation; requires users to refer to external API documentation for field definitions; potential for misuse.

## Rationale

- **AzAPI Inspiration:** The Azure AzAPI provider follows a similar pattern for managing generic Azure resources without needing an explicit Terraform provider update for each new feature. The reduced maintenance cost and time-to-market for new features were key aspects of this decision.
- **Dataverse Configuration:** Many administrative settings in Power Platform are stored in Dataverse tables. By implementing `powerplatform_data_record` and `powerplatform_data_records`, Terraform users can configure these settings without requiring Go-based contributions.
- **Flexibility for Contributors:** The loosely typed approach enables Terraform module developers to define complex Dataverse configurations purely in HCL, rather than requiring modifications to the provider’s Go code.
- **Limited Scope:** We are limiting the scope of this decision to Dataverse configuration and the generic REST providers.  Strongly typed resources and data sources are still the best option for most resources in Power Platform.  

## Future Impact

- **Scalability:** Reduces provider maintenance as Microsoft evolves Dataverse APIs.
- **Validation Challenges:** Since Terraform does not perform deep schema validation for dynamic attributes, users may experience runtime errors if incorrect attributes are provided.
- **Potential Enhancements:** Future improvements could include auto-generated schemas from API metadata or validation plugins to assist users.

## Conclusion

This approach balances flexibility and ease of contribution with the need for stability. While strongly typed resources remain valuable for well-defined entities, the loosely typed pattern enables Terraform to manage Power Platform configurations dynamically. This decision aligns with best practices from AzAPI and supports a broader Terraform user base.
