# Title

Insufficient Resource Naming Consistency in Test Resource Definitions

##

internal/services/data_record/datasource_data_record_test.go

## Problem

Within the `BootstrapDataRecordTest` function and elsewhere in test setups, resource names ("contact1", "contact2", etc.) are directly interpolated from function arguments (`name`) but others are hardcoded or use numeric suffixes. For larger test suites or parallelized tests, this makes it easier to have conflicting names, and it is less clear what each resource represents.

## Impact

- Can lead to naming collisions in more complex or parallelized test runs.
- Reduces clarity/intent for readers who may not know the significance of each resource.
- Slightly increases brittleness of test maintenance or reuse.

Severity: Low (mainly a readability and maintainability issue).

## Location

In the output of `BootstrapDataRecordTest(name string) string`:

## Code Issue

```go
display_name     = "` + name + `"
...
columns = {
    ...
    firstname = "contact1"
    lastname  = "contact1"
    ...
}
```

Similar patterns for `contact2`, `contact3`, `"account1"`, `"contact5"` and so on.

## Fix

Use a consistent naming strategy for all test resources. Possible improvements:

- Use the `name` argument as a prefix/suffix for all dynamic names.
- If tests are ever parallelized or dynamically generated, make sure all resources include a run/test-specific identifier.
- For clarity, use a format like:

```go
firstname = sprintf("%s_contact1", name)
lastname  = sprintf("%s_contact1", name)
```

Or generate all test resource names programmatically from the function argument. This aids in mapping resources to test cases and reduces the risk of name collisions.

Save as a structure/naming/readability issue.
