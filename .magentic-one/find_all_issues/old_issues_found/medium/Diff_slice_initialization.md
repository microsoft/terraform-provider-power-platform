# Issue 2: Inefficient initialization of slices in `Diff`

##

/workspaces/terraform-provider-power-platform/internal/helpers/array/array.go

## Problem

The `addedElements` and `removedElements` slices are initialized via `make`, but their capacity is not preallocated. This results in potential performance issues due to frequent reallocations during appending operations, especially for large inputs.

## Impact

Inefficient slice initialization might degrade performance in contexts with large arrays being passed to the `Diff` function. While the severity is not critical, addressing this issue can lead to better code optimization and memory management.

**Severity:** Medium

## Location

File location: `Diff` function in `/workspaces/terraform-provider-power-platform/internal/helpers/array/array.go`.

## Code Issue

Current code:

```go
	addedElements := make([]string, 0)
	removedElements := make([]string, 0)
```

## Fix

Provide a length hint to the slices by initializing them with a capacity equal to the length of the input arrays.

```go
	addedElements := make([]string, 0, len(newArr))
	removedElements := make([]string, 0, len(oldArr))
```

By preallocating capacity, the slices reduce internal resizing and reallocation during the `append` operations, boosting performance.
