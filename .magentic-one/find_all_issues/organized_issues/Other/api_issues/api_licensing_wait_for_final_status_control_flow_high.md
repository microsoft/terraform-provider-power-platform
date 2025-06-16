# Tight Control Flow Loop in DoWaitForFinalStatus (No Timeout or Limit)

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/api_licensing.go

## Problem

In `DoWaitForFinalStatus`, the function waits indefinitely for the policy status to reach a terminal state ("Enabled" or "Disabled") but does not provide a timeout or maximum number of retries.

## Impact

High. This can cause the function to block forever if the state never reaches a terminal status, leading to resource exhaustion or stuck operations.

## Location

```go
for {
    billingPolicy, err := client.GetBillingPolicy(ctx, billingId)
    ...
    if billingPolicy.Status == "Enabled" || billingPolicy.Status == "Disabled" {
        return billingPolicy, nil
    }
    ...
}
```

## Fix

Introduce a timeout or a maximum number of retries. Example:

```go
maxRetries := 30
for i := 0; i < maxRetries; i++ {
    billingPolicy, err := c.GetBillingPolicy(ctx, billingId)
    if err != nil {
        return nil, err
    }
    if billingPolicy.Status == "Enabled" || billingPolicy.Status == "Disabled" {
        return billingPolicy, nil
    }
    if err := c.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
        return nil, err
    }
    tflog.Debug(ctx, fmt.Sprintf("Billing Policy Operation State: '%s'", billingPolicy.Status))
}
return nil, fmt.Errorf("wait for final billing policy status exceeded maximum retries")
```
