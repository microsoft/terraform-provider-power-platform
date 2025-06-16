# Issue: Public-Like Method Names for Unexported Methods

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/api_tenant_isolation_policy.go

## Problem

The code uses methods such as `getTenantIsolationPolicy`, `createOrUpdateTenantIsolationPolicy`, and `doWaitForLifecycleOperationStatus` that are unexported (start with a lowercase letter), yet are named in a style (starting with a verb without an underscore or indication of being private) that is often reserved for exported/public methods in idiomatic Go. While not a compiler issue, Go best practices encourage more clear distinction for unexported helpers (sometimes starting with an underscore or less user-facing function prefixes).

## Impact

Low. Does not affect execution but can create confusion about method visibility and intended API for new collaborators, and slightly reduces readability.

## Location

Method declarations throughout the file:

```go
func (client *Client) getTenantIsolationPolicy(...)
func (client *Client) createOrUpdateTenantIsolationPolicy(...)
func (client *Client) doWaitForLifecycleOperationStatus(...)
```

## Code Issue

```go
func (client *Client) getTenantIsolationPolicy(...)
func (client *Client) createOrUpdateTenantIsolationPolicy(...)
func (client *Client) doWaitForLifecycleOperationStatus(...)
```

## Fix

For clarity, either:

- Prefix with an underscore to make intent obvious that these are internal (e.g. `_getTenantIsolationPolicy`).  
- Or ensure public/exported methods use capital letter and internal helpers are private and better documented for intent.

Best Go practice is not to prefix with underscores, but to restrict their visibility with lower-case and clearly document their usage in code comments. Add clear doc comments and keep as is, unless there is a need to make these public.

---
