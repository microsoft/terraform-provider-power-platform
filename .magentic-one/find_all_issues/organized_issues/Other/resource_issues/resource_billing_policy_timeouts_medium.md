# Issue: Incorrect or Incomplete Timeout Management

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy.go

## Problem

Resource-level timeouts are defined in the schema via `timeouts.Attributes`, but the resource methods (`Create`, `Read`, `Update`, `Delete`) do not retrieve or use the user-configured timeouts from the resource data or context. Terraform may pass extended or custom timeouts, but these are ignored, potentially causing premature context cancellations (default 60m may not be honored).

## Impact

Severity: **Medium**

Custom timeout values given by users in their configuration will be ignored, possibly leading to timeouts on long-running operations, or overruling user intent. This could cause resource creation or deletion to fail even though the API is still processing.

## Location

- `Schema`:  
  ```go
  "timeouts": timeouts.Attributes(ctx, timeouts.Opts{
      Create: true,
      Update: true,
      Delete: true,
      Read:   true,
  }),
  ```
- `Create`/`Read`/`Update`/`Delete`:  
  Not using or extracting timeouts from the request/input.

## Code Issue

```go
// Example from Create, but same in all methods
var plan *BillingPolicyResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
// NO logic for extracting timeout or using timeout context
```

## Fix

Retrieve the timeout duration (if present) and prepare the child context accordingly:

```go
import (
    "time"
    "github.com/hashicorp/terraform-plugin-framework/resource/timeouts"
)

...

timeout, diags := timeouts.Get(ctx, req.Plan, "create") // or req.State, or the method name ("update", "delete", "read")
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
    return
}
if timeout != nil {
    var cancel context.CancelFunc
    ctx, cancel = context.WithTimeout(ctx, *timeout)
    defer cancel()
}
```

Incorporate this pattern for each method, using the correct timeout for the verb.
