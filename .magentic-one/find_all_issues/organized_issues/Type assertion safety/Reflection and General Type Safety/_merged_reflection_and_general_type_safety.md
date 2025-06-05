# Reflection and General Type Safety Issues

This document contains type assertion safety issues related to reflection usage, general type safety violations, nil checks, and runtime type validation in the codebase.

## ISSUE 1

# Title

Improper use of reflection to check for nil interface

##

/workspaces/terraform-provider-power-platform/internal/api/client.go

## Problem

The code in the `prepareRequestBody` function uses `reflect.ValueOf(body).Kind()` and `reflect.ValueOf(body).IsNil()` to determine if `body` is nil. This is inefficient and error-prone because it uses reflection instead of Go's native type checking. It can panic if `ValueOf` is called on a non-pointer or a non-nil value for certain types. Go's type assertion and interface nil checks are preferred.

## Impact

Improper reflection use may cause panics at runtime and makes the code less readable. Severity: **high**

## Location

`prepareRequestBody` function

## Code Issue

```go
if body != nil && (reflect.ValueOf(body).Kind() != reflect.Ptr || !reflect.ValueOf(body).IsNil()) {
 if strp, ok := body.(*string); ok {
  bodyBuffer = strings.NewReader(*strp)
 } else {
  bodyBytes, err := json.Marshal(body)
  if err != nil {
   return nil, err
  }
  bodyBuffer = bytes.NewBuffer(bodyBytes)
 }
}
```

## Fix

Simplify to check if the interface is nil by a safer pattern for optional pointers, and check for `*string` value directly:

```go
if body != nil {
 if strp, ok := body.(*string); ok && strp != nil {
  bodyBuffer = strings.NewReader(*strp)
 } else {
  bodyBytes, err := json.Marshal(body)
  if err != nil {
   return nil, err
  }
  bodyBuffer = bytes.NewBuffer(bodyBytes)
 }
}
```

---

## ISSUE 2

# Title

Missing Type Validations and Unchecked Nil Returns in Conversion Logic

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy.go

## Problem

Throughout the resource lifecycle methods (Read, Create, Update), model conversion functions such as `convertToAttrValueConnectorsGroup`, `convertToAttrValueCustomConnectorUrlPatternsDefinition`, and corresponding setters/readers are invoked. There is no type assertion or nil-check logic to handle conversion failures, bad data, or type mismatches. If any conversion returns nil or has a type mismatch, this may cause panics or incorrect assignment to state/plan, and may leave the Terraform state in an inconsistent or corrupted state.

## Impact

Severity: Critical

Unvalidated conversions, unchecked nil errors, or assignment of nil to required schema attributes pose a risk of provider panics and corruption of the user's infrastructure state. This is a critical risk especially during upgrades or with malformed upstream API responses.

## Location

```go
state.BusinessGeneralConnectors = convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
state.NonBusinessConfidentialConnectors = convertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
state.BlockedConnectors = convertToAttrValueConnectorsGroup("Blocked", policy.ConnectorGroups)
...
policyToCreate.Environments = convertToDlpEnvironment(ctx, plan.Environments)
policyToCreate.CustomConnectorUrlPatternsDefinition = convertToDlpCustomConnectorUrlPatternsDefinition(ctx, resp.Diagnostics, plan.CustomConnectorsPatterns)
```

## Code Issue

```go
state.BusinessGeneralConnectors = convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
...
policyToCreate.CustomConnectorUrlPatternsDefinition = convertToDlpCustomConnectorUrlPatternsDefinition(ctx, resp.Diagnostics, plan.CustomConnectorsPatterns)
```

## Fix

- Add explicit type checks and nil handling for all conversion functions.  
- If a conversion returns nil or fails, append a diagnostic error rather than assigning nil to required fields.
- Ensure all assignments to required attributes are validated before state updates.

```go
result := convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
if result == nil {
 resp.Diagnostics.AddError("Connector Group Conversion Failed", "Failed to convert 'Confidential' connector group to attribute value.")
 return
}
state.BusinessConnectors = result
```

---
Save as:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/type_safety/resource_dlp_policy_missing_nil_checks_critical.md

---

## ISSUE 3

# Untyped Error Code Check in Error Handling

##

/workspaces/terraform-provider-power-platform/internal/services/environment_wave/resource_environment_wave.go

## Problem

The code checks if an error corresponds to `customerrors.ERROR_OBJECT_NOT_FOUND` by using the `customerrors.Code()` function, but does not use Go 1.13+ idiomatic error wrapping and type assertions for robust error handling. Using error codes as strings or constants can introduce bugs and makes handling error types less safe and readable.

## Impact

This can lead to fragile error handling logic if the error string or code changes. It is less type-safe and can hide error context, making troubleshooting more difficult. Severity: **critical** because it can cause logic to silently misbehave if the underlying error value changes or is wrapped.

## Location

```go
 if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
  resp.State.RemoveResource(ctx)
  return
 }
```

## Code Issue

```go
if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
 resp.State.RemoveResource(ctx)
 return
}
```

## Fix

Use Go's standard library error wrapping and unwrapping with errors.Is or custom sentinel error variables for robust error comparison:

```go
import "errors"

if errors.Is(err, customerrors.ErrObjectNotFound) {
 resp.State.RemoveResource(ctx)
 return
}
```

**Explanation:**

- This approach uses type-safe error handling and supports Go's error-wrapping best practices (introduced in Go 1.13).
- Adjust `customerrors` to export a proper error variable if not already present (e.g., `var ErrObjectNotFound = errors.New("object not found")`).
- Benefits include maintainability, improved debugging, and correctness.

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
