# No Nil Check Before Using DlpPolicyClient

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/datasource_dlp_policy.go

## Problem

In the `Read` method:

```go
policies, err := d.DlpPolicyClient.GetPolicies(ctx)
```

The code assumes that `d.DlpPolicyClient` will always be properly initialized in `Configure`. If for any reason `Configure` fails or is not called, `DlpPolicyClient` will be nil and calling a method on a nil pointer will panic at runtime.

## Impact

This can cause Terraform operations to panic, which is critical since it can halt all resource operations and result in primary provider failure. Severity: **critical**.

## Location

`Read` function, just before calling `GetPolicies(ctx)`.

## Code Issue

```go
policies, err := d.DlpPolicyClient.GetPolicies(ctx)
```

## Fix

Add a nil check and provide an appropriate error diagnostic message if the client is not configured:

```go
if d.DlpPolicyClient == nil {
    resp.Diagnostics.AddError(
        "Client not initialized",
        "The DLP Policy client was not configured. Ensure the provider Configure method has run correctly.",
    )
    return
}
policies, err := d.DlpPolicyClient.GetPolicies(ctx)
```

**Explanation:**  
Protects against panics and gives clear diagnostic output to the user, helping them quickly detect provider misconfiguration.
