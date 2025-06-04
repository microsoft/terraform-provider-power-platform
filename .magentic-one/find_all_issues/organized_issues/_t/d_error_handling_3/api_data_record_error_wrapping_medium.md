# Title

Control flow: Missing error wrap/explicit cause in many error returns

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/api_data_record.go

## Problem

Errors from called functions and API responses are often returned directly, without wrapping or providing helpful context for consumers. For example, when `json.Unmarshal` fails, or when a required field is not found, the error is returned as-is. This makes debugging and error tracing harder, especially in a complex service where multiple underlying operations can fail for the same high-level reason.

## Impact

**Severity: Medium**

- Errors in Terraform logs may not provide enough context to understand what failed, especially for unmarshalling, HTTP, or type assertion problems.
- Can slow debugging cycles for both developers and users.

## Location

```go
err = json.Unmarshal(response.BodyAsBytes, &mapResponse)
if err != nil {
    return nil, err
}
```

Also, for most errors returned directly from utility calls, e.g.:
```go
if err != nil {
    return nil, err
}
```

## Code Issue

```go
err = json.Unmarshal(response.BodyAsBytes, &mapResponse)
if err != nil {
    return nil, err
}
```

## Fix

Wrap errors with context using `fmt.Errorf` and `%w` to capture the cause:

```go
err = json.Unmarshal(response.BodyAsBytes, &mapResponse)
if err != nil {
    return nil, fmt.Errorf("unmarshalling response bytes for [operation context]: %w", err)
}
```

Apply this best practice for all returned errors, especially those returned from other helpers, type assertions, or third-party packages (json, url, etc).

---

File:  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/api_data_record_error_wrapping_medium.md`
