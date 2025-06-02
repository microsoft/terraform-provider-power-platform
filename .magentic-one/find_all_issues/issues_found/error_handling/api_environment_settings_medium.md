# Title
Non-idiomatic Error Wrapping

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/api_environment_settings.go

## Problem

Throughout the code, errors from `client.Api.Execute` and related functions are returned directly. If additional context is needed, idiomatic error wrapping using `fmt.Errorf` or `errors.Wrap` (from `pkg/errors`, if used), should be considered to provide more context about the error source.

## Impact

Low/Medium. Lacking error context reduces the ease of debugging errors in higher layers.

## Location

Example in `getEnvironment`:

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
if err != nil {
    return nil, err
}
```

## Code Issue

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
if err != nil {
    return nil, err
}
```

## Fix

Wrap the error with additional context:

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
if err != nil {
    return nil, fmt.Errorf("failed to execute API request for environment %s: %w", environmentId, err)
}
```
