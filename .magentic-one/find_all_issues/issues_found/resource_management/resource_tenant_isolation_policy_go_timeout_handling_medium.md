# Potential Timeout Handling Omission

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy.go

## Problem

The resource schema defines a `timeouts` attribute using the Terraform plugin framework, but the Create/Read/Update/Delete functions do not refer to or make use of the timeout values when calling the underlying API. This could result in provider operations not respecting user- or operator-defined timeouts, which might in turn cause operations to hang indefinitely if the remote API misbehaves.

## Impact

- **Severity: Medium**
- User-defined operation timeouts are ignored, possibly leading to hanging/deadlocked resource changes or poor user experience for long-running operations.

## Location

```go
"timeouts": timeouts.Attributes(ctx, timeouts.Opts{}),
// ...
// In Create/Read/Update/Delete: state.Timeouts and plan.Timeouts are carried, but not used to set context deadlines or passed to client calls.
```

## Code Issue

```go
// Example from Create:
var plan TenantIsolationPolicyResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
	return
}
...
tenantInfo, err := r.Client.TenantApi.GetTenant(ctx) // No timeout used here
...
policy, err := r.Client.createOrUpdateTenantIsolationPolicy(ctx, tenantInfo.TenantId, *policyDto) // No timeout used here
```

## Fix

Extract the relevant timeout value from the plan/state, and use it to set a context with deadline for the duration of the API call. Pass this context to the API operations. This ensures user timeouts are respected.

```go
import "time"

// Example timeout handling (pseudo-code, as 'timeouts' model may vary):
timeout := plan.Timeouts.Create // Default, or read the correct timeout per op
apiCtx, cancel := context.WithTimeout(ctx, timeout)
defer cancel()
// Use apiCtx in all API calls
tenantInfo, err := r.Client.TenantApi.GetTenant(apiCtx)
if err != nil { ... }
policy, err := r.Client.createOrUpdateTenantIsolationPolicy(apiCtx, tenantInfo.TenantId, *policyDto)
if err != nil { ... }
```
