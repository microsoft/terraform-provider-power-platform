# Missing Diagnostic Handling for State Unmarshalling

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/datasource_dlp_policy.go

## Problem

In the `Read` function, after retrieving state:

```go
resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
if resp.Diagnostics.HasError() {
    return
}
```

This error check is immediate, which is correct. However, there is a missing diagnostic log or report in the event of an error—execution simply returns. This may make debugging hard, as nothing is logged, and the user may see a silent failure.

## Impact

When errors arise during state retrieval, there is no indication in logs as to why `Read` has returned early, which damages debuggability for both maintainers and end-users. Severity: **medium**.

## Location

`Read` function, state unmarshalling block.

## Code Issue

```go
resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
if resp.Diagnostics.HasError() {
    return
}
```

## Fix

Log or otherwise inform the user that state retrieval failed:

```go
resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
if resp.Diagnostics.HasError() {
    tflog.Error(ctx, "failed to unmarshal current state in Read: "+resp.Diagnostics.Errors()[0].Summary)
    return
}
```

**Explanation:**  
This provides better traceability in logs/debugging, helping developers quickly identify at what point the function has exited.
