# Missing error detail/wrapping on GetConnectionShares API error

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connection_shares.go

## Problem

When the GetConnectionShares API call fails, only the error string is given in diagnostics, without parameter context (environment, connector, connection).

## Impact

**Severity: medium**

Makes debugging harder for operators and maintainers.

## Location

```go
if err != nil {
	resp.Diagnostics.AddError("Failed to get connection shares", err.Error())
	return
}
```

## Code Issue

```go
if err != nil {
	resp.Diagnostics.AddError("Failed to get connection shares", err.Error())
	return
}
```

## Fix

Wrap/extend the error with more parameter information:

```go
if err != nil {
	resp.Diagnostics.AddError(
		fmt.Sprintf(
			"Failed to get connection shares for environment_id '%s', connector_name '%s', connection_id '%s'",
			state.EnvironmentId.ValueString(), state.ConnectorName.ValueString(), state.ConnectionId.ValueString(),
		),
		err.Error(),
	)
	return
}
```
