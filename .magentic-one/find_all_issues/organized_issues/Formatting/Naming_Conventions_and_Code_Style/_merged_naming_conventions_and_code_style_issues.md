# Naming Conventions and Code Style Issues

This document contains merged issues related to naming conventions and code style in the Power Platform Terraform provider.

## ISSUE 1

**Title:** Helper Function Name: `convertFromRuleDto` Should be Idiomatic

**File:** `/workspaces/terraform-provider-power-platform/internal/services/solution_checker_rules/models.go`

**Problem:**
The helper function `convertFromRuleDto` does not follow Go idioms for conversion or transformation functions, which typically use the `ToX` or `FromX` pattern to indicate which type is converting and to what, for example `RuleDtoToModel`. This improves readability and discoverability, aiding in onboarding and consistency.

**Impact:**
Non-idiomatic naming impacts consistency, onboarding, and makes the intentions of the function less clear. This could add confusion when scanning the codebase or using helper functions generally. Severity: **low**

**Location:**

```go
// Helper function to convert from DTO to Model.
func convertFromRuleDto(rule ruleDto) RuleModel {
```

**Code Issue:**

```go
func convertFromRuleDto(rule ruleDto) RuleModel {
```

**Fix:**
Rename the function to a more idiomatic Go style, for example:

```go
func RuleDtoToModel(rule RuleDto) RuleModel {
```

Or, if sticking to idiomatic abbreviations:

```go
func RuleDTOToModel(rule RuleDTO) RuleModel {
```

Also ensure parameter and type names match the naming conventions.

## ISSUE 2

**Title:** Function Naming: Non-Idiomatic Name 'RequireReplaceStringFromNonEmptyPlanModifier'

**File:** `/workspaces/terraform-provider-power-platform/internal/modifiers/requires_replace_string_from_non_empty_modifier.go`

**Problem:**
The function name `RequireReplaceStringFromNonEmptyPlanModifier` does not follow Go idioms (like PascalCase for exported functions), and is excessively verbose and difficult to quickly parse. Good function names should be concise but meaningful and generally avoid repeating information unnecessarily. This can make function usages more readable and maintainable.

**Impact:**
**Severity: Low**

Verbose or unclear naming can reduce readability and make the code slightly harder to maintain or consume, especially as your codebase or team grows. Go recommends being brief but accurate in naming, relying on package context.

**Location:**

```go
func RequireReplaceStringFromNonEmptyPlanModifier() planmodifier.String {
 return &requireReplaceStringFromNonEmptyPlanModifier{}
}
```

**Code Issue:**

```go
func RequireReplaceStringFromNonEmptyPlanModifier() planmodifier.String {
 return &requireReplaceStringFromNonEmptyPlanModifier{}
}
```

**Fix:**
Rename the function and corresponding struct to follow Go best practices. For example:

```go
func ReplaceOnNonEmptyStringChange() planmodifier.String {
 return &replaceOnNonEmptyStringChange{}
}
```

This aligns with Go naming conventions (concise and descriptive), and pairs well with revised struct naming. Update usages in the project as needed.

## ISSUE 3

**Title:** Struct Naming: Lack of Meaningful Name for `requireReplaceStringFromNonEmptyPlanModifier`

**File:** `/workspaces/terraform-provider-power-platform/internal/modifiers/requires_replace_string_from_non_empty_modifier.go`

**Problem:**
The struct `requireReplaceStringFromNonEmptyPlanModifier` uses an excessively long and verbose name, which can lead to decreased readability and maintainability. A struct name should be concise while still giving enough information about its purpose. This specific name makes code harder to read and unnecessarily repeats information that could be derived from context or documentation.

**Impact:**
**Severity: Low**

Long and verbose names can hinder code readability and make maintenance harder, especially when such names are used throughout the codebase. While not a critical issue, improving naming conventions contributes to cleaner, more professional code.

**Location:**

```go
type requireReplaceStringFromNonEmptyPlanModifier struct {
}
```

**Code Issue:**

```go
type requireReplaceStringFromNonEmptyPlanModifier struct {
}
```

**Fix:**
Consider renaming the struct to a shorter, more meaningful name such as `replaceOnNonEmptyChangeModifier` or `replaceOnNonEmptyStringChange`. Here is an example:

```go
type replaceOnNonEmptyStringChange struct {
}
```

If you change the struct name, also update any references (including constructor function name) for consistency, such as:

```go
func ReplaceOnNonEmptyStringChangeModifier() planmodifier.String {
 return &replaceOnNonEmptyStringChange{}
}
```

This improves readability and aligns with Go naming best practices.

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
