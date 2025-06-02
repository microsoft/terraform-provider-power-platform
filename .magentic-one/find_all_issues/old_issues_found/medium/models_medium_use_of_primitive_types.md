# Issue Found: Use of Primitive Types Instead of Wrapper Types

## Severity
Medium

## File
`/workspaces/terraform-provider-power-platform/internal/services/environment_wave/models.go`

## Problem
Structs such as `FeatureDto` and `OrganizationDto` directly use primitive types (e.g., `string`, `bool`, `int`) for their fields instead of the `types` wrapper types from the Terraform plugin framework. This approach can lead to inconsistencies when interacting with Terraform, as these wrapper types provide more robust handling for null values and other Terraform-specific semantics.

## Impact
This can cause integration issues when parsing or validating data from Terraform configurations. It may also lead to extra boilerplate code to convert between primitives and wrapper types.

## Code Example
An example from `FeatureDto`:

```go
type FeatureDto struct {
	FeatureName      string `json:"featureName"`
	DisplayName      string `json:"displayName"`
	CanBeReset       bool   `json:"canBeReset"`
	Enabled          bool   `json:"enabled"`
	IsAllowed        bool   `json:"isAllowed"`
	NotBefore        string `json:"notBefore"`
	NotAfter         string `json:"notAfter"`
	MinVersion       string `json:"minVersion"`
	MaxVersion       string `json:"maxVersion"`
	State            string `json:"state"`
	AppsUpgradeState string `json:"appsUpgradeState"`
}
```

## Recommendation
Apply Terraform's `types` wrapper types to ensure compatibility and nullability handling. For example:

### Fixed Code Example:
```go
type FeatureDto struct {
	FeatureName      types.String `json:"featureName"`
	DisplayName      types.String `json:"displayName"`
	CanBeReset       types.Bool   `json:"canBeReset"`
	Enabled          types.Bool   `json:"enabled"`
	IsAllowed        types.Bool   `json:"isAllowed"`
	NotBefore        types.String `json:"notBefore"`
	NotAfter         types.String `json:"notAfter"`
	MinVersion       types.String `json:"minVersion"`
	MaxVersion       types.String `json:"maxVersion"`
	State            types.String `json:"state"`
	AppsUpgradeState types.String `json:"appsUpgradeState"`
}
```