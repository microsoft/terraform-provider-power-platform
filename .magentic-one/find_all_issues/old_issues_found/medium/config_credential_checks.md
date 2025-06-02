# Title

Inconsistent credential checks in configuration validation

##

/workspaces/terraform-provider-power-platform/internal/config/config.go

## Problem

The `IsUserManagedIdentityProvided`, `IsSystemManagedIdentityProvided`, and `IsClientSecretCredentialsProvided` methods use inconsistent checks to determine if certain credentials are provided. These inconsistencies could lead to logical errors in determining credential validity, especially for conditions involving overlapping attributes such as `UseMsi` and `ClientId`.

For example:
- The `IsUserManagedIdentityProvided` checks if `UseMsi` is true and `ClientId` is non-empty.
- The `IsSystemManagedIdentityProvided` only checks if `UseMsi` is true but adds that `ClientId` should be empty, relying heavily on the implicit mutual exclusivity of both methods.

## Impact

This inconsistency introduces potential logical errors, ambiguity, and lack of clarity in credential validation, increasing the likelihood of bugs and misconfigurations. This issue is of medium severity since it primarily affects developer productivity and configuration correctness.

## Location

Functions:
- `IsUserManagedIdentityProvided`
- `IsSystemManagedIdentityProvided`
- `IsClientSecretCredentialsProvided`

## Code Issue

```go
func (model *ProviderConfig) IsUserManagedIdentityProvided() bool {
	return model.UseMsi && model.ClientId != ""
}

func (model *ProviderConfig) IsSystemManagedIdentityProvided() bool {
	return model.UseMsi && model.ClientId == ""
}

func (model *ProviderConfig) IsClientSecretCredentialsProvided() bool {
	return model.ClientId != "" && model.ClientSecret != "" && model.TenantId != ""
}
```

## Fix

Refactor these methods to eliminate ambiguity and adhere to consistent validation patterns. For example:

```go
// Checks if Managed Identity is used and if the user-provided Client ID is valid
func (model *ProviderConfig) IsUserManagedIdentityProvided() bool {
	return model.UseMsi && model.ClientId != ""
}

// Checks if Managed Identity is used and if no user-provided Client ID is given (system identity)
func (model *ProviderConfig) IsSystemManagedIdentityProvided() bool {
	return model.UseMsi && model.ClientId == "" // Checks explicitly for system identity
}

// Checks the presence of all components of the client secret credentials
func (model *ProviderConfig) IsClientSecretCredentialsProvided() bool {
	return model.ClientId != "" && model.ClientSecret != "" && model.TenantId != ""
}

// Optional: Add logging or error reporting to detect unintended overlapping conditions
```

This change improves readability and maintains a robust design by clearly separating various identity types and their checks. It should be ensured that the consuming modules respect these expectations.
