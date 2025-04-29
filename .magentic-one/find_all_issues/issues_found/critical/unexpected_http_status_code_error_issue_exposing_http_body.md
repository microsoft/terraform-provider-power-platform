# Title

Unnecessary Inclusion of `Body` in the Error Message

##

`/workspaces/terraform-provider-power-platform/internal/customerrors/unexpected_http_status_code_error.go`

## Problem

The `Error` method of the `UnexpectedHttpStatusCodeError` struct includes the `Body` field when formatting the error message. This could expose sensitive data or make logs unnecessarily verbose. Including the body of an HTTP response in an error message is generally discouraged unless absolutely necessary.

## Impact

Potential risk of data leaks if sensitive information is embedded within the HTTP response body, especially in production environments with logging enabled. This increases debugging complexity and creates a security risk. Severity is **critical**.

## Location

Line: The `Error` method implementation.

## Code Issue

```go
func (e UnexpectedHttpStatusCodeError) Error() string {
	return fmt.Sprintf("Unexpected HTTP status code. Expected: %v, received: [%d] %s | %s", e.ExpectedStatusCodes, e.StatusCode, e.StatusText, e.Body)
}
```

## Fix

Limit the error message to include only essential information (such as expected status codes, actual status code, and status text). Remove the `Body` from the output, ensuring the error message is succinct and non-sensitive.

```go
func (e UnexpectedHttpStatusCodeError) Error() string {
	return fmt.Sprintf("Unexpected HTTP status code. Expected: %v, received: [%d] %s", e.ExpectedStatusCodes, e.StatusCode, e.StatusText)
}
```