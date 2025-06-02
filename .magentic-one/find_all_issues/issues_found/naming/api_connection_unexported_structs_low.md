# Title

Unexported Structs and Methods Reduce Testability and Reuse

##

/workspaces/terraform-provider-power-platform/internal/services/connection/api_connection.go

## Problem

The `client` struct and its methods are unexported (start with a lowercase letter), which might limit testability and code reuse outside this package, especially in integration tests or external helpers.

## Impact

If these types and methods are meant to be reused or tested externally, their unexported status would block those use-cases. This may complicate testing and composition. Severity: Low.

## Location

For example:

```go
type client struct {
	Api *api.Client
}
```
and
```go
func (client *client) CreateConnection(...)
```

## Code Issue

```go
type client struct {
	Api *api.Client
}
```

## Fix

If these are intended to be public (for wider use), capitalize the struct and method names:

```go
type Client struct {
	Api *api.Client
}
```
and method receivers:
```go
func (c *Client) CreateConnection(...)
```

---

*This file should be saved in:*
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/api_connection_unexported_structs_low.md
