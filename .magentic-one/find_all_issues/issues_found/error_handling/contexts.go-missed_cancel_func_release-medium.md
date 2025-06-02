# Issue: Control Flow—Missed CancelFunc Release on Early Return

##

/workspaces/terraform-provider-power-platform/internal/helpers/contexts.go

## Problem

The returned closure function from `EnterRequestContext` always calls `(*cancel)()` if `cancel` is not nil. However, there are code paths within `enterTimeoutContext` where `cancel` will be `nil` (for example, if an error occurs). While `(*cancel)()` is only called if `cancel != nil`, this places a subtle responsibility on the caller to always check for nil, which is risky in future refactoring.

Additionally, it's idiomatically clearer in Go to return a no-op function if cancellation isn't needed, and to never return a potentially nil function pointer.

## Impact

- **Impact:** Medium  
  If someone refactors or copies this idiom incorrectly, it could lead to code that unintentionally dereferences nil, or misses the opportunity to clean up resources correctly.

## Location

Closure function in `EnterRequestContext`'s return statement:

## Code Issue

```go
return ctx, func() {
	tflog.Debug(ctx, fmt.Sprintf("%s END: %s", reqType, name))
	if cancel != nil {
		(*cancel)()
	}
}
```

## Fix

Change `enterTimeoutContext` to return a no-op cancel function if there was no cancellation set, and always call it.

```go
func enterTimeoutContext[T AllowedRequestTypes](ctx context.Context, req T) (context.Context, context.CancelFunc) {
	// instead of *context.CancelFunc return type, use just context.CancelFunc and return context.CancelFunc(func(){}) when nil
}
```

Change `EnterRequestContext` closure to always call the returned cancel function.

```go
ctx, cancel := enterTimeoutContext(ctx, req)

return ctx, func() {
	tflog.Debug(ctx, fmt.Sprintf("%s END: %s", reqType, name))
	cancel()
}
```
