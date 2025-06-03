# Title

Function and Method Naming Isn't Clear About Export Status

##

/workspaces/terraform-provider-power-platform/internal/services/currencies/api_currencies.go

## Problem

The factory function `newCurrenciesClient` uses camelCase, matching an unexported function in Go. However, if your usage intends for it to be used outside the package, it should be exported (`NewCurrenciesClient`). Otherwise, stick with the convention, but ensure this fits your intended use.

## Impact

**Low**. Code readability and package API clarity.

## Location

Definition of the constructor function:

```go
func newCurrenciesClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}
```

## Fix

If you intend to export the client:

```go
func NewCurrenciesClient(apiClient *api.Client) CurrenciesClient {
	return CurrenciesClient{
		Api: apiClient,
	}
}
```

If intended for internal use only, the current convention is correct.
