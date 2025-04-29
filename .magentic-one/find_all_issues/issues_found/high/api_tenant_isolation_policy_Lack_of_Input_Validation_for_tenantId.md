# Title

Lack of Input Validation for `tenantId`

##

`/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/api_tenant_isolation_policy.go`

## Problem

The `tenantId` parameter used in both `getTenantIsolationPolicy` and `createOrUpdateTenantIsolationPolicy` functions is directly interpolated into the API URL without any validation. This can lead to potential security vulnerabilities such as injection attacks or malformed requests.

## Impact

This issue can lead to unexpected behavior, including but not limited to:
- Injection attacks where malicious input exploits the URL structure.
- Malformed requests causing unnecessary API traffic and potential system downtime.
  
Severity: **High**

## Location

```go
func (client *Client) getTenantIsolationPolicy(ctx context.Context, tenantId string) (*TenantIsolationPolicyDto, error) {
    apiUrl := &url.URL{
        Scheme: constants.HTTPS,
        Host:   client.Api.GetConfig().Urls.BapiUrl,
        Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", tenantId),
    }

    ...
}
```

```go
func (client *Client) createOrUpdateTenantIsolationPolicy(ctx context.Context, tenantId string, policy TenantIsolationPolicyDto) (*TenantIsolationPolicyDto, error) {
    apiUrl := &url.URL{
        Scheme: constants.HTTPS,
        Host:   client.Api.GetConfig().Urls.BapiUrl,
        Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", tenantId),
    }

    ...
}
```

## Code Issue

```go
Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", tenantId),
```

## Fix

Introduce a helper function to validate `tenantId` format before using it in API requests. This ensures no invalid or malicious values are interpolated.

```go
func validateTenantID(tenantId string) error {
    if tenantId == "" {
        return errors.New("tenantId cannot be empty")
    }

    if len(tenantId) < 5 || len(tenantId) > 50 { // Example constraint
        return errors.New("tenantId must be between 5 and 50 characters")
    }

    if !regexp.MustCompile(`^[a-zA-Z0-9\-]+$`).MatchString(tenantId) {
        return errors.New("tenantId contains invalid characters")
    }

    return nil
}

func (client *Client) getTenantIsolationPolicy(ctx context.Context, tenantId string) (*TenantIsolationPolicyDto, error) {
    if err := validateTenantID(tenantId); err != nil {
        return nil, err
    }

    apiUrl := &url.URL{
        Scheme: constants.HTTPS,
        Host:   client.Api.GetConfig().Urls.BapiUrl,
        Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", tenantId),
    }

    ...
}
```