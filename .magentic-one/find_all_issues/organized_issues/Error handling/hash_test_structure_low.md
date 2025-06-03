# Use of Magic Strings in File Content

##

/workspaces/terraform-provider-power-platform/internal/helpers/hash_test.go

## Problem

The same string literals like `"same"` and `"different"` are used multiple times as file content. This can lead to errors and inconsistencies if used in multiple places or changed in a single location but not others.

## Impact

Low severity. This is a minor maintainability issue, but improving this makes content changes less error-prone and increases overall readability.

## Location

```go
err := os.WriteFile(file1, []byte("same"), 0644)
...
err = os.WriteFile(file2, []byte("same"), 0644)
...
err = os.WriteFile(file3, []byte("different"), 0644)
```

## Code Issue

```go
err := os.WriteFile(file1, []byte("same"), 0644)
```

## Fix

Declare file content as constants at the beginning of the test:

```go
const (
    sameContent      = "same"
    differentContent = "different"
)
```

Then use:

```go
err := os.WriteFile(file1, []byte(sameContent), 0644)
...
err = os.WriteFile(file3, []byte(differentContent), 0644)
```
