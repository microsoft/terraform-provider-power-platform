# Resource Redundant Code Issues

This document consolidates all redundant code issues found in resource components of the Terraform Power Platform provider.

## ISSUE 1

**File:** `/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record.go`

### Problem

Inside both `Create` and `Update` functions, the following lines reassign plan fields to themselves, resulting in no effective operation:

```go
plan.Id = types.StringValue(plan.Id.ValueString())
plan.EnvironmentId = types.StringValue(plan.EnvironmentId.ValueString())
plan.TableLogicalName = types.StringValue(plan.TableLogicalName.ValueString())
plan.Columns = types.DynamicValue(plan.Columns)
```

These assignments are likely vestigial, left-over from previous implementations or copy-paste. They are redundant since the `plan` fields already hold these values, and this pattern is repeated in both the `Create` and `Update` methods.

### Impact

- **Severity:** Low
- Causes unnecessary confusion and reduces code clarity.
- May cause maintainers to question if a side-effect is expected, leading to misunderstandings.
- Minor performance impact, though negligible.

### Location

In both `Create` (lines ~100-107) and `Update` (lines ~193-200).

### Code Issue

```go
plan.Id = types.StringValue(plan.Id.ValueString())
plan.EnvironmentId = types.StringValue(plan.EnvironmentId.ValueString())
plan.TableLogicalName = types.StringValue(plan.TableLogicalName.ValueString())
plan.Columns = types.DynamicValue(plan.Columns)
```

### Fix

Remove these unnecessary reassignments; the fields are already populated via `req.Plan.Get()` and do not need to be reset.

```go
// Remove these lines in both Create and Update methods:
// plan.Id = types.StringValue(plan.Id.ValueString())
// plan.EnvironmentId = types.StringValue(plan.EnvironmentId.ValueString())
// plan.TableLogicalName = types.StringValue(plan.TableLogicalName.ValueString())
// plan.Columns = types.DynamicValue(plan.Columns)
```

## ISSUE 2

**File:** `/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy.go`

### Problem

In both the `Create` and `Update` resource methods, there are redundant initializations of the `ConnectorGroups` and `CustomConnectorUrlPatternsDefinition` slices within the `dlpPolicyModelDto` struct. After initializing the slice as empty, you immediately overwrite it using the output of append operations. This is unnecessary and may lead to confusion or, if not overwritten, bugs arising from data loss or accidental use of stale slices.

### Impact

Severity: Low

This has a low performance impact and code clarity risk. While not functionally harmful due to the immediate overwrite, it clutters the code and can cause confusion about intentional initialization vs assignment semantics.

### Location

```go
policyToCreate := dlpPolicyModelDto{
 DefaultConnectorsClassification:      plan.DefaultConnectorsClassification.ValueString(),
 DisplayName:                          plan.DisplayName.ValueString(),
 EnvironmentType:                      plan.EnvironmentType.ValueString(),
 Environments:                         []dlpEnvironmentDto{},
 ConnectorGroups:                      []dlpConnectorGroupsModelDto{}, // unnecessary
 CustomConnectorUrlPatternsDefinition: []dlpConnectorUrlPatternsDefinitionDto{}, // unnecessary
}
// ...
policyToCreate.ConnectorGroups = make([]dlpConnectorGroupsModelDto, 0)
policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, convertToDlpConnectorGroup(ctx, resp.Diagnostics, "Confidential", plan.BusinessGeneralConnectors))
// repeats for \"General\" and \"Blocked\"
```

### Fix

Remove the redundant slice initializations in the struct literal and the redundant `make()` call:

```go
policyToCreate := dlpPolicyModelDto{
 DefaultConnectorsClassification: plan.DefaultConnectorsClassification.ValueString(),
 DisplayName:                     plan.DisplayName.ValueString(),
 EnvironmentType:                 plan.EnvironmentType.ValueString(),
 Environments:                    []dlpEnvironmentDto{},
 // Remove these redundant initializations:
 // ConnectorGroups:                      []dlpConnectorGroupsModelDto{},
 // CustomConnectorUrlPatternsDefinition: []dlpConnectorUrlPatternsDefinitionDto{},
}
// ...
// Remove: policyToCreate.ConnectorGroups = make([]dlpConnectorGroupsModelDto, 0)
policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, ...)
```

## ISSUE 3

**File:** `/workspaces/terraform-provider-power-platform/internal/services/application/resource_environment_application_package_install.go`

### Problem

In the `Create` method, after fetching the plan values, the code redundantly reassigns `EnvironmentId` and `UniqueName` to themselves, already as `types.StringValue(...)`. This is unnecessary unless the values are being normalized, which does not appear to be the case here.

### Impact

Severity: **Low**

Redundant assignments increase noise and can mislead readers to believe that the values are being processed when they're simply copied, slightly impacting code readability and maintainability.

### Location

Within `Create`:

```go
state.Id = types.StringValue(fmt.Sprintf("%s_%s", state.EnvironmentId.ValueString(), strings.ReplaceAll(strings.ToLower(state.UniqueName.ValueString()), " ", "_")))
state.EnvironmentId = types.StringValue(state.EnvironmentId.ValueString())
state.UniqueName = types.StringValue(state.UniqueName.ValueString())
```

### Fix

Remove the redundant assignments unless normalization or transformation is required. Only set these values if actual conversion, validation, or business logic is needed.

```go
state.Id = types.StringValue(fmt.Sprintf("%s_%s", state.EnvironmentId.ValueString(), strings.ReplaceAll(strings.ToLower(state.UniqueName.ValueString()), " ", "_")))
// Remove unnecessary copies of EnvironmentId and UniqueName unless transformation is needed
```

## ISSUE 4

**File:** `/workspaces/terraform-provider-power-platform/internal/services/solution/resource_solution_test.go`

### Problem

Inside `TestUnitSolutionResource_Validate_Create_And_Force_Recreate`, `httpmock.Activate()` and `defer httpmock.DeactivateAndReset()` are unnecessarily called twice at the beginning of the function. The second pair is redundant.

### Impact

Doesn't break functionality, but introduces code duplication, increases cognitive load, and may confuse readers or future maintainers. Severity: **low** (structure/maintainability).

### Location

Lines (approx.):

```go
func TestUnitSolutionResource_Validate_Create_And_Force_Recreate(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()
    ...
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()
    ...
}
```

### Fix

Remove the second pair of calls. Only one activation at the beginning and one defer at the end is needed:

```go
func TestUnitSolutionResource_Validate_Create_And_Force_Recreate(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()
    ...
    // Remove duplicate activate/deactivate block
}
```

## ISSUE 5

**File:** `/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go`

### Problem

In the `Create` function, `req.Plan.Get(ctx, &plan)` is called and diagnostics appended twice, one shortly after the other. This is redundant and unnecessary. The planned state should only be extracted once in a single code path unless there's a compelling reason to refresh `plan` between two stages (which does not seem to be the case here).

### Impact

This reduces code clarity and can lead to confusion about which instance of `plan` is in use. Severity: low.

### Location

Line ~170 ("resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)..." is called twice.)

### Code Issue

```go
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
 return
}

// ... (unrelated code skipped for clarity)

// Get the plan
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
 return
}
```

### Fix

Remove the redundant second call to `req.Plan.Get(ctx, &plan)`, keeping only the first instance.

```go
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
 return
}

// ... continue with rest of function, do not repeat req.Plan.Get(ctx, &plan)
```

## ISSUE 6

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment_settings/resources_environment_settings.go`

### Problem

In the Create, Read, and Update methods, after converting DTOs to resource models, both `state.Id` and `state.EnvironmentId` are explicitly set to `plan.EnvironmentId` or `state.EnvironmentId`. This appears redundant given that these should already be set correctly from the DTO mapping function. This practice can mask errors in DTO-to-model conversion, causing maintainers to not notice bugs if the conversion is wrong.

### Impact

Assigning these fields redundantly may reduce code clarity and hide bugs in DTO conversion. The severity is low, but it affects maintainability and could delay debugging type mapping mistakes.

### Location

```go
 state, err := convertFromEnvironmentSettingsDto[EnvironmentSettingsResourceModel](envSettings, plan.Timeouts)
 if err != nil {
  resp.Diagnostics.AddError("Error converting environment settings", err.Error())
  return
 }
 state.Id = plan.EnvironmentId
 state.EnvironmentId = plan.EnvironmentId
```

### Fix

Remove the redundant assignment unless DTO conversion does not (and should not) update these fields. Instead, ensure that the `convertFromEnvironmentSettingsDto` sets them properly during mapping.

```go
 state, err := convertFromEnvironmentSettingsDto[EnvironmentSettingsResourceModel](envSettings, plan.Timeouts)
 if err != nil {
  resp.Diagnostics.AddError("Error converting environment settings", err.Error())
  return
 }
 // Only set these fields here if there is a valid reason not to do so in conversion
```

Review and correct the DTO conversion if it does not set these essential identifiers correctly.

---

## To finish the task you have to

1. Run linter and fix any issues
2. Run UnitTest and fix any of failing ones
3. Generate docs
4. Run Changie

## Changie Instructions

Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```

Where:

- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number

Apply this fix to the whole codebase
