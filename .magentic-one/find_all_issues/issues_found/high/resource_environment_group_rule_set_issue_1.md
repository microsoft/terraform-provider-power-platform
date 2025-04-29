# Title

Improper Validation Logic for `sharing_controls` Attributes

##

`/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/resource_environment_group_rule_set.go`

## Problem

The validation logic for the `sharing_controls` attribute does not properly handle cases where required fields should be `null`. For instance:
- When `share_mode` is `"no limit"`, the corresponding `share_max_limit` should not accept any value except `null`. 
- The validation implemented in the `ValidateConfig` method does allow such incorrect configurations as it fails to check for specific values properly when `share_mode` is `"exclude sharing with security groups"`.

## Impact

- It will create inconsistent configurations in the Terraform state, leading to unpredictable failures or misconfigurations.
- Misleading validation errors for users due to improperly defined rules.

**Severity: High**

## Location

- Method: `ValidateConfig`
- Code Snippet:

## Code Issue

Current faulty implementation:
```go
	sharingControlsObj := config.Rules.Attributes()["sharing_controls"]
	if !sharingControlsObj.IsNull() && !sharingControlsObj.IsUnknown() {
		var sharingControl environmentGroupRuleSetSharingControlsModel
		sharingControlsObj.(basetypes.ObjectValue).As(ctx, &sharingControl, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

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
	}
```

## Fix

Updated conditional logic and validation for clarity:
```go
	sharingControlsObj := config.Rules.Attributes()["sharing_controls"]
	if !sharingControlsObj.IsNull() && !sharingControlsObj.IsUnknown() {
		var sharingControl environmentGroupRuleSetSharingControlsModel
		sharingControlsObj.(basetypes.ObjectValue).As(ctx, &sharingControl, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		switch sharingControl.ShareMode.ValueString() {
		case "no limit":
			if sharingControl.ShareMaxLimit.Value != nil && !sharingControl.ShareMaxLimit.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("rules"),
					"Invalid Sharing Controls",
					"When 'share_mode' is 'no limit', 'share_max_limit' must be null.",
				)
			}
		case "exclude sharing with security groups":
			if sharingControl.ShareMaxLimit.IsNull() || sharingControl.ShareMaxLimit.Equal(basetypes.NewFloat64Value(-1)) {
				resp.Diagnostics.AddAttributeError(
					path.Root("rules"),
					"Invalid Sharing Controls",
					"'share_max_limit' must be between 0 and 99 when 'share_mode' is 'exclude sharing with security groups'.",
				)
			}
		default:
			resp.Diagnostics.AddAttributeError(
				path.Root("rules"),
				"Unsupported Share Mode",
				fmt.Sprintf("'%s' is not a valid value for 'share_mode'. Allowed values are 'no limit', 'exclude sharing with security groups'.", sharingControl.ShareMode.ValueString()),
			)
		}
	}
```