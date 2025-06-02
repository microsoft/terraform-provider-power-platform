# Title

Global maps for static configurations are not synchronized

##

/workspaces/terraform-provider-power-platform/internal/config/config.go

## Problem

The `configuration` variable in the `GetCurrentCloudConfiguration` function is a global map containing nested maps for cloud-specific configurations. When used in concurrent workloads, this code could lead to race conditions if the map is concurrently read and modified.

## Impact

Race conditions can cause the program to behave unexpectedly, leading to significant issues in production environments, including corrupted state, unexpected behavior, or server crashes. This issue rates a high severity since it directly affects runtime reliability.

## Location

Function: `GetCurrentCloudConfiguration`

## Code Issue

```go
func (model *ProviderConfig) GetCurrentCloudConfiguration(key CloudTypeConfigurationKey) *string {
	configuration := map[string]map[string]*string{
		string(CloudTypePublic): {
			string(FirstReleaseClusterName): helpers.StringPtr("FirstRelease"),
			// Add more cloud specific configurations here
		},
		string(CloudTypeGcc): {
			string(FirstReleaseClusterName): helpers.StringPtr("GovFR"),
		},
		// ...other similar configurations omitted for brevity
	}

	return configuration[string(model.CloudType)][string(key)]
}
```

## Fix

To resolve the issue, "static configurations" such as the ones in the `configuration` map should be declared at the package level and wrapped within a `sync.Map` or explicitly locked using `sync.Mutex` to ensure thread safety.

```go
var cloudConfigurations = map[string]map[string]*string{
	string(CloudTypePublic): {
		string(FirstReleaseClusterName): helpers.StringPtr("FirstRelease"),
	},
	string(CloudTypeGcc): {
		string(FirstReleaseClusterName): helpers.StringPtr("GovFR"),
	},
	// Other definitions...
}

func (model *ProviderConfig) GetCurrentCloudConfiguration(key CloudTypeConfigurationKey) *string {
	// Access global thread-safe configurations here
	return cloudConfigurations[string(model.CloudType)][string(key)]
}
```

By declaring the map globally but restricting modification operations, the structure becomes inherently thread-safe for concurrent read operations.
