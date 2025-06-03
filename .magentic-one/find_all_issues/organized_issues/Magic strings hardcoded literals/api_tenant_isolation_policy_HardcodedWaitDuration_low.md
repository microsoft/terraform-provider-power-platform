# Issue: Hardcoding of Wait Duration Bounds in getRetryAfterDuration

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/api_tenant_isolation_policy.go

## Problem

The `getRetryAfterDuration` function hardcodes the minimum (2 seconds) and maximum (60 seconds) wait duration for polling retries directly in the implementation. Hardcoded magic numbers reduce flexibility, can violate project constants standards, and make adjustment harder if retry policy needs tuning. These bounds are duplicated here instead of being documented or externally configured.

## Impact

Low. While it does not directly cause incorrect results, it makes future adjustments riskier and less transparent, and reduces maintainability.

## Location

Within `getRetryAfterDuration`:

```go
if duration < 2*time.Second {
    duration = 2 * time.Second
} else if duration > 60*time.Second {
    duration = 60 * time.Second
}
```

## Code Issue

```go
if duration < 2*time.Second {
    duration = 2 * time.Second
} else if duration > 60*time.Second {
    duration = 60 * time.Second
}
```

## Fix

Define constants at the top of the file or in a shared package, and reference these in the function for clarity and centralized control. Example:

```go
const (
    MinRetryAfterDuration = 2 * time.Second
    MaxRetryAfterDuration = 60 * time.Second
)
```

Then:

```go
if duration < MinRetryAfterDuration {
    duration = MinRetryAfterDuration
} else if duration > MaxRetryAfterDuration {
    duration = MaxRetryAfterDuration
}
```
---
