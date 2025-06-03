# Misspelled Field Name in Struct

application/models.go

## Problem

The struct field `ApplicationDescprition` in `TenantApplicationPackageDataSourceModel` is misspelled. It should be `ApplicationDescription`.

## Impact

Having a typo in the struct field name can cause confusion, reduce code readability, and may also lead to inconsistent behavior when the code expects the correct field name. This can introduce bugs when interacting with APIs, terraform schema, or refactoring code in the future.  
Severity: Low

## Location

Line defining `ApplicationDescprition` in the struct `TenantApplicationPackageDataSourceModel`.

## Code Issue

```go
	ApplicationDescprition types.String                                   `tfsdk:"application_descprition"`
```

## Fix

Correct the spelling in both the struct field name and the struct tag:

```go
	ApplicationDescription types.String                                   `tfsdk:"application_description"`
```
