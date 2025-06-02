# Title

Improperly Handled File Close on Error in `CalculateSHA256` Function

##

`/workspaces/terraform-provider-power-platform/internal/helpers/hash.go`

## Problem

The `file.Close()` method is deferred immediately after the file is successfully opened. However, if the `io.Copy` operation fails later in the function, the deferred `file.Close()` may not execute until after an error is propagated, leading to potentially improper resource cleanup or delayed close.

## Impact

This can cause resource leaks, as the file descriptor may remain open longer than intended. While its impact in most environments would be low, in high-frequency or critical systems, resource exhaustion could occur, resulting in operational issues.

*Severity*: **Medium**

## Location

The code can be found here:

```go
file, err := os.Open(filePath)
if err != nil {
    if os.IsNotExist(err) {
        return "", nil
    }
    return "", fmt.Errorf("failed to open file '%s': %w", filePath, err)
}
defer file.Close()
```

## Code Issue

```go
defer file.Close()
```

## Fix

Encapsulate `file.Close()` inside explicit cleanup logic (such as in a `defer` or directly after `io.Copy`). Handle the error conditions for `Close()` during success or failure paths:

```go
func CalculateSHA256(filePath string) (string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        if os.IsNotExist(err) {
            return "", nil
        }
        return "", fmt.Errorf("failed to open file '%s': %w", filePath, err)
    }

    var closeErr error
    defer func() {
        closeErr = file.Close()
    }()

    hash := sha256.New()
    if _, err := io.Copy(hash, file); err != nil {
        return "", fmt.Errorf("failed to read file '%s': %w", filePath, err)
    }

    if closeErr != nil {
        return "", fmt.Errorf("failed to close file '%s': %w", filePath, closeErr)
    }

    return hex.EncodeToString(hash.Sum(nil)), nil
}
```