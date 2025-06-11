# General Nil Handling Issues

This document contains all identified nil handling issues related to general components in the Terraform Power Platform provider codebase.

## ISSUE 1

<!-- Source: contexts.go-missed_cancel_func_release-medium.md -->

# Issue: Control Flow—Missed CancelFunc Release on Early Return

##

/workspaces/terraform-provider-power-platform/internal/helpers/contexts.go

## Problem

The returned closure function from `EnterRequestContext` always calls `(*cancel)()` if `cancel` is not nil. However, there are code paths within `enterTimeoutContext` where `cancel` will be `nil` (for example, if an error occurs). While `(*cancel)()` is only called if `cancel != nil`, this places a subtle responsibility on the caller to always check for nil, which is risky in future refactoring.

Additionally, it's idiomatically clearer in Go to return a no-op function if cancellation isn't needed, and to never return a potentially nil function pointer.

## Impact

- **Impact:** Medium  
  If someone refactors or copies this idiom incorrectly, it could lead to code that unintentionally dereferences nil, or misses the opportunity to clean up resources correctly.

## Location

Closure function in `EnterRequestContext`'s return statement:

## Code Issue

```go
return ctx, func() {
 tflog.Debug(ctx, fmt.Sprintf("%s END: %s", reqType, name))
 if cancel != nil {
  (*cancel)()
 }
}
```

## Fix

Change `enterTimeoutContext` to return a no-op cancel function if there was no cancellation set, and always call it.

```go
func enterTimeoutContext[T AllowedRequestTypes](ctx context.Context, req T) (context.Context, context.CancelFunc) {
 // instead of *context.CancelFunc return type, use just context.CancelFunc and return context.CancelFunc(func(){}) when nil
}
```

Change `EnterRequestContext` closure to always call the returned cancel function.

```go
ctx, cancel := enterTimeoutContext(ctx, req)

return ctx, func() {
 tflog.Debug(ctx, fmt.Sprintf("%s END: %s", reqType, name))
 cancel()
}
```

## ISSUE 2

<!-- Source: other_field_required_when_value_of_validator_structure_low.md -->

# Unnecessary use of pointer for string type

##

/workspaces/terraform-provider-power-platform/internal/validators/other_field_required_when_value_of_validator.go

## Problem

The `otherFieldValue` variable is declared as `new(string)`, resulting in a pointer to a string, but there's no functional need for a pointer since the value is being set directly. This introduces extra indirection that is not idiomatic in Go for basic types unless mutability (across function calls) or explicit nil/no value distinctions are needed. The `GetAttribute` call is also inconsistent in how receiver variables are being used (`currentFieldValue` is a value; `otherFieldValue` is a pointer).

## Impact

Unnecessary use of pointers can make the code harder to read and maintain, and introduces potential for subtle bugs. Severity: **low**.

## Location

```go
otherFieldValue := new(string)
d := req.Config.GetAttribute(ctx, paths[0], &otherFieldValue)
```

## Code Issue

```go
 otherFieldValue := new(string)
 d := req.Config.GetAttribute(ctx, paths[0], &otherFieldValue)
```

## Fix

Just declare as a value variable and pass its address:

```go
 var otherFieldValue string
 d := req.Config.GetAttribute(ctx, paths[0], &otherFieldValue)
```

If you need to check for empty string, just use `otherFieldValue == ""`.

## ISSUE 3

<!-- Source: request_body_bytes_assign_on_error_medium.md -->

# Body as bytes is set even if `io.ReadAll` fails

##

/workspaces/terraform-provider-power-platform/internal/api/request.go

## Problem

The code reads body without handling `err` before assigning its value. If `err` is not nil, `body` will be meaningless or incomplete, but it will still be assigned to `resp.BodyAsBytes`. Caller might receive an invalid or partial body.

## Impact

Severity: Medium

This could lead to attempts to unmarshal invalid/incomplete data or propagate partial/corrupted information.

## Location

```go
 body, err := io.ReadAll(apiResponse.Body)
 resp.BodyAsBytes = body
```

## Code Issue

```go
 body, err := io.ReadAll(apiResponse.Body)
 resp.BodyAsBytes = body
```

## Fix

Check for `err` from `io.ReadAll` before setting `resp.BodyAsBytes`:

```go
 body, err := io.ReadAll(apiResponse.Body)
 if err != nil {
  return resp, err
 }
 resp.BodyAsBytes = body
```

## ISSUE 4

<!-- Source: request_error_return_nil_http_response_medium_high.md -->

# Error return combined with possibly nil http.Response

##

/workspaces/terraform-provider-power-platform/internal/api/request.go

## Problem

In `doRequest`, if an error is returned by `httpClient.Do`, a `Response` containing a possibly nil `apiResponse` is returned. Later code (e.g., the caller) could try to access fields on `resp.HttpResponse` without checking for nil, causing a panic.

## Impact

Severity: Medium/High

This could result in runtime panics if not handled, reducing the stability of the codebase.

## Location

```go
 apiResponse, err := httpClient.Do(request)
 resp := &Response{
  HttpResponse: apiResponse,
 }

 if err != nil {
  return resp, err
 }

 if apiResponse == nil {
  return resp, errors.New("unexpected nil response without error")
 }
```

## Code Issue

```go
 resp := &Response{
  HttpResponse: apiResponse,
 }

 if err != nil {
  return resp, err
 }
```

## Fix

Only return a non-nil Response if apiResponse is non-nil, otherwise return nil for Response:

```go
 apiResponse, err := httpClient.Do(request)

 if err != nil {
  return nil, err
 }
 if apiResponse == nil {
  return nil, errors.New("unexpected nil response without error")
 }
 resp := &Response{
  HttpResponse: apiResponse,
 }
```

## ISSUE 5

<!-- Source: request_getheader_nil_dereference_high.md -->

# Possible nil dereference in Response.GetHeader

##

/workspaces/terraform-provider-power-platform/internal/api/request.go

## Problem

If `apiResponse.HttpResponse` is nil (e.g., on error responses), calling methods on it will panic.

## Impact

Severity: High

This can easily result in runtime panics if callers do not check for nil Response before calling `GetHeader`.

## Location

```go
func (apiResponse *Response) GetHeader(name string) string {
 return apiResponse.HttpResponse.Header.Get(name)
}
```

## Code Issue

```go
 return apiResponse.HttpResponse.Header.Get(name)
```

## Fix

Check for nil `HttpResponse` before accessing its fields:

```go
func (apiResponse *Response) GetHeader(name string) string {
 if apiResponse.HttpResponse == nil {
  return ""
 }
 return apiResponse.HttpResponse.Header.Get(name)
}
```

## ISSUE 6

<!-- Source: requires_replace_string_from_non_empty_modifier_go_lack_of_error_handling_medium.md -->

# Issue: Lack of Error Handling in PlanModifyString Method

##

/workspaces/terraform-provider-power-platform/internal/modifiers/requires_replace_string_from_non_empty_modifier.go

## Problem

The `PlanModifyString` method in the `requireReplaceStringFromNonEmptyPlanModifier` struct does not handle or report any errors. While the logic checks certain conditions, it is possible that calls like `IsNull()`, `IsUnknown()`, or especially `ValueString()` could, depending on their implementation, encounter an error or unexpected state (such as operating on a nil or malformed value). The method signature does not allow for passing errors or diagnostics to the response, which is typically needed in terraform planmodifier pattern.

## Impact

**Severity: Medium**

If an error occurs within the plan modification logic (for example, if the value accessors panic or return an invalid value), the lack of error handling will result in silent failures, possible panics, or, more insidiously, incorrect plan behavior that can go undetected.

## Location

```go
func (d *requireReplaceStringFromNonEmptyPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
 if req.PlanValue != req.StateValue && (!req.StateValue.IsNull() && !req.StateValue.IsUnknown() && req.StateValue.ValueString() != "") {
  resp.RequiresReplace = true
 }
}
```

## Code Issue

```go
func (d *requireReplaceStringFromNonEmptyPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
 if req.PlanValue != req.StateValue && (!req.StateValue.IsNull() && !req.StateValue.IsUnknown() && req.StateValue.ValueString() != "") {
  resp.RequiresReplace = true
 }
}
```

## Fix

Add error handling by checking if the `ValueString()` call (and other accessors, if needed) provides a way to detect errors. Typically in Terraform plugin framework, you should add errors to the diagnostics in the response if something goes wrong. For example:

```go
func (d *requireReplaceStringFromNonEmptyPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
 valueStr, err := req.StateValue.ToStringValue()
 if err != nil {
  resp.Diagnostics.AddError(
   "Unable to convert state value to string",
   err.Error(),
  )
  return
 }
 if req.PlanValue != req.StateValue && (!req.StateValue.IsNull() && !req.StateValue.IsUnknown() && valueStr != "") {
  resp.RequiresReplace = true
 }
}
```

## ISSUE 7

<!-- Source: uuid_value_valueuuid_pointer_medium.md -->

# Inconsistent Value and Diagnostic Return in ValueUUID

##

/workspaces/terraform-provider-power-platform/internal/customtypes/uuid_value.go

## Problem

In the `ValueUUID` method, the code returns an empty `UUIDValue` struct and diagnostics when the value is null, unknown, or invalid. However, returning an empty struct as the value might mislead the consumer into thinking this is a valid, but "zero" UUID. It is more idiomatic in Go to return a pointer and return `nil` instead, or to document that the returned struct must not be used if diagnostics are non-empty.

## Impact

**Medium Severity**: This may lead to code that trusts the returned value even when errors are present, resulting in subtle bugs or panics in downstream operations.

## Location

Method `ValueUUID`

## Code Issue

```go
func (v UUIDValue) ValueUUID() (UUIDValue, diag.Diagnostics) {
 var diags diag.Diagnostics

 if v.IsNull() {
  diags.Append(diag.NewErrorDiagnostic(UUIDTypeErrorInvalidStringHeader, "UUID string value is null"))

  return UUIDValue{}, diags
 }

 if v.IsUnknown() {
  diags.Append(diag.NewErrorDiagnostic(UUIDTypeErrorInvalidStringHeader, "UUID string value is unknown"))

  return UUIDValue{}, diags
 }

 _, err := uuid.ParseUUID(v.ValueString())
 if err != nil {
  diags.Append(diag.NewErrorDiagnostic(
   UUIDTypeErrorInvalidStringHeader,
   fmt.Sprintf(UUIDTypeErrorInvalidStringDetails, v.ValueString()),
  ))

  return UUIDValue{}, diags
 }

 return v, nil
}
```

## Fix

Consider returning a pointer (e.g., `*UUIDValue`) so that errors are clearly reflected by returning `nil`, or update documentation and usages to ensure downstream code validates diagnostics before using the value. Example using a pointer:

```go
func (v UUIDValue) ValueUUID() (*UUIDValue, diag.Diagnostics) {
 var diags diag.Diagnostics

 if v.IsNull() {
  diags.Append(diag.NewErrorDiagnostic(UUIDTypeErrorInvalidStringHeader, "UUID string value is null"))
  return nil, diags
 }

 if v.IsUnknown() {
  diags.Append(diag.NewErrorDiagnostic(UUIDTypeErrorInvalidStringHeader, "UUID string value is unknown"))
  return nil, diags
 }

 _, err := uuid.ParseUUID(v.ValueString())
 if err != nil {
  diags.Append(diag.NewErrorDiagnostic(
   UUIDTypeErrorInvalidStringHeader,
   fmt.Sprintf(UUIDTypeErrorInvalidStringDetails, v.ValueString()),
  ))
  return nil, diags
 }

 return &v, nil
}
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

Apply this fix to the whole codebase
