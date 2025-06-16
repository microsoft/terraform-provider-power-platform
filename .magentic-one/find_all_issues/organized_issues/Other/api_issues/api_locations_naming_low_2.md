# Title

Struct field `Api` does not follow Go naming conventions

##

/workspaces/terraform-provider-power-platform/internal/services/locations/api_locations.go

## Problem

Struct field `Api` should follow Go's ID capitalization rules; abbreviations should appear as full capitals (`API` instead of `Api`). This makes the code more idiomatic and improves readability.

## Impact

Low severity. Reduces code consistency and may make searching for fields less predictable.

## Location

Definition of the `client` struct:

## Code Issue

```go
type client struct {
	Api *api.Client
}
```

## Fix

Change the field name from `Api` to `API`.

```go
type client struct {
	API *api.Client
}
```

Also, update all references in the file accordingly (for example, `client.Api` should become `client.API`).
