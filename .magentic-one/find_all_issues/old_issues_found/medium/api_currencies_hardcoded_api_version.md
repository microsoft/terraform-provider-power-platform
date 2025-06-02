# Title

Hardcoded API version in `RawQuery`

## 

`/workspaces/terraform-provider-power-platform/internal/services/currencies/api_currencies.go`

## Problem

The API version `2023-06-01` is hardcoded in the query parameters within the `apiUrl.RawQuery` construction. This creates maintainability issues because changes to the API version will require direct modification of the code.

Hardcoding such configuration values violates the principle of centralizing configurations and parameters, making the program less flexible and error-prone.

## Impact

- **Severity:** **Medium**
- The maintainability of the program is reduced as updating the API version necessitates code changes instead of configuration updates.
- The code becomes less reusable.

## Location

Problem exists at the following code block:

## Code Issue

```go
apiUrl.RawQuery = url.Values{
    "api-version": []string{"2023-06-01"},
}.Encode()
```

## Fix

Externalize the version string into a configuration constant or environment variable, so future changes do not require direct modification of the code.

Below is one possible fix:

```go
const APIVersion = "2023-06-01"

func (client *client) GetCurrenciesByLocation(ctx context.Context, location string) (currenciesDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/locations/%s/environmentCurrencies", location),
	}
	apiUrl.RawQuery = url.Values{
		"api-version": []string{APIVersion},
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

Using `APIVersion` allows the software to accept dynamic inputs via environment variables or configuration files. If variations are needed for testing or deployment, they can easily be provided without re-compiling code.