# Title

Using `require.Contains` Lead to Unnecessary Complexity

## Path

`/workspaces/terraform-provider-power-platform/internal/provider/provider_test.go`

## Problem

`require.Contains` is used within a loop to ensure all registered data sources match `expectedDataSources`. This results in redundant iterations and unnecessary complexity because it fails to use hashing for quick validation.

```go
for _, d := range datasources {
	require.Contains(t, expectedDataSources, d(), "An unexpected data source was registered")
}
```

## Impact

- Performance inefficiency due to repeated linear searches within `expectedDataSources` (O(n^2) complexity for n datasources).
- Potential subtle bugs if the `ExpectedSources` changes and size increases significantly.
- Medium severity since the code still works but at suboptimal performance.

## Location

Line 65 in the file `provider_test.go`

## Code Issue

```go
for _, d := range datasources {
	require.Contains(t, expectedDataSources, d(), "An unexpected data source was registered")
}
```

## Fix

Instead of linear iteration, convert `expectedDataSources` to a set or a map and validate against that. This brings down the complexity to O(n) for n datasources.

```go
// Convert expectedDataSources into a map for fast lookup.
expectedSourcesMap := make(map[datasource.DataSource]bool)
for _, ds := range expectedDataSources {
	expectedSourcesMap[ds] = true
}

// Validate each datasource efficiently.
for _, d := range datasources {
	require.True(t, expectedSourcesMap[d], "An unexpected data source was registered")
}
```