# Title

Missing Use of Interfaces or Potential Abstraction for Various Connector Types

##
`/workspaces/terraform-provider-power-platform/internal/services/connectors/dto.go`

## Problem

Across the file, different connector types (such as `connectorDto`, `virtualConnectorDto`, `unblockableConnectorDto`) have similar fields (`Id`, `Metadata`, etc.), suggesting possible reuse. However, no abstraction through common interfaces or structs has been applied.

## Impact

This results in repetitive code and hampers maintainability and extensibility. Future additions or modifications to the connector types may require changes in multiple places, leading to brittle code. Severity: **medium**

## Location

This issue spans the redundant definitions in the following structures:
- `connectorDto`
- `virtualConnectorDto`
- `unblockableConnectorDto`

## Code Issue

Here are parts of the problematic structure definitions:

```go
type connectorDto struct {
	Name       string                 `json:"name"`
	Id         string                 `json:"id"`
	Type       string                 `json:"type"`
	Properties connectorPropertiesDto `json:"properties"`
}

type unblockableConnectorDto struct {
	Id       string                          `json:"id"`
	Metadata unblockableConnectorMetadataDto `json:"metadata"`
}

type virtualConnectorDto struct {
	Id       string                      `json:"id"`
	Metadata virtualConnectorMetadataDto `json:"metadata"`
}
```

## Fix

Introduce an interface or base struct to encapsulate common fields (`Id`, `Type`, `Metadata`). This will promote reuse and simplify the code. Here's the corrected implementation:

```go
// Define a Connector interface:
type Connector interface {
	GetID() string
}

// Define a base struct for common fields:
type BaseConnector struct {
	Id string `json:"id"`
}

func (bc BaseConnector) GetID() string {
	return bc.Id
}

// Update individual connector structs:
type connectorDto struct {
	BaseConnector
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Properties connectorPropertiesDto `json:"properties"`
}

type unblockableConnectorDto struct {
	BaseConnector
	Metadata unblockableConnectorMetadataDto `json:"metadata"`
}

type virtualConnectorDto struct {
	BaseConnector
	Metadata virtualConnectorMetadataDto `json:"metadata"`
}
```

This approach minimizes redundancy and simplifies adding new connector types in the future.