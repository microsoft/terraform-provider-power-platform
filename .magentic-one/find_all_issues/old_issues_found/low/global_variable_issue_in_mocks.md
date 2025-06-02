# Title

Potential Overuse of Global Variables for Test Configuration

##

/workspaces/terraform-provider-power-platform/internal/mocks/mocks.go

## Problem

The file contains two global variables: `TestUnitTestProtoV6ProviderFactories` and `TestAccProtoV6ProviderFactories`. These serve as test configuration maps. While this approach can centralize configuration, the use of global variables leads to tightly coupled code and less modularity. Additionally, it can cause side effects or unintended interactions in tests.

## Impact

This approach reduces test isolation and makes it difficult to modify tests independently. It can also lead to subtle bugs that are hard to trace. The severity of this issue is **low**, as it only reduces code quality and doesn't directly break functionality.

## Location

- Line 22: `TestUnitTestProtoV6ProviderFactories`
- Line 27: `TestAccProtoV6ProviderFactories`

## Code Issue

```go
var TestUnitTestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"powerplatform": func() (tfprotov6.ProviderServer, error) {
		return providerserver.NewProtocol6WithError(provider.NewPowerPlatformProvider(helpers.UnitTestContext(context.Background(), ""), true))
	},
}

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"powerplatform": func() (tfprotov6.ProviderServer, error) {
		return providerserver.NewProtocol6WithError(provider.NewPowerPlatformProvider(context.Background(), false))
	},
}
```

## Fix

Encapsulate these factory maps within functions to limit their scope and encourage modularity:

```go
func GetUnitTestProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"powerplatform": func() (tfprotov6.ProviderServer, error) {
			return providerserver.NewProtocol6WithError(provider.NewPowerPlatformProvider(helpers.UnitTestContext(context.Background(), ""), true))
		},
	}
}

func GetAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"powerplatform": func() (tfprotov6.ProviderServer, error) {
			return providerserver.NewProtocol6WithError(provider.NewPowerPlatformProvider(context.Background(), false))
		},
	}
}
```

This change ensures that the configuration is only accessible where explicitly required, improving modularity and test isolation.