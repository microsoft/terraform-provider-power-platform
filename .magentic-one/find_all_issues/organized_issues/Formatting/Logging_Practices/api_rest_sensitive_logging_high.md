# Possible exposure of sensitive response data in logs

##

/workspaces/terraform-provider-power-platform/internal/services/rest/api_rest.go

## Problem

The code logs the entire response body and status at trace level, which, depending on the API being wrapped, could expose sensitive or confidential information. Even if logs are not enabled in production, this poses a risk if trace logging is accidentally switched on or logs are accessed by unauthorized individuals.

## Impact

Severity: High. Logging sensitive data can lead to accidental exposure, data leaks, or violations of data privacy/security policies.

## Location

Lines ~53-56:

## Code Issue

```go
if res != nil && res.HttpResponse != nil {
	tflog.Trace(ctx, fmt.Sprintf("SendOperation Response: %v", res.BodyAsBytes))
	tflog.Trace(ctx, fmt.Sprintf("SendOperation Response Status: %v", res.HttpResponse.Status))
}
```

## Fix

Redact or avoid logging full response bodies. Limit logging to non-sensitive metadata or ensure explicit configuration for safe debug output.

```go
if res != nil && res.HttpResponse != nil {
	tflog.Trace(ctx, fmt.Sprintf("SendOperation Response Status: %v", res.HttpResponse.Status))
	// tflog.Trace(ctx, fmt.Sprintf("SendOperation Response: %v", res.BodyAsBytes)) // Avoid or redact sensitive bodies
}
```
Or, if logging the body is critical, add guards or filters to ensure sensitive fields are removed.
