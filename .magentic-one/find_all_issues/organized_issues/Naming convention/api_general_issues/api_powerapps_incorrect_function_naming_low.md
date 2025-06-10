# Incorrect Function Naming Convention

##

/workspaces/terraform-provider-power-platform/internal/services/powerapps/api_powerapps.go

## Problem

The function is named `newPowerAppssClient`, which is inconsistent with Go naming conventions and likely a typo (double 's' in "Appss").

## Impact

This can cause confusion for maintainers and may introduce subtle bugs when this function is called elsewhere. Severity: Low.

## Location

Line 14 in the file.

## Code Issue

```go
func newPowerAppssClient(apiClient *api.Client) client {
	return client{
		Api:               apiClient,
		environmentClient: environment.NewEnvironmentClient(apiClient),
	}
}
```

## Fix

Correct the function name to use the singular "PowerApps":

```go
func newPowerAppsClient(apiClient *api.Client) client {
	return client{
		Api:               apiClient,
		environmentClient: environment.NewEnvironmentClient(apiClient),
	}
}
```
