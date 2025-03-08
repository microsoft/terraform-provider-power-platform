# Terraform Resource Implementation Template for Power Platform Provider

## OVERVIEW

Your task is to implement a new Terraform resource for the Power Platform Terraform Provider based on the specification in the attached user story.

## PRINCIPLES & RULES

- Your code must adhere to the established coding patterns of the repository
- Must use the [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
- Include the appropriate license header
- Use the correct package name (matching the folder except for tests) at the beginning of source code files
- Follow the naming conventions (e.g. `powerplatform_${resource_name}`).
- Reuse helper functions and client methods already available in the repository.
- Use `resp.Diagnostics.AddError` to report errors when available

## Common Resource Implementation Patterns

### Resource Organization

- Resources are organized in packages by service domain (e.g., `environment`, `connection`)
- Each resource has its own file named `resource_<resource_name>.go`
- Associated DTOs and models are often in separate files
- Resources are registered in `provider.go` via the `Resources()` method

### Resource Structure

- Resources implement the `resource.Resource` interface
- Most resources also implement `resource.ResourceWithImportState` for import capability
- Some implement additional interfaces like `resource.ResourceWithValidateConfig`
- Standard method implementation order: `Metadata`, `Schema`, `Create`, `Read`, `Update`, `Delete`, `ImportState`

### Schema Definition

- Each attribute includes a detailed `MarkdownDescription` with links to official documentation
- Attributes clearly marked as `Required`, `Optional`, or `Computed`
- Nested schemas for complex data structures (objects, lists, sets). Prefer sets over lists.
- Plan modifiers control behavior during state transitions:
  - `RequiresReplace()` for attributes that force resource recreation
  - `UseStateForUnknown()` for computed values that shouldn't be reset
  - Custom modifiers for special cases

### Error Handling

- Errors are reported via `resp.Diagnostics.AddError(title, detail)`
- "Not found" errors in Read operations lead to state removal
- Errors include context about the resource and operation

### Client Interactions

- Resources contain a client field for API communication
- API operations are delegated to the client
- State conversion happens in the resource methods

### Testing

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

### Documentation

- Markdown documentation in `/docs/resources/`
- Example configurations showing how to use the resource
- Complete schema documentation with descriptions

### Common Features

- Import support via the `ImportState` method
- Timeouts support using the `timeouts` package
- ID-based resource identification
- Context-aware logging with `tflog`

This consistent structure makes the codebase maintainable and helps ensure a consistent user experience across different resources.
