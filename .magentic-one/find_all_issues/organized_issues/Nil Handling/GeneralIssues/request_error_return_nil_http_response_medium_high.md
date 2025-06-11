# Error return combined with possibly nil http.Response

##

/workspaces/terraform-provider-power-platform/internal/api/request.go

## Problem

In `doRequest`, if an error is returned by `httpClient.Do`, a `Response` containing a possibly nil `apiResponse` is returned. Later code (e.g., the caller) could try to access fields on `resp.HttpResponse` without checking for nil, causing a panic.

## Impact

Severity: Medium/High

This could result in runtime panics if not handled, reducing the stability of the codebase.

## Location

```go
	apiResponse, err := httpClient.Do(request)
	resp := &Response{
		HttpResponse: apiResponse,
	}

	if err != nil {
		return resp, err
	}

	if apiResponse == nil {
		return resp, errors.New("unexpected nil response without error")
	}
```

## Code Issue

```go
	resp := &Response{
		HttpResponse: apiResponse,
	}

	if err != nil {
		return resp, err
	}
```

## Fix

Only return a non-nil Response if apiResponse is non-nil, otherwise return nil for Response:

```go
	apiResponse, err := httpClient.Do(request)

	if err != nil {
		return nil, err
	}
	if apiResponse == nil {
		return nil, errors.New("unexpected nil response without error")
	}
	resp := &Response{
		HttpResponse: apiResponse,
	}
```
