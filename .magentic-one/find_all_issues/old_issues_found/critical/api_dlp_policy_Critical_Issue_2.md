# Title

Unnecessary variable reinitialization in iteration inside `GetPolicies`

## Path

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/api_dlp_policy.go

## Problem

In the loop within the `GetPolicies` method, the variable `policy` is reinitialized unnecessarily, which affects readability and introduces potential bugs.

## Impact

This may lead to an incorrect representation of data where the same variable name is being redeclared multiple times unnecessarily. Severity is **Critical** because reinitializing the same variable creates confusion and risks overriding another instance within the scope of iteration.

## Location

`func (client *client) GetPolicies(ctx context.Context) ([]dlpPolicyModelDto, error)`

## Code Issue

```go
policy := dlpPolicyDto{}
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &policy)
if err != nil {
    return nil, err
}
```

## Fix

Use unique variable names or refactor the code to avoid variable masking.

```go
var singlePolicy dlpPolicyDto
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &singlePolicy)
if err != nil {
    return nil, err
}
```