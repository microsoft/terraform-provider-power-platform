# Title

Lack of Subtests for Better Test Structure and Output

##

/workspaces/terraform-provider-power-platform/internal/services/environment/datasource_environments_test.go

## Problem

All assertions are done in single monolithic test functions. Using Go's t.Run or subtests can help organize assertions, clarify intent, and improve debuggability by making it clear which field or operation failed, especially as more edge cases or additional properties are added. Grouping related assertions as subtests gives more granular output from `go test` and enables easier selective re-running and more readable test output.

## Impact

Severity: Low

- Harder to pinpoint failing assertions when many are grouped in a single function/block.
- Large test functions become cumbersome as new scenarios/fields are added.
- Misses some advanced features of the testing package designed to support complex or data-driven test logic.

## Location

Applies across the main test functions:
- TestAccEnvironmentsDataSource_Basic
- TestUnitEnvironmentsDataSource_Validate_Read

## Code Issue

```go
func TestAccEnvironmentsDataSource_Basic(t *testing.T) {
    resource.Test(t, resource.TestCase{
        // ... All assertions in one block
    })
}

func TestUnitEnvironmentsDataSource_Validate_Read(t *testing.T) {
    resource.Test(t, resource.TestCase{
        // ... All assertions in one block
    })
}
```

## Fix

- Use t.Run for each significant scenario or field cluster:

```go
func TestAccEnvironmentsDataSource_Basic(t *testing.T) {
    t.Run("basic attributes", func(t *testing.T) {
        // main field assertions
    })
    t.Run("dataverse fields", func(t *testing.T) {
        // dataverse sub-assertions
    })
}
```

- Alternatively, use helper assertion functions composed with t.Run or resource.ComposeTestCheckFunc as appropriate, to break up logic and output by intention.
