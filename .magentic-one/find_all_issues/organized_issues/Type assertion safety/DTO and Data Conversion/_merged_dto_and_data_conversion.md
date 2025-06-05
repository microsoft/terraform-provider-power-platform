# Type Assertion Safety Issues

This document contains merged type assertion safety issues found in the codebase.


## ISSUE 1

# Title

Missing Error Handling When Using Type Assertion (`.(basetypes.ObjectValue)`)

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/dto.go

## Problem

Multiple conversion functions in the file perform type assertions like `object.(basetypes.ObjectValue)` without checking for assertion failures. If the value does not implement `basetypes.ObjectValue`, a runtime panic will occur.

## Impact

Failure to check the result of the type assertion may cause a panic that crashes the provider at runtime. This is critical for a Terraform provider, as user input or upstream changes may invalidate the assumed interface, causing unexplained provider failures and loss of state. Severity: high.

## Location

- Functions like `convertTeamsIntegrationModel`, `convertPowerAppsModel`, `convertPowerAutomateModel`, etc.

## Code Issue

```go
var teamsIntegrationSettings TeamsIntegrationSettings
teamIntegrationObject.(basetypes.ObjectValue).As(ctx, &teamsIntegrationSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
```

## Fix

Use a type assertion with `ok` check and handle the error gracefully.

```go
objectValue, ok := teamIntegrationObject.(basetypes.ObjectValue)
if !ok {
    // handle the error: return, or log as appropriate
    return // or possibly log or add diagnostics
}
objectValue.As(ctx, &teamsIntegrationSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
```

Apply similar checks in all conversion functions where type assertions are used.


---

## ISSUE 2

# Unchecked Type Assertions can Panic

## 
/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/dto.go

## Problem

The code contains unchecked type assertions such as:

```go
aiGenerativeSettingsObj.(basetypes.ObjectValue).As(...)
```

If the value in `attrs["ai_generative_settings"]` is not of the expected type (`basetypes.ObjectValue`), this will cause a runtime panic, crashing the program. There are no checks (e.g., type assertion with second return value) to ensure the assertion succeeded before using the result.

## Impact

- **Severity:** High
- Unhandled panic paths can bring down the process (Terraform provider), possibly causing data loss or inconsistent state.
- Difficult to debug, as panics are abrupt and may not provide meaningful error messages to users or logs.
- Not idiomatic Go error handling.

## Location

Across all conversion functions. Example:

```go
if diags := aiGenerativeSettingsObj.(basetypes.ObjectValue).As(ctx, &aiGenerativeSettings, ...); diags != nil {
    ...
}
```

Similar patterns occur in other conversion helpers.

## Code Issue

```go
obj := attrs["ai_generative_settings"]
if !obj.IsNull() && !obj.IsUnknown() {
    var aiGenerativeSettings environmentGroupRuleSetAiGenerativeSettingsModel
    if diags := obj.(basetypes.ObjectValue).As(ctx, &aiGenerativeSettings, ...); diags != nil {
        ...
    }
}
```

## Fix

Use the "comma ok" type assertion form to ensure the type is correct before proceeding:

```go
obj := attrs["ai_generative_settings"]
if !obj.IsNull() && !obj.IsUnknown() {
    objVal, ok := obj.(basetypes.ObjectValue)
    if !ok {
        return fmt.Errorf("expected ai_generative_settings to be of type ObjectValue, got %T", obj)
    }
    var aiGenerativeSettings environmentGroupRuleSetAiGenerativeSettingsModel
    if diags := objVal.As(ctx, &aiGenerativeSettings, ...); diags != nil {
        ...
    }
}
```

Apply this change wherever type assertions are performed on `attr.Value` or similar interfaces.

---

---

## ISSUE 3

# Title

Potential type safety issue with `ElementsAs` usage and diagnostics

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/helpers.go

## Problem

Repeated patterns in this file use the `ElementsAs` method on `basetypes.SetValue`, and either ignore diagnostics, do not check error returns, or fail to return/gate the output on that basis. This is present in `convertToDlpEnvironment`, `getConnectorGroup`, and possibly other functions. This risks panics or data inconsistencies if the underlying type assertion/marshalling fails, but processing continues.

## Impact

Medium severity. Unexpected type assertion or unmarshalling failures could crash at runtime (`panic`), or result in data inconsistencies further down the stack.

## Location

Line 154-159, convertToDlpEnvironment:

```go
func convertToDlpEnvironment(ctx context.Context, environmentsInPolicy basetypes.SetValue) []dlpEnvironmentDto {
    envs := []string{}
    environmentsInPolicy.ElementsAs(ctx, &envs, true)
    ...
}
```

Line 110-113, getConnectorGroup:

```go
func getConnectorGroup(ctx context.Context, connectorsAttr basetypes.SetValue) (*dlpConnectorGroupsModelDto, error) {
    var connectors []dataLossPreventionPolicyResourceConnectorModel
    if diags := connectorsAttr.ElementsAs(ctx, &connectors, true); diags != nil {
        return nil, fmt.Errorf("error converting elements: %v", diags)
    }
```

## Code Issue

```go
environmentsInPolicy.ElementsAs(ctx, &envs, true)
```

## Fix

Always check the error return or diagnostics and propagate as needed, or at minimum, gracefully handle failures before further usage. For example, for `convertToDlpEnvironment`:

```go
func convertToDlpEnvironment(ctx context.Context, environmentsInPolicy basetypes.SetValue) ([]dlpEnvironmentDto, error) {
    envs := []string{}
    if err := environmentsInPolicy.ElementsAs(ctx, &envs, true); err != nil {
        return nil, err
    }
    ...
    return environments, nil
}
```

Be consistent in all uses of `ElementsAs` throughout the file.

---

## ISSUE 4

# Partial Error Handling for Type Assertions in caseX Functions

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record.go

## Problem

The helper functions `caseMapStringOfAny`, `caseArrayOfAny`, etc., only populate values if the type assertion succeeds but do not log or report if assertions fail. The calling code has no idea if the value could not be set. This may result in silent data loss or state drift.

## Impact

- **Severity:** Medium
- Can lead to silent errors or ignored/empty attributes.
- Reduces code robustness and makes debugging more difficult.

## Location

Functions: `caseMapStringOfAny`, `caseArrayOfAny`, etc.

## Code Issue

```go
value, ok := columnValue.(string)
if ok {
	// ...
}
```

## Fix

Consider logging or returning an error (if critical), or at minimum, logging via tflog for unexpected value types:

```go
value, ok := columnValue.(string)
if !ok {
    tflog.Debug(context.TODO(), "caseMapStringOfAny: failed to cast value to string", map[string]interface{}{ "key": key })
    return // or capture error for diagnostic
}
// ... (continue existing logic)
```

---

## ISSUE 5

# Title

Non-Idiomatic Type Assertion in Equal Method

##

/workspaces/terraform-provider-power-platform/internal/customtypes/uuid_type.go

## Problem

In the `Equal` method, the type assertion is performed against `UUIDType` (a value type), not against a pointer (`*UUIDType`). This will fail if the `attr.Type` instance is a pointer (which is common in Go when working with interfaces). This may result in false negatives when checking equality, even if both are logically the same type.

## Impact

Severity: **Medium**

- Equality checks may erroneously return false even when types are effectively equal, leading to subtle bugs, especially as pointer/value receivers are mixed or if interface implementations are refactored.
- May cause unexpected behavior when the framework expects equality logic to work reliably.

## Location

`func (t UUIDType) Equal(o attr.Type) bool`

```go
	other, ok := o.(UUIDType)
	if !ok {
		return false
	}
```

## Code Issue

```go
	other, ok := o.(UUIDType)
	if !ok {
		return false
	}
```

## Fix

- Use pointer receivers for the method and assert to `*UUIDType`, or enhance support for both pointer and value types:

Example using pointer receiver and asserting both ways:

```go
func (t *UUIDType) Equal(o attr.Type) bool {
	other, ok := o.(*UUIDType)
	if !ok {
		return false
	}
	return t.StringType.Equal(other.StringType)
}
```

Or, support both forms:

```go
	var other *UUIDType
	switch v := o.(type) {
	case UUIDType:
		other = &v
	case *UUIDType:
		other = v
	default:
		return false
	}
	return t.StringType.Equal(other.StringType)
```

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
