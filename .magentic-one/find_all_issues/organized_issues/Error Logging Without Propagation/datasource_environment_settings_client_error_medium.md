# Title

Potential Missing Return After Error in Client Check

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/datasource_environment_settings.go

## Problem

When calling `DataverseExists`, if there is an error, you log a diagnostic but *do not return*. This may result in further logic running on bad data, and possible misleading or cascading errors.

## Impact

Medium. If an error occurs but you continue, the logic may malfunction, and this could confuse end users and make debugging more difficult.

## Location

Lines:

```go
dvExits, err := d.EnvironmentSettingsClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
if err != nil {
	resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
}
```

## Code Issue

```go
dvExits, err := d.EnvironmentSettingsClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
if err != nil {
	resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
}
```

## Fix

Return after adding the diagnostic to avoid further execution.

```go
dvExits, err := d.EnvironmentSettingsClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
if err != nil {
	resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
	return
}
```
