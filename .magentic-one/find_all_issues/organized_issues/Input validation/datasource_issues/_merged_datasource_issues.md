# Data Source Issues - Input Validation

This document contains all data source-related input validation issues found in the terraform-provider-power-platform codebase.


## ISSUE 1

# Type Safety: Lack of Validation on Downstream API Data

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/datasource_analytics_data_exports.go

## Problem

After calling the downstream API to retrieve analytics data exports, there is no validation of the content, shape, or completeness of the data before mapping it to the model. If the downstream API changes (e.g., omits required fields, returns nulls/unexpected values), this could lead to panics, zero-values, or subtle bugs communicated to end-users.

## Impact

- Could cause runtime panics if downstream returns malformed data
- Silent loss of information if fields become missing
- Reduces robustness, and user-facing errors become harder to diagnose

**Severity:** Medium

## Location

```go
analyticsDataExport, err := d.analyticsExportClient.GetAnalyticsDataExport(ctx)
// ...
for _, export := range *analyticsDataExport {
    if model := convertDtoToModel(&export); model != nil {
        exports = append(exports, *model)
    }
}
```

## Code Issue

No explicit type or value checks after obtaining data from downstream API.

## Fix

Add validation or sanity checks before mapping data (inside or prior to `convertDtoToModel`). Example:

```go
for _, export := range *analyticsDataExport {
    // Validate required fields are present and sane
    if export.ID == "" || export.Sink == nil {
        resp.Diagnostics.AddWarning("Incomplete data", fmt.Sprintf("Found analytics data export with nil/empty ID or Sink: %+v", export))
        continue
    }
    if model := convertDtoToModel(&export); model != nil {
        exports = append(exports, *model)
    }
}
```

Consider making `convertDtoToModel` return errors for invalid or inconsistent structures, or introducing a validation helper inline.


## ISSUE 2

# Missing Validation on Empty Required Attributes

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/datasource_data_record.go

## Problem

In the schema, fields such as `environment_id` and `entity_collection` are required, but there is no explicit check or validation for these attributes in the `Read` function beyond what Terraform does upstream. If the config retrieves an empty value for these attributes, subsequent calls to the API will fail with possibly less user-friendly errors, or downstream panics if those values are nil.

## Impact

- Weakens guarantees about input data consistency and robustness.
- Downstream errors may be less obvious to the user, resulting in poor UX and harder debugging.
- Severity: **medium**.

## Location

In `Read`, while processing inputs from config/state:

```go
tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", d.FullTypeName()))

resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

if resp.Diagnostics.HasError() {
	return
}

query, headers, err := BuildODataQueryFromModel(&config)
// ... continues
```

## Code Issue

No explicit check for empty required fields:

```go
queryRespnse, err := d.DataRecordClient.GetDataRecordsByODataQuery(ctx, config.EnvironmentId.ValueString(), query, headers)
```

## Fix

Explicitly check that required string attribute values are non-empty when reading and before making API calls:

```go
if config.EnvironmentId.IsUnknown() || config.EnvironmentId.IsNull() || config.EnvironmentId.ValueString() == "" {
	resp.Diagnostics.AddError(
		"Missing Environment ID",
		"The `environment_id` field must be provided and non-empty.",
	)
	return
}
if config.EntityCollection.IsUnknown() || config.EntityCollection.IsNull() || config.EntityCollection.ValueString() == "" {
	resp.Diagnostics.AddError(
		"Missing Entity Collection",
		"The `entity_collection` field must be provided and non-empty.",
	)
	return
}
```
Add similar checks for other required fields as needed. This provides better diagnostics and robustness.


## ISSUE 3

# Issue with Data Consistency When Building Rows

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/datasource_data_record.go

## Problem

When building the `rows` attribute in the `Read` method, the code constructs a slice of elements from records and builds a `TupleValue` using a parallel slice of their types. However, if `elements` has a length of zero (no records returned), the code calls `types.TupleValue(elementsTypes, elements)`, which may not behave as expected. Additionally, error return values from `types.TupleValue` and similar constructors are ignored, which could result in subtle or silent failures, especially if the input data has inconsistencies.

## Impact

- If the row count is zero, an invalid or nil value may be stored (unclear how framework handles zero length).
- Any internal error in attribute/type construction will be silently ignored (since errors are discarded).
- This weakens data consistency guarantees for consumers. 
- Severity: **medium**.

## Location

In `Read`:
```go
elementTypes := []attr.Type{}
for range elements {
	elementTypes = append(elementTypes, types.DynamicType)
}
rows, _ := types.TupleValue(elementTypes, elements)
state.Rows = types.DynamicValue(rows)
```

## Code Issue

```go
rows, _ := types.TupleValue(elementTypes, elements)
state.Rows = types.DynamicValue(rows)
```

## Fix

- Always check the error returned by `types.TupleValue` and handle it (at a minimum, surface the error to diagnostics and return).
- If there are no elements, set the value to empty appropriately using the Terraform SDK type helpers.

```go
rows, err := types.TupleValue(elementTypes, elements)
if err != nil {
	resp.Diagnostics.AddError("Failed to build tuple value for rows", err.Error())
	return
}
state.Rows = types.DynamicValue(rows)
```
This ensures that any error when constructing the attribute is reported, and consumers won't receive silently broken state.


## ISSUE 4

# Validators Applied to Computed-Only Schema Attributes

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/datasource_dlp_policy.go

## Problem

Attributes such as `default_connectors_classification` and fields nested under `custom_connectors_patterns` have validators defined, but are `Computed: true`, e.g.:

```go
"default_connectors_classification": schema.StringAttribute{
    MarkdownDescription: "Default classification for connectors (\"General\", \"Confidential\", \"Blocked\")",
    Computed:            true,
    Validators: []validator.String{
        stringvalidator.OneOf("General", "Confidential", "Blocked"),
    },
},
```

Validators are ignored for attributes that are not user-supplied (`Computed: true` and not `Optional`/`Required`). Including them is misleading and may signal to maintainers that some user input is being checked, which is not the case.

## Impact

Reduces clarity and can confuse other developers, who may believe the validator has runtime significance. It may also increase provider memory usage slightly and fails static checks in schema lint tools. Severity: **low**.

## Location

- Any attribute that is only `Computed` but includes validators.

## Code Issue

```go
"default_connectors_classification": schema.StringAttribute{
    ...
    Computed:            true,
    Validators: []validator.String{
        stringvalidator.OneOf("General", "Confidential", "Blocked"),
    },
},
// Also under custom_connectors_patterns.data_group, etc.
```

## Fix

Remove the `Validators` from `Computed`-only fields:

```go
"default_connectors_classification": schema.StringAttribute{
    MarkdownDescription: "Default classification for connectors (\"General\", \"Confidential\", \"Blocked\")",
    Computed:            true,
},
```

**Explanation:**  
Validators only affect user-supplied attributes, so don't add them unless the field is user-modifiable.


## ISSUE 5

# Title

Lack of input validation for `location` attribute in `Read` method

##

/workspaces/terraform-provider-power-platform/internal/services/languages/datasource_languages.go

## Problem

The `location` attribute from the Terraform state is used directly in API requests without validating if it is empty or malformed before making the call. This can result in unnecessary or erroneous API requests if input is incorrect or not set.

## Impact

This could cause confusing errors or unexpected results from the API, which would be harder to debug by the end user. It is a type safety and input validation issue, with **medium** severity.

## Location

```go
languages, err := d.LanguagesClient.GetLanguagesByLocation(ctx, state.Location.ValueString())
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
    return
}
```

## Code Issue

```go
languages, err := d.LanguagesClient.GetLanguagesByLocation(ctx, state.Location.ValueString())
```

## Fix

Add a check to ensure `state.Location.ValueString()` is not empty or invalid before calling the API. Provide a diagnostic error if invalid.

```go
location := state.Location.ValueString()
if location == "" {
    resp.Diagnostics.AddError(
        "Missing or invalid location",
        "The `location` attribute must be set and non-empty.",
    )
    return
}

languages, err := d.LanguagesClient.GetLanguagesByLocation(ctx, location)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
    return
}
```


## ISSUE 6

# Title

Missing error handling when extracting attribute from config

##

/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity.go

## Problem

In the `Read` function, the code calls `req.Config.GetAttribute(ctx, path.Root("tenant_id"), &tenantId)` to extract the tenant ID from the input configuration. However, it does not check or handle the returned diagnostics or errors from this method. If extraction fails, the function may proceed with an empty or invalid tenantId, resulting in unpredictable runtime errors or API rejections.

## Impact

High severity. Ignoring the result can produce misleading errors, cause API calls to fail, or make debugging difficult. Error handling of configuration extraction is crucial to ensure stable control flow.

## Location

Function: `Read`, during config attribute extraction

## Code Issue

```go
var tenantId string
req.Config.GetAttribute(ctx, path.Root("tenant_id"), &tenantId)
```

## Fix

Capture and handle the diagnostics appropriately:

```go
diags := req.Config.GetAttribute(ctx, path.Root("tenant_id"), &tenantId)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
    return
}
```

This ensures any extraction error halts further processing before making external calls.


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
