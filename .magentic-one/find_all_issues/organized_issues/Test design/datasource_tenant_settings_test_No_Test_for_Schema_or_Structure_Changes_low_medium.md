# No Test for Schema or Structure Changes

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/datasource_tenant_settings_test.go

## Problem

The current tests check a large set of fixed attributes, but do not provide a mechanism to easily maintain or assert changes when the data source schema is updated (attributes added, renamed, or removed). This can lead to missing attribute regressions or incomplete testing when the schema evolves, because changes can be made in the implementation without corresponding updates to the tests.

## Impact

- **Severity:** Low-Medium
- This can cause silent test omissions for new or removed attributes, reducing test reliability and maintainability.
- Over time, diverging schema and test coverage may erode confidence in the resource's correctness.

## Location

Across the test attribute assertions in both test functions.

## Code Issue

```go
// Long flat list of assertions, but no structural check or summary validation
resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "disable_capacity_allocation_by_environment_admins", regexp.MustCompile(helpers.BooleanRegex)),
// ... many more
```

## Fix

Consider adding a dynamic check which validates:
- All expected schema attributes are checked (e.g., by maintaining a slice of attribute names and iterating over it)
- Or, use a `resource.TestCheckNoResourceAttr` call for unexpected attributes

Example refactor:

```go
expectedAttrs := []string{
    "disable_capacity_allocation_by_environment_admins",
    "disable_environment_creation_by_non_admin_users",
    // ...
}

for _, attr := range expectedAttrs {
    t.Run(attr, func(t *testing.T) {
        resource.TestCheckResourceAttr("data.powerplatform_tenant_settings.settings", attr, "some_expected_value")
    })
}
```

Or, if using regex, ensure new attributes get tested by auto-expanding the list.

---

This will help future-proof your tests and quickly reveal untested additions or removals.
