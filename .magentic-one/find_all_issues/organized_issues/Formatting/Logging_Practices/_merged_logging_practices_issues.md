# Logging Practices Issues

This document contains merged issues related to logging practices in the Power Platform Terraform provider.

## ISSUE 1

**Title:** Possible exposure of sensitive response data in logs

**File:** `/workspaces/terraform-provider-power-platform/internal/services/rest/api_rest.go`

**Problem:**
The code logs the entire response body and status at trace level, which, depending on the API being wrapped, could expose sensitive or confidential information. Even if logs are not enabled in production, this poses a risk if trace logging is accidentally switched on or logs are accessed by unauthorized individuals.

**Impact:**
Severity: High. Logging sensitive data can lead to accidental exposure, data leaks, or violations of data privacy/security policies.

**Location:**
Lines ~53-56:

**Code Issue:**

```go
if res != nil && res.HttpResponse != nil {
 tflog.Trace(ctx, fmt.Sprintf("SendOperation Response: %v", res.BodyAsBytes))
 tflog.Trace(ctx, fmt.Sprintf("SendOperation Response Status: %v", res.HttpResponse.Status))
}
```

**Fix:**
Redact or avoid logging full response bodies. Limit logging to non-sensitive metadata or ensure explicit configuration for safe debug output.

```go
if res != nil && res.HttpResponse != nil {
 tflog.Trace(ctx, fmt.Sprintf("SendOperation Response Status: %v", res.HttpResponse.Status))
 // tflog.Trace(ctx, fmt.Sprintf("SendOperation Response: %v", res.BodyAsBytes)) // Avoid or redact sensitive bodies
}
```

Or, if logging the body is critical, add guards or filters to ensure sensitive fields are removed.

## ISSUE 2

**Title:** Potential Inefficiency: Logging with Sprintf Instead of Key/Value in tflog

**File:** `/workspaces/terraform-provider-power-platform/internal/services/solution/datasource_solutions.go`

**Problem:**
In the `Metadata` function, `tflog.Debug` is called with a formatted string, rather than structured fields. The terraform-plugin-log library encourages structured logging using attributes/fields for better log query, filtering, and searching.

**Impact:**
**Low**. Purely affects logging, but may make logs less queryable and less useful for debugging in complex environments.

**Location:**
Line 38:

**Code Issue:**

```go
tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
```

**Fix:**
Use structured key/value logging:

```go
tflog.Debug(ctx, "METADATA", map[string]any{
 "type_name": resp.TypeName,
})
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
