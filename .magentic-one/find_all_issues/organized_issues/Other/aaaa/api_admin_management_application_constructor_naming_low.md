# Title

Constructor Function Naming Not Idiomatic

##

/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/api_admin_management_application.go

## Problem


## Impact

Low. This is a stylistic and maintainability issue but does not affect runtime behavior. Conventions matter in Go as they signal intent to other contributors.

## Location

Function declaration near the top:

## Code Issue

```go
func newAdminManagementApplicationClient(clientApi *api.Client) client {
	return client{
		Api: clientApi,
	}
}
```

## Fix

If you intend the constructor to be unexported (for internal use only), that's fine. If it should be exported for wider package use, consider renaming to `NewAdminManagementApplicationClient`. Otherwise, consider a comment explaining the unexported status.

```go
// For unexported constructor (supply clear comment about intention)
func newAdminManagementApplicationClient(clientApi *api.Client) client {
	return client{
		Api: clientApi,
	}
}

// Or, if exported is intended:
func NewAdminManagementApplicationClient(clientApi *api.Client) client {
	return client{
		Api: clientApi,
	}
}
```
