# Title

Error Wrapping and Propagation Inconsistency

##

`/workspaces/terraform-provider-power-platform/internal/helpers/hash.go`

## Problem

The error wrapping used in the `CalculateSHA256` function for logging issues (e.g., `fmt.Errorf("failed to open file '%s': %w", filePath, err)` or `fmt.Errorf("failed to read file '%s': %w", filePath, err)`) is inconsistent and does not provide a highly descriptive message with actionable insights. While `%w` is used for wrapping, the custom prefix message is generic and does not help identify the root issue explicitly.

## Impact

Incomplete error propagation details can obscure the origin of the problem, making diagnosis harder when debugging failures. 
- **Severity**: **High**

## Location

The problematic lines are:

```go
return "", fmt.Errorf("failed to open file '%s': %w", filePath, err)
...
return "", fmt.Errorf("failed to read file '%s': %w", filePath, err)
```

## Code Issue

```go
return "", fmt.Errorf("failed to open file '%s': %w", filePath, err)
```

## Fix

Enhance the error comment for clarity. Ensure that each error wrapping reports actionable insights (e.g., which stage failed and why):

```go
if err != nil {
    if os.IsNotExist(err) {
        return "", nil
    }
    return "", fmt.Errorf("CalculateSHA256: error opening file '%s' for hashing: %w", filePath, err)
}

...

if _, err := io.Copy(hash, file); err != nil {
    return "", fmt.Errorf("CalculateSHA256: error processing file '%s' for SHA256 hash computation: %w", filePath, err)
}
```