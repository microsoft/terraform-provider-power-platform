// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant_isolation_policy

// TenantIsolationPolicyDto represents the API data transfer object for tenant isolation policy.
type TenantIsolationPolicyDto struct {
	Properties TenantIsolationPolicyPropertiesDto `json:"properties"`
}

// TenantIsolationPolicyPropertiesDto represents the properties of a tenant isolation policy.
type TenantIsolationPolicyPropertiesDto struct {
	TenantId       string             `json:"tenantId"`
	IsDisabled     *bool              `json:"isDisabled,omitempty"`
	AllowedTenants []AllowedTenantDto `json:"allowedTenants"`
}

// AllowedTenantDto represents an allowed tenant in the API response.
type AllowedTenantDto struct {
	TenantId  string       `json:"tenantId"`
	Direction DirectionDto `json:"direction"`
}

// DirectionDto represents the direction settings for an allowed tenant.
type DirectionDto struct {
	Inbound  *bool `json:"inbound"`
	Outbound *bool `json:"outbound"`
}
