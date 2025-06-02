# Title

Missing Documentation for Public Methods

##

/workspaces/terraform-provider-power-platform/internal/services/environment_wave/api_environment_wave.go

## Problem

All public methods in the `environmentWaveClient` struct (`GetOrgEnvironmentId`, `UpdateFeature`, `GetFeature`) are missing comments or documentation.

## Impact

- This makes the code less readable and less maintainable.
- Developers unfamiliar with the code will find it harder to understand the purpose and behavior of these methods.
- Severity: **Low**

## Location

Methods inside the `environmentWaveClient` struct.

## Code Issues

```go
func (client *environmentWaveClient) GetOrgEnvironmentId(ctx context.Context, environmentId string) (*OrganizationDto, error)
// Similar lack of comments for UpdateFeature() and GetFeature() methods
```

## Fix

Add documentation comments for each method following the Go conventions.

```go
// GetOrgEnvironmentId retrieves the organization ID associated with the provided environment.
// It uses the environment ID to fetch details and find an organization that matches the environment.
// Returns an OrganizationDto or an error if the environment ID is invalid.
func (client *environmentWaveClient) GetOrgEnvironmentId(ctx context.Context, environmentId string) (*OrganizationDto, error) {
	// Function body...
}

// Provide similar comments for UpdateFeature() and GetFeature().
```