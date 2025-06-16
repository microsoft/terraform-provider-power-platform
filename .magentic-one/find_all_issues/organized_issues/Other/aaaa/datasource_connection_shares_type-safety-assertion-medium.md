# Type safety risk: interface{} assertion without type check

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connection_shares.go

## Problem

The `ConvertFromConnectionSharesDto` function directly indexes into a field assumed to be `map[string]interface{}`. If the API changes (or the property is malformed), this will panic at runtime.

## Impact

**Severity: medium**

Could lead to runtime panics/crashes, especially with unknown API changes.

## Location

```go
if displayName, ok := connection.Properties.Principal["displayName"].(string); ok {
	share.Principal.DisplayName = types.StringValue(displayName)
} else {
	share.Principal.DisplayName = types.StringValue("")
}

if entraId, ok := connection.Properties.Principal["id"].(string); ok {
	share.Principal.EntraId = types.StringValue(entraId)
} else {
	share.Principal.EntraId = types.StringValue("")
}
```

## Code Issue

```go
if displayName, ok := connection.Properties.Principal["displayName"].(string); ok {
	share.Principal.DisplayName = types.StringValue(displayName)
} else {
	share.Principal.DisplayName = types.StringValue("")
}

if entraId, ok := connection.Properties.Principal["id"].(string); ok {
	share.Principal.EntraId = types.StringValue(entraId)
} else {
	share.Principal.EntraId = types.StringValue("")
}
```

## Fix

Verify type before accessing – use a type assertion for the map and only then index into the fields:

```go
principal := connection.Properties.Principal
var displayName, entraId string

if m, ok := principal.(map[string]interface{}); ok {
	if n, ok := m["displayName"].(string); ok {
		displayName = n
	}
	if eid, ok := m["id"].(string); ok {
		entraId = eid
	}
}
share.Principal.DisplayName = types.StringValue(displayName)
share.Principal.EntraId = types.StringValue(entraId)
```
