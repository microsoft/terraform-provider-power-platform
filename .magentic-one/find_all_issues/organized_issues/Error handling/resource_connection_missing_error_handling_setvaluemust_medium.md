# Title
Missing error handling for `types.SetValueMust` usages

##
/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection.go

## Problem
In several places, the function `types.SetValueMust` is used when setting the `Status` attribute:

```go
plan.Status = types.SetValueMust(types.StringType, statuses)
```

The `SetValueMust` variant panics if it encounters an error (instead of returning an error), which can cause uncontrolled panics and crash the process. There is no error handling or recovery from such panics in the `Create`, `Read`, and `Update` methods.

## Impact
Medium: If the data is invalid, the provider can panic and crash Terraform, making it less robust and failing ungracefully, potentially losing user data/state.

## Location
Create, Read, Update methods:

## Code Issue
```go
plan.Status = types.SetValueMust(types.StringType, statuses)
// ... and
state.Status = types.SetValueMust(types.StringType, statuses)
```

## Fix
Replace `SetValueMust` with proper error handling using `SetValue` and check its error return. Example:

```go
planStatus, err := types.SetValue(types.StringType, statuses)
if err != nil {
  resp.Diagnostics.AddError("Failed to set status", err.Error())
  return
}
plan.Status = planStatus
```
Do likewise for state assignment in other methods.

