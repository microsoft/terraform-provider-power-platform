# Test data path construction is not robust

##

/workspaces/terraform-provider-power-platform/internal/services/environment_wave/resource_environment_wave_test.go

## Problem

The `loadTestResponse` function constructs its file path using constant, hardcoded directory names ("test", "resource", etc.) without checking if these directories exist. This can lead to brittle tests that fail unexpectedly when the directory structure changes, or when running in different environments or test runners. There is no clear handling for cases where the constructed path is incorrect due to relative path differences in various environments.

## Impact

Medium. Inconsistently constructed file paths can cause tests to fail with confusing error messages (`Failed to read test response file ...`), especially when run from different locations or CI environments. Test maintainability and debugging become harder.

## Location

```go
func loadTestResponse(t *testing.T, testFolder string, filename string) string {
	path := filepath.Join("test", "resource", testFolder, filename)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test response file %s: %v", filename, err)
	}
	return string(content)
}
```

## Fix

Establish a configurable test data root, handle missing files more gracefully, and allow for explicit base paths (possibly via an environment variable). Example fix:

```go
func loadTestResponse(t *testing.T, testFolder string, filename string) string {
	basePath := os.Getenv("TESTDATA_ROOT")
	if basePath == "" {
		basePath = "test/resource"
	}
	path := filepath.Join(basePath, testFolder, filename)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test response file %s: %v\nTried path: %s", filename, err, path)
	}
	return string(content)
}
```

This fix allows the test data root path to be overridden and improves diagnostics.
