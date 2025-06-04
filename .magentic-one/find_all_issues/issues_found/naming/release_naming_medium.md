# Unexported Variables Naming Convention

## 
/workspaces/terraform-provider-power-platform/common/release.go

## Problem

The variables `Commit` and `Branch` are exported (start with a capital letter) but do not have a package documentation comment explaining their purpose or intended usage, given that they are set during build processes. As package-level exported variables, Go conventions recommend adding documentation. Furthermore, these are only used for internal build purposes and it is unclear if they are meant for public API consumption.

## Impact

Medium. Not documenting exported identifiers can lead to confusion for users of the package and can impact code maintainability, especially if someone tries to use/modify these outside of intended scenarios.

## Location

Lines 8-12

## Code Issue

```go
	ProviderVersion = "0.0.0-dev" // Default value for development builds
	Commit          = "dev"       // Default value for development builds
	Branch          = "dev"       // Default value for development builds
```

## Fix

Add package comments for each exported variable explaining their purpose, or unexport (make lowercase) if they are not intended for public use.

```go
// ProviderVersion is the version of the released provider, set during build/release process via ldflags.
// Default: "0.0.0-dev" for development builds.
var ProviderVersion = "0.0.0-dev"

// Commit is the git commit hash set during the build/release process via ldflags.
// Default: "dev" for development builds.
var Commit = "dev"

// Branch is the git branch name set during the build/release process via ldflags.
// Default: "dev" for development builds.
var Branch = "dev"
```
Or, if not for export, use:

```go
var providerVersion = "0.0.0-dev"
var commit = "dev"
var branch = "dev"
```
And update usage throughout the codebase accordingly.
