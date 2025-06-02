# Title

Hard-coded Default Retry Duration in `retryAfter`

##

`/workspaces/terraform-provider-power-platform/internal/api/request.go`

## Problem

The `retryAfter` function uses a hard-coded default retry duration of 5-10 seconds when the `Retry-After` header cannot be parsed. This approach does not provide flexibility or configurability, which could be a problem in environments requiring different retry durations depending on service behaviors.

## Impact

Using hard-coded retry durations reduces the adaptability of the function, particularly for environments where default retry durations must be adjusted dynamically. It limits the code’s ability to cater to different API rate-limiting strategies. Severity: Medium.

## Location

Line: `return DefaultRetryAfter()` in the `retryAfter` function.

## Code Issue

```go
retryAfter, err := time.ParseDuration(retryHeader)
if err != nil {
    // default retry after 5-10 seconds
    return DefaultRetryAfter()
}
```

## Fix

Introduce a configurable parameter or settings value for default retry durations to avoid hard-coded values. This provides flexibility and allows different retry durations to be specified dynamically.

```go
retryAfter, err := time.ParseDuration(retryHeader)
if err != nil {
    // Use a configurable retry duration instead of hard-coding
    return client.Config.DefaultRetryAfterDuration
}
```

Ensure the `Client` struct has a `Config` field with a property `DefaultRetryAfterDuration` to enable dynamic configuration. For example:

```go
type ClientConfig struct {
    DefaultRetryAfterDuration time.Duration
}

type Client struct {
    Config ClientConfig
}
```