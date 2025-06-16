# Overuse of `interface{}`-like Values in Attribute Type Declarations

##

/workspaces/terraform-provider-power-platform/internal/services/environment/models.go

## Problem

In the attribute type declarations for the Dataverse object (see `attrTypesDataverseObject` and related pattern), field types are set with package types, but all values in `attrValuesProductProperties` are of type `attr.Value{}` (effectively an untyped map). This allows incorrect or unintended types to be assigned, reducing compile-time safety and possibly breaking type-sensitive consumers.

## Impact

- **Severity:** Medium
- Type mismatches may only be caught at runtime.
- Increases risk of subtle bugs that are hard to trace.
- Decreases code maintainability and safety.

## Location

```go
attrTypesDataverseObject := map[string]attr.Type{
	"url":                          types.StringType,
	"domain":                       types.StringType,
	//...
}
attrValuesProductProperties := map[string]attr.Value{}
model.Dataverse = types.ObjectNull(attrTypesDataverseObject)
```

Throughout the block:

```go
attrValuesProductProperties["linked_app_type"] = types.StringValue("")
// etc.
```

## Code Issue

```go
attrTypesDataverseObject := map[string]attr.Type{ ... }
attrValuesProductProperties := map[string]attr.Value{}
```

## Fix

Wherever possible, use strongly-typed attribute types and values that match field definitions. If code logic produces objects dynamically, enforce value assignment with clear and minimal conversion logic, or wrap in helper constructors to improve type guarantees.

Additionally, ensure that all `attr.Value` assignments and extractions honor expected types and handle errors/edge cases accordingly, e.g.:

```go
attrValuesProductProperties["url"] = types.StringValue(environmentDto.Properties.LinkedEnvironmentMetadata.InstanceURL)
// For lists:
templ, err := types.ListValue(types.StringType, values)
```

Use Go's type system to enforce correct types, even in dynamic settings.

---

**This markdown should be saved as:**  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/type_safety/models_attr_value_medium.md`
