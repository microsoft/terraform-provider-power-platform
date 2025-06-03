# Lack of Validation for Return Value Limits in Fuzz Test

##

/workspaces/terraform-provider-power-platform/internal/helpers/hash_fuzz_test.go

## Problem

The fuzz test currently only checks for errors returned from `helpers.CalculateSHA256` but does not validate the correctness or expected characteristics of the hash returned (e.g., non-empty, fixed length for successful reads, or proper handling of invalid input). As a result, a bug in the hash function that returns unexpected hash values would go undetected by this test.

## Impact

Severity: Medium

Without validating the return value, the test can give a false sense of security: the function may not panic or error, but it could still exhibit incorrect behavior (e.g., wrong hash values, empty string on success cases, or returning sensitive error details in the hash value).

## Location

```go
		_, err := helpers.CalculateSHA256(filePath)

		// Ensure the function does not panic and handles errors gracefully
		if err != nil {
			t.Logf("Expected error for input '%s': %v", filePath, err)
		}
```

## Code Issue

```go
		_, err := helpers.CalculateSHA256(filePath)

		// Ensure the function does not panic and handles errors gracefully
		if err != nil {
			t.Logf("Expected error for input '%s': %v", filePath, err)
		}
```

## Fix

Assert properties of the output hash, such as correct length or format, for successful calls. For known-good inputs add checks for exact hash value. For error conditions, ensure returned value is properly documented (e.g., empty string).

```go
		hash, err := helpers.CalculateSHA256(filePath)
		if err != nil {
			t.Logf("Expected error for input '%s': %v", filePath, err)
			if hash != "" {
				t.Errorf("On error, expected empty hash for input '%s', got: %q", filePath, hash)
			}
			return
		}
		// For success, check hash is hex, length is correct (e.g., 64 for SHA-256), etc.
		if len(hash) != 64 {
			t.Errorf("Expected hash length 64, got %d for filePath '%s'", len(hash), filePath)
		}
		// Optionally, assert hash value for known-good inputs.
```
