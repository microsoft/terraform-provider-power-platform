# Title

Potential null pointer dereference in `GetCurrentCloudConfiguration`

##

/workspaces/terraform-provider-power-platform/internal/config/config.go

## Problem

The `configuration` map in the function `GetCurrentCloudConfiguration` uses nested maps to store the cloud type configurations. There is no explicit handling for cases where the requested cloud type or configuration key does not exist, leading to a potential null pointer dereference when trying to access `configuration[string(model.CloudType)][string(key)]`.

## Impact

The null pointer dereference may cause runtime panics, crashing the application and affecting service availability in production environments. This issue rates as high severity due to its potential impact on application stability.

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

Introduce explicit error handling to ensure the function gracefully handles cases where the requested cloud type or configuration key does not exist.

```go
func (model *ProviderConfig) GetCurrentCloudConfiguration(key CloudTypeConfigurationKey) *string {
	configuration := map[string]map[string]*string{
		string(CloudTypePublic): {
			string(FirstReleaseClusterName): helpers.StringPtr("FirstRelease"),
		},
		string(CloudTypeGcc): {
			string(FirstReleaseClusterName): helpers.StringPtr("GovFR"),
		},
	}

	if cloudConfig, exists := configuration[string(model.CloudType)]; exists {
		if keyValue, exists := cloudConfig[string(key)]; exists {
			return keyValue
		}
	}

	return nil // Return nil or appropriate error handling in case of missing key
}
```

This fix ensures that the function gracefully handles missing keys and avoids crashing due to null pointer dereferences.
