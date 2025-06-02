# Title

Invalid Map Initialization with Function Invocation

##

/workspaces/terraform-provider-power-platform/internal/mocks/mocks.go

## Problem

The `TestUnitTestProtoV6ProviderFactories` map attempts to initialize its value using an inline function call `providerserver.NewProtocol6WithError`. However, the returned value from `provider.NewPowerPlatformProvider` does not correctly match the expected signature `(tfprotov6.ProviderServer, error)` due to incorrect usage of the additional function call `true)()`.

## Impact

This causes a compilation error for improperly matching the type signature required in the map. As this issue blocks the test setup, the impact is **high**.

## Location

Line 22 of the file `/workspaces/terraform-provider-power-platform/internal/mocks/mocks.go`.

## Code Issue

```go
var TestUnitTestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"powerplatform": providerserver.NewProtocol6WithError(provider.NewPowerPlatformProvider(helpers.UnitTestContext(context.Background(), ""), true)()),
}
```

## Fix

Replace the inline invocation to match the required function signature. For example:

```go
var TestUnitTestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"powerplatform": func() (tfprotov6.ProviderServer, error) {
		return providerserver.NewProtocol6WithError(provider.NewPowerPlatformProvider(helpers.UnitTestContext(context.Background(), ""), true))
	},
}
```

This fix adds a lambda function that wraps the inner call to `providerserver.NewProtocol6WithError`, thereby ensuring the proper function signature required by the map initializer.