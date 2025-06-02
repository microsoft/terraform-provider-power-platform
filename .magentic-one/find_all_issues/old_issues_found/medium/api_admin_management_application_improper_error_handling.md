# Title

Improper Error Handling in API Execution

##

`/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/api_admin_management_application.go`

## Problem

In the following methods:
- `GetAdminApplication`
- `RegisterAdminApplication`
- `UnregisterAdminApplication`

The returned errors from `client.Api.Execute` calls are directly passed to the calling function without additional context. This makes it difficult for callers to determine the root cause of the error, as the error lacks context about which operation or API failed.

## Impact

- **Severity:** **Medium**
- Reduced debuggability and maintainability of the code.
- Makes tracing issues logged by the client more complex in case of multiple failures.

## Location

- Method `GetAdminApplication`: Line 21-25
- Method `RegisterAdminApplication`: Line 31-35
- Method `UnregisterAdminApplication`: Line 45-49

## Code Issue

```go
func (client *client) GetAdminApplication(ctx context.Context, clientId string) (*adminManagementApplicationDto, error) {
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &adminApp)

	return &adminApp, err // Error passed directly without context
}

func (client *client) RegisterAdminApplication(ctx context.Context, clientId string) (*adminManagementApplicationDto, error) {
	_, err := client.Api.Execute(ctx, nil, "PUT", apiUrl.String(), nil, nil, []int{http.StatusOK}, &adminApp)

	return &adminApp, err // Error passed directly without context
}

func (client *client) UnregisterAdminApplication(ctx context.Context, clientId string) error {
	_, err := client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNoContent}, nil)

	return err // Error passed directly without context
}
```

## Fix

Wrap the returned error with additional context that specifies the operation being performed and the corresponding API endpoint. This helps significantly in debugging and log tracing during runtime.

```go
func (client *client) GetAdminApplication(ctx context.Context, clientId string) (*adminManagementApplicationDto, error) {
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &adminApp)

	if err != nil {
		// Wrap error with additional context
		return nil, fmt.Errorf("failed to execute GET request for admin application: %v, error: %w", apiUrl.String(), err)
	}

	return &adminApp, nil
}

func (client *client) RegisterAdminApplication(ctx context.Context, clientId string) (*adminManagementApplicationDto, error) {
	_, err := client.Api.Execute(ctx, nil, "PUT", apiUrl.String(), nil, nil, []int{http.StatusOK}, &adminApp)

	if err != nil {
		// Wrap error with additional context
		return nil, fmt.Errorf("failed to execute PUT request for admin application registration: %v, error: %w", apiUrl.String(), err)
	}

	return &adminApp, nil
}

func (client *client) UnregisterAdminApplication(ctx context.Context, clientId string) error {
	_, err := client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNoContent}, nil)

	if err != nil {
		// Wrap error with additional context
		return fmt.Errorf("failed to execute DELETE request for admin application unregistration: %v, error: %w", apiUrl.String(), err)
	}

	return nil
}
```