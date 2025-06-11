# Body as bytes is set even if `io.ReadAll` fails

##

/workspaces/terraform-provider-power-platform/internal/api/request.go

## Problem

The code reads body without handling `err` before assigning its value. If `err` is not nil, `body` will be meaningless or incomplete, but it will still be assigned to `resp.BodyAsBytes`. Caller might receive an invalid or partial body.

## Impact

Severity: Medium

This could lead to attempts to unmarshal invalid/incomplete data or propagate partial/corrupted information.

## Location

```go
	body, err := io.ReadAll(apiResponse.Body)
	resp.BodyAsBytes = body
```

## Code Issue

```go
	body, err := io.ReadAll(apiResponse.Body)
	resp.BodyAsBytes = body
```

## Fix

Check for `err` from `io.ReadAll` before setting `resp.BodyAsBytes`:

```go
	body, err := io.ReadAll(apiResponse.Body)
	if err != nil {
		return resp, err
	}
	resp.BodyAsBytes = body
```
