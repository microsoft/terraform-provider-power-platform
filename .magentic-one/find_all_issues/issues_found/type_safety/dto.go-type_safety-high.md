# Potential JSON Marshalling/Unmarshalling Inconsistency Due to Unexported Struct Fields

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/dto.go

## Problem

Field `InstanceURL` in `linkedEnvironmentIdMetadataDto` is not exported (it’s not tagged for JSON, but naming still matters for serialization). While the use of JSON tags is not present, unexported struct fields are not accessible during JSON (un)marshalling. If the intent is to serialize/deserialize this field, it will be skipped.

## Impact

High severity: breaks serialization/deserialization logic if the field is supposed to be used over API boundaries and results in lost or omitted data.

## Location

- Line 18 (`linkedEnvironmentIdMetadataDto`)

## Code Issue

```go
type linkedEnvironmentIdMetadataDto struct {
    InstanceURL string
}
```

## Fix

Export the field, and consider adding JSON struct tags if they are needed.

```go
type LinkedEnvironmentIdMetadataDto struct {
    InstanceURL string `json:"instanceUrl"`
}
```
