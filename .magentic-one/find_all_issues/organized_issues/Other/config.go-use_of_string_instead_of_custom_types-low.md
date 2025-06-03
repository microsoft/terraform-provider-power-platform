# Title

Use of string type for CloudType and CloudTypeConfigurationKey instead of custom types

##

/workspaces/terraform-provider-power-platform/internal/config/config.go

## Problem

The code defines `CloudType` and `CloudTypeConfigurationKey` as type aliases for `string`, but later uses them as strongly typed keys for maps, and interconverts them back to `string` via `string(x)`. This can lead to confusion, undermine the benefits of using custom types, and enable silent type conversion bugs.

## Impact

This is a low severity issue, as it does not cause immediate runtime problems but can reduce type safety and make the code harder to maintain or refactor in the future.

## Location

Lines: type declarations and usage in map keys (GetCurrentCloudConfiguration):

```go
type CloudType string
type CloudTypeConfigurationKey string

...

configuration := map[string]map[string]*string{
	string(CloudTypePublic): {
		string(FirstReleaseClusterName): helpers.StringPtr("FirstRelease"),
		// ...
	},
	// ...
}

return configuration[string(model.CloudType)][string(key)]
```

## Code Issue

```go
type CloudType string
type CloudTypeConfigurationKey string

configuration := map[string]map[string]*string{
	string(CloudTypePublic): {
		string(FirstReleaseClusterName): helpers.StringPtr("FirstRelease"),
		// Add more cloud specific configurations here
	},
	// ...
}

return configuration[string(model.CloudType)][string(key)]
```

## Fix

Instead of using string as map keys, use the custom types directly (where possible). If you must convert, document why. 

The more maintainable approach is to use the custom types as map keys:

```go
configuration := map[CloudType]map[CloudTypeConfigurationKey]*string{
	CloudTypePublic: {
		FirstReleaseClusterName: helpers.StringPtr("FirstRelease"),
		// Add more cloud specific configurations here
	},
	CloudTypeGcc: {
		FirstReleaseClusterName: helpers.StringPtr("GovFR"),
	},
	// ...
}

return configuration[model.CloudType][key]
```

This will enforce type safety and improve code readability. If you need to interoperate with string-based APIs, handle conversions only at boundaries, not throughout the code.
