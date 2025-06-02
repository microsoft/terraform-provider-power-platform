# No Validation or Checks on ConnectorsClient Construction

##

/workspaces/terraform-provider-power-platform/internal/services/connectors/datasource_connectors.go

## Problem

The `Configure` method assigns `d.ConnectorsClient = newConnectorsClient(client.Api)` but does not validate whether the result is non-nil or if the provided API client is valid. If `newConnectorsClient` can return a nil value or an improperly initialized client (depending on API evolutions or failures), later API calls could panic or fail unclearly.

## Impact

**Medium severity:** Possible panics or unclear errors at later points in execution, especially if any dependency involved in constructing the client fails or changes behavior in the future.

## Location

```go
d.ConnectorsClient = newConnectorsClient(client.Api)
```

## Fix

Validate that `newConnectorsClient` returns a non-nil instance (or proper values), and add a check to append a diagnostic error and return if the client was not constructed properly:

```go
d.ConnectorsClient = newConnectorsClient(client.Api)
if d.ConnectorsClient == nil {
    resp.Diagnostics.AddError(
        "Failed to create connectors client",
        "Connectors client returned nil. Check provider configuration and upstream client logic.",
    )
    return
}
```
