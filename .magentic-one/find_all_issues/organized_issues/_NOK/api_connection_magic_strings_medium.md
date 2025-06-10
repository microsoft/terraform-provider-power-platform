# Title

Use of Magic Strings for HTTP Method and Status Codes in API Calls

##

/workspaces/terraform-provider-power-platform/internal/services/connection/api_connection.go

## Problem

Throughout the file, HTTP methods (e.g., `"GET"`, `"PUT"`, `"POST"`, `"DELETE"`) and status codes (e.g., `http.StatusOK`, `http.StatusCreated`) are used as hardcoded string literals and slices instead of leveraging enums or constants specific to the project.

## Impact

This practice introduces a risk of typos, makes the code less maintainable, and reduces readability. If the same values are used in multiple places and need updating, it increases the maintenance overhead. Severity: Medium.

## Location

For example, in multiple methods such as:

```go
_, err := client.Api.Execute(ctx, nil, "PUT", apiUrl.String(), nil, connectionToCreate, []int{http.StatusCreated}, &connection)
```

## Code Issue

```go
_, err := client.Api.Execute(ctx, nil, "PUT", apiUrl.String(), nil, connectionToCreate, []int{http.StatusCreated}, &connection)
```

## Fix

Define package-level constants for HTTP methods and status codes, and use them in API calls. Example:

```go
const (
    httpMethodPut    = "PUT"
    httpMethodPost   = "POST"
    httpMethodGet    = "GET"
    httpMethodDelete = "DELETE"
    apiStatusCreated = http.StatusCreated
    apiStatusOK      = http.StatusOK
)
```

Update calls such as:

```go
_, err := client.Api.Execute(ctx, nil, httpMethodPut, apiUrl.String(), nil, connectionToCreate, []int{apiStatusCreated}, &connection)
```

---

*This file should be saved in:*
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/api_connection_magic_strings_medium.md
