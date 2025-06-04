# Inconsistent Use of Logging for Error Reporting

##

/workspaces/terraform-provider-power-platform/main.go

## Problem

The code uses `log.Fatalf()` to handle errors during provider serving. Using `log.Fatalf()` will print the error and call `os.Exit(1)`, which abruptly terminates the process. This makes graceful shutdown and deferred resource cleanup impossible, and is discouraged in libraries or shared binaries, though in a `main()` package it's marginally acceptable.

A better practice is to return a proper status code and ensure that any deferred cleanup (if needed) is executed.

## Impact

If deferred functions are ever added to `main`, they will not be run. Abrupt process termination also makes integration testing and diagnosis harder.

Severity: **Medium**

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

Replace with log output and call to `os.Exit(1)` after any possible deferred cleanup, or at least make it clear. Example:

```go
import (
	// ... other imports
	"os"
)

// ... main function

	if err != nil {
		log.Printf("Error serving provider: %s", err)
		os.Exit(1)
	}
```
