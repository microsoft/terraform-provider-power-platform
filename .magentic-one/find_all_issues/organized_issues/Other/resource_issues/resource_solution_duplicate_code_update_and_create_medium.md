# Code Duplication: Nearly Identical Logic in Create and Update Functions

##
/workspaces/terraform-provider-power-platform/internal/services/solution/resource_solution.go

## Problem
The logic in `Create` and `Update` functions is almost identical regarding how they:
- Call `importSolution`
- Update `plan` fields based on the result
- Re-calculate and set the checksums

This code duplication increases maintenance burden, as future changes must be synchronized between both locations. It also can be a source of inconsistencies and bugs if additional logic is ever added to just one function (by mistake), causing divergence. A shared helper would increase maintainability.

## Impact
- **Severity:** Medium
- Higher ongoing maintenance effort
- Greater risk of subtle bugs creeping in due to unsynchronized refactors
- Reduced clarity/readability for future contributors

## Location
The code bodies of `Create` and `Update` functions.

## Code Issue
```go
// Both functions have nearly identical logic for checksum calculation and field update
```

## Fix
Extract the common logic into a helper/private method, for example `populatePlanFieldsAfterImport()` and call it from both `Create` and `Update`.

```go
func (r *Resource) populatePlanFieldsAfterImport(ctx context.Context, plan *ResourceModel, solution *SolutionDto, resp diagnosticsAdder) {
    // Implement shared field and checksum logic here
}
```
Then, in `Create` and `Update`, replace duplicated segments with a call to this helper.
