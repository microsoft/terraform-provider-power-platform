# Magic Strings and Hardcoded Literals - Data Source Test Issues

This document consolidates all magic strings and hardcoded literals issues found in data source test files.


## ISSUE 1

# Hardcoded Magic Strings in Test Configuration

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/datasource_analytics_data_exports_test.go

## Problem

Several resource attribute names, resource paths, and test data are hardcoded as strings throughout the test, both in the test configuration and in attribute checks. Using magic strings hampers maintainability and increases the risk of typos or mismatches if the schema evolves.

## Impact

Medium. Maintenance cost rises, and silent test failures can occur if underlying resource or provider attribute names change. This risk is amplified in a plugin ecosystem where provider schemas may evolve.

## Location

```go
resource.TestMatchResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.id", regexp.MustCompile(helpers.GuidRegex))
```

## Code Issue

```go
resource.TestMatchResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.id", regexp.MustCompile(helpers.GuidRegex)),
resource.TestMatchResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.source", regexp.MustCompile(helpers.StringRegex)),
// ...etc.
```

## Fix

Define constants for attribute names and resource paths in a dedicated section or file. This centralizes modifications and reduces risk when renaming or refactoring schema fields.

```go
const (
    dataSourceName = "data.powerplatform_analytics_data_exports.test"
    attrExportsID = "exports.0.id"
    attrExportsSource = "exports.0.source"
    // ...etc.
)

resource.TestMatchResourceAttr(dataSourceName, attrExportsID, regexp.MustCompile(helpers.GuidRegex)),
resource.TestMatchResourceAttr(dataSourceName, attrExportsSource, regexp.MustCompile(helpers.StringRegex)),
// ...etc.
```

---

## ISSUE 2

# Title

Hardcoded File Paths Reduce Test Portability

##

/workspaces/terraform-provider-power-platform/internal/services/connectors/datasource_connectors_test.go

## Problem

File paths in the HTTP mock responders are hardcoded. This can break tests when directory structures change, or if the tests are run from a location where the paths are invalid.

## Impact

Reduces test portability and flexibility, leading to fragile tests. Severity: medium.

## Location

Lines 46–62

## Code Issue

```go
httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/connectors/metadata/virtual`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Read/get_virtual.json").String()), nil
    })
```

## Fix

Store files in locations relative to the test file or use test helpers to resolve absolute paths, and document the requirement.

```go
path := filepath.Join("tests", "Validate_Read", "get_virtual.json")
data := httpmock.File(path).String()
return httpmock.NewStringResponse(http.StatusOK, data), nil
```

---

## ISSUE 3

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

---

## ISSUE 4

# Title

Reliance on Magic Strings and Hardcoded IDs in Test Setup

##

internal/services/data_record/datasource_data_record_test.go

## Problem

Tests frequently use repeated, hardcoded UUIDs and other literal values (e.g., `"00000000-0000-0000-0000-000000000001"`) for environment IDs, entity IDs, and resource names. These "magic strings" recur many times, both in test resource configs and mock HTTP endpoints, increasing the chance of copy-paste errors and making it harder to update or refactor test infrastructure.

## Impact

- Difficult refactoring due to multiple places to update values (low-to-medium impact).
- Higher risk of inconsistencies and accidental mismatches.
- Obfuscates meaning and intent, reducing test readability and maintainability.

## Location

Examples:
- In HCL configuration strings (e.g., `environment_id = "00000000-0000-0000-0000-000000000001"`)
- In mock URLs (e.g., `"https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts?...`)
- In JSON paths, responders, etc.

## Code Issue

```go
// test HCL snippet
environment_id = "00000000-0000-0000-0000-000000000001"
// endpoint registration
httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts?$select=...", ...)
// repeated throughout the file
```

## Fix

Replace repeated magic strings with clearly-named constants at the top level of the test file, or with a function to generate standard/randomized test UUIDs. For example:

```go
const testEnvID = "00000000-0000-0000-0000-000000000001"
// Then use testEnvID everywhere instead of repeated literals
```
Or:
```go
func testEnvID() string { return \"00000000-0000-0000-0000-000000000001\" }
// Or return a randomized value if possible
```
This makes maintaining/generating test configuration and updating test resource IDs easier, and improves code readability and robustness.

Save as a structure/readability/maintainability issue.

---

## ISSUE 5

# Magic Strings for File Paths in Mocks

##

/workspaces/terraform-provider-power-platform/internal/services/application/datasource_environment_application_packages_test.go

## Problem

File paths to test fixtures are hardcoded as raw strings repeatedly, risking typos and making refactoring harder. Multiple places use the string "tests/datasource/environment_application_packages/..." directly, which is also duplicated across test functions.

## Impact

Severity: Low

Centralizing and reusing constants makes renaming and refactoring easier (single source of truth), reduces risk of typos, and improves readability.

## Location

Example:

```go
httpmock.File("tests/datasource/environment_application_packages/Validate_Read/get_applications.json")
```

## Fix

Introduce a constant or helper at the top of the test file:

```go
const testDataDir = "tests/datasource/environment_application_packages"

// then use:
httpmock.File(fmt.Sprintf("%s/Validate_Read/get_applications.json", testDataDir))
```

---

## ISSUE 6

# Use of Hardcoded Test Data in Test Configuration

##

/workspaces/terraform-provider-power-platform/internal/services/application/datasource_environment_application_packages_test.go

## Problem

The test cases use hardcoded values for resources and data sources, such as specific GUIDs and display names, which reduces the flexibility and maintainability of the test cases. This can make updates more difficult if the underlying data changes. It also makes parallelization and test isolation harder.

## Impact

Severity: Low

Using static values for identifiers and resource names can interfere with other tests running in parallel, may mask real failures, and increases maintenance effort as changes are needed across the tests if the specifics change.

## Location

Lines throughout test cases, e.g.:

```go
display_name      = "` + mocks.TestName() + `"
location          = "unitedstates"
environment_type  = "Sandbox"
// ...
security_group_id = "00000000-0000-0000-0000-000000000000"
```

## Fix

Generate dynamic test data where possible by using helper functions or introducing variables. For example, you may create random identifiers or use `mocks.TestName()` consistently for all identifiers.

```go
display_name      = "` + mocks.TestName() + `"
location          = "${var.test_location}"
environment_type  = "Sandbox"
dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "${var.test_security_group_id}"
}
```

Define and use variables or helper methods to create test-specific, unique values.

---

## ISSUE 7

# Title

Potential Flakiness due to Hardcoded File Paths in Test Mocks

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/datasource_environment_settings_test.go

## Problem

Test mocks use hardcoded file paths to test data (such as `"tests/datasource/Validate_Read/get_environment_00000000-0000-0000-0000-000000000001.json"`). This can cause test flakiness if tests are run from a different working directory, or if test data is moved/renamed without refactoring tests. There is no check/validation for file existence, so a missing file would cause a panic or unexpected test failure.

## Impact

Tests may fail unexpectedly if run in different environments, CI setups, or if the test data is restructured. This is a medium-severity issue as it can lead to unreliable test pipelines or local developer confusion.

## Location

Within the registered responders, for example:

```go
httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
    })
```

## Code Issue

```go
httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_environment_00000000-0000-0000-0000-000000000001.json").String())
```

## Fix

- Use a helper function to resolve paths relative to the test file (not working directory).
- Add explicit checks/asserts for file existence and clearer error reporting if the file is missing.

Example of path resolution and error handling:

```go
import (
    "os"
    "path/filepath"
    "runtime"
)

func testFileName(name string) string {
    _, filename, _, _ := runtime.Caller(0)
    return filepath.Join(filepath.Dir(filename), name)
}

// Usage in tests:
data, err := os.ReadFile(testFileName("tests/datasource/Validate_Read/get_environment_00000000-0000-0000-0000-000000000001.json"))
if err != nil {
    t.Fatalf("unable to read test fixture: %v", err)
}
return httpmock.NewStringResponse(http.StatusOK, string(data)), nil
```

---

This file will be saved as:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/testing/datasource_environment_settings_test.go_hardcoded_file_paths_medium.md`

---

## ISSUE 8

# Title

Hardcoded category names in `appendToList` calls

##

/workspaces/terraform-provider-power-platform/internal/services/environment_templates/datasource_environment_templates.go

## Problem

The `Read` method contains multiple hardcoded calls to `appendToList` with both the property name and the `"category"` string repeated, e.g.:

```go
appendToList(environment_templates.Standard, "standard", &state.Templates)
```

This pattern is repeated for `"premium"`, `"developer"`, etc. If the struct changes or new categories are added/removed, this block must be manually kept in sync, leading to error-prone and non-scalable code.

## Impact

Severity: Medium

- Maintainability issue—difficult to extend or refactor;
- Possible bugs if struct fields and hardcoded categories diverge;
- Reduces readability for future contributors.

## Location

`Read` method, lines with:

```go
appendToList(environment_templates.Standard, "standard", &state.Templates)
appendToList(environment_templates.Premium, "premium", &state.Templates)
...
```

## Code Issue

```go
appendToList(environment_templates.Standard, "standard", &state.Templates)
appendToList(environment_templates.Premium, "premium", &state.Templates)
...
```

## Fix

Use a static slice or map to express the property–category relationship once, then loop, e.g.:

```go
categories := []struct{
    items []itemDto
    name string
}{
    {environment_templates.Standard, "standard"},
    {environment_templates.Premium, "premium"},
    {environment_templates.Developer, "developer"},
    {environment_templates.Basic, "basic"},
    {environment_templates.Production, "production"},
    {environment_templates.Sandbox, "sandbox"},
    {environment_templates.Trial, "trial"},
    {environment_templates.Default, "default"},
    {environment_templates.Support, "support"},
    {environment_templates.SubscriptionBasedTrial, "subscriptionBasedTrial"},
    {environment_templates.Teams, "teams"},
    {environment_templates.Platform, "platform"},
}

for _, c := range categories {
    appendToList(c.items, c.name, &state.Templates)
}
```

Or use reflection if extensibility is needed.


---

## ISSUE 9

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

---


# To finish the task you have to 
1. Run linter and fix any issues 
2. Run UnitTest and fix any of failing ones
3. Generate docs 
4. Run Changie

# Changie Instructions
Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```
Where:
- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for "copilot-commit-message-instructions.md" how to write description.
- `<issue_number>` pick the issue number or PR number
```
