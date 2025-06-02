# Title

Ambiguity in `Type` Field: Lack of Consistent Enum Usage and Validation

##
`/workspaces/terraform-provider-power-platform/internal/services/connectors/dto.go`

## Problem

The `Type` field in multiple structs (e.g., `connectorDto` and `virtualConnectorMetadataDto`) lacks validation and constraints, resulting in potential ambiguity. For example, `Type` fields are represented as plain strings without enum definition, making it error-prone and difficult to maintain or validate across the application.

## Impact

The absence of strong constraints for `Type` can lead to inconsistent or incorrect data being stored or sent via serialization. As this field is likely critical for defining the behavior or categorization of connectors, its mismanagement can result in data corruption or runtime errors. Severity: **high**

## Location

- `Type` field in `connectorDto`
- `Type` field in `virtualConnectorMetadataDto`.

## Code Issue

Here are the problematic usages of `Type`:

```go
type connectorDto struct {
	Name       string                 `json:"name"`
	Id         string                 `json:"id"`
	Type       string                 `json:"type"`
	Properties connectorPropertiesDto `json:"properties"`
}

// Another usage:
type virtualConnectorMetadataDto struct {
	VirtualConnector bool   `json:"virtualConnector"`
	Name             string `json:"name"`
	Type             string `json:"type"` // Reused with likely no validation
	DisplayName      string `json:"displayName"`
}
```

## Fix

To resolve this issue, define an enumeration for valid types and ensure the `Type` field is constrained to these enumerations during validation. Here's the fixed implementation:

```go
// Define enum for valid 'Type' values:
type ConnectorType string

const (
	StandardConnector ConnectorType = "Standard"
	VirtualConnector  ConnectorType = "Virtual"
)

// Use the enum in the structs:
type connectorDto struct {
	Name       string                 `json:"name"`
	Id         string                 `json:"id"`
	Type       ConnectorType          `json:"type"` 
	Properties connectorPropertiesDto `json:"properties"`
}

type virtualConnectorMetadataDto struct {
	VirtualConnector bool          `json:"virtualConnector"`
	Name             string        `json:"name"`
	Type             ConnectorType `json:"type"`
	DisplayName      string        `json:"displayName"`
}
```

This ensures strict type checking, reduces ambiguity, and prevents invalid data from being processed.