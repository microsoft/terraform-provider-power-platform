# Title

Lack of UUID format validation on NewUUIDValue

##

/workspaces/terraform-provider-power-platform/internal/customtypes/uuid.go

## Problem

The function `NewUUIDValue` accepts any string without validating whether the provided value is a valid UUID format. This can potentially allow invalid UUIDs to enter the system, which may lead to bugs or data inconsistencies later.

## Impact

This is a medium-severity type safety and data consistency issue. Storing invalid UUIDs can cause downstream errors, integration issues, or data corruption.

## Location

```go
func NewUUIDValue(value string) UUID {
	return UUID{
		StringValue: basetypes.NewStringValue(value),
	}
}
```

## Code Issue

```go
func NewUUIDValue(value string) UUID {
	return UUID{
		StringValue: basetypes.NewStringValue(value),
	}
}
```

## Fix

Validate the UUID format using a regular expression or a UUID parsing library before allowing the value. You may return an error or handle diagnostics if the value is invalid.

```go
import (
	"regexp"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func NewUUIDValue(value string) UUID {
	if !uuidRegex.MatchString(value) {
		// handle as you see fit, you could return a Null/Unknown UUID
		return UUID{
			StringValue: basetypes.NewStringUnknown(),
		}
	}
	return UUID{
		StringValue: basetypes.NewStringValue(value),
	}
}
```
