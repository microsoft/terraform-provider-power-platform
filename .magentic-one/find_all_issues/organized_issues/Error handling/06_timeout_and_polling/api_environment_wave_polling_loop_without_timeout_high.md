# Polling Loop Without Timeout or Delay Customization

##

/workspaces/terraform-provider-power-platform/internal/services/environment_wave/api_environment_wave.go

## Problem

The polling loop inside `UpdateFeature` will continue indefinitely until the state changes or an error occurs. There is no timeout or maximum number of retries.

## Impact

Potential for the function to hang indefinitely if the upgrade never completes, which is a significant operational risk. Severity: **high**

## Location

Within the `UpdateFeature` method, the polling loop:

```go
	for {
		feature, err := client.GetFeature(ctx, environmentId, featureName)
		...
		err = client.Api.SleepWithContext(ctx, retryAfter)
		...
	}
```

## Code Issue

```go
	retryAfter := api.DefaultRetryAfter()
	for {
		feature, err := client.GetFeature(ctx, environmentId, featureName)
		if err != nil {
			return nil, err
		}

		if feature != nil && feature.AppsUpgradeState != "Upgrading" {
			tflog.Info(ctx, fmt.Sprintf("Feature %s  with state: %s", featureName, feature.AppsUpgradeState))
			return feature, nil
		}

		err = client.Api.SleepWithContext(ctx, retryAfter)
		if err != nil {
			return nil, err
		}

		tflog.Debug(ctx, fmt.Sprintf("Feature %s not yet enabled, polling...", featureName))
	}
```

## Fix

Add context deadline/timeout or a maximum number of retries to prevent infinite loops:

```go
	attempts := 0
	maxAttempts := 20
	retryAfter := api.DefaultRetryAfter()

	for attempts < maxAttempts {
		feature, err := client.GetFeature(ctx, environmentId, featureName)
		if err != nil {
			return nil, err
		}

		if feature != nil && feature.AppsUpgradeState != "Upgrading" {
			tflog.Info(ctx, fmt.Sprintf("Feature %s  with state: %s", featureName, feature.AppsUpgradeState))
			return feature, nil
		}

		err = client.Api.SleepWithContext(ctx, retryAfter)
		if err != nil {
			return nil, err
		}

		tflog.Debug(ctx, fmt.Sprintf("Feature %s not yet enabled, polling...", featureName))
		attempts++
	}
	return nil, fmt.Errorf("timed out waiting for feature %s to be enabled in environment %s", featureName, environmentId)
```

Or use the context deadline with `ctx.Done()` in the loop.
