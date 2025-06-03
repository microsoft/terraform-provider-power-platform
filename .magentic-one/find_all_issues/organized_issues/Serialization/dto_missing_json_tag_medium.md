# Missing JSON Tag on DTO Field

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/dto.go

## Problem

The struct `linkedEnvironmentIdMetadataDto` has a field `InstanceURL` with no JSON tag. This may cause incorrect mapping if the JSON payload uses a different casing or name for the instance URL, which breaks deserialization/serialization.

## Impact

May cause JSON parsing issues when the field name does not match exactly what is in the JSON payload (Go's default is to lowercase the struct field name for JSON mapping, but if the server responds with a different case, this will break).  
**Severity:** medium.

## Location

```go
type linkedEnvironmentIdMetadataDto struct {
    InstanceURL string
}
```

## Code Issue

```go
InstanceURL string
```

## Fix

Add a JSON tag reflecting the actual key in the JSON response (e.g. `instanceUrl`). Adjust according to the actual API.

```go
InstanceURL string `json:"instanceUrl"`
```
