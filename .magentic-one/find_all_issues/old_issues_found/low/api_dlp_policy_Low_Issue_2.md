# Title

Missing Documentation Comments on Public Methods

## Path

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/api_dlp_policy.go

## Problem

Public methods such as `GetPolicies`, `GetPolicy`, `DeletePolicy`, `UpdatePolicy`, and `CreatePolicy` lack Go-style documentation comments.

## Impact

Missing comments reduce code understandability and adherence to Go standards. Severity marked as **Low** as this does not affect functionality but impacts maintainability.

## Location

All public methods:

- `func (client *client) GetPolicies`
- `func (client *client) GetPolicy`
- `func (client *client) DeletePolicy`
- `func (client *client) UpdatePolicy`
- `func (client *client) CreatePolicy`

## Code Issue

```go
func (client *client) GetPolicies(ctx context.Context) ([]dlpPolicyModelDto, error) {
...
}

func (client *client) GetPolicy(ctx context.Context, name string) (*dlpPolicyModelDto, error) {
...
}
...
```

## Fix

Add proper documentation comments to each method.

```go
// GetPolicies retrieves all DLP policies as a list
func (client *client) GetPolicies(ctx context.Context) ([]dlpPolicyModelDto, error) {
...
}

// GetPolicy retrieves a specific DLP policy by name
func (client *client) GetPolicy(ctx context.Context, name string) (*dlpPolicyModelDto, error) {
...
}
...
```