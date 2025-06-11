# Naming: Struct Field Name Not Idiomatic

##

/workspaces/terraform-provider-power-platform/internal/services/languages/api_languages.go

## Problem

The field `Api` in the struct should be named `API` as per Go naming conventions for acronyms.

## Impact

Non-standard naming can reduce code readability and maintainability, especially in large codebases. Severity: **low**.

## Location

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

Rename the field to use the all-caps acronym:

```go
type client struct {
	API *api.Client
}
```
