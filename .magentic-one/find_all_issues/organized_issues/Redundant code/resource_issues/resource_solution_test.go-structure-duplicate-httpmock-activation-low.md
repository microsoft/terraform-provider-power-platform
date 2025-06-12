# Title
Redundant Activation and Deactivation of httpmock

##
/workspaces/terraform-provider-power-platform/internal/services/solution/resource_solution_test.go

## Problem
Inside `TestUnitSolutionResource_Validate_Create_And_Force_Recreate`, `httpmock.Activate()` and `defer httpmock.DeactivateAndReset()` are unnecessarily called twice at the beginning of the function. The second pair is redundant.

## Impact
Doesn't break functionality, but introduces code duplication, increases cognitive load, and may confuse readers or future maintainers. Severity: **low** (structure/maintainability).

## Location
Lines (approx.):

```go
func TestUnitSolutionResource_Validate_Create_And_Force_Recreate(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()
    ...
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()
    ...
}
```

## Code Issue
```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()
...
httpmock.Activate()
defer httpmock.DeactivateAndReset()
```

## Fix
Remove the second pair of calls. Only one activation at the beginning and one defer at the end is needed:

```go
func TestUnitSolutionResource_Validate_Create_And_Force_Recreate(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()
    ...
    // Remove duplicate activate/deactivate block
}
```
