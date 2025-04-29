# Title

Unexported struct `locationProperties`

## Path

/workspaces/terraform-provider-power-platform/internal/services/locations/dto.go

## Problem

The struct `locationProperties` is unexported, yet its fields are designed for JSON serialization. This strongly implies a requirement for external use, making it incompatible with being unexported.

## Impact

Leaving `locationProperties` unexported can result in JSON serialization/deserialization issues for its parent structs and objects that rely on it. **Severity: High**

## Location

Line containing `type locationProperties struct`.

## Code Issue

```go
type locationProperties struct {
    DisplayName                            string   `json:"displayName"`
    Code                                   string   `json:"code"`
    IsDefault                              bool     `json:"isDefault"`
    IsDisabled                             bool     `json:"isDisabled"`
    CanProvisionDatabase                   bool     `json:"canProvisionDatabase"`
    CanProvisionCustomerEngagementDatabase bool     `json:"canProvisionCustomerEngagementDatabase"`
    AzureRegions                           []string `json:"azureRegions"`
}
```

## Fix

Export the struct to allow external usage and compatibility with JSON serialization.

```go
type LocationProperties struct {
    DisplayName                            string   `json:"displayName"`
    Code                                   string   `json:"code"`
    IsDefault                              bool     `json:"isDefault"`
    IsDisabled                             bool     `json:"isDisabled"`
    CanProvisionDatabase                   bool     `json:"canProvisionDatabase"`
    CanProvisionCustomerEngagementDatabase bool     `json:"canProvisionCustomerEngagementDatabase"`
    AzureRegions                           []string `json:"azureRegions"`
}
```

Refactor all references throughout the codebase accordingly.
