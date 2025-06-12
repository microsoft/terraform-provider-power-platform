# API Response and DTO Nil Pointer Handling Issues

This document consolidates all identified nil pointer dereference issues related to API responses and DTO handling in the Terraform Provider for Power Platform.

## ISSUE 1

**Title**: Error Handling for `resp` Possibly Nil

**File**: `/workspaces/terraform-provider-power-platform/internal/services/environment_settings/api_environment_settings.go`

**Problem**: In the `UpdateEnvironmentSettings` method, there is a check on `resp != nil && resp.HttpResponse.StatusCode == http.StatusInternalServerError`. However, `resp` may be nil, but subsequent functions such as `client.Api.HandleForbiddenResponse(resp)` and `client.Api.HandleNotFoundResponse(resp)` are called without verifying this, which may lead to nil pointer dereferences.

**Impact**: High. This can cause runtime panics if `resp` is nil.

**Location**:

```go
if resp != nil && resp.HttpResponse.StatusCode == http.StatusInternalServerError {
    return nil, customerrors.WrapIntoProviderError(nil, customerrors.ERROR_ENVIRONMENT_SETTINGS_FAILED, string(resp.BodyAsBytes))
}
if err := client.Api.HandleForbiddenResponse(resp); err != nil {
    return nil, err
}
if err := client.Api.HandleNotFoundResponse(resp); err != nil {
    return nil, err
}
if err != nil {
    return nil, err
}
```

**Fix**: Return early if `resp` is nil to avoid calling methods on a nil pointer.

```go
if resp == nil {
    return nil, fmt.Errorf("response is nil")
}
if resp.HttpResponse.StatusCode == http.StatusInternalServerError {
    return nil, customerrors.WrapIntoProviderError(nil, customerrors.ERROR_ENVIRONMENT_SETTINGS_FAILED, string(resp.BodyAsBytes))
}
if err := client.Api.HandleForbiddenResponse(resp); err != nil {
    return nil, err
}
if err := client.Api.HandleNotFoundResponse(resp); err != nil {
    return nil, err
}
if err != nil {
    return nil, err
}
```

## ISSUE 2

**Title**: Missing Nil Check for `savedRoleData` Result

**File**: `/workspaces/terraform-provider-power-platform/internal/services/authorization/api_user.go`

**Problem**: In the `RemoveEnvironmentUserSecurityRoles` function, the variable `savedRoleData` is assigned by calling `array.Find` for each role in `securityRoles`. The code proceeds to use `savedRoleData.Name` without checking whether `savedRoleData` is nil. If `array.Find` does not find a matching role, this will result in a runtime panic (nil pointer dereference).

**Impact**: Severity: Critical

If the input lists are inconsistent or data is missing/corrupt, the provider can crash the process, causing cascading failures and overall instability of the Terraform execution and user workflow. This is especially important in infrastructure automation where reliability is paramount.

**Code Issue**:

```go
for _, role := range securityRoles {
 savedRoleData := array.Find(savedRoles, func(roleDto securityRoleDto) bool {
  return roleDto.RoleId == role
 })

 remove.Remove = append(remove.Remove, RoleDefinitionDto{
  Id: fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s/roleAssignments/%s", environmentId, savedRoleData.Name),
 })
}
```

**Fix**: Check for nil before accessing fields and return an appropriate error if `savedRoleData` is not found.

```go
for _, role := range securityRoles {
 savedRoleData := array.Find(savedRoles, func(roleDto securityRoleDto) bool {
  return roleDto.RoleId == role
 })
 if savedRoleData == nil {
  return nil, fmt.Errorf("security role with ID %s not found in savedRoles", role)
 }
 remove.Remove = append(remove.Remove, RoleDefinitionDto{
  Id: fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s/roleAssignments/%s", environmentId, savedRoleData.Name),
 })
}
```

## ISSUE 3

**Title**: Missing nil checks for nested structs in response handling

**File**: `/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity.go`

**Problem**: When mapping values from `capacity.Consumption` in the Read function, the code assumes that `capacity.Consumption` is never nil. If the API response contains a capacity object without a `Consumption` object, this will cause a runtime panic due to dereferencing a nil pointer.

**Impact**: High severity. This can cause the entire provider operation to fail with a panic if upstream data does not guarantee non-nil `Consumption` fields. This is both a stability and safety concern.

**Location**: Function: `Read`, during state mapping

**Code Issue**:

```go
Consumption: ConsumptionDataSourceModel{
    Actual:          types.Float32Value(capacity.Consumption.Actual),
    Rated:           types.Float32Value(capacity.Consumption.Rated),
    ActualUpdatedOn: types.StringValue(capacity.Consumption.ActualUpdatedOn),
    RatedUpdatedOn:  types.StringValue(capacity.Consumption.RatedUpdatedOn),
},
```

**Fix**: Add a nil check before accessing fields of `capacity.Consumption`:

```go
var consumption ConsumptionDataSourceModel
if capacity.Consumption != nil {
    consumption = ConsumptionDataSourceModel{
        Actual:          types.Float32Value(capacity.Consumption.Actual),
        Rated:           types.Float32Value(capacity.Consumption.Rated),
        ActualUpdatedOn: types.StringValue(capacity.Consumption.ActualUpdatedOn),
        RatedUpdatedOn:  types.StringValue(capacity.Consumption.RatedUpdatedOn),
    }
}
// ...
Consumption: consumption,
```

## ISSUE 4

**Title**: Lack of Null Checks for Nested Pointers in DTO Access

**File**: `/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/dto.go`

**Problem**: When accessing nested pointer fields (e.g., `tenantSettingsDto.PowerPlatform.Search.DisableDocsSearch`), many places in the code presume non-nil pointers. If a parent pointer (`PowerPlatform` or `Search`) is nil, this will cause a panic due to dereferencing a nil pointer.

**Impact**: Dereferencing nil pointers leads to panics, causing a provider crash and potential Terraform workflow interruption. This is a critical reliability and user experience problem. Severity: high.

**Location**: Example in `convertSearchSettings` and similar conversion functions:

```go
if tenantSettingsDto.PowerPlatform == nil || tenantSettingsDto.PowerPlatform.Search == nil {
    return types.ObjectType{AttrTypes: attrTypesSearchProperties}, types.ObjectNull(attrTypesSearchProperties)
}
attrValuesSearchProperties := map[string]attr.Value{
    "disable_docs_search":       types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Search.DisableDocsSearch),
    ...
}
```

**Code Issue**:

```go
tenantSettingsDto.PowerPlatform.Search.DisableDocsSearch
```

**Fix**: Always check for nil at each pointer access level before dereferencing:

```go
if tenantSettingsDto.PowerPlatform != nil && tenantSettingsDto.PowerPlatform.Search != nil {
    // safe to access tenantSettingsDto.PowerPlatform.Search.DisableDocsSearch
}
```

## ISSUE 5

**Title**: Lack of Error Handling in Data Conversion

**File**: `/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/models.go`

**Problem**: The function `convertDtoToModel` directly accesses fields of pointer arguments (e.g., `dto.ID`, `dto.Sink.ID`, etc.) and slices of pointers (e.g., `dto.Environments`, `dto.Scenarios`) without additional nil-checks or error handling. If fields inside nested structs or slices are unexpectedly nil, this could lead to panics due to nil pointer dereferencing.

**Impact**: Severity: **High**

Any unexpected nil fields returned from backend APIs could cause the provider to panic and crash, resulting in errors in Terraform runs and possibly leading to loss of provider state.

**Code Issue**:

```go
 return &AnalyticsDataModel{
  ID:           types.StringValue(dto.ID),
  Source:       types.StringValue(dto.Source),
  Environments: environments,
  Status:       status,
  Sink: SinkModel{
   ID:                types.StringValue(dto.Sink.ID),
   Type:              types.StringValue(dto.Sink.Type),
   SubscriptionId:    types.StringValue(dto.Sink.SubscriptionId),
   ResourceGroupName: types.StringValue(dto.Sink.ResourceGroupName),
   ResourceName:      types.StringValue(dto.Sink.ResourceName),
   Key:               types.StringValue(dto.Sink.Key),
  },
  PackageName:      types.StringValue(dto.PackageName),
  Scenarios:        scenarios,
  ResourceProvider: types.StringValue(dto.ResourceProvider),
  AiType:           types.StringValue(dto.AiType),
 }
```

**Fix**: Check for nil pointers before accessing struct fields:

```go
 var sink SinkModel
 if dto.Sink != nil {
  sink = SinkModel{
   ID:                types.StringValue(dto.Sink.ID),
   Type:              types.StringValue(dto.Sink.Type),
   SubscriptionId:    types.StringValue(dto.Sink.SubscriptionId),
   ResourceGroupName: types.StringValue(dto.Sink.ResourceGroupName),
   ResourceName:      types.StringValue(dto.Sink.ResourceName),
   Key:               types.StringValue(dto.Sink.Key),
  }
 } else {
  sink = SinkModel{
   ID:                types.StringNull(),
   Type:              types.StringNull(),
   SubscriptionId:    types.StringNull(),
   ResourceGroupName: types.StringNull(),
   ResourceName:      types.StringNull(),
   Key:               types.StringNull(),
  }
 }

 return &AnalyticsDataModel{
  ID:           types.StringValue(dto.ID),
  Source:       types.StringValue(dto.Source),
  Environments: environments,
  Status:       status,
  Sink:         sink,
  PackageName:      types.StringValue(dto.PackageName),
  Scenarios:        scenarios,
  ResourceProvider: types.StringValue(dto.ResourceProvider),
  AiType:           types.StringValue(dto.AiType),
 }
```

## ISSUE 6

**Title**: Missing error return after detecting nil `share` in Create, Read, and Update methods

**File**: `/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection_share.go`

**Problem**: In the `Create`, `Read`, and `Update` methods, when a nil `share` or `newShare` is detected, an error is logged to `resp.Diagnostics` but the control flow continues to run and may attempt to access properties on a nil pointer. For example, after checking `if share == nil`, it should immediately return rather than continuing on to `convertFromConnectionResourceSharesDto(plan, share)` which will panic if `share` is `nil`.

**Impact**: This introduces a risk of panics due to nil pointer dereferences if the `share` or `newShare` objects are missing. The severity is **high** since it will cause the provider to crash during runtime.

**Code Issue**:

```go
 if share == nil {
  resp.Diagnostics.AddError("Error getting connection share", "Connection share not found")
 }
 // code continues and uses "share" without returning

 // ... also in Read and Update, same problem with "share" and "newShare"
```

**Fix**: Add a `return` statement immediately after adding the error to diagnostics for these nil checks:

```go
 if share == nil {
  resp.Diagnostics.AddError("Error getting connection share", "Connection share not found")
  return
 }
```

## ISSUE 7

**Title**: Possible nil pointer dereference if GovernanceConfiguration.Settings is nil

**File**: `/workspaces/terraform-provider-power-platform/internal/services/managed_environment/resource_managed_environment.go`

**Problem**: Within the Create and Update handlers, after a successful client call, the returned environment object's `env.Properties.GovernanceConfiguration.Settings` is dereferenced and fields accessed directly, assuming it is always non-nil. However, this assumption may not hold if the backend data does not provide configuration settings as expected—e.g., in freshly created environments, misprovisioned states, or certain error conditions. A nil pointer dereference would result in a panic and crash the Terraform provider process.

**Impact**: High severity. A nil pointer panic will crash Terraform runs and could destroy in-progress state. The risk is particularly notable during environment churn, failures or API changes.

**Code Issue**:

```go
maxLimitUserSharing, _ := strconv.ParseInt(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MaxLimitUserSharing, 10, 64)
plan.IsUsageInsightsDisabled = types.BoolValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.ExcludeEnvironmentFromAnalysis == "true")
// ... similar lines follow
```

**Fix**: Check `env.Properties.GovernanceConfiguration.Settings` for nil before dereferencing:

```go
settings := env.Properties.GovernanceConfiguration.Settings
if settings == nil {
    resp.Diagnostics.AddError("Missing GovernanceConfiguration.Settings after environment provisioning", "API response did not include configuration settings. This might indicate an API or state consistency issue. Please inspect the backend state or retry.")
    return
}
// then proceed to access settings.ExtendedSettings ...
```

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
