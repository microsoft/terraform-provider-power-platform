# Testing Guidance

1. **Test organization**:
   - Tests are placed in the same directory as the implementation files
   - The test package name has a `_test` suffix (e.g., `package solution_checker_rules_test`)
   - Test files follow a naming convention: `<resource_or_datasource_name>_test.go`

2. **Two types of tests**:
   - **Acceptance Tests**: Named with a `TestAcc` prefix (e.g., `TestAccSolutionCheckerRulesDataSource_Basic`)
     - These tests run against actual infrastructure
     - They use `SkipUnlessAcceptanceTestMode(t)` to only run in acceptance test mode
     - They use `testhelpers.TestAccProtoV6ProviderFactories` for provider setup

   - **Unit Tests**: Named with a `TestUnit` prefix (e.g., `TestUnitSolutionCheckerRulesDataSource_Validate_Read`)
     - These use HTTP mocking to simulate API responses
     - They use `httpmock.Activate()` and `httpmock.DeactivateAndReset()` for HTTP mocking
     - They use `testhelpers.RegisterHTTPResponse()` to register HTTP responses
     - They use `testhelpers.TestUnitTestProtoV6ProviderFactories` for provider setup
     - They reference JSON files in a `tests/` directory rather than embedding JSON responses directly in the test file

3. **Test structure**:
   - Tests use `resource.ParallelTest()` or `resource.UnitTest()` with a `resource.TestCase` structure
   - Test configurations are defined inline as strings
   - `resource.ComposeAggregateTestCheckFunc()` is used to assert on the test results
   - Various check functions like `resource.TestCheckResourceAttr` are used to verify resource attributes

4. **Mocking**:
   - HTTP responses are mocked using `httpmock` from the `jarcoal/httpmock` package
   - Mock JSON responses are embedded in the test file for HTTP mock responses

The unit tests follow a clear pattern of setting up HTTP mocks to simulate API responses and then verifying that the data source correctly processes those responses. The acceptance tests are designed to run against real infrastructure and verify that the data source works correctly with live APIs.
