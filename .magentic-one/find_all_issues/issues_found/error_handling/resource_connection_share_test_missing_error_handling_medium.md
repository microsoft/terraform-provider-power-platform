# Title

Missing error handling in HTTP mock responder

##

/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection_share_test.go

## Problem

In both HTTP mock responders (`POST` and `GET`), the registered responder functions ignore error responses or error simulation, always returning nil as the error value. The code does not provide the ability to simulate or test error cases (network issues, bad responses) which is important for robust test coverage.

## Impact

This restricts the ability to test error handling logic in the provider code. The test only validates successful (happy-path) scenarios, potentially missing regressions or missed code branches in actual error pathways. Severity: medium.

## Location

Lines around the following responders (search: `RegisterResponder`):

```go
httpmock.RegisterResponder("POST", `https://.../modifyPermissions?...`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, ""), nil
    })

httpmock.RegisterResponder("GET", `https://.../permissions?...`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/connection_shares/Validate_Create/get_connection_shares.json").String()), nil
    })
```

## Code Issue

```go
httpmock.RegisterResponder("POST", `https://000000000000000000000000000000.00.environment.api.powerplatform.com/connectivity/connectors/shared_commondataserviceforapps/connections/00000000-0000-0000-0000-000000000001/modifyPermissions?%24filter=environment+eq+%2700000000-0000-0000-0000-000000000000%27&api-version=1`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, ""), nil
    })

httpmock.RegisterResponder("GET", `https://000000000000000000000000000000.00.environment.api.powerplatform.com/connectivity/connectors/shared_commondataserviceforapps/connections/00000000-0000-0000-0000-000000000001/permissions?%24filter=environment+eq+%2700000000-0000-0000-0000-000000000000%27&api-version=1`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/connection_shares/Validate_Create/get_connection_shares.json").String()), nil
    })
```

## Fix

Update the test to optionally simulate error scenarios. This could involve parameterizing the test with additional steps to simulate error HTTP responses. At a minimum, create an additional test or step to handle and assert on error responses:

```go
// Example for error case testing
httpmock.RegisterResponder("POST", `https://.../modifyPermissions?...`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusInternalServerError, "Internal Server Error"), nil
    })
// Then, in test steps, test expected error handling logic
```

This structure allows expanding the suite to test failure cases, improving test coverage of error handling.
