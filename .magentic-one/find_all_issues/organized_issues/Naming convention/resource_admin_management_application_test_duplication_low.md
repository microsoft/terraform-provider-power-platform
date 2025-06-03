# Test Function Duplication with Minimal Variation

##

/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/resource_admin_management_application_test.go

## Problem

There are two different test functions with similar names and redundant purposes:

- `TestAccAdminManagementApplicationResource_Validate_Create`
- `TestUnitAdminManagementApplicationResource_Validate_Create`

Both are designed to test creation scenarios—one as an acceptance test (which calls out to actual Terraform provider infrastructure with mocked external provider), and the other as a unit test by using mocked HTTP endpoints.

While it is reasonable to cover both integration and unit testing, the heavy similarity in logic and naming may cause confusion, redundancy, and extra maintenance effort. New contributors may not easily discern the scope or difference between “Acc” and “Unit” prefixes.

## Impact

- **Severity: Low**
- Increases maintenance cost as test logic must be kept in sync in multiple places.
- Confuses readers about the scope and meaning of tests.

## Location

Both test definitions:

```go
func TestAccAdminManagementApplicationResource_Validate_Create(t *testing.T) { ... }
func TestUnitAdminManagementApplicationResource_Validate_Create(t *testing.T) { ... }
```

## Code Issue

```go
func TestAccAdminManagementApplicationResource_Validate_Create(t *testing.T) {
...
}
func TestUnitAdminManagementApplicationResource_Validate_Create(t *testing.T) {
...
}
```

## Fix

- Document clearly (using regular comments) the difference in scope: acceptance vs. unit.
- Consider consolidating any shared setup or validation logic into helper functions.
- Consider naming using Go’s sub-test conventions or unique descriptive names:

```go
func TestAdminManagementApplicationResource(t *testing.T) {
    t.Run("Acceptance_Create", func(t *testing.T) {
        // acceptance logic
    })
    t.Run("Unit_Create", func(t *testing.T) {
        // unit logic
    })
}
```
