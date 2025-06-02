# Title

Misnamed Field: `application_descprition` should be `application_description`

## File Path

`/workspaces/terraform-provider-power-platform/internal/services/application/datasource_tenant_application_packages_test.go`

## Problem

There is a typo in the field name `application_descprition`, which should be `application_description`. This seems to be incorrectly used in multiple test cases.

## Impact

The issue could cause unexpected behavior or errors in validation, especially during tests related to application metadata fields. While this might not affect the main code functionality directly, it reduces code readability and correctness in tests.

**Severity: Medium**

## Location

Occurrences of incorrect usage of the field `application_descprition`.

## Code Issue

```go
resource.TestMatchResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.application_descprition", regexp.MustCompile(helpers.StringRegex))
resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.application_descprition", "An easier way to get and manage approvals.")
```

## Fix

The field name should be corrected to `application_description`. Here is the corrected code:

```go
resource.TestMatchResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.application_description", regexp.MustCompile(helpers.StringRegex))
resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.application_description", "An easier way to get and manage approvals.")
```