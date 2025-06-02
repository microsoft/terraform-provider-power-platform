# Title

Repeated validation checks lead to poor readability and maintainability.

# Path

`/workspaces/terraform-provider-power-platform/internal/services/connectors/datasource_connectors_test.go`

## Problem

There are repeated validation checks for each connector resource attribute, making the code harder to maintain and update. Using a helper function or aggregate validation approach can improve readability and maintainability.

## Impact

Poor readability and potential bugs in case of updates to validation logic due to repeated checks across tests. Severity: medium.

## Location

Functions: `TestAccConnectorsDataSource_Validate_Read`, Line: 21-39; `TestUnitConnectorsDataSource_Validate_Read`, Line: 78-114

## Code Issue

```go
resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.description", regexp.MustCompile(helpers.StringRegex)),
resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.display_name", regexp.MustCompile(helpers.StringRegex)),
resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.id", regexp.MustCompile(helpers.ApiIdRegex)),
resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.name", regexp.MustCompile(helpers.StringRegex)),
resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.publisher", regexp.MustCompile(helpers.StringRegex)),
resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.tier", regexp.MustCompile(helpers.StringRegex)),
resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.type", regexp.MustCompile(helpers.ApiIdRegex)),
```

## Fix

Use a helper function for validation checks:

```go
func validateConnector(index int, connectors map[string]string) []resource.TestCheckFunc {
    return []resource.TestCheckFunc{
        resource.TestMatchResourceAttr(fmt.Sprintf("data.powerplatform_connectors.all.connectors.%d.description", index), regexp.MustCompile(helpers.StringRegex)),
        resource.TestMatchResourceAttr(fmt.Sprintf("data.powerplatform_connectors.all.connectors.%d.display_name", index), regexp.MustCompile(helpers.StringRegex)),
        resource.TestMatchResourceAttr(fmt.Sprintf("data.powerplatform_connectors.all.connectors.%d.id", index), regexp.MustCompile(helpers.ApiIdRegex)),
        resource.TestMatchResourceAttr(fmt.Sprintf("data.powerplatform_connectors.all.connectors.%d.name", index), regexp.MustCompile(helpers.StringRegex)),
        resource.TestMatchResourceAttr(fmt.Sprintf("data.powerplatform_connectors.all.connectors.%d.publisher", index), regexp.MustCompile(helpers.StringRegex)),
        resource.TestMatchResourceAttr(fmt.Sprintf("data.powerplatform_connectors.all.connectors.%d.tier", index), regexp.MustCompile(helpers.StringRegex)),
        resource.TestMatchResourceAttr(fmt.Sprintf("data.powerplatform_connectors.all.connectors.%d.type", index), regexp.MustCompile(helpers.ApiIdRegex)),
    }
}

// Usage in tests:
Steps: []resource.TestStep{
    {
        Config: `
data "powerplatform_connectors" "all" {}`,
        Check: resource.ComposeAggregateTestCheckFunc(validateConnector(0, connectors)),
    }
}
```