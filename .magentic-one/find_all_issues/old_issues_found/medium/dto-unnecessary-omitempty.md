# Title

Unnecessary use of `omitempty` for fields that are mandatory

##

`/workspaces/terraform-provider-power-platform/internal/services/environment_groups/dto.go`

## Problem

The `omitempty` option in the `json` struct tags is being applied to fields that are supposed to be mandatory, such as `Id` in `environmentGroupPrincipalDto`. This tag causes the field to be excluded from the JSON representation if its value is empty, which can lead to unexpected behaviors during serialization and deserialization.

## Impact

The `omitempty` tag can lead to incomplete or incorrect data being transmitted, particularly in the case of required fields. This issue can cause downstream systems or processes to fail when they expect these fields to be present. **Severity: Medium**

## Location

Struct definitions across the file.

## Code Issue

```go
type environmentGroupPrincipalDto struct {
	Id   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}
```

## Fix

Remove the `omitempty` tag from fields that are mandatory. If `Id` is required, it should always be included in the JSON representation, regardless of whether it has a value.

```go
type environmentGroupPrincipalDto struct {
	Id   string `json:"id"`
	Type string `json:"type,omitempty"`
}
```
