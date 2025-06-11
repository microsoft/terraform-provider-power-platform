# Contexts Client Issues - Merged Issues

## ISSUE 1

# Title

Shadowing the standard library package `net/url` with alias `neturl`

##

/workspaces/terraform-provider-power-platform/internal/api/client.go

## Problem

The import statement `neturl "net/url"` unnecessarily aliases the `net/url` package as `neturl`. This is not idiomatic and may lead to confusion, as all Go code uses `url.Parse()`, not `neturl.Parse()`, unless there is a naming conflict—which is not evident here.

## Impact

Lowers readability, makes it harder for new contributors familiar with Go's standard library. Severity: **low**

## Location

Imports at file top and throughout code wherever `neturl.Parse` is used.

## Code Issue

```go
import (
    //...
    neturl "net/url"
    //...
)

// ...
if u, e := neturl.Parse(url); e != nil || !u.IsAbs() {
    // ...
}
```

## Fix

Remove the alias and use the standard package name:

```go
import (
    //...
    "net/url"
    //...
)

// ...
if u, e := url.Parse(url); e != nil || !u.IsAbs() {
    // ...
}
```


---

## ISSUE 2

# Issue: Unexported Helper Function Naming

##

/workspaces/terraform-provider-power-platform/internal/helpers/contexts.go

## Problem

The function `enterTimeoutContext` is unexported (lowercase "e") but performs a significant action fundamental to handling request contexts and timeouts. In Go, unexported helper functions should follow a clear, precise, and consistent naming convention, and it should be clear that they're only internally used. The name `enterTimeoutContext` is reasonable, but it might be missed as an important utility if not documented or if the naming doesn't follow local idioms/context.

## Impact

- **Impact:** Low  
  Minor readability and consistency issue, but does not introduce bugs or errors.

## Location

Function definition for `enterTimeoutContext`:

## Code Issue

```go
func enterTimeoutContext[T AllowedRequestTypes](ctx context.Context, req T) (context.Context, *context.CancelFunc) {
  // implementation
}
```

## Fix

Ensure clear documentation and that the function is only used internally. Optionally, add a comment or prefix the function name with an underscore (used sometimes in Go for internal helpers, though not strictly necessary).

```go
// enterTimeoutContext creates a derived context with a timeout, based on the request type and configured defaults.
// It must only be used within this package.
func enterTimeoutContext[T AllowedRequestTypes](ctx context.Context, req T) (context.Context, *context.CancelFunc) {
  // implementation
}
```


---

# To finish the task you have to 
1. Run linter and fix any issues 
2. Run UnitTest and fix any of failing ones
3. Generate docs 
4. Run Changie

# Changie Instructions
Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```
Where:
- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
