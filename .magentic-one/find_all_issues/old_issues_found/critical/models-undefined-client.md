# Title

Undefined Type `client` in Structs

##

`/workspaces/terraform-provider-power-platform/internal/services/application/models.go`

## Problem

The `ApplicationClient` field in both `TenantApplicationPackagesDataSource` and `EnvironmentApplicationPackageInstallResource` structs refers to an undefined type `client`. This type is not imported or declared anywhere in the file.

## Impact

This results in compilation errors due to undefined types. The functionality dependent on these structs will fail unless the type is defined or imported. Severity: **critical**.

## Location

Structs `TenantApplicationPackagesDataSource` and `EnvironmentApplicationPackageInstallResource`.

## Code Issue

```go
ApplicationClient client
```

## Fix

Import or define the `client` type appropriately.

For example, import the correct package if it exists:
```go
import "github.com/microsoft/terraform-provider-power-platform/internal/client"
```

Or define the `client` type if missing:
```go
type client struct {
    // Define relevant fields required for ApplicationClient
}
```
