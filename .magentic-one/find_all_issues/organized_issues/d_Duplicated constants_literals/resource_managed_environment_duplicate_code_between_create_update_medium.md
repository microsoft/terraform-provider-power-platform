# Title

Duplicate logic in Create and Update methods

##

/workspaces/terraform-provider-power-platform/internal/services/managed_environment/resource_managed_environment.go

## Problem

The Create and Update methods share a very large block of nearly-identical logic: they both fetch valid solution checker rules, validate user overrides, construct a DTO, call EnableManagedEnvironment, and then reconstruct state from the returned environment. This violates DRY (Don't Repeat Yourself) and reduces maintainability—any fix/enhancement to this logic must be duplicated, increasing the risk of subtle bugs or divergence.

## Impact

Medium. Maintenance burden increases and bug fixes can be accidentally applied in only one place, causing inconsistent behavior between resource creation and update.

## Location

Both Create and Update, from solution checker rules validation through to setting state from the API response.

## Code Issue

Example of duplicated logic:
```go
// Fetch the available solution checker rules
validRules, err := r.ManagedEnvironmentClient.FetchSolutionCheckerRules(ctx, plan.EnvironmentId.ValueString())
if err != nil {
    resp.Diagnostics.AddError("Failed to fetch solution checker rules", err.Error())
    return
}
// ...
// Validate the provided solutionCheckerRuleOverrides
// ...
// Construct DTO
// ...
// Set state from env
```

## Fix

Extract the shared logic to helper functions, e.g.:

```go
type SolutionCheckerResult struct {
    Dto environment.GovernanceConfigurationDto
    RuleOverrides *string
}

func (r *ManagedEnvironmentResource) validateAndBuildGovernanceDTO(ctx context.Context, plan *ManagedEnvironmentResourceModel) (*SolutionCheckerResult, diag.Diagnostics) {
    // move shared code here
}
```
Call the helper from both Create and Update. This ensures single-responsibility and centralization for bug fixing or future code expansion.
