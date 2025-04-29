# Title

Absence of unit tests for `main()` function

##

`/workspaces/terraform-provider-power-platform/main.go`

## Problem

The `main()` function lacks any form of unit tests or testability hooks. It directly defines its behavior without abstraction, making it difficult to write unit tests. Testing the `main()` function in its current state requires starting the entire program, which is not ideal for modular or isolated testing.

## Impact

- Reduces test coverage and increases the risk of introducing bugs in application startup logic.
- Makes the startup flow harder to modify and troubleshoot due to lack of unit test hooks.
- Severity: High.

## Location

Entire `main` function:
Lines 13-32:
`main.go`

## Code Issue

```go
func main() {
	log.Printf("[INFO] Starting the Power Platform Terraform Provider %s %s", common.ProviderVersion, common.Branch)

	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()
	ctx := context.Background()

	serveOpts := providerserver.ServeOpts{
		Debug:   debug,
		Address: "registry.terraform.io/microsoft/power-platform",
	}

	err := providerserver.Serve(ctx, provider.NewPowerPlatformProvider(ctx), serveOpts)

	if err != nil {
		log.Fatalf("Error serving provider: %s", err)
	}
}
```

## Fix

To make the `main()` function testable, it is recommended to encapsulate its logic in another function (e.g., `RunServer()`) that can be independently tested or mocked for unit tests.

```go
// Encapsulate logic within RunServer
func RunServer(ctx context.Context, debug bool) error {
	log.Printf("[INFO] Starting the Power Platform Terraform Provider %s %s", common.ProviderVersion, common.Branch)

	serveOpts := providerserver.ServeOpts{
		Debug:   debug,
		Address: "registry.terraform.io/microsoft/power-platform",
	}

	return providerserver.Serve(ctx, provider.NewPowerPlatformProvider(ctx), serveOpts)
}

// main function delegates its logic
func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()
	ctx := context.Background()

	if err := RunServer(ctx, debug); err != nil {
		log.Fatalf("Error serving provider: %s", err)
	}
}
```

This modification allows for testability by writing isolated unit tests for `RunServer()`. 