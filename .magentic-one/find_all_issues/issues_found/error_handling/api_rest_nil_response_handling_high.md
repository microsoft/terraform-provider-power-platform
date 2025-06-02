# Insufficient error handling and control flow for nil response

##

/workspaces/terraform-provider-power-platform/internal/services/rest/api_rest.go

## Problem

In `SendOperation`, there's a conditional block handling the case when `res == nil && err != nil`, but the general logic outside this block assumes that `res.BodyAsBytes` is always safe to access, which could result in a runtime panic if `res` is nil.

## Impact

Potential high-severity runtime panic due to dereferencing a nil pointer, impacting server reliability and correctness.

## Location

Lines ~60-68:

## Code Issue

```go
	if res == nil && err != nil {
		output["body"] = types.StringValue(err.Error())
	} else {
		if len(res.BodyAsBytes) > 0 {
			output["body"] = types.StringValue(string(res.BodyAsBytes))
		}
	}
```

## Fix

Check for `res` being non-nil before accessing its members. This prevents runtime panics and makes the code more robust.

```go
	if res == nil && err != nil {
		output["body"] = types.StringValue(err.Error())
	} else if res != nil && len(res.BodyAsBytes) > 0 {
		output["body"] = types.StringValue(string(res.BodyAsBytes))
	}
```
