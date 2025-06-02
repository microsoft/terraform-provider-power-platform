# Issue Found: Unnecessary JSON Tag Names That Match Field Names

## Severity
Low

## File
`/workspaces/terraform-provider-power-platform/internal/services/environment_wave/models.go`

## Problem
Several structs (e.g., `FeatureDto`, `OrganizationDto`) include JSON tags that match the field names exactly (e.g., `FeatureName` tagged as `json:"featureName"`). This is redundant because Go will encode/decode such fields with identical names by default.

## Impact
This increases code verbosity unnecessarily, making the struct definition longer without adding functional benefits.

## Code Example
Here's an example of the redundant code:

```go
type FeatureDto struct {
	FeatureName      string `json:"featureName"` // Redundant JSON tag
	DisplayName      string `json:"displayName"` // Redundant JSON tag
	// Other fields...
}
```

## Recommendation
Remove the JSON tags where the field name exactly matches the tag name. For fields requiring customization, retain the JSON tag.

### Fixed Code:
```go
type FeatureDto struct {
	FeatureName      string // Use default JSON encoding
	DisplayName      string // Use default JSON encoding
	// Other fields...
}
```
