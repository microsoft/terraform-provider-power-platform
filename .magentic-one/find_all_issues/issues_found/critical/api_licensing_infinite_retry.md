# Inconsistent Retry Mechanism in `DoWaitForFinalStatus`

**Path:** `/workspaces/terraform-provider-power-platform/internal/services/licensing/api_licensing.go`

---

## Description

The `DoWaitForFinalStatus` method uses a "sleep with context" mechanism to retry fetching the billing policy until it reaches a terminal state (`Enabled` or `Disabled`). However, the method lacks a configurable timeout or maximum number of retries. If a billing policy fails to reach a terminal state due to an external system issue, the method could potentially retry indefinitely, causing resource blocking.

---

## Observed Code

```go
func (client *Client) DoWaitForFinalStatus(ctx context.Context, billingPolicyDto *BillingPolicyDto) (*BillingPolicyDto, error) {
    billingId := billingPolicyDto.Id

    for {
        billingPolicy, err := client.GetBillingPolicy(ctx, billingId)

        if err != nil {
            return nil, err
        }

        if billingPolicy.Status == "Enabled" || billingPolicy.Status == "Disabled" {
            return billingPolicy, nil
        }

        if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
            return nil, err
        }

        tflog.Debug(ctx, fmt.Sprintf("Billing Policy Operation State: '%s'", billingPolicy.Status))
    }
}
```

---

## Impact

Retrying without a maximum limit or timeout can cause the application to hang indefinitely, leading to resource wastage.

**Severity:** Critical

---

## Suggested Fix

Introduce a retry limit or configurable timeout to avoid infinite retries:

```go
func (client *Client) DoWaitForFinalStatus(ctx context.Context, billingPolicyDto *BillingPolicyDto) (*BillingPolicyDto, error) {
    billingId := billingPolicyDto.Id
    maxRetries := 10 // Configurable retry limit.
    retryCount := 0

    for {
        if retryCount >= maxRetries {
            return nil, fmt.Errorf("reached maximum retries while waiting for billing policy '%s' to reach terminal state", billingId)
        }

        billingPolicy, err := client.GetBillingPolicy(ctx, billingId)

        if err != nil {
            return nil, err
        }

        if billingPolicy.Status == "Enabled" || billingPolicy.Status == "Disabled" {
            return billingPolicy, nil
        }

        if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
            return nil, err
        }

        retryCount++
        tflog.Debug(ctx, fmt.Sprintf("Billing Policy Operation State: '%s' (Retry %d/%d)", billingPolicy.Status, retryCount, maxRetries))
    }
}
```

---

## Recommended Actions

- Ensure all long-running methods have a retry limit or timeout mechanism to avoid blocking resources indefinitely.
- Consider externalizing the `maxRetries` value to a configuration file or environment variable for flexibility.

---