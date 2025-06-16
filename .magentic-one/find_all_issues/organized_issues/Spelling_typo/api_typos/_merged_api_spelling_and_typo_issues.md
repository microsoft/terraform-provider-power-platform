# API Spelling and Typo Issues

This document contains all spelling and typo issues found in API-related files within the terraform-provider-power-platform codebase.

## ISSUE 1

### Incorrect Spelling in Log Message

**File:** `/workspaces/terraform-provider-power-platform/internal/services/application/api_application.go`

**Problem:** There is a typo in the debug logging statement `"Opeartion Location Header"`; it should be `"Operation Location Header"`.

**Impact:** This issue has a **low** severity. While it doesn't impact program logic or functionality directly, spelling mistakes in log messages can make logs harder to search and lead to confusion or mistakes during debugging and troubleshooting.

**Location:** Line inside `InstallApplicationInEnvironment`, in the following code:

**Code Issue:**

```go
tflog.Debug(ctx, "Opeartion Location Header: "+operationLocationHeader)
```

**Fix:** Correct the spelling of "Opeartion" to "Operation":

```go
tflog.Debug(ctx, "Operation Location Header: "+operationLocationHeader)
```

## ISSUE 2

### Misspelled Variable Name: "virutualConnector"

**File:** `/workspaces/terraform-provider-power-platform/internal/services/connectors/api_connectors.go`

**Problem:** The variable `virutualConnector` in this section is misspelled; it should be `virtualConnector` for consistency and clarity.

**Impact:** Misspelled variable names decrease code readability, can cause confusion during maintenance or code reviews, and may introduce bugs if similarly-named variables are later introduced. Severity: **low**.

**Location:** Within the `GetConnectors` method, this line:

**Code Issue:**

```go
for _, virutualConnector := range virtualConnectorArray {
 connectorArray.Value = append(connectorArray.Value, connectorDto{
  Id:   virutualConnector.Id,
  Name: virutualConnector.Metadata.Name,
  Type: virutualConnector.Metadata.Type,
  Properties: connectorPropertiesDto{
   DisplayName: virutualConnector.Metadata.DisplayName,
   Unblockable: false,
   Tier:        "Built-in",
   Publisher:   "Microsoft",
   Description: "",
  },
 })
}
```

**Fix:** Rename the variable for correct spelling:

```go
for _, virtualConnector := range virtualConnectorArray {
 connectorArray.Value = append(connectorArray.Value, connectorDto{
  Id:   virtualConnector.Id,
  Name: virtualConnector.Metadata.Name,
  Type: virtualConnector.Metadata.Type,
  Properties: connectorPropertiesDto{
   DisplayName: virtualConnector.Metadata.DisplayName,
   Unblockable: false,
   Tier:        "Built-in",
   Publisher:   "Microsoft",
   Description: "",
  },
 })
}
```

## ISSUE 3

### Typo in Struct and Variable Naming

**File:** `/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/api_enterprise_policy.go`

**Problem:** The struct and variable names `linkEnterprosePolicyDto` are incorrectly spelled. The correct spelling should be `linkEnterprisePolicyDto`. This typo appears at variable declarations and literal assignments in both `LinkEnterprisePolicy` and `UnLinkEnterprisePolicy` functions.

**Impact:** This impacts both code readability and maintainability. Developers unfamiliar with the code might be confused, and searches for `enterprise` will not find these entries. Severity: **low** (does not directly affect functionality but affects clarity and can be error-prone in future maintenance).

**Location:** Lines where the DTO struct and variables are named.

**Code Issue:**

```go
linkEnterprosePolicyDto := linkEnterprosePolicyDto{
 SystemId: systemId,
}
```

**Fix:** Rename all instances of `linkEnterprosePolicyDto` to `linkEnterprisePolicyDto` (struct and variable names). Confirm the type definition is also correctly named elsewhere.

```go
linkEnterprisePolicyDto := linkEnterprisePolicyDto{
 SystemId: systemId,
}
```

## ISSUE 4

### Misspelled struct name: `enironmentDeleteDto`

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment/api_environment.go`

**Problem:** There is a misspelled struct name `enironmentDeleteDto` instead of the likely intended `environmentDeleteDto`. This could be confusing to maintainers and can easily lead to inconsistencies or further typos in referencing this type throughout the codebase. If this typo also occurs in the definition, it could have broader effects on code readability elsewhere.

**Impact:**

- Severity: Low
- This is primarily a readability and code quality issue. It does not directly cause program failures but creates confusion for future maintainers or contributors.

**Location:** Line in function `DeleteEnvironment`:

**Code Issue:**

```go
environmentDelete := enironmentDeleteDto{
 Code:    "7", // Application.
 Message: "Deleted using Power Platform Terraform Provider",
}
```

**Fix:** Change all instances of `enironmentDeleteDto` to `environmentDeleteDto` to improve naming clarity.

```go
environmentDelete := environmentDeleteDto{
 Code:    "7", // Application.
 Message: "Deleted using Power Platform Terraform Provider",
}
```

And make sure the struct is defined with the correct name as well, if it exists in this package or imported.

## ISSUE 5

### Misspelled Function Name for Removing Environments

**File:** `/workspaces/terraform-provider-power-platform/internal/services/licensing/api_licensing.go`

**Problem:** There is a typo in the function name: `RemoveEnvironmentsToBillingPolicy`. This should instead be `RemoveEnvironmentsFromBillingPolicy` to correctly reflect the operation semantics.

**Impact:** Low. This is a semantic and clarity problem, but may cause confusion for developers and users of the API.

**Location:**

```go
func (client *Client) RemoveEnvironmentsToBillingPolicy(ctx context.Context, billingId string, environmentIds []string) error {
```

**Fix:** Correct the method name and any references:

```go
func (client *Client) RemoveEnvironmentsFromBillingPolicy(ctx context.Context, billingId string, environmentIds []string) error {
```

## ISSUE 6

### String-based Comparison for Role Names Instead of Constants or Enums

**File:** `/workspaces/terraform-provider-power-platform/internal/services/authorization/api_user.go`

**Problem:** In the function `GetEnvironmentUserByAadObjectId`, the code determines role assignment by comparing role names using string equality with untyped literals ("EnvironmentAdmin", "EnvironmentMaker"). If the backend changes these string literals or a typo occurs, it would result in subtle runtime errors. This approach couples logic to string values and reduces type safety.

**Impact:** Severity: Low

This is mostly a maintainability and reliability issue. While it is unlikely to cause immediate failure, changes in API contract (role name typo, casing change, etc.), or misspellings, will break functionality without compiler errors and could be hard to detect.

**Location:** Within GetEnvironmentUserByAadObjectId:

**Code Issue:**

```go
isAdminRole := roleAssignment.Properties.RoleDefinition.Name == "EnvironmentAdmin"
isMakerRole := roleAssignment.Properties.RoleDefinition.Name == "EnvironmentMaker"
```

**Fix:** Define role constants or an enum-like structure for these commonly used string values, and reference those instead. Example:

```go
const (
    RoleEnvironmentAdmin = "EnvironmentAdmin"
    RoleEnvironmentMaker = "EnvironmentMaker"
)

...

isAdminRole := roleAssignment.Properties.RoleDefinition.Name == RoleEnvironmentAdmin
isMakerRole := roleAssignment.Properties.RoleDefinition.Name == RoleEnvironmentMaker
```

This enhances reliability, facilitates searchability, and limits bugs from string typos.

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
