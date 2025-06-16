# Manual Resource Dependency Management Instead of Using Terraform's Built-in Mechanisms

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record_test.go

## Problem

Some resources use the `depends_on` HCL argument with resources that reference other resources' attributes (such as `table_logical_name` or `id`), but, in Terraform, any explicit attribute reference already creates a dependency. The use of `depends_on` in these cases may signal either cargo-culted or misunderstood dependency management or is simply unnecessary.

## Impact

- **Redundant Code**: Clutters configurations and creates maintenance overhead.
- **Maintainability**: Gives false impression that dependencies are more complex than they really are.
- **Potential for Inconsistency**: If the real attribute references change, the `depends_on` list can become out-of-sync or misleading.
- **Severity**: Low

## Location

For example, in:

```go
resource "powerplatform_data_record" "data_record_sample_contact2" {
	environment_id     = powerplatform_environment.test_env.id
	table_logical_name = "contact"
	columns = {
      contactid = "00000000-0000-0000-0000-000000000020"
	  firstname          = "contact2"
	}

	depends_on = [powerplatform_data_record.data_record_sample_contact1]
}
```

And similar blocks.

## Code Issue

Redundant use of `depends_on` with resources when attribute references (such as `data_record_sample_contact1.id`) already express necessary order.

## Fix

Remove any `depends_on` arguments where there is already a direct attribute dependency:

```go
resource "powerplatform_data_record" "data_record_sample_contact2" {
	environment_id     = powerplatform_environment.test_env.id
	table_logical_name = "contact"
	columns = {
      contactid = "00000000-0000-0000-0000-000000000020"
	  firstname          = "contact2"
	}

	// depends_on not needed if referencing resource attributes elsewhere
}
```

This results in cleaner config and avoids misleading implicit ordering.

---
