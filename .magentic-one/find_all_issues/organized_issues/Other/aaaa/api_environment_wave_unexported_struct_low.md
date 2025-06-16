# Unexported Struct `environmentWaveClient` Used Across Packages

##

/workspaces/terraform-provider-power-platform/internal/services/environment_wave/api_environment_wave.go

## Problem

The struct `environmentWaveClient` (and its corresponding constructor `newEnvironmentWaveClient`) are unexported (start with a lowercase letter), but the file appears to be part of a larger internal API. If this client is intended to be used outside this package, it should be exported.

## Impact

Limits usability and discoverability in larger codebases. Severity: **low**

## Location

Declaration:

```go
type environmentWaveClient struct {
```

Constructor:

```go
func newEnvironmentWaveClient(apiClient *api.Client) *environmentWaveClient {
```

## Code Issue

```go
type environmentWaveClient struct {
```

## Fix

Export the struct and constructor if they are meant for use outside the package:

```go
type EnvironmentWaveClient struct {
    ...
}

func NewEnvironmentWaveClient(apiClient *api.Client) *EnvironmentWaveClient {
    ...
}
```

If not intended for export, no action is needed.
