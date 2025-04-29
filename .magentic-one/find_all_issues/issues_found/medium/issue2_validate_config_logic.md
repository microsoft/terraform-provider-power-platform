# Title

Improper validation logic for attributes.

##

`/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy.go`

## Problem

The validation logic for `default_action_rule_behavior` and `action_rules` attributes in the `ValidateConfig` method is incorrect. Specifically, the logic expects `default_action_rule_behavior` to be empty if `action_rules` are empty, which is validated using an `if` condition:

```go
if (c.DefaultActionRuleBehavior != "" && len(c.ActionRules) == 0) || (c.DefaultActionRuleBehavior == "" && len(c.ActionRules) > 0) {
    resp.Diagnostics.AddAttributeError(
        path.Empty(),
        "Incorrect attribute Configuration",
        "Expected 'default_action_rule_behavior' to be empty if 'action_rules' are empty.",
    )
}
```

## Impact

This approach tightly couples the validation logic for these attributes, making it difficult to extend or modify in the future. Additionally, error messages are generic and do not provide meaningful guidance for correcting invalid configurations.
Severity: **medium**

## Location

Validation logic in `ValidateConfig` method.

## Code Issue

```go
if (c.DefaultActionRuleBehavior != "" && len(c.ActionRules) == 0) || (c.DefaultActionRuleBehavior == "" && len(c.ActionRules) > 0) {
    resp.Diagnostics.AddAttributeError(
        path.Empty(),
        "Incorrect attribute Configuration",
        "Expected 'default_action_rule_behavior' to be empty if 'action_rules' are empty.",
    )
}
```

## Fix

Improve validation logic and make error messages more context-specific. For example:

```go
if len(c.ActionRules) == 0 && c.DefaultActionRuleBehavior != "" {
    resp.Diagnostics.AddAttributeError(
        path.Empty(),
        "Invalid Configuration",
        "'default_action_rule_behavior' must be empty if there are no 'action_rules'.",
    )
} else if len(c.ActionRules) > 0 && c.DefaultActionRuleBehavior == "" {
    resp.Diagnostics.AddAttributeError(
        path.Empty(),
        "Invalid Configuration",
        "'default_action_rule_behavior' cannot be empty if 'action_rules' are specified.",
    )
}
```

This simplifies the validation logic and provides clearer error messages.