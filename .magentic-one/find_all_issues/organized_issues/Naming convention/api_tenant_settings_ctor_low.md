# Unexported Constructor Naming Inconsistency

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/api_tenant_settings.go

## Problem

The constructor function `newTenantSettingsClient` uses camel case, which is fine for unexported functions, but for consistency and clarity, Go typically uses the form `new<Type>` or, for exported functions, `New<Type>`. If the struct is made public (`Client`), the constructor should be renamed accordingly (e.g., `NewClient`). This helps maintain uniform naming standards and clarity across the codebase.

## Impact

- **Severity: Low**
- Minor readability and consistency impact.
- Slightly increases friction for future maintainers or external contributors.
- Could hinder the use of Go’s standard conventions if the API is later exported.

## Location

```go
func newTenantSettingsClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}
```

## Code Issue

```go
func newTenantSettingsClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}
```

## Fix

If intended to be exported, rename to `NewClient` and update all usages. If internal, consider simplifying to `newClient` or matching the struct name if exported.

```go
func NewClient(apiClient *api.Client) *Client {
	return &Client{
		Api: apiClient,
	}
}
```
If keeping `client` unexported, prefer `newClient`. If exporting, use `NewClient` for maximal clarity and Go idiomatic usage.
