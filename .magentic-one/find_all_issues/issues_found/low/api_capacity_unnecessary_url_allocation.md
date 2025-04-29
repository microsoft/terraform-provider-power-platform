# Title

Unnecessary allocation of `&url.URL` struct in API URL construction

##

/workspaces/terraform-provider-power-platform/internal/services/capacity/api_capacity.go

## Problem

The usage of `&url.URL{}` to allocate a new struct during API URL construction is unnecessary. Using the literal form without the ampersand works fine because Go handles struct copying efficiently. This improves readability and avoids mistaken assumptions about pointer handling.

## Impact

This issue has a **low severity** as it does not impact program correctness or performance but contributes to less clear code.

## Location

Line(s) of code involving allocation:
```go
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.LicensingUrl,
		Path:   fmt.Sprintf("/v0.1-alpha/tenants/%s/TenantCapacity", tenantId),
	}
```

## Code Issue

```go
apiUrl := &url.URL{
	Scheme: constants.HTTPS,
	Host:   client.Api.GetConfig().Urls.LicensingUrl,
	Path:   fmt.Sprintf("/v0.1-alpha/tenants/%s/TenantCapacity", tenantId),
}
```

## Fix

Refactor and use struct literals directly without explicit pointer allocation:
```go
apiUrl := url.URL{
	Scheme: constants.HTTPS,
	Host:   client.Api.GetConfig().Urls.LicensingUrl,
	Path:   fmt.Sprintf("/v0.1-alpha/tenants/%s/TenantCapacity", tenantId),
}
```

Explanation:

- The `url.URL` struct is passed as a copy to the `String()` method anyway, and the use of a pointer here does not offer any additional benefits.
- This results in better readability without compromising functionality.