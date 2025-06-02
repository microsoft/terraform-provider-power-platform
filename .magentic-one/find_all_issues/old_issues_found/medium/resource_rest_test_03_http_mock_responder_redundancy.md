# Title
Unnecessary Re-registration of HTTP Mock Responder

## Problem
The `httpmock.RegisterResponder` methods are repeatedly used to register the same mocked URLs within the `TestUnitTestRest_Validate_Create()` function. This could be optimized for reusable mocks or placed within a setup function that ensures mock consistency across test cases.

## Impact
Excessive redundancy in re-registering responders reduces maintainability and increases the difficulty of debugging. If changes are required in the mock setup, developers may need to update every occurrence. This issue is of **medium severity**.

## Location
This issue is found in the `TestUnitTestRest_Validate_Create()` function of the file `resource_rest_test.go`.

## Code Issue
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
Create a helper function `registerMockResponders()` to centralize the HTTP mock setup logic. This reduces redundancy and ensures all mocks are consistent across tests.

```go
func registerMockResponders() {
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
}

// Use the helper function within the test case
func TestUnitTestRest_Validate_Create(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    registerMockResponders()

    resource.Test(t, resource.TestCase{
        ...
    })
}
```