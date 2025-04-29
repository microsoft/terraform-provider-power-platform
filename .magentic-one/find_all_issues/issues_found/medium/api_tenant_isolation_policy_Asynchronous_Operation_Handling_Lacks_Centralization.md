# Title

Asynchronous Operation Handling Lacks Centralization

##

`/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/api_tenant_isolation_policy.go`

## Problem

The logic for handling asynchronous operations using the `doWaitForLifecycleOperationStatus` method appears custom to each usage. This introduces duplicate or near-duplicate logic throughout the codebase when similar operations may be handled elsewhere.  

Having multiple implementations of a task (polling async operations) increases maintenance efforts and encourages inconsistencies. Additionally, the method itself should potentially be placed in a more generic package, such as `api`.

## Impact

1. Code duplication wastes resources and causes inconsistencies.
2. Polling behavior must be updated across multiple locations when changes arise, increasing technical debt.
3. Central logic for retries and error handling makes debugging easier.

Severity: **Medium**

## Location

```go
func (client *Client) doWaitForLifecycleOperationStatus(ctx context.Context, response *api.Response) (*TenantIsolationPolicyDto, error) {
    locationHeader := response.GetHeader(constants.HEADER_LOCATION)
    if locationHeader == "" {
        locationHeader = response.GetHeader(constants.HEADER_OPERATION_LOCATION)
    }
    tflog.Debug(ctx, fmt.Sprintf("Location Header: %s", locationHeader))

    if locationHeader == "" {
        return nil, errors.New("no Location or Operation-Location header found in async response")
    }

    for {
        apiResp, err := client.Api.Execute(ctx, nil, "GET", locationHeader, nil, nil, []int{http.StatusOK, http.StatusAccepted}, nil)
        if err != nil {
            return nil, fmt.Errorf("error polling operation status: %v", err)
        }

        tflog.Debug(ctx, fmt.Sprintf("Operation poll status: %s", apiResp.HttpResponse.Status))

        if apiResp.HttpResponse.StatusCode == http.StatusOK {
            return nil, nil
        }

        waitTime := getRetryAfterDuration(apiResp.HttpResponse)
        tflog.Debug(ctx, fmt.Sprintf("Waiting for %s before polling again", waitTime))

        err = client.Api.SleepWithContext(ctx, waitTime)
        if err != nil {
            return nil, fmt.Errorf("polling interrupted: %v", err)
        }
    }
}
```

## Fix

Refactor the `doWaitForLifecycleOperationStatus` function to be a utility function shared across services. Move it into the `api` package where it can serve as a universal helper for polling async operations.

### Example Refactoring 

Place the centralized function in the `api` package:

```go
package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"time"
)

// WaitForAsyncOperation polls an asynchronous operation until completion using the provided context, location headers & retry logic.
func WaitForAsyncOperation(ctx context.Context, apiClient *Client, locationHeaders map[string]string) error {
	locationHeader := locationHeaders[constants.HEADER_LOCATION]
	if locationHeader == "" {
		locationHeader = locationHeaders[constants.HEADER_OPERATION_LOCATION]
	}

	if locationHeader == "" {
		return errors.New("no Location or Operation-Location header found in async response")
	}

	for {
		apiResp, err := apiClient.Execute(ctx, nil, "GET", locationHeader, nil, nil, []int{http.StatusOK, http.StatusAccepted}, nil)
		if err != nil {
			return fmt.Errorf("error polling operation status: %v", err)
		}

		if apiResp.HttpResponse.StatusCode == http.StatusOK {
			return nil
		}

		waitTime := GetRetryAfterDuration(apiResp.HttpResponse)
		time.Sleep(waitTime)
	}
}
```

Finally, remove `doWaitForLifecycleOperationStatus` from `api_tenant_isolation_policy.go` and instead rely on the function exported in the `api` package.