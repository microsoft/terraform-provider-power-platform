# Title

Inconsistent naming: `Api` should be `API` to conform with Go conventions

##

/workspaces/terraform-provider-power-platform/internal/api/client.go

## Problem

The struct `ProviderClient` uses the field name `Api`, which is inconsistent with Go's convention of using initialisms in all caps, e.g., `API`.

## Impact

Violating Go naming conventions can reduce code readability and may confuse developers who expect a consistent style. Severity: **low**

## Location

Declaration of `ProviderClient`

## Code Issue

```go
type ProviderClient struct {
	Config *config.ProviderConfig
	Api    *Client
}
```

## Fix

Change the field name to `API`:

```go
type ProviderClient struct {
	Config *config.ProviderConfig
	API    *Client
}
```

And update all references of `Api` to `API`.
