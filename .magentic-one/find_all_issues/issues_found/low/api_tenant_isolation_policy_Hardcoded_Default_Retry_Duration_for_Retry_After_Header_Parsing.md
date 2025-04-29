# Title

Hardcoded Default Retry Duration for Retry-After Header Parsing

##

`/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/api_tenant_isolation_policy.go`

## Problem

The function `getRetryAfterDuration` relies on a hardcoded default retry duration (`5 seconds`), as well as hardcoded minimum (`2 seconds`) and maximum (`60 seconds`). While this ensures immediate functionality, these constants may require frequent modifications based on external API rules. Currently, they are scattered in the function, making changes error-prone.

It would be better to centralize these default values and make them configurable.

## Impact

1. Hardcoded values introduce technical debt and reduce code scalability.
2. If external API's retry expectations change, updating values across multiple instances requires additional effort.
3. Lack of centralized configuration makes understanding and modifying behavior difficult.

Severity: **Low**

## Location

```go
func getRetryAfterDuration(resp *http.Response) time.Duration {
    defaultDuration := 5 * time.Second // Hardcoded default
    if resp == nil {
        return defaultDuration
    }

    retryAfter := resp.Header.Get(constants.HEADER_RETRY_AFTER)
    if retryAfter == "" {
        return defaultDuration
    }

    seconds, err := strconv.Atoi(retryAfter)
    if err == nil && seconds > 0 {
        duration := time.Duration(seconds) * time.Second
        if duration < 2*time.Second { // Hardcoded minimum
            duration = 2 * time.Second
        } else if duration > 60*time.Second { // Hardcoded maximum
            duration = 60 * time.Second
        }

        return duration
    }

    return defaultDuration
}
```

## Fix

Introduce centralized constants or configuration values for default, minimum, and maximum retry durations. Place these constants in a shared configuration file or package (e.g., `constants`).

### Example Configuration:

```go
const (
    RETRY_DEFAULT_DURATION = 5 * time.Second
    RETRY_MIN_DURATION = 2 * time.Second
    RETRY_MAX_DURATION = 60 * time.Second
)
```

### Refactored Code:

Modify `getRetryAfterDuration` to use centralized constants:

```go
func getRetryAfterDuration(resp *http.Response) time.Duration {
    if resp == nil {
        return constants.RETRY_DEFAULT_DURATION
    }

    retryAfter := resp.Header.Get(constants.HEADER_RETRY_AFTER)
    if retryAfter == "" {
        return constants.RETRY_DEFAULT_DURATION
    }

    seconds, err := strconv.Atoi(retryAfter)
    if err == nil && seconds > 0 {
        duration := time.Duration(seconds) * time.Second
        if duration < constants.RETRY_MIN_DURATION {
            duration = constants.RETRY_MIN_DURATION
        } else if duration > constants.RETRY_MAX_DURATION {
            duration = constants.RETRY_MAX_DURATION
        }

        return duration
    }

    return constants.RETRY_DEFAULT_DURATION
}
```