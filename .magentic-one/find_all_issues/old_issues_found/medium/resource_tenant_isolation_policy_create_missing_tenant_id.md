# Title
Error handling for missing tenant ID in the `Create` function is insufficient.

## 
/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy.go

## Problem
The `tenantInfo.TenantId` is retrieved using `r.Client.TenantApi.GetTenant(ctx)`, but there is no validation immediately after the API call to check if `tenantInfo.TenantId` is empty or unset. A missing tenant ID can lead to downstream errors and additional debugging challenges.

## Impact
The absence of a tenant ID may:
1. Cause other function calls like `convertToDto` and `r.Client.createOrUpdateTenantIsolationPolicy` to fail unexpectedly.
2. Introduce undefined or unpredictable behavior later in the process.
3. Inflate the difficulty of tracing errors back to the root cause. Severity: `Medium`

## Location
- File: `/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy.go`
- Function: `func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse)`

## Code Issue
```go
// Get the current tenant ID
tenantInfo, err := r.Client.TenantApi.GetTenant(ctx)
if err != nil {
	resp.Diagnostics.AddError(
		"Error retrieving tenant information",
		fmt.Sprintf("Could not retrieve tenant information: %s", err.Error()),
	)
	return
}
```

## Fix
Add a validation block to confirm that `tenantInfo.TenantId` is not null or empty following its retrieval. This step ensures smooth execution of subsequent logic.

```go
// Get the current tenant ID
tenantInfo, err := r.Client.TenantApi.GetTenant(ctx)
if err != nil {
	resp.Diagnostics.AddError(
		"Error retrieving tenant information",
		fmt.Sprintf("Could not retrieve tenant information: %s", err.Error()),
	)
	return
}

if tenantInfo.TenantId == "" {
	resp.Diagnostics.AddError(
		"Missing Tenant ID",
		"Tenant ID was not returned from the API. Unable to perform the creation of tenant isolation policy.",
	)
	return
}
```

This validation guarantees that the missing `tenantInfo.TenantId` is identified early. 