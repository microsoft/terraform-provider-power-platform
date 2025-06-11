# Models Naming Issues - Merged Issues

## ISSUE 1

# Title

Mixed Naming Conventions for Field Names

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/models.go

## Problem

Go recommends using consistent naming conventions, typically CamelCase for struct fields and acronyms. In this file, `SubscriptionId` does not match Go's convention, where it should be `SubscriptionID`. Other fields such as `ID`, `AiType`, or `PackageName` are correct (ID/AI both in all-caps as is canonical for Go), but `SubscriptionId` uses a lowercase `d`.

## Impact

Severity: **Low**

This is mostly a readability and maintainability issue, but inconsistent naming can be confusing to contributors and breaks Go idioms, which may affect tooling and code generation further down the line.

## Location

Across all struct definitions in the file:

```go
	SubscriptionId    types.String `tfsdk:"subscription_id"`
```

## Fix

Rename the field to `SubscriptionID`. Ensure all code referencing this field is updated accordingly.

```go
	SubscriptionID    types.String `tfsdk:"subscription_id"`
```

Apply for the whole codebase


---

## ISSUE 2

# Inconsistent Struct Field Names: "Id" vs "ID"

application/models.go

## Problem

In `EnvironmentApplicationPackageInstallResourceModel`, the field is named `Id` instead of `ID`, which is inconsistent with Go naming conventions (acronyms should be all uppercase, e.g., `ID`). Other similar fields like `EnvironmentId`, `ApplicationId`, and `PublisherId` also do not follow the Go convention of all-uppercase "ID".

## Impact

This breaks Go idiomatic naming conventions, makes the code less readable and maintainable, and can cause subtle issues with automated tools or code generation that expect conventional struct field names.  
Severity: Low

## Location

All affected structs:

- `EnvironmentApplicationPackageInstallResourceModel` (field `Id`)
- `TenantApplicationPackageDataSourceModel` (fields `ApplicationId`, `PublisherId`)

## Code Issue

```go
	Id            types.String   `tfsdk:"id"`
	UniqueName    types.String   `tfsdk:"unique_name"`
	EnvironmentId types.String   `tfsdk:"environment_id"`
```

## Fix

Rename all "Id" suffixes to "ID" to follow Go conventions:

```go
	ID            types.String   `tfsdk:"id"`
	UniqueName    types.String   `tfsdk:"unique_name"`
	EnvironmentID types.String   `tfsdk:"environment_id"`
```

And similarly update other fields throughout the file.


---

## ISSUE 3

# Issue: Naming Convention for Field in Struct

##

/workspaces/terraform-provider-power-platform/internal/services/environment_groups/models.go

## Problem

The field `Id` in the `EnvironmentGroupResourceModel` struct does not follow Go naming conventions according to commonly accepted practices for acronyms and initialisms. In Go, initialisms and acronyms should be written in all caps, so `Id` should be `ID`.

## Impact

Low. The code will work as expected but does not adhere to Go naming conventions, which may cause minor confusion or inconsistencies throughout the codebase, especially when integrating with other packages or when using code generation tools.

## Location

```go
type EnvironmentGroupResourceModel struct {
	Id          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
}
```

## Code Issue

```go
	Id          types.String `tfsdk:"id"`
```

## Fix

Rename the field from `Id` to `ID` to adhere to Go's naming conventions for acronyms and initialisms.

```go
	ID          types.String `tfsdk:"id"`
```

Apply on the whole codebase 


---

## ISSUE 4

# Issue with Inconsistent Naming: "ComponetType" instead of "ComponentType"

##

/workspaces/terraform-provider-power-platform/internal/services/managed_environment/models.go

## Problem

In the `SolutionCheckerRuleDto` struct, the field is named `ComponetType`, which appears to be a typo. The correct spelling should be `ComponentType`. This could cause confusion for anyone consuming the JSON API or managing this struct, as the field does not accurately describe its purpose and introduces risk of inconsistent access or bugs related to misnaming.

## Impact

Incorrect naming impacts code clarity and maintainability. It can also lead to potential serialization/deserialization issues, making it harder for consumers to understand or utilize this struct. The issue severity is **low**, but addressing this improves professionalism and reduces future technical debt.

## Location

```go
type SolutionCheckerRuleDto struct {
    // ...
    ComponetType    string `json:"componetType,omitempty"`
    // ...
}
```

## Code Issue

```go
ComponetType    string `json:"componetType,omitempty"`
```

## Fix

The field name and its JSON tag should be corrected to `ComponentType`:

```go
ComponentType   string `json:"componentType,omitempty"`
```

Update references to this field everywhere in your codebase, not just in this file, to prevent mismatches between naming in code and serialized JSON data.


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
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
