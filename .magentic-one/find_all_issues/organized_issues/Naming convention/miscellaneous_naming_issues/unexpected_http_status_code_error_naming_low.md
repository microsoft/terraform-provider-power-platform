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
