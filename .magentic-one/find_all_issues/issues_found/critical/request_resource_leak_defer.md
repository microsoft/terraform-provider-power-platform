# Title

Potential Resource Leak with `defer apiResponse.Body.Close()`

##

`/workspaces/terraform-provider-power-platform/internal/api/request.go`

## Problem

The statement `defer apiResponse.Body.Close()` is executed after checking if an error occurred in the `doRequest` function. If `apiResponse` is nil due to an error during the request, calling `apiResponse.Body.Close()` will result in a runtime panic since `apiResponse.Body` will be dereferenced without verifying its presence.

## Impact

This issue can lead to runtime panics and crash the program during error conditions when `apiResponse` is nil. Severity: Critical.

## Location

Line: `defer apiResponse.Body.Close()` in the `doRequest` function.

## Code Issue

```go
apiResponse, err := httpClient.Do(request)

resp := &Response{
    HttpResponse: apiResponse,
}

if err != nil {
    return resp, err
}

defer apiResponse.Body.Close()
body, err := io.ReadAll(apiResponse.Body)
resp.BodyAsBytes = body
```

## Fix

Introduce a check to ensure `apiResponse` and `apiResponse.Body` are non-nil before attempting to close the Body in a deferred statement.

```go
apiResponse, err := httpClient.Do(request)

resp := &Response{
    HttpResponse: apiResponse,
}

if err != nil {
    return resp, err
}

if apiResponse != nil && apiResponse.Body != nil {
    defer apiResponse.Body.Close()
}

body, err := io.ReadAll(apiResponse.Body)
resp.BodyAsBytes = body
```

This ensures that the Body is only closed if it actually exists, thereby preventing runtime panics during the error handling process.
