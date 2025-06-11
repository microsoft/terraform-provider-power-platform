# Error Wrapping and Context Loss

##

/workspaces/terraform-provider-power-platform/internal/services/powerapps/api_powerapps.go

## Problem

On error, functions return raw errors from called functions (`return nil, err`), losing valuable context about the operation that failed.

## Impact

Makes debugging more difficult, as it’s harder to trace the origin and cause of errors. Severity: Medium.

## Location

Within the `GetPowerApps` function:

```go
	if err != nil {
		return nil, err
	}
```

Happens in two places in the method.

## Code Issue

```go
	envs, err := client.environmentClient.GetEnvironments(ctx)
	if err != nil {
		return nil, err
	}
	...
		_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &appsArray)
		if err != nil {
			return nil, err
		}
```

## Fix

Wrap the errors to add context:

```go
	envs, err := client.environmentClient.GetEnvironments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get environments: %w", err)
	}
	...
		_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &appsArray)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch power apps for environment %s: %w", env.Name, err)
		}
```
