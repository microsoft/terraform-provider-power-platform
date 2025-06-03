# Issue: Struct Field Naming Can Be Improved for Go Idioms

## 
/workspaces/terraform-provider-power-platform/internal/services/capacity/api_capacity.go

## Problem

The `client` struct has a field named `Api` (capitalized). In Go, commonly, struct fields are named using lowerCamelCase unless they need to be exported. If this is only used internally (as implied by the package scoping), it should be named with a lowercase `a` to follow Go idioms, improving readability and distinguishability from types.

## Impact

Severity: **low**

This is a style/readability issue, but following conventions ensures code is idiomatic and easier for Go engineers to read and maintain.

## Location

Struct definition:

```go
type client struct {
    Api *api.Client
}
```

## Code Issue

```go
type client struct {
    Api *api.Client
}
```

## Fix

Use lowerCamelCase for private struct fields unless there's a requirement for it to be exported:

```go
type client struct {
    api *api.Client
}
```

And update all corresponding usages in methods (e.g. `client.Api` → `c.api`).
