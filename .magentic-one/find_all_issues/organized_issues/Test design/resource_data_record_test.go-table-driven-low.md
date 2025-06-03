# No Table-Driven or Subtest Organization

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record_test.go

## Problem

The test cases are implemented as a series of giant individual test functions, each covering multiple scenarios within single functions and large configuration blocks, rather than organized using table-driven tests or subtests (which are a Go idiom).

Table-driven and subtest structure improves both readability and maintainability by grouping variations under a single test function and leveraging `t.Run` for scenario variations.

## Impact

- **Readability**: Difficult to see what variants/scenarios exist and which inputs are covered.
- **Test Reporting**: Cannot see fine-grained pass/fail for individual scenarios.
- **Extensibility**: Adding new cases can mean large copy/paste blocks rather than a new table entry.
- **Severity**: Low

## Location

Example (every major function):

```go
func TestAccDataRecordResource_Validate_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		// steps...
	})
}

func TestAccDataRecordResource_Validate_Create(t *testing.T) { ... }
// etc.
```

## Fix

Prefer a single `TestAccDataRecordResource` function with a table of scenarios, using `t.Run` or a table-driven pattern. This makes it easier to extend, and each subtest is reported on its own in `go test` output.

```go
func TestAccDataRecordResource(t *testing.T) {
	tests := []struct{
		name   string
		config string
		checks []resource.TestCheckFunc
	}{
		{
			name: "validate create",
			config: ...,  // test HCL
			checks: ...,
		},
		{
			name: "validate update",
			config: ...,
			checks: ...,
		},
		// more cases...
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				// ...
			})
		})
	}
}
```

This approach is much more scalable for a large test matrix.

---
