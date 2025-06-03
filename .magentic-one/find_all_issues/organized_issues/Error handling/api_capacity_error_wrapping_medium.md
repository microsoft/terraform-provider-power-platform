# Issue: Lack of Error Wrapping on API Invocation

## 
/workspaces/terraform-provider-power-platform/internal/services/capacity/api_capacity.go

## Problem

In the `GetTenantCapacity` method, if an error occurs during the execution of the API call, it is returned directly without wrapping or contextualizing. This makes it harder to trace where the error originated from when debugging, especially in larger codebases with many API calls. Proper error wrapping (with `fmt.Errorf("...: %w", err)`) allows for easier and more informative debugging.

## Impact

Severity: **medium**

Directly returning low-level errors without additional context reduces maintainability and makes future debugging and log tracing more cumbersome.

## Location

Line(s):  
```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &dto)
if err != nil {
    return nil, err
}
```

## Code Issue

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &dto)
if err != nil {
    return nil, err
}
```

## Fix

Wrap the error to provide function context:

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &dto)
if err != nil {
    return nil, fmt.Errorf("failed to get tenant capacity: %w", err)
}
```
