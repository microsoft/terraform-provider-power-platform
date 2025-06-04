# Issue: Lack of Proper Error Wrapping and Propagation in getTenantIsolationPolicy

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/api_tenant_isolation_policy.go

## Problem

The function `getTenantIsolationPolicy` directly returns the error it receives from the API call without any context or error wrapping. This can make it hard for callers to determine the source of the error, especially in larger codebases where multiple API calls might propagate generic errors up the stack.

## Impact

Medium. While the error is returned, debugging and tracing the error's origin can become difficult, affecting maintainability and troubleshooting. Developers and support engineers may struggle to identify the failure's context quickly.

## Location

Line in `getTenantIsolationPolicy`:
```go
	if err != nil {
		return nil, err
	}
```

## Code Issue

```go
	if err != nil {
		return nil, err
	}
```

## Fix

Wrap the error with a relevant context to improve traceability:

```go
	if err != nil {
		return nil, fmt.Errorf("could not retrieve tenant isolation policy for tenant %s: %w", tenantId, err)
	}
```
