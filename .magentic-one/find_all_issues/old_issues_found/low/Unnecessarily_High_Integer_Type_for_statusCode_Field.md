# Title

Unnecessarily High Integer Type for statusCode Field

##

/workspaces/terraform-provider-power-platform/internal/services/application/dto.go

## Problem

The `statusCode` field in the `tenantApplicationErrorDetailsDto` struct is defined as `int64`. However, API error codes typically fall within a smaller range, making the use of `int64` excessive and unnecessary.

## Impact

Using an overly large integer type consumes more memory and may cause unnecessary overhead in serialized representations (e.g., JSON or network transmissions). This is a minor optimization issue but can accumulate impact when scaled over many instances (e.g., during high-volume data processing).

Severity: **Low**

## Location

tenantApplicationErrorDetailsDto struct.

## Code Issue

Current field definition:

```go
type tenantApplicationErrorDetailsDto struct {
    ErrorCode  string `json:"errorCode"`
    ErrorName  string `json:"errorName"`
    Message    string `json:"message"`
    Source     string `json:"source"`
    StatusCode int64  `json:"statusCode"`
    Type       string `json:"type"`
}
```

## Fix

Change the `statusCode` field type from `int64` to `int`, which is sufficient to handle error codes.

```go
type tenantApplicationErrorDetailsDto struct {
    ErrorCode  string `json:"errorCode"`
    ErrorName  string `json:"errorName"`
    Message    string `json:"message"`
    Source     string `json:"source"`
    StatusCode int    `json:"statusCode"` // Updated from int64 to int
    Type       string `json:"type"`
}
```

This adjustment reduces memory usage and serialization size slightly while maintaining the functionality. Save the markdown file with the above analysis.