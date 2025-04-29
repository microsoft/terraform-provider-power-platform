# Title

Missing context cancellation in `DoWaitForLifecycleOperationStatus`.

##

/workspaces/terraform-provider-power-platform/internal/api/lifecycle.go

## Problem

The `DoWaitForLifecycleOperationStatus` function does not handle the cancellation of the provided context (`ctx`). This omission can lead to hanging operations if the parent context is canceled or exceeds its timeout.

## Impact

Failure to respect context cancellation can result in resource leaks, hanging requests, or inefficient cancellation handling. It is particularly problematic in long-running operations like waiting for lifecycle status, where the operation might continue even though the caller has canceled the context. Severity: **Critical**

## Location

The issue is found in the `DoWaitForLifecycleOperationStatus` function:

## Code Issue

```go
func (client *Client) DoWaitForLifecycleOperationStatus(ctx context.Context, response *Response) (*LifecycleDto, error) {
    locationHeader := response.GetHeader(constants.HEADER_LOCATION)
    if locationHeader == "" {
        locationHeader = response.GetHeader(constants.HEADER_OPERATION_LOCATION)
    }
    tflog.Debug(ctx, "Location Header: "+locationHeader)

    _, err := url.Parse(locationHeader)
    if err != nil {
        tflog.Error(ctx, "Error parsing location header: "+err.Error())
    }

    waitFor := retryAfter(ctx, response.HttpResponse)

    for {
        lifecycleResponse := LifecycleDto{}
        response, err = client.Execute(ctx, nil, "GET", locationHeader, nil, nil, []int{http.StatusOK}, &lifecycleResponse)
        if err != nil {
            return nil, err
        }

        err = client.SleepWithContext(ctx, waitFor)
        if err != nil {
            return nil, err
        }

        tflog.Debug(ctx, "Lifecycle Operation State: '"+lifecycleResponse.State.Id+"'")
        tflog.Debug(ctx, "Lifecycle Operation HTTP Status: '"+response.HttpResponse.Status+"'")

        if lifecycleResponse.State.Id == "Succeeded" || lifecycleResponse.State.Id == "Failed" {
            return &lifecycleResponse, nil
        }
    }
}
```

## Fix

To ensure proper context cancellation handling, the loop should check for the context's completion status and exit appropriately when the context is canceled. Below is the adjusted function:

```go
func (client *Client) DoWaitForLifecycleOperationStatus(ctx context.Context, response *Response) (*LifecycleDto, error) {
    locationHeader := response.GetHeader(constants.HEADER_LOCATION)
    if locationHeader == "" {
        locationHeader = response.GetHeader(constants.HEADER_OPERATION_LOCATION)
    }
    tflog.Debug(ctx, "Location Header: "+locationHeader)

    _, err := url.Parse(locationHeader)
    if err != nil {
        tflog.Error(ctx, "Error parsing location header: "+err.Error())
    }

    waitFor := retryAfter(ctx, response.HttpResponse)

    for {
        select {
        case <-ctx.Done():
            tflog.Warn(ctx, "Context canceled: exiting lifecycle operation wait")
            return nil, ctx.Err()
        default:
            lifecycleResponse := LifecycleDto{}
            response, err = client.Execute(ctx, nil, "GET", locationHeader, nil, nil, []int{http.StatusOK}, &lifecycleResponse)
            if err != nil {
                return nil, err
            }

            err = client.SleepWithContext(ctx, waitFor)
            if err != nil {
                return nil, err
            }

            tflog.Debug(ctx, "Lifecycle Operation State: '"+lifecycleResponse.State.Id+"'")
            tflog.Debug(ctx, "Lifecycle Operation HTTP Status: '"+response.HttpResponse.Status+"'")

            if lifecycleResponse.State.Id == "Succeeded" || lifecycleResponse.State.Id == "Failed" {
                return &lifecycleResponse, nil
            }
        }
    }
}
```

This change ensures that the function will check for context cancellation and terminate the operation appropriately, preventing potential hangs and resource issues.