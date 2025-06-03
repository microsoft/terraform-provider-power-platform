# Title

Missing JSON Tags on Public Field

##

/workspaces/terraform-provider-power-platform/internal/services/solution/dto.go

## Problem

The only field without a JSON tag in this file appears to be `InstanceURL` in the `linkedEnvironmentIdMetadataDto` struct. All other struct fields are explicitly tagged for marshaling/unmarshaling. Omitting a tag risks inconsistent casing/behavior for JSON consumers.

## Impact

May cause issues with (de)serialization or integration, especially if relying on standard lowerCamelCase-to-snake_case auto-conversion, and is inconsistent with the rest of the DTO file. Severity: low.

## Location

Line 127

## Code Issue

```go
type linkedEnvironmentIdMetadataDto struct {
	InstanceURL string
}
```

## Fix

Add a JSON tag for the field, e.g.:

```go
type linkedEnvironmentIdMetadataDto struct {
	InstanceURL string `json:"instanceUrl"`
}
```

Replace `instanceUrl` with whatever the actual remote JSON uses for this field.
