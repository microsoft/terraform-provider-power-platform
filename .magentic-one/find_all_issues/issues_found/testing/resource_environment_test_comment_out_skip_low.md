# Use of t.Skip with Commented-Out Test Logic

##

/workspaces/terraform-provider-power-platform/internal/services/environment/resource_environment_test.go

## Problem

In `TestAccEnvironmentsResource_Validate_CreateDeveloperEnvironment`, there is a `t.Skip("creating dev environments with SP is NOT yet supported")` but the majority of the test function is left in place, including the test logic and configuration. This makes it unclear if the test is eventually expected to run or if the code/dependencies are obsolete.

## Impact

- **Severity: Low**
- The test will always be skipped, so the logic will never be run. This could mask bitrot (stale config or breaking changes elsewhere), and leaves dead code in the codebase.
- Editors/lint tools may not warn about dead/skip logic if it is wrapped with `t.Skip`, so rot may accumulate unnoticed.

## Location

```go
func TestAccEnvironmentsResource_Validate_CreateDeveloperEnvironment(t *testing.T) {
    t.Skip("creating dev environments with SP is NOT yet supported")
    resource.Test(t, resource.TestCase{
        ...
    })
}
```

## Code Issue

Skipping a test should ideally only wrap the skip logic, or the rest of the logic should be commented/removed until it is needed, with a clear ticket, TODO, or status note.

## Fix

Either:
- Clearly comment that the entire function is deprecated or pending future support with a reference to an issue or planned re-work,
- Or wrap the test contents in an `if false { ... }` or move the full test code behind a flag/conditional so it is not compiled/run until reactivated,
- Or simply remove the function and rely on the VCS/review pipeline to restore it when support is added.

```go
func TestAccEnvironmentsResource_Validate_CreateDeveloperEnvironment(t *testing.T) {
    t.Skip("creating dev environments with SP is NOT yet supported; see ISSUE-1234")
    // Test logic removed for now. See VCS history and tracking ticket for restoration when dev envs are supported.
}
```
