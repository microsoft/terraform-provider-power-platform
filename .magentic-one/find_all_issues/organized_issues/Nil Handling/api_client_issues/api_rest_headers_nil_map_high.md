# Headers map may be nil

##

/workspaces/terraform-provider-power-platform/internal/services/rest/api_rest.go

## Problem

In `SendOperation`, the `headers` map is only initialized when `operation.Headers` has items. If there are no headers, it remains `nil` and is passed to subsequent functions, which expect a map and iterate over it, potentially leading to a panic when assigning headers in `ExecuteApiRequest`.

## Impact

Potential high-severity runtime panic from operations on a nil map. This can cause hard-to-debug issues and service downtime.

## Location

Lines ~27-48, relevant in both `SendOperation` and `ExecuteApiRequest`:

## Code Issue

```go
	var headers map[string]string
	// ...
	if len(operation.Headers) > 0 {
		headers = make(map[string]string)
		for _, h := range operation.Headers {
			headers[h.Name.ValueString()] = h.Value.ValueString()
		}
	}
    // ...
	for k, v := range headers {
		h.Add(k, v)
	}
```

## Fix

Initialize the `headers` map as an empty map by default, ensuring it is never nil. 

```go
	headers := make(map[string]string)
	if len(operation.Headers) > 0 {
		for _, h := range operation.Headers {
			headers[h.Name.ValueString()] = h.Value.ValueString()
		}
	}
```
