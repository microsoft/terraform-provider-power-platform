# Title

Potential Redundant Error Addition in Case of No Dataverse

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/datasource_environment_settings.go

## Problem

The code adds an error diagnostic with empty details if no Dataverse exists in the environment. This empty error lacks specific context, which could make debugging difficult.

## Impact

Confusing diagnostic messages hinder troubleshooting efforts. Severity is medium, as it exclusively affects error handling but could make debugging harder.

## Location

```go
	if !dvExits {
		resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), "")
		return
	}
```

## Code Issue

Adding an error diagnostic without meaningful details:

```go
	resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), "")
```

## Fix

Provide additional context in the error message.

```go
	resp.Diagnostics.AddError(
		"Dataverse Not Found",
		fmt.Sprintf("No Dataverse exists in the specified environment '%s'. Ensure that the environment ID is correct.", state.EnvironmentId.ValueString()),
	)
```