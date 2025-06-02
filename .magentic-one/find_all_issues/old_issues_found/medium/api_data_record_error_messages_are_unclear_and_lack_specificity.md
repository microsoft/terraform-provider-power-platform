# Title

Error messages are unclear and lack specificity

##

Path to the file `/workspaces/terraform-provider-power-platform/internal/services/data_record/api_data_record.go`

## Problem

Many error messages throughout the file, such as in `GetEnvironmentHostById`, `GetEntityRelationDefinitionInfo`, and other places, are generic and do not provide enough context for diagnosing problems. For example, returning `"value field is not of type []any"` or `"ReferencedEntity field is not of type string"` does not specify the location or add useful details about what went wrong.

## Impact

Generic error messages can make debugging difficult, lead to poor user experiences, and impede the ability to understand the root cause of an issue. Severity: **Medium**

## Location

Example location in the function `GetEntityRelationDefinitionInfo`:

## Code Issue

```go
return "", errors.New("ReferencedEntity field is not of type string")
```

## Fix

Replace the generic message with one that provides more details. Include variable values and relevant context for better debugging.

```go
return "", errors.New(fmt.Sprintf("ReferencedEntity field is not of type string. RelationLogicalName: %s, EntityLogicalName: %s", relationLogicalName, entityLogicalName))
```

This approach ensures that users and developers have better insight into what caused an error.

---

Next, I will search for additional issues and document them.