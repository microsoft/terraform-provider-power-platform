# Testing Guidance

1. **Test organization**:
   - Tests are placed in the same directory as the implementation files
   - The test package name has a `_test` suffix (e.g., `package solution_checker_rules_test`)
   - Test files follow a naming convention: `<resource_or_datasource_name>_test.go`

2. **Two types of tests**:
   - **Acceptance Tests**: Named with a `TestAcc` prefix (e.g., `TestAccSolutionCheckerRulesDataSource_Validate_Create`)
     - These tests run against actual infrastructure
     - Acceptance tests should have test steps for testing the happy path for Create, Update, and Destroy
     - They use `mocks.TestAccProtoV6ProviderFactories` for provider setup

   - **Unit Tests**: Named with a `TestUnit` prefix (e.g., `TestUnitSolutionCheckerRulesDataSource_Validate_Read`)
     - These use HTTP mocking to simulate API responses
     - Unit tests should test boundry conditions and error handling
     - They use `httpmock.Activate()` and `httpmock.DeactivateAndReset()` for HTTP mocking
     - They use `testhelpers.RegisterHTTPResponse()` to register HTTP responses
     - They use `testhelpers.TestUnitTestProtoV6ProviderFactories` for provider setup
     - HTTP responses are mocked using `httpmock` from the `jarcoal/httpmock` package
     - Mock JSON responses should only be stored in json files in the `tests/` folder of each service. Do not embed mock JSON responses directly in the _test.go file
     - Calls to httpmock should not be nested.  Register the responders in a serial fashion to improve the readability of the tests.
     - Do not change mocked json data files to fix a broken test without explicit human instruction to do so

3. **Test structure**:
   - Tests use `resource.UnitTest()` with a `resource.TestCase` structure
   - Tests should not use `resource.ParallelTest()`
   - Test terraform configurations are defined inline as strings and do not include provider declarations without human approval for exceptional cases
   - `resource.ComposeAggregateTestCheckFunc()` is used to assert on the test results
   - Various check functions like `resource.TestCheckResourceAttr` are used to verify resource attributes
   - Be thourough in testing resource attributes

The unit tests follow a clear pattern of setting up HTTP mocks to simulate API responses and then verifying that the data source correctly processes those responses. The acceptance tests are designed to run against real infrastructure and verify that the data source works correctly with live APIs.
