# Inconsistent Naming: `PowerAppssClient` should be `PowerAppsClient`

##

/workspaces/terraform-provider-power-platform/internal/services/powerapps/datasource_environment_powerapps.go

## Problem

The field is named `PowerAppssClient` (with a double "s") instead of the expected `PowerAppsClient`. This inconsistency can cause confusion for developers and increases the risk of typos in other parts of the codebase.

## Impact

Severity: Low  
This is a minor naming issue. However, inconsistent naming can reduce code readability and maintainability, and may cause subtle bugs if developers mistakenly use the wrong identifier.

## Location

```go
d.PowerAppssClient = newPowerAppssClient(client.Api)
apps, err := d.PowerAppssClient.GetPowerApps(ctx)
```

## Code Issue

```go
d.PowerAppssClient = newPowerAppssClient(client.Api)
apps, err := d.PowerAppssClient.GetPowerApps(ctx)
```

## Fix

Change all instances of `PowerAppssClient` to `PowerAppsClient`. Also, ensure that the struct declares this property with the correct name.

```go
// Rename struct field and all its references
d.PowerAppsClient = newPowerAppsClient(client.Api)
apps, err := d.PowerAppsClient.GetPowerApps(ctx)
```
