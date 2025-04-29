# Issue 1: Inefficient map initialization in `Diff`

##

/workspaces/terraform-provider-power-platform/internal/helpers/array/array.go

## Problem

In the `Diff` function, maps `oldMap` and `newMap` are initialized without using the length of the input arrays (`oldArr` and `newArr`) to preallocate capacity. This can cause a performance overhead due to repeated map resizing during insertion operations.

## Impact

Without preallocating capacity based on the input size, the code may experience performance inefficiencies, especially when dealing with large arrays. This impact, while not severe, could be noticeable in high-performance scenarios.

**Severity:** Medium

## Location

File location: `Diff` function in `/workspaces/terraform-provider-power-platform/internal/helpers/array/array.go`.

## Code Issue

Current code:

```go
	oldMap := make(map[string]bool)
	for _, item := range oldArr {
		oldMap[item] = true
	}
	newMap := make(map[string]bool)
	for _, item := range newArr {
		newMap[item] = true
	}
```

## Fix

Update map initialization to include preallocated capacity using the length of the input arrays. This will improve performance when working with large arrays.

```go
	oldMap := make(map[string]bool, len(oldArr))
	for _, item := range oldArr {
		oldMap[item] = true
	}
	newMap := make(map[string]bool, len(newArr))
	for _, item := range newArr {
		newMap[item] = true
	}
```

This change utilizes the length of the input arrays, reducing unnecessary memory reallocations during the map's growth.
