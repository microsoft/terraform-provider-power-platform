# Title

Missing Error Check after convertSourceModelFromEnvironmentDto in Create Method

## Path

/workspaces/terraform-provider-power-platform/internal/services/environment/resource_environment.go

## Problem

In the `Create` method, after calling `convertSourceModelFromEnvironmentDto`, the `err` returned is not checked before proceeding further in the code. This can lead to undefined behavior if the conversion fails due to invalid data, leaving `newState` in an inconsistent state.

## Impact

- Impacts the reliability of the resource creation process.
- May lead to runtime errors or incorrect data being written to the state.
- Severity: **Critical**, as this directly affects the integrity of newly created resources.

## Location

File: `/workspaces/terraform-provider-power-platform/internal/services/environment/resource_environment.go`

Code Issue:

```go
newState, err := convertSourceModelFromEnvironmentDto(*envDto, &currencyCode, plan.OwnerId.ValueStringPointer(), templateMetadata, templates, plan.Timeouts, *r.EnvironmentClient.Api.Config)

if !plan.AzureRegion.IsNull() && plan.AzureRegion.ValueString() != "" && (plan.AzureRegion.ValueString() != newState.AzureRegion.ValueString()) {
    resp.Diagnostics.AddAttributeError(path.Root("azure_region"), fmt.Sprintf("Provisioning environment in azure region '%s' failed", plan.AzureRegion.ValueString()), "Provisioning environment in azure region was not successful, please try other region in that location or try again later")
    return
}
```

## Fix

Add an error check immediately after calling `convertSourceModelFromEnvironmentDto` to verify if the function executed successfully.

```go
newState, err := convertSourceModelFromEnvironmentDto(*envDto, &currencyCode, plan.OwnerId.ValueStringPointer(), templateMetadata, templates, plan.Timeouts, *r.EnvironmentClient.Api.Config)
if err != nil {
    resp.Diagnostics.AddError("Error when converting environment to source model", err.Error())
    return
}

if !plan.AzureRegion.IsNull() && plan.AzureRegion.ValueString() != "" && (plan.AzureRegion.ValueString() != newState.AzureRegion.ValueString()) {
    resp.Diagnostics.AddAttributeError(path.Root("azure_region"), fmt.Sprintf("Provisioning environment in azure region '%s' failed", plan.AzureRegion.ValueString()), "Provisioning environment in azure region was not successful, please try other region in that location or try again later")
    return
}
```