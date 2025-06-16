# Unexported Field in Exported Struct

##

/workspaces/terraform-provider-power-platform/internal/services/environment_templates/models.go

## Problem

The `EnvironmentTemplatesDataSource` struct is exported, but its field `EnvironmentTemplatesClient` is unexported (`client`). This inconsistency can hinder users from properly accessing the `EnvironmentTemplatesClient` when using the struct outside of the package, which may cause confusion or restrict embedding/composition as intended.

## Impact

Limited package usability and potential confusion for library consumers. This is a **medium** severity issue because it affects public API consistency but does not break existing functionality.

## Location

```go
type EnvironmentTemplatesDataSource struct {
    helpers.TypeInfo
    EnvironmentTemplatesClient client
}
```

## Code Issue

```go
type EnvironmentTemplatesDataSource struct {
    helpers.TypeInfo
    EnvironmentTemplatesClient client
}
```

## Fix

Change the field name to be exported if intended to be used outside the package, and ensure the type (`client`) is also exported or accessible as needed.

```go
type EnvironmentTemplatesDataSource struct {
    helpers.TypeInfo
    EnvironmentTemplatesClient Client // Ensure 'Client' type is exported and imported
}
```

If `client` should not be exported, consider marking the struct itself as unexported, or document usage constraints clearly.

