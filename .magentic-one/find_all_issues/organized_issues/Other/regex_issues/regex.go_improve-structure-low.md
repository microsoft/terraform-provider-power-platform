# Title
Regex Constants Not Grouped Using Custom Type for Clarity

##
/workspaces/terraform-provider-power-platform/internal/helpers/regex.go

## Problem
All regex literal strings are declared as top-level constants. While functional, grouping these related regex constants with a custom type or dedicated struct could enhance maintainability, organization and discoverability.

## Impact
Low. This is a maintainability/readability improvement; it does not cause correctness or runtime issues but could help as the codebase grows, making the regex set easier to reason about and extend.

## Location
Lines 4-13

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
Consider grouping the regexes within a type or using a struct for better organization. For example:
```go
type RegexPatterns struct {
	Guid             string
	GuidOrEmptyValue string
	UrlValidString   string
	ApiId            string
	String           string
	Version          string
	Time             string
	Boolean          string
}

var Regex = RegexPatterns{
	Guid:             "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$",
	GuidOrEmptyValue: "^(?:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})?$",
	UrlValidString:   "(?i)^[A-Za-z0-9-._~%/:/?=]+$",
	ApiId:            "^[0-9a-zA-Z/._]*$",
	String:           "^.*$",
	Version:          "^[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+$",
	Time:             "^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z)$",
	Boolean:          "^(true|false)$",
}
```

This makes regex access (
`Regex.Guid`, etc.) clearer and groups related values.
