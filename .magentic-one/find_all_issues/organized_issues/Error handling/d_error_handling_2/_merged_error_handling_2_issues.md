# Error Handling Issues: Recursive Retry Loops and Diagnostic Control Flow

This file contains error handling issues focused on recursive retry mechanisms, missing error context wrapping, and critical diagnostic control flow problems that can lead to data corruption and infinite loops in the environment service operations.

## ISSUE 1

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment/api_environment.go`

### Title

Redundant error checking in recursive error-handling branches

### Problem

Several methods in this file recursively call themselves in error-handling branches (e.g., on HTTP 409/Conflict), but may not always preserve the original stack trace or diagnostics context. Furthermore, the recursion is performed even if the error is not recoverable or could result in an infinite loop if the API repeatedly returns the same error.

### Impact

- Severity: Medium
- Can create the risk of infinite retry loops in pathological cases (API/Service bug or throttling).
- Makes debugging and observing error context more difficult.
- Decreases maintainability by duplicating error retry logic.

### Location

Example from `DeleteEnvironment` and similar in other methods:

```go
if response.HttpResponse.StatusCode == http.StatusConflict {
    err := client.handleHttpConflict(ctx, response)
    if err != nil {
        return err
    }
    return client.DeleteEnvironment(ctx, environmentId)
}
```

Similar patterns are in `CreateEnvironment`, `UpdateEnvironment`, and `UpdateEnvironmentAiFeatures` methods.

### Code Issue

```go
if response.HttpResponse.StatusCode == http.StatusConflict {
    err := client.handleHttpConflict(ctx, response)
    if err != nil {
        return err
    }
    return client.DeleteEnvironment(ctx, environmentId)
}
```

### Fix

Add some form of retry limit (e.g., maximum number of tries or a time budget) or utilize a proper exponential backoff/retry library. Consider returning an error if a conflict persists after several retries, to avoid infinite loops.

Pseudo-code example with a retry limit:

```go
const maxRetries = 10
func (client *Client) DeleteEnvironment(ctx context.Context, environmentId string, retryCount int) error {
    // ...
    if response.HttpResponse.StatusCode == http.StatusConflict {
        if retryCount >= maxRetries {
            return fmt.Errorf("maximum retries reached for DeleteEnvironment on conflict")
        }
        err := client.handleHttpConflict(ctx, response)
        if err != nil {
            return err
        }
        return client.DeleteEnvironment(ctx, environmentId, retryCount+1)
    }
    // ...
}
```

Alternatively, use a for-loop with a retry budget instead of recursion, and propagate retry state cleanly.

## ISSUE 2

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment_settings/api_environment_settings.go`

### Title

Non-idiomatic Error Wrapping

### Problem

Throughout the code, errors from `client.Api.Execute` and related functions are returned directly. If additional context is needed, idiomatic error wrapping using `fmt.Errorf` or `errors.Wrap` (from `pkg/errors`, if used), should be considered to provide more context about the error source.

### Impact

Low/Medium. Lacking error context reduces the ease of debugging errors in higher layers.

### Location

Example in `getEnvironment`:

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
if err != nil {
    return nil, err
}
```

### Code Issue

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
if err != nil {
    return nil, err
}
```

### Fix

Wrap the error with additional context:

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
if err != nil {
    return nil, fmt.Errorf("failed to execute API request for environment %s: %w", environmentId, err)
}
```

## ISSUE 3

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment_templates/api_environment_templates.go`

### Title

Error Handling: Lack of Specific Error Wrapping and Context in API and Unmarshal Errors

### Problem

The function `GetEnvironmentTemplatesByLocation` returns errors directly from either the `client.Api.Execute()` or `json.Unmarshal()` calls. These errors lack contextual wrapping, which would help track where and why the error occurred, especially in complex codebases or logging environments.

### Impact

If an error is propagated up the stack, it would be less informative and harder to diagnose, making debugging more challenging. Severity: **Medium**.

### Location

```go
response, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
if err != nil {
    return templates, err
}
...
err = json.Unmarshal(response.BodyAsBytes, &templates)
if err != nil {
    return templates, err
}
```

### Code Issue

```go
response, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
if err != nil {
    return templates, err
}

defer response.HttpResponse.Body.Close()

err = json.Unmarshal(response.BodyAsBytes, &templates)
if err != nil {
    return templates, err
}
```

### Fix

Wrap the errors with `fmt.Errorf` to provide more context.

```go
response, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
if err != nil {
    return templates, fmt.Errorf("failed to execute API request for environment templates: %w", err)
}

defer response.HttpResponse.Body.Close()

err = json.Unmarshal(response.BodyAsBytes, &templates)
if err != nil {
    return templates, fmt.Errorf("failed to unmarshal environment templates response: %w", err)
}
```

## ISSUE 4

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment/resource_environment.go`

### Title

Error and Warning Handling: Missing Immediate Return After `AddError` or `AddWarning`

### Problem

In several methods (notably, `Read`), after calling `AddError` or `AddWarning` on the diagnostics object, the function does not return immediately. This leads to continued execution, which may act on invalid or incomplete state and possibly cause further, cascading errors.

### Impact

- **Severity**: High
- **Explanation**: Can lead to accessing nil, corrupted, or inconsistent data or duplicate/inconsistent warnings/errors in diagnostics.

### Location

For example, in the `Read` function:

```go
defaultCurrency, err := r.EnvironmentClient.GetDefaultCurrencyForEnvironment(ctx, envDto.Name)
if err != nil {
    if customerrors.Code(err) != customerrors.ERROR_ENVIRONMENT_URL_NOT_FOUND {
        // This is only a warning because you may have BAPI access to the environment but not WebAPI access to dataverse to get currency.
        resp.Diagnostics.AddWarning(fmt.Sprintf("Error when reading default currency for environment %s", envDto.Name), err.Error())
    }

    if !state.Dataverse.IsNull() && !state.Dataverse.IsUnknown() {
        var dataverseSourceModel DataverseSourceModel
        state.Dataverse.As(ctx, &dataverseSourceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
        currencyCode = dataverseSourceModel.CurrencyCode.ValueString()
    }
} else {
    currencyCode = defaultCurrency.IsoCurrencyCode
}

var templateMetadata *createTemplateMetadataDto
var templates []string
if !state.Dataverse.IsNull() && !state.Dataverse.IsUnknown() {
    dv, err := convertEnvironmentCreateLinkEnvironmentMetadataDtoFromDataverseSourceModel(ctx, state.Dataverse)
    if err != nil {
        resp.Diagnostics.AddError("Error when converting dataverse source model to create link environment metadata", err.Error())
        return
    }
    ...
}
```

### Code Issue

```go
dv, err := convertEnvironmentCreateLinkEnvironmentMetadataDtoFromDataverseSourceModel(ctx, state.Dataverse)
if err != nil {
    resp.Diagnostics.AddError("Error when converting dataverse source model to create link environment metadata", err.Error())
    // should return here, else code acts on a nil dv
}
if dv != nil {
    // ...
}
```

**Several similar cases exist throughout the file.**

### Fix

Add explicit `return` statements immediately after `AddError` or, if appropriate, after critical `AddWarning`s that are expected to not continue processing due to potential invalid state.

```go
dv, err := convertEnvironmentCreateLinkEnvironmentMetadataDtoFromDataverseSourceModel(ctx, state.Dataverse)
if err != nil {
    resp.Diagnostics.AddError("Error when converting dataverse source model to create link environment metadata", err.Error())
    return
}
if dv != nil {
    // ...
}
```

---

This must be fixed wherever diagnostic error/warning calls are made and continued execution could lead to misuse of the returned data or incorrect resource state.

---

**Total Issues Found:** 4

**Summary:**

- Low severity: 0 issues
- Medium severity: 3 issues
- High severity: 1 issue

**Primary Categories:**

- Infinite retry loops and recursion issues
- Missing error context and wrapping
- Lack of proper error handling in API calls
- Missing return statements after diagnostic errors
- Continued execution after critical errors

**Focus Areas:**

- Environment service operations have critical error handling gaps
- Consistent pattern of insufficient error wrapping across multiple services
- High-severity issue with diagnostic handling that can cause data corruption
- Need for proper retry mechanisms and error propagation
- Immediate attention required for the high-severity diagnostic return issue
