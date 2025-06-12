# API Redundant Code Issues

This document consolidates all redundant code issues found in API-related components of the Terraform Power Platform provider.

## ISSUE 1

**File:** `/workspaces/terraform-provider-power-platform/internal/services/capacity/api_capacity.go`

### Problem

Throughout the file, the receiver for the `client` struct methods is named `client`, which is both verbose and redundant. Go convention is to use a one- or two-letter receiver name (commonly derived from the struct type). Long or generic receiver names can reduce code readability, especially when distinguishing between the receiver and local variables or types that have similar names.

### Impact

Severity: **low**

This is primarily a readability/naming issue but can make code harder to read and follow, especially as the struct grows or if there are naming clashes with local variables.

### Location

All method receivers for `client`:

```go
func (client *client) GetTenantCapacity(ctx context.Context, tenantId string) (*capacityDto, error) {
    // ...
}
```

### Code Issue

```go
func (client *client) GetTenantCapacity(ctx context.Context, tenantId string) (*capacityDto, error) {
    // ...
}
```

### Fix

Use a concise, idiomatic receiver name (such as `c`):

```go
func (c *client) GetTenantCapacity(ctx context.Context, tenantId string) (*capacityDto, error) {
    // replace all occurrences of 'client.' in the method body with 'c.'
}
```

## ISSUE 2

**File:** `/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/api_copilot_studio_application_insights.go`

### Problem

Within `getCopilotStudioAppInsightsConfiguration`, the environment is fetched twice for the same `environmentId`:

1. In `client.getCopilotStudioEndpoint` (which calls `EnvironmentClient.GetEnvironment`)
2. Directly afterward in the same function

This results in redundant network/API calls.

### Impact

This causes unnecessary network traffic, latency, and increased risk of API throttling with a potential **medium** impact, especially on large or repeated requests.

### Location

```go
copilotStudioEndpoint, err := client.getCopilotStudioEndpoint(ctx, environmentId)
if err != nil {
 return nil, err
}

env, err := client.EnvironmentClient.GetEnvironment(ctx, environmentId)
if err != nil {
 return nil, err
}
```

### Fix

Retrieve the environment once, then use it for both endpoint extraction and property checking.

```go
env, err := client.EnvironmentClient.GetEnvironment(ctx, environmentId)
if err != nil {
 return nil, err
}
copilotStudioEndpoint, err := extractCopilotStudioEndpoint(env)
// ... implement extractCopilotStudioEndpoint to avoid duplicate code.
```

Or, refactor `getCopilotStudioEndpoint` to accept an environment object if already known.

---

## To finish the task you have to

1. Run linter and fix any issues
2. Run UnitTest and fix any of failing ones
3. Generate docs
4. Run Changie

## Changie Instructions

Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```

Where:

- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number

Apply this fix to the whole codebase
