# Type Safety and Formatting Issues

This document consolidates issues related to type safety, string formatting, and validation that need to be addressed to improve code reliability and user experience.

## ISSUE 1

**File:** `/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy.go`

**Problem:** Missing Validation for Unknown/Null String Attribute in Schema

In the resource schema, required attributes such as `"name"`, `"location"`, and those in `"billing_instrument"` allow an empty string by default since there are no explicit `stringvalidator.LengthAtLeast(1)` (or similar) validators. Terraform might send an empty or whitespace string, resulting in failed API calls or unclear errors.

**Impact:** Severity: **Medium**
Missing explicit validation could cause confusing errors for users, potential API rejections, and an inconsistent user experience if empty strings are passed to the backend.

**Location:** `Schema` method:  

```go
"name": schema.StringAttribute{
    MarkdownDescription: "The name of the billing policy",
    Required:            true,
},
// ... and similar for others
```

**Code Issue:**

```go
"name": schema.StringAttribute{
    MarkdownDescription: "The name of the billing policy",
    Required:            true,
},
```

**Fix:** Add `stringvalidator.LengthAtLeast(1)` in the `Validators` for each required string-type attribute to prevent empty string assignments.

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

## ISSUE 2

**File:** `/workspaces/terraform-provider-power-platform/internal/customerrors/url_format_error.go`

**Problem:** Interface Implementation Assertion as Value Instead of Pointer

The variable `_ error = UrlFormatError{}` is used to assert that `UrlFormatError` implements the `error` interface. While this is valid if the `Error()` method has a value receiver (as in the present code), it can be a potential maintainability concern: if the receiver of `Error()` is later changed to a pointer, this assertion will no longer be valid. In most Go codebases, error types are implemented with pointer receivers to allow for more flexibility, such as mutability, and in such cases, the interface assertion should also use a pointer.

**Impact:** If the receiver type changes in the future to a pointer (e.g. `func (e *UrlFormatError) Error() string`), this assertion will silently fail to assert interface satisfaction at compile time, possibly introducing errors later on. Severity is **low**, but impacts future-proofing and maintainability.

**Location:** Global variable: `var _ error = UrlFormatError{}`

**Code Issue:**

```go
var _ error = UrlFormatError{}
```

**Fix:** If the receiver will always be a value receiver, leave as is. However, best practice is to implement the error interface with pointer receivers and assert interface implementation using pointer as well:

```go
var _ error = (*URLFormatError)(nil)
```

The `Error()` method would then look like:

```go
func (e *URLFormatError) Error() string {
    // implementation
}
```

This is more idiomatic for future extensibility.

## ISSUE 3

**File:** `/workspaces/terraform-provider-power-platform/internal/customtypes/uuid_value.go`

**Problem:** Incorrect Use of Format String with Escaped Newlines in Error Message Constant

The constant `UUIDTypeErrorInvalidStringDetails` is defined as a raw string literal with literal `\n` escape characters, not as actual newlines, due to the use of backticks and explicit `\n` (not actual newlines). This will result in error messages containing literal `\\n\\n` instead of real line breaks when used with `fmt.Sprintf`.

**Impact:** **Low Severity**: The user experience in diagnostics and error logs will be affected due to improperly formatted error messages. This does not impact functionality but makes logs harder to read, which impairs debugging and support.

**Location:** Definition of `UUIDTypeErrorInvalidStringDetails`

**Code Issue:**

```go
const (
        UUIDTypeErrorInvalidStringHeader  = "Invalid UUID String Value"
        UUIDTypeErrorInvalidStringDetails = `A string value was provided that is not valid UUID string format.\\n\\nGiven Value: %s\\n`
)
```

**Fix:** Use a regular string literal with proper newlines or, if using a raw string with backticks, include actual blank lines:

```go
const (
        UUIDTypeErrorInvalidStringHeader  = "Invalid UUID String Value"
        UUIDTypeErrorInvalidStringDetails = "A string value was provided that is not valid UUID string format.\n\nGiven Value: %s\n"
)
```

Or with backticks (no escapes needed):

```go
const (
        UUIDTypeErrorInvalidStringHeader  = "Invalid UUID String Value"
        UUIDTypeErrorInvalidStringDetails = `A string value was provided that is not valid UUID string format.

Given Value: %s
`
)
```

---

## Task Completion Instructions

After implementing these fixes:

1. **Run the linter:** `make lint` to ensure code style compliance
2. **Run unit tests:** `make unittest` to verify functionality  
3. **Generate documentation:** `make userdocs` to update auto-generated docs
4. **Add changelog entry:** Use `changie new` to document the changes

**Changie Entry Template:**

```yaml
kind: fixed
body: Fixed type safety and formatting issues including missing string validation, interface assertions, and escaped newline formatting
time: [current-timestamp]
custom:
  Issue: "[ISSUE_NUMBER_IF_APPLICABLE]"
```

Replace `[ISSUE_NUMBER_IF_APPLICABLE]` with the relevant GitHub issue number, or remove the custom section if no specific issue exists.
