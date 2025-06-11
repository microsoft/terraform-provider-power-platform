# Issue: Panic Risk if LicensingClient is Nil

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy.go

## Problem

Throughout the CRUD functions (`Create`, `Read`, `Update`, `Delete`) the `LicensingClient` is called without checking if it is initialized (i.e., nil). If `.Configure` was never called with valid `ProviderData`, `LicensingClient` would be nil, leading to a runtime panic.

## Impact

Severity: **High**

If this situation happens, the provider will panic, causing Terraform runs to crash, user frustration, and possibly lost state. The failure is abrupt and non-recoverable.

## Location

- All CRUD methods (`Create`, `Read`, `Update`, `Delete`)
- Usage:  
  ```go
  policy, err := r.LicensingClient.CreateBillingPolicy(ctx, billingPolicyToCreate)
  ```  
  ...and similar

## Code Issue

```go
policy, err := r.LicensingClient.CreateBillingPolicy(ctx, billingPolicyToCreate)
// and similar usage elsewhere
```

## Fix

Add a check for nil `LicensingClient` at the beginning of each method that uses it, and produce a meaningful error if it isn’t initialized.

```go
if r.LicensingClient == nil {
	resp.Diagnostics.AddError("Uninitialized LicensingClient", "Could not access LicensingClient; the provider may not be configured properly. Please review your provider configuration.")
	return
}
```

Add this check to the start of each CRUD method and any other using `LicensingClient`, just after `defer exitContext()`.
