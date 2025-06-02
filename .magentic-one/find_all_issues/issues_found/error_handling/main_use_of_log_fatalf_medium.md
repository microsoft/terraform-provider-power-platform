# Use of `log.Fatalf` Causes Application Exit on Error 

##

/workspaces/terraform-provider-power-platform/main.go

## Problem

The `log.Fatalf` function is used to handle errors returned from `providerserver.Serve()`. This function logs the error message and then immediately terminates the application using `os.Exit(1)`. While this is common for command-line Go tools, abrupt termination may not allow for resource clean-up or defer execution, and is not always idiomatic Go error handling. It also makes testing this path more cumbersome.

## Impact

Severity: **Medium**

- Immediate exit can obscure more graceful error management or logging/metrics collection.
- Harder to test this failure mode programmatically (since it exits the process).
- It prevents any deferred clean-up from running.

## Location

```go
	if err != nil {
		log.Fatalf("Error serving provider: %s", err)
	}
```

## Code Issue

```go
	if err != nil {
		log.Fatalf("Error serving provider: %s", err)
	}
```

## Fix

Consider using `log.Printf` and returning an error exit code with `os.Exit` so that defers run, if any. If special clean-up is needed, handle it before process exit:

```go
import (
	// ...
	"os"
)

// ...

if err != nil {
	log.Printf("Error serving provider: %s", err)
	os.Exit(1)
}
```

Or, for a more robust and testable approach, consider refactoring `main()` to delegate to a function returning `error`:

```go
func realMain() error {
	log.Printf("[INFO] Starting the Power Platform Terraform Provider %s %s", common.ProviderVersion, common.Branch)
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()
	ctx := context.Background()
	serveOpts := providerserver.ServeOpts{
		Debug:   debug,
		Address: "registry.terraform.io/microsoft/power-platform",
	}
	return providerserver.Serve(ctx, provider.NewPowerPlatformProvider(ctx), serveOpts)
}

func main() {
	if err := realMain(); err != nil {
		log.Printf("Error serving provider: %s", err)
		os.Exit(1)
	}
}
```
