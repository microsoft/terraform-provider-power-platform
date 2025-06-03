# Missing cancellation defer in tests

##

/workspaces/terraform-provider-power-platform/internal/api/client_test.go

## Problem

In the tests where a cancellable context is created (using `context.WithTimeout`), the `cancel()` function is sometimes called at the end of the test, but not immediately after its creation with a `defer`. If any failure or panic occurs before `cancel()` is reached, the cancellation may be skipped, possibly leading to resource leaks or confusion in future tests.

## Impact

This is a low severity issue as the Go runtime will eventually reclaim resources and the tests are relatively short-lived, but best practice is to ensure cancellation is always called by deferring it immediately after creation, ensuring proper resource cleanup in all code paths including failure early returns.

## Location

Several places, for example:

```go
ctx := context.Background()
ctx, cancel := context.WithTimeout(ctx, time.Duration(1)*time.Second)
err := a.SleepWithContext(ctx, time.Duration(5)*time.Second)
...
cancel()
```

## Code Issue

```go
ctx := context.Background()
ctx, cancel := context.WithTimeout(ctx, time.Duration(1)*time.Second)
err := a.SleepWithContext(ctx, time.Duration(5)*time.Second)
if err == nil {
	t.Error("Expected an error but got nil error")
}

if err.Error() != "context deadline exceeded" {
	t.Errorf("Expected error message %s but got %s", "context deadline exceeded", err.Error())
}

cancel()
```

## Fix

Defer the cancellation immediately after creating the context with timeout, for example:

```go
ctx := context.Background()
ctx, cancel := context.WithTimeout(ctx, time.Duration(1)*time.Second)
defer cancel()
err := a.SleepWithContext(ctx, time.Duration(5)*time.Second)
if err == nil {
	t.Error("Expected an error but got nil error")
}

if err.Error() != "context deadline exceeded" {
	t.Errorf("Expected error message %s but got %s", "context deadline exceeded", err.Error())
}
```
