# Datasource General Issues - Merged Issues

## ISSUE 1

# Functions and variable naming consistency

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connection_shares.go

## Problem

There are some naming inconsistencies and potential confusion in function and variable names, such as `ConvertFromConnectionSharesDto` (using PascalCase for a converting function, and mixing DTO and model names), `SharesDataSourceModel`/`SharesListDataSourceModel` (potentially redundant or confusing), and the difference between singular/plural (`share` vs `shares`).

## Impact

**Severity: low**

This may lead to confusion for future contributors, reduce maintainability, and increase the risk of subtle bugs or misunderstandings in usage.

## Location

- `ConvertFromConnectionSharesDto`
- `SharesDataSourceModel` vs `SharesListDataSourceModel`
- `NewConnectionSharesDataSource`

## Code Issue

```go
func ConvertFromConnectionSharesDto(connection shareConnectionResponseDto) SharesDataSourceModel
```

## Fix

Adopt consistent, idiomatic Go naming conventions:

- Use camelCase for local variables and PascalCase for exported types/functions.
- Stick to singular/plural conventions (e.g., `ShareDataSourceModel`, `SharesListDataSourceModel`).
- Use clearer function names, such as `convertConnectionShareDtoToModel`.

Example:

```go
func convertConnectionShareDTOToModel(dto shareConnectionResponseDto) ShareDataSourceModel
```


---

## ISSUE 2

# Issue: Function Naming - `ConvertFromConnectionDto` Is Not Idiomatic

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connections.go

## Problem

The function `ConvertFromConnectionDto` does not follow Go idioms for naming conversion methods. Go typically uses the form `toX` or `fromX` for conversion helpers, e.g., `connectionDtoToModel` or `toConnectionsDataSourceModel`.

## Impact

Severity: **Low**

While not a functional problem, non-idiomatic naming can make the codebase less consistent and harder to navigate for Go developers, especially in larger codebases.

## Location

```go
func ConvertFromConnectionDto(connection connectionDto) ConnectionsDataSourceModel
```

## Code Issue

```go
func ConvertFromConnectionDto(connection connectionDto) ConnectionsDataSourceModel
```

## Fix

Rename the function to follow Go conventions, such as:

```go
func connectionDtoToModel(connection connectionDto) ConnectionsDataSourceModel
```


---

## ISSUE 3

# Inconsistent Category Mapping for Connector Groups

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/datasource_dlp_policy.go

## Problem

In the iteration over `policies` in the `Read` function, there is a possible mistake in how the conversion helpers are called:

```go
policyModel.BusinessGeneralConnectors = convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
policyModel.NonBusinessConfidentialConnectors = convertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
```

Naming implies that "BusinessGeneralConnectors" should relate to "General" and "NonBusinessConfidentialConnectors" to "Confidential". However, the code uses the opposite mapping, potentially leading to incorrect data population.

## Impact

This could result in connectors being assigned to the wrong data group in the Terraform state, leading to incorrect resource configuration and potentially violating data loss prevention logic. Severity: **high**.

## Location

Function: `Read`
Lines where `convertToAttrValueConnectorsGroup` is called.

## Code Issue

```go
policyModel.BusinessGeneralConnectors = convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
policyModel.NonBusinessConfidentialConnectors = convertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
```

## Fix

Swap the category strings where the conversion functions are called, so the variable names match the values:

```go
policyModel.BusinessGeneralConnectors = convertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
policyModel.NonBusinessConfidentialConnectors = convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
```

**Explanation:**  
This aligns the data being set with the correct expected category as described by the field names, ensuring the right connectors are mapped.


---

## ISSUE 4

# Title

Ambiguous Type Name: `DataverseWebApiDatasourceModel`

##

/workspaces/terraform-provider-power-platform/internal/services/rest/datasource_rest_query.go

## Problem

The `DataverseWebApiDatasourceModel` type is referenced in the code (in the `Read` function) but its declaration is not included in the provided file. Assuming it is declared elsewhere or imported, the name is verbose and could be misaligned with Go's naming conventions or with its usage context. Names should be descriptive but concise and fit their purpose (state or schema model, etc.).

## Impact

Ambiguous or overly verbose names slow understanding and burden maintenance, especially when interleaved with other "model" types. This is a low-severity issue, mostly about maintainability and readability.

## Location

```go
var state DataverseWebApiDatasourceModel
```

## Code Issue

```go
var state DataverseWebApiDatasourceModel
```

## Fix

Rename to a more concise and context-specific name, such as:

```go
var state DataverseWebAPIState
```

Or, if the model is for schema state, prefix/suffix accordingly. Also ensure consistency with any type aliases or field names elsewhere in the provider.

---

This issue will be saved in:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/datasource_rest_query.go_model_naming_low.md`


---

## ISSUE 5

# Title

Type Naming: DataverseWebApiDatasource Has Inconsistent Casing

##

/workspaces/terraform-provider-power-platform/internal/services/rest/datasource_rest_query.go

## Problem

The main type in this file is named `DataverseWebApiDatasource`, which mixes "Api" and "Datasource" casing. In Go, the recommended convention based on acronym usage is either "API" or "Datasource"/"DataSource" but not partial. Additionally, the naming is inconsistent with the commonly used "DataSource" in the Terraform Plugin SDK community.

## Impact

Inconsistent or unconventional naming can cause confusion and decrease maintainability, as contributors might be unsure whether the correct form is "API", "Api", "Datasource", or "DataSource". This is a low-severity issue related to code clarity.

## Location

Throughout the file as the type name and references.

## Code Issue

```go
type DataverseWebApiDatasource struct {
  //...
}
```

## Fix

Rename the struct and relevant usages to `DataverseWebAPIDatasource` (or, if following the Go/TF convention strongly, `DataverseWebAPIDataSource`):

```go
type DataverseWebAPIDatasource struct {
  //...
}
```

And update all references accordingly. This will bring clarity and consistency to the codebase.

---

This issue will be saved in:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/datasource_rest_query.go_type_naming_low.md`


---

## ISSUE 6

# Title

Struct fields may not conform to Go naming conventions

##

/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity.go

## Problem

The struct fields used in the `state.TenantCapacities` slice and the corresponding `TenantCapacityDataSourceModel` and `ConsumptionDataSourceModel` structs (presumably defined elsewhere) use names like `CapacityType`, `CapacityUnits`, and `TotalCapacity`, which match the upstream JSON keys but should be capitalized and formatted according to Go naming conventions when used in Go struct definitions. 

Additionally, field names like `TenantId` should ideally be `TenantID` (per Go's initialism guidelines).

## Impact

Medium severity. While the code functions, improper naming can hinder readability and long-term maintainability, and doesn't conform to Go style guidelines, reducing clarity for future maintainers.

## Location

Throughout the read function (struct assignments)

## Code Issue

```go
state.TenantId = types.StringValue(tenantCapacityDto.TenantId)
// ...
TenantCapacityDataSourceModel{
    CapacityType:  types.StringValue(capacity.CapacityType),
    CapacityUnits: types.StringValue(capacity.CapacityUnits),
    // ...
}
```

## Fix

Follow Go's naming conventions for initialisms. For example, rename `TenantId` to `TenantID`, and similarly update struct field definitions and variable assignments. Ensure all references match the revised names.

```go
state.TenantID = types.StringValue(tenantCapacityDto.TenantID)
// ...
TenantCapacityDataSourceModel{
    CapacityType:  types.StringValue(capacity.CapacityType),
    CapacityUnits: types.StringValue(capacity.CapacityUnits),
    // ...
}
```

Note: This example assumes that the corresponding struct definitions are updated to match the revised naming.


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
