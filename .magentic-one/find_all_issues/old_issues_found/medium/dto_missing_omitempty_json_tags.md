# Title

Missing usage of `omitempty` in JSON tags for optional fields

## Path to the file

`/workspaces/terraform-provider-power-platform/internal/services/data_record/dto.go`

## Problem

The current JSON tags for fields do not use the `omitempty` option. This leads to inclusion of fields with zero values (`""` for strings, `nil` for slices) during JSON serialization, which can result in verbose outputs, unnecessary storage overhead, and potential API compatibility issues when zero values are not expected.

## Impact

Including zero-value fields in serialized JSON can create issues when interacting with APIs that only expect fields to be sent if they contain meaningful data. This may lead to failures or unpredictable behavior in integrated systems. **Severity:** Medium

## Location

Affects multiple struct definitions, including:

1. `dataRecordDto`
2. `environmentIdDto`
3. `entityDefinitionsDto`
4. `relationApiResponseDto`
5. `relationApiBodyDto`
6. `attributesApiResponseDto`
7. `attributesApiBodyDto`

## Code Issue

Example of current implementation in `dataRecordDto` struct:

```go
type dataRecordDto struct {
	Id           string `json:"id"`
	OdataContext string `json:"@odata.context"`
	OdataEtag    string `json:"@odata.etag"`
}
```

## Fix

Add `omitempty` to JSON tags for all fields that do not always have values. This ensures only non-zero-value fields are serialized.

```go
type dataRecordDto struct {
	Id           string `json:"id,omitempty"`
	OdataContext string `json:"@odata.context,omitempty"`
	OdataEtag    string `json:"@odata.etag,omitempty"`
}
```

Similarly, update all other structs to include `omitempty` for optional fields.