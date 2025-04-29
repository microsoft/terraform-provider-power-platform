# Title

Improper use of `t.Fatalf` leading to potential abrupt termination in tests

##

`/workspaces/terraform-provider-power-platform/internal/services/environment_wave/resource_environment_wave_test.go`

## Problem

The `loadTestResponse` function uses `t.Fatalf` to handle errors when reading a test response file. Using `t.Fatalf` exits the current test immediately, which might not be necessary or appropriate in cases where multiple test cases need to be executed sequentially.

Using `t.Fatalf` instead of a more graceful error handling method could terminate the test prematurely and obscure other test failures.

## Impact

Abrupt termination of test execution could lead to incomplete test runs and reduce visibility into how the code behaves under different scenarios. This negatively impacts test coverage and debugging efforts.

**Severity:** Medium

## Location

This issue occurs in the `loadTestResponse` function as shown below.

## Code Issue

```go
content, err := os.ReadFile(path)
if err != nil {
    t.Fatalf("Failed to read test response file %s: %v", filename, err)
}
```

## Fix

Refactor the error handling to make it more graceful, enabling subsequent tests to run even if one fails. Returning an error could be a better option instead of terminating the test.

Example fix:

```go
func loadTestResponse(t *testing.T, testFolder string, filename string) (string, error) {
    path := filepath.Join("test", "resource", testFolder, filename)
    content, err := os.ReadFile(path)
    if err != nil {
        t.Logf("Warning: Failed to read test response file %s: %v", filename, err)
        return "", err
    }
    return string(content), nil
}
```

By refactoring the function in this manner, failed reads will be logged as warnings, with the error returned for further handling by the calling code. It ensures that test execution can continue, leading to more robust validation and debugging.