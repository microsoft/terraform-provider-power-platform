# API Response Handling Issues

This document consolidates all issues related to API response handling and error processing in the Terraform Provider for Power Platform.

## ISSUE 1

### Unused HTTP Response Value from API Call

**File:** `/workspaces/terraform-provider-power-platform/internal/services/locations/api_locations.go`

**Problem:** The unused HTTP response value returned by `client.Api.Execute` may leave out valuable information (headers, status code) or error handling opportunities.

**Impact:** Low severity. While not always critical, it is generally better to consider if the returned response may have diagnostic value.

**Location:** `GetLocations` method.

**Code Issue:**

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &locations)
```

**Fix:** If the response value is not needed, name it `_` as is done now, or use documentation to explain, or optionally examine and return it for better caller observability:

```go
resp, err := client.API.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &locations)
// Optionally use resp for additional validation or logging
return locations, err
```

## ISSUE 2

### Inconsistent Error Messages for Invalid API URL

**File:** `/workspaces/terraform-provider-power-platform/internal/services/managed_environment/api_managed_environment.go`

**Problem:** The method `FetchSolutionCheckerRules` returns a generic error `"PowerAppsAdvisor URL is empty"` if the URL field is empty, and wraps errors on URL parsing. However, the distinction between an empty, nil, or invalid URL is not consistent and could be made clearer for maintainability, observability, and debugging.

**Impact:**

- **Low severity**  
- Makes troubleshooting harder as errors are not specific.
- In case of issues, logs/messaging do not provide the exact place of failure or configuration problem.

**Code Issue:**

```go
if env.Properties.RuntimeEndpoints.PowerAppsAdvisor == "" {
    return nil, errors.New("PowerAppsAdvisor URL is empty")
}

powerAppsAdvisorUrl, err := url.Parse(env.Properties.RuntimeEndpoints.PowerAppsAdvisor)
if err != nil {
    return nil, fmt.Errorf("failed to parse PowerAppsAdvisor URL: %w", err)
}
```

**Fix:** Standardize and clarify error handling. Include environmentId or more context:

```go
if env.Properties.RuntimeEndpoints.PowerAppsAdvisor == "" {
    return nil, fmt.Errorf("missing PowerAppsAdvisor endpoint for environment %s", environmentId)
}
powerAppsAdvisorUrl, err := url.Parse(env.Properties.RuntimeEndpoints.PowerAppsAdvisor)
if err != nil {
    return nil, fmt.Errorf("invalid PowerAppsAdvisor URL '%s' for environment %s: %w", env.Properties.RuntimeEndpoints.PowerAppsAdvisor, environmentId, err)
}
```

# To finish the task you have to

1. Run linter and fix any issues
2. Run UnitTest and fix any of failing ones
3. Generate docs
4. Run Changie

# Changie Instructions

Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```

Where:

- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
