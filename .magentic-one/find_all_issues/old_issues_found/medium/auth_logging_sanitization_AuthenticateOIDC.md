# Title

Inappropriate placement of logging in the `AuthenticateOIDC` method.

##

/workspaces/terraform-provider-power-platform/internal/api/auth.go

## Problem

The error message logged using `tflog.Error` in the `AuthenticateOIDC` method seems to be more of an error handling detail rather than a meaningful log for debugging purposes. Logging this directly might reveal sensitive information if properly not filtered.

Example problematic log:
```go
tflog.Error(ctx, fmt.Sprintf("newDefaultAzureCredential failed to initialize oidc credential:\n\t%s", err.Error()))
```

## Impact

Including sensitive or irrelevant technical data in production logs can lead to sensitive data exposure and make logs harder to consume for operational or security personnel. **Severity: Medium**.

## Location

The `AuthenticateOIDC` method at line 318 in `/workspaces/terraform-provider-power-platform/internal/api/auth.go`.

## Code Issue

```go
tflog.Error(ctx, fmt.Sprintf("newDefaultAzureCredential failed to initialize oidc credential:\n\t%s", err.Error()))
```

## Fix

Improve logging by logging a meaningful, high-level message without exposing sensitive error details unnecessarily.

```go
tflog.Error(ctx, "Failed to initialize OIDC credential", map[string]interface{}{
    "error": "Initialization error - details omitted",
})
```

Explanation:
- Instead of directly logging `err.Error()`, a sanitized message is logged.
- By using a higher-level error message, sensitive details are avoided and logs remain concise yet helpful.
