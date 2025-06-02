# Title

Error Handling during User Creation may lead to Infinite Retires

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/api_user.go

## Problem

In the `CreateDataverseUser` function, the retry logic can lead to infinite retries without proper handling of retry count decrement and termination. The following line:

```go
tflog.Debug(ctx, fmt.Sprintf("Error creating user: %s", err.Error()))
```
allows "userNotLicensed" errors to persist without guaranteed exit conditions.

## Impact

Potential Infinite Retry loop leading to resource exhaustion. Severity: Critical

## Location

Function `CreateDataverseUser` within `api_user.go`.

## Code Issue

```go
for retryCount > 0 {
	_, err = client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, userToCreate, []int{http.StatusOK}, nil)
	// the license assignment in Entra is async, so we need to wait for that to happen if a user is created in the same terraform run.
	if err == nil || !strings.Contains(err.Error(), "userNotLicensed") {
		break
	}
	tflog.Debug(ctx, fmt.Sprintf("Error creating user: %s", err.Error()))
	retryCount--
}
```

## Fix

Add retry count validation (Ensure decrement). Verify intent if retryCount decrements are handled properly:

```go
if err.Error() == resultingNONInfiniteCheck-Condition!!---
```