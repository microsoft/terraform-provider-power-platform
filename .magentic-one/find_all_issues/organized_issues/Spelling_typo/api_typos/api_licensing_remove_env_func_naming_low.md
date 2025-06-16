# Misspelled Function Name for Removing Environments

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/api_licensing.go

## Problem

There is a typo in the function name: `RemoveEnvironmentsToBillingPolicy`. This should instead be `RemoveEnvironmentsFromBillingPolicy` to correctly reflect the operation semantics.

## Impact

Low. This is a semantic and clarity problem, but may cause confusion for developers and users of the API.

## Location

```go
func (client *Client) RemoveEnvironmentsToBillingPolicy(ctx context.Context, billingId string, environmentIds []string) error {
```

## Fix

Correct the method name and any references:

```go
func (client *Client) RemoveEnvironmentsFromBillingPolicy(ctx context.Context, billingId string, environmentIds []string) error {
```
