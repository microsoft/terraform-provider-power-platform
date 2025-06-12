# Missing 'omitempty' on Required JSON Field

## /workspaces/terraform-provider-power-platform/internal/services/environment/dto.go

## Problem

In the structure `EnviromentPropertiesDto`, the field `Description` does not have an `omitempty` in its JSON tag, which means it will always be present in the marshaled JSON, even if set to the empty string. Most other fields use `omitempty`, suggesting this was likely unintentional.

## Impact

This causes inconsistent API responses and can be confusing. API consumers may expect similar field presence/absence semantics for all optional fields. This is a medium-severity data consistency issue.

## Location

- Struct `EnviromentPropertiesDto`, field `Description`, likely around line 51

## Code Issue

```go
Description string `json:"description"`
```

## Fix

Update this field to use `omitempty` for consistency:

```go
Description string `json:"description,omitempty"`
```
