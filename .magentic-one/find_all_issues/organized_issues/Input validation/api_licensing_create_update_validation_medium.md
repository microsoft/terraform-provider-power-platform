# No Validation of Input Data in Create/Update Functions

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/api_licensing.go

## Problem

The functions `CreateBillingPolicy` and `UpdateBillingPolicy` do not validate `policyToCreate` and `policyToUpdate` input parameters, respectively. If a nil or zero-value struct is passed, the API call may result in an unexpected error.

## Impact

Medium. This could cause errors from the API or make debugging input-related problems more difficult.

## Location

```go
func (client *Client) CreateBillingPolicy(ctx context.Context, policyToCreate billingPolicyCreateDto) (*BillingPolicyDto, error) {
...
func (client *Client) UpdateBillingPolicy(ctx context.Context, billingId string, policyToUpdate BillingPolicyUpdateDto) (*BillingPolicyDto, error) {
```

## Fix

Add validation for required fields before making the API call. For example:

```go
if policyToCreate.Name == "" {
    return nil, fmt.Errorf("policy name is required")
}
```

(Similar validation for other required fields and update function.)
