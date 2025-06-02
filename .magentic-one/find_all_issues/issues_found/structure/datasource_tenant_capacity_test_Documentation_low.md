# Title

No Documentation for Test Purpose

##

/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity_test.go

## Problem

The test functions lack descriptive comments explaining the intent, expectation, and context of each test scenario.

## Impact

Low. While the code is readable, comments would improve maintainability by making it obvious what aspect of the datasource or API behavior is being verified (especially useful for future contributors).

## Location

Before each test function.

## Code Issue

```go
func TestUnitTenantCapacityDataSource_Validate_Read(t *testing.T) {
...
func TestAccTenantCapacityDataSource_Validate_Read(t *testing.T) {
```

## Fix

Add a doc comment above each test explaining what it verifies.

```go
// TestUnitTenantCapacityDataSource_Validate_Read verifies that the tenant capacity datasource
// reads and unmarshals mock tenant capacity responses as expected.
func TestUnitTenantCapacityDataSource_Validate_Read(t *testing.T) {
...

// TestAccTenantCapacityDataSource_Validate_Read verifies reading real tenant capacity information
// in an acceptance environment using live API calls.
func TestAccTenantCapacityDataSource_Validate_Read(t *testing.T) {
```
