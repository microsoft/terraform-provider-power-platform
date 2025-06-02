# Title

Improper Error Handling in Function `GetConnectors`

##

`/workspaces/terraform-provider-power-platform/internal/services/connectors/api_connectors.go`

## Problem

The error handling in the `GetConnectors` function does not provide sufficient context when returning errors. If an error occurs during the HTTP request execution, the returned error does not include information about which API call failed or why the request failed.

## Impact

The lack of detailed error information impacts debugging and troubleshooting efforts. Without sufficient context, developers will struggle to identify the root cause of errors. This is a **high severity** issue because error handling and logging are critical for identifying issues during execution.

## Location

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connectorArray)
if err != nil {
    return nil, err
}

_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &unblockableConnectorArray)
if err != nil {
    return nil, err
}

_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &virtualConnectorArray)
if err != nil {
    return nil, err
}
```

## Code Issue

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connectorArray)
if err != nil {
    return nil, err
}
```

## Fix

Provide more context to the errors when returning them. This can be done by using the `fmt.Errorf()` function or wrapping the errors with additional information.

```go
import (
    "fmt"
)

_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connectorArray)
if err != nil {
    return nil, fmt.Errorf("failed to fetch connectors from %s: %w", apiUrl.String(), err)
}

_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &unblockableConnectorArray)
if err != nil {
    return nil, fmt.Errorf("failed to fetch unblockable connectors from %s: %w", apiUrl.String(), err)
}

_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &virtualConnectorArray)
if err != nil {
    return nil, fmt.Errorf("failed to fetch virtual connectors from %s: %w", apiUrl.String(), err)
}
```

This will give developers a better understanding of the issue and the exact API call that caused the error.