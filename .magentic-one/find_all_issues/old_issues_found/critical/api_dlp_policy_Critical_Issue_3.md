# Title

Error Handling for `GetPolicy` lacks context wrapping

## Path

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/api_dlp_policy.go

## Problem

The error returned from `client.Api.Execute` in `GetPolicy` lacks context wrapping, making it harder to debug or understand the source of the error.

## Impact

It impairs error traceability and debugging efficiency, making troubleshooting difficult for both developers and users. Severity is critical since poor error handling affects reliability.

## Location

`func (client *client) GetPolicy(ctx context.Context, name string) (*dlpPolicyModelDto, error)`

## Code Issue

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &policy)
if err != nil {
    var httpError *customerrors.UnexpectedHttpStatusCodeError
    if errors.As(err, &httpError) && httpError.StatusCode == http.StatusNotFound {
        return nil, customerrors.WrapIntoProviderError(err, customerrors.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("DLP Policy '%s' not found", name))
    }
    return nil, err
}
```

## Fix

Wrap the errors using insights/contextual messages to better identify errors.

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &policy)
if err != nil {
    var httpError *customerrors.UnexpectedHttpStatusCodeError
    if errors.As(err, &httpError) && httpError.StatusCode == http.StatusNotFound {
        return nil, customerrors.WrapIntoProviderError(err, customerrors.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("DLP Policy '%s' not found", name))
    }
    return nil, fmt.Errorf("unexpected error while fetching DLP Policy '%s': %w", name, err)
}
```