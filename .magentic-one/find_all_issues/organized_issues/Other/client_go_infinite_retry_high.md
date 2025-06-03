# Title

Error handling in the HTTP request/response loop risks masking underlying issues

##

/workspaces/terraform-provider-power-platform/internal/api/client.go

## Problem

The `Execute` function attempts to retry HTTP requests for certain status codes but only logs retry attempts and never limits the number of retries or implements backoff strategies. This can cause a potential infinite retry loop for a persistently failing endpoint, leading to wasted resources or a stuck process.

## Impact

May cause resource exhaustion, hanging processes, or unresponsiveness. Severity: **high**

## Location

Within `Client.Execute` for-loop, no maximum retry, nor context deadline check before retry.

## Code Issue

```go
for {
    // ... various request code ...
    if !isRetryable {
        return resp, customerrors.NewUnexpectedHttpStatusCodeError(acceptableStatusCodes, resp.HttpResponse.StatusCode, resp.HttpResponse.Status, resp.BodyAsBytes)
    }

    waitFor := retryAfter(ctx, resp.HttpResponse)
    tflog.Debug(ctx, fmt.Sprintf("Received status code %d for request %s, retrying after %s", resp.HttpResponse.StatusCode, url, waitFor))

    err = client.SleepWithContext(ctx, waitFor)
    if err != nil {
        return resp, err
    }
}
```

## Fix

Introduce a maximum retry count and/or a total elapsed timeout, and return an error if exceeded. Example:

```go
maxRetries := 5
retries := 0
for {
    // ... request logic ...
    if !isRetryable {
        return resp, customerrors.NewUnexpectedHttpStatusCodeError(acceptableStatusCodes, resp.HttpResponse.StatusCode, resp.HttpResponse.Status, resp.BodyAsBytes)
    }
    if retries >= maxRetries {
        return resp, fmt.Errorf("Maximum retries exceeded for request to %s", url)
    }
    retries++
    waitFor := retryAfter(ctx, resp.HttpResponse)
    tflog.Debug(ctx, fmt.Sprintf("Received status code %d for request %s, retrying after %s (attempt %d/%d)", resp.HttpResponse.StatusCode, url, waitFor, retries, maxRetries))
    err = client.SleepWithContext(ctx, waitFor)
    if err != nil {
        return resp, err
    }
}
```
