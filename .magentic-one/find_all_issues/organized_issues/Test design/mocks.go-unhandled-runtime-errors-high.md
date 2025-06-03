# Issue: Unhandled error return values on runtime.Caller and runtime.FuncForPC

##

/workspaces/terraform-provider-power-platform/internal/mocks/mocks.go

## Problem

In the `TestName` function, the `runtime.Caller` function returns four values, but the error value is ignored (the last, `ok`). Similarly, `runtime.FuncForPC(pc)` can return `nil`, but its result is used directly. This risks panics if `FuncForPC` returns `nil`.

## Impact

**Severity: High**

If the `runtime.Caller` call fails or the returned `pc` doesn't reference a valid function, `FuncForPC(pc)` may return `nil`. Subsequently, calling `.Name()` on a nil pointer will cause a panic, resulting in abrupt test failures that are confusing and may be hard to debug.

## Location

Line: inside `func TestName() string { ... }`

## Code Issue

```go
func TestName() string {
	pc, _, _, _ := runtime.Caller(1)
	nameFull := runtime.FuncForPC(pc).Name()
	nameEnd := filepath.Ext(nameFull)
	name := strings.TrimPrefix(nameEnd, ".")
	return name
}
```

## Fix

Check `ok` from `runtime.Caller`, and check for `nil` from `runtime.FuncForPC(pc)`:

```go
func TestName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}
	nameFull := fn.Name()
	nameEnd := filepath.Ext(nameFull)
	name := strings.TrimPrefix(nameEnd, ".")
	return name
}
```

This prevents potential panics and clearly handles situations where function name information is unavailable.

---

**Save this file as:**  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/mocks.go-unhandled-runtime-errors-high.md`
