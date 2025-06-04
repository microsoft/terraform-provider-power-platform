# Issue: Missing Validation for Unknown/Null String Attribute in Schema

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy.go

## Problem

In the resource schema, required attributes such as `"name"`, `"location"`, and those in `"billing_instrument"` allow an empty string by default since there are no explicit `stringvalidator.LengthAtLeast(1)` (or similar) validators. Terraform might send an empty or whitespace string, resulting in failed API calls or unclear errors.

## Impact

Severity: **Medium**

Missing explicit validation could cause confusing errors for users, potential API rejections, and an inconsistent user experience if empty strings are passed to the backend.

## Location

- `Schema` method:  
  ```go
  "name": schema.StringAttribute{
      MarkdownDescription: "The name of the billing policy",
      Required:            true,
  },
  // ... and similar for others
  ```

## Code Issue

```go
"name": schema.StringAttribute{
    MarkdownDescription: "The name of the billing policy",
    Required:            true,
},
```

## Fix

Add `stringvalidator.LengthAtLeast(1)` in the `Validators` for each required string-type attribute to prevent empty string assignments.

```go
"name": schema.StringAttribute{
    MarkdownDescription: "The name of the billing policy",
    Required:            true,
    Validators: []validator.String{
        stringvalidator.LengthAtLeast(1),
    },
},
```

Repeat for `"location"`, `"billing_instrument.resource_group"`, `"billing_instrument.subscription_id"`.
