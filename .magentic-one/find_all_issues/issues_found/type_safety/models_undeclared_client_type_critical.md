# Issue: Missing Imports or Type Declaration for `client`

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/models.go

## Problem

The type `client` used in the definitions of the `TenantSettingsDataSource` and `TenantSettingsResource` structs is not declared or imported anywhere in the file.

## Impact

If this type is not declared elsewhere or if not imported, this will result in a compilation failure. Severity: **critical**.

## Location

Anywhere `client` is referenced in structs.

## Code Issue

```go
TenantSettingsClient client
...
TenantSettingClient client
```

## Fix

Ensure that `client` is either defined in this package or imported from the correct location. For example:

```go
// If client is a type in another package, e.g., "github.com/microsoft/terraform-provider-power-platform/internal/client"
import "github.com/microsoft/terraform-provider-power-platform/internal/client"

// Then the struct fields can be updated as:
TenantSettingsClient client.Client
...
TenantSettingsClient client.Client
```
