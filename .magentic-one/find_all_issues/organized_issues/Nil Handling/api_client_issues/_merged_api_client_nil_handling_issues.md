# API Client Nil Handling Issues

This document contains all identified nil handling issues related to API client components in the Terraform Power Platform provider codebase.

## ISSUE 1

<!-- Source: api_connectors-medium-3.md -->

# Missing Type Safety and Response Validation in Connectors API

##

/workspaces/terraform-provider-power-platform/internal/services/connectors/api_connectors.go

## Problem

After API calls (`client.Api.Execute`), the code assumes correct types and expected structures are always returned (such as `connectorArray.Value`, `virtualConnectorArray`, etc.). There is no check that fields parsed from the response are valid, not nil, or have the expected structure. This is risky in loosely-typed API scenarios.

## Impact

If the API changes or unexpected/malformed data is returned, this might cause nil dereferences and panics, or result in silent data corruption/invalid output. Severity: **medium** (can break production at runtime if API contracts drift or errors occur).

## Location

Everywhere that code uses fields like `connectorArray.Value` and properties on elements of the API arrays directly, especially in loops and appends, such as:

```go
for inx, connector := range connectorArray.Value { ... }
```

## Code Issue

```go
for inx, connector := range connectorArray.Value {
  // Assumes connectorArray.Value is present and fully valid
}
```

## Fix

Validate types and content before iterating or dereferencing. For example:

```go
if connectorArray.Value == nil {
  return nil, fmt.Errorf("connectorArray.Value is nil or missing in API response")
}
```

Further, implement error handling for missing/null/invalid fields everywhere these are used, and ensure test coverage.

## ISSUE 2

<!-- Source: api_environment_add_dataverse_nil_pointer_high.md -->

# Issue: Unhandled error value in AddDataverseToEnvironment

##

/workspaces/terraform-provider-power-platform/internal/services/environment/api_environment.go

## Problem

In the `AddDataverseToEnvironment` method, errors from API calls and header parsing are logged but not always handled properly. For example, after calling `client.Api.Execute`, an error is logged but execution continues. If `apiResponse` is nil due to an error, dereferencing it later will cause a panic.

## Impact

- Severity: High
- This may lead to nil pointer dereference panics during runtime and inconsistent or unexpected execution flow.
- Logging the error is not sufficient: the calling function may expect a valid return value when the request actually failed.

## Location

Within the `AddDataverseToEnvironment` function:

```go
apiResponse, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, environmentCreateLinkEnvironmentMetadata, []int{http.StatusAccepted}, nil)
if err != nil {
    tflog.Error(ctx, "Error adding Dataverse to environment: "+err.Error())
}

tflog.Debug(ctx, "Environment Creation Operation HTTP Status: '"+apiResponse.HttpResponse.Status+"'")
```

## Code Issue

```go
apiResponse, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, environmentCreateLinkEnvironmentMetadata, []int{http.StatusAccepted}, nil)
if err != nil {
    tflog.Error(ctx, "Error adding Dataverse to environment: "+err.Error())
}

tflog.Debug(ctx, "Environment Creation Operation HTTP Status: '"+apiResponse.HttpResponse.Status+"'")
```

## Fix

Return immediately after logging the error to prevent further operations on a possibly nil `apiResponse`.

```go
apiResponse, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, environmentCreateLinkEnvironmentMetadata, []int{http.StatusAccepted}, nil)
if err != nil {
    tflog.Error(ctx, "Error adding Dataverse to environment: "+err.Error())
    return nil, err
}
if apiResponse == nil {
    return nil, errors.New("unexpected nil response from AddDataverseToEnvironment")
}

tflog.Debug(ctx, "Environment Creation Operation HTTP Status: '"+apiResponse.HttpResponse.Status+"'")
```

## ISSUE 3

<!-- Source: api_environment_delete_nil_pointer_high.md -->

# Issue: Incorrect error handling and nil pointer access in DeleteEnvironment

##

/workspaces/terraform-provider-power-platform/internal/services/environment/api_environment.go

## Problem

The function `DeleteEnvironment` does not always check for `err != nil` immediately after the `Api.Execute` method call. It proceeds to perform logic on the `response` object (accessing `response.HttpResponse.StatusCode`, etc.) that could be nil if an error occurs, which can cause a panic at runtime.

## Impact

- Severity: High
- Can lead to panics and crashes at runtime if `response` is nil and error is not handled immediately.
- Having incorrect error-handling logic also makes the code harder to maintain and debug.

## Location

Within the `DeleteEnvironment` function:

```go
response, err := client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, environmentDelete, []int{http.StatusNoContent, http.StatusAccepted, http.StatusConflict, http.StatusNotFound}, nil)

// Handle HTTP 404 case - if the environment is not found, consider it already deleted
if response != nil && response.HttpResponse.StatusCode == http.StatusNotFound {
    tflog.Info(ctx, fmt.Sprintf("Environment '%s' not found. Treating as successfully deleted.", environmentId))
    return nil
}

if response.HttpResponse.StatusCode == http.StatusConflict {
    err := client.handleHttpConflict(ctx, response)
    if err != nil {
        return err
    }
    return client.DeleteEnvironment(ctx, environmentId)
}
```

## Code Issue

```go
response, err := client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, environmentDelete, []int{http.StatusNoContent, http.StatusAccepted, http.StatusConflict, http.StatusNotFound}, nil)

if response != nil && response.HttpResponse.StatusCode == http.StatusNotFound {
    // ...
}

if response.HttpResponse.StatusCode == http.StatusConflict {
    // ...
}
```

## Fix

Check `err` immediately after the API call, and only continue if it is `nil`. Ensure `response` is not nil before dereferencing, and consolidate the error and response handling for clarity and safety.

```go
response, err := client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, environmentDelete, []int{http.StatusNoContent, http.StatusAccepted, http.StatusConflict, http.StatusNotFound}, nil)
if err != nil {
    return err
}
if response == nil { // paranoid check, should not happen if err is nil, but for robustness
    return errors.New("unexpected nil response in DeleteEnvironment")
}

if response.HttpResponse.StatusCode == http.StatusNotFound {
    tflog.Info(ctx, fmt.Sprintf("Environment '%s' not found. Treating as successfully deleted.", environmentId))
    return nil
}

if response.HttpResponse.StatusCode == http.StatusConflict {
    herr := client.handleHttpConflict(ctx, response)
    if herr != nil {
        return herr
    }
    return client.DeleteEnvironment(ctx, environmentId)
}
```

## ISSUE 4

<!-- Source: api_environment_group_critical.md -->

# Incomplete Error Handling for nil HttpResponse in GetEnvironmentGroup

##

/workspaces/terraform-provider-power-platform/internal/services/environment_groups/api_environment_group.go

## Problem

In the `GetEnvironmentGroup` method, the code accesses `httpResponse.HttpResponse.StatusCode` without checking if `httpResponse` is `nil`. This can cause a panic if `client.Api.Execute` returns an error and a nil httpResponse, violating the contract of robust Go error handling.

## Impact

This introduces a potential runtime panic, specifically a "nil pointer dereference," which is a critical runtime error and can cause the process to crash.

**Severity:** Critical

## Location

```go
func (client *client) GetEnvironmentGroup(ctx context.Context, environmentGroupId string) (*environmentGroupDto, error) {
 // ...
 httpResponse, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNotFound}, &environmentGroup)
 if httpResponse.HttpResponse.StatusCode == http.StatusNotFound {
  return nil, nil
 } else if err != nil {
  return nil, err
 }
 // ...
}
```

## Fix

Add a nil check for `httpResponse` and its `HttpResponse` field before dereferencing.

```go
func (client *client) GetEnvironmentGroup(ctx context.Context, environmentGroupId string) (*environmentGroupDto, error) {
 // ...
 httpResponse, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNotFound}, &environmentGroup)
 if httpResponse != nil && httpResponse.HttpResponse != nil && httpResponse.HttpResponse.StatusCode == http.StatusNotFound {
  return nil, nil
 } else if err != nil {
  return nil, err
 }
 // ...
}
```

## ISSUE 5

<!-- Source: api_languages_Not_Checking_for_nil_Before_Closing_Response_Body_high.md -->

# Error Handling and Resource Management: Not Checking for `nil` Before Closing Response Body

##

/workspaces/terraform-provider-power-platform/internal/services/languages/api_languages.go

## Problem

The code calls `defer response.HttpResponse.Body.Close()` without ensuring that `response.HttpResponse` or its `Body` field are not nil. If `response.HttpResponse` is nil, a panic may occur.

## Impact

May cause the application to panic at runtime if `HttpResponse` is nil. Severity: **high**.

## Location

```go
defer response.HttpResponse.Body.Close()
```

## Code Issue

```go
defer response.HttpResponse.Body.Close()
```

## Fix

Check for `nil` before deferring the closure:

```go
if response.HttpResponse != nil && response.HttpResponse.Body != nil {
    defer response.HttpResponse.Body.Close()
}
```

## ISSUE 6

<!-- Source: api_locations_error_handling_medium.md -->

# Title

Missing nil check on parameter `apiClient` in `newLocationsClient`

##

/workspaces/terraform-provider-power-platform/internal/services/locations/api_locations.go

## Problem

The `newLocationsClient` function does not check if its parameter `apiClient` is nil. This could lead to runtime panics when the returned `client` is used.

## Impact

Medium severity. Could lead to runtime panics if an invalid `apiClient` pointer is passed.

## Location

`newLocationsClient` function

## Code Issue

```go
func newLocationsClient(apiClient *api.Client) client {
 return client{
  Api: apiClient,
 }
}
```

## Fix

Add a check and optionally return an error or panic with a clear message.

```go
func newClient(apiClient *api.Client) (client, error) {
 if apiClient == nil {
  return client{}, fmt.Errorf("apiClient cannot be nil")
 }
 return client{
  API: apiClient,
 }, nil
}
```

## ISSUE 7

<!-- Source: api_powerapps_error_wrapping_context_loss_medium.md -->

# Error Wrapping and Context Loss

##

/workspaces/terraform-provider-power-platform/internal/services/powerapps/api_powerapps.go

## Problem

On error, functions return raw errors from called functions (`return nil, err`), losing valuable context about the operation that failed.

## Impact

Makes debugging more difficult, as it's harder to trace the origin and cause of errors. Severity: Medium.

## Location

Within the `GetPowerApps` function:

```go
 if err != nil {
  return nil, err
 }
```

Happens in two places in the method.

## Code Issue

```go
 envs, err := client.environmentClient.GetEnvironments(ctx)
 if err != nil {
  return nil, err
 }
 ...
  _, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &appsArray)
  if err != nil {
   return nil, err
  }
```

## Fix

Wrap the errors to add context:

```go
 envs, err := client.environmentClient.GetEnvironments(ctx)
 if err != nil {
  return nil, fmt.Errorf("failed to get environments: %w", err)
 }
 ...
  _, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &appsArray)
  if err != nil {
   return nil, fmt.Errorf("failed to fetch power apps for environment %s: %w", env.Name, err)
  }
```

## ISSUE 8

<!-- Source: api_rest_headers_nil_map_high.md -->

# Headers map may be nil

##

/workspaces/terraform-provider-power-platform/internal/services/rest/api_rest.go

## Problem

In `SendOperation`, the `headers` map is only initialized when `operation.Headers` has items. If there are no headers, it remains `nil` and is passed to subsequent functions, which expect a map and iterate over it, potentially leading to a panic when assigning headers in `ExecuteApiRequest`.

## Impact

Potential high-severity runtime panic from operations on a nil map. This can cause hard-to-debug issues and service downtime.

## Location

Lines ~27-48, relevant in both `SendOperation` and `ExecuteApiRequest`:

## Code Issue

```go
 var headers map[string]string
 // ...
 if len(operation.Headers) > 0 {
  headers = make(map[string]string)
  for _, h := range operation.Headers {
   headers[h.Name.ValueString()] = h.Value.ValueString()
  }
 }
    // ...
 for k, v := range headers {
  h.Add(k, v)
 }
```

## Fix

Initialize the `headers` map as an empty map by default, ensuring it is never nil.

```go
 headers := make(map[string]string)
 if len(operation.Headers) > 0 {
  for _, h := range operation.Headers {
   headers[h.Name.ValueString()] = h.Value.ValueString()
  }
 }
```

## ISSUE 9

<!-- Source: api_rest_nil_response_handling_high.md -->

# Insufficient error handling and control flow for nil response

##

/workspaces/terraform-provider-power-platform/internal/services/rest/api_rest.go

## Problem

In `SendOperation`, there's a conditional block handling the case when `res == nil && err != nil`, but the general logic outside this block assumes that `res.BodyAsBytes` is always safe to access, which could result in a runtime panic if `res` is nil.

## Impact

Potential high-severity runtime panic due to dereferencing a nil pointer, impacting server reliability and correctness.

## Location

Lines ~60-68:

## Code Issue

```go
 if res == nil && err != nil {
  output["body"] = types.StringValue(err.Error())
 } else {
  if len(res.BodyAsBytes) > 0 {
   output["body"] = types.StringValue(string(res.BodyAsBytes))
  }
 }
```

## Fix

Check for `res` being non-nil before accessing its members. This prevents runtime panics and makes the code more robust.

```go
 if res == nil && err != nil {
  output["body"] = types.StringValue(err.Error())
 } else if res != nil && len(res.BodyAsBytes) > 0 {
  output["body"] = types.StringValue(string(res.BodyAsBytes))
 }
```

## ISSUE 10

<!-- Source: api_rest_scope_string_medium.md -->

# Counterintuitive parameter design for `scope`

##

/workspaces/terraform-provider-power-platform/internal/services/rest/api_rest.go

## Problem

The `scope` parameter for `ExecuteApiRequest` is a pointer to a string but must not be nil (function returns an error otherwise). This creates unnecessary complexity and forces the caller into handling pointer logic, with little gain since nil is an error.

## Impact

Severity: Medium. Promotes confusing API and parameter handling, possibly propagating nil pointer patterns unnecessarily throughout the code base.

## Location

Within `ExecuteApiRequest`:

## Code Issue

```go
 if scope == nil {
  return nil, errors.New("invalid input: scope must be provided")
 }

 return client.Api.Execute(ctx, []string{*scope}, method, url, h, body, expectedStatusCodes, nil)
```

## Fix

Accept `scope` as a string (not a pointer), and enforce presence at compile-time via type signature:

```go
func (client *Client) ExecuteApiRequest(ctx context.Context, scope string, url, method string, body *string, headers map[string]string, expectedStatusCodes []int) (*api.Response, error) {
 h := http.Header{}
 for k, v := range headers {
  h.Add(k, v)
 }
 return client.Api.Execute(ctx, []string{scope}, method, url, h, body, expectedStatusCodes, nil)
}
```

## ISSUE 11

<!-- Source: api_solution_code_duplication_low.md -->

# Title

Code Duplication in API Response Handling

##

internal/services/solution/api_solution.go

## Problem

In several methods (e.g., `GetSolutionUniqueName`, `GetSolutionById`, `GetSolutions`, `CreateSolution`, `DeleteSolution`, `GetTableData`, `validateSolutionImportResult`), there exists repeated code for handling forbidden and not found HTTP responses right after each `Execute` call:

```go
if err := client.Api.HandleForbiddenResponse(resp); err != nil {
    return nil, err
}
if err := client.Api.HandleNotFoundResponse(resp); err != nil {
    return nil, err
}
```

This repetition across almost every method reduces maintainability and increases the risk of inconsistency if the error handling logic ever changes.

## Impact

Severity: **low**. While this does not present an immediate bug, it decreases maintainability and contributes to code bloat.

## Location

Most functions, e.g.,

```go
resp, err := client.Api.Execute(...)
if err != nil {
    return nil, err
}
if err := client.Api.HandleForbiddenResponse(resp); err != nil {
    return nil, err
}
if err := client.Api.HandleNotFoundResponse(resp); err != nil {
    return nil, err
}
```

## Code Issue

```go
if err := client.Api.HandleForbiddenResponse(resp); err != nil {
    return nil, err
}
if err := client.Api.HandleNotFoundResponse(resp); err != nil {
    return nil, err
}
```

## Fix

Extract the error handling into a helper function, for example:

```go
func handleCommonApiErrors(api *api.Client, resp *http.Response) error {
    if err := api.HandleForbiddenResponse(resp); err != nil {
        return err
    }
    if err := api.HandleNotFoundResponse(resp); err != nil {
        return err
    }
    return nil
}
```

And then in each method:

```go
if err := handleCommonApiErrors(client.Api, resp); err != nil {
    return nil, err
}
```

## ISSUE 12

<!-- Source: api_solution_response_handling_high.md -->

# Title

Incorrect HTTP Response Handling: Reusing Response from Previous API Call

##

internal/services/solution/api_solution.go

## Problem

In `CreateSolution`, after the initial "StageSolution" POST, subsequent POST and GET requests are made (most notably to `ImportSolutionAsync` and the `asyncoperations` endpoint). After each such request, the code runs:

```go
if err := client.Api.HandleForbiddenResponse(resp); err != nil {
    return nil, err
}
if err := client.Api.HandleNotFoundResponse(resp); err != nil {
    return nil, err
}
```

However, for the second and subsequent API invocations (to `ImportSolutionAsync` and inside the async polling loop), the variable `resp` is not updated with the result of those `Execute` calls—only the error variable `err` is. This means the response object being inspected by the forbidden/notfound handlers is stale and may lead to wrong error handling, masking HTTP errors and resulting in undetected failures.

## Impact

Severity: **high**. This results in incorrect error handling control flow after asynchronous POST and GET requests and can conceal HTTP errors, resulting in misleading function success or masked failures.

## Location

Main problematic location(s):

```go
_, err = client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, importSolutionRequestBody, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &importSolutionResponse)
if err != nil {
    return nil, err
}
if err := client.Api.HandleForbiddenResponse(resp); err != nil {
    return nil, err
}
if err := client.Api.HandleNotFoundResponse(resp); err != nil {
    return nil, err
}

// and inside the for loop:
_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &asyncSolutionPullResponse)
if err != nil {
    return nil, err
}
if err := client.Api.HandleForbiddenResponse(resp); err != nil {
    return nil, err
}
if err := client.Api.HandleNotFoundResponse(resp); err != nil {
    return nil, err
}
```

## Code Issue

```go
_, err = client.Api.Execute(...)
if err != nil {
    return nil, err
}
if err := client.Api.HandleForbiddenResponse(resp); err != nil {
    return nil, err
}
if err := client.Api.HandleNotFoundResponse(resp); err != nil {
    return nil, err
}
```

where `resp` is not updated by the most recent Execute call.

## Fix

Correctly capture and use the response object returned by `Execute` each time, instead of using a stale or previously set reference:

```go
resp, err = client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, importSolutionRequestBody, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &importSolutionResponse)
if err != nil {
    return nil, err
}
if err := client.Api.HandleForbiddenResponse(resp); err != nil {
    return nil, err
}
if err := client.Api.HandleNotFoundResponse(resp); err != nil {
    return nil, err
}

// In the polling loop:
resp, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &asyncSolutionPullResponse)
if err != nil {
    return nil, err
}
if err := client.Api.HandleForbiddenResponse(resp); err != nil {
    return nil, err
}
if err := client.Api.HandleNotFoundResponse(resp); err != nil {
    return nil, err
}
```

## ISSUE 13

<!-- Source: api_tenant_settings_reflection_panic_high.md -->

# Unsafe Reflection Handling and Potential Panic in filterDto

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/api_tenant_settings.go

## Problem

In the `filterDto` function, the use of reflection for nil checking and field access can panic if assumptions about types or interface values are not consistently held. Particularly, the repetitive use of `.Elem()` and kind checking without verification of pointer-ness or underlying value validity can cause runtime panics. For instance:

```go
if !configuredFieldValue.IsNil() && !backendFieldValue.IsNil() && backendFieldValue.IsValid() && outputField.CanSet() {
    if configuredFieldValue.Kind() == reflect.Pointer && configuredFieldValue.Elem().Kind() == reflect.Struct {
        outputStruct := filterDto(ctx, configuredFieldValue.Elem().Interface(), backendFieldValue.Elem().Interface())
        outputField.Set(reflect.ValueOf(outputStruct))
    } else if configuredFieldValue.Kind() == reflect.Pointer && configuredFieldValue.Elem().Kind() == reflect.Bool {
        boolValue := backendFieldValue.Elem().Bool()
        newBool := bool(boolValue)
        outputField.Set(reflect.ValueOf(&newBool))
    }
    // ... similar for string and int64
}
```

If any `configuredFieldValue` or `backendFieldValue` are not valid pointers, or are nil when calling `.Elem()`, this will panic.

## Impact

- **Severity: High**
- Causes runtime panics if any value is not set as expected, which can crash the provider.
- Makes code fragile and difficult to maintain.
- Type safety violation due to unchecked kind and value handling.

## Location

```go
if !configuredFieldValue.IsNil() && !backendFieldValue.IsNil() && backendFieldValue.IsValid() && outputField.CanSet() {
    if configuredFieldValue.Kind() == reflect.Pointer && configuredFieldValue.Elem().Kind() == reflect.Struct {
        outputStruct := filterDto(ctx, configuredFieldValue.Elem().Interface(), backendFieldValue.Elem().Interface())
        outputField.Set(reflect.ValueOf(outputStruct))
    } // ...
}
```

## Code Issue

```go
if !configuredFieldValue.IsNil() && !backendFieldValue.IsNil() && backendFieldValue.IsValid() && outputField.CanSet() {
    if configuredFieldValue.Kind() == reflect.Pointer && configuredFieldValue.Elem().Kind() == reflect.Struct {
        outputStruct := filterDto(ctx, configuredFieldValue.Elem().Interface(), backendFieldValue.Elem().Interface())
        outputField.Set(reflect.ValueOf(outputStruct))
    }
    // ... etc
}
```

## Fix

Refactor to ensure correct order of nil, validity, and kind checking before dereferencing pointers or calling `.Elem()`. E.g.:

```go
if configuredFieldValue.IsValid() && configuredFieldValue.Kind() == reflect.Pointer && !configuredFieldValue.IsNil() &&
   backendFieldValue.IsValid() && backendFieldValue.Kind() == reflect.Pointer && !backendFieldValue.IsNil() &&
   outputField.CanSet() {

   elemKind := configuredFieldValue.Elem().Kind()
   switch elemKind {
   case reflect.Struct:
       outputStruct := filterDto(ctx, configuredFieldValue.Elem().Interface(), backendFieldValue.Elem().Interface())
       outputField.Set(reflect.ValueOf(outputStruct))
   case reflect.Bool:
       boolValue := backendFieldValue.Elem().Bool()
       newBool := bool(boolValue)
       outputField.Set(reflect.ValueOf(&newBool))
   case reflect.String:
       stringValue := backendFieldValue.Elem().String()
       newString := string(stringValue)
       outputField.Set(reflect.ValueOf(&newString))
   case reflect.Int64:
       int64Value := backendFieldValue.Elem().Int()
       newInt64 := int64(int64Value)
       outputField.Set(reflect.ValueOf(&newInt64))
   default:
       tflog.Debug(ctx, fmt.Sprintf("Skipping unknown field type %s", elemKind))
   }
}
```

This ensures only valid, non-nil pointers of expected kinds are dereferenced, preventing panics and maintaining type safety.

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
