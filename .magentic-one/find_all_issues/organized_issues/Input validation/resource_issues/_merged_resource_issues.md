# Resource Issues - Input Validation

This document contains all resource-related input validation issues found in the terraform-provider-power-platform codebase.


## ISSUE 1

# Inadequate Handling of nil/Zero State in Read, Delete

##

/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/resource_admin_management_application.go

## Problem

In both the `Read` and `Delete` methods, the code assumes that the state object is successfully populated by:

```go
resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
if resp.Diagnostics.HasError() {
    return
}
```

However, beyond checking for `Diagnostics.HasError()`, there is no validation that `state.Id` is actually set or that it passes UUID format requirements. This could enable runtime errors (nil dereference, bad API call) if the state is malformed or incomplete.

## Impact

Severity: **medium**

This could result in crashes or unpredictable calls to the API client, especially if Terraform’s state is externally modified, corrupted, or not correctly managed.

## Location

Read and Delete methods:

```go
var state AdminManagementApplicationResourceModel
resp.Diagnostics.Append(req.State.Get(ctx, &state)...) 
if resp.Diagnostics.HasError() { return }

// no further check of state.Id
adminApp, err := r.AdminManagementApplicationClient.GetAdminApplication(ctx, state.Id.ValueString())
err := r.AdminManagementApplicationClient.UnregisterAdminApplication(ctx, state.Id.ValueString())
```

## Fix

After retrieving state, validate that `state.Id` is set and is non-zero (and optionally validate UUID correctness):

```go
if state.Id.IsNull() || state.Id.ValueString() == "" {
    resp.Diagnostics.AddError("Missing or invalid ID in state", "Cannot perform operation: resource ID is not set or invalid.")
    return
}
```

This prevents awkward downstream errors and provides clearer diagnostics.


## ISSUE 2

# No Assertion Comments or Grouping in resource.ComposeAggregateTestCheckFunc

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy_test.go

## Problem

Assertions are listed one after another with minimal distinction, making it difficult to see which assertion validates what scenario in the resource.

## Impact

- **Readability**: Difficult to trace which assertion relates to which resource.
- **Debuggability**: Harder to track test failures to actual intent.

**Severity: Low**

## Location

All `Check: resource.ComposeAggregateTestCheckFunc(...)` invocations.

## Code Issue

```go
Check: resource.ComposeAggregateTestCheckFunc(
    resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", "00000000-0000-0000-0000-000000000001"),
    resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "name", "payAsYouGoBillingPolicyExample"),
    // ... more attributes
),
```

## Fix

Group assertions by logical sections and add comments, or, better yet, split into mini helper check functions for better intent expression.

```go
Check: resource.ComposeAggregateTestCheckFunc(
    // ID Checks
    resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", "expected-id"),
    // Core Attribute Checks
    resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "name", "expected-name"),
    resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "status", "Enabled"),
    // Billing Instrument Checks
    resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "billing_instrument.resource_group", "expected-group"),
),
```


## ISSUE 3

# Title

Lack of helper-based field extraction or resource model validation causes repetitive code

##

/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/resource_copilot_studio_application_insights.go

## Problem

In all CRUD methods, code that extracts values from the plan or state is repeated, and type/semantic validation and error reporting is manual. This creates repetitive, verbose, and less maintainable code, and makes it easy for developers to miss validations or introduce subtle bugs.

## Impact

Severity: **Low**

This is a maintainability and readability problem. Without a field validation or extraction helper or resource model validator, the codebase is harder to audit for correctness and harder to update/typesafe in the future.

## Location

- `Create`, `Read`, `Update`, `Delete` methods.

## Code Issue

```go
var plan *ResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
	return
}
appInsightsConfigToCreate, err := createAppInsightsConfigDtoFromSourceModel(*plan)
if err != nil {
	resp.Diagnostics.AddError("Error when converting source model to create Copilot Studio Application Insights configuration dto", err.Error())
	return
}
```

## Fix

Centralize validation and extraction logic for resource models. For example, add a `Validate()` method to `ResourceModel`, or write a resource model validation helper:

```go
func (m *ResourceModel) Validate() error {
	if m.BotId.IsUnknown() || m.BotId.ValueString() == "" {
		return fmt.Errorf("bot_id is required")
	}
	// Check other fields as necessary
	return nil
}

// Usage:
if err := plan.Validate(); err != nil {
	resp.Diagnostics.AddError("Invalid resource model", err.Error())
	return
}
```

This reduces code repetition, increases clarity, and ensures that all validation paths are handled consistently, simplifying later refactorings.

---

File to save:  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/resource_copilot_studio_application_insights_model_validation_low.md`


## ISSUE 4

# Title

No validation of required resource attributes before use

##

/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/resource_copilot_studio_application_insights.go

## Problem

Throughout the CRUD methods, the plan and state models are accessed via pointer dereference (e.g., `plan.BotId.ValueString()` and `state.BotId.ValueString()`), but there is no explicit validation that these required fields are non-empty or valid before use, especially just after fetching from the plan or state.

If an invalid Terraform configuration, a provider bug, or a state migration issue leads to these being empty, subsequent downstream client calls could fail in an uncontrolled way.

## Impact

Severity: **Medium**

Lack of validation can lead to cryptic API errors or panics that do not give a user-friendly error message in Diagnostics. A more robust approach that validates all required attributes before invoking an API improves maintainability, user experience, and robustness.

## Location

- Methods: `Create`, `Read`, `Update`, `Delete`
- Usage: Immediately after plan/state is loaded and before API client methods

## Code Issue

```go
var plan *ResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

// No checks here!
// Directly using plan.BotId.ValueString()
```

## Fix

Explicitly check all required fields in the plan (or state) after loading and before use, and return a helpful error if they are not set.

```go
if plan.BotId.IsUnknown() || plan.BotId.ValueString() == "" {
	resp.Diagnostics.AddError("Missing Bot ID", "The bot_id field is required but was not provided.")
	return
}
if plan.EnvironmentId.IsUnknown() || plan.EnvironmentId.ValueString() == "" {
	resp.Diagnostics.AddError("Missing Environment ID", "The environment_id field is required but was not provided.")
	return
}
// Repeat as needed for other critical fields
```

Repeat for the `state` variable in Read and Delete, and for all critical variables wherever appropriate.

---

File to save:  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/type_safety/resource_copilot_studio_application_insights_required_fields_medium.md`


## ISSUE 5

# Title

Incorrect usage of field names in ValidateConfig (BusinessGeneralConnectors, NonBusinessConfidentialConnectors)

##
/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy.go

## Problem

The function `ValidateConfig` references struct fields `BusinessGeneralConnectors`, `NonBusinessConfidentialConnectors`, and `BlockedConnectors`. However, the model schema, as defined in `Schema`, uses the field names `business_connectors`, `non_business_connectors`, and `blocked_connectors`. The use of incorrect field names will result in runtime errors or unexpected behavior as the fields may not be populated or marshaled from Terraform input properly.

## Impact

Severity: High

This issue impacts the correctness of validation logic and may prevent configuration validation on Terraform plans. Any validation using these incorrect names will not operate as expected and may result in misconfigured DLP policies.

## Location

```go
var connectors []dlpConnectorModelDto
conn, err := getConnectorGroup(ctx, config.BusinessGeneralConnectors)
...
conn, err = getConnectorGroup(ctx, config.NonBusinessConfidentialConnectors)
...
conn, err = getConnectorGroup(ctx, config.BlockedConnectors)
...
```

## Code Issue

```go
var connectors []dlpConnectorModelDto
conn, err := getConnectorGroup(ctx, config.BusinessGeneralConnectors)
if err != nil {
	resp.Diagnostics.AddError("BusinessGeneralConnectors validation error", err.Error())
}
connectors = append(connectors, conn.Connectors...)

conn, err = getConnectorGroup(ctx, config.NonBusinessConfidentialConnectors)
if err != nil {
	resp.Diagnostics.AddError("NonBusinessConfidentialConnectors validation error", err.Error())
}
connectors = append(connectors, conn.Connectors...)

conn, err = getConnectorGroup(ctx, config.BlockedConnectors)
if err != nil {
	resp.Diagnostics.AddError("BlockedConnectors validation error", err.Error())
}
connectors = append(connectors, conn.Connectors...)
```

## Fix

Replace `BusinessGeneralConnectors` with `BusinessConnectors`, `NonBusinessConfidentialConnectors` with `NonBusinessConnectors` to align with the schema and struct field names.

```go
var connectors []dlpConnectorModelDto
conn, err := getConnectorGroup(ctx, config.BusinessConnectors)
if err != nil {
	resp.Diagnostics.AddError("BusinessConnectors validation error", err.Error())
}
connectors = append(connectors, conn.Connectors...)

conn, err = getConnectorGroup(ctx, config.NonBusinessConnectors)
if err != nil {
	resp.Diagnostics.AddError("NonBusinessConnectors validation error", err.Error())
}
connectors = append(connectors, conn.Connectors...)

conn, err = getConnectorGroup(ctx, config.BlockedConnectors)
if err != nil {
	resp.Diagnostics.AddError("BlockedConnectors validation error", err.Error())
}
connectors = append(connectors, conn.Connectors...)
```
---
Save as:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/resource_dlp_policy_incorrect_field_name_high.md


## ISSUE 6

# Title

Struct Field and Schema Attribute Name Mismatches

##
/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy.go

## Problem

Across multiple CRUD operations (Create, Read, Update, Delete), and the validation logic, there is a recurring mismatch between Go struct/model field names and Terraform resource schema attributes. Sometimes the code uses names like `BusinessGeneralConnectors`, `NonBusinessConfidentialConnectors` which do not correspond to the schema attributes (`business_connectors`, `non_business_connectors`). This is likely due to Go struct fields being named differently than the schema attributes, and the downstream usage (converters, API DTOs) required field names to match the schema. This increases cognitive complexity, introduces subtle bugs, and makes future maintenance difficult.

## Impact

Severity: High

Severe maintainability and reliability risk. Using mismatched names between Terraform schema and Go struct fields will lead to bugs, field omission, data conversion failures, and runtime panics that are difficult to debug.

## Location

Throughout the file, especially:

- In ValidateConfig: using `config.BusinessGeneralConnectors`
- In Read/Create/Update: using `convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)` (semantic confusion)
- In Schema: attribute is `business_connectors`, struct field may be `BusinessGeneralConnectors` or similar.

## Code Issue

```go
// Schema uses "business_connectors"
"business_connectors": schema.SetNestedAttribute{...}

// Elsewhere, uses BusinessGeneralConnectors, NonBusinessConfidentialConnectors, etc.
plan.BusinessGeneralConnectors // or config.BusinessGeneralConnectors
```

## Fix

Ensure that struct field names in the resource model (`dataLossPreventionPolicyResourceModel`) exactly match the resource schema attribute names (`business_connectors`, `non_business_connectors`, etc.) and update all usage accordingly. Refactor conversion helpers to align semantically and reduce confusion by creating a direct mapping between schema, struct, and API DTO.

```go
// In your resourceModel struct
type dataLossPreventionPolicyResourceModel struct {
  // ...
  BusinessConnectors types.Set   `tfsdk:\"business_connectors\"`
  NonBusinessConnectors types.Set `tfsdk:\"non_business_connectors\"`
  BlockedConnectors types.Set     `tfsdk:\"blocked_connectors\"`
  // ...
}
```

---
Save as:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/resource_dlp_policy_schema_model_mismatch_high.md


## ISSUE 7

# Title

Unchecked return values and silent error ignoring in strconv.ParseInt

##

/workspaces/terraform-provider-power-platform/internal/services/managed_environment/resource_managed_environment.go

## Problem

In multiple parts of the file, the `strconv.ParseInt` function's error return value is ignored by assigning it to `_`. Ignoring the error may lead to unexpected behaviors if the input string is not a valid integer, leading to unintentional zero values being interpreted and used further in the logic of the provider. This can result in incorrect resource state propagation or masking underlying data/formatting bugs.

## Impact

If the conversion fails, the `maxLimitUserSharing` will be zero and errors will be silently swallowed, potentially producing incorrect Terraform state or misconfigurations without indication to the user or system maintainers. This is a high severity issue because it can corrupt infrastructure state.

## Location

Notably seen in the following code snippet in Create, Update, and Read methods:

## Code Issue

```go
maxLimitUserSharing, _ := strconv.ParseInt(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MaxLimitUserSharing, 10, 64)
```
And similarly in:
```go
maxLimitUserSharing, _ := strconv.ParseInt(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MaxLimitUserSharing, 10, 64)
```

## Fix

Check and handle the error case appropriately and propagate a diagnostic error if parsing fails. For example:

```go
maxLimitUserSharing, err := strconv.ParseInt(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MaxLimitUserSharing, 10, 64)
if err != nil {
    resp.Diagnostics.AddError("Error parsing MaxLimitUserSharing as integer", err.Error())
    return
}
```
Add this pattern in Create, Update, and Read methods where parsing is performed.


## ISSUE 8

# Title

Assumption about string length and format when changing case for LimitSharingMode and SolutionCheckerMode

##

/workspaces/terraform-provider-power-platform/internal/services/managed_environment/resource_managed_environment.go

## Problem

When setting and reading the state for `LimitSharingMode` and `SolutionCheckerMode`, the code assumes the string is non-empty and at least 1 character long for expressions like:
```go
strings.ToLower(plan.LimitSharingMode.ValueString()[:1]) + plan.LimitSharingMode.ValueString()[1:]
```
If the value is empty or too short (due to user input, changes in validator definitions, or API evolution), this will cause a runtime panic (slice out of range) or produce an invalid string.

## Impact

High. A panic here during plan/apply will break the provider and result in user-facing errors and failed deployments. Defensive coding is required to ensure non-empty, valid strings before indexing.

## Location

Within set and read logic for LimitSharingMode and SolutionCheckerMode in Create, Update, and Read:

## Code Issue

```go
LimitSharingMode: strings.ToLower(plan.LimitSharingMode.ValueString()[:1]) + plan.LimitSharingMode.ValueString()[1:],
SolutionCheckerMode: strings.ToLower(plan.SolutionCheckerMode.ValueString()),
// ... when reconstructing state (uses [0:1] style substring without prior length check)
```

## Fix

Add checks to ensure safety:

```go
limitSharing := plan.LimitSharingMode.ValueString()
if len(limitSharing) > 0 {
    LimitSharingMode: strings.ToLower(limitSharing[:1]) + limitSharing[1:]
} else {
    // Handle empty (error, default, or skip)
}
```
Apply similar pattern where any substringing or indexing is performed.


## ISSUE 9

# Lack of Error Checking and Defensive Logic in getSolutionId

##
/workspaces/terraform-provider-power-platform/internal/services/solution/resource_solution.go

## Problem
The `getSolutionId` function splits the `id` string using `_` and returns the last element, assuming the split array always has elements. There is no check for an empty string or malformed `id` values (for example, if the input is an empty string or does not contain any `_`). If the string does not contain an underscore, then `split[len(split)-1]` will still work, but if the string is empty, `split[0]` will be the empty string, which may not be desirable. A stronger check would provide clearer code intention and better defensive coding, protecting against malformed IDs.

## Impact
- **Severity:** Medium
- Could lead to silent bugs or ingesting/processing invalid IDs which may cause downstream errors or logic issues, especially if the code evolves.
- Not strict about data consistency/parsing at code boundary edges.

## Location
End of the file:

## Code Issue
```go
func getSolutionId(id string) string {
	split := strings.Split(id, "_")
	return split[len(split)-1]
}
```

## Fix
Consider error handling for empty strings and possibly returning an error in those (rare) cases, or at least logging an unexpected format:

```go
func getSolutionId(id string) string {
    if id == "" {
        // Optionally log or handle error
        return ""
    }
    split := strings.Split(id, "_")
    return split[len(split)-1]
}
```

Or, if desired, a more defensive pattern:

```go
func getSolutionId(id string) string {
    if id == "" {
        return ""
    }
    split := strings.Split(id, "_")
    if len(split) == 0 {
        return ""
    }
    return split[len(split)-1]
}
```


## ISSUE 10

# Insufficient Validation of Required Attributes

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy.go

## Problem

While the schema specifies required attributes (`is_disabled`, `allowed_tenants`, etc.), there is no explicit logic in the resource code ensuring these attributes are not empty within Create/Update paths beyond what the framework validator does. There is only validation that at least one of `inbound` or `outbound` is set per tenant, but no guard that `allowed_tenants` is not an empty set (which could represent an always-blocking configuration that may be unintentional).

## Impact

- **Severity: Medium**
- Could allow a situation where a user accidentally submits an empty `allowed_tenants` set (or doesn't provide required directions) and the API call goes through, possibly breaking tenant connectivity in production. Behavioral differences between the Terraform provider's required-vs-empty semantics and actual API expectations can cause state drift or unintended outages.

## Location

```go
// Schema only has Required: true, but in the ValidateConfig and CRUD, no explicit check for allowed_tenants not being empty
```

## Code Issue

```go
"allowed_tenants": schema.SetNestedAttribute{
	Required:            true,
	MarkdownDescription: "List of tenants that are allowed to connect with your tenant.",
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the tenant that is allowed to connect.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"inbound": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether inbound connections from this tenant are allowed.",
			},
			"outbound": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether outbound connections to this tenant are allowed.",
			},
		},
	},
},
// ...
// In ValidateConfig()
if resp.Diagnostics.HasError() {
	return
}
// No check: if len(modelTenants) == 0 { /* error */ }
```

## Fix

Add explicit validation to ensure that the list/set of allowed tenants is not empty, both at the schema validator level and in ValidateConfig, so intent is clear and errors are user-friendly.

```go
// Add validator in schema if supported, otherwise in ValidateConfig:
if len(modelTenants) == 0 {
	resp.Diagnostics.AddAttributeError(
		path.Root("allowed_tenants"),
		"Empty allowed_tenants set",
		"'allowed_tenants' must not be empty. At least one outbound or inbound allowed tenant must be specified.",
	)
	return
}
```


## ISSUE 11

# Title

Error not returned or propagated in environment security role validation

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user.go

## Problem

In various method contexts (`Create`, `Update`, etc.), when calling `validateEnvironmentSecurityRoles`, the error is added using `resp.Diagnostics.AddError` but code execution continues instead of returning immediately upon error detection.

## Impact

This causes functions to possibly proceed with an invalid state, which could result in further unexpected errors downstream or cause partial/misleading updates to the resource. Errors should immediately halt further processing in these lifecycle operations. Severity: **High**.

## Location

For example, in `Create`, inside the else block:

```go
err := validateEnvironmentSecurityRoles(plan.SecurityRoles)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err.Error())
}
// execution continues even when error is detected!
```

## Code Issue

```go
err := validateEnvironmentSecurityRoles(plan.SecurityRoles)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err.Error())
}
// ... code continues
```

## Fix

Return immediately after adding an error to diagnostics to prevent further execution after a fatal validation:

```go
err := validateEnvironmentSecurityRoles(plan.SecurityRoles)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err.Error())
    return
}
```

_This applies to similar usages in `Update` and other places where validation errors are encountered._

---

This issue should be saved in:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/resource_user.go-validation_immediate_return-high.md.


## ISSUE 12

# Title

Ambiguous return value handling for validation and CRUD path

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user.go

## Problem

In "environment" user mode in both `Create` and `Update`, the code validates the list of security roles via `validateEnvironmentSecurityRoles`, but the validation is potentially run multiple times (once per invocation of the function, e.g., both in `Create` and again in `Update`). There's logic drift on how to ensure only valid roles are sent to API or in plan/state. Additionally, after some validation failures, code paths for CRUD operations still attempt to marshal or update user data, which should not occur for invalid input.

## Impact

This introduces potential for inconsistent runtime/resource states, repeated API calls for known-invalid operations, and blurs separation of validation versus mutation. It also increases cognitive load for maintainers. Severity: **Medium**.

## Location

Example from `Update`:

```go
err := validateEnvironmentSecurityRoles(plan.SecurityRoles)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
}
if len(addedSecurityRoles) > 0 {
    userDto, err := r.UserClient.AddEnvironmentUserSecurityRoles(ctx, plan.EnvironmentId.ValueString(), plan.AadId.ValueString(), addedSecurityRoles)
    ...
}
```

No explicit short-circuit/return after error, and logic continues as if validation passed.

## Code Issue

```go
err := validateEnvironmentSecurityRoles(plan.SecurityRoles)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
}
// further update/CRUD logic follows, potentially using invalid state
```

## Fix

Always return immediately after a validation error to prevent invalid state propagating to CRUD:

```go
err := validateEnvironmentSecurityRoles(plan.SecurityRoles)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
    return
}
// safe to continue with CRUD logic here
```

---

This issue should be saved in:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/resource_user.go-validation_vs_crud-medium.md.


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
