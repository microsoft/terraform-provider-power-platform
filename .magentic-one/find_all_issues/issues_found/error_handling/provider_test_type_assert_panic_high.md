# Title

Potential Panic If Type Assertion Fails

##

internal/provider/provider_test.go

## Problem

The tests use this pattern:

```go
provider.NewPowerPlatformProvider(context.Background())().(*provider.PowerPlatformProvider)
```

Here, the result of `provider.NewPowerPlatformProvider(context.Background())()` is asserted as `*provider.PowerPlatformProvider`. If the provider signature ever changes or a nil value is returned, this will panic, causing the test to crash rather than fail with an informative error. This is dangerous in test code as it might mask root causes for failures.

## Impact

High. If an implementation detail changes, the test will panic (abort) rather than fail gracefully, making CI/CD troubleshooting harder.

## Location

```go
datasources := provider.NewPowerPlatformProvider(context.Background())().(*provider.PowerPlatformProvider).DataSources(context.Background())

resources := provider.NewPowerPlatformProvider(context.Background())().(*provider.PowerPlatformProvider).Resources(context.Background())
```

## Fix

Use a "comma ok" idiom to check for assertion safety and provide meaningful errors:

```go
ppp, ok := provider.NewPowerPlatformProvider(context.Background())().(*provider.PowerPlatformProvider)
require.True(t, ok, "Expected *PowerPlatformProvider, got %T", ppp)
datasources := ppp.DataSources(context.Background())
```
