# Inconsistent JSON Struct Tag Naming

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/dto.go

## Problem

JSON struct tag names are inconsistently cased. Some are camelCase ("LogicalName"), others are PascalCase ("MetadataId"), and others "@odata.context". Conventionally, JSON properties should follow lowerCamelCase for external APIs.

## Impact

Medium severity: Inconsistent JSON output can confuse users/consumers of the API, is error-prone, and is counter to established best practices.

## Location

- All struct field tags in this file

## Code Issue

```go
type entityDefinitionsDto struct {
	OdataContext          string `json:"@odata.context"`
	PrimaryIDAttribute    string `json:"PrimaryIdAttribute"`
	LogicalCollectionName string `json:"LogicalCollectionName"`
	MetadataID            string `json:"MetadataId"`
}
// ...etc
```

## Fix

Decide on a consistent naming convention (typically lowerCamelCase for JSON) and update tags accordingly.

```go
type EntityDefinitionsDto struct {
	ODataContext          string `json:"@odata.context"`
	PrimaryIdAttribute    string `json:"primaryIdAttribute"`
	LogicalCollectionName string `json:"logicalCollectionName"`
	MetadataId            string `json:"metadataId"`
}
```
