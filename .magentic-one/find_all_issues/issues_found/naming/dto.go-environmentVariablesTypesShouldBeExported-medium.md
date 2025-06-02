# Title

Environment Variables DTO Types Should Be Exported

##

/workspaces/terraform-provider-power-platform/internal/services/solution/dto.go

## Problem

The types `settingsEnvironmentVariableDto` and `settingsConnectionReferencesDto` are unexported because their names begin with a lowercase letter. If these are used outside the package (which is likely for DTOs), they should be exported.

## Impact

These types cannot be used from other packages, limiting code reusability and API clarity. Severity: medium.

## Location

Lines 11–20

## Code Issue

```go
type settingsEnvironmentVariableDto struct {
	SchemaName string `json:"schemaname"`
	Value      string `json:"value"`
}

type settingsConnectionReferencesDto struct {
	LogicalName  string `json:"logicalname"`
	ConnectionId string `json:"connectionid"`
	ConnectorId  string `json:"connectorid"`
}
```

## Fix

Capitalize the first letter to export the types:

```go
type SettingsEnvironmentVariableDto struct {
	SchemaName string `json:"schemaname"`
	Value      string `json:"value"`
}

type SettingsConnectionReferencesDto struct {
	LogicalName  string `json:"logicalname"`
	ConnectionId string `json:"connectionid"`
	ConnectorId  string `json:"connectorid"`
}
```
