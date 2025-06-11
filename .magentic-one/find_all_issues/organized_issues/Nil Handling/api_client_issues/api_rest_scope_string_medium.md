# Counterintuitive parameter design for `scope`

##

/workspaces/terraform-provider-power-platform/internal/services/rest/api_rest.go

## Problem

The `scope` parameter for `ExecuteApiRequest` is a pointer to a string but must not be nil (function returns an error otherwise). This creates unnecessary complexity and forces the caller into handling pointer logic, with little gain since nil is an error.

## Impact

Severity: Medium. Promotes confusing API and parameter handling, possibly propagating nil pointer patterns unnecessarily throughout the code base.

## Location

Within `ExecuteApiRequest`:

## Code Issue

```go
	if scope == nil {
		return nil, errors.New("invalid input: scope must be provided")
	}

	return client.Api.Execute(ctx, []string{*scope}, method, url, h, body, expectedStatusCodes, nil)
```

## Fix

Accept `scope` as a string (not a pointer), and enforce presence at compile-time via type signature:

```go
func (client *Client) ExecuteApiRequest(ctx context.Context, scope string, url, method string, body *string, headers map[string]string, expectedStatusCodes []int) (*api.Response, error) {
	h := http.Header{}
	for k, v := range headers {
		h.Add(k, v)
	}
	return client.Api.Execute(ctx, []string{scope}, method, url, h, body, expectedStatusCodes, nil)
}
```
You may add a call-time check and documentation, but using pointer is unnecessary if nil is always an error.
