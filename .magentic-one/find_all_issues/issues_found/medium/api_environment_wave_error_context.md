# Title

Error Returned Without Contextual Information

##

/workspaces/terraform-provider-power-platform/internal/services/environment_wave/api_environment_wave.go

## Problem

Errors returned from several places lack proper contextual information to debug effectively. Examples:
  - `return nil, err` statements don't include meaningful error messages explaining the failure.

## Impact

- Lack of context in error messages can make debugging difficult.
- Severity: **Medium**

## Location

Multiple error return points in the methods:

## Code Issues

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &organizations)
if err != nil {
	return nil, err
}

// Similar issue in other methods
```

## Fix

Wrap errors with meaningful context using `fmt.Errorf` or `errors.Wrap`.

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &organizations)
if err != nil {
	return nil, fmt.Errorf("failed to execute GET request to %s: %w", apiUrl.String(), err)
}
```