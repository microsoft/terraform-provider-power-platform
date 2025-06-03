# Title
Missing Negative Test Cases for File Operations

##
/workspaces/terraform-provider-power-platform/internal/services/solution/resource_solution_test.go

## Problem
Most acceptance and unit tests only check the happy path for file reading and writing. For example, the test cases assume solution files always exist and are readable/writable. There are no explicit tests validating error paths when files are missing or permissions are denied. This omits important failure state verification.

## Impact
Absence of negative test cases means regressions or platform/environment issues could slip through undetected, potentially leading to crashes or subtle bugs in CI/CD pipelines. Severity: **high** for QA coverage.

## Location
Example lines: 27-36, seen throughout tests e.g.

## Code Issue
```go
solutionFileBytes1, err := os.ReadFile(SOLUTION_1_RELATIVE_PATH)
if err != nil {
    t.Fatalf("Failed to read solution file: %s", err.Error())
}
... (other file reads/writes, no tests for fail path)
```

## Fix
Add explicit negative tests for the following situations:
- Solution file does not exist (expect error)
- Cannot write file due to permissions (expect error)
- I/O failure (simulate with a mock or wrong path)

Example:
```go
func Test_ReadNonExistentSolutionFile(t *testing.T) {
    _, err := os.ReadFile("nonexistent_file.zip")
    if err == nil {
        t.Fatalf("expected error when reading nonexistent file, got none")
    }
}
```
Include similar tests for file write operations and settings file errors.
