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
