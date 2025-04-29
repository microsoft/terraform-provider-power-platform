# Title

Plan Modifier is Commented Out Without Explanation in `power_apps` Attribute

##

`/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go`

## Problem

In the `power_apps` section of the schema definition, the plan modifier `objectplanmodifier.UseStateForUnknown()` is commented out without any explanation or justification. This results in unclear behavior regarding state management for this resource attribute, as the effective plan modifier logic for this attribute is missing.

## Impact

Commenting out critical configuration components without explanation can lead to unintended behavior in state management, especially for complex attributes like `power_apps`. Developers maintaining the code may make erroneous assumptions regarding the intentionality of this omission. Severity: **Medium**.

## Location

Line 121-125: In the schema definition for `power_apps` under `power_platform`.

## Code Issue

```go
"power_apps": schema.SingleNestedAttribute{
    Description:   "Power Apps",
    Optional:      true,
    PlanModifiers: []planmodifier.Object{
        // objectplanmodifier.UseStateForUnknown(),
    },
    Attributes: map[string]schema.Attribute{
```

## Fix

Provide an explanation for why the plan modifier is commented out, or re-enable it if its omission is accidental. Including comments or notes improves code clarity and reduces the chances of future errors.

```go
"power_apps": schema.SingleNestedAttribute{
    Description:   "Power Apps",
    Optional:      true,
    PlanModifiers: []planmodifier.Object{
        objectplanmodifier.UseStateForUnknown(), // State tracking is essential if expected mutations occur.
    },
    Attributes: map[string]schema.Attribute{
```

Alternatively, if commenting was intentional, state reasoning:

```go
"power_apps": schema.SingleNestedAttribute{
    Description:   "Power Apps",
    Optional:      true,
    PlanModifiers: []planmodifier.Object{
        // objectplanmodifier.UseStateForUnknown(), // Refer to issue #XYZ - Removed to address X problem in state resolution.
    },
    Attributes: map[string]schema.Attribute{
```