# Title

Delete uses `state.Id` instead of `state.Principal.EntraObjectId` for principal identification

##

/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection_share.go

## Problem

In the `Delete` method:

```go
err := r.ConnectionsClient.DeleteConnectionShare(ctx, state.EnvironmentId.ValueString(), state.ConnectorName.ValueString(), state.ConnectionId.ValueString(), state.Id.ValueString())
```

The last parameter is `state.Id.ValueString()`. Throughout the create/read/update cycle, the code always passes `plan.Principal.EntraObjectId.ValueString()` or similar for principal identification. Here, however, `state.Id` (likely the connection share's unique identifier) is used instead of the principal's ID. This is inconsistent and may cause API errors if the backend expects the principal's object ID.

## Impact

Severity is **high**, as this may lead to failing deletes or deleting the wrong share, causing data inconsistencies.

## Location

Delete method:

```go
err := r.ConnectionsClient.DeleteConnectionShare(ctx, state.EnvironmentId.ValueString(), state.ConnectorName.ValueString(), state.ConnectionId.ValueString(), state.Id.ValueString())
```

## Code Issue

```go
err := r.ConnectionsClient.DeleteConnectionShare(ctx, state.EnvironmentId.ValueString(), state.ConnectorName.ValueString(), state.ConnectionId.ValueString(), state.Id.ValueString())
```

## Fix

Use the principal's EntraObjectId as in other lifecycle methods:

```go
err := r.ConnectionsClient.DeleteConnectionShare(ctx, state.EnvironmentId.ValueString(), state.ConnectorName.ValueString(), state.ConnectionId.ValueString(), state.Principal.EntraObjectId.ValueString())
```
