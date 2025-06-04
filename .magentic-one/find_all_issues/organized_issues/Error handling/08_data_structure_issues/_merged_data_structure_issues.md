# Data Structure Issues

This document consolidates all issues related to data structure organization, schema definitions, and conversion logic in the Terraform Provider for Power Platform.

## ISSUE 1

### Structure and Maintainability: Large Method for Schema Definition

**File:** `/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/datasource_analytics_data_exports.go`

**Problem:** The `Schema` method contains a deeply nested and lengthy literal schema definition with 8+ levels of indentation and many attributes. This can be difficult to read, maintain, and extend, as it is easy for errors or copy-paste mistakes to occur. It is also hard to quickly grasp the schema structure, especially for more complex objects.

**Impact:**

- Maintenance is challenging; adding or removing fields is error prone.
- Code review and comprehension are harder for complex, deeply nested object schemas.
- Can introduce subtle inconsistencies in style, docs, or required/computed flags over time.
- **Severity:** Medium

**Location:** The method `func (d *AnalyticsExportDataSource) Schema(...)` and its schema literal.

**Code Issue:**

```go
func (d *AnalyticsExportDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
    defer exitContext()
    resp.Schema = schema.Schema{
        // ... huge nested literal ...
    }
}
```

**Fix:** Move repeated or deeply nested attributes into helper functions or variables:

```go
var sinkAttribute = schema.SingleNestedAttribute{
    MarkdownDescription: "The sink configuration for analytics data",
    Required: true,
    Attributes: map[string]schema.Attribute{
        // ...
    },
}
// In Schema method:
"sink": sinkAttribute,
```

This modularizes the schema definition, improves readability, enables reuse, and reduces risk of copy-paste or indentation errors.

## ISSUE 2

### Large Schema Function, Highly Nested Structure

**File:** `/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/datasource_tenant_settings.go`

**Problem:** The `Schema` method of `TenantSettingsDataSource` is extremely large and consists of deeply nested attribute declarations to describe all possible tenant settings. This results in a method that is very difficult to read, scan, debug, and extend. It increases the risk of unintentional errors (omitting commas, incorrect nesting, hard-to-locate settings) and makes consistency and future modifications daunting for maintainers.

**Impact:** Complicates maintenance and readability, increases cognitive overhead for new contributors, and heightens the risk of subtle bugs due to difficult-to-review code. **Severity: Medium**

**Location:** Main body of the `Schema` method, especially:

```go
func (d *TenantSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
   ...
        Attributes: map[string]schema.Attribute{
            "timeouts": timeouts.Attributes(ctx, timeouts.Opts{
                Create: false,
                Update: false,
                Delete: false,
                Read:   false,
            }),
            ...
            "power_platform": schema.SingleNestedAttribute{
                ...
                Attributes: map[string]schema.Attribute{
                    ... // deep nesting for dozens of settings
                },
            },
        },
   ...
}
```

**Fix:** Factor the deepest/nested groups of attributes into their own helper functions that return `map[string]schema.Attribute` values:

```go
func powerPlatformAttributes() map[string]schema.Attribute {
    return map[string]schema.Attribute{
        // place Power Platform nested attributes here...
    }
}
// Then in Schema:
"power_platform": schema.SingleNestedAttribute{
    MarkdownDescription: "Power Platform",
    Computed:            true,
    Attributes:          powerPlatformAttributes(),
},
```

This reduces function size and nesting, and helps enforce separation of concerns. Repeat for other major top-level groups ("governance", "power_apps", etc.) where applicable.

## ISSUE 3

### Duplication of Conversion Logic Across Similar Functions

**File:** `/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/dto.go`

**Problem:** Many conversion functions are highly similar, each translating a specific DTO segment, but repeating the same patterns of attribute type creation, null checks, and attribute value assignment. This duplication leads to a proliferation of boilerplate, increases the maintenance burden, and makes changes prone to error.

**Impact:**

- Harder to maintain: any fixes or enhancements must be applied in multiple places.
- Higher risk of inconsistency and subtle bugs, as some functions may drift apart over time.
- Difficult to refactor at a systemic level because of scattered duplications.
- Reduces clarity, especially for newcomers or reviewers. Severity: medium.

**Location:** Functions like: `convertUserManagementSettings`, `convertCatalogSettings`, `convertPowerAppsSettings`, `convertLicensingSettings`, etc.

**Code Issue:**

```go
func convertUserManagementSettings(tenantSettingsDto tenantSettingsDto) (basetypes.ObjectType, basetypes.ObjectValue) {
    attrTypesUserManagementSettings := map[string]attr.Type{
        "enable_delete_disabled_user_in_all_environments": types.BoolType,
    }
    if tenantSettingsDto.PowerPlatform == nil || tenantSettingsDto.PowerPlatform.UserManagementSettings == nil {
        return types.ObjectType{AttrTypes: attrTypesUserManagementSettings}, types.ObjectNull(attrTypesUserManagementSettings)
    }
    attrValuesUserManagementSettings := map[string]attr.Value{
        "enable_delete_disabled_user_in_all_environments": types.BoolPointerValue(tenantSettingsDto.PowerPlatform.UserManagementSettings.EnableDeleteDisabledUserinAllEnvironments),
    }
    return types.ObjectType{AttrTypes: attrTypesUserManagementSettings}, types.ObjectValueMust(attrTypesUserManagementSettings, attrValuesUserManagementSettings)
}
```

**Fix:** Create a set of utility functions or higher-order abstractions that can factor out common conversion patterns:

```go
func convertSingleBoolField(dto interface{}, attrName string, fieldPtr *bool) (basetypes.ObjectType, basetypes.ObjectValue) {
    attrTypes := map[string]attr.Type{attrName: types.BoolType}
    if fieldPtr == nil {
        return types.ObjectType{AttrTypes: attrTypes}, types.ObjectNull(attrTypes)
    }
    attrValues := map[string]attr.Value{attrName: types.BoolPointerValue(fieldPtr)}
    return types.ObjectType{AttrTypes: attrTypes}, types.ObjectValueMust(attrTypes, attrValues)
}
```

Apply this principle to remove boilerplate and improve code quality throughout the conversion functions.

## ISSUE 4

### Incomplete Mapping of Feature State Strings

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment_wave/resource_environment_wave.go`

**Problem:** The function `mapFeatureStateToSchemaState` provides an incomplete mapping of potential API states to the possible schema states. The function maps only two specific string values from the API (`Upgrading` to `upgrading`, `ON` to `enabled`), and any other value is mapped to `error`, including potentially valid but unknown or new values. This will result in all new/unknown states being treated as errors.

**Impact:** If Microsoft adds states to the API, the provider will interpret all of them as `error`, making troubleshooting or feature evolution difficult. Severity: **medium** because it risks misrepresentation of provider state and can cause downstream automation or monitoring to misbehave.

**Location:**

```go
func mapFeatureStateToSchemaState(apiState string) string {
        switch apiState {
        case "Upgrading":
                return "upgrading"
        case "ON":
                return "enabled"
        default:
                return "error"
        }
}
```

**Code Issue:**

```go
func mapFeatureStateToSchemaState(apiState string) string {
        switch apiState {
        case "Upgrading":
                return "upgrading"
        case "ON":
                return "enabled"
        default:
                return "error"
        }
}
```

**Fix:** Implement logging for unmapped/unknown values and/or a passthrough or explicit error state with diagnostics:

```go
func mapFeatureStateToSchemaState(apiState string) string {
        switch apiState {
        case "Upgrading":
                return "upgrading"
        case "ON":
                return "enabled"
        case "Error", "ERROR":
                return "error"
        default:
                tflog.Warn(context.TODO(), fmt.Sprintf("Unknown feature state from API: %s", apiState))
                return "error"
        }
}
```

- Or consider returning an explicit error or separate status for unmapped states
- Add documentation for maintainers that API state values might change

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
