# Miscellaneous Spelling and Typo Issues

This document contains all spelling and typo issues found in miscellaneous files within the terraform-provider-power-platform codebase.

## ISSUE 1

### Misspelled function name `covertDlpPolicyToPolicyModelDto`

**File:** `/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/helpers.go`

**Problem:** The function name `covertDlpPolicyToPolicyModelDto` contains a typo (`covert` instead of `convert`). This could be confusing for maintainers and reduces the clarity and discoverability of your code.

**Impact:** Low severity. While this doesn't break code functionality, it can confuse developers and reduce code readability and maintainability due to the inconsistency in naming and possible misspelling.

**Location:** Line 16 - Function definition:

**Code Issue:**

```go
func covertDlpPolicyToPolicyModelDto(policy dlpPolicyDto) (*dlpPolicyModelDto, error) {
```

**Fix:** Rename the function to correct the spelling, and update all internal references to this function accordingly:

```go
func convertDlpPolicyToPolicyModelDto(policy dlpPolicyDto) (*dlpPolicyModelDto, error) {
```

## ISSUE 2

### Inconsistent Resource Naming: Typo in NewEnterpisePolicyResource

**File:** `/workspaces/terraform-provider-power-platform/internal/provider/provider.go`

**Problem:** There is a typo in the resource provider function `enterprise_policy.NewEnterpisePolicyResource()`. The word "Enterpise" should be "Enterprise".

**Impact:** Incorrect naming could cause confusion and violates naming consistency. It may also lead to import errors or break code referencing this function elsewhere if fixed without care. Severity: **medium**.

**Location:** In the Resources registration function, specifically this line:

**Code Issue:**

```go
func() resource.Resource { return enterprise_policy.NewEnterpisePolicyResource() },
```

**Fix:** Rename the function call to use "Enterprise" instead of "Enterpise":

```go
func() resource.Resource { return enterprise_policy.NewEnterprisePolicyResource() },
```

Be sure to fix the actual function name in the implementation file under the enterprise_policy package as well to maintain consistency.

---

Apply this fix to the whole codebase

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
