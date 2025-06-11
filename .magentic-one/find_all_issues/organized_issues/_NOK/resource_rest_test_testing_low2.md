# Title

Insufficient unit/acceptance split in resource test structure

##

/workspaces/terraform-provider-power-platform/internal/services/rest/resource_rest_test.go

## Problem

The test functions `TestAccTestRest_Validate_Create` and `TestUnitTestRest_Validate_Create` have similar logic and structure, but the separation between acceptance and unit test setup is not clear and may invite confusion. Specifically, the naming and the presence of both acceptance and unit test frameworks (e.g., use of real provider factories versus mocks) are interleaved with only limited differences. Stronger and clearer structure/enforcement should exist between acceptance (live/real infrastructure, full flow) and unit (mocked APIs, in-process logic) testing.

## Impact

This structural duality and lack of explicit distinction could cause confusion for test contributors and reviewers, increase the risk of misconfiguration or the accidental use of real infrastructure in unit tests, and generally makes the test intent harder to read. This is a low severity issue for small teams but becomes more critical as the codebase and contributor set grows.

## Location

Definition and structure of `TestAccTestRest_Validate_Create` vs `TestUnitTestRest_Validate_Create` as seen here:

```go
func TestAccTestRest_Validate_Create(t *testing.T) {
    ...
    resource.Test(t, resource.TestCase{
        ...
    })
}

func TestUnitTestRest_Validate_Create(t *testing.T) {
    ...
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()
    ...
    resource.Test(t, resource.TestCase{
        ...
        IsUnitTest: true,
        ...
    })
}
```

## Fix

Refactor to clarify the split between acceptance and unit test infrastructure, by:
- Clearly commenting/explaining the difference at the start of each function.
- Structuring the setup code for acceptance and unit tests in separately reusable helpers where possible.
- Naming test configurations and helper functions to clearly indicate unit/acceptance context.
- Ensuring mocks and provider factories are strictly scoped to the correct test.

For example:

```go
func setupAccEnv() string {
    // returns config string for acceptance environment
}
func setupUnitTestEnv() string {
    // returns config string for unit/mocked environment
}

func TestRest_ValidateCreate(t *testing.T) {
    config := setupAccEnv()
    // ...
}

func TestRest_ValidateCreate_Unit(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()
    config := setupUnitTestEnv()
    // ...
}
```
This makes the distinction explicit, easier to maintain, and less confusing for new developers.
