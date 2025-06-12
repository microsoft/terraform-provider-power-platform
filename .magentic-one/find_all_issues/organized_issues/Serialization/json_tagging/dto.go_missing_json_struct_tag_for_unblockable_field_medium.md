# Issue 2

Missing JSON struct tag for Unblockable field

##

/workspaces/terraform-provider-power-platform/internal/services/connectors/dto.go

## Problem

The `Unblockable` field in `connectorPropertiesDto` lacks a JSON struct tag. This can lead to inconsistent marshaling and unmarshaling behaviors when working with JSON, potentially causing bugs if this struct is used with JSON APIs.

## Impact

Medium severity: Omitting the JSON tag makes this field invisible when serializing or deserializing, which can lead to subtle bugs in API integration or data storage.

## Location

Line:

```go
type connectorPropertiesDto struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Tier        string `json:"tier"`
	Publisher   string `json:"publisher"`
	Unblockable bool
}
```

## Code Issue

```go
	Unblockable bool
```

## Fix

Add an appropriate JSON struct tag to the field:

```go
	Unblockable bool `json:"unblockable"`
```

---
