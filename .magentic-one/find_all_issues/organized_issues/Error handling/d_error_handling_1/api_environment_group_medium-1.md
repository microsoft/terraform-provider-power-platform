# Missing Error Wrapping with Context

##

/workspaces/terraform-provider-power-platform/internal/services/environment_groups/api_environment_group.go

## Problem

In several methods (for example, `RemoveEnvironmentFromEnvironmentGroup`, `CreateEnvironmentGroup`, `UpdateEnvironmentGroup`, `GetEnvironmentsInEnvironmentGroup`), errors from dependencies or API executions are simply returned, losing important context about the operation in which they occurred.

## Impact

This makes it much harder to debug or trace errors, especially when the same error (e.g., HTTP 500) can occur in multiple functions. It results in reduced maintainability and supportability.

**Severity:** Medium

## Location

```go
tenantDto, err := client.TenantApi.GetTenant(ctx)
if err != nil {
    return err
}
```
...and elsewhere (`return err` with no context wrapping).

## Fix

Wrap errors with additional context using `fmt.Errorf("context: %w", err)`:

```go
tenantDto, err := client.TenantApi.GetTenant(ctx)
if err != nil {
    return fmt.Errorf("failed to get tenant: %w", err)
}
```
Similarly, wrap all errors returned by API calls and downstream dependencies with a helpful message.
