# Schema Field Naming vs. Struct Field Naming Mismatch

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/datasource_dlp_policy.go

## Problem

Terraform Schema attributes use "business_connectors" and "non_business_connectors", but the model code in `Read` references:

```go
policyModel.BusinessGeneralConnectors = ...
policyModel.NonBusinessConfidentialConnectors = ...
```

Similarly, the call to `convertToAttrValueConnectorsGroup` uses these two fields, but the schema does not define "business_general_connectors" or "non_business_confidential_connectors".

## Impact

If the struct fields and schema field names do not match, there will be state setting/getting errors and fields may not be mapped correctly between Terraform and state. This can lead to broken resource behavior or even provider panics. Severity: **high**.

## Location

- Schema definition in `Schema` function.
- Assignment in `Read` function.

## Code Issue

```go
// Schema (correct):
"business_connectors": schema.SetNestedAttribute{ ... },
"non_business_connectors": schema.SetNestedAttribute{ ... },

// In Read:
policyModel.BusinessGeneralConnectors = ...
policyModel.NonBusinessConfidentialConnectors = ...
```

## Fix

Ensure the struct (`dataLossPreventionPolicyDatasourceModel`) fields match the exact naming style/fields from the schema, i.e., use `BusinessConnectors` and `NonBusinessConnectors`. Refactor accordingly:

```go
policyModel.BusinessConnectors = convertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
policyModel.NonBusinessConnectors = convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
```

**Explanation:**  
Field names in the model must align with the schema for Terraform-generated code to function properly.
