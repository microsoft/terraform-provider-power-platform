# Possible nil dereference in Response.GetHeader

##

/workspaces/terraform-provider-power-platform/internal/api/request.go

## Problem

If `apiResponse.HttpResponse` is nil (e.g., on error responses), calling methods on it will panic.

## Impact

Severity: High

This can easily result in runtime panics if callers do not check for nil Response before calling `GetHeader`.

## Location

```go
func (apiResponse *Response) GetHeader(name string) string {
	return apiResponse.HttpResponse.Header.Get(name)
}
```

## Code Issue

```go
	return apiResponse.HttpResponse.Header.Get(name)
```

## Fix

Check for nil `HttpResponse` before accessing its fields:

```go
func (apiResponse *Response) GetHeader(name string) string {
	if apiResponse.HttpResponse == nil {
		return ""
	}
	return apiResponse.HttpResponse.Header.Get(name)
}
```
