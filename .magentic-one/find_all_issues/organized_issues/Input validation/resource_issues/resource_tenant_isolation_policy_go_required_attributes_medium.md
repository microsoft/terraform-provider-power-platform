# Insufficient Validation of Required Attributes

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy.go

## Problem

While the schema specifies required attributes (`is_disabled`, `allowed_tenants`, etc.), there is no explicit logic in the resource code ensuring these attributes are not empty within Create/Update paths beyond what the framework validator does. There is only validation that at least one of `inbound` or `outbound` is set per tenant, but no guard that `allowed_tenants` is not an empty set (which could represent an always-blocking configuration that may be unintentional).

## Impact

- **Severity: Medium**
- Could allow a situation where a user accidentally submits an empty `allowed_tenants` set (or doesn't provide required directions) and the API call goes through, possibly breaking tenant connectivity in production. Behavioral differences between the Terraform provider's required-vs-empty semantics and actual API expectations can cause state drift or unintended outages.

## Location

```go
// Schema only has Required: true, but in the ValidateConfig and CRUD, no explicit check for allowed_tenants not being empty
```

## Code Issue

```go
"allowed_tenants": schema.SetNestedAttribute{
	Required:            true,
	MarkdownDescription: "List of tenants that are allowed to connect with your tenant.",
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the tenant that is allowed to connect.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"inbound": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether inbound connections from this tenant are allowed.",
			},
			"outbound": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether outbound connections to this tenant are allowed.",
			},
		},
	},
},
// ...
// In ValidateConfig()
if resp.Diagnostics.HasError() {
	return
}
// No check: if len(modelTenants) == 0 { /* error */ }
```

## Fix

Add explicit validation to ensure that the list/set of allowed tenants is not empty, both at the schema validator level and in ValidateConfig, so intent is clear and errors are user-friendly.

```go
// Add validator in schema if supported, otherwise in ValidateConfig:
if len(modelTenants) == 0 {
	resp.Diagnostics.AddAttributeError(
		path.Root("allowed_tenants"),
		"Empty allowed_tenants set",
		"'allowed_tenants' must not be empty. At least one outbound or inbound allowed tenant must be specified.",
	)
	return
}
```
