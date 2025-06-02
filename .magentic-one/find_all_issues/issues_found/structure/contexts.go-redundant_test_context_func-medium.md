# Issue: Redundant Function for Test Context Creation

##

/workspaces/terraform-provider-power-platform/internal/helpers/contexts.go

## Problem

The file contains two functions `UnitTestContext` and `TestContext` which perform the exact same operation: creating a new context with the test context value set. This causes unnecessary redundancy in the codebase and may lead to confusion about which function should be used.

## Impact

- **Impact:** Medium  
  Redundant code can cause maintenance issues and confusion for other developers. It can also increase the cognitive load when debugging or implementing features related to test contexts.

## Location

Function definitions for `UnitTestContext` and `TestContext`:

## Code Issue

```go
func UnitTestContext(ctx context.Context, testName string) context.Context {
	return context.WithValue(ctx, TEST_CONTEXT_KEY, TestContextValue{IsTestMode: true, TestName: testName})
}

// ...

func TestContext(ctx context.Context, testName string) context.Context {
	return context.WithValue(ctx, TEST_CONTEXT_KEY, TestContextValue{IsTestMode: true, TestName: testName})
}
```

## Fix

Remove one of the redundant functions to keep only one function for clarity and simplicity. Update all references in the codebase to use the remaining function.

```go
// Remove UnitTestContext, keep only TestContext as the public-facing API.

func TestContext(ctx context.Context, testName string) context.Context {
	return context.WithValue(ctx, TEST_CONTEXT_KEY, TestContextValue{IsTestMode: true, TestName: testName})
}
```
