// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package licensing

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type BillingPoliciesEnvironmetsDataSource struct {
	helpers.TypeInfo
	LicensingClient Client
}

type BillingPoliciesEnvironmetsListDataSourceModel struct {
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
	BillingPolicyId string         `tfsdk:"billing_policy_id"`
	Environments    []string       `tfsdk:"environments"`
}

type BillingPoliciesDataSource struct {
	helpers.TypeInfo
	LicensingClient Client
}

type BillingPoliciesListDataSourceModel struct {
	Timeouts        timeouts.Value                 `tfsdk:"timeouts"`
	BillingPolicies []BillingPolicyDataSourceModel `tfsdk:"billing_policies"`
}

type BillingPolicyDataSourceModel struct {
	Id                types.String                     `tfsdk:"id"`
	Name              types.String                     `tfsdk:"name"`
	Location          types.String                     `tfsdk:"location"`
	Status            types.String                     `tfsdk:"status"`
	BillingInstrument BillingInstrumentDataSourceModel `tfsdk:"billing_instrument"`
}

type BillingInstrumentDataSourceModel struct {
	Id             types.String `tfsdk:"id"`
	ResourceGroup  types.String `tfsdk:"resource_group"`
	SubscriptionId types.String `tfsdk:"subscription_id"`
}

type BillingPolicyEnvironmentResource struct {
	helpers.TypeInfo
	LicensingClient Client
}

type BillingPolicyEnvironmentResourceModel struct {
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
	BillingPolicyId string         `tfsdk:"billing_policy_id"`
	Environments    []string       `tfsdk:"environments"`
}

type BillingPolicyResource struct {
	helpers.TypeInfo
	LicensingClient Client
}

type BillingPolicyResourceModel struct {
	Timeouts          timeouts.Value                 `tfsdk:"timeouts"`
	Id                types.String                   `tfsdk:"id"`
	Name              types.String                   `tfsdk:"name"`
	Location          types.String                   `tfsdk:"location"`
	Status            types.String                   `tfsdk:"status"`
	BillingInstrument BillingInstrumentResourceModel `tfsdk:"billing_instrument"`
}

type BillingInstrumentResourceModel struct {
	Id             types.String `tfsdk:"id"`
	ResourceGroup  types.String `tfsdk:"resource_group"`
	SubscriptionId types.String `tfsdk:"subscription_id"`
}
