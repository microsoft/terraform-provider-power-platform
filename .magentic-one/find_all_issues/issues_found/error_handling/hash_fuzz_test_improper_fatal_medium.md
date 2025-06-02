# Improper use of `t.Fatal` in Fuzz Test Leads to Crash

##

/workspaces/terraform-provider-power-platform/internal/helpers/hash_fuzz_test.go

## Problem

The call to `f.Fatal(err)` inside the initial setup of the fuzz test may cause the test to exit immediately upon error, rather than providing diagnostic information for fuzzing. Fuzz tests are meant to continue running and gather as much information as possible, not to abort abruptly. Use `t.Fatalf` and move this to within the fuzz function for better test resilience.

## Impact

Severity: Medium

Abruptly halting the fuzz test may mask underlying issues and prevent the execution of further fuzzing, reducing the effectiveness of fuzz detection and coverage.

## Location

```go
	err := os.WriteFile(expected, []byte("same"), 0644)
	if err != nil {
		f.Fatal(err)
	}
```

## Code Issue

```go
	err := os.WriteFile(expected, []byte("same"), 0644)
	if err != nil {
		f.Fatal(err)
	}
```

## Fix

Move the write operation and the associated error handling into the fuzzing function so that errors can be handled on a per-test-case basis without aborting all fuzz runs.

```go
func FuzzCalculateSHA256(f *testing.F) {
    tmp := f.TempDir()
    expected := tmp + "/test.txt"

    // Add initial seed corpus
    f.Add(expected)
    // ...add other cases...

    f.Fuzz(func(t *testing.T, filePath string) {
        // Write file only for the specific expected case
        if filePath == expected {
            if err := os.WriteFile(expected, []byte("same"), 0644); err != nil {
                t.Fatalf("Failed to write test file: %v", err)
            }
        }

        // Call the function with the fuzzed input
        _, err := helpers.CalculateSHA256(filePath)

        // Ensure the function does not panic and handles errors gracefully
        if err != nil {
            t.Logf("Expected error for input '%s': %v", filePath, err)
        }
    })
}
```
