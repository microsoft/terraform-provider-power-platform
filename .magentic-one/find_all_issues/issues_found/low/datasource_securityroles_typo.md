# Title

The error message for `environment_id` typo in diagnostics text

##

`/workspaces/terraform-provider-power-platform/internal/services/authorization/datasource_securityroles.go`

## Problem

The Diagnostics error incorrectly spells the word "cannot" as "connot" in the text of the error message. Specifically, it says: `"environment_id connot be an empty string"`. This is an obvious typographical error.

## Impact

- **Severity**: **Low**

While this doesn't break the functionality, it can cause confusion or diminish the professional appearance of the codebase. It also impacts readability and clarity for users and developers encountering the error.

## Location

Inside the `Read` method of the `SecurityRolesDataSource`:

## Code Issue

```go
if state.EnvironmentId.ValueString() == "" {
	resp.Diagnostics.AddError("environment_id connot be an empty string", "environment_id connot be an empty string")
	return
}
```

## Fix

Correct the typo by updating `"connot"` to `"cannot"` in both instances of the error message text.

```go
if state.EnvironmentId.ValueString() == "" {
	resp.Diagnostics.AddError("environment_id cannot be an empty string", "environment_id cannot be an empty string")
	return
}
```

This makes the message accurate and readable for developers and users.