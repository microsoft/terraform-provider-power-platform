# Unnecessary Use of Pointer for receiver

##

/workspaces/terraform-provider-power-platform/internal/services/powerapps/api_powerapps.go

## Problem

The `GetPowerApps` method is defined on a pointer receiver `*client` despite the struct’s fields either being pointers themselves or unexported, and no apparent mutation occurs.

## Impact

Unnecessarily using pointer receivers can be avoided for better clarity and to signal immutability. Severity: Low.

## Location

```go
func (client *client) GetPowerApps(ctx context.Context) ([]powerAppBapiDto, error)
```

## Code Issue

```go
func (client *client) GetPowerApps(ctx context.Context) ([]powerAppBapiDto, error)
```

## Fix

If mutation isn’t needed, use value receiver:

```go
func (client client) GetPowerApps(ctx context.Context) ([]powerAppBapiDto, error)
```

Alternatively, if mutation is (or will be) needed, this can be left as-is. Consider reviewing if pointer receivers are necessary.
