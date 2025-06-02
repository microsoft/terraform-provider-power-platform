# Naming: Field Name Uses Lowercase Type Name

##

/workspaces/terraform-provider-power-platform/internal/services/rest/models.go

## Problem

The struct field `DataRecordClient client` uses a lowercase type (`client`) instead of a properly qualified type. It's unclear whether `client` refers to a type defined in this or another package, and its lowercase name suggests it is not exported, which could impact usability elsewhere.

## Impact

**Severity: Medium**

Ambiguous or non-exported type usage can hinder code understanding, cause confusion about accessibility, and create possible bugs if exported struct fields don't use exported types.

## Location

```go
type DataverseWebApiDatasource struct {
    helpers.TypeInfo
    DataRecordClient client
}
```

## Code Issue

```go
    DataRecordClient client
```

## Fix

Ensure that `client` is an exported type, properly qualified with its package, and appropriately named for clarity. For example, if `client` is defined in this package and is intended to be used here, it should be `Client`.

```go
    DataRecordClient Client
```

Or if it comes from another package, qualify with the package name:

```go
    DataRecordClient rest.Client
```

Replace all usages and definitions accordingly.

