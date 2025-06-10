# Error Logging Without Propagation - Plan Modifier Issues

This document consolidates all issues related to error logging without proper propagation found in plan modifier implementations across the Terraform Provider for Power Platform.

## ISSUE 1

# Unhandled error from `GetKey` function call

##

/workspaces/terraform-provider-power-platform/internal/modifiers/force_string_value_unknown_modifier.go

## Problem

In the `PlanModifyString` method, the return value `err` from `req.Private.GetKey(ctx, "force_value_unknown")` is ignored. By using the blank identifier `_`, any error from this function call would be discarded, possibly masking runtime issues that would be useful to log or to fail gracefully.

## Impact

If `GetKey` returns an error, it is silently ignored. This could cause the function to behave incorrectly or make debugging more difficult, especially if the expected key is missing or an internal deserialization failed. Severity: medium.

## Location

```go
func (d *forceStringValueUnknownModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
 r, _ := req.Private.GetKey(ctx, "force_value_unknown")
 if r == nil || !bytes.Equal(r, []byte("true")) {
  return
 }
 resp.PlanValue = types.StringUnknown()
}
```

## Fix

Handle the error returned by `GetKey`. Depending on the desired behavior, you may want to log, abort, or set diagnostic information in the response.

```go
func (d *forceStringValueUnknownModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
 r, err := req.Private.GetKey(ctx, "force_value_unknown")
 if err != nil {
  // Optionally, add error to diagnostics (if appropriate for the response in your framework)
  // resp.Diagnostics.AddError("Error reading private key", err.Error())
  return
 }
 if r == nil || !bytes.Equal(r, []byte("true")) {
  return
 }
 resp.PlanValue = types.StringUnknown()
}
```

## ISSUE 2

# Title

Potentially Confusing Conditional in PlanModifyInt64

##

/workspaces/terraform-provider-power-platform/internal/modifiers/require_replace_int_attribute_modifier.go

## Problem

The conditional in `PlanModifyInt64` is currently implemented as a single dense line:

```go
if req.PlanValue != req.StateValue && (!req.StateValue.IsNull() && !req.StateValue.IsUnknown() && req.StateValue.ValueInt64() != 0) {
    resp.RequiresReplace = true
}
```

This logic is not self-explanatory and could easily cause confusion or lead to errors if maintained in the future without clear documentation or refactoring. It can be made more readable by splitting the condition into well-named variables and adding proper documentation.

## Impact

Medium: Reduced code readability and maintainability, as future changes to the condition might introduce mistakes. It's also difficult to debug or audit the business logic.

## Location

Within the method `PlanModifyInt64`:

```go
func (d *requireReplaceIntAttributePlanModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
    if req.PlanValue != req.StateValue && (!req.StateValue.IsNull() && !req.StateValue.IsUnknown() && req.StateValue.ValueInt64() != 0) {
        resp.RequiresReplace = true
    }
}
```

## Code Issue

```go
if req.PlanValue != req.StateValue && (!req.StateValue.IsNull() && !req.StateValue.IsUnknown() && req.StateValue.ValueInt64() != 0) {
    resp.RequiresReplace = true
}
```

## Fix

Refactor the conditional to use well-named intermediate variables, improving readability and making logic modification safer.

```go
isValueChanged := req.PlanValue != req.StateValue
hasValidPreviousState := !req.StateValue.IsNull() && !req.StateValue.IsUnknown() && req.StateValue.ValueInt64() != 0

if isValueChanged && hasValidPreviousState {
    resp.RequiresReplace = true
}
```

Add a comment to document why this logic is in place to prevent confusion for future maintainers.

## ISSUE 3

# Diagnostics May Contain Errors But Processing Continues

##

/workspaces/terraform-provider-power-platform/internal/modifiers/set_bool_value_unknown_if_checksum_change_modifier.go

## Problem

In `hasChecksumChanged`, after appending diagnostics from `GetAttribute`, there is no check to see if an error occurred before proceeding with the rest of the logic. If `GetAttribute` fails, `attribute` or `attributeChecksum` might not be populated with valid data, which could lead to incorrect calculations or misleading diagnostics.

## Impact

This can result in attempts to calculate a SHA256 value or compare checksums using invalid or empty data, making error detection and debugging harder. The severity is medium, as it may lead to non-obvious erroneous state propagating downstream.

## Location

```go
diags := req.Plan.GetAttribute(ctx, path.Root(attributeName), &attribute)
resp.Diagnostics.Append(diags...)
...
diags = req.State.GetAttribute(ctx, path.Root(checksumAttributeName), &attributeChecksum)
resp.Diagnostics.Append(diags...)
...
value, err := helpers.CalculateSHA256(attribute.ValueString())
```

## Code Issue

```go
diags := req.Plan.GetAttribute(ctx, path.Root(attributeName), &attribute)
resp.Diagnostics.Append(diags...)

var attributeChecksum types.String
diags = req.State.GetAttribute(ctx, path.Root(checksumAttributeName), &attributeChecksum)
resp.Diagnostics.Append(diags...)

value, err := helpers.CalculateSHA256(attribute.ValueString())
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Error calculating MD5 checksum for %s", attribute), err.Error())
}
```

## Fix

After each `Append` of diagnostics, check if there was an error and halt further processing if needed:

```go
diags := req.Plan.GetAttribute(ctx, path.Root(attributeName), &attribute)
resp.Diagnostics.Append(diags...)
if diags.HasError() {
    return false
}

var attributeChecksum types.String
diags = req.State.GetAttribute(ctx, path.Root(checksumAttributeName), &attributeChecksum)
resp.Diagnostics.Append(diags...)
if diags.HasError() {
    return false
}
```

# To finish the task you have to

1. Run linter and fix any issues
2. Run UnitTest and fix any of failing ones
3. Generate docs
4. Run Changie

# Changie Instructions

Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```

Where:

- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
