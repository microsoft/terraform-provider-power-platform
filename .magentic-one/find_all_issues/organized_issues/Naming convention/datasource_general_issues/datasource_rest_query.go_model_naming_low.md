# Title

Ambiguous Type Name: `DataverseWebApiDatasourceModel`

##

/workspaces/terraform-provider-power-platform/internal/services/rest/datasource_rest_query.go

## Problem

The `DataverseWebApiDatasourceModel` type is referenced in the code (in the `Read` function) but its declaration is not included in the provided file. Assuming it is declared elsewhere or imported, the name is verbose and could be misaligned with Go's naming conventions or with its usage context. Names should be descriptive but concise and fit their purpose (state or schema model, etc.).

## Impact

Ambiguous or overly verbose names slow understanding and burden maintenance, especially when interleaved with other "model" types. This is a low-severity issue, mostly about maintainability and readability.

## Location

```go
var state DataverseWebApiDatasourceModel
```

## Code Issue

```go
var state DataverseWebApiDatasourceModel
```

## Fix

Rename to a more concise and context-specific name, such as:

```go
var state DataverseWebAPIState
```

Or, if the model is for schema state, prefix/suffix accordingly. Also ensure consistency with any type aliases or field names elsewhere in the provider.

---

This issue will be saved in:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/datasource_rest_query.go_model_naming_low.md`
