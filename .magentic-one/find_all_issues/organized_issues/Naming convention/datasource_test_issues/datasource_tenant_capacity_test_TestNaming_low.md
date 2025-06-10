# Title

Test Function Naming Inconsistency

##

/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity_test.go

## Problem

The test function name `TestUnitTenantCapacityDataSource_Validate_Read` partially describes the test, but the pattern is inconsistently applied; the second test uses `TestAccTenantCapacityDataSource_Validate_Read` (for acceptance tests). If the naming convention relies on `Unit` or `Acc` as prefixes to distinguish test types, this is fine, but should be documented and consistent across the whole package/test suite.

## Impact

Low. This is a readability and maintainability issue. An inconsistent or unclear naming scheme makes it more difficult for developers to identify the purpose/nature of a test by its function name alone.

## Location

Line 10 and line 56.

## Code Issue

```go
func TestUnitTenantCapacityDataSource_Validate_Read(t *testing.T) {
...
func TestAccTenantCapacityDataSource_Validate_Read(t *testing.T) {
```

## Fix

Establish and document a clear testing naming convention (i.e., TestUnit*, TestAcc*, or TestX_Unit, TestX_Acc, etc.). If following the existing scheme, just ensure strict consistency.

```go
// For clarity and consistency, perhaps document that each test follows the format:
// Test<Type><Resource><Operation>
func TestUnitTenantCapacityDataSource_ValidateRead(t *testing.T) {
...
func TestAccTenantCapacityDataSource_ValidateRead(t *testing.T) {
```

Or, if keeping underscores:

```go
func Test_Unit_TenantCapacityDataSource_Validate_Read(t *testing.T) {
...
func Test_Acc_TenantCapacityDataSource_Validate_Read(t *testing.T) {
```
