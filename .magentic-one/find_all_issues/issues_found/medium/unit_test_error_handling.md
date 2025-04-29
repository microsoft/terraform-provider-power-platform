# Title

Unit Test Error Handling Could Be Improved

##

/workspaces/terraform-provider-power-platform/internal/helpers/hash_test.go

## Problem

The unit test `TestUnitCalculateSHA256_FileDoesNotExist` does not correctly handle the case where a file does not exist. This test expects an output of `""` for the `CalculateSHA256` function but does not check for an error, making debugging difficult.

## Impact

If the `CalculateSHA256` function's behavior changes or the input file path is modified inadvertently in future versions, the test might incorrectly pass or fail without adequate insight.
Severity: Medium

## Location

The issue is located in the test running under the name `TestUnitCalculateSHA256_FileDoesNotExist`.

## Code Issue

```go
	t.Run("TestUnitCalculateSHA256_FileDoesNotExist", func(t *testing.T) {
		t.Parallel()

		f4, err := helpers.CalculateSHA256(file4)
		if err != nil {
			t.Fatal(err)
		}

		if f4 != "" {
			t.Errorf("Expected %s to be empty", f4)
		}
	})
```

## Fix

The test should include a condition to verify that an error message is properly returned when the file does not exist.

```go
	t.Run("TestUnitCalculateSHA256_FileDoesNotExist", func(t *testing.T) {
		t.Parallel()

		f4, err := helpers.CalculateSHA256(file4)
		if err == nil {
			t.Fatal("Expected an error but got nil")
		}

		if f4 != "" {
			t.Errorf("Expected %s to be empty", f4)
		}
	})
```

This way, the test validates both that an error is raised and that the output string is empty, covering two aspects of the function's behavior.