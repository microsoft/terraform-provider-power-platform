# Incorrect Error Handling and Missing Return After Error Addition

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/datasource_securityroles.go

## Problem

If `d.UserClient.DataverseExists` returns an error, the function adds a diagnostic error but does not return immediately, resulting in the following code blocks still being executed. This may result in confusing or misleading additional error messages, as both the error and the "Dataverse does not exist" error could be reported even if the root cause was a client error.

## Impact

Medium severity. Poor error handling flow can result in misleading diagnostics and a harder debugging experience for users and developers.

## Location

```go
	dvExits, err := d.UserClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
	tflog.Debug(ctx, fmt.Sprintf("Environment Id: %s", state.EnvironmentId.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
	}
```

## Code Issue

```go
	dvExits, err := d.UserClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
	tflog.Debug(ctx, fmt.Sprintf("Environment Id: %s", state.EnvironmentId.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
	}
```

## Fix

Return immediately after adding the diagnostic, to prevent further execution if an error occurred:

```go
	dvExits, err := d.UserClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
	tflog.Debug(ctx, fmt.Sprintf("Environment Id: %s", state.EnvironmentId.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
		return
	}
```
