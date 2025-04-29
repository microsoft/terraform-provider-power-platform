# Title

Improper Regular Expression for Attribute Validation in Tests

##

`internal/services/environment_templates/datasource_environment_templates_test.go`

## Problem

In the acceptance test, the regular expression `helpers.StringRegex` is used for validating multiple attributes (such as `id`, `name`, `display_name`, `category`) but doesn't account for specific structures and constraints that these values might have—making the test overly broad and reducing validation accuracy.

## Impact

- **Severity**: Medium
- Reduces test precision, allowing potentially invalid data to pass validations unnoticed.
- May lead to undetected issues in production since the test is not verifying the validity of specific attributes correctly.

## Location

**File:**
`/workspaces/terraform-provider-power-platform/internal/services/environment_templates/datasource_environment_templates_test.go`

**Code Block:**
```go
resource.TestMatchResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.0.id", regexp.MustCompile(helpers.StringRegex)),
resource.TestMatchResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.0.name", regexp.MustCompile(helpers.StringRegex)),
resource.TestMatchResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.0.display_name", regexp.MustCompile(helpers.StringRegex)),
resource.TestMatchResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.0.category", regexp.MustCompile(helpers.StringRegex)),
```

## Code Issue

```go
regexp.MustCompile(helpers.StringRegex)
```

## Fix

Update the regular expressions to represent the actual format of each attribute to ensure more comprehensive testing. For example:

```go
resource.TestMatchResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.0.id", regexp.MustCompile(`^/providers/Microsoft.BusinessAppPlatform/locations/[^/]+/environmentTemplates/[^/]+$`)),
resource.TestMatchResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.0.name", regexp.MustCompile(`^[A-Za-z0-9_]+$`)),
resource.TestMatchResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.0.display_name", regexp.MustCompile(`^\w[\w\s]*$`)),
resource.TestMatchResourceAttr("data.powerplatform_environment_templates.all_environment_templates_for_unitedstates", "environment_templates.0.category", regexp.MustCompile(`^(developer|production|sandbox)$`)),
```

These specific validations will prevent unintended matches and improve the reliability of the test cases.
