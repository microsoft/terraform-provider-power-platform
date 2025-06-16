# Title

Confusing or Inaccurate Variable Naming: dvExits

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/datasource_environment_settings.go

## Problem

The variable `dvExits` is likely a typo—it should probably be `dvExists`. This affects readability.

## Impact

Low. Naming clarity impacts maintainability, but the logic still works as expected.

## Location

All of:

```go
dvExits, err := d.EnvironmentSettingsClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
if err != nil {
	resp.Diagnostics.AddError(...)
}
if !dvExits {
	resp.Diagnostics.AddError(...)
	return
}
```

## Code Issue

```go
dvExits, err := d.EnvironmentSettingsClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
if err != nil {
	resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
}
if !dvExits {
	resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), "")
	return
}
```

## Fix

Rename variable `dvExits` to `dvExists`.

```go
dvExists, err := d.EnvironmentSettingsClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
if err != nil {
	resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
	return
}
if !dvExists {
	resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), "")
	return
}
```
