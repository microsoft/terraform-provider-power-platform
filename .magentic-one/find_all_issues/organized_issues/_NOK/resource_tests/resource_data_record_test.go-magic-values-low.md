# Inconsistent Use of Table/Resource Names and Identifiers

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record_test.go

## Problem

The file mixes the use of hardcoded UUIDs, table names, and resource references (e.g., `"00000000-0000-0000-0000-000000000001"` versus referential expressions like `powerplatform_data_record.data_record_sample_contact1.id` and `table_logical_name = "contact"` versus referencing another resource's attribute. There is no abstraction or helper to ensure these references are consistent across tests, which increases the risk of subtle errors and magic-value bugs. It also makes tests tightly coupled to the literal values.

## Impact

- **Maintainability**: Changing a logical table naming convention or UUID means many edits throughout the file.
- **Readability**: Hard to quickly see which values are meant to be static and which are dynamically referenced.
- **Portability**: Reusing or moving parts of a config or resource definition between tests is harder.
- **Severity**: Low

## Location

Throughout the test file, for example:

```go
environment_id     = "00000000-0000-0000-0000-000000000001"
table_logical_name = "contact"
columns = {
  // ...
}
```
and, elsewhere,
```go
environment_id     = powerplatform_environment.test_env.id
table_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
```

## Code Issue

Use of hardcoded UUIDs and inconsistent referencing.

## Fix

Encapsulate UUIDs, logical table names, and resource references in config helpers or constants wherever possible. Optionally, introduce helper functions/constants where IDs or names are reused, and/or consider leveraging test fixtures to inject the required values.

```go
const envID = "00000000-0000-0000-0000-000000000001"
const contactTableName = "contact"

// Usage:
environment_id     = envID
table_logical_name = contactTableName

// Or dynamically in test helpers/config factories
```

This ensures consistency and helps make changes in one place if requirements evolve.

---
