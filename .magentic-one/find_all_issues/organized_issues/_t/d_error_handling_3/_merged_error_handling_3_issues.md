# Error Handling Issues: API Error Wrapping and Data Conversion Context

This document consolidates API error wrapping issues, data conversion error handling, and missing contextual information in error propagation across multiple services including admin management, analytics, billing, and copilot studio components.

---

## ISSUE 1

**File Path:** `/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/api_admin_management_application.go`

**Title:** No Error Context Wrapping or Logging for API Calls

**Problem:**
When API errors occur, they are returned directly from the method without any log or error wrapping/context. It can make debugging difficult, as the caller will not know which API call or parameters led to the error, particularly important with multiple similar methods.

**Impact:**
Medium. Affects debuggability and observability.

**Code Issue:**

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &adminApp)
return &adminApp, err
```

**Fix:**
Wrap errors with context, using e.g., `fmt.Errorf` or the `%w` verb, or log them if appropriate.

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &adminApp)
if err != nil {
    return nil, fmt.Errorf("failed to get admin app %s: %w", clientId, err)
}
return &adminApp, nil
```

---

## ISSUE 2

**File Path:** `/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/api_analytics_data_exports.go`

**Title:** Inconsistent or Missing Error Wrapping for API Error in GetGatewayCluster

**Problem:**
In `GetGatewayCluster`, an error from the API call is returned as-is, losing the opportunity to wrap it with meaningful context (as is done elsewhere in the code). Consistent error wrapping ensures error traces are useful at every level.

**Impact:**
Without wrapping the error, debugging and error diagnosis by downstream consumers is impaired. Missing context can make troubleshooting harder and error roots unclear. Severity: Low.

**Code Issue:**

```go
 _, err = client.Api.Execute(ctx, nil, "GET", tenantApiUrl.String(), nil, nil, []int{http.StatusOK}, &gatewayCluster)
 if err != nil {
  return nil, err
 }
```

**Fix:**
Wrap the error to include API context, such as:

```go
 _, err = client.Api.Execute(ctx, nil, "GET", tenantApiUrl.String(), nil, nil, []int{http.StatusOK}, &gatewayCluster)
 if err != nil {
  return nil, fmt.Errorf("failed to execute GetGatewayCluster API request: %w", err)
 }
```

---

## ISSUE 3

**File Path:** `/workspaces/terraform-provider-power-platform/internal/services/capacity/api_capacity.go`

**Title:** Lack of Error Wrapping on API Invocation

**Problem:**
In the `GetTenantCapacity` method, if an error occurs during the execution of the API call, it is returned directly without wrapping or contextualizing. This makes it harder to trace where the error originated from when debugging, especially in larger codebases with many API calls. Proper error wrapping (with `fmt.Errorf("...: %w", err)`) allows for easier and more informative debugging.

**Impact:**
Severity: **medium**

Directly returning low-level errors without additional context reduces maintainability and makes future debugging and log tracing more cumbersome.

**Code Issue:**

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &dto)
if err != nil {
    return nil, err
}
```

**Fix:**
Wrap the error to provide function context:

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &dto)
if err != nil {
    return nil, fmt.Errorf("failed to get tenant capacity: %w", err)
}
```

---

## ISSUE 4

**File Path:** `/workspaces/terraform-provider-power-platform/internal/services/connection/api_connection.go`

**Title:** Potential Error Wrapping Issue Missing in API Calls

**Problem:**
Not all API errors are wrapped with contextual information. While some errors are wrapped using `customerrors.WrapIntoProviderError`, others are directly returned. Without consistent error wrapping, debugging and tracing issues becomes more difficult.

**Impact:**
This inconsistency in error handling can make it harder to identify the source and context of errors, impacting debuggability and support. Severity: Medium.

**Code Issue:**

```go
_, err := client.Api.Execute(ctx, nil, "PUT", apiUrl.String(), nil, connectionToCreate, []int{http.StatusCreated}, &connection)
if err != nil {
    return nil, err
}
```

**Fix:**
Consistently wrap errors returned from API calls with additional context using the existing error wrapping strategy or custom error types.

```go
if err != nil {
    return nil, fmt.Errorf("failed to create connection: %w", err)
}
```

Or using the project's custom error approach if desired:

```go
if err != nil {
    return nil, customerrors.WrapIntoProviderError(err, customerrors.ERROR_OBJECT_CREATION_FAILED, "Failed to create connection")
}
```

---

## ISSUE 5

**File Path:** `/workspaces/terraform-provider-power-platform/internal/services/connectors/api_connectors.go`

**Title:** Ineffective Error Handling: Lost Context in API Calls

**Problem:**
When `client.Api.Execute` returns an error, the function simply returns the error as-is without any additional context about which API call failed or what part of the `GetConnectors` operation was unsuccessful. This means the error will lack crucial debugging information.

**Impact:**
It can be difficult to determine which of the several API calls failed. This can complicate debugging and error reporting, making support and maintenance harder, especially when end users report API failures. Severity: **medium**.

**Code Issue:**

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connectorArray)
if err != nil {
 return nil, err
}
```

**Fix:**
Wrap the errors with context to indicate which API call failed:

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connectorArray)
if err != nil {
 return nil, fmt.Errorf("failed to fetch PowerApps connectors: %w", err)
}
```

---

## ISSUE 6

**File Path:** `/workspaces/terraform-provider-power-platform/internal/services/data_record/api_data_record.go`

**Title:** Control flow: Missing error wrap/explicit cause in many error returns

**Problem:**
Errors from called functions and API responses are often returned directly, without wrapping or providing helpful context for consumers. For example, when `json.Unmarshal` fails, or when a required field is not found, the error is returned as-is. This makes debugging and error tracing harder, especially in a complex service where multiple underlying operations can fail for the same high-level reason.

**Impact:**
**Severity: Medium**

- Errors in Terraform logs may not provide enough context to understand what failed, especially for unmarshalling, HTTP, or type assertion problems.
- Can slow debugging cycles for both developers and users.

**Code Issue:**

```go
err = json.Unmarshal(response.BodyAsBytes, &mapResponse)
if err != nil {
    return nil, err
}
```

**Fix:**
Wrap errors with context using `fmt.Errorf` and `%w` to capture the cause:

```go
err = json.Unmarshal(response.BodyAsBytes, &mapResponse)
if err != nil {
    return nil, fmt.Errorf("unmarshalling response bytes for [operation context]: %w", err)
}
```

Apply this best practice for all returned errors, especially those returned from other helpers, type assertions, or third-party packages (json, url, etc).

---

## ISSUE 7

**File Path:** `/workspaces/terraform-provider-power-platform/internal/services/environment_wave/api_environment_wave.go`

**Title:** Insufficient Error Context in Returned Errors

**Problem:**
The errors returned from the API and environment calls (such as `client.Api.Execute` and `client.environmentClient.GetEnvironment`) are returned directly without additional context or wrapping. This makes it harder to trace the origin of the error, especially in a larger codebase where similar errors might occur in multiple places.

**Impact:**
This can make debugging and support more challenging, as the caller has less information about where and why an error occurred. Severity: **medium**

**Code Issue:**

```go
 if err != nil {
  return nil, err
 }
```

**Fix:**
Wrap errors with additional context using `fmt.Errorf` (or `errors.Wrap` if using the pkg/errors library):

```go
 if err != nil {
  return nil, fmt.Errorf("failed to execute API call for organizations: %w", err)
 }
```

Do this for all error returns where additional context could be valuable.

---

## ISSUE 8

**File Path:** `/workspaces/terraform-provider-power-platform/internal/services/solution/api_solution.go`

**Title:** Improper Use of Error Variable in Custom Error Wrapping for Not Found Solution (GetSolutionUniqueName)

**Problem:**
In the `GetSolutionUniqueName` method, if `len(solutions.Value) == 0`, the returned error is created by wrapping `err` (which is `nil` at this point) with a new error via `customerrors.WrapIntoProviderError`. Passing `nil` as the error argument can be misleading and is not idiomatic Go practice. The same pattern occurs in `GetSolutionById`.

**Impact:**
This can result in Go errors whose underlying cause is `nil`, reducing code clarity and making debugging harder. Severity: **medium**, because it impairs error traceability and could cause confusion in diagnostics.

**Code Issue:**

```go
if len(solutions.Value) == 0 {
 return nil, customerrors.WrapIntoProviderError(err, customerrors.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("solution with unique name '%s' not found", name))
}
```

**Fix:**
Create a new standard error to wrap, so the underlying error is meaningful.

```go
if len(solutions.Value) == 0 {
 baseErr := fmt.Errorf("solution with unique name '%s' not found", name)
 return nil, customerrors.WrapIntoProviderError(baseErr, customerrors.ERROR_OBJECT_NOT_FOUND, baseErr.Error())
}
```

And for GetSolutionById:

```go
if len(solutions.Value) == 0 {
 baseErr := fmt.Errorf("solution with id '%s' not found", solutionId)
 return nil, customerrors.WrapIntoProviderError(baseErr, customerrors.ERROR_OBJECT_NOT_FOUND, baseErr.Error())
}
```

---

## ISSUE 9

**File Path:** `/workspaces/terraform-provider-power-platform/internal/services/tenant/api_tenant.go`

**Title:** Possible Error Handling Omission by Not Wrapping Errors

**Problem:**
When `client.Api.Execute` returns an error, it is forwarded as-is. Wrapping errors with context (e.g., using `fmt.Errorf("...: %w", err)`) provides better stack traces and debugging information, aiding in tracking the source of an error throughout the codebase.

**Impact:**
Without contextual error wrapping, debugging is harder, and error logs may not provide enough information about where or why failures occur. **Severity:** Medium

**Code Issue:**

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &dto)
if err != nil {
 return nil, err
}
```

**Fix:**
Wrap errors with extra context information before returning:

```go
import "fmt"

//...

_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &dto)
if err != nil {
 return nil, fmt.Errorf("failed to execute GET tenant API call: %w", err)
}
```

---

## ISSUE 10

**File Path:** `/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/api_tenant_isolation_policy.go`

**Title:** Lack of Proper Error Wrapping and Propagation in getTenantIsolationPolicy

**Problem:**
The function `getTenantIsolationPolicy` directly returns the error it receives from the API call without any context or error wrapping. This can make it hard for callers to determine the source of the error, especially in larger codebases where multiple API calls might propagate generic errors up the stack.

**Impact:**
Medium. While the error is returned, debugging and tracing the error's origin can become difficult, affecting maintainability and troubleshooting. Developers and support engineers may struggle to identify the failure's context quickly.

**Code Issue:**

```go
 if err != nil {
  return nil, err
 }
```

**Fix:**
Wrap the error with a relevant context to improve traceability:

```go
 if err != nil {
  return nil, fmt.Errorf("could not retrieve tenant isolation policy for tenant %s: %w", tenantId, err)
 }
```

---

## ISSUE 11

**File Path:** `/workspaces/terraform-provider-power-platform/internal/helpers/cert.go`

**Title:** Unwrapped Errors on I/O and Certificate Decoding

**Problem:**
Throughout the file, errors are returned directly from lower-level library/system functions (like `os.ReadFile`, `pkcs12.DecodeChain`, base64 decode). This exposes raw error messages, which can be less useful for consumers of this package, making debugging less contextual and error handling less robust. Wrapping these errors with higher-level context provides a clearer indication of where and why the failure occurred.

**Impact:**
Returning unwrapped errors in exported functions leads to less maintainable and debuggable code, especially when this package is integrated with larger projects. Severity: **medium**.

**Code Issue:**

```go
pfx, err := os.ReadFile(certificateFilePath)
if err != nil {
    return "", err
}
```

**Fix:**
Add context to error messages using `fmt.Errorf("context: %w", err)` especially in exported functions.

```go
pfx, err := os.ReadFile(certificateFilePath)
if err != nil {
    return "", fmt.Errorf("failed to read certificate file '%s': %w", certificateFilePath, err)
}
```

And in `convertByteToCert`:

```go
key, cert, _, err := pkcs12.DecodeChain(certData, password)
if err != nil {
    return nil, nil, fmt.Errorf("failed to decode PKCS12 certificate chain: %w", err)
}
```

---

## ISSUE 12

**File Path:** `/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/models.go`

**Title:** Insufficient Error Wrapping in Data Conversion

**Problem:**
Both `createAppInsightsConfigDtoFromSourceModel` and `convertAppInsightsConfigModelFromDto` return a simple error without rich context. If errors are ever generated in these functions (such as after adding validation, or if ValueString/ValueBool might error in the future), they should be wrapped or annotated to provide meaningful call-site information.

**Impact:**
Severity: **Medium**

Lack of context in errors makes debugging and tracing issues more difficult in complex workflows.

**Code Issue:**

```go
return nil, fmt.Errorf("EnvironmentId cannot be empty")
```

**Fix:**
Use Go 1.13+ error wrapping for context where appropriate.

```go
return nil, fmt.Errorf("failed to create AppInsightsConfigDto: EnvironmentId cannot be empty")
```

And similarly for other errors.

---

## Summary

**Total Issues:** 12

**Severity Breakdown:**

- High: 0 issues
- Medium: 11 issues
- Low: 1 issue

**Categories:**

- API error wrapping: 9 issues
- Custom error handling: 1 issue
- Data conversion errors: 1 issue
- Certificate/file handling errors: 1 issue

All issues relate to the lack of proper error context wrapping throughout the codebase. The consistent pattern is that API calls and other operations return errors directly without adding contextual information, making debugging and troubleshooting more difficult. The recommended fix across all issues is to use `fmt.Errorf` with the `%w` verb to wrap errors with meaningful context.
