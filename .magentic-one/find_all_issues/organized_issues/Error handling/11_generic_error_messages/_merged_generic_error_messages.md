# Generic Error Messages and Formatting Issues

This document consolidates issues related to generic, inconsistent, or poorly formatted error messages throughout the codebase that need improvement for better user experience and maintainability.

## ISSUE 1

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment_templates/datasource_environment_templates.go`

**Problem:** Inconsistent error message formatting in `Configure` method

The error message in the `Configure` method, within the `AddError` call, is user-directed but loses some helpful information and its structure could be improved for consistency.

Currently:

```go
resp.Diagnostics.AddError(
    "Unexpected ProviderData Type",
    fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
)
```

**Impact:** Severity: Low

- This is a usability and UX issue, but not strictly a correctness problem. Terse and actionable error messages improve developer experience.

**Location:** Method `Configure`, in the block where `ok := req.ProviderData.(*api.ProviderClient)` fails.

**Code Issue:**

```go
resp.Diagnostics.AddError(
    "Unexpected ProviderData Type",
    fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
)
```

**Fix:** Expand the message with instructions, and structure it for clarity, for example:

```go
resp.Diagnostics.AddError(
    "Invalid Provider Configuration",
    fmt.Sprintf("The provider data was not of the expected type '*api.ProviderClient' (got: %T). "+
        "This is likely a bug in the provider. Please file a bug report with the configuration you used.", req.ProviderData),
)
```

## ISSUE 2

**File:** `/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy.go`

**Problem:** API Client Error Message Consistency

The diagnostic error messages for API client operations in CRUD actions are written as plain strings with string concatenation. There is the risk of inconsistent messaging or leaking sensitive error content from the API directly into diagnostics. It does not sanitize or provide meaningful context, especially should the API client error structure change or include nested errors.

**Impact:** Severity: Low

- Can result in unclear, unstructured, or over-verbose error messages for end-users of the provider.
- Might expose sensitive or internal exception text from the underlying API call.

**Location:**

```go
resp.Diagnostics.AddError(
        "Error creating tenant isolation policy",
        fmt.Sprintf("Could not create tenant isolation policy: %s", err.Error()),
)
```

**Code Issue:**

```go
resp.Diagnostics.AddError(
        "Error creating tenant isolation policy",
        fmt.Sprintf("Could not create tenant isolation policy: %s", err.Error()),
)
```

**Fix:** Consider using a centralized helper for error message formatting and sanitization. Always provide user-oriented error context before including technical details (possibly truncated or parsed for relevance). Example:

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

## ISSUE 3

**File:** `/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go`

**Problem:** Use of generic error messages in AddError/AddWarning

Some error and warning diagnostics are generic and repetitive, especially around "Error converting tenant settings", "Error reading tenant", "Error applying corrections", etc. There's potential for more actionable, contextual messages for the end user.

**Impact:** Reduced effectiveness of errors/warnings: less actionable for users and more difficult for provider maintainers to distinguish between error sources during support or debugging. Severity: low.

**Location:** Throughout, e.g.:

- `resp.Diagnostics.AddError("Error converting tenant settings", err.Error())`
- `resp.Diagnostics.AddWarning("Tenant Settings are not deleted", ...)`, etc.

**Code Issue:**

```go
if err != nil {
        resp.Diagnostics.AddError("Error converting tenant settings", err.Error())
        return
}
```

**Fix:** Provide context-rich error titles and messages, e.g.:

```go
if err != nil {
        resp.Diagnostics.AddError(
                "Unable to Convert Tenant Settings in resource_tenant_settings Create",
                fmt.Sprintf("Could not convert planned tenant settings model to DTO: %s", err.Error()),
        )
        return
}
```

Use similar context in AddWarning messages as well.

---

## Task Completion Instructions

After implementing these fixes:

1. **Run the linter:** `make lint` to ensure code style compliance
2. **Run unit tests:** `make unittest` to verify functionality  
3. **Generate documentation:** `make userdocs` to update auto-generated docs
4. **Add changelog entry:** Use `changie new` to document the changes

**Changie Entry Template:**

```yaml
kind: changed
body: Improved error message formatting and consistency for better user experience
time: [current-timestamp]
custom:
  Issue: "[ISSUE_NUMBER_IF_APPLICABLE]"
```

Replace `[ISSUE_NUMBER_IF_APPLICABLE]` with the relevant GitHub issue number, or remove the custom section if no specific issue exists.
