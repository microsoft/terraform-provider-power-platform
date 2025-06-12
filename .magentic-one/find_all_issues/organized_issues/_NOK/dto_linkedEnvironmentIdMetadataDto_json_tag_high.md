# Field Not Tagged for JSON Serialization

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/dto.go

## Problem

The `InstanceURL` field in the `linkedEnvironmentIdMetadataDto` struct does not have a `json` struct tag. All other DTO fields are correctly tagged, ensuring proper JSON (un)marshaling, but this omission could cause failures in marshalling and unmarshalling JSON payloads involving that field.

## Impact

**Severity: High**

Omitting a `json` tag on DTO fields leads to unexpected behavior or broken features when (un)marshaling, potentially resulting in missing or misnamed JSON fields in API calls or responses. This can create data loss, broken contract, or defects.

## Location

```go
type linkedEnvironmentIdMetadataDto struct {
	InstanceURL string
}
```

## Code Issue

```go
type linkedEnvironmentIdMetadataDto struct {
	InstanceURL string
}
```

## Fix

Add the appropriate `json` struct tag to `InstanceURL` to match the expected JSON property name.

```go
type linkedEnvironmentIdMetadataDto struct {
	InstanceURL string `json:"instanceUrl"`
}
```

Choose the property name (`instanceUrl` or similar) that matches the API contract. If the JSON field is snake_case use `instance_url`, or use whatever matches your API/consumer expectations.
