# Lack of Parallel Test Execution Restriction Is Not Explicit

##

/workspaces/terraform-provider-power-platform/internal/helpers/config_test.go

## Problem

Environment variables are process global. You rightly avoid `t.Parallel()` and state so via comments. However, Go test runner may run other files concurrently. For truly safe isolation when tests in other files/packages might also manipulate environment variables, you should document this clearly in a package-level comment and/or synchronize such tests, particularly if CI reliability is a concern.

## Impact

Severity: **Medium**  
If other tests run in parallel and use the same environment variables, this may result in rare flakiness, especially in a larger codebase or if parallel package test execution is enabled.

## Location

Anywhere environment variables are set and read in tests.

## Code Issue

```go
// No code visible -- this is a missing defense/documentation issue regarding parallelism.
```

## Fix

- Add a package-level comment or test file doc stating:  
  `// Note: tests in this file set environment variables and must not run concurrently with other such tests.`

- Consider renaming environment variables to be more unique per test if/when possible.
