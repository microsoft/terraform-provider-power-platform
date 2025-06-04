# Ineffective Error Handling: Lost Context in API Calls

##

/workspaces/terraform-provider-power-platform/internal/services/connectors/api_connectors.go

## Problem

When `client.Api.Execute` returns an error, the function simply returns the error as-is without any additional context about which API call failed or what part of the `GetConnectors` operation was unsuccessful. This means the error will lack crucial debugging information.

## Impact

It can be difficult to determine which of the several API calls failed. This can complicate debugging and error reporting, making support and maintenance harder, especially when end users report API failures. Severity: **medium**.

## Location

Each `if err != nil { return nil, err }` after the three API calls within the `GetConnectors` method:

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connectorArray)
if err != nil {
	return nil, err
}
...
_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &unblockableConnectorArray)
if err != nil {
	return nil, err
}
...
_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &virtualConnectorArray)
if err != nil {
	return nil, err
}
```

## Code Issue

```go
if err != nil {
	return nil, err
}
```

## Fix

Wrap the errors with context to indicate which API call failed (requires Go 1.13+):

```go
import (
	"fmt"
	...
)

...

_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connectorArray)
if err != nil {
	return nil, fmt.Errorf("failed to fetch PowerApps connectors: %w", err)
}
...
_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &unblockableConnectorArray)
if err != nil {
	return nil, fmt.Errorf("failed to fetch unblockable connectors: %w", err)
}
...
_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &virtualConnectorArray)
if err != nil {
	return nil, fmt.Errorf("failed to fetch virtual connectors: %w", err)
}
```
