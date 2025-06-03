# Typo in Struct Attribute: application_descprition

##

/workspaces/terraform-provider-power-platform/internal/services/application/datasource_tenant_application_packages_test.go

## Problem

There is a repeated typo in the attribute name `application_descprition` (should be `application_description`). The misspelled attribute appears both in resource checks in both acceptance and unit tests. This may lead to confusion, reduce maintainability, and cause integration issues with other parts of the codebase expecting the correctly spelled attribute.

## Impact

Severity: **Medium**.  
The typo reduces code readability, may propagate confusion across tests and production code, and could result in errors if the schema changes or is referenced elsewhere with the correct spelling.

## Location

- Several instances in the resource attribute checks, e.g., lines containing `applications.0.application_descprition`.

## Code Issue

```go
resource.TestMatchResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.application_descprition", regexp.MustCompile(helpers.StringRegex))
...
resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.application_descprition", "An easier way to get and manage approvals.")
```

## Fix

Update all resource attribute checks and references in tests to use the correct spelling `application_description`.

```go
resource.TestMatchResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.application_description", regexp.MustCompile(helpers.StringRegex))
...
resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.application_description", "An easier way to get and manage approvals.")
```
