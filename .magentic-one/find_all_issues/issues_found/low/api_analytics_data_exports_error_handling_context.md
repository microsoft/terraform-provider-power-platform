# Issue: Error Handling Without Adequate Context

## Problem
Error Handling in certain portions of the file lacks enough context to understand what operation failed.

## Impact
Debugging and troubleshooting become harder, requiring developers to manually trace failures. Low severity but poor practice.

## Affected Code
```go
_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &adr)
if err != nil {
	return nil, fmt.Errorf("failed to get analytics data export: %w", err)
}
```

## Recommended Fix
Include operation-specific details in the error context message for better debugging ability.

### Example Solution:
```go
_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &adr)
if err != nil {
	return nil, fmt.Errorf("failed to send GET request to %s for analytics data export: %w", apiUrl.String(), err)
}
```