# Unexported Type Name Not Idiomatic

##

/workspaces/terraform-provider-power-platform/internal/services/languages/api_languages.go

## Problem

The `client` type is defined as a struct with a lowercase name, making it unexported. In Go, if the intention is to use this type outside the `languages` package, it should be exported (i.e., named `Client`). If it is deliberately unexported, this is not an issue, but the naming should be reviewed for intent and clarity.

## Impact

If the `client` type is supposed to be used by other packages, keeping it unexported prevents access from outside the package. Severity: **low** (unless package boundaries require it to be exported).

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

If the type should be exported for reuse, capitalize its name:

```go
type Client struct {
	Api *api.Client
}
```
