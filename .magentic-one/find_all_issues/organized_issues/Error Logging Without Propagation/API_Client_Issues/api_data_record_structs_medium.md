# Title

Data consistency: Weak type-specific Dto definitions (use of map[string]any)

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/api_data_record.go

## Problem

In multiple places, function or method signatures and logic use `map[string]any` to represent structured data coming from API responses (OData, attributes, environments, etc). This reduces type-safety, means that field names are not checked at compile-time, and makes code harder to maintain.

Examples:
- `GetDataRecordsByODataQuery` returns `(*ODataQueryResponse, error)` but the records themselves are `map[string]any`.
- Various helpers for relations/attributes use generic maps, not typed structs.

## Impact

**Severity: Medium**

- Compiler cannot catch misnamed fields or wrong types, risking silent runtime errors.
- Maintainers and IDEs cannot provide autocomplete, refactoring, or doc hints.
- Testing is harder and bug-prone as mocks must construct arbitrary maps.
- If APIs change, breakage won't be immediately visible.

## Location

```go
func (client *client) GetDataRecord(ctx context.Context, recordId, environmentId, tableName string) (map[string]any, error) {
	// ...
	result := make(map[string]any, 0)
	// ...
}
```

## Code Issue

```go
func (client *client) GetDataRecord(ctx context.Context, recordId, environmentId, tableName string) (map[string]any, error)
```

## Fix

Define properly typed Dto structs (with explicit fields rather than dynamic maps):

```go
type DataRecord struct {
    // Example fields, fill as per actual API
    ID   string `json:\"id\"`
    Name string `json:\"name\"`
    // ...
}
```

Then update function signatures:

```go
func (client *Client) GetDataRecord(ctx context.Context, recordId, environmentId, tableName string) (*DataRecord, error)
```

And parse responses using `json.Unmarshal` directly into these types.

---

File:  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/type_safety/api_data_record_structs_medium.md`
