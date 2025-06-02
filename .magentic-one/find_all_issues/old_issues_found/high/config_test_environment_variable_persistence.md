# Title

Environment variable state persistence issue in tests

##

/workspaces/terraform-provider-power-platform/internal/helpers/config_test.go

## Problem

The environment variables set during the test runs are not reset after individual test cases. This is problematic because environment variables persist across multiple test runs, leading to potential test interference and flaky test results due to residual environment variable states from previous tests.

## Impact

This issue can lead to inconsistent test results when tests are run repeatedly or together with other tests that rely on environment variables. It increases debugging difficulty and can mask genuine issues. Severity: High.

## Location

Line starting at `t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME, *testCase.environmentValue)` in the `TestUnitGetConfigString_Matrix` function.

## Code Issue

```go
if testCase.environmentValue != nil {
    t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME, *testCase.environmentValue)
}
```

## Fix

Each test case must ensure that environment variables are unset or reset back to their original state at the end of the test run to prevent leakage.

```go
t.Run(testCase.name, func(t *testing.T) {
    if testCase.environmentValue != nil {
        t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME, *testCase.environmentValue)
        defer t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME, "") // Ensure cleanup
    }
    
    if testCase.configValue != nil {
        configValue = basetypes.NewStringValue(*testCase.configValue)
    } else {
        configValue = basetypes.NewStringNull()
    }

    ctx := context.Background()

    result := helpers.GetConfigString(ctx, configValue, TEST_ENVIRONMENT_VARIABLE_NAME, testCase.defaultValue)

    if result != testCase.expectedValue {
        t.Errorf("Expected '%s', got '%s'", testCase.expectedValue, result)
    }
})
```

This fix uses `defer` to revert the environment variable state to ensure that the next test case does not use a modified state.