# Title

Direct State Mutation in Mocking (`isOdataQueryRun` variable pattern)

##

internal/services/data_record/datasource_data_record_test.go

## Problem

Multiple unit tests use a package-level or function-local Boolean variable (`isOdataQueryRun`) that gets mutated inside an HTTP responder registered by the httpmock package, to determine if the expected query was actually executed (i.e., the responder was invoked). While functional for single-threaded, non-parallelized tests, this pattern is fragile:  
- It does not scale safely to parallel test execution (`t.Parallel()`).
- If any test fails before the state is checked, the mutation is lost and the assertion is never reached.
- Direct mutation inside a responder closure is more error-prone than using explicit hooks or channels.

## Impact

- May make tests non-thread-safe or non-deterministic if run in parallel with others.
- Makes the tests fragile to refactoring—state can be missed if multiple responders are set.
- Makes the test logic harder to follow than more standard assertion patterns.

Severity: Low to Medium (as long as tests are not parallelized, issue is minor; if parallelized, can cause flakiness).

## Location

Pattern found in most `TestUnitDataRecordDatasource_...` methods, e.g.:

## Code Issue

```go
isOdataQueryRun := false

httpmock.RegisterResponder("GET", "...", func(req *http.Request) (*http.Response, error) {
    isOdataQueryRun = true
    return httpmock.NewStringResponse(http.StatusOK, `...`), nil
})

// ... later in the test
if !isOdataQueryRun {
    t.Errorf("Odata query should have been run in '%s' unit test", mocks.TestName())
}
```

## Fix

Replace with a safer test assertion strategy. For example, use the built-in counting features of httpmock or a custom test wrapper to count invocations:

```go
var callCount int
httpmock.RegisterResponder("GET", "...", func(req *http.Request) (*http.Response, error) {
    callCount++;
    return httpmock.NewStringResponse(http.StatusOK, `...`), nil
})
// At the end:
if callCount < 1 {
    t.Errorf("Expected Odata query responder to be called at least once in '%s' unit test", mocks.TestName())
}
```

Or, use the features httpmock provides for asserting call counts, e.g.:
```go
httpmock.RegisterResponder(...)
...
calls := httpmock.GetCallCountInfo()
if calls["GET ..."] < 1 {
    t.Errorf("Expected call count, etc")
}
```

This will ensure thread safety and more robust test assertions, especially as the test suite grows or is parallelized.

Save as a testing/code structure issue.
