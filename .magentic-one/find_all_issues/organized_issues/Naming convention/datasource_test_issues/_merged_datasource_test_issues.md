# Datasource Test Issues - Merged Issues

## ISSUE 1

# Title

Potentially Redundant or Misleading Test Name "TestUnitTestBillingPoliciesDataSource_Validate_Read"

## 

/workspaces/terraform-provider-power-platform/internal/services/licensing/datasource_billing_policies_test.go

## Problem

The naming of the test `TestUnitTestBillingPoliciesDataSource_Validate_Read` is confusing. By Go's convention, `UnitTest` is often not included in the function name itself but specified via build tags, file names, or test framework configuration (`IsUnitTest: true` is already used in the `resource.TestCase`). Adding `UnitTest` to both the function and the provider factories (`mocks.TestUnitTestProtoV6ProviderFactories`) is redundant and can cause confusion or difficulty with automation and filtering test types.

## Impact

Severity: **low**

Poor or redundant naming can make it harder for contributors and automated tools to differentiate between test types, or filter/select specific tests for running. Inconsistent/verbose naming is a maintainability issue.

## Location

Second test function, line 57:

## Code Issue

```go
func TestUnitTestBillingPoliciesDataSource_Validate_Read(t *testing.T) {
```

## Fix

Rename the test function and supporting types to be more clear and consistent with Go's standard testing idioms. Prefer:

```go
func TestBillingPoliciesDataSource_Unit_Read(t *testing.T) {
    // or, more simply:
    // func TestBillingPoliciesDataSource_Read_Unit(t *testing.T) {

    // setup
}
```
Or, simply rely on the file name (ending with `_test.go`), comments, and the usage of `IsUnitTest: true` rather than having `UnitTest` in the function name.

Also, ensure supporting variables/types are named clearly, e.g., `mocks.TestProtoV6ProviderFactories` for the unit variant.

Apply for whole codebase


---

## ISSUE 2

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
func Test_Unit_TenantCapacityDataSource_Validate_Read(t *testing.T) {
...
func Test_Acc_TenantCapacityDataSource_Validate_Read(t *testing.T) {
```
apply for whole code base


---

# To finish the task you have to 
1. Run linter and fix any issues 
2. Run UnitTest and fix any of failing ones
3. Generate docs 
4. Run Changie

# Changie Instructions
Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```
Where:
- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
