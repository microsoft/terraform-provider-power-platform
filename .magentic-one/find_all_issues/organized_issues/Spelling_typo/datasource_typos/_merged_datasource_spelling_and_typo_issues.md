# Data Source Spelling and Typo Issues

This document contains all spelling and typo issues found in data source-related files within the terraform-provider-power-platform codebase.

## ISSUE 1

### Typo in `ApplicaitonId` and Typo in Key

**File:** `/workspaces/terraform-provider-power-platform/internal/services/application/datasource_environment_application_packages.go`

**Problem:** There is a typo in the schema attribute's markdown description for `application_id` field: "ApplicaitonId" instead of "ApplicationId".

**Impact:** Severity: **Low**

This typo does not affect code execution but may confuse users reading the generated provider documentation and lower the perceived API quality.

**Location:**

```go
"application_id": schema.StringAttribute{
    MarkdownDescription: "ApplicaitonId",
    Computed:            true,
},
```

**Code Issue:**

```go
"application_id": schema.StringAttribute{
    MarkdownDescription: "ApplicaitonId",
    Computed:            true,
},
```

**Fix:** Fix the typo in the markdown:

```go
"application_id": schema.StringAttribute{
    MarkdownDescription: "ApplicationId",
    Computed:            true,
},
```

## ISSUE 2

### Inconsistent Naming: `PowerAppssClient` should be `PowerAppsClient`

**File:** `/workspaces/terraform-provider-power-platform/internal/services/powerapps/datasource_environment_powerapps.go`

**Problem:** The field is named `PowerAppssClient` (with a double "s") instead of the expected `PowerAppsClient`. This inconsistency can cause confusion for developers and increases the risk of typos in other parts of the codebase.

**Impact:** Severity: Low  
This is a minor naming issue. However, inconsistent naming can reduce code readability and maintainability, and may cause subtle bugs if developers mistakenly use the wrong identifier.

**Location:**

```go
d.PowerAppssClient = newPowerAppssClient(client.Api)
apps, err := d.PowerAppssClient.GetPowerApps(ctx)
```

**Code Issue:**

```go
d.PowerAppssClient = newPowerAppssClient(client.Api)
apps, err := d.PowerAppssClient.GetPowerApps(ctx)
```

**Fix:** Change all instances of `PowerAppssClient` to `PowerAppsClient`. Also, ensure that the struct declares this property with the correct name.

```go
// Rename struct field and all its references
d.PowerAppsClient = newPowerAppsClient(client.Api)
apps, err := d.PowerAppsClient.GetPowerApps(ctx)
```

## ISSUE 3

### Confusing or Inaccurate Variable Naming: dvExits

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment_settings/datasource_environment_settings.go`

**Problem:** The variable `dvExits` is likely a typo—it should probably be `dvExists`. This affects readability.

**Impact:** Low. Naming clarity impacts maintainability, but the logic still works as expected.

**Location:** All of:

```go
dvExits, err := d.EnvironmentSettingsClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
if err != nil {
 resp.Diagnostics.AddError(...)
}
if !dvExits {
 resp.Diagnostics.AddError(...)
 return
}
```

**Code Issue:**

```go
dvExits, err := d.EnvironmentSettingsClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
if err != nil {
 resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
}
if !dvExits {
 resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), "")
 return
}
```

**Fix:** Rename variable `dvExits` to `dvExists`.

```go
dvExists, err := d.EnvironmentSettingsClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
if err != nil {
 resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
 return
}
if !dvExists {
 resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), "")
 return
}
```

## ISSUE 4

### Incorrect Spelling in Field Names and Documentation

**File:** `/workspaces/terraform-provider-power-platform/internal/services/application/datasource_tenant_application_packages.go`

**Problem:** Throughout the code and schema definition, the word "Description" is misspelled as "Descprition" and "Applicaiton" instead of "Application". This is present both in struct fields and their documentation/markdown. Having incorrect spelling decreases code readability and causes confusion, especially in APIs, data models, and schema attributes that need to be referenced elsewhere or mapped to upstream/downstream systems.

**Impact:**

- Low to Medium:
  - Source code and schema attribute misspellings cause user confusion and difficulty for maintainers.
  - API consumers or integrators might reference incorrect or inconsistent field names, leading to bugs or undocumented behavior.
  - Inconsistent spelling increases chances of runtime errors if fields are referenced dynamically or via reflection.

**Location:**

- Schema definition under attributes for "application_descprition" and markdown.
- Data model mapping in the Read method and related structs.

**Code Issue:**

```go
"application_id": schema.StringAttribute{
 MarkdownDescription: "ApplicaitonId",
 Computed:            true,
},
"application_descprition": schema.StringAttribute{
 MarkdownDescription: "Applicaiton Description",
 Computed:            true,
},
...
ApplicationDescprition: types.StringValue(application.ApplicationDescription),
```

**Fix:** Update attribute names, markdown, and struct fields to correct spelling:

```go
"application_id": schema.StringAttribute{
 MarkdownDescription: "Application ID",
 Computed:            true,
},
"application_description": schema.StringAttribute{
 MarkdownDescription: "Application Description",
 Computed:            true,
},
...
ApplicationDescription: types.StringValue(application.ApplicationDescription),
```

- Update the struct fields/properties to use `ApplicationDescription` consistently.
- Change attribute map keys and MarkdownDescription values to avoid inconsistent spelling.
- Confirm these corrections in all model, mapping, and schema places.

## ISSUE 5

### Incorrect function name in constructor

**File:** `/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity.go`

**Problem:** The constructor function for the data source is named `NewTenantCapcityDataSource`, which appears to be a typo – it should likely be `NewTenantCapacityDataSource` (missing the "a" in "Capacity").

**Impact:** Could cause confusion and reduce codebase maintainability. This is a high-severity naming issue because consumers of the function may not expect the typo and it could lead to errors or misunderstandings.

**Location:** Line 19

**Code Issue:**

```go
func NewTenantCapcityDataSource() datasource.DataSource {
    return &DataSource{
        TypeInfo: helpers.TypeInfo{
            TypeName: "tenant_capacity",
        },
    }
}
```

**Fix:** Rename the function to correct the typo:

```go
func NewTenantCapacityDataSource() datasource.DataSource {
    return &DataSource{
        TypeInfo: helpers.TypeInfo{
            TypeName: "tenant_capacity",
        },
    }
}
```

This helps ensure clarity and correct usage in the codebase.

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
