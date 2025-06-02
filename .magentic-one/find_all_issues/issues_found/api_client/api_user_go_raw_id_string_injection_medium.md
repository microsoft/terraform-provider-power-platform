# Title

Potential String Injection Using Raw IDs in HTTP Path

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/api_user.go

## Problem

In multiple functions (such as `GetDataverseUserBySystemUserId`, `UpdateDataverseUser`, `DeleteDataverseUser`, `RemoveDataverseSecurityRoles`, `AddDataverseSecurityRoles`), the `systemUserId` (or `roleId`) is interpolated directly into the HTTP path string without any URL path escaping. This can result in malformed URLs if the IDs contain unexpected or invalid URL path characters, or in edge cases, could be vulnerable to injection if the ID is not properly filtered (for example, a specially crafted ID could alter the intended HTTP request).

## Impact

Severity: Medium

This could cause failures if IDs contain special characters, and potentially be a vector for path injection or ambiguous logs. While the risk of controlled injection is somewhat low if IDs are always GUIDs, this is not explicitly enforced, and defensive programming is preferred.

## Location

Any place where code like this appears (e.g., in GetDataverseUserBySystemUserId):

## Code Issue

```go
Path: "/api/data/v9.2/systemusers(" + systemUserId + ")",
```

## Fix

Use `url.PathEscape` to safely encode path parameters when building URLs with variable user input.

```go
Path: "/api/data/v9.2/systemusers(" + url.PathEscape(systemUserId) + ")",
```

Repeat this fix for all paths that interpolate IDs in a similar manner. This ensures safe, predictable URL formation and protects against malformed requests or path confusion.
