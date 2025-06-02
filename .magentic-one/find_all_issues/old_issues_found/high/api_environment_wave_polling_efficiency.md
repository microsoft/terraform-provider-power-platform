# Title

Inefficient Polling Mechanism in `UpdateFeature`

##

/workspaces/terraform-provider-power-platform/internal/services/environment_wave/api_environment_wave.go

## Problem

The polling logic in the `UpdateFeature` method continuously loops until a condition is met. This might cause unnecessary resource usage.

## Impact

- Inefficient CPU and memory usage during the loop.
- Poor scalability of this code in case of extended retries.
- Severity: **High**

## Location

Polling loop inside the `UpdateFeature` method.

## Code Issues

```go
retryAfter := api.DefaultRetryAfter()
for {
	// Pending feature check...
	err = client.Api.SleepWithContext(ctx, retryAfter)
	if err != nil {
		return nil, err
	}
}
```

## Fix

Introduce a retry limit and log a timeout error if the condition isn't met within a specified number of attempts.

```go
retryAfter := api.DefaultRetryAfter()
maxRetries := 10 // or any suitable upper limit
retries := 0

for retries < maxRetries {
	// Check feature state...

	err = client.Api.SleepWithContext(ctx, retryAfter)
	if err != nil {
		return nil, err
	}

	retries++
}

return nil, fmt.Errorf("failed to enable feature %s within retry limit", featureName)
```