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
