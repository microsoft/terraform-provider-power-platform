# Title

Incorrect usage of required/optional attribute and validation in sharing_controls

##

/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/resource_environment_group_rule_set.go

## Problem

The `share_max_limit` attribute in `sharing_controls` is marked as both `Optional` and `Computed`. In Terraform Plugin Framework, it's usually discouraged to have an attribute that is both `Optional` and `Computed` unless it's specifically meant for legacy compatibility or as a computed default. Furthermore, validation and `PlanModifiers` are combined such that if the user provides the wrong value, it might not be surfaced as clearly as desired.

Also, the value -1 is valid only when "share_mode" is "no limit", but this is only enforced at validation. However, with the current schema, the drift between configuration and actual API state could be allowed silently.

## Impact

Medium severity.

- The schema definition may lead to ambiguous provider behavior, including drift, unclear error reporting, and possibly failing plans or applies.
- Might allow users to create configurations that don't reflect the actual underlying system's constraints, making the provider less robust.

## Location

Lines referencing the `share_max_limit` attribute definition and its related validation logic:

```go
"share_max_limit": schema.NumberAttribute{
    MarkdownDescription: "Maximum total of individual who can be shared to: (-1..99). If `share_mode` is `No limit`, this value must be -1.",
    Optional:            true,
    Computed:            true,
    PlanModifiers: []planmodifier.Number{
        numberplanmodifier.UseStateForUnknown(),
    },
    Validators: []validator.Number{
        // validation for -1..99
        numbervalidator.OneOf(maxSharingRange...),
    },
},
```
and the validation logic:
```go
if sharingControl.ShareMode.ValueString() == "no limit" {
    if !sharingControl.ShareMaxLimit.IsNull() {
        resp.Diagnostics.AddAttributeError(
            path.Root("rules"),
            "sharing_controls validation error",
            "'share_max_limit' must be null when 'share_mode' is 'no limit'",
        )
    }
} else {
    if sharingControl.ShareMaxLimit.IsNull() || sharingControl.ShareMaxLimit.Equal(basetypes.NewFloat64Value(-1)) {
        resp.Diagnostics.AddAttributeError(
            path.Root("rules"),
            "sharing_controls validation error",
            "'share_max_limit' must be a value between 0 and 99 when 'share_mode' is 'exclude sharing with security groups'",
        )
    }
}
```

## Fix

Update the attribute and validation to more clearly reflect the semantic relationship, and avoid potential pitfalls of optional+computed unless truly intended.

**Proposed schema and logic:**

```go
"share_max_limit": schema.NumberAttribute{
    MarkdownDescription: "Maximum number of individuals who can be shared to: (-1..99). Required unless share_mode is `no limit`.",
    Optional:            true,
    PlanModifiers: []planmodifier.Number{
        numberplanmodifier.UseStateForUnknown(),
    },
    Validators: []validator.Number{
        numbervalidator.OneOf(maxSharingRange...),
    },
},
```
Update the validation to also set nulls correctly, and consider setting a default if necessary, or removing "Computed" to make explicit what happens:

```go
if sharingControl.ShareMode.ValueString() == "no limit" {
    if !sharingControl.ShareMaxLimit.IsNull() && !sharingControl.ShareMaxLimit.Equal(basetypes.NewFloat64Value(-1)) {
        resp.Diagnostics.AddAttributeError(
            path.Root("rules"),
            "sharing_controls validation error",
            "'share_max_limit' must be -1 or null when 'share_mode' is 'no limit'",
        )
    }
} else if sharingControl.ShareMode.ValueString() == "exclude sharing with security groups" {
    if sharingControl.ShareMaxLimit.IsNull() || sharingControl.ShareMaxLimit.Equal(basetypes.NewFloat64Value(-1)) {
        resp.Diagnostics.AddAttributeError(
            path.Root("rules"),
            "sharing_controls validation error",
            "'share_max_limit' must be a value between 0 and 99 when 'share_mode' is 'exclude sharing with security groups'",
        )
    }
}
```
