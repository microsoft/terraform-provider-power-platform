# Models Spelling and Typo Issues

This document contains all spelling and typo issues found in model files within the terraform-provider-power-platform codebase.

## ISSUE 1

### Misspelled Field Name in Struct

**File:** `/workspaces/terraform-provider-power-platform/internal/services/application/models.go`

**Problem:** The struct field `ApplicationDescprition` in `TenantApplicationPackageDataSourceModel` is misspelled. It should be `ApplicationDescription`.

**Impact:** Having a typo in the struct field name can cause confusion, reduce code readability, and may also lead to inconsistent behavior when the code expects the correct field name. This can introduce bugs when interacting with APIs, terraform schema, or refactoring code in the future.  
Severity: Low

**Location:** Line defining `ApplicationDescprition` in the struct `TenantApplicationPackageDataSourceModel`.

**Code Issue:**

```go
ApplicationDescprition types.String                                   `tfsdk:"application_descprition"`
```

**Fix:** Correct the spelling in both the struct field name and the struct tag:

```go
ApplicationDescription types.String                                   `tfsdk:"application_description"`
```

## ISSUE 2

### Misspelling: `Catergory` field name in `ClusterDto`

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment/models.go`

**Problem:** The field `Catergory` in `ClusterDto` is misspelled; it should be `Category`. This typo appears in multiple places (both reading and assignment). Using misspelled identifiers decreases code clarity and increases the risk that future contributors may misinterpret or overlook the field. Furthermore, it reduces consistency with other APIs or DTOs that use the correct spelling.

**Impact:**

- **Severity:** Medium
- Reduces code readability.
- May introduce subtle bugs if misreferenced elsewhere or during API integration.
- Decreases maintainability and professionalism.

**Location:**

```go
if environmentSource.ReleaseCycle.ValueString() == ReleaseCycleTypesEarly {
 value := conf.GetCurrentCloudConfiguration(config.FirstReleaseClusterName)
 if value != nil {
  environmentDto.Properties.Cluster = &ClusterDto{
   Catergory: *value,
  }
 }
}
```

And

```go
func convertReleaseCycleModelFromDto(environmentDto EnvironmentDto, model *SourceModel, providerConfig config.ProviderConfig) {
 value := providerConfig.GetCurrentCloudConfiguration(config.FirstReleaseClusterName)
 if environmentDto.Properties.Cluster != nil && value != nil && environmentDto.Properties.Cluster.Catergory == *value {
  model.ReleaseCycle = types.StringValue(ReleaseCycleTypesEarly)
 } else {
  model.ReleaseCycle = types.StringValue(ReleaseCycleTypesStandard)
 }
}
```

**Code Issue:**

```go
environmentDto.Properties.Cluster = &ClusterDto{
 Catergory: *value,
}
...
if environmentDto.Properties.Cluster != nil && value != nil && environmentDto.Properties.Cluster.Catergory == *value {
```

**Fix:** Correct all usages from `Catergory` to `Category`:

```go
environmentDto.Properties.Cluster = &ClusterDto{
 Category: *value,
}
...
if environmentDto.Properties.Cluster != nil && value != nil && environmentDto.Properties.Cluster.Category == *value {
```

Be sure to rename the field definition in the `ClusterDto` type itself and update all affected usages throughout the codebase for consistency.

## ISSUE 3

### Incorrect Interface Name for the Environment Settings Client Field

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment_settings/models.go`

**Problem:** In the `EnvironmentSettingsDataSource` and `EnvironmentSettingsResource` struct definitions, the struct tags use different interface/struct names for the client field: `EnvironmentSettingsClient` and `EnvironmentSettingClient` (missing 's'), which likely leads to a typo and logical error if the interfaces are not both defined and consistent.

**Impact:** This inconsistency will cause compile-time errors if `client` variable or interface is not defined with both names, and can lead to confusion or bugs if the wrong client is injected. This is a medium-severity issue due to reliability and maintenance impacts.

**Location:** Lines:

```go
type EnvironmentSettingsDataSource struct {
 helpers.TypeInfo
 EnvironmentSettingsClient client
}

type EnvironmentSettingsResource struct {
 helpers.TypeInfo
 EnvironmentSettingClient client
}
```

**Code Issue:**

```go
type EnvironmentSettingsDataSource struct {
 helpers.TypeInfo
 EnvironmentSettingsClient client
}

type EnvironmentSettingsResource struct {
 helpers.TypeInfo
 EnvironmentSettingClient client
}
```

**Fix:** Ensure that both structs use the same, correct interface/type for the client field. For example, if the type should be `EnvironmentSettingsClient`, update both definitions as such:

```go
type EnvironmentSettingsDataSource struct {
 helpers.TypeInfo
 EnvironmentSettingsClient client
}

type EnvironmentSettingsResource struct {
 helpers.TypeInfo
 EnvironmentSettingsClient client
}
```

This change makes the naming consistent and less error-prone. If the missing 's' is intentional, clarify it in documentation or comments for future maintainers.

---

Apply this fix to the whole codebase

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
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
