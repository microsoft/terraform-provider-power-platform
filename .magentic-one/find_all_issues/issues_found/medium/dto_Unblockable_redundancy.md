# Title

Field `Unblockable` is Redundant in `connectorPropertiesDto`, and Overlaps with `unblockableConnectorMetadataDto`

##
`/workspaces/terraform-provider-power-platform/internal/services/connectors/dto.go`

## Problem

The `Unblockable` field in `connectorPropertiesDto` has redundancy due to the presence of `unblockableConnectorMetadataDto`, which is explicitly designed to handle `Unblockable` as a JSON field named `"unblockable"`. This discrepancy can lead to data collisions or ambiguity in the implementation.

## Impact

Redundant fields introduce confusion and increase the potential for bugs. They also inflate the memory footprint of the application unnecessarily while complicating the business logic. Severity: **medium**

## Location

Found in the struct definition `connectorPropertiesDto` and overlaps conceptually with `unblockableConnectorMetadataDto`.

## Code Issue

Here is the problematic code:

```go
type connectorPropertiesDto struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Tier        string `json:"tier"`
	Publisher   string `json:"publisher"`
	Unblockable bool   `json:"unblockable"` // This seems redundant
}
```

## Fix

Evaluate the need for the `Unblockable` field in `connectorPropertiesDto`. If this field serves no unique purpose and overlaps with `unblockableConnectorMetadataDto`, it should be omitted entirely. Here is the corrected implementation:

```go
type connectorPropertiesDto struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Tier        string `json:"tier"`
	Publisher   string `json:"publisher"`
	// Removed redundant Unblockable field
}
```

This removal simplifies the structure and avoids unnecessary duplication.