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
