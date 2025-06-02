# Title

Error Message Details Are Insufficient

## 

`/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/api_copilot_studio_application_insights.go`

## Problem

The error messages being returned by the functions contain little to no additional details. For instance, in the functions such as `getCopilotStudioEndpoint`, if there is an error due to an invalid `environmentId`, the error merely mentions "power virtual agents runtime endpoint is not available in the environment" without specifying the actual `environmentId` or context for debugging. This can make debugging unnecessarily difficult.

## Impact

- **Severity**: Medium
- Limited error details lead to more time spent troubleshooting.
- Provides inadequate context to end users or calling functions.

## Location

```go
return "", errors.New("power virtual agents runtime endpoint is not available in the environment")
```

## Fix

Augment the error messages to include dynamic information such as the `environmentId`.

```go
if env == nil || env.Properties == nil || env.Properties.RuntimeEndpoints == nil || env.Properties.RuntimeEndpoints.PowerVirtualAgents == "" {
    return "", fmt.Errorf("power virtual agents runtime endpoint is not available in the environment with ID: %s", environmentId)
}
```