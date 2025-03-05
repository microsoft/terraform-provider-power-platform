// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant_isolation_policy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Resource model for tenant isolation policy.
type TenantIsolationPolicyResourceModel struct {
	Timeouts       timeouts.Value `tfsdk:"timeouts"`
	Id             types.String   `tfsdk:"id"`
	IsDisabled     types.Bool     `tfsdk:"is_disabled"`
	AllowedTenants types.Set      `tfsdk:"allowed_tenants"`
}

type AllowedTenantModel struct {
	TenantId types.String `tfsdk:"tenant_id"`
	Inbound  types.Bool   `tfsdk:"inbound"`
	Outbound types.Bool   `tfsdk:"outbound"`
}
