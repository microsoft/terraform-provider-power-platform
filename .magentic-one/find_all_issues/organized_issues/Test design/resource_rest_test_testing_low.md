# Title

No error assertion on responder registration in HTTP mock

##

/workspaces/terraform-provider-power-platform/internal/services/rest/resource_rest_test.go

## Problem

In the unit test `TestUnitTestRest_Validate_Create`, the code registers HTTP responders using `httpmock.RegisterResponder` but does not assert or check for errors returned by this call. The `RegisterResponder` method can return an error if, for example, the responder pattern is invalid or already registered. Not handling the error is a missed opportunity for defense-in-depth in the test setup.

## Impact

If a responder is not properly registered due to an error, the test will not behave as expected, possibly leading to misleading test results. Errors can be introduced silently—especially as the test suite grows or HTTP patterns are adjusted. This is a low severity issue in the context of this provider test, but it's important for robust test code.

## Location

Within function `TestUnitTestRest_Validate_Create(t *testing.T)`, in these lines:
```go
httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Web_Api_Validate_Create/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
    })

httpmock.RegisterResponder("POST", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts?$select=name,accountid`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/Web_Api_Validate_Create/post_account.json").String()), nil
    })

httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts(00000000-0000-0000-0000-000000000001)?$select=name,accountid`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/Web_Api_Validate_Create/post_account.json").String()), nil
    })
```

## Fix

Capture and check the error on each `RegisterResponder` call, failing the test immediately if one occurs. For example:

```go
if err := httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Web_Api_Validate_Create/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
    }); err != nil {
    t.Fatalf("failed to register responder: %v", err)
}
```

Repeat the pattern for every registration within the test function. This makes failures explicit and catches setup issues early.
