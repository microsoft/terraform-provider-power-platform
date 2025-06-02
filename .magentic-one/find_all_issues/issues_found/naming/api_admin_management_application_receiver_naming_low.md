# Title

Method Receiver Naming: Collides With Type Name

##

/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/api_admin_management_application.go

## Problem

The receiver for methods on `client` is named `client`, i.e. `func (client *client)`. Idiomatic Go prefers short, distinctive receiver names (e.g., `c` for `client`) to avoid confusion and unnecessary typing.

## Impact

Low. Readability and maintainability: makes it easier to reason about code, especially in larger files where variable naming could overlap.

## Location

All `client` method declarations.

## Code Issue

```go
func (client *client) GetAdminApplication...
func (client *client) RegisterAdminApplication...
func (client *client) UnregisterAdminApplication...
```

## Fix

Change receiver name to a short, idiomatic name:

```go
func (c *client) GetAdminApplication...
func (c *client) RegisterAdminApplication...
func (c *client) UnregisterAdminApplication...
```
