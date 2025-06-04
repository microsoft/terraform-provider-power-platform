# Diagnostic Usage Issues

This document consolidates all issues related to improper diagnostic usage and error handling in the Terraform Provider for Power Platform.

## ISSUE 1

### Missing Error Detail/Wrapping on GetConnectionShares API Error

**File:** `/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connection_shares.go`

**Problem:** When the GetConnectionShares API call fails, only the error string is given in diagnostics, without parameter context (environment, connector, connection).

**Impact:** **Severity: medium** - Makes debugging harder for operators and maintainers.

**Location:**

```go
if err != nil {
 resp.Diagnostics.AddError("Failed to get connection shares", err.Error())
 return
}
```

**Code Issue:**

```go
if err != nil {
 resp.Diagnostics.AddError("Failed to get connection shares", err.Error())
 return
}
```

**Fix:** Wrap/extend the error with more parameter information:

```go
if err != nil {
 resp.Diagnostics.AddError(
  fmt.Sprintf(
   "Failed to get connection shares for environment_id '%s', connector_name '%s', connection_id '%s'",
   state.EnvironmentId.ValueString(), state.ConnectorName.ValueString(), state.ConnectionId.ValueString(),
  ),
  err.Error(),
 )
 return
}
```

## ISSUE 2

### Improper Diagnostic Usage for Error Handling in Helper Functions

**File:** `/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/helpers.go`

**Problem:** In both `convertToDlpConnectorGroup` and `convertToDlpCustomConnectorUrlPatternsDefinition`, diagnostic errors are added using `diags.AddError` when decoding attributes, but the diagnostic is not checked or returned properly. This can lead to partially incorrect data being created if an error occurs, as the function will proceed and return an incomplete or default structure.

**Impact:** High severity. Errors encountered during data marshalling or transformation are not propagated or handled meaningfully. This could lead to incorrect or incomplete outputs, and can cause subtle bugs which are difficult to track down during usage.

**Location:** Lines 128–150 (example from `convertToDlpConnectorGroup`):

**Code Issue:**

```go
if err != nil {
 diags.AddError("Client error when converting DlpConnectorGroups", "")
}
...
return connectorGroup
```

```go
if err != nil {
 diags.AddError("Client error when converting DlpCustomConnectorUrlPatternsDefinition", "")
}
...
return customConnectorUrlPatternsDefinition
```

**Fix:** Return early or propagate the error/diagnostics when an error is encountered. Example for `convertToDlpConnectorGroup`:

```go
func convertToDlpConnectorGroup(ctx context.Context, diags diag.Diagnostics, classification string, connectorsAttr basetypes.SetValue) (dlpConnectorGroupsModelDto, error) {
 var connectors []dataLossPreventionPolicyResourceConnectorModel
 err := connectorsAttr.ElementsAs(ctx, &connectors, true)
 if err != nil {
  diags.AddError("Client error when converting DlpConnectorGroups", err.Error())
  return dlpConnectorGroupsModelDto{}, err
 }
 // ...
 return connectorGroup, nil
}
```

## ISSUE 3

### Inconsistent Error Handling in convertColumnsToState

**File:** `/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record.go`

**Problem:** The `convertColumnsToState` function has some places where error values are deliberately ignored using `_` or omitted via no checks. This can hide bugs or introduce hard-to-debug issues.

**Impact:**

- **Severity:** High
- Can cause subtle or silent state bugs.
- Undermines error visibility and diagnosis.

**Location:** `convertColumnsToState`, especially in:

```go
columnField, _ := types.ObjectValue(attributeTypes, attributes)
```

**Code Issue:**

```go
columnField, _ := types.ObjectValue(attributeTypes, attributes)
result := types.DynamicValue(columnField)
return &result, nil
```

**Fix:** Check for the error and propagate or report it via diagnostics:

```go
columnField, err := types.ObjectValue(attributeTypes, attributes)
if err != nil {
    return nil, fmt.Errorf("failed to create object value: %w", err)
}
result := types.DynamicValue(columnField)
return &result, nil
```

## ISSUE 4

### Misleading Error Message When ProviderData is Not *api.ProviderClient

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment_settings/resources_environment_settings.go`

**Problem:** Within the `Configure` method, the code assumes that `req.ProviderData` is always of the type `*api.ProviderClient`. If a different type is provided, it tries to access its `Api` property which would lead to a panic. The error message in this case refers to an `*http.Client`, which is misleading and could cause confusion when debugging.

**Impact:** If the ProviderData is not the expected type, this will result in a panic. The error message refers to `*http.Client`, which is misleading because the expected value is a `*api.ProviderClient`. This could lead to confusion and make debugging harder. Severity: high.

**Code Issue:**

```go
 client := req.ProviderData.(*api.ProviderClient).Api

 if client == nil {
  resp.Diagnostics.AddError(
   "Unexpected Resource Configure Type",
   fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
  )

  return
 }
```

**Fix:** Add a type assertion check and improve the error message to clearly state the expected type:

```go
 providerClient, ok := req.ProviderData.(*api.ProviderClient)
 if !ok {
  resp.Diagnostics.AddError(
   "Unexpected Resource Configure Type",
   fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
  )
  return
 }

 client := providerClient.Api
 if client == nil {
  resp.Diagnostics.AddError(
   "Unexpected Resource Configure Type",
   "Provider client Api is nil. Please report this issue to the provider developers.",
  )
  return
 }
```

This fix avoids a panic and produces more accurate error messages, improving error handling and developer experience.

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
