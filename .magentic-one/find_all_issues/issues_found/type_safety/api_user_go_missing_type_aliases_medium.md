# Title

Missing Type Aliases/Wrapper Types for Magic String Parameters and Return Types

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/api_user.go

## Problem

Multiple functions pass around and return bare `string` parameters/fields representing different concepts (e.g., environment ID, system user ID, AAD object ID, role ID, domain name, etc.). Using the base string type blurs semantic context and increases risk of mismatched or swapped usage. Go supports custom type aliases/wrappers to provide type safety without extra runtime cost.

## Impact

Severity: Medium

Reducing everything to `string` leads to:
- Poor API self-documentation
- Harder static analysis for bugs
- Increased risk of passing the wrong string to a function (e.g. AAD object ID instead of system user ID)

## Location

Signatures throughout the file, for example:

## Code Issue

```go
func (client *client) GetDataverseUserBySystemUserId(ctx context.Context, environmentId, systemUserId string) (*userDto, error)
```

## Fix

Define meaningful type aliases for major string concepts:

```go
type EnvironmentID string
type SystemUserID string
type AadObjectID string
type RoleID string
```

Update function signatures and struct fields accordingly:

```go
func (client *client) GetDataverseUserBySystemUserId(ctx context.Context, environmentId EnvironmentID, systemUserId SystemUserID) (*userDto, error)
```

This improves compiler checking, documentation, and overall confidence in code correctness.