# Skipped Test Disabled Without Explanation in CI

##

/workspaces/terraform-provider-power-platform/internal/services/environment_groups/resource_environment_group_test.go

## Problem

The `TestNoManagementAppPermissionsMessage` test uses `t.Skip` to disable the test, indicating it should only run under certain credentials (an SP without management app permissions). However, this disables verification for a critical error scenario in CI and provides no automation path or guidance for periodically running it.

## Impact

This reduces test coverage, particularly for important error conditions, because the test will never run unless explicitly enabled by a maintainer. The severity is **medium** (loss of automated regression verification).

## Location

```go
func TestNoManagementAppPermissionsMessage(t *testing.T) {
	t.Skip("Skipping test. This should run with an sp that does not have management app permissions.")
	// ...
```

## Code Issue

```go
	t.Skip("Skipping test. This should run with an sp that does not have management app permissions.")
```

## Fix

If possible, convert the skip to a conditional that runs only when the appropriate environment is detected. Or, add a build/tag or environment variable to allow enabling the test in CI:

```go
	if os.Getenv("TEST_NO_MANAGEMENT_APP_PERMISSIONS") != "1" {
		t.Skip("Skipping test: set TEST_NO_MANAGEMENT_APP_PERMISSIONS=1 to enable.")
	}
```

And document in the project how to run this test in these situations.
