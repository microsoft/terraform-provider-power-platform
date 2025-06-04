# Title

No Error Context Wrapping or Logging for API Calls

##

/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/api_admin_management_application.go

## Problem

When API errors occur, they are returned directly from the method without any log or error wrapping/context. It can make debugging difficult, as the caller will not know which API call or parameters led to the error, particularly important with multiple similar methods.

## Impact

Medium. Affects debuggability and observability.

## Location

Each function that returns an error from `client.Api.Execute`.

## Code Issue

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &adminApp)
return &adminApp, err
```

## Fix

Wrap errors with context, using e.g., `fmt.Errorf` or the `%w` verb, or log them if appropriate.

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &adminApp)
if err != nil {
    return nil, fmt.Errorf("failed to get admin app %s: %w", clientId, err)
}
return &adminApp, nil
```
