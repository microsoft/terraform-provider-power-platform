# Inconsistent Go Struct Tag Formatting

## /workspaces/terraform-provider-power-platform/internal/services/environment/dto.go

## Problem

Some struct tags are missing a space after the field type and before the struct tag, which is the idiomatic Go style. For example, many lines look like:

```go
X string`json:"x,omitempty"`
```

Rather than:

```go
X string `json:"x,omitempty"`
```

## Impact

While this does not break compilation or runtime correctness, it is less readable and less idiomatic, and may annoy code reviewers or trigger style linters. Severity: **low**.

## Location

Check all data structs for missing spaces before struct tags throughout the file.

## Code Issue

```go
Id string`json:"id"`
```

## Fix

Add a space between the field definition and the struct tag in all instances.

```go
Id string `json:"id"`
```
