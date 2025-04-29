# Title

Incomplete Validation in Create Function for Dataverse Existence

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/resources_environment_settings.go

## Problem

The `Create` function does not handle cases where the Dataverse check fails due to unexpected errors. While it checks for Dataverse existence, it does not provide sufficient context or fallback mechanisms in cases where `DataverseExists` encounters unexpected behavior.

## Impact

Critical failure in handling Dataverse validation may lead to resource creation attempts in invalid environments, potentially causing security risks or integrity issues. Users may also face confusing error diagnostics without sufficient context.

Severity: **Critical**

## Location

Line ~140: Function `Create`.

## Code Issue

```go
    dvExits, err := r.EnvironmentSettingClient.DataverseExists(ctx, plan.EnvironmentId.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), err.Error())
        return
    }

    if !dvExits {
        resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), "")
        return
    }
```

## Fix

Introduce robust error handling for the Dataverse existence check, and provide fallback mechanisms to gracefully handle unexpected behavior.

```go
    dvExits, err := r.EnvironmentSettingClient.DataverseExists(ctx, plan.EnvironmentId.ValueString())
    if err != nil {
        resp.Diagnostics.AddError("Client Error: Dataverse Check Failed", fmt.Sprintf("Client error occurred when verifying Dataverse existence in environment '%s': %s", plan.EnvironmentId.ValueString(), err.Error()))
        return
    }

    if !dvExits {
        resp.Diagnostics.AddError("Validation Error: Dataverse Missing", fmt.Sprintf("Dataverse does not exist in the specified environment: '%s'. Ensure the environment ID is correct or create a Dataverse instance.", plan.EnvironmentId.ValueString()))
        return
    }
```
