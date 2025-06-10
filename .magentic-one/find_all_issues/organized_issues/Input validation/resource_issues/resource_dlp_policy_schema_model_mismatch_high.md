# Title

Struct Field and Schema Attribute Name Mismatches

##
/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy.go

## Problem

Across multiple CRUD operations (Create, Read, Update, Delete), and the validation logic, there is a recurring mismatch between Go struct/model field names and Terraform resource schema attributes. Sometimes the code uses names like `BusinessGeneralConnectors`, `NonBusinessConfidentialConnectors` which do not correspond to the schema attributes (`business_connectors`, `non_business_connectors`). This is likely due to Go struct fields being named differently than the schema attributes, and the downstream usage (converters, API DTOs) required field names to match the schema. This increases cognitive complexity, introduces subtle bugs, and makes future maintenance difficult.

## Impact

Severity: High

Severe maintainability and reliability risk. Using mismatched names between Terraform schema and Go struct fields will lead to bugs, field omission, data conversion failures, and runtime panics that are difficult to debug.

## Location

Throughout the file, especially:

- In ValidateConfig: using `config.BusinessGeneralConnectors`
- In Read/Create/Update: using `convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)` (semantic confusion)
- In Schema: attribute is `business_connectors`, struct field may be `BusinessGeneralConnectors` or similar.

## Code Issue

```go
// Schema uses "business_connectors"
"business_connectors": schema.SetNestedAttribute{...}

// Elsewhere, uses BusinessGeneralConnectors, NonBusinessConfidentialConnectors, etc.
plan.BusinessGeneralConnectors // or config.BusinessGeneralConnectors
```

## Fix

Ensure that struct field names in the resource model (`dataLossPreventionPolicyResourceModel`) exactly match the resource schema attribute names (`business_connectors`, `non_business_connectors`, etc.) and update all usage accordingly. Refactor conversion helpers to align semantically and reduce confusion by creating a direct mapping between schema, struct, and API DTO.

```go
// In your resourceModel struct
type dataLossPreventionPolicyResourceModel struct {
  // ...
  BusinessConnectors types.Set   `tfsdk:\"business_connectors\"`
  NonBusinessConnectors types.Set `tfsdk:\"non_business_connectors\"`
  BlockedConnectors types.Set     `tfsdk:\"blocked_connectors\"`
  // ...
}
```

---
Save as:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/resource_dlp_policy_schema_model_mismatch_high.md
