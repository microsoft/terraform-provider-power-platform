# Title

Code Lacks Null or Empty Slice Check in `Contains`

## File

`/workspaces/terraform-provider-power-platform/internal/helpers/contains.go`

## Problem

The `Contains` function does not include a check for whether the provided `slice` parameter is `nil` or empty. If the slice is `nil` and the loop attempt occurs, it may lead to undefined behavior in some cases or unnecessary inefficiencies.

## Impact

- **Severity**: Medium
- Impact includes:
  - Potential risk of runtime issues if relying clients provide a `nil` slice.
  - Unnecessary iteration even when the slice is empty.

## Location

```go
func Contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}
```

## Fix

### Proposed Fix:

You can add an extra check at the beginning of the function to handle `nil` or empty slices immediately. This ensures the function explicitly handles such cases.

```go
func Contains(slice []string, item string) bool {
    // Check if the slice is nil or empty right away
    if slice == nil || len(slice) == 0 {
        return false
    }

    // Iterate through the slice
    for _, element := slice {
        if element == item {
            return true
        }
    }
    return false
}
```