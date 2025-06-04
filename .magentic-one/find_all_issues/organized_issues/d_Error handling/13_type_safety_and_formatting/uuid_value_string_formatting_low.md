# Incorrect Use of Format String with Escaped Newlines in Error Message Constant

##

/workspaces/terraform-provider-power-platform/internal/customtypes/uuid_value.go

## Problem

The constant `UUIDTypeErrorInvalidStringDetails` is defined as a raw string literal with literal `\n` escape characters, not as actual newlines, due to the use of backticks and explicit `\n` (not actual newlines). This will result in error messages containing literal `\\n\\n` instead of real line breaks when used with `fmt.Sprintf`.

## Impact

**Low Severity**: The user experience in diagnostics and error logs will be affected due to improperly formatted error messages. This does not impact functionality but makes logs harder to read, which impairs debugging and support.

## Location

Definition of `UUIDTypeErrorInvalidStringDetails`

## Code Issue

```go
const (
	UUIDTypeErrorInvalidStringHeader  = "Invalid UUID String Value"
	UUIDTypeErrorInvalidStringDetails = `A string value was provided that is not valid UUID string format.\\n\\nGiven Value: %s\\n`
)
```

## Fix

Use a regular string literal with proper newlines or, if using a raw string with backticks, include actual blank lines:

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
