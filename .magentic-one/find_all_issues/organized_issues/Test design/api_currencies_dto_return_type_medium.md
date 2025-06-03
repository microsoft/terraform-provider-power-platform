# Title

Direct Exposure of DTO in Public Method

##

/workspaces/terraform-provider-power-platform/internal/services/currencies/api_currencies.go

## Problem

The method `GetCurrenciesByLocation` returns a `currenciesDto` type directly. If this type is not exported (`currenciesDto` is likely a package-private type), returning it directly in a public API can cause problems with extensibility, testability, or code comprehension—especially if users from other packages wish to use this method. Also, returning non-standard or unexported types is generally discouraged in Go unless strictly for internal usage.

## Impact

**Medium**. Encapsulation, maintainability, and potential for API misuse if used outside the intended scope.

## Location

Signature of the `GetCurrenciesByLocation` method:

```go
func (client *client) GetCurrenciesByLocation(ctx context.Context, location string) (currenciesDto, error)
```

## Fix

Either export the `currenciesDto` type or introduce a standard public type for the external method signature.

```go
// Example: Exported DTO
type CurrenciesDTO struct {
    // fields...
}

func (client *CurrenciesClient) GetCurrenciesByLocation(ctx context.Context, location string) (CurrenciesDTO, error)
```
