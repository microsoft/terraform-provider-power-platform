# Title

Missing Unit Testing for `GetCurrenciesByLocation` Function

## 

`/workspaces/terraform-provider-power-platform/internal/services/currencies/api_currencies.go`

## Problem

The `GetCurrenciesByLocation` function retrieves currency information based on location context, constructs URLs dynamically, and depends on external API calls. However, the code lacks associated unit tests or error handling mocks to safeguard this function.

Key points for unit testing missed:
- Dependency injection for the `Client` structure isn't leveraged effectively for testable mocking.
- Dynamic URL construction and HTTP responses are not validated with expected inputs and corresponding outputs.
- Error scenarios (e.g., malformed JSON, API call failures) are not simulated.

## Impact

- **Severity:** **Critical**
- Lack of unit tests increases the likelihood of unseen bugs in production and decreases confidence in the quality of the implementation.
- Test coverage ensures proper functionality when refactoring or extending the code.
- Edge cases (e.g., invalid location input or server-side errors) aren't validated, increasing runtime risk.

## Location

Entire function `GetCurrenciesByLocation` lacks associated unit tests. Highlighted fragment:

## Code Issue

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

	defer response.HttpResponse.Body.Close()

	err = json.Unmarshal(response.BodyAsBytes, &currencies)

	if err != nil {
		return currencies, err
	}

	return currencies, nil
}
```

## Fix

Add unit tests and mocking for API calls using a popular Go mocking framework like `mockery` or `gomock`. Below is a suggestion for a basic unit test framework setup:

### Unit Test
```go
package currencies_test

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/microsoft/terraform-provider-power-platform/internal/services/currencies"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockApiClient struct {
	mock.Mock
}

// Mock the Execute function
func (m *mockApiClient) Execute(ctx context.Context, headers map[string]string, method, url string, body []byte, extraHeaders map[string]string, acceptedStatus []int, queryParameters map[string]string) (*api.Response, error) {
	args := m.Called(ctx, headers, method, url, body, extraHeaders, acceptedStatus, queryParameters)
	responseArg := args.Get(0)
	if responseArg != nil {
		return responseArg.(*api.Response), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestGetCurrenciesByLocation(t *testing.T) {
	mockClient := new(mockApiClient)
	currencyClient := currencies.NewCurrenciesClient(mockClient)

	mockResponse := &api.Response{
		BodyAsBytes: []byte(`{"currencies": []}`),
	}

	mockClient.On("Execute", mock.AnythingOfType("*context.emptyCtx"), nil, "GET", mock.AnythingOfType("string"), nil, nil, mock.Anything, nil).Return(mockResponse, nil)

	result, err := currencyClient.GetCurrenciesByLocation(context.Background(), "mockLocation")

	assert.Nil(t, err)
	assert.NotNil(t, result)
	mockClient.AssertExpectations(t)
}
```

By adding tests, we can validate interactions under multiple scenarios, such as:
1. Valid location input.
2. Invalid API responses.
3. Downstream HTTP errors.

Having testing as a core improvement ensures maintainable and extendable code.