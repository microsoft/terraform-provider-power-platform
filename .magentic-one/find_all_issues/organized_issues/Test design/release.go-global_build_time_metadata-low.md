# Title

Use of global variables for build-time metadata

##

/workspaces/terraform-provider-power-platform/common/release.go

## Problem

The file declares `ProviderVersion`, `Commit`, and `Branch` as global variables for storing build-time metadata. While this is a common Go pattern for build-time vars set by `ldflags`, using global variables can potentially make unit testing harder and can lead to issues if the values are accidentally modified elsewhere in the codebase.

## Impact

Severity: low  
Impact: Mostly maintainability and testability. There are no functional errors, but if these values are modified at runtime, it could cause inconsistencies.

## Location

Lines 7–14

## Code Issue

```go
var (
	ProviderVersion = "0.0.0-dev" // Default value for development builds
	Commit          = "dev"       // Default value for development builds
	Branch          = "dev"       // Default value for development builds
)
```

## Fix

You may consider using getter functions with unexported variables instead of exported globals, or marking variables as `const` where possible. However, since these are set via `ldflags`, leave as they are, but add documentation and perhaps code comment warnings so contributors do not write to them at runtime.

```go
// DO NOT MODIFY these variables at runtime – meant to be set only via build process (ldflags).
var (
	ProviderVersion = "0.0.0-dev"
	Commit          = "dev"
	Branch          = "dev"
)
```
