# Title

Missing Default Value Handling for `ExpectedHttpStatus`

##

/workspaces/terraform-provider-power-platform/internal/services/rest/datasource_rest_query.go

## Problem

The `Read` function lacks default value handling for the `ExpectedHttpStatus` field. If `ExpectedHttpStatus` is not defined, the function fails to provide a sensible default (e.g., HTTP status code 200). This omission can lead to unexpected behavior or failures when the field isn’t explicitly set by the user.

## Impact

The issue can affect runtime operations, potentially causing unnecessary failures. Severity is **high**, as it directly impacts the functionality and reliability of the API request execution.

## Location

The issue is in the `Read` function, within the block where `ExpectedHttpStatus` is processed.

## Code Issue

```go
// // If the expected status code is not provided, default to 200
// if state.ExpectedHttpStatus == nil {
// 	state.ExpectedHttpStatus = []int{200}
// }
```

Commented-out code suggests that this feature was intended but is not implemented.

## Fix

Uncomment and properly implement the default value logic for `ExpectedHttpStatus`.

```go
if state.ExpectedHttpStatus == nil {
    state.ExpectedHttpStatus = []int{200}
}
```

This ensures that requests without an explicitly defined `ExpectedHttpStatus` will default to the standard success status code, improving reliability and reducing user error.