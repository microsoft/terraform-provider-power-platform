# Unexported DTO Types Restrict Usage Outside the Package

##

/workspaces/terraform-provider-power-platform/internal/services/locations/dto.go

## Problem

The DTO structs (`locationDto`, `locationsArrayDto`, `locationProperties`) are unexported (using lowercase names). DTOs (Data Transfer Objects) are typically used to transfer data between layers/packages, and as such, these should be exported if they need to be used by other packages in your codebase. Keeping them unexported restricts their use to the `locations` package only.

## Impact

Severity: **Medium**

This impacts maintainability and extensibility. If other packages need to interact with location DTOs, they cannot access these types due to their unexported status. It can also create unnecessary wrappers just to access these types elsewhere.

## Location

Top of the file (type declarations).

## Code Issue

```go
type locationDto struct {
	Value []locationsArrayDto `json:"value"`
}

type locationsArrayDto struct {
	ID         string             `json:"id"`
	Type       string             `json:"type"`
	Name       string             `json:"name"`
	Properties locationProperties `json:"properties"`
}

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

Export all DTO types that might be needed by other packages by capitalizing the first letter of the type name.

```go
type LocationDTO struct {
	Value []LocationsArrayDTO `json:"value"`
}

type LocationsArrayDTO struct {
	ID         string             `json:"id"`
	Type       string             `json:"type"`
	Name       string             `json:"name"`
	Properties LocationProperties `json:"properties"`
}

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

