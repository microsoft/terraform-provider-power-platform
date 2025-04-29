# Title

Missing error wrapping for `client.Api.Execute` in `GetPolicies` method

## Path

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/api_dlp_policy.go

## Problem

In the `GetPolicies` method, error handling for `client.Api.Execute` does not wrap the error into a context message. Not wrapping the error makes debugging difficult and prevents providing context-specific messages when failures occur.

## Impact

The issue makes debugging harder and impacts usability by providing insufficient context when an error occurs. Severity is marked as **Critical** because this affects error handling integrity.

## Location

`func (client *client) GetPolicies(ctx context.Context) ([]dlpPolicyModelDto, error)`

## Code Issue

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &policiesArray)
if err != nil {
    return nil, err
}
```

## Fix

Wrap the error with additional context information using the `fmt.Errorf` or custom error wrapping.

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &policiesArray)
if err != nil {
    return nil, fmt.Errorf("error executing API for policies retrieval: %w", err)
}
```