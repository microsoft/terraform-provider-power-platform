# Incomplete or Overly Generic Resource Attribute Checks in Tests

##

/workspaces/terraform-provider-power-platform/internal/services/powerapps/datasource_environment_powerapps_test.go

## Problem

The tests contain checks for only specific attribute instances (such as `powerapps.0`, `powerapps.2`), and do not iterate or validate all returned results in a generalized/parameterized way. For example, if more powerapps are returned, hard-coded indexes in checks will not scale or validate all items, potentially missing issues if ordering or presence changes.

## Impact

- **Test quality**: Medium – Rigid checks may cause false positives or negatives if the underlying data changes (e.g., order, count).
- **Extensibility**: Makes it harder to extend/add new test cases or validate a variable number of results.

## Location

```go
resource.TestCheckResourceAttr("data.powerplatform_environment_powerapps.all", "powerapps.0.name", "00000000-0000-0000-0000-000000000001"),
// ... others ...
resource.TestCheckResourceAttr("data.powerplatform_environment_powerapps.all", "powerapps.2.name", "00000000-0000-0000-0000-000000000002"),
```

## Code Issue

```go
resource.TestCheckResourceAttr("data.powerplatform_environment_powerapps.all", "powerapps.0.name", "00000000-0000-0000-0000-000000000001"),
resource.TestCheckResourceAttr("data.powerplatform_environment_powerapps.all", "powerapps.2.name", "00000000-0000-0000-0000-000000000002"),
```

## Fix

If possible, loop over the expected list to programmatically check all expected items, or decouple the check logic to increase flexibility and maintainability. For example:

```go
expectedApps := []struct {
    index int
    name string
    id string
    displayName string
    createdTime string
}{
    {0, "00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000001", "Overview", "2023-09-27T07:08:47.1964785Z"},
    {2, "00000000-0000-0000-0000-000000000002", "00000000-0000-0000-0000-000000000002", "Overview", "2023-09-27T07:08:47.1964785Z"},
}

for _, app := range expectedApps {
    t.Run(fmt.Sprintf("Check App Index %d", app.index), func(t *testing.T) {
        resource.TestCheckResourceAttr(fmt.Sprintf("data.powerplatform_environment_powerapps.all", "powerapps.%d.name", app.index), app.name)
        resource.TestCheckResourceAttr(fmt.Sprintf("data.powerplatform_environment_powerapps.all", "powerapps.%d.id", app.index), app.id)
        resource.TestCheckResourceAttr(fmt.Sprintf("data.powerplatform_environment_powerapps.all", "powerapps.%d.display_name", app.index), app.displayName)
        resource.TestCheckResourceAttr(fmt.Sprintf("data.powerplatform_environment_powerapps.all", "powerapps.%d.created_time", app.index), app.createdTime)
    })
}
```
