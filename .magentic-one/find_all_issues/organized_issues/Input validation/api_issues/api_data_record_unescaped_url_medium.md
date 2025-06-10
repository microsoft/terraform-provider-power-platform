# Title

Possible HTTP API misuse: Unchecked/unvalidated URL construction

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/api_data_record.go

## Problem

Several API calls use `fmt.Sprintf` and string concatenation to construct URLs that then get passed to `url.URL` or directly sent in requests, e.g.,

```go
apiUrl := fmt.Sprintf("https://%s/api/data/%s/%s", environmentHost, constants.DATAVERSE_API_VERSION, query)
```

or

```go
apiPath := fmt.Sprintf("/api/data/%s/%s", constants.DATAVERSE_API_VERSION, entityDefinition.LogicalCollectionName)
```

These constructs do not validate or escape URL path and query parameters, which could result in malformed requests, injection attacks, or undefined errors if input variables (like `environmentHost`, `query`, `tableName`, or `recordId`) contain special URL characters.

## Impact

**Severity: Medium**

- If input variables are tainted (possibly from external sources), this can be an injection vulnerability.
- If variables contain reserved URL characters, API requests may break or behave unexpectedly.
- Can lead to difficult-to-diagnose bugs when requests return 404s or fail randomly.

## Location

Examples:

```go
apiUrl := fmt.Sprintf("https://%s/api/data/%s/%s", environmentHost, constants.DATAVERSE_API_VERSION, query)

Path:   fmt.Sprintf("/api/data/%s/%s(%s)", constants.DATAVERSE_API_VERSION, tableEntityDefinition.LogicalCollectionName, recordId),
```

## Code Issue

```go
apiUrl := fmt.Sprintf("https://%s/api/data/%s/%s", environmentHost, constants.DATAVERSE_API_VERSION, query)
```

## Fix

Always use `url.URL` for path assembly and `url.PathEscape` for dynamic path segments:

```go
apiUrl := &url.URL{
    Scheme: "https",
    Host:   environmentHost,
    Path:   fmt.Sprintf("/api/data/%s/%s", constants.DATAVERSE_API_VERSION, url.PathEscape(query)),
}
```

Where query parameters, logical names, or IDs are included, use `url.PathEscape(variable)` or `url.QueryEscape(variable)` as appropriate for path or query context.

Review all API URL construction points and wrap dynamic segments with `url.PathEscape` to prevent malformed URLs and potential security issues.

---

File:  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/api_client/api_data_record_unescaped_url_medium.md`
