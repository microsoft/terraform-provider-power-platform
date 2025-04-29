# Title

Environment variable state persistence issue in multi-string configuration tests

##

/workspaces/terraform-provider-power-platform/internal/helpers/config_test.go

## Problem

The environment variables set during `TestUnitGetConfigMultiString_Matrix` are not reset or cleaned up after individual test cases. This can lead to inconsistent results when tests interact or are run in sequence due to residual state.

## Impact

Persistent states from previous test runs could result in unreliable or flaky tests, making it difficult to ensure accurate test outcomes, especially in collaborative environments with multiple contributors. Severity: High.

## Location

Line starting at `t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME, value)` in the `TestUnitGetConfigMultiString_Matrix` function.

## Code Issue

```go
for _, value := range testCase.environmentValues {
    t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME, value)
}
```

## Fix

Introduce logic to unset or reset the environment variables after each test. Use `defer` within each test case to clean up after execution.

```go
t.Run(testCase.name, func(t *testing.T) {
    for _, value := range testCase.environmentValues {
        t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME, value)
        defer t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME, "") // Ensure cleanup
    }
    
    var configValue basetypes.StringValue

    if testCase.configValue != nil {
        configValue = basetypes.NewStringValue(*testCase.configValue)
    } else {
        configValue = basetypes.NewStringNull()
    }

    ctx := context.Background()

    result := helpers.GetConfigMultiString(ctx, configValue, testCase.environmentValues, testCase.defaultValue)

    if result != testCase.expectedValue {
        t.Errorf("Expected '%s', got '%s'", testCase.expectedValue, result)
    }
})
```

This fix ensures that environment variables are properly unset, preserving the integrity of the test environment for other tests.