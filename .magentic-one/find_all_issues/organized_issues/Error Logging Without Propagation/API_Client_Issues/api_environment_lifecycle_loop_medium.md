# Issue: Nested and duplicated logic blocks in lifecycle loops reduce maintainability

##

/workspaces/terraform-provider-power-platform/internal/services/environment/api_environment.go

## Problem

The lifecycle-wait loops inside `AddDataverseToEnvironment` and `UpdateEnvironment` contain deep/nested branches, repeated logging, repeated calls to `SleepWithContext`, with similar but slightly differing retry and continuation conditions. This reduces readability and makes future changes risky or error-prone.

## Impact

- Severity: Medium
- Makes the logic harder to refactor, test, or diagnose if lifecycle APIs change.
- Repeated code invites future inconsistencies.

## Location

Example from `AddDataverseToEnvironment`:

```go
for {
    lifecycleEnv := &EnvironmentDto{}
    lifecycleResponse, err := client.Api.Execute(ctx, nil, "GET", locationHeader, nil, nil, []int{http.StatusOK, http.StatusAccepted, http.StatusConflict}, &lifecycleEnv)
    if err != nil {
        return nil, err
    }

    tflog.Debug(ctx, fmt.Sprintf("Dataverse Creation Operation HTTP Status: '%s'", lifecycleResponse.HttpResponse.Status))
    if lifecycleResponse.HttpResponse.StatusCode == http.StatusConflict {
        continue
    }

    if lifecycleEnv == nil || lifecycleEnv.Properties == nil {
        tflog.Debug(ctx, fmt.Sprintf("The environment lifecycle response body did not match expected format. Response status code: %s", lifecycleResponse.HttpResponse.Status))
        continue
    }

    err = client.Api.SleepWithContext(ctx, retryAfter)
    if err != nil {
        return nil, err
    }

    tflog.Debug(ctx, fmt.Sprintf("Dataverse Creation Operation State: '%s'", lifecycleEnv.Properties.ProvisioningState))

    if lifecycleEnv.Properties.ProvisioningState == "Succeeded" {
        return lifecycleEnv, nil
    } else if lifecycleEnv.Properties.ProvisioningState != "LinkedDatabaseProvisioning" && lifecycleEnv.Properties.ProvisioningState != "Succeeded" {
        return lifecycleEnv, fmt.Errorf("dataverse creation failed. provisioning state: %s", lifecycleEnv.Properties.ProvisioningState)
    }
}
```

## Code Issue

```go
// This pattern with multiple continues, sleep, triple-checks, and repeated log/delay blocks is seen in multiple places.
// The same applies to similar lifecycle-wait loops in UpdateEnvironment and related code.
```

## Fix

Refactor to flatten branches, extract common logging and delay logic to helper functions, and document states being checked for clarity:

```go
for {
    lifecycleEnv := &EnvironmentDto{}
    lifecycleResponse, err := client.Api.Execute(ctx, nil, "GET", locationHeader, nil, nil, []int{http.StatusOK, http.StatusAccepted, http.StatusConflict}, &lifecycleEnv)
    if err != nil {
        return nil, err
    }
    tflog.Debug(ctx, fmt.Sprintf("Dataverse Creation Operation HTTP Status: '%s'", lifecycleResponse.HttpResponse.Status))
    
    if lifecycleResponse.HttpResponse.StatusCode == http.StatusConflict {
        // Optionally count conflicts/retries, or sleep and continue
        continue
    }
    if lifecycleEnv == nil || lifecycleEnv.Properties == nil {
        tflog.Debug(ctx, fmt.Sprintf("lifecycle response malformed; code: %s", lifecycleResponse.HttpResponse.Status))
        continue
    }
    tflog.Debug(ctx, fmt.Sprintf("Dataverse Creation State: '%s'", lifecycleEnv.Properties.ProvisioningState))
    if lifecycleEnv.Properties.ProvisioningState == ProvisioningStateSucceeded {
        return lifecycleEnv, nil
    }
    if lifecycleEnv.Properties.ProvisioningState != ProvisioningStateLinkedDatabaseProvisioning {
        return lifecycleEnv, fmt.Errorf("dataverse creation failed, state: %s", lifecycleEnv.Properties.ProvisioningState)
    }
    if err := client.Api.SleepWithContext(ctx, retryAfter); err != nil {
        return nil, err
    }
}
```
Extract provisioning checks, logging, and delays to avoid code duplication.

---

This issue should be saved under:  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/api_environment_lifecycle_loop_medium.md`
