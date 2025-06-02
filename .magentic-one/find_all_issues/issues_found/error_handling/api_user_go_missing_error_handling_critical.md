# Title

Missing Error Handling in Retry Loop for CreateDataverseUser

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/api_user.go

## Problem

In the `CreateDataverseUser` function, the code is retrying API calls up to a certain count, attempting to handle a specific error ("userNotLicensed") via string matching. However, there is no upper-bound or critical alert when all retry attempts are exhausted. If the retries are exhausted (retryCount drops to 0) but the error persists, the final error is only returned, which may lose context regarding the exhaustive nature of the retries and lacks a clear error message for failed user creation after maximum attempts.

## Impact

Severity: Critical

This can lead to silent failures or confusing error experiences for the caller. Users may get a generic API error rather than understanding that all retry attempts have been exhausted. Operators/debuggers may have difficulty diagnosing race conditions, transient failures, or persistent issues with user creation and licensing propagation.

## Location

Within CreateDataverseUser, section for retrying creation on userNotLicensed error:

## Code Issue

```go
	// 9 minutes of retries.
	retryCount := 6 * 9
	var err error

	for retryCount > 0 {
		_, err = client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, userToCreate, []int{http.StatusOK}, nil)
		// the license assignment in Entra is async, so we need to wait for that to happen if a user is created in the same terraform run.
		if err == nil || !strings.Contains(err.Error(), "userNotLicensed") {
			break
		}
		tflog.Debug(ctx, fmt.Sprintf("Error creating user: %s", err.Error()))
		err = client.Api.SleepWithContext(ctx, 10*time.Second)
		if err != nil {
			return nil, err
		}

		retryCount--
	}
	if err != nil {
		return nil, err
	}
```

## Fix

Log an explicit error or wrap the error with additional context when retries are exhausted. Example:

```go
// 9 minutes of retries.
retryCount := 6 * 9
var err error

for retryCount > 0 {
	_, err = client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, userToCreate, []int{http.StatusOK}, nil)
	if err == nil || !strings.Contains(err.Error(), "userNotLicensed") {
		break
	}
	tflog.Debug(ctx, fmt.Sprintf("Error creating user: %s", err.Error()))
	err = client.Api.SleepWithContext(ctx, 10*time.Second)
	if err != nil {
		return nil, err
	}

	retryCount--
}
if err != nil {
	if retryCount == 0 {
		return nil, fmt.Errorf("failed to create Dataverse user after maximum retries: %w", err)
	}
	return nil, err
}
```

This provides a more descriptive error flow for consumers and makes debugging easier.