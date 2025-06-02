# Title

Possible Data Race: Use of Shared p.Config in Concurrent Provider Methods

##

/workspaces/terraform-provider-power-platform/internal/provider/provider.go

## Problem

The `PowerPlatformProvider` struct contains a pointer field `Config *config.ProviderConfig`, which is mutated in the `Configure` method. Other provider methods (`Resources`, `DataSources`, etc.) may be called concurrently by Terraform after configuration. The direct, unsynchronized assignment to `p.Config` fields in `Configure` (and the subsidiary configure functions) opens up the possibility of data races if these fields are read or modified concurrently.

## Impact

Potential race conditions could cause unpredictable behavior in multithreaded environments, including crashes or incorrect configuration being used for resource or data source operations. Severity: **medium**.

## Location

Assignments such as:

```go
p.Config.Urls = *providerConfigUrls
p.Config.Cloud = *cloudConfiguration
p.Config.TelemetryOptout = telemetryOptOut
p.Config.EnableContinuousAccessEvaluation = enableCae
p.Config.TerraformVersion = req.TerraformVersion
```

and similar assignments in other configure* functions.

## Fix

To ensure thread safety, protect access to the configuration with a mutex or use per-operation configuration instances only. Example fix with a mutex:

```go
type PowerPlatformProvider struct {
    Config *config.ProviderConfig
    Api    *api.Client
    mu     sync.RWMutex
}

// When writing:
p.mu.Lock()
p.Config.Urls = *providerConfigUrls
// ...
p.mu.Unlock()

// When reading elsewhere:
p.mu.RLock()
// use p.Config
p.mu.RUnlock()
```

Alternatively, consider making the provider configuration immutable after Configure, or passing a deep copy to every resource/data source operation.
