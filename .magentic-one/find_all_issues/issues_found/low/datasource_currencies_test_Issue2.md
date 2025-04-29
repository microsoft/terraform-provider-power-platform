# Title

Potential Redundancy in HTTP Mock Setup

##

/workspaces/terraform-provider-power-platform/internal/services/currencies/datasource_currencies_test.go

## Problem

In the `TestUnitCurrenciesDataSource_Validate_Read` test case, `httpmock.Activate()` and `httpmock.DeactivateAndReset()` are used to manage the lifecycle of HTTP mock setups. While this approach is functional, it could result in redundancy or errors if further tests depend on the same API or forget to deactivate the mocks.

## Impact

Although the current tests seem correct, potential issues could arise if multiple test cases rely on HTTP mocks without proper isolation. This may lead to conflicts or inconsistencies in testing results.

Severity: Low

## Location

`TestUnitCurrenciesDataSource_Validate_Read` test case.

## Code Issue

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()

httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentCurrencies?api-version=2023-06-01`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_currencies.json").String()), nil
    })
```

## Fix

Consider wrapping `httpmock.Activate()` and `httpmock.DeactivateAndReset()` within reusable helper functions to ensure clarity and consistency:

```go
func SetupHTTPMock(t *testing.T, responderFunc httpmock.Responder) {
    httpmock.Activate()
    httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentCurrencies?api-version=2023-06-01`, responderFunc)

    t.Cleanup(func() {
        httpmock.DeactivateAndReset()
    })
}
```