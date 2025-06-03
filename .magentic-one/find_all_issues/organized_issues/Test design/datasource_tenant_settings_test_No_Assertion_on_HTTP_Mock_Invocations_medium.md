# No Assertion on HTTP Mock Invocations or Coverage

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/datasource_tenant_settings_test.go

## Problem

The tests use `httpmock` to mock the HTTP responses of the external API, but they do not assert that all registered responders are actually invoked, nor do they check the count of calls. This can lead to false positives where a test passes because the resource logic never actually issues an HTTP request, or the registered responder is not hit due to a URL mismatch or bug. Proper verification of mock invocation is essential in unit tests using mocks.

## Impact

- **Severity:** Medium
- Lax assertion on HTTP mock usage allows for broken or skipping code paths to appear as "passing" tests, resulting in logical gaps and undetected regressions.
- This reduces the reliability of the test suite, especially as mocking is critical for correctness in provider tests.

## Location

In the `TestUnitTestTenantSettingsDataSource_Validate_Read` function.

## Code Issue

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()

httpmock.RegisterResponder("POST", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/listTenantSettings?api-version=2023-06-01`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/post_list_tenant_settings.json").String()), nil
    })

resource.Test(t, resource.TestCase{
    // ...
})
```

## Fix

Add a verification step after the test logic to assert the number of mock HTTP responder invocations for the registered endpoint.

```go
// After the test run
info := httpmock.GetCallCountInfo()
expectedKey := "POST https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/listTenantSettings?api-version=2023-06-01"
if info[expectedKey] == 0 {
    t.Errorf("Expected API endpoint to be called, but it was not")
}
```

Place this after `resource.Test(...)` in your test.

---

This ensures your mocks are used and will fail if the API is not called as designed.
