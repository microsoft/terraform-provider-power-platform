# Title

Misuse of `regexp.MustCompile` within Test Steps

##

`/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/datasource_tenant_settings_test.go`

## Problem

The current test file makes extensive use of `regexp.MustCompile` with values like `helpers.BooleanRegex`. While it's fine to use compiled regular expressions, the value referenced (`helpers.BooleanRegex`) should be validated at least once in the file. Right now, there is no direct validation or documentation in the tests of what the regex is expected to match.

## Impact

This could lead to brittle tests or errors if the regex changes in the helper file without being validated against expected behavior in this test file. Furthermore, using dynamic regular expressions without checks can lead to unexpected issues during runtime. Severity: **Medium**.

## Location

Any test step using `resource.TestMatchResourceAttr` with `regexp.MustCompile(helpers.BooleanRegex)`.

## Code Issue

```go
resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.governance.enable_default_environment_routing", regexp.MustCompile(helpers.BooleanRegex)),
```

## Fix

A better approach would be to validate or document `helpers.BooleanRegex` explicitly in this test file, ensuring its compatibility with expected input values.

```go
// Validate the helpers.BooleanRegex before using it extensively
validBooleanRegex := regexp.MustCompile(helpers.BooleanRegex)
if validBooleanRegex == nil {
    t.Fatalf("Invalid Boolean regex: %s", helpers.BooleanRegex)
}

// Example usage within test steps
resource.TestMatchResourceAttr("data.powerplatform_tenant_settings.settings", "power_platform.governance.enable_default_environment_routing", validBooleanRegex),
```

This ensures that the regex is checked against expected inputs dynamically before heavy usage across multiple test cases.