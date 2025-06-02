# Title

Redundant HTTP Status Handling in `GetEnvironmentGroupRuleSet`

##

/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/api_environment_group_rule_set.go

## Problem

The method redundantly checks the `StatusNoContent` HTTP response status after determining it indicates a missing resource. Instead of wrapping the error into `customerrors.WrapIntoProviderError` and separately having another `fmt.Errorf` for the same condition, this duplication adds unnecessary complexity and redundancy.

## Impact

The redundant handling of `StatusNoContent` increases code complexity and leads to poor readability. Severity: Low.

## Location

In the `GetEnvironmentGroupRuleSet`, lines 49-56.

## Code Issue

```go
if resp.HttpResponse.StatusCode == http.StatusNoContent {
	return nil, customerrors.WrapIntoProviderError(err, customerrors.ERROR_OBJECT_NOT_FOUND, "rule set '%s' not found")
}

if len(environmentGroupRuleSet.Value) == 0 {
	return nil, fmt.Errorf("no environment group ruleset found for environment group id %s", environmentGroupId)
}
```

## Fix

Combine the handling of both conditions into a single block to eliminate redundancy:

```go
if resp.HttpResponse.StatusCode == http.StatusNoContent || len(environmentGroupRuleSet.Value) == 0 {
	return nil, customerrors.WrapIntoProviderError(nil, customerrors.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("no environment group ruleset found for environment group id '%s'", environmentGroupId))
}
```

This cleanup simplifies the logic while maintaining the correct behavior.

---
