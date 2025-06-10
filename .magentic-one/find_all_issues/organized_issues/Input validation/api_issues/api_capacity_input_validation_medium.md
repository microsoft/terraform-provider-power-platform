# Issue: Lack of Input Validation on Function Argument

## 
/workspaces/terraform-provider-power-platform/internal/services/capacity/api_capacity.go

## Problem

The `GetTenantCapacity` function does not validate its `tenantId` parameter before using it to create a URL path. If an invalid, empty, or malformed tenant ID is passed, the constructed URL could be invalid or could result in unexpected behavior. Input validation ensures early detection of incorrect usage and can prevent subtle bugs and security issues.

## Impact

Severity: **medium**

Allowing invalid input unchecked can lead to failed API requests, developer confusion, or even potential security vulnerabilities if URL paths can be manipulated.

## Location

```go
func (client *client) GetTenantCapacity(ctx context.Context, tenantId string) (*capacityDto, error) {
// ... no validation on tenantId
}
```

## Code Issue

```go
func (client *client) GetTenantCapacity(ctx context.Context, tenantId string) (*capacityDto, error) {
    apiUrl := &url.URL{
        Scheme: constants.HTTPS,
        Host:   client.Api.GetConfig().Urls.LicensingUrl,
        Path:   fmt.Sprintf("/v0.1-alpha/tenants/%s/TenantCapacity", tenantId),
    }
    // ...
}
```

## Fix

Add a check to ensure that `tenantId` is not empty and consider additional formatting/length validation if applicable:

```go
func (client *client) GetTenantCapacity(ctx context.Context, tenantId string) (*capacityDto, error) {
    if tenantId == "" {
        return nil, fmt.Errorf("tenantId cannot be empty")
    }
    apiUrl := &url.URL{
        Scheme: constants.HTTPS,
        Host:   client.Api.GetConfig().Urls.LicensingUrl,
        Path:   fmt.Sprintf("/v0.1-alpha/tenants/%s/TenantCapacity", tenantId),
    }
    // ...
}
```
