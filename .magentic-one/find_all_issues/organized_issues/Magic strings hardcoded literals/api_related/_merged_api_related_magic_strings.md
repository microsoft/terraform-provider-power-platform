# Magic Strings and Hardcoded Literals - API Related Issues

This document consolidates all magic strings and hardcoded literals issues found in API-related files.


## ISSUE 1

# Issue: Magic constants used for retry/backoff without documentation or configuration

##

/workspaces/terraform-provider-power-platform/internal/services/environment/api_environment.go

## Problem

Several places use hardcoded retry timings via `api.DefaultRetryAfter()` and constant strings for provisioning state (e.g., `"Succeeded"`, `"Failed"`, `"LinkedDatabaseProvisioning"`, etc) and fixed HTTP status code lists (e.g., `[]int{http.StatusOK, http.StatusAccepted, http.StatusConflict}`). These are magic values and would benefit from documentation, centralization, or configuration for easier adjustments and maintainability.

## Impact

- Severity: Low
- Makes maintenance and debugging harder if a retry interval or provisioning logic changes.
- Hard to find widespread usages for tuning.

## Location

Examples throughout the file, e.g.,

```go
if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
    return err
}
```

And provisioning state checks:

```go
if lifecycleEnv.Properties.ProvisioningState == "Succeeded" {
    // ...
} else if lifecycleEnv.Properties.ProvisioningState != "LinkedDatabaseProvisioning" && lifecycleEnv.Properties.ProvisioningState != "Succeeded" {
    // ...
}
```

## Code Issue

```go
if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
    return err
}
```

Or provisioning states:

```go
if lifecycleEnv.Properties.ProvisioningState == "Succeeded" {
    // ...
} else if lifecycleEnv.Properties.ProvisioningState != "LinkedDatabaseProvisioning" && lifecycleEnv.Properties.ProvisioningState != "Succeeded" {
    // ...
}
```

## Fix

Document the rationale for backoff/retry values and magic strings; consider a centralized definition:

```go
const (
    ProvisioningStateSucceeded = "Succeeded"
    ProvisioningStateFailed = "Failed"
    ProvisioningStateLinkedDatabaseProvisioning = "LinkedDatabaseProvisioning"
    // ...more as needed
)
```
And for retry intervals:

```go
const DefaultRetryAfter = 15 * time.Second  // document why!
```
Or, link documentation directly above with context for the value.

---

This issue should be saved under:  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/api_environment_magic_constants_low.md`

---

## ISSUE 2

# Issue: Hardcoding of Wait Duration Bounds in getRetryAfterDuration

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/api_tenant_isolation_policy.go

## Problem

The `getRetryAfterDuration` function hardcodes the minimum (2 seconds) and maximum (60 seconds) wait duration for polling retries directly in the implementation. Hardcoded magic numbers reduce flexibility, can violate project constants standards, and make adjustment harder if retry policy needs tuning. These bounds are duplicated here instead of being documented or externally configured.

## Impact

Low. While it does not directly cause incorrect results, it makes future adjustments riskier and less transparent, and reduces maintainability.

## Location

Within `getRetryAfterDuration`:

```go
if duration < 2*time.Second {
    duration = 2 * time.Second
} else if duration > 60*time.Second {
    duration = 60 * time.Second
}
```

## Code Issue

```go
if duration < 2*time.Second {
    duration = 2 * time.Second
} else if duration > 60*time.Second {
    duration = 60 * time.Second
}
```

## Fix

Define constants at the top of the file or in a shared package, and reference these in the function for clarity and centralized control. Example:

```go
const (
    MinRetryAfterDuration = 2 * time.Second
    MaxRetryAfterDuration = 60 * time.Second
)
```

Then:

```go
if duration < MinRetryAfterDuration {
    duration = MinRetryAfterDuration
} else if duration > MaxRetryAfterDuration {
    duration = MaxRetryAfterDuration
}
```
---

---

## ISSUE 3

# Issue: Repeated Logic in API URL Construction

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/api_tenant_isolation_policy.go

## Problem

Both `getTenantIsolationPolicy` and `createOrUpdateTenantIsolationPolicy` construct the `apiUrl` in an almost identical way, using the same hardcoded path pattern. This duplicated logic makes it more likely for inconsistencies and maintenance challenges to arise. If the URL structure ever changes, updates would need to be made in multiple places, increasing risk of errors.

## Impact

Low to Medium. The impact is primarily on maintainability and code clarity. While not immediately breaking, this code duplication violates DRY (Don't Repeat Yourself) principles and creates unnecessary maintenance burden.

## Location

Multiple locations—repeated construction of:

```go
apiUrl := &url.URL{
	Scheme: constants.HTTPS,
	Host:   client.Api.GetConfig().Urls.BapiUrl,
	Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", tenantId),
}
```

## Code Issue

```go
apiUrl := &url.URL{
	Scheme: constants.HTTPS,
	Host:   client.Api.GetConfig().Urls.BapiUrl,
	Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", tenantId),
}
```

## Fix

Extract this URL construction logic into a helper method on the `Client` type to centralize the logic and support future changes in a single place:

```go
func (client *Client) getTenantIsolationPolicyURL(tenantId string) string {
	return (&url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", tenantId),
	}).String()
}
```

Then call this helper method in both locations:

```go
apiUrl := client.getTenantIsolationPolicyURL(tenantId)
```

---


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
- `<description>` is a clear explanation of what was fixed/changed search for "copilot-commit-message-instructions.md" how to write description.
- `<issue_number>` pick the issue number or PR number
```
