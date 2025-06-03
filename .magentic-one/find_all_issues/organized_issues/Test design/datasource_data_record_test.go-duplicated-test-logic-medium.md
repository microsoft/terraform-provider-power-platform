# Title

Duplicated Test Logic and Potential for Test Fragility

##

internal/services/data_record/datasource_data_record_test.go

## Problem

Multiple test functions in this suite repeat nearly identical logic—mocking API responders, activating/deactivating mocks, building similar configurations, and verifying essentially the same pattern. This causes code duplication and increases the maintenance burden. Additionally, directly manipulating or assuming the global state of mocks can make tests fragile and more prone to flakiness, especially when testing in parallel or when sharing state.

## Impact

Too much duplication increases the cost of test maintenance (medium impact). Any change or enhancement in how responses are mocked or how test configs are built necessitates updating many places. It also increases the risk that tests can fail for non-functional reasons (such as state bleed), affecting reliability and confidence in the test suite.

## Location

Present in (but not limited to):

- `TestUnitDataRecordDatasource_Validate_Expand_Query`
- `TestUnitDataRecordDatasource_Validate_Single_Record_Expand_Query`
- `TestUnitDataRecordDatasource_Validate_Top`
- `TestUnitDataRecordDatasource_Validate_Apply`
- `TestUnitDataRecordDatasource_Validate_OrderBy`
- `TestUnitDataRecordDatasource_Validate_SavedQuery`
- `TestUnitDataRecordDatasource_Validate_UserQuery`
- `TestUnitDataRecordDatasource_Validate_Expand_Lookup`

## Code Issue

Example pattern:
```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()
isOdataQueryRun := false

httpmock.RegisterResponder("GET", "....", ...)
// More responders

resource.Test(t, resource.TestCase{
    ...
})

// Last check:
if !isOdataQueryRun {
    t.Errorf("Odata query should have been run in '%s' unit test", mocks.TestName())
}
```

All of these code blocks are almost the same with only slight differences.

## Fix

Refactor common test boilerplate and mocking setup into reusable helper functions. For instance, you could use a function like this for API endpoint setup:

```go
func setupCommonMocks(endpoint string, resp string) (cleanup func(), *bool) {
    httpmock.Activate()
    isOdataQueryRun := false

    httpmock.RegisterResponder("GET", endpoint, func(req *http.Request) (*http.Response, error) {
        isOdataQueryRun = true
        return httpmock.NewStringResponse(http.StatusOK, resp), nil
    })
    // Register other common responders...
    cleanup = func() { httpmock.DeactivateAndReset() }
    return cleanup, &isOdataQueryRun
}
```
Use this in test functions:
```go
cleanup, isOdataQueryRun := setupCommonMocks(endpoint, response)
defer cleanup()
...
if !*isOdataQueryRun {
    t.Errorf("Odata query should have been run in '%s' unit test", mocks.TestName())
}
```
This makes the test suite smaller, clearer, and reduces the risk of introducing bugs in the testing logic as the provider/test suite evolves.
