# Title

Improper Error Handling in DeleteEnvironmentGroup

##

`/workspaces/terraform-provider-power-platform/internal/services/environment_groups/api_environment_group.go`

## Problem

In the `DeleteEnvironmentGroup` function, the error handling design assumes certain conditions for the HTTP response body and status code. If the response body is empty, it defaults to a generic hardcoded error message. Additionally, the conditions for `strings.Contains` check whether specific error scenarios occurred but fail to handle unknown errors effectively.

## Impact

- **Severity:** High  
- This might cause confusion when debugging the issue as generic error messages are provided for unknown cases. Incomplete or uninformative error handling could lead to a poor developer experience and reliability issues in production.

## Location

### `DeleteEnvironmentGroup` Method

```go
	if len(resp.BodyAsBytes) == 0 {
		return errors.New("failed to delete environment group")
	}

	body := string(resp.BodyAsBytes[:])
	if strings.Contains(body, "EnvironmentsInEnvironmentGroup") {
		return customerrors.WrapIntoProviderError(err, customerrors.ERROR_ENVIRONMENTS_IN_ENV_GROUP, "Failed to delete environment group because it contains environments")
	} else if strings.Contains(body, "PolicyAssignedToEnvironmentGroup") {
		return customerrors.WrapIntoProviderError(err, customerrors.ERROR_POLICY_ASSIGNED_TO_ENV_GROUP, "Failed to delete environment group because it has a policy assigned")
	}
	return errors.New(body)
```

## Code Issue

```go
	if len(resp.BodyAsBytes) == 0 {
		return errors.New("failed to delete environment group")
	}

	body := string(resp.BodyAsBytes[:])
	if strings.Contains(body, "EnvironmentsInEnvironmentGroup") {
		return customerrors.WrapIntoProviderError(err, customerrors.ERROR_ENVIRONMENTS_IN_ENV_GROUP, "Failed to delete environment group because it contains environments")
	} else if strings.Contains(body, "PolicyAssignedToEnvironmentGroup") {
		return customerrors.WrapIntoProviderError(err, customerrors.ERROR_POLICY_ASSIGNED_TO_ENV_GROUP, "Failed to delete environment group because it has a policy assigned")
	}
	return errors.New(body)
```

## Fix

Refactor the error handling to consider unknown scenarios and provide informative error outputs.

```go
	if len(resp.BodyAsBytes) == 0 {
		return customerrors.WrapIntoProviderError(err, customerrors.ERROR_UNKNOWN, "Response body is empty while deleting environment group")
	}

	body := string(resp.BodyAsBytes[:])
	switch {
	case strings.Contains(body, "EnvironmentsInEnvironmentGroup"):
		return customerrors.WrapIntoProviderError(err, customerrors.ERROR_ENVIRONMENTS_IN_ENV_GROUP, "Failed to delete environment group because it contains environments")
	case strings.Contains(body, "PolicyAssignedToEnvironmentGroup"):
		return customerrors.WrapIntoProviderError(err, customerrors.ERROR_POLICY_ASSIGNED_TO_ENV_GROUP, "Failed to delete environment group because it has a policy assigned")
	default:
		return customerrors.WrapIntoProviderError(err, customerrors.ERROR_UNKNOWN, body)
	}
```

- Using a `switch` statement improves readability.
- Added a default case to handle unknown error messages properly.
