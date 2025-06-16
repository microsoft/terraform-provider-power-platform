# Unclear Type Names (DTO Suffix Unnecessary)

##

/workspaces/terraform-provider-power-platform/internal/services/environment_wave/api_environment_wave.go

## Problem

Types like `OrganizationDto`, `OrganizationsArrayDto`, `FeatureDto`, and `FeaturesArrayDto` use the "Dto" suffix, which is not idiomatic in Go. It's better to have clear, concise type names, such as `Organization`, `Organizations`, `Feature`, and `Features`.

## Impact

Reduces code readability and clarity for Go developers. Severity: **low**

## Location

Throughout the file wherever these types are used.

## Code Issue

```go
organizations := OrganizationsArrayDto{}
...
func (client *environmentWaveClient) GetOrgEnvironmentId(ctx context.Context, environmentId string) (*OrganizationDto, error)
```

## Fix

Rename types to more idiomatic Go names:

```go
organizations := Organizations{}
...
func (client *environmentWaveClient) GetOrgEnvironmentId(ctx context.Context, environmentId string) (*Organization, error)
```

Make sure to update the type definitions accordingly.
