# Magic Numbers Used In Test Assertions

##

/workspaces/terraform-provider-power-platform/internal/services/environment_templates/datasource_environment_templates_test.go

## Problem

Test case uses hardcoded values for the number of environment templates ("53") rather than describing their intent.

## Impact

Reduces readability and makes refactoring brittle; intent of the magic number is not clear. Severity: Low.

## Location

```go
resource.TestCheckResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.#", "53")
```

## Code Issue

```go
resource.TestCheckResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.#", "53"),
```

## Fix

Assign the count to a named constant:

```go
const expectedEnvironmentTemplatesCount = "53"
...
resource.TestCheckResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.#", expectedEnvironmentTemplatesCount),
```
