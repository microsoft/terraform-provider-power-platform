# Title

Missing Validation for Errors in Helpers Function Calls

##

/workspaces/terraform-provider-power-platform/internal/helpers/hash_test.go

## Problem

While calling `helpers.CalculateSHA256`, error values returned by the function are often ignored. These errors should be validated, and appropriate assertions should be added in the test cases.

## Impact

Failing to validate errors hides potential bugs and leaves the code vulnerable to failures or unexpected behavior in future versions. Ensuring proper validation avoids silently broken tests and improves reliability and correctness. Severity: High

## Location

The issue is present throughout the test file whenever `helpers.CalculateSHA256` is called without validation of the returned error.

## Code Issue

One example of a problematic code snippet:

```go
	t.Run("TestUnitCalculateSHA256_SameFile", func(t *testing.T) {
		t.Parallel()

		// Test code here
		f1, err := helpers.CalculateSHA256(file1)
		if err != nil {
			t.Fatal(err)
		}

		f1b, err := helpers.CalculateSHA256(file1)
		if err != nil {
			t.Fatal(err)
		}

		if f1 != f1b {
			t.Errorf("Expected %s to equal %s", f1, f1b)
		}
	})
```

## Fix

Add assertions to validate errors returned from helper function calls. Example:

```go
	t.Run("TestUnitCalculateSHA256_SameFile", func(t *testing.T) {
		t.Parallel()

		f1, err := helpers.CalculateSHA256(file1)
		if err != nil {
			t.Fatalf("Error calculating SHA256 for file1: %v", err)
		}

		f1b, err := helpers.CalculateSHA256(file1)
		if err != nil {
			t.Fatalf("Error calculating SHA256 for file1: %v", err)
		}

		if f1 != f1b {
			t.Errorf("Expected %s to equal %s", f1, f1b)
		}
	})
```

This example illustrates how errors should be validated and not ignored to ensure comprehensive test coverage and reliability.