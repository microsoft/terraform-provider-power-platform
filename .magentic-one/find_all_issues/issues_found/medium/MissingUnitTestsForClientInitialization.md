# Title

Missing Unit Test Coverage for `newPowerAppssClient`

##

`/workspaces/terraform-provider-power-platform/internal/services/powerapps/api_powerapps.go`

## Problem

The `newPowerAppssClient` function is a basic factory function for creating a new `client` instance, yet there is no indication of unit tests or validation coverage for this function to ensure it initializes fields correctly.

## Impact

- Severity: **Medium**
- Risk: Bugs in initialization may go unnoticed, leading to runtime errors that could be easily caught during testing.
- Prevents early detection of misconfigured or nil values for `Api` and `environmentClient`.

## Location

Line: Function `newPowerAppssClient`

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

Write a unit test to validate proper initialization of the `client` structure, ensuring `Api` and `environmentClient` fields are set correctly.

```go
func TestNewPowerAppsClient(t *testing.T) {
    mockApiClient := &api.Client{} // Mock or stub API client
    result := newPowerAppssClient(mockApiClient)
    
    if result.Api != mockApiClient {
        t.Errorf("expected Api client to be %v, got %v", mockApiClient, result.Api)
    }
    
    if result.environmentClient == nil {
        t.Errorf("expected environmentClient to not be nil")
    }
}
```