# Title

Environment variable state persistence issue in tests for boolean configuration

##

/workspaces/terraform-provider-power-platform/internal/helpers/config_test.go

## Problem

Similar to `TestUnitGetConfigString_Matrix`, the environment variables set during the test runs in `TestUnitGetConfigBool_Matrix` are not cleaned up or reset after individual test cases. This results in persistent states that could affect subsequent tests.

## Impact

This issue can lead to flaky and unreliable tests, making it difficult to accurately validate test outcomes in scenarios where multiple tests depend on environment variables. Severity: High.

## Location

Line starting at `t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME, *testCase.environmentValue)` in the `TestUnitGetConfigBool_Matrix` function.

## Code Issue

```go
if testCase.environmentValue != nil {
    t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME, *testCase.environmentValue)
}
```

## Fix

Add cleanup logic using `defer` to ensure that each test case in the function properly resets the environment variable state.

```go
t.Run(testCase.name, func(t *testing.T) {
    if testCase.environmentValue != nil {
        t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME, *testCase.environmentValue)
        defer t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME, "") // Ensure cleanup
    }
    
    if testCase.configValue != nil {
        configValue = basetypes.NewBoolValue(*testCase.configValue)
    } else {
        configValue = basetypes.NewBoolNull()
    }

    ctx := context.Background()

    result := helpers.GetConfigBool(ctx, configValue, TEST_ENVIRONMENT_VARIABLE_NAME, testCase.defaultValue)

    if result != testCase.expectedValue {
        t.Errorf("Expected '%t', got '%t'", testCase.expectedValue, result)
    }
})
```

This fix ensures a consistent environment variable state across all test cases, preventing interference between test cases.