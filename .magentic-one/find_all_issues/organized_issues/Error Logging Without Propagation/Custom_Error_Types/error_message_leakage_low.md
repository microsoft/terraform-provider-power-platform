# Title

Error String Formatting and Leakage of Internal Details

##

/workspaces/terraform-provider-power-platform/internal/services/solution_checker_rules/client.go

## Problem

Some error messages returned include specific endpoint URLs and internal identifiers (e.g., raw environment IDs, endpoint URLs). While this can aid debugging, it can also unintentionally leak sensitive implementation information if surfaced to end-users or logs, depending on where errors are handled.

## Impact

Potential information leakage; confusion for users unfamiliar with internal endpoints; inconsistent error messages. Severity: low (unless propagated to users, could rise to medium).

## Location

In `GetSolutionCheckerRules`:

```go
if err != nil {
	return nil, fmt.Errorf("failed to get environment details for %s: %w", environmentId, err)
}

if env.Properties.RuntimeEndpoints.PowerAppsAdvisor == "" {
	return nil, fmt.Errorf("could not find PowerAppsAdvisor endpoint for environment %s", environmentId)
}

...

advisorURL, err := url.Parse(env.Properties.RuntimeEndpoints.PowerAppsAdvisor)
if err != nil {
	return nil, fmt.Errorf("failed to parse PowerAppsAdvisor URL: %w", err)
}

...

rulesUrl, err := url.Parse(rulesBaseUrl)
if err != nil {
	return nil, fmt.Errorf("failed to parse rules URL: %w", err)
}
```

## Fix

Sanitize error messages to be user-friendly and avoid exposing internal details unless required for debugging (and ideally, toggle detail level via logging config).

```go
if err != nil {
	return nil, fmt.Errorf("failed to get environment details: %w", err)
}

if env.Properties.RuntimeEndpoints.PowerAppsAdvisor == "" {
	return nil, fmt.Errorf("PowerAppsAdvisor endpoint not found for the target environment")
}

advisorURL, err := url.Parse(env.Properties.RuntimeEndpoints.PowerAppsAdvisor)
if err != nil {
	return nil, fmt.Errorf("invalid advisor endpoint URL: %w", err)
}

rulesUrl, err := url.Parse(rulesBaseUrl)
if err != nil {
	return nil, fmt.Errorf("invalid rules base URL: %w", err)
}
```
