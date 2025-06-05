# Title

Error Message May Leak Underlying Error Detail

##
/workspaces/terraform-provider-power-platform/internal/customerrors/url_format_error.go

## Problem

The error message formatting in the `Error()` method directly includes the underlying error message returned from `e.Err.Error()`. This can potentially leak internal or sensitive error details to the caller or logs, which might not be appropriate for end-users and could pose a security concern.

## Impact

In cases where `e.Err` contains internal context (such as stack traces, credentials, or sensitive configuration), this could inadvertently expose sensitive information outside of expected logging channels. The severity is **medium** because this can have consequences in production environments and public logs.

## Location

Method: `func (e UrlFormatError) Error() string`
File: /workspaces/terraform-provider-power-platform/internal/customerrors/url_format_error.go

## Code Issue

```go
func (e UrlFormatError) Error() string {
	errorMsg := ""
	if e.Err != nil {
		errorMsg = e.Err.Error()
	}

	return fmt.Sprintf("Request url must be an absolute url: '%s' : '%s'", e.Url, errorMsg)
}
```

## Fix

Carefully sanitize or wrap error information. Only propagate user-friendly or non-sensitive messages in user-facing error strings. If additional debugging is needed, use logging (not error messages) for sensitive details.

```go
func (e UrlFormatError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Request URL must be an absolute URL: '%s' : error occurred during URL parsing/validation.", e.Url)
	}
	return fmt.Sprintf("Request URL must be an absolute URL: '%s'", e.Url)
}
```

If internal error details are useful, expose them only in logs under a debug mode, not in the error string returned.
