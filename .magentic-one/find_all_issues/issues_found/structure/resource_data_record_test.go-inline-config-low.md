# Large Inline Configuration Blocks Reduce Readability

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record_test.go

## Problem

Each test step often contains a massive multiline configuration string written inline within the test, spanning many lines. This clutters the body of the test logic, makes the logic of the test hard to follow, and leads to hard-to-diff code reviews when either config or test logic needs to change. It's also error-prone if any string interpolation is needed, and makes it very hard to reuse even small pieces of configuration.

## Impact

- **Readability**: Quickly understanding the real purpose and coverage of each test is difficult because of HCL config noise.
- **Refactoring Pain**: Small changes require changing many locations in the code.
- **Reusability**: Inline strings can't be easily reused elsewhere.
- **Severity**: Low

## Location

Most/all test steps with `Config: ` blocks.

## Code Issue

```go
{
	Config: `
	resource "powerplatform_environment" "test_env" {
		display_name     = "` + mocks.TestName() + `"
		location         = "unitedstates"
		// ...
	}

	resource "powerplatform_data_record" "data_record_sample_contact1" {
		// ...
	}
	`,
	// ...
}
```

## Fix

Refactor any large config string blocks into constants or, where variable interpolation is required, template helpers. This makes test logic clearer and configuration maintainable.

```go
const envConfig = `
resource "powerplatform_environment" "test_env" {
	display_name     = "%s"
	location         = "unitedstates"
	// ...
}
`
const contact1Config = `
resource "powerplatform_data_record" "data_record_sample_contact1" {
	// ...
}
`

Config: fmt.Sprintf(envConfig, mocks.TestName()) + contact1Config,
```

If test steps require only small changes, generate config with simple helpers or string replace/template.

---
