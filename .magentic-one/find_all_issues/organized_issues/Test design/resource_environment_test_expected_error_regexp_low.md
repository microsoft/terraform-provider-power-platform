# Overly Permissive Regex in Expected Error Testing

##

/workspaces/terraform-provider-power-platform/internal/services/environment/resource_environment_test.go

## Problem

Expected error assertions in many tests use `regexp.MustCompile(".*<error>.*")` (catch-all regex), matching any error that contains the substring. This is overly broad, and may match the wrong error or mask legitimate misbehaviors in the codebase. For example, any error with the substring `InsufficientCapacity_StorageDriven` will pass the check, regardless of whether the actual error is correct for that scenario.

## Impact

- **Severity: Low**
- Errors may go undetected if the wrong error message text is returned but still matches the loose regular expression.
- Can hide bugs if upstream code changes the format or location of the error message, leading to a false sense of test coverage.

## Location

```go
ExpectError: regexp.MustCompile(".*InsufficientCapacity_StorageDriven.*"),
// ...
ExpectError: regexp.MustCompile(".*InvalidDomainName.*"),
```

## Code Issue

```go
ExpectError: regexp.MustCompile(".*InsufficientCapacity_StorageDriven.*"),
```

## Fix

Where possible, use more specific, anchored, or exact error message checks, for example:

```go
ExpectError: regexp.MustCompile(`Error: \"InsufficientCapacity_StorageDriven\"`),
// Or match the whole expected error message or a known unique format.
```

If exact matches are not feasible (i.e., dynamic error text), at least anchor or make the regexes stricter, or better yet, refactor the code to raise structured errors that can be compared directly.
