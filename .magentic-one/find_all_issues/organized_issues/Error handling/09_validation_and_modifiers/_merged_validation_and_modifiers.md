# Validation and Modifiers Issues

This document consolidates all issues related to validation logic and plan modifiers in the Terraform Provider for Power Platform.

## ISSUE 1

### Unchecked Error from config.GetAttribute

**File:** `/workspaces/terraform-provider-power-platform/internal/services/data_record/dynamic_columns_validator.go`

**Problem:** The returned error from `config.GetAttribute` is ignored (assigned to `_`). If `GetAttribute` fails, the subsequent code may operate on invalid or undefined data.

**Impact:** Can potentially lead to data inconsistency or mask actual configuration errors the user should be aware of. Severity: **medium**.

**Code Issue:**

```go
_ = config.GetAttribute(ctx, matchedPaths[0], &dynamicColumns)
```

**Fix:** Check the error and add to diagnostics if present:

```go
if err := config.GetAttribute(ctx, matchedPaths[0], &dynamicColumns); err != nil {
    diags.AddError("Failed to get dynamic columns attribute", err.Error())
    return diags
}
```

## ISSUE 2

### Diagnostic Message Is Not Actionable

**File:** `/workspaces/terraform-provider-power-platform/internal/validators/make_field_required_when_other_field_does_not_have_value_validator.go`

**Problem:** The diagnostic error message returned when path matches are not found is not actionable or descriptive to the end user. Messages such as `"Other field required when value of validator should have exactly one match"` do not provide enough information to guide users to resolve the configuration error.

**Impact:** Low severity. Affects user experience and troubleshooting, but not runtime correctness.

**Code Issue:**

```go
res.Diagnostics.AddError("Other field required when value of validator should have exactly one match", "")
```

**Fix:** Improve the diagnostic messages to provide context and guidance:

```go
res.Diagnostics.AddError(
        "Validator Configuration Error: Other field match failed",
        "The validator could not uniquely locate the other field in the configuration. Ensure that 'OtherFieldExpression' matches exactly one attribute.",
)
```

## ISSUE 3

### Duplicate Error Message in AddError Usage

**File:** `/workspaces/terraform-provider-power-platform/internal/validators/other_field_required_when_value_of_validator.go`

**Problem:** In the line where diagnostics are added for an invalid "other field" value, the `ErrorMessage` is used as both the summary and the detail, causing redundancy. Ideally, the detail section should provide extra context or actionable information, not just reiterate the summary.

**Impact:** This reduces clarity and effectiveness of error diagnostics, especially for users and maintainers debugging complex configurations. Severity: **low**.

**Code Issue:**

```go
        res.Diagnostics.AddError(av.ErrorMessage, av.ErrorMessage)
```

**Fix:** Change to provide actionable/contextual detail:

```go
        res.Diagnostics.AddError(
                av.ErrorMessage,
                "Field \""+paths[0].String()+"\" does not meet required value conditions.",
        )
```

## ISSUE 4

### Error Message Mentions Wrong Hash Algorithm

**File:** `/workspaces/terraform-provider-power-platform/internal/modifiers/set_bool_value_unknown_if_checksum_change_modifier.go`

**Problem:** In the `hasChecksumChanged` method, when the error message is generated for a checksum calculation failure, it states "Error calculating MD5 checksum" even though the code actually uses SHA256 via `helpers.CalculateSHA256`.

**Impact:** This discrepancy can cause confusion during debugging and mislead engineers regarding the type of hash being calculated. The severity is low as it only affects log/error clarity, but should be corrected for accuracy.

**Code Issue:**

```go
resp.Diagnostics.AddError(fmt.Sprintf("Error calculating MD5 checksum for %s", attribute), err.Error())
```

**Fix:** Replace "MD5" with "SHA256" in the error string:

```go
resp.Diagnostics.AddError(fmt.Sprintf("Error calculating SHA256 checksum for %s", attribute), err.Error())
```

## ISSUE 5

### Inadequate Diagnostic Context for Attribute in Checksum Calculation

**File:** `/workspaces/terraform-provider-power-platform/internal/modifiers/set_string_value_unknown_if_checksum_change_modifier.go`

**Problem:** When reporting an error in `hasChecksumChanged`, the diagnostic message uses `attribute` (which could be a zero value if unmarshalling fails) rather than the name of the attribute being processed.

**Impact:** Error messages do not clearly specify which attribute caused the error, making debugging more difficult. Severity: medium.

**Code Issue:**

```go
resp.Diagnostics.AddError(fmt.Sprintf("Error calculating MD5 checksum for %s", attribute), err.Error())
```

**Fix:** Report the name of the attribute, not its value:

```go
resp.Diagnostics.AddError(fmt.Sprintf("Error calculating SHA256 checksum for attribute %q", attributeName), err.Error())
```

## ISSUE 6

### Redundant Diagnostic Appending

**File:** `/workspaces/terraform-provider-power-platform/internal/modifiers/set_string_value_unknown_if_checksum_change_modifier.go`

**Problem:** In `hasChecksumChanged`, diagnostics are appended with `resp.Diagnostics.Append(diags...)` even if the diagnostics are empty. It is usually better to check for errors and return early if any diagnostic occurs, as further steps might not be meaningful/valid on failed get operations.

**Impact:** Potential misleading diagnostics and performing unnecessary operations if attribute fetching fails. Severity: medium.

**Code Issue:**

```go
diags := req.Plan.GetAttribute(ctx, path.Root(attributeName), &attribute)
resp.Diagnostics.Append(diags...)
diags = req.State.GetAttribute(ctx, path.Root(checksumAttributeName), &attributeChecksum)
resp.Diagnostics.Append(diags...)
```

**Fix:** Check if getting the attribute caused errors, and return early if necessary:

```go
diags := req.Plan.GetAttribute(ctx, path.Root(attributeName), &attribute)
resp.Diagnostics.Append(diags...)
if diags.HasError() {
    return false
}
diags = req.State.GetAttribute(ctx, path.Root(checksumAttributeName), &attributeChecksum)
resp.Diagnostics.Append(diags...)
if diags.HasError() {
    return false
}
```

## ISSUE 7

### Unnecessary Separate Handling for Empty Checksum Value

**File:** `/workspaces/terraform-provider-power-platform/internal/modifiers/sync_attribute_plan_modifier.go`

**Problem:** After calculating the checksum, the code checks `if value == ""` and sets `resp.PlanValue = types.StringUnknown()`. Normally, a checksum calculation failure (which could result in an empty string) should already be handled by the error check above; reaching this branch is ambiguous and may mask edge cases.

**Impact:** This check may hide potential errors in checksum generation or underlying issues, making it harder to diagnose problems. It also doesn't document why an empty hash should be interpreted as "unknown." Severity: **medium**.

**Code Issue:**

```go
if value == "" {
        resp.PlanValue = types.StringUnknown()
} else {
        resp.PlanValue = types.StringValue(value)
}
```

**Fix:** Document explicitly why an empty hash might occur, or treat this as an error/diagnostic:

```go
if value == "" {
        resp.Diagnostics.AddError(fmt.Sprintf("Checksum is empty for %s", d.syncAttribute), "Calculated SHA256 checksum resulted in an empty value, which is unexpected.")
        resp.PlanValue = types.StringUnknown()
} else {
        resp.PlanValue = types.StringValue(value)
}
```

## ISSUE 8

### Incorrect Error Message Refers to MD5 Instead of SHA256

**File:** `/workspaces/terraform-provider-power-platform/internal/modifiers/sync_attribute_plan_modifier.go`

**Problem:** The error message in the error handling references "MD5 checksum" while the function actually computes a SHA256 hash using `helpers.CalculateSHA256()`. This introduces confusion and misleads maintainers or users as to which hashing algorithm is being used.

**Impact:** This can cause confusion during debugging or auditing, potentially leading to diagnostic errors or incorrect assumptions about security. Severity: **low**.

**Code Issue:**

```go
resp.Diagnostics.AddError(fmt.Sprintf("Error calculating MD5 checksum for %s", d.syncAttribute), err.Error())
```

**Fix:** Correct the error message to reference SHA256 instead of MD5:

```go
resp.Diagnostics.AddError(fmt.Sprintf("Error calculating SHA256 checksum for %s", d.syncAttribute), err.Error())
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
