# Error Logging Without Propagation - Model Type Safety

This document consolidates all issues related to error logging without proper propagation found in model type safety implementations across the Terraform Provider for Power Platform.

## ISSUE 1

# Error Handling: Potential Loss of Context in `convertCreateEnvironmentDtoFromSourceModel`

##

/workspaces/terraform-provider-power-platform/internal/services/environment/models.go

## Problem

In `convertCreateEnvironmentDtoFromSourceModel`, if retrieving the tenant ID from `r.EnvironmentClient.tenantClient.GetTenant(ctx)` fails, the error is returned directly without additional context or logging, making it harder to diagnose the original call's context and origin during debugging.

## Impact

- **Severity:** Medium
- Makes tracing errors difficult, especially in complex workflows or remote API calls.
- May frustrate operators/users since log output may be unclear.
- May prevent proper root-cause diagnosis if similar errors occur in related library functions.

## Location

```go
if !environmentSource.OwnerId.IsNull() && !environmentSource.OwnerId.IsUnknown() {
 tenantId, err := r.EnvironmentClient.tenantClient.GetTenant(ctx)
 if err != nil {
  return nil, err
 }
 environmentDto.Properties.UsedBy = &UsedByDto{
  Id:       environmentSource.OwnerId.ValueString(),
  Type:     "1",
  TenantId: tenantId.TenantId,
 }
}
```

## Code Issue

```go
tenantId, err := r.EnvironmentClient.tenantClient.GetTenant(ctx)
if err != nil {
 return nil, err
}
```

## Fix

Wrap errors with context for better traceability, using `fmt.Errorf` or `%w`:

```go
tenantId, err := r.EnvironmentClient.tenantClient.GetTenant(ctx)
if err != nil {
 return nil, fmt.Errorf("failed to retrieve tenant for OwnerId %s: %w", environmentSource.OwnerId.ValueString(), err)
}
```

This will help operators quickly understand what failed if this surface-level call errors.

# To finish the task you have to

1. Run linter and fix any issues
2. Run UnitTest and fix any of failing ones
3. Generate docs
4. Run Changie

# Changie Instructions

Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```

Where:

- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
