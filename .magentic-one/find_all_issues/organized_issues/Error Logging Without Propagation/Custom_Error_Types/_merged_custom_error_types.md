# Error Logging Without Propagation - Custom Error Types

This document consolidates all issues related to error logging without proper propagation found in custom error type implementations across the Terraform Provider for Power Platform.

## ISSUE 1

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

## ISSUE 2

# Title

Error Message May Leak Underlying Error Detail

##

/workspaces/terraform-provider-power-platform/internal/customerrors/url_format_error.go

## Problem

The error message formatting in the `Error()` method directly includes the underlying error message returned from `e.Err.Error()`. This can potentially leak internal or sensitive error details to the caller or logs, which might not be appropriate for end-users and could pose a security concern.

## Impact

In cases where `e.Err` contains internal context (such as stack traces, credentials, or sensitive configuration), this could inadvertently expose sensitive information outside of expected logging channels. The severity is **medium** because this can have consequences in production environments and public logs.

## Location

Method: `func (e UrlFormatError) Error() string`
File: /workspaces/terraform-provider-power-platform/internal/customerrors/url_format_error.go

## Code Issue

```go
func (e UrlFormatError) Error() string {
 errorMsg := ""
 if e.Err != nil {
  errorMsg = e.Err.Error()
 }

 return fmt.Sprintf("Request url must be an absolute url: '%s' : '%s'", e.Url, errorMsg)
}
```

## Fix

Carefully sanitize or wrap error information. Only propagate user-friendly or non-sensitive messages in user-facing error strings. If additional debugging is needed, use logging (not error messages) for sensitive details.

```go
func (e UrlFormatError) Error() string {
 if e.Err != nil {
  return fmt.Sprintf("Request URL must be an absolute URL: '%s' : error occurred during URL parsing/validation.", e.Url)
 }
 return fmt.Sprintf("Request URL must be an absolute URL: '%s'", e.Url)
}
```

If internal error details are useful, expose them only in logs under a debug mode, not in the error string returned.

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
