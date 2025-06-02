# Environment Variable Pollution Across Test Cases

##

/workspaces/terraform-provider-power-platform/internal/helpers/config_test.go

## Problem

The tests in this file, especially those involving environment variables (such as `TestUnitGetConfigString_Matrix`, `TestUnitGetConfigBool_Matrix`, and `TestUnitGetConfigMultiString_Matrix`), set environment variables but do not unset or restore their previous values after each subtest. This can cause unexpected behaviors because environment variables are global and persist across `t.Run` subtests and, potentially, across other tests running in the same process, leading to test pollution and flaky tests.

## Impact

Severity: **High**  
Pollution of environment variables can lead to side effects between test cases, making test outcomes order-dependent and non-deterministic. This diminishes the reliability of the test suite, leading to possible false positives or negatives during test runs.

## Location

Example in `TestUnitGetConfigString_Matrix`:

```go
if testCase.environmentValue != nil {
    t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME, *testCase.environmentValue)
}
```

And in `TestUnitGetConfigMultiString_Matrix`:

```go
t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME1, environmentValue1)
t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME2, environmentValue2)
...
for (_, value := range testCase.environmentValues) {
    t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME, value)
}
```

## Code Issue

```go
if testCase.environmentValue != nil {
    t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME, *testCase.environmentValue)
}
// No code to unset or reset environment variables after test.
```

## Fix

Leverage the fact that `t.Setenv` (in Go 1.17+) automatically restores the environment variable when the test and any subtests complete. However, environments should be isolated per test and not set globally before running the test cases.  
Move any setting of environment variables inside the subtest closure, and ensure that environment variables are NOT set globally outside of `t.Run`. If different subtests need different environment variables set, set them individually within each subtest, not globally.

```go
t.Run(testCase.name, func(t *testing.T) {
    if testCase.environmentValue != nil {
        t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME, *testCase.environmentValue)
    }
    // ... rest of the test ...
})
```

And for `TestUnitGetConfigMultiString_Matrix`, do **not** set environment variables with `t.Setenv` before entering the loop. Instead, set them only as appropriate inside each subtest:

```go
t.Run(testCase.name, func(t *testing.T) {
    if testCase.environmentValues != nil {
        for idx, envVar := range testCase.environmentValues {
            switch idx {
            case 0:
                t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME1, environmentValue1)
            case 1:
                t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME2, environmentValue2)
            // add more cases as needed...
            }
        }
    }
    // ... rest of test
})
```
