# Plan and State Nil Pointer Handling Issues

This document consolidates all identified nil pointer dereference issues related to plan and state handling in the Terraform Provider for Power Platform.

## ISSUE 1

**Title**: Insufficient validation for required attributes and potential for nil pointer dereference

**File**: `/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy_environment.go`

**Problem**: While the Terraform schema enforces that `"billing_policy_id"` and `"environments"` are required attributes, the resource logic does not appear to validate that these values are non-empty before using them. If for any reason (migration, schema change, state corruption, or upstream bug) `plan.BillingPolicyId` or `plan.Environments` is empty or nil, downstream functions (such as `GetEnvironmentsForBillingPolicy`, `RemoveEnvironmentsToBillingPolicy`, `AddEnvironmentsToBillingPolicy`) may receive empty or malformed values, leading to unclear errors, API rejections, or even nil pointer dereference panics.

**Impact**: Severity: high

A nil pointer dereference can cause a panic, which will crash the Terraform provider plugin and stop the user's operation abruptly, causing user distrust and potential data loss. Passing invalid input downstream without validation can lead to late, unclear errors or inconsistent state.

**Code Issue**:

```go
var plan *BillingPolicyEnvironmentResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
 return
}
// plan could be nil here if deserialization failed, or attributes empty!
environments, err := r.LicensingClient.GetEnvironmentsForBillingPolicy(ctx, plan.BillingPolicyId)
```

**Fix**: Add robust nil checks and field value validation after loading state or plan. If required values are missing, add descriptive diagnostic errors before attempting any API call or downstream logic.

```go
// Example logic after loading from state/plan
if plan == nil {
 resp.Diagnostics.AddError("Invalid plan", "Plan could not be loaded; aborting resource operation.")
 return
}
if plan.BillingPolicyId == "" {
 resp.Diagnostics.AddError("Missing required attribute", "\"billing_policy_id\" cannot be empty.")
 return
}
if len(plan.Environments) == 0 {
 resp.Diagnostics.AddError("Missing required attribute", "\"environments\" cannot be empty.")
 return
}
```

## ISSUE 2

**Title**: Error handling for plan/state Get does not differentiate between error types

**File**: `/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection_share.go`

**Problem**: The calls to `req.Plan.Get` and `req.State.Get` in resource lifecycle methods append diagnostics but do not differentiate between errors caused by user input, type incompatibility, or partial failures. There's no separate logging or error return except for a generic check with `HasError()`. Additionally, nil pointer dereferences could occur if a `nil` value slips into the plan or state, since they are immediately dereferenced later.

**Impact**: Severity is **medium**. While most standard errors will be caught by the diagnostics check, more granular error handling or protective checks for nil plan/state would make error handling more robust.

**Code Issue**:

```go
 var plan *ShareResourceModel

 resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

 if resp.Diagnostics.HasError() {
  return
 }
```

**Fix**: Check for a nil value on state or plan before using it, in addition to diagnostics:

```go
 var plan *ShareResourceModel
 resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
 if resp.Diagnostics.HasError() || plan == nil {
  return
 }
```

## ISSUE 3

**Title**: No Defensive Checks for nil plan or state in CRUD Methods

**File**: `/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy.go`

**Problem**: In the CRUD methods (`Read`, `Create`, `Update`, `Delete`), the code directly assigns into fields of potentially nil pointers (`plan`, `state`). If the value from `req.Plan.Get` or `req.State.Get` is nil due to misconfiguration or changes in Terraform core, a nil pointer dereference panic will occur.

**Impact**: Severity: Critical

Critical stability and reliability bug. Any runtime change from Terraform, or upstream changes in schema/parser behavior, could cause a panic and crash the provider, leading to loss of in-flight state, partial infrastructure, and requiring manual operator intervention.

**Code Issue**:

```go
var plan *dataLossPreventionPolicyResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
 return
}
// plan may still be nil here!
plan.Id = types.StringValue(policy.Name)
```

**Fix**: Add a nil check after unmarshal and before dereferencing/assigning to fields.

```go
var plan *dataLossPreventionPolicyResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
 return
}
if plan == nil {
 resp.Diagnostics.AddError("Internal Provider Error", "Plan was nil after reading configuration.")
 return
}
```

## ISSUE 4

**Title**: Struct Initialization Issue in Resource Methods

**File**: `/workspaces/terraform-provider-power-platform/internal/services/application/resource_environment_application_package_install.go`

**Problem**: In the `Read`, `Update`, and `Delete` methods, the code retrieves resource state into a pointer to a struct (`*EnvironmentApplicationPackageInstallResourceModel`). If the state is empty or not present, this will result in a nil pointer, leading to a potential `nil` dereference and runtime panic if any field is accessed.

**Impact**: Severity: **High**  
This can lead to runtime panics which can stop the provider execution abruptly and cause irrecoverable state within a Terraform operation.

**Code Issue**:

```go
var state *EnvironmentApplicationPackageInstallResourceModel

resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

if resp.Diagnostics.HasError() {
    return
}

tflog.Debug(ctx, fmt.Sprintf("READ: %s with application_name %s", r.FullTypeName(), state.UniqueName.ValueString()))

resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
```

**Fix**: Initialize the struct as a non-pointer to avoid nil pointer panics when reading from the state, and use the address of the struct when passing to the `Get` and `Set` methods:

```go
var state EnvironmentApplicationPackageInstallResourceModel

resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

if resp.Diagnostics.HasError() {
    return
}

tflog.Debug(ctx, fmt.Sprintf("READ: %s with application_name %s", r.FullTypeName(), state.UniqueName.ValueString()))

resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
```

## ISSUE 5

**Title**: Multiple Use of Pointers for ResourceModel in CRUD Functions Reduces Consistency and Clarity

**File**: `/workspaces/terraform-provider-power-platform/internal/services/environment_groups/resource_environment_group.go`

**Problem**: In your CRUD operations, you sometimes use a value (`EnvironmentGroupResourceModel{}` in `Read`), and sometimes a pointer (`var plan *EnvironmentGroupResourceModel` in `Create`, `Update`). This inconsistency can cause unexpected nil pointer dereference panics or require unnecessary allocation and pointer indirection.

**Impact**:

- Potential for panic if the pointer is not properly set by the framework (medium impact).
- Reduces code readability and consistency.
- May complicate mocking and testing.

**Severity:** medium

**Code Issue**:
Current pattern in `Create`:

```go
var plan *EnvironmentGroupResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
    return
}

// usage
environmentGroupToCreate := environmentGroupDto{
    DisplayName: plan.DisplayName.ValueString(),
    Description: plan.Description.ValueString(),
}
```

Pattern in `Read`:

```go
state := EnvironmentGroupResourceModel{}
resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
```

**Fix**: Prefer consistent use—if the struct is not very large and ownership does not need to be transferred, use value types:

```go
plan := EnvironmentGroupResourceModel{}
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
    return
}
...
// use plan.X not plan->X
```

## ISSUE 6

**Title**: Use of pointer-to-struct for resource plan/state model may risk nil dereference

**File**: `/workspaces/terraform-provider-power-platform/internal/services/managed_environment/resource_managed_environment.go`

**Problem**: The code repeatedly uses a pointer (`var plan *ManagedEnvironmentResourceModel`) for the resource plan/state model throughout Create, Update, and Read. In most cases, this is populated by `.Get()` methods, but any future code changes, test scaffolding, or framework changes could result in nil being assigned or improperly mocked, yielding nil pointer panics on field access. Best practice is to prefer value structs unless mutation of the struct pointer itself is needed (rare in Terraform provider logic) or nil is a valid state. Type safety is improved by avoiding unnecessary indirection.

**Impact**: Medium. Could cause panics or hidden bugs, especially in future test or code refactoring where plan is accidentally left nil after .Get or during mock initialization.

**Code Issue**:

```go
var plan *ManagedEnvironmentResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
// use plan immediately: plan.EnvironmentId.ValueString()
```

**Fix**: Declare the struct as a value, not a pointer:

```go
var plan ManagedEnvironmentResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
// use plan.EnvironmentId.ValueString(), etc.
```

## ISSUE 7

**Title**: Use of Pointers for ResourceModel in Plan and State May Cause Nil Dereferences

**File**: `/workspaces/terraform-provider-power-platform/internal/services/solution/resource_solution.go`

**Problem**: Throughout the code (e.g., in Create, Update, Read), the `plan` and `state` variables are declared as pointers (`*ResourceModel`), and populated via `req.Plan.Get()` or `req.State.Get()`, but code later dereferences fields without checking if the pointer is non-nil. If, for any reason, the plan or state is not properly assigned (e.g., decoding failure or data bug), this will cause a runtime panic with nil pointer dereference.

**Impact**:

- **Severity:** High
- Potential for runtime panics if Terraform sends an invalid or unexpected state or plan, leading to provider crashes.
- Makes the code unsafe to refactor or test.
- Reduces type safety and defensiveness, especially in a plugin context.

**Code Issue**:

```go
var plan *ResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

// ...

if !plan.SettingsFile.IsNull() && !plan.SettingsFile.IsUnknown() {
    value, err := helpers.CalculateSHA256(plan.SettingsFile.ValueString())
// ...
```

**Fix**: Immediately after extracting the plan or state, validate that it is non-nil before use, and emit a diagnostic (or handle gracefully) if it is nil:

```go
if plan == nil {
    resp.Diagnostics.AddError("Invalid plan received", "Resource plan is nil after decoding. This is likely an internal bug or provider incompatibility.")
    return
}

// Similar for 'state' as used in other functions
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
