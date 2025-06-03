# Title
Temporary Test Files Are Not Cleaned Up

##
/workspaces/terraform-provider-power-platform/internal/services/solution/resource_solution_test.go

## Problem
Test functions create several files directly in the working directory (e.g., solution and settings files) using `os.WriteFile` and `os.Create`. These files are not deleted after tests, which can pollute the workspace and affect subsequent test runs.

## Impact
Persistent test files may result in unexpected test behavior on repeated runs, increase workspace size, and possibly leak sensitive information. Severity: **medium**, as it affects local/dev workflow, CI hygiene, and repeatability.

## Location
Multiple test functions (examples: lines 28, 36, 86, 93, etc.)

## Code Issue
```go
err = os.WriteFile(SOLUTION_1_NAME, solutionFileBytes1, 0644)
if err != nil {
    t.Fatalf("Failed to write solution file: %s", err.Error())
}
...
file, err := os.Create(fileName)
// ...
defer file.Close()
_, err = file.Write([]byte(content))
```

## Fix
Ensure that files written during tests are removed using `defer os.Remove()` after their creation, or use temporary directories/files via `os.WriteTemp` or `t.TempDir()` if possible. Example:

```go
name := filepath.Join(t.TempDir(), "test_solution.zip")
err := os.WriteFile(name, solutionFileBytes, 0644)
if err != nil {
    t.Fatalf("Failed to write solution file: %s", err.Error())
}
// automatically deleted after test
```
For files created outside `t.TempDir()`, add:
```go
defer os.Remove(name)
```
where relevant, close to the file’s creation.
