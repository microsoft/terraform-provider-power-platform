# Issue 6: Variable Naming Inconsistency

##

/workspaces/terraform-provider-power-platform/internal/services/application/api_application.go

## Problem

Some receiver names use the full word `client` (`func (client *client)`), which is against Go idiomatic naming—receivers should typically be a single letter (e.g., `c *client`).

## Impact

**Low** severity. Inconsistent receiver naming can reduce code readability and goes against Go best practices.

## Location

Every method defined on `client` struct:

## Code Issue

```go
func (client *client) MethodName(...) {...}
```

## Fix

Refactor all receiver names from `client` to `c`:

```go
func (c *client) MethodName(...) {...}
```
