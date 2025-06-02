# Title

Potential Issue with `TypeInfo` Struct Implementation

##

/workspaces/terraform-provider-power-platform/internal/helpers/typeinfo.go

## Problem

The `FullTypeName` method does not handle situations where either `t.TypeName` or `t.ProviderTypeName` inadvertently contains spaces or special characters. This could lead to unexpected behaviors or inconsistencies when generating full type names.

## Impact

If `TypeName` or `ProviderTypeName` includes invalid characters or formats, the resulting `FullTypeName` could cause errors downstream in the code where consistency in formatting is expected. This could be of medium severity because it would only result in issues in specific edge cases unless validated elsewhere.

## Location

`TypeInfo` struct and its `FullTypeName()` method.

## Code Issue

```go
type TypeInfo struct {
	ProviderTypeName string
	TypeName         string
}

// FullTypeName returns the full type name in the format provider_type.
func (t *TypeInfo) FullTypeName() string {
	if t.ProviderTypeName == "" {
		return fmt.Sprintf("powerplatform_%s", t.TypeName)
	}

	return fmt.Sprintf("%s_%s", t.ProviderTypeName, t.TypeName)
}
```

## Fix

To ensure proper formatting and avoid potential issues, you can employ validation or sanitization functions to ensure the consistency of `ProviderTypeName` and `TypeName` before formatting.

```go
import (
	"fmt"
	"regexp"
)

// sanitizeName ensures that the provided name contains only valid characters.
func sanitizeName(name string) string {
	// Replace any space or special characters with underscores
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	return re.ReplaceAllString(name, "_")
}

// FullTypeName returns the full type name in the format provider_type.
func (t *TypeInfo) FullTypeName() string {
	sanitizedTypeName := sanitizeName(t.TypeName)
	sanitizedProviderTypeName := sanitizeName(t.ProviderTypeName)

	if sanitizedProviderTypeName == "" {
		return fmt.Sprintf("powerplatform_%s", sanitizedTypeName)
	}

	return fmt.Sprintf("%s_%s", sanitizedProviderTypeName, sanitizedTypeName)
}
```

By using the `sanitizeName` function, the `FullTypeName` will always produce a consistent, valid format, reducing the probability of errors.
