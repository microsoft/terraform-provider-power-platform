# Inconsistent Naming of Struct Fields in client

##

/workspaces/terraform-provider-power-platform/internal/services/environment_groups/api_environment_group.go

## Problem

The `client` struct uses mixed naming conventions for its fields (e.g., `Api` vs. `RuleSetApi` vs. `TenantApi`). In Go, consistent naming (all abbreviations either cased as ID/URL/API or not) enhances code readability and maintainability.

## Impact

This reduces code clarity, may lead to confusion, and makes refactoring and code review more difficult.

**Severity:** Low

## Location

```go
type client struct {
	Api        *api.Client
	TenantApi  tenant.Client
	RuleSetApi environment_group_rule_set.Client
}
```

## Fix

Choose a consistent convention for abbreviations (preferably "API" as all-caps, conforming to Go idioms for initialisms).

```go
type client struct {
	API        *api.Client
	TenantAPI  tenant.Client
	RuleSetAPI environment_group_rule_set.Client
}
```

And update all usages accordingly in the whole codebase
