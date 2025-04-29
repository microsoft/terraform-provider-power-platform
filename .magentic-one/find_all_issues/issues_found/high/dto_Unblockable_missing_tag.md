# Title

Omission of JSON Tags for Struct Field `Unblockable`

##
`/workspaces/terraform-provider-power-platform/internal/services/connectors/dto.go`

## Problem

The field `Unblockable` in the `connectorPropertiesDto` struct lacks a JSON tag. This will cause issues during JSON serialization and deserialization. Any data corresponding to this field will not be mapped correctly.

## Impact

This issue impacts the correctness of JSON handling for the struct. Serialization and deserialization may fail, or data might not be included in the JSON result, leading to potential bugs and data loss when using APIs. Severity: **high**

## Location

Found in the struct definition of `connectorPropertiesDto`.

## Code Issue

Here is the problematic code:

```go
type connectorPropertiesDto struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Tier        string `json:"tier"`
	Publisher   string `json:"publisher"`
	Unblockable bool
}
```

## Fix

To fix this issue, add a JSON tag to the `Unblockable` field. Here is the corrected code:

```go
type connectorPropertiesDto struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Tier        string `json:"tier"`
	Publisher   string `json:"publisher"`
	Unblockable bool   `json:"unblockable"`
}
```

The addition of the JSON annotation ensures that `Unblockable` maps correctly when converting to or from JSON.