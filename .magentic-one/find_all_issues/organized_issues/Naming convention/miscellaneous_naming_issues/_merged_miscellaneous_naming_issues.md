# Miscellaneous Naming Issues - Merged Issues

## ISSUE 1

# Title
TokenExpiredError follows Go naming but message field should be unexported

##
/workspaces/terraform-provider-power-platform/internal/api/auth.go

## Problem
The `TokenExpiredError` struct is exported and follows Go naming. However, the field `Message` is exported but is only used internally to the error struct and is not used as part of any JSON serialization, nor is required to be visible outside of the package. Per Go convention, unexport fields unless there is a clear need for exporting. The field should be lowercase (`message`).

## Impact
This is a **low severity** maintainability/naming issue. It increases unnecessary API surface and can be confusing for package users or generate extra documentation noise.

## Location
At the top of the file:

## Code Issue
```go
type TokenExpiredError struct {
	Message string
}
```

## Fix
Change `Message` to `message`:

```go
type TokenExpiredError struct {
	message string
}

func (e *TokenExpiredError) Error() string {
	return e.message
}
```

If you want to construct it with a public message, consider a constructor function, e.g.:

```go
func NewTokenExpiredError(msg string) error {
	return &TokenExpiredError{message: msg}
}
```


---

## ISSUE 2

# Issue 1

Unexpected Type Naming Convention

##

/workspaces/terraform-provider-power-platform/internal/customerrors/unexpected_http_status_code_error.go

## Problem

The struct and function use the name "Http" instead of the Go standard "HTTP" for abbreviations. In Go, common initialisms and acronyms should use all uppercase letters (e.g., "HTTP" not "Http") per Go naming conventions.

## Impact

- Reduces code readability and maintainability for developers familiar with Go naming conventions.
- Potential confusion/inconsistency when interacting with other code following the convention.
- Severity: Low

## Location

- Struct: `UnexpectedHttpStatusCodeError`
- Functions: `NewUnexpectedHttpStatusCodeError`

## Code Issue

```go
type UnexpectedHttpStatusCodeError struct {
	ExpectedStatusCodes []int
	StatusCode          int
	StatusText          string
	Body                []byte
}

func (e UnexpectedHttpStatusCodeError) Error() string {
	return fmt.Sprintf("Unexpected HTTP status code. Expected: %v, received: [%d] %s | %s", e.ExpectedStatusCodes, e.StatusCode, e.StatusText, e.Body)
}

func NewUnexpectedHttpStatusCodeError(expectedStatusCodes []int, statusCode int, statusText string, body []byte) error {
	return UnexpectedHttpStatusCodeError{
		ExpectedStatusCodes: expectedStatusCodes,
		StatusCode:          statusCode,
		StatusText:          statusText,
		Body:                body,
	}
}
```

## Fix

Refactor names to use `HTTP` throughout for better adherence to Go standards.

```go
type UnexpectedHTTPStatusCodeError struct {
	ExpectedStatusCodes []int
	StatusCode          int
	StatusText          string
	Body                []byte
}

func (e UnexpectedHTTPStatusCodeError) Error() string {
	return fmt.Sprintf("Unexpected HTTP status code. Expected: %v, received: [%d] %s | %s", e.ExpectedStatusCodes, e.StatusCode, e.StatusText, e.Body)
}

func NewUnexpectedHTTPStatusCodeError(expectedStatusCodes []int, statusCode int, statusText string, body []byte) error {
	return UnexpectedHTTPStatusCodeError{
		ExpectedStatusCodes: expectedStatusCodes,
		StatusCode:          statusCode,
		StatusText:          statusText,
		Body:                body,
	}
}
```
---
 
Apply for the whole code base


---

## ISSUE 3

# Title

Naming Convention Does Not Follow Go Standards for "Url"

##

/workspaces/terraform-provider-power-platform/internal/customerrors/url_format_error.go

## Problem

The type and field names use `Url` instead of the standard Go naming of `URL` (all caps), such as in `UrlFormatError` and the struct field `Url string`. According to Go conventions and the Go Code Review Comments (see: https://github.com/golang/go/wiki/CodeReviewComments#initialisms), initialisms should be all uppercase.

## Impact

Not following naming conventions impacts maintainability and readability, causing confusion for new maintainers and making it harder to search for similar constructs using standard names (`URL` instead of `Url`). This is a **low** severity issue but does affect professional code quality.

## Location

- Type name: `UrlFormatError`
- Struct field: `Url string`
- File: /workspaces/terraform-provider-power-platform/internal/customerrors/url_format_error.go

## Code Issue

```go
type UrlFormatError struct {
	Url string
	Err error
}

func NewUrlFormatError(url string, err error) error {
	return UrlFormatError{
		Err: err,
		Url: url,
	}
}
```

## Fix

Change all occurrences from `Url` to `URL` in type names, variable names, and struct fields, and update references accordingly.

```go
type URLFormatError struct {
	URL string
	Err error
}

func NewURLFormatError(url string, err error) error {
	return URLFormatError{
		Err: err,
		URL: url,
	}
}
```

All usages should be updated to follow this correction.


---

## ISSUE 4

# Title

Inconsistent Naming: Go Types Should Not Be Suffixed with 'Type'

##

/workspaces/terraform-provider-power-platform/internal/customtypes/uuid_type.go

## Problem

The struct is named `UUIDType`, which is redundant and non-idiomatic in Go. According to Go naming conventions, the “Type” suffix should be avoided unless it specifically disambiguates. Since this code is in a custom types package, a better name would be simply `UUID` or similar.

## Impact

Severity: **Low**

- May reduce readability and burden cross-reference and refactoring processes.
- May complicate code, especially when using code navigation tools or documentation.

## Location

```go
type UUIDType struct {
	basetypes.StringType
}
```

## Code Issue

```go
type UUIDType struct {
	basetypes.StringType
}
```

## Fix

- Rename the struct to `UUID`, and update all associated usages within the codebase.

Example:

```go
type UUID struct {
	basetypes.StringType
}
```


---

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
