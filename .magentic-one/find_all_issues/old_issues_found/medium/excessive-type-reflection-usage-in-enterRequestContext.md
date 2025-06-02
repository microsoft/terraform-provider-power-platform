# Title

Excessive Type Reflection Usage in `EnterRequestContext`

##

`/workspaces/terraform-provider-power-platform/internal/helpers/contexts.go`

## Problem

Type reflection is heavily used (`reflect.TypeOf(req).String()` and `typ.FullTypeName()`), which adds cognitive overhead and has performance implications.

## Impact

Using reflection-based type operations causes performance degradation and can make the code less maintainable and harder to debug. While Go's reflection capabilities can be powerful, they should be used judiciously. Severity: **Medium**.

## Location

Function `EnterRequestContext`, located in `/workspaces/terraform-provider-power-platform/internal/helpers/contexts.go`.

## Code Issue

```go
reqType := reflect.TypeOf(req).String()
name := typ.FullTypeName()
```

## Fix

Consider alternatives like type assertions or passing explicit type data to reduce reliance on runtime reflection.

**Example Fix:**
```go
reqType := typ.GetRequestType() // Define explicitly instead of relying on reflection
name := typ.GetFullName()
```