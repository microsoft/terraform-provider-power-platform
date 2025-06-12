# Miscellaneous Redundant Code Issues

This document consolidates all miscellaneous redundant code issues found in various components of the Terraform Power Platform provider.

## ISSUE 1

**File:** `/workspaces/terraform-provider-power-platform/internal/provider/provider.go`

### Problem

The `validateProviderAttribute` function provides a confusing error message when the `environmentVariableName` string is non-empty. The advice "Target apply the source of the value first, set the value statically in the configuration." is unclear, and when `environmentVariableName` is passed, it redundantly says, "Either target apply the source of the value first, set the value statically in the configuration, or use the ... environment variable." This message could be clearer and more actionable.

### Impact

Ambiguous messages may confuse users and hinder troubleshooting. Severity: **low**.

### Location

Lines surrounding this fragment:

```go
environmentVariableText := "Target apply the source of the value first, set the value statically in the configuration."
if environmentVariableName != "" {
    environmentVariableText = fmt.Sprintf("Either target apply the source of the value first, set the value statically in the configuration, or use the %s environment variable.", environmentVariableName)
}

if value == "" {
    resp.Diagnostics.AddAttributeError(
        attrPath,
        fmt.Sprintf("Unknown %s", name),
        fmt.Sprintf("The provider cannot create the API client as there is an unknown configuration value for %s. %s", name, environmentVariableText))
}
```

### Fix

Use a concise and actionable message, e.g.:

```go
if environmentVariableName != "" {
    environmentVariableText = fmt.Sprintf("Set the value in the provider configuration or via the environment variable %s.", environmentVariableName)
} else {
    environmentVariableText = "Set the value in the provider configuration."
}
```

## ISSUE 2

**File:** `/workspaces/terraform-provider-power-platform/internal/customtypes/uuid_value.go`

### Problem

In the `StringSemanticEquals` method, if errors occur when parsing the old or new UUID values using `uuid.ParseUUID`, the diagnostics are populated by calling `diags.AddError`, but the method only checks if diagnostics exist with `diags.HasError()` and then returns `false, diags`. However, the flow does not stop parsing the new UUID if the old UUID check fails, potentially leading to misleading or redundant errors.

### Impact

**Low Severity**: Combining both errors together may confuse users, as the error with the old value is as significant as the new value and error reporting could be clearer if the function returned immediately after the first encountered error.

### Location

Method `StringSemanticEquals`

### Code Issue

```go
 oldUUID, err := uuid.ParseUUID(v.ValueString())
 if err != nil {
  diags.AddError("expected old value to be a valid UUID", err.Error())
 }

 newUUID, err := uuid.ParseUUID(newValue.ValueString())
 if err != nil {
  diags.AddError("expected new value to be a valid UUID", err.Error())
 }

 if diags.HasError() {
  return false, diags
 }

 return reflect.DeepEqual(oldUUID, newUUID), diags
```

### Fix

Return diagnostics immediately after encountering the first error (optional: for stricter validation/reporting), or document the intent to collect all parsing errors. Here's an example of returning on the first error (recommended for clarity):

```go
 oldUUID, err := uuid.ParseUUID(v.ValueString())
 if err != nil {
  diags.AddError("expected old value to be a valid UUID", err.Error())
  return false, diags
 }

 newUUID, err := uuid.ParseUUID(newValue.ValueString())
 if err != nil {
  diags.AddError("expected new value to be a valid UUID", err.Error())
  return false, diags
 }

 return reflect.DeepEqual(oldUUID, newUUID), diags
```

---

## To finish the task you have to

1. Run linter and fix any issues
2. Run UnitTest and fix any of failing ones
3. Generate docs
4. Run Changie

## Changie Instructions

Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```

Where:

- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number

Apply this fix to the whole codebase
