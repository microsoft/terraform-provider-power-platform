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
