# Issue 1: Invalid Error Handling in `Create` Method

##

`/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/resource_enterprise_policy.go`

## Problem

The `Create` method does not differentiate between transient and non-transient errors when interacting with `LinkEnterprisePolicy`. As a result, no retry mechanism is implemented, which can lead to failures related to transient network issues or temporary client-side errors.

## Impact

This could affect the reliability of the `Create` function if called under conditions where transient errors are common (e.g., network interruptions). The lack of retries could cause unnecessary failures for operations that could succeed upon retry.

Severity: **High**

## Location

```go
err := r.EnterprisePolicyClient.LinkEnterprisePolicy(ctx, plan.EnvironmentId.ValueString(), plan.PolicyType.ValueString(), plan.SystemId.ValueString())
if err != nil {
	resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err.Error())
	return
}
```

## Fix

Introduce a retry mechanism for transient errors using an exponential backoff strategy.

```go
import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"time"
)

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// ... existing code

	var retries = 3
	for i := 0; i < retries; i++ {
		err := r.EnterprisePolicyClient.LinkEnterprisePolicy(ctx, plan.EnvironmentId.ValueString(), plan.PolicyType.ValueString(), plan.SystemId.ValueString())
		if err == nil {
			break
		}
		if i < retries-1 {
			tflog.Warn(ctx, fmt.Sprintf("Retrying operation for Client error: %s", err.Error()))
			time.Sleep(time.Duration(2^i) * time.Second) // Exponential backoff
		} else {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err.Error())
			return
		}
	}
	// ... remaining code
}
```