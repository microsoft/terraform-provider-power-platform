# API Client Error Message Consistency

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy.go

## Problem

The diagnostic error messages for API client operations in CRUD actions are written as plain strings with string concatenation. There is the risk of inconsistent messaging or leaking sensitive error content from the API directly into diagnostics. It does not sanitize or provide meaningful context, especially should the API client error structure change or include nested errors.

## Impact

- **Severity: Low**
- Can result in unclear, unstructured, or over-verbose error messages for end-users of the provider.
- Might expose sensitive or internal exception text from the underlying API call.

## Location

```go
resp.Diagnostics.AddError(
	"Error creating tenant isolation policy",
	fmt.Sprintf("Could not create tenant isolation policy: %s", err.Error()),
)
```

## Code Issue

```go
resp.Diagnostics.AddError(
	"Error creating tenant isolation policy",
	fmt.Sprintf("Could not create tenant isolation policy: %s", err.Error()),
)
```

## Fix

Consider using a centralized helper for error message formatting and sanitization. Always provide user-oriented error context before including technical details (possibly truncated or parsed for relevance). Example:

```go
func humanizeApiError(prefix string, err error) string {
	// Truncate or format as needed to avoid leaking overly technical or sensitive info.
	return fmt.Sprintf("%s: %s", prefix, err.Error())
}

// Usage:
resp.Diagnostics.AddError(
	"Error creating tenant isolation policy",
	humanizeApiError("Could not create tenant isolation policy", err),
)
```
This approach standardizes error output, improves maintainability, and allows central control of sensitive details.
