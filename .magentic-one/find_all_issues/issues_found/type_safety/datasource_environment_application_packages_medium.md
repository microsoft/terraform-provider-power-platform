# Issue 4: Potential Data Consistency in Setting State ID

##

Path: /workspaces/terraform-provider-power-platform/internal/services/application/datasource_environment_application_packages.go

## Problem

The line setting `state.Id` uses the length of `applications` (`len(applications)`) which is the count *before* applying the `$filter` logic for Name and PublisherName. This means that if filtering in the for-loop reduces the number of `Applications` placed in the state, the ID may not match the real content.

## Impact

Severity: **Medium**

This could lead to inconsistent IDs, and unnecessary refreshes in Terraform when the actual number of set applications does not match the ID, because filtering reduces the returned set.

## Location

```go
state.Id = types.StringValue(fmt.Sprintf("%s_%d", state.EnvironmentId.ValueString(), len(applications)))
```

## Code Issue

```go
state.Id = types.StringValue(fmt.Sprintf("%s_%d", state.EnvironmentId.ValueString(), len(applications)))
```

## Fix

Use the length of `state.Applications` (the filtered list) instead:

```go
state.Id = types.StringValue(fmt.Sprintf("%s_%d", state.EnvironmentId.ValueString(), len(state.Applications)))
```
