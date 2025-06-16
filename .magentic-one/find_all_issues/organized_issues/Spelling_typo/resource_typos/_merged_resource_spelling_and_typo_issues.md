# Resource Spelling and Typo Issues

This document contains all spelling and typo issues found in resource files within the terraform-provider-power-platform codebase.

## ISSUE 1

### Duplicated schema documentation for "connection_parameters" and "connection_parameters_set"

**File:** `/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection.go`

**Problem:** Both the `connection_parameters` and `connection_parameters_set` attributes in the schema have essentially identical documentation strings, including the same note with a typo ("requried in-place-update"). This introduces ambiguity about their distinct purposes, hinders readability, and may confuse users and maintainers.

**Impact:** Low: Does not affect correctness, but reduces overall documentation clarity for both code maintainers and resource users.

**Location:** Resource Schema definition:

**Code Issue:**

```go
"connection_parameters": schema.StringAttribute{
    MarkdownDescription: "Connection parameters. Json string containing the authentication connection parameters ...",
    ...
},
"connection_parameters_set": schema.StringAttribute{
    MarkdownDescription: "Set of connection parameters. Json string containing the authentication connection parameters ...",
    ...
},
```

**Fix:** Review and clarify the documentation for each field so that their specific roles and any differences are apparent, removing typos and duplicated information where not relevant.

```go
"connection_parameters": schema.StringAttribute{
    MarkdownDescription: "Connection parameters. JSON string with authentication details, used when ...",
    ...
},
"connection_parameters_set": schema.StringAttribute{
    MarkdownDescription: "(Advanced) An explicit set of connection parameters (JSON string), for ... [explain use case and distinction].",
    ...
},
```

Also, fix the typo: "requried in-place-update" → "required in-place update".

## ISSUE 2

### Incorrect variable naming: `conectionState` typo (should be `connectionState`)

**File:** `/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection.go`

**Problem:** There is a typo in the variable name `conectionState` (should be `connectionState`) throughout the Create, Read, and Update methods.

**Impact:** Medium severity: Typos in variable names can lead to confusion and reduce code readability and maintainability. This doesn't break functionality but is a code quality concern.

**Location:** Multiple occurrences in methods: Create, Read, Update:

**Code Issue:**

```go
conectionState := ConvertFromConnectionDto(*connection)
plan.Id = types.String(conectionState.Id)
// ... and similar lines in Read and Update
```

**Fix:** Rename all occurrences of `conectionState` to `connectionState` for clarity and consistency.

```go
connectionState := ConvertFromConnectionDto(*connection)
plan.Id = types.String(connectionState.Id)
// ... and similar fixes in Read and Update
```

## ISSUE 3

### Typo in Attribute Descriptions: Markdown and Enforcement

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/resource_environment_group_rule_set.go`

**Problem:** There are repeated typographical errors in Markdown descriptions for schema attributes (such as `MarkdownDescription: "Solution checker enforceemnt mode: none, warm, block"` and others like `"Enable AI generated description"`). These should be corrected for professionalism and clarity.

**Impact:** Low.

- Minor: only documentation is affected, not code logic.
- However, typos can affect user trust and documentation usability.

**Location:** Lines such as:

```go
MarkdownDescription: "Solution checker enforceemnt mode: none, warm, block",
```

and

```go
MarkdownDescription: "Agree to enable Bing search features",
```

(and others containing "enbaled", "Inculde", etc.)

**Fix:** Correct the MarkdownDescriptions:

```go
MarkdownDescription: "Solution checker enforcement mode: none, warn, block",
MarkdownDescription: "Enable AI generated description",
MarkdownDescription: "Include insights for all Managed Environments in this group in weekly email digest.",
MarkdownDescription: "Agree to enable Bing search features",
```

## ISSUE 4

### Variable Naming: `dvExits` should be `dvExists`

**File:** `/workspaces/terraform-provider-power-platform/internal/services/solution/resource_solution.go`

**Problem:** A variable is named `dvExits`, but the intended word is `Exists` (as in "does it exist?"). This appears to be a typographical error and could lead to misunderstanding or reduced readability.

**Impact:**

- **Severity:** Low
- Minor impact on readability, but could cause confusion during code reviews or maintenance, especially for non-native English speakers or new contributors.

**Location:** `importSolution` function, line using:

**Code Issue:**

```go
dvExits, err := r.SolutionClient.DataverseExists(ctx, plan.EnvironmentId.ValueString())
if err != nil {
    diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), err.Error())
}

if !dvExits {
    diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), "")
    return nil
}
```

**Fix:** Rename the variable to `dvExists` for clarity:

```go
dvExists, err := r.SolutionClient.DataverseExists(ctx, plan.EnvironmentId.ValueString())
if err != nil {
    diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), err.Error())
}

if !dvExists {
    diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), "")
    return nil
}
```

## ISSUE 5

### Minor code readability and docstring/description typos in schema block

**File:** `/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user.go`

**Problem:** Several docstrings and Markdown descriptions in the `Schema` method (especially for `disable_delete` and `security_roles` attributes) contain typos, markdown errors, or unclear phrasing. Example issues include:  

- Spelling: "Delte" instead of "Delete", "propertyto" instead of "property to"  
- Markdown syntax error: `(Disable Delte)[URL]` should be `[Disable Delete](URL)`  
- Unclear documentation logic on relationship of options.

**Impact:** These documentation typos and doc issues may confuse users, reduce end-user confidence, and undermine the usability of provider-generated docs in the Terraform Registry. Severity: **Low**.

**Location:** Within the `Schema` method, for multiple attribute docstrings, e.g.:

```go
"disable_delete": schema.BoolAttribute{
    MarkdownDescription: "Disable delete. When set to `True` is expects that (Disable Delte)[https://learn.microsoft.com/power-platform/admin/delete-users..." +
        "... If you just want to remove the resource and not delete the user from Dataverse, set this propertyto `False`\n\n" +
        ...
},
```

**Code Issue:**

```go
MarkdownDescription: "Disable delete. When set to `True` is expects that (Disable Delte)[https://learn.microsoft.com/power-platform/admin/delete-users?WT.mc_id=ppac_inproduct_settings#soft-delete-users-in-power-platform] feature to be enabled." +
    "Removing resource will try to delete the systemuser from Dataverse. This is the default behaviour. If you just want to remove the resource and not delete the user from Dataverse, set this propertyto `False`\n\n" +
    "**This attribute applies only when working with dataverse users.**",
```

**Fix:** Correct all spelling and markdown errors for clarity and better user-facing documentation:

```go
MarkdownDescription: "Disable delete. When set to `True`, it expects that [Disable Delete](https://learn.microsoft.com/power-platform/admin/delete-users?WT.mc_id=ppac_inproduct_settings#soft-delete-users-in-power-platform) feature to be enabled." +
    "Removing the resource will try to delete the system user from Dataverse (default behavior). If you just want to remove the resource and not delete the user from Dataverse, set this property to `False`.\n\n" +
    "**This attribute applies only when working with Dataverse users.**",
```

Perform similar corrections for other documentation fields as needed.

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
