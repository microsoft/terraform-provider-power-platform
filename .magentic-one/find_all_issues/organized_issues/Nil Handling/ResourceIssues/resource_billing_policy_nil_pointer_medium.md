# Issue: Inconsistent Use of pointer versus non-pointer Resource Model

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy.go

## Problem

Across the CRUD methods, variables representing the Terraform resource model are consistently declared as pointers (e.g. `var plan *BillingPolicyResourceModel`), but there is no explicit check for `nil`. If Terraform or the framework ever passes in `nil`, field dereferencing will panic.

## Impact

Severity: **Medium**

This is a latent bug – if `Get` ever returns a `nil` pointer, the next field access will panic the provider. Even if the current framework never does this, future updates or subtle bugs could introduce this situation.

## Location

- In CRUD methods:
  ```go
  var plan *BillingPolicyResourceModel
  resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
  // no check for plan == nil before dereferencing
  ```

## Code Issue

```go
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
    return
}
billingPolicyToCreate := billingPolicyCreateDto{
    BillingInstrument: BillingInstrumentDto{
        ResourceGroup:  plan.BillingInstrument.ResourceGroup.ValueString(),
        SubscriptionId: plan.BillingInstrument.SubscriptionId.ValueString(),
    },
    Location: plan.Location.ValueString(),
    Name:     plan.Name.ValueString(),
}
```

## Fix

Add an explicit `nil` check after reading into the pointer, and add a diagnostic if `nil` is unexpected:

```go
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
    return
}
if plan == nil {
    resp.Diagnostics.AddError(
        "Internal Error: Missing plan",
        "The plan received by the resource was nil. Please report this bug to the provider maintainers.",
    )
    return
}
```

Repeat for `state` and any other pointers fetched similarly.
