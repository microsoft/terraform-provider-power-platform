# Issue 4

Test Function Naming Convention Not Consistent

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connections_test.go

## Problem

The test functions are named inconsistently regarding their test category. The first test uses `TestAccConnectionsDataSource_Validate_Read` (implying acceptance test convention), while the second uses `TestUnitConnectionsDataSource_Validate_Read` (unit test). However, both functions are in the same file and use similar logic and dependencies.

Go’s standard test naming conventions usually recommend using `TestXxx` and clearly annotating or placing acceptance (`Acc`) and unit tests in their respective files/folders; the difference should be systematic and unambiguous, especially when running with test tags.

## Impact

- **Severity:** Low
- Test discovery, filtering, and categorization will be less straightforward.
- Contributors may be confused where to add/edit new tests.
- Harder to run only acceptance or only unit tests reliably by convention.

## Location

```go
func TestAccConnectionsDataSource_Validate_Read(t *testing.T) {
...
func TestUnitConnectionsDataSource_Validate_Read(t *testing.T) {
```

## Fix

Align test function naming with common conventions, e.g.:

```go
func TestAcc_ConnectionsDataSource_ValidateRead(t *testing.T) { ... }
func TestUnit_ConnectionsDataSource_ValidateRead(t *testing.T) { ... }
```
Alternatively, split acceptance and unit tests into different files or packages, 
e.g., `datasource_connections_acc_test.go` and `datasource_connections_unit_test.go`, 
and use consistent prefixes/suffixes to allow selective runs.
