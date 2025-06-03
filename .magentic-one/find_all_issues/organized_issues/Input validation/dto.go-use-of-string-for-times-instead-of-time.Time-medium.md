# Use of `string` for Times Instead of `time.Time`

##

/workspaces/terraform-provider-power-platform/internal/services/powerapps/dto.go

## Problem

The fields `CreatedTime`, `LastModifiedTime`, and `LastPublishTime` in `powerAppPropertiesBapiDto` are typed as `string`. In Go, it's best practice to use the `time.Time` type for date/time data to ensure proper parsing, validation, and time operations, instead of using raw strings.

## Impact

Severity: Medium

This may lead to data inconsistency, missed validation (like parsing), inability to use Go's time utilities, and bugs related to time formatting, localization, or comparison.

## Location

```go
type powerAppPropertiesBapiDto struct {
    // ...
    CreatedTime      string                 `json:"createdTime"`
    LastModifiedTime string                 `json:"lastModifiedTime"`
    LastPublishTime  string                 `json:"lastPublishTime"`
    // ...
}
```

## Code Issue

```go
CreatedTime      string                 `json:"createdTime"`
LastModifiedTime string                 `json:"lastModifiedTime"`
LastPublishTime  string                 `json:"lastPublishTime"`
```

## Fix

Use `time.Time` instead of `string` and, if necessary, provide custom JSON unmarshal/marshal logic if the time format is not RFC3339.

```go
import "time"

type powerAppPropertiesBapiDto struct {
    // ...
    CreatedTime      time.Time              `json:"createdTime"`
    LastModifiedTime time.Time              `json:"lastModifiedTime"`
    LastPublishTime  time.Time              `json:"lastPublishTime"`
    // ...
}
```

If a custom date/time format is used, implement `UnmarshalJSON` and `MarshalJSON` methods for proper parsing.
