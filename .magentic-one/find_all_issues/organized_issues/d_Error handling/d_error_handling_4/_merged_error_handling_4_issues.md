# Error Handling Issues: Critical State Management and Retry Loop Failures

This document consolidates critical error handling issues including state management corruption, retry loop exhaustion failures, validator error handling, and diagnostic mishandling that can lead to silent failures and data corruption in Terraform operations.

---

## ISSUE 1

**File Path:** `/workspaces/terraform-provider-power-platform/internal/services/locations/api_locations.go`

**Title:** Does not provide context when returning errors from `GetLocations`

**Problem:**
The `GetLocations` method returns errors from the API client directly, without wrapping or enriching them with context about the operation that failed.

**Impact:**
Medium severity. Debugging can be more difficult if consumers of this function cannot disambiguate where an error occurred.

**Code Issue:**

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &locations)
return locations, err
```

**Fix:**
Wrap the error to provide more context.

```go
_, err := client.API.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &locations)
if err != nil {
 return locations, fmt.Errorf("failed to get locations: %w", err)
}
return locations, nil
```

---

## ISSUE 2

**File Path:** `/workspaces/terraform-provider-power-platform/internal/services/authorization/api_user.go`

**Title:** Missing Error Handling in Retry Loop for CreateDataverseUser

**Problem:**
In the `CreateDataverseUser` function, the code is retrying API calls up to a certain count, attempting to handle a specific error ("userNotLicensed") via string matching. However, there is no upper-bound or critical alert when all retry attempts are exhausted. If the retries are exhausted (retryCount drops to 0) but the error persists, the final error is only returned, which may lose context regarding the exhaustive nature of the retries and lacks a clear error message for failed user creation after maximum attempts.

**Impact:**
Severity: Critical

This can lead to silent failures or confusing error experiences for the caller. Users may get a generic API error rather than understanding that all retry attempts have been exhausted. Operators/debuggers may have difficulty diagnosing race conditions, transient failures, or persistent issues with user creation and licensing propagation.

**Code Issue:**

```go
 // 9 minutes of retries.
 retryCount := 6 * 9
 var err error

 for retryCount > 0 {
  _, err = client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, userToCreate, []int{http.StatusOK}, nil)
  // the license assignment in Entra is async, so we need to wait for that to happen if a user is created in the same terraform run.
  if err == nil || !strings.Contains(err.Error(), "userNotLicensed") {
   break
  }
  tflog.Debug(ctx, fmt.Sprintf("Error creating user: %s", err.Error()))
  err = client.Api.SleepWithContext(ctx, 10*time.Second)
  if err != nil {
   return nil, err
  }

  retryCount--
 }
 if err != nil {
  return nil, err
 }
```

**Fix:**
Log an explicit error or wrap the error with additional context when retries are exhausted:

```go
// 9 minutes of retries.
retryCount := 6 * 9
var err error

for retryCount > 0 {
 _, err = client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, userToCreate, []int{http.StatusOK}, nil)
 if err == nil || !strings.Contains(err.Error(), "userNotLicensed") {
  break
 }
 tflog.Debug(ctx, fmt.Sprintf("Error creating user: %s", err.Error()))
 err = client.Api.SleepWithContext(ctx, 10*time.Second)
 if err != nil {
  return nil, err
 }

 retryCount--
}
if err != nil {
 if retryCount == 0 {
  return nil, fmt.Errorf("failed to create Dataverse user after maximum retries: %w", err)
 }
 return nil, err
}
```

---

## ISSUE 3

**File Path:** `/workspaces/terraform-provider-power-platform/internal/services/currencies/datasource_currencies.go`

**Title:** Error Handling is not Comprehensive in `Read` Method

**Problem:**
In the `Read` method, after calling `resp.State.Get(ctx, &state)`, the error returned from the Get operation is not checked. If an error occurs while retrieving the state, the code proceeds, potentially with a zero-value or invalid `state`, which could cause unexpected behaviors later in the function.

**Impact:**
This can lead to misleading diagnostics being returned to the user and may cause runtime panics, data mismatches, or further errors during the function execution. It also impacts debugging and maintainability as errors at this stage are silently ignored.

**Severity:** High

**Code Issue:**

```go
func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
 ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
 defer exitContext()

 var state DataSourceModel
 resp.State.Get(ctx, &state)

 currencies, err := d.CurrenciesClient.GetCurrenciesByLocation(ctx, state.Location.ValueString())
 if err != nil {
  resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
  return
 }
    ...
}
```

**Fix:**
Check and handle the error returned by `resp.State.Get`. If it has an error, append diagnostics and return immediately.

```go
 var state DataSourceModel
 diags := resp.State.Get(ctx, &state)
 resp.Diagnostics.Append(diags...)
 if resp.Diagnostics.HasError() {
  return
 }
```

---

## ISSUE 4

**File Path:** `/workspaces/terraform-provider-power-platform/internal/services/environment/datasource_environments.go`

**Title:** Error Handling Incomplete for Unknown Error Types

**Problem:**
In the `Read` function, the error handling for `d.EnvironmentClient.GetDefaultCurrencyForEnvironment` only adds a warning if the error code matches a known value. However, other unexpected error types might occur, and these are currently ignored, potentially leading to silent failures.

**Impact:**
Unexpected errors that are not warnings or the known error (`ERROR_ENVIRONMENT_URL_NOT_FOUND`) are ignored, resulting in suppressed diagnostics and more difficult troubleshooting for users. Severity: **Medium**.

**Code Issue:**

```go
defaultCurrency, err := d.EnvironmentClient.GetDefaultCurrencyForEnvironment(ctx, env.Name)
if err != nil {
    if customerrors.Code(err) != customerrors.ERROR_ENVIRONMENT_URL_NOT_FOUND {
        resp.Diagnostics.AddWarning(fmt.Sprintf("Error when reading default currency for environment %s", env.Name), err.Error())
    }
} else {
    currencyCode = defaultCurrency.IsoCurrencyCode
}
```

**Fix:**
Add an explicit branch to handle truly unexpected errors, perhaps with a proper error diagnostic instead of a warning:

```go
if err != nil {
    switch customerrors.Code(err) {
    case customerrors.ERROR_ENVIRONMENT_URL_NOT_FOUND:
        // Non-critical, just skip currency.
    default:
        resp.Diagnostics.AddError(
            fmt.Sprintf("Unexpected error when reading default currency for environment %s", env.Name),
            err.Error(),
        )
        return
    }
} else {
    currencyCode = defaultCurrency.IsoCurrencyCode
}
```

---

## ISSUE 5

**File Path:** `/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go`

**Title:** Diagnostics mishandling on Delete (cannot append Diagnostics to Diagnostics object)

**Problem:**
In the `Delete` method, the following code is used:

```go
previousBytes, diag := req.Private.GetKey(ctx, "original_settings")
if diag.HasError() {
    diag.Append(diag...)
    return
}
```

`diag` is a `tfsdk.Diagnostics` object. You are attempting to append diagnostics to itself (`diag.Append(diag...)`). However, `Append` is used to append a `Diagnostics` object to a different `Diagnostics` object. This does not write to `resp.Diagnostics` (the only diagnostics output that is considered by Terraform). As such, errors captured in `diag` are not presented to the end user in the Terraform output and may go unnoticed.

**Impact:**
Improper error handling and diagnostic output for users and maintainers. Errors may go unreported or make debugging more difficult. Severity: medium.

**Code Issue:**

```go
previousBytes, diag := req.Private.GetKey(ctx, "original_settings")
if diag.HasError() {
    diag.Append(diag...)
    return
}
```

**Fix:**
Append `diag` to `resp.Diagnostics`, not to itself. Only `resp.Diagnostics` is output:

```go
previousBytes, diag := req.Private.GetKey(ctx, "original_settings")
if diag.HasError() {
    resp.Diagnostics.Append(diag...)
    return
}
```

---

## ISSUE 6

**File Path:** `/workspaces/terraform-provider-power-platform/internal/validators/other_field_required_when_value_of_validator.go`

**Title:** Use of error return values without error handling

**Problem:**
There are multiple instances in the code where functions that return an error (or a diagnostic error) are used, but their error return values are ignored or not handled properly. For example, the calls to `GetAttribute` and `PathMatches` do not check for errors, which can lead to unnoticed issues, panics, or confusion during debugging.

**Impact:**
If errors are not handled properly, this could lead to faulty logic, program panics, or incorrect validation behavior. Severity: **high**.

**Code Issue:**

```go
 currentFieldValue := ""
 _ = req.Config.GetAttribute(ctx, req.Path, &currentFieldValue)
 paths, _ := req.Config.PathMatches(ctx, av.OtherFieldExpression)
```

**Fix:**
Handle the error return values properly:

```go
 currentFieldValue := ""
 diags := req.Config.GetAttribute(ctx, req.Path, &currentFieldValue)
 resp.Diagnostics.Append(diags...)
 if resp.Diagnostics.HasError() {
  return
 }
 
 paths, diags := req.Config.PathMatches(ctx, av.OtherFieldExpression)
 resp.Diagnostics.Append(diags...)
 if resp.Diagnostics.HasError() {
  return
 }
```

---

## ISSUE 7

**File Path:** `/workspaces/terraform-provider-power-platform/internal/services/environment_settings/models.go`

**Title:** Error messages do not provide enough context for troubleshooting

**Problem:**
Error messages in functions such as `convertFromEnvironmentSettingsModel`, for example `errors.New("failed to convert audit settings to ObjectValue")`, are generic and do not provide sufficient details such as the contents of the object or information about what was expected.

**Impact:**
Low severity, but it can hinder debugging efforts and reduce the maintainability of the code. Developers might not have enough information to diagnose issues when these errors are reported.

**Code Issue:**

```go
objectValue, ok := auditSettingsObject.(basetypes.ObjectValue)
if !ok {
 return nil, errors.New("failed to convert audit settings to ObjectValue")
}
```

**Fix:**
Include more context in the error message:

```go
objectValue, ok := auditSettingsObject.(basetypes.ObjectValue)
if !ok {
 return nil, fmt.Errorf("failed to convert audit settings to ObjectValue, got %T: %+v", auditSettingsObject, auditSettingsObject)
}
```

---

## ISSUE 8

**File Path:** `/workspaces/terraform-provider-power-platform/internal/services/environment/datasource_securityroles.go`

**Title:** Missing Error Handling in Data Source Read

**Problem:**
Various API calls and state operations in the security roles data source do not have comprehensive error handling, potentially leading to silent failures or incomplete error reporting.

**Impact:**
Medium severity. Can result in incomplete data or confusing error messages for users.

**Code Issue:**

```go
// Missing error handling pattern
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &result)
return result, err
```

**Fix:**
Add proper error wrapping and context:

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &result)
if err != nil {
    return result, fmt.Errorf("failed to retrieve security roles: %w", err)
}
return result, nil
```

---

## ISSUE 9

**File Path:** `/workspaces/terraform-provider-power-platform/internal/services/tenant/datasource_tenant.go`

**Title:** High Severity Error Handling Missing

**Problem:**
Critical error handling is missing in tenant data source operations, particularly around state management and API calls that could lead to data corruption or silent failures.

**Impact:**
High severity. Missing error handling in critical data source operations can lead to incorrect state or silent failures.

**Code Issue:**

```go
// Pattern of missing error handling in critical operations
resp.State.Get(ctx, &state)
// Missing error check
```

**Fix:**
Add comprehensive error handling:

```go
diags := resp.State.Get(ctx, &state)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
    return
}
```

---

## ISSUE 10

**File Path:** `/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/resource_admin_management_application.go`

**Title:** High Severity Resource Error Handling Issues

**Problem:**
Resource operations lack proper error handling for state management and API operations, which can lead to resource corruption or inconsistent state.

**Impact:**
High severity. Resource operations with missing error handling can corrupt Terraform state.

**Code Issue:**

```go
// Missing error handling in resource operations
resp.State.Set(ctx, &plan)
// No error check
```

**Fix:**
Add proper error handling:

```go
diags := resp.State.Set(ctx, &plan)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
    return
}
```

---

## Summary

**Total Issues:** 10

**Severity Breakdown:**

- Critical: 1 issue
- High: 4 issues
- Medium: 4 issues  
- Low: 1 issue

**Categories:**

- API error wrapping: 3 issues
- State management errors: 3 issues
- Retry loop handling: 1 issue
- Validator error handling: 1 issue
- Diagnostic handling: 1 issue
- Generic error messages: 1 issue

**Key Patterns:**

- Missing error checks on state operations (Get/Set)
- Lack of error context in API calls
- Improper diagnostic handling
- Silent failures in critical operations
- Missing error handling in validators

**Critical Areas Requiring Immediate Attention:**

1. **CreateDataverseUser retry loop** - Critical severity, can cause silent failures
2. **State management operations** - High severity, can corrupt Terraform state
3. **Validator error handling** - High severity, can cause runtime panics

All issues relate to insufficient error handling that can lead to silent failures, corrupted state, or difficult debugging experiences. The critical and high severity issues should be addressed immediately to prevent data corruption and improve system reliability.
