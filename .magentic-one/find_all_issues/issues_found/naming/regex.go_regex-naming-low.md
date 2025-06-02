# Title
Regex Constant Naming Lacks Standard `RE` Suffix

##
/workspaces/terraform-provider-power-platform/internal/helpers/regex.go

## Problem
The regex constants are named with the `Regex` suffix (e.g., `GuidRegex`). While this is readable, Go conventionally uses the `RE` suffix (uppercase, e.g., `GuidRE`) for regular expressions. This improves consistency across Go codebases and aids quick identification.

## Impact
Low. This is a readability and convention adherence issue. It does not affect runtime or behavior but unifies with community standard best practices for naming regex patterns in Go.

## Location
Lines 5-13

## Code Issue
```go
const (
	GuidRegex             = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
	GuidOrEmptyValueRegex = "^(?:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})?$"
	UrlValidStringRegex   = "(?i)^[A-Za-z0-9-._~%/:/?=]+$"
	ApiIdRegex            = "^[0-9a-zA-Z/._]*$"
	StringRegex           = "^.*$"
	VersionRegex          = "^[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+$"
	TimeRegex             = "^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z)$"
	BooleanRegex          = "^(true|false)$"
)
```

## Fix
Rename constants to use the `RE` suffix as per Go idioms, for example:
```go
const (
	GuidRE             = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
	GuidOrEmptyValueRE = "^(?:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})?$"
	UrlValidStringRE   = "(?i)^[A-Za-z0-9-._~%/:/?=]+$"
	ApiIdRE            = "^[0-9a-zA-Z/._]*$"
	StringRE           = "^.*$"
	VersionRE          = "^[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+$"
	TimeRE             = "^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z)$"
	BooleanRE          = "^(true|false)$"
)
```
