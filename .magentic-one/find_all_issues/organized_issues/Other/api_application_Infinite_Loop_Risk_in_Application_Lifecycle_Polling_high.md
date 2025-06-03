# Issue 5: Infinite Loop Risk in Application Lifecycle Polling

##

/workspaces/terraform-provider-power-platform/internal/services/application/api_application.go

## Problem

The install logic uses an indefinite `for { ... }` loop to poll for lifecycle completion, but there is no timeout or delay, which could exhaust resources or hang the system if the API never responds with a terminal state.

## Impact

This issue has a **high** severity as it could cause tight CPU loops and hangs, negatively impacting performance and system reliability.

## Location

In `InstallApplicationInEnvironment`, here:

## Code Issue

```go
for {
    lifecycleResponse := environmentApplicationLifecycleDto{}
    response, err := client.Api.Execute(...)
    // ...
    if response.HttpResponse.StatusCode == http.StatusConflict {
        tflog.Debug(ctx, "Lifecycle Operation HTTP Status: '"+response.HttpResponse.Status+"'")
        continue
    }
    // ...
    if lifecycleResponse.Status == "Succeeded" { ... }
    else if lifecycleResponse.Status == "Failed" { ... }
}
```

## Fix

Introduce a sleep/break condition, and optionally a timeout:

```go
maxAttempts := 60
waitInterval := time.Second * 5
for i := 0; i < maxAttempts; i++ {
    lifecycleResponse := environmentApplicationLifecycleDto{}
    response, err := client.Api.Execute(...)
    // ...
    if response.HttpResponse.StatusCode == http.StatusConflict {
        tflog.Debug(ctx, "Lifecycle Operation HTTP Status: '"+response.HttpResponse.Status+"'")
        time.Sleep(waitInterval) // Add delay between checks
        continue
    }
    // ...
    if lifecycleResponse.Status == "Succeeded" { ... }
    else if lifecycleResponse.Status == "Failed" { ... }
}
return "", errors.New("application installation polling timed out")
```
