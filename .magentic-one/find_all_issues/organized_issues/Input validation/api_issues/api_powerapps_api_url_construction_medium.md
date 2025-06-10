# API URL Construction Not Resilient to Trailing Slashes

##

/workspaces/terraform-provider-power-platform/internal/services/powerapps/api_powerapps.go

## Problem

Manual construction of `apiUrl.Path` using `fmt.Sprintf` risks malformed URLs if input values accidentally contain slashes.

## Impact

Potential for broken URLs, especially if `env.Name` contains unexpected characters. Severity: Medium.

## Location

```go
Path:   fmt.Sprintf("/providers/Microsoft.PowerApps/scopes/admin/environments/%s/apps", env.Name),
```

## Code Issue

```go
Path:   fmt.Sprintf("/providers/Microsoft.PowerApps/scopes/admin/environments/%s/apps", env.Name),
```

## Fix

Validate and sanitize `env.Name` or use path join utilities to assemble URLs safely. Example fix:

```go
Path:   path.Join("/providers/Microsoft.PowerApps/scopes/admin/environments", env.Name, "apps"),
```

Add `"path"` to imports.
