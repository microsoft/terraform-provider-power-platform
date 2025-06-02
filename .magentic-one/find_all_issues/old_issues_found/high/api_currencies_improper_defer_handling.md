# Title

Improper Handling of Deferred Body Closure for `response.HttpResponse.Body.Close()`

## 

`/workspaces/terraform-provider-power-platform/internal/services/currencies/api_currencies.go`

## Problem

The deferred statement `defer response.HttpResponse.Body.Close()` is called in the function `GetCurrenciesByLocation`, but there is no validation of whether `response.HttpResponse` is non-nil before attempting to access the `Body` property.

If an error occurred during `client.Api.Execute`, the `response.HttpResponse` could potentially be `nil`, and calling `defer response.HttpResponse.Body.Close()` in that scenario can cause a runtime panic.

## Impact

- **Severity:** **High**
- Risk of **runtime panic** if `HttpResponse` is nil, leading to unexpected termination of the application at runtime.
  
## Location

This issue occurs in the `GetCurrenciesByLocation` function within the following code block:

## Code Issue

```go
defer response.HttpResponse.Body.Close()
```

## Fix

Ensure that `response.HttpResponse` is non-nil before deferring closure. Add a conditional check to validate its existence. For example:

```go
if response != nil && response.HttpResponse != nil {
    defer response.HttpResponse.Body.Close()
}
```

Full fixed code snippet:

```go
func (client *client) GetCurrenciesByLocation(ctx context.Context, location string) (currenciesDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/locations/%s/environmentCurrencies", location),
	}
	apiUrl.RawQuery = url.Values{
		"api-version": []string{"2023-06-01"},
	}.Encode()

	currencies := currenciesDto{}

	response, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return currencies, err
	}

	if response != nil && response.HttpResponse != nil {
		defer response.HttpResponse.Body.Close()
	}

	err = json.Unmarshal(response.BodyAsBytes, &currencies)

	if err != nil {
		return currencies, err
	}

	return currencies, nil
}
```