# Unexported Client Type with Exported Methods

##

/workspaces/terraform-provider-power-platform/internal/services/powerapps/api_powerapps.go

## Problem

The `client` type is unexported (lowercase `client`), but has methods attached that may be intended for use outside the package.

## Impact

This limits usability of the type outside this package, even though methods like `GetPowerApps` may be expected to be used by other packages. Severity: Medium.

## Location

Line 19 in the file.

## Code Issue

```go
type client struct {
	Api               *api.Client
	environmentClient environment.Client
}
```

## Fix

Export the `Client` type:

```go
type Client struct {
	Api               *api.Client
	environmentClient environment.Client
}
```

And adjust all constructors and uses accordingly.
