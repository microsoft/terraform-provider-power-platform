# Timeout Attributes Are Not Used in API Calls

##

/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/resource_admin_management_application.go

## Problem

Although the schema defines a `timeouts` attribute (with create, read, and delete timeouts), these timeouts are not passed to or honored by any API call or request context. This means that user-supplied timeouts have no effect on the actual operations the resource performs, which could cause hangs or timeouts to be managed incorrectly by the Terraform core engine.

## Impact

Severity: **low**

While not critical, this impacts expectations: users define timeouts that have no real effect. In long or unreliable operations, this may mislead users.

## Location

```go
"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
    Create: true,
    Delete: true,
    Read:   true,
}),
// ... but not used in Create, Delete, or Read logic
```

## Fix

Consider applying the configured timeout values by enhancing API calls with derived context with timeout, e.g.:

```go
timeout, diags := timeouts.ValueFrom(ctx, plan.Timeouts, CreateTimeout)
if diags.HasError() {
    resp.Diagnostics.Append(diags...)
    return
}
childCtx, cancel := context.WithTimeout(ctx, timeout)
defer cancel()
// Pass childCtx to client call
```

This enables proper resource timeouts based on user configuration.
