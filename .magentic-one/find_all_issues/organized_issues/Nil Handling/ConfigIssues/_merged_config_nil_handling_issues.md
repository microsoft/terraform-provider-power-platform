# Configuration Nil Handling Issues

This document contains all identified nil handling issues related to configuration components in the Terraform Power Platform provider codebase.

## ISSUE 1

<!-- Source: config.go-potential_panic_accessing_map-high.md -->

# Title

Potential panic when accessing configuration map without existence check

##

/workspaces/terraform-provider-power-platform/internal/config/config.go

## Problem

The `GetCurrentCloudConfiguration` method returns a value directly from the nested `configuration` map without checking if the key exists. If `model.CloudType` or `key` do not match any defined configuration, this will cause a runtime panic (due to accessing a map with a missing key and immediately indexing a nil map).

## Impact

Severity: High

This is a high severity issue because it may cause the application to crash unexpectedly if an unsupported CloudType or configuration key is used. This breaks reliability and developer/user trust.

## Location

Method: `GetCurrentCloudConfiguration`

```go
func (model *ProviderConfig) GetCurrentCloudConfiguration(key CloudTypeConfigurationKey) *string {
 configuration := map[string]map[string]*string{
  string(CloudTypePublic): {
   string(FirstReleaseClusterName): helpers.StringPtr("FirstRelease"),
   // Add more cloud specific configurations here
  },
  // ...
 }

 return configuration[string(model.CloudType)][string(key)]
}
```

## Code Issue

```go
return configuration[string(model.CloudType)][string(key)]
```

## Fix

Check for existence of keys before accessing the map, and return `nil` (or an appropriate error) if the configuration does not exist:

```go
func (model *ProviderConfig) GetCurrentCloudConfiguration(key CloudTypeConfigurationKey) *string {
 configuration := map[string]map[string]*string{
  string(CloudTypePublic): {
   string(FirstReleaseClusterName): helpers.StringPtr("FirstRelease"),
  },
  // ...
 }

 cloudConfig, ok := configuration[string(model.CloudType)]
 if !ok {
  return nil
 }
 val, keyOk := cloudConfig[string(key)]
 if !keyOk {
  return nil
 }
 return val
}
```

This avoids a panic and provides a graceful handling of missing configuration.

---

# To finish the task you have to

1. Run linter and fix any issues
2. Run UnitTest and fix any of failing ones
3. Generate docs
4. Run Changie

# Changie Instructions

Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```

Where:

- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number

Apply this fix to the whole codebase
