# Title
Use of panic in Helper Function createFile

##
/workspaces/terraform-provider-power-platform/internal/services/solution/resource_solution_test.go

## Problem
The helper function `createFile` uses `panic(err)` for error handling if file creation or writing fails. While idiomatic in test helpers for unexpected failures, it's considered less controlled than returning the error or failing the test directly.

## Impact
`panic` will abruptly halt the test process and may not show up in test results as clearly as using `t.Fatalf`. Stack traces might be less clear, and panics can mask errors in parallel test scenarios. Severity: **medium** (test code robustness, debuggability).

## Location
End of the file:

## Code Issue
```go
func createFile(fileName string, content string) string {
    file, err := os.Create(fileName)

    if err != nil {
        panic(err)
    }

    defer file.Close()

    _, err = file.Write([]byte(content))
    if err != nil {
        panic(err)
    }

    fileChecksum, _ := helpers.CalculateSHA256(fileName)
    return fileChecksum
}
```

## Fix
Refactor the helper to accept a `t *testing.T` parameter, and use `t.Fatalf` for failures instead:

```go
func createFile(t *testing.T, fileName, content string) string {
    file, err := os.Create(fileName)
    if err != nil {
        t.Fatalf("failed to create file %s: %v", fileName, err)
    }
    defer file.Close()
    if _, err := file.Write([]byte(content)); err != nil {
        t.Fatalf("failed to write to file %s: %v", fileName, err)
    }
    fileChecksum, _ := helpers.CalculateSHA256(fileName)
    return fileChecksum
}
```
Update all callers to pass `t` as the first argument.
