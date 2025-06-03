# Unexported Types Limit Reuse and Testing

application/models.go

## Problem

The structs defined in this file (`TenantApplicationPackagesDataSource`, `TenantApplicationPackagesListDataSourceModel`, etc.) use unexported type `client`. If `client` is defined locally and not exported, external packages cannot construct or test these structs directly. This limits testability and code reuse, especially for mocking or dependency injection.

## Impact

Unexported types reduce flexibility, hinder proper testing (such as table-driven or integration testing with mock clients), and restrict the code from being used or extended outside its current package.  
Severity: Medium

## Location

Fields with type `client` in `TenantApplicationPackagesDataSource` and `EnvironmentApplicationPackageInstallResource`.

## Code Issue

```go
	ApplicationClient client
```

## Fix

Make sure the `client` type is exported if it is commonly used across package boundaries:

```go
	ApplicationClient Client
```

Alternatively, ensure mockability by using interfaces, e.g.:

```go
	ApplicationClient ClientInterface
```
