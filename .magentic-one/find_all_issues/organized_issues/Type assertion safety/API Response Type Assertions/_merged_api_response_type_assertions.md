# API Response Type Assertion Safety Issues

This document contains type assertion safety issues related to API response handling, data extraction, and response parsing in the codebase.

## ISSUE 1

# Title

Possible Panic Due to Type Assertion Without Existence/Type Check in getPrincipalString

##

/workspaces/terraform-provider-power-platform/internal/services/connection/api_connection.go

## Problem

The function `getPrincipalString` performs a type assertion on `principal[key]` without checking if the key exists or the actual interface value type. While it returns an error on failure, if `principal` is `nil` or lacks the specified key, it might cause an unexpected program state.

## Impact

This could cause a program panic during iteration or when extracting data from the response if the structure of the map changes or if unexpected data is returned from the API. Severity: Medium.

## Location

```go
func getPrincipalString(principal map[string]any, key string) (string, error) {
 value, ok := principal[key].(string)
 if !ok {
  return "", fmt.Errorf("failed to convert principal %s to string", key)
 }
 return value, nil
}
```

## Code Issue

```go
value, ok := principal[key].(string)
if !ok {
 return "", fmt.Errorf("failed to convert principal %s to string", key)
}
```

## Fix

Add an existence check before type assertion:

```go
func getPrincipalString(principal map[string]any, key string) (string, error) {
 raw, exists := principal[key]
 if !exists {
  return "", fmt.Errorf("principal key %s does not exist", key)
 }
 value, ok := raw.(string)
 if !ok {
  return "", fmt.Errorf("failed to convert principal %s to string", key)
 }
 return value, nil
}
```

---

*This file should be saved in:*
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/api_connection_principal_type_assertion_medium.md

---

## ISSUE 2

# Title

Panic risk: unchecked type assertions on interface{} values in map context

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/api_data_record.go

## Problem

Throughout the file, there are a number of unchecked type assertions from the `any` (i.e. `interface{}`) type, e.g.:

- `response["@odata.context"].(string)`
- `response["value"].([]any)`
- `item.(map[string]any)`
If these type assertions fail (because the map is missing a key, or the value is of a different type), the code will panic.

Examples include:

- `pluralName := strings.Split(response["@odata.context"].(string), "#")[1]`
- `valueSlice, ok := response["value"].([]any)`
- `if value, ok := mapResponse["value"].([]any)[0].(map[string]any); ok { ... }`

In most locations only a basic check for `nil` is present, or in some places, a type assertion is not protected by an `ok` check at all.

## Impact

- **Severity: High**
- Can cause runtime panics if the remote API response changes structure or isn't as expected, crashing the calling Terraform operation.
- Causes nondeterministic error conditions which are hard to diagnose and fix in production.
- Introduces security and stability risk.

## Location

Example location:

```go
pluralName := strings.Split(response["@odata.context"].(string), "#")[1]
```

But similar issues are present in various locations in the file.

## Code Issue

```go
pluralName := strings.Split(response["@odata.context"].(string), "#")[1]

if mapResponse["value"] != nil && len(mapResponse["value"].([]any)) > 0 {
    if value, ok := mapResponse["value"].([]any)[0].(map[string]any); ok {
        if logicalName, ok := value["LogicalName"].(string); ok {
            result = logicalName
        }
    }
} else if logicalName, ok := mapResponse["LogicalName"].(string); ok {
    result = logicalName
} else {
    return nil, errors.New("logicalName field not found in result when retrieving table singular name")
}
```

## Fix

Wrap all type assertions in a two-part "comma ok" assertion and check for presence before use. Return a meaningful error instead of panicking. For example:

```go
odataCtxRaw, exists := response["@odata.context"]
if !exists {
    return nil, errors.New("@odata.context field missing from response")
}
odataCtx, ok := odataCtxRaw.(string)
if !ok {
    return nil, errors.New("@odata.context field is not a string")
}
splitParts := strings.Split(odataCtx, "#")
if len(splitParts) < 2 {
    return nil, errors.New("@odata.context string is malformed")
}
pluralName := splitParts[1]
if index := strings.IndexAny(pluralName, "(/"); index != -1 {
    pluralName = pluralName[:index]
}
```

And similarly for usages of slices and maps:

```go
if rawVal, ok := response["value"]; ok && rawVal != nil {
    valueSlice, ok := rawVal.([]any)
    if !ok {
        return nil, errors.New("value field is not of type []any")
    }
    for _, item := range valueSlice {
        valueMap, ok := item.(map[string]any)
        if !ok {
            return nil, errors.New("item is not of type map[string]any")
        }
        records = append(records, valueMap)
    }
} else {
    records = append(records, response)
}
```

Apply this pattern at all locations where a map or slice is type asserted.

---

Path for saving:  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/api_data_record_panic_high.md`

---

## ISSUE 3

# Type Assertion without Proper Error Handling in applyCorrections

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/api_tenant_settings.go

## Problem

Within the `applyCorrections` function, the code uses a type assertion after a call to `filterDto`:

```go
corrected, ok := correctedFilter.(*tenantSettingsDto)
if !ok {
    tflog.Error(ctx, "Type assertion to failed in applyCorrections")
    return nil
}
```

While an error is logged if the assertion fails, this surface-level error handling potentially obscures the root cause, allows `nil` values to be returned silently, and fails to provide actionable information to upstream callers. It is preferable to return an explicit error and allow calling functions to react accordingly. The signature of `applyCorrections` should be changed to account for this possibility.

## Impact

- **Severity: Medium**
- Returning `nil` silently can result in unexpected panics or misbehavior further up the call stack.
- Insufficient transparency for debugging and error propagation.
- Reduces reliability and maintainability.

## Location

```go
func applyCorrections(ctx context.Context, planned tenantSettingsDto, actual tenantSettingsDto) *tenantSettingsDto {
    correctedFilter := filterDto(ctx, planned, actual)
    corrected, ok := correctedFilter.(*tenantSettingsDto)
    if !ok {
        tflog.Error(ctx, "Type assertion to failed in applyCorrections")
        return nil
    }
    ...
}
```

## Code Issue

```go
corrected, ok := correctedFilter.(*tenantSettingsDto)
if !ok {
    tflog.Error(ctx, "Type assertion to failed in applyCorrections")
    return nil
}
```

## Fix

Return an error from the function rather than a `nil` pointer, and update callers accordingly.

```go
func applyCorrections(ctx context.Context, planned tenantSettingsDto, actual tenantSettingsDto) (*tenantSettingsDto, error) {
    correctedFilter := filterDto(ctx, planned, actual)
    corrected, ok := correctedFilter.(*tenantSettingsDto)
    if !ok {
        tflog.Error(ctx, "Type assertion failed in applyCorrections")
        return nil, fmt.Errorf("type assertion to *tenantSettingsDto failed in applyCorrections")
    }

    // ... (rest of function unchanged)

    return corrected, nil
}
```

Callers of `applyCorrections` must now handle the error explicitly, which will promote more robust error propagation.

---

## ISSUE 4

# Issue with Error Handling on Dynamic Type Assertion

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/datasource_data_record.go

## Problem

In multiple methods, especially in type assertion blocks (e.g., in `convertColumnsToState`, `buildObjectValueFromX`, `buildExpandObject`), the code is using direct type assertions like `columns[key].(bool)` without handling cases where the assertion might fail. If the assertion fails, a panic will be triggered. This could cause the program to crash unexpectedly if the data returned by the API isn't as expected.

## Impact

Crashing due to a failed type assertion is a critical issue for stability and error handling, as it makes the provider brittle in the face of unexpected data. Severity: **critical**.

## Location

Multiple methods, including:

- `convertColumnsToState`
- `buildObjectValueFromX`
- `buildExpandObject`

## Code Issue

```go
switch value.(type) {
case bool:
    caseBool(columns[key].(bool), attributes, attributeTypes, key)
...
// repeated in all functions with type assertions
```

## Fix

You should use the "comma ok" idiom on type assertion to avoid a panic and handle the error gracefully. For example:

```go
switch v := value.(type) {
case bool:
    caseBool(v, attributes, attributeTypes, key)
case int64:
    caseInt64(v, attributes, attributeTypes, key)
case float64:
    caseFloat64(v, attributes, attributeTypes, key)
case string:
    caseString(v, attributes, attributeTypes, key)
case map[string]any:
    typ, val, _ := d.buildObjectValueFromX(v)
    tupleElementType := types.ObjectType{
        AttrTypes: typ,
    }
    objVal, _ := types.ObjectValue(typ, val)
    attributes[key] = objVal
    attributeTypes[key] = tupleElementType
case []any:
    typeObj, valObj := d.buildExpandObject(v)
    attributeTypes[key] = typeObj
    attributes[key] = valObj
default:
    // optionally log or handle unexpected type
}
```

Add similar handling in other relevant methods. This will prevent panics due to unexpected types and let you report or skip unhandled data types without crashing the provider.

---

## ISSUE 5

# Title

Unchecked Type Assertion on ProviderData in Configure

##

/workspaces/terraform-provider-power-platform/internal/services/solution_checker_rules/datasource_solution_checker_rules.go

## Problem

In the `Configure` method, there is a type assertion on `req.ProviderData` to `*api.ProviderClient`. If this assertion fails, an error is logged via `resp.Diagnostics.AddError` and the function returns. While a diagnostic is added, the code after this point assumes the type assertion has succeeded (e.g., subsequent usage of `client.Api`). This pattern, while accepted in Terraform providers, may lead to subtle failures if the error is not handled thoroughly downstream, or if additional initialization logic is added later. Defensive programming suggests treating failed type assertions as fatal errors or using more robust error handling.

## Impact

Potential for future control flow bugs or nil pointer dereference if refactoring occurs and additional code after the assertion assumes a valid client. Severity is **high**, since improper error handling here impacts provider setup, potentially resulting in failed provider initialization or misleading error messages.

## Location

Line starting:

```go
client, ok := req.ProviderData.(*api.ProviderClient)
```

## Code Issue

```go
client, ok := req.ProviderData.(*api.ProviderClient)
if !ok {
    resp.Diagnostics.AddError(
        "Unexpected ProviderData Type",
        fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
    )
    return
}
d.SolutionCheckerRulesClient = newSolutionCheckerRulesClient(client.Api)
```

## Fix

Perform the assignment to `d.SolutionCheckerRulesClient` only if the assertion is successful, and consider adding a test for this branch. Document that a failed assertion is considered a terminal configuration error.

```go
client, ok := req.ProviderData.(*api.ProviderClient)
if !ok {
    resp.Diagnostics.AddError(
        "Unexpected ProviderData Type",
        fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
    )
    // Optionally: panic or return here, as expected.
    return
}
// Proceed safely knowing client is valid
if client != nil {
    d.SolutionCheckerRulesClient = newSolutionCheckerRulesClient(client.Api)
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
