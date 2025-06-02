# Long/Complex Test Functions Decrease Maintainability

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record_test.go

## Problem

Many test functions in this file are extremely long and complex, integrating multiple setup blocks, extensive inline resource configuration (as string literals), nested inline configuration, and many inline assertions. This makes the tests difficult to read, hard to review, and prone to errors or duplication. In some cases the size and complexity discourages necessary additions and refactoring, and can hide subtle bugs. This affects maintainability and future extensibility.

## Impact

- **Readability**: It is difficult for future maintainers to understand what each test does, why, and how to adjust or extend it safely.
- **Test Coverage**: Long tests tend to hide untested cases, or can lead to duplicated test logic that is hard to spot.
- **Debuggability**: Failures within large test cases are harder to isolate and debug.
- **Extensibility**: Reusing test setup or behaviors is hard without refactoring.
- **Severity**: Medium (Maintainability, readability, code quality)

## Location

Numerous locations throughout the file, but for example:

## Code Issue

```go
func TestAccDataRecordResource_Validate_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				// ...hundreds of lines of HCL config...
				`,
				ConfigStateChecks: []statecheck.StateCheck{ /* ... */ },
				Check: resource.ComposeAggregateTestCheckFunc( /* ... */ ),
			},
		},
	})
}
```

## Fix

Refactor long tests into smaller, focused, and reusable helper functions. Move repeated mock and configuration setup code into reusable helpers. Where possible, move large string literals or repeated config blocks into constants or fixtures. For instance:

```go
const envConfig = `
resource "powerplatform_environment" "test_env" {
	// ...
}
`

func validAccountConfig(contactResourceName string) string {
	return fmt.Sprintf(`
resource "powerplatform_data_record" "data_record_account" {
	table_logical_name = "account"
	columns = {
		name = "Sample Account"
		primarycontactid = {
			table_logical_name = %s.table_logical_name
			data_record_id = %s.id
		}
	}
}
`, contactResourceName, contactResourceName)
}

func TestAccDataRecordResource_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		// ...
		Steps: []resource.TestStep{
			{
				Config: envConfig + validAccountConfig("powerplatform_data_record.data_record_sample_contact1"),
				// ...
			},
		},
	})
}

// Also extract repeated HTTP mock setups to helpers, e.g.:
func setupEntityDefinitionMocks() {
	httpmock.RegisterResponder("GET", `...`, ...)
	// etc.
}
```

---

Consider splitting out test types into separate files (acceptance/unit), grouping similar test cases, and keeping each to a manageable size for clarity and maintainability.

---
