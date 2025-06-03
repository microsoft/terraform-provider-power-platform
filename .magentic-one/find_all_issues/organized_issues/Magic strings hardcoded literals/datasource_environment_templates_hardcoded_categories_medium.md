# Title

Hardcoded category names in `appendToList` calls

##

/workspaces/terraform-provider-power-platform/internal/services/environment_templates/datasource_environment_templates.go

## Problem

The `Read` method contains multiple hardcoded calls to `appendToList` with both the property name and the `"category"` string repeated, e.g.:

```go
appendToList(environment_templates.Standard, "standard", &state.Templates)
```

This pattern is repeated for `"premium"`, `"developer"`, etc. If the struct changes or new categories are added/removed, this block must be manually kept in sync, leading to error-prone and non-scalable code.

## Impact

Severity: Medium

- Maintainability issue—difficult to extend or refactor;
- Possible bugs if struct fields and hardcoded categories diverge;
- Reduces readability for future contributors.

## Location

`Read` method, lines with:

```go
appendToList(environment_templates.Standard, "standard", &state.Templates)
appendToList(environment_templates.Premium, "premium", &state.Templates)
...
```

## Code Issue

```go
appendToList(environment_templates.Standard, "standard", &state.Templates)
appendToList(environment_templates.Premium, "premium", &state.Templates)
...
```

## Fix

Use a static slice or map to express the property–category relationship once, then loop, e.g.:

```go
categories := []struct{
    items []itemDto
    name string
}{
    {environment_templates.Standard, "standard"},
    {environment_templates.Premium, "premium"},
    {environment_templates.Developer, "developer"},
    {environment_templates.Basic, "basic"},
    {environment_templates.Production, "production"},
    {environment_templates.Sandbox, "sandbox"},
    {environment_templates.Trial, "trial"},
    {environment_templates.Default, "default"},
    {environment_templates.Support, "support"},
    {environment_templates.SubscriptionBasedTrial, "subscriptionBasedTrial"},
    {environment_templates.Teams, "teams"},
    {environment_templates.Platform, "platform"},
}

for _, c := range categories {
    appendToList(c.items, c.name, &state.Templates)
}
```

Or use reflection if extensibility is needed.

