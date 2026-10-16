// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package licensing

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
	Meters            types.Object                   `tfsdk:"meters"`
}

type BillingInstrumentResourceModel struct {
	Id             types.String `tfsdk:"id"`
	ResourceGroup  types.String `tfsdk:"resource_group"`
	SubscriptionId types.String `tfsdk:"subscription_id"`
}

type BillingPolicyMetersResourceModel struct {
	Dataverse           types.Bool `tfsdk:"dataverse"`
	PowerApps           types.Bool `tfsdk:"power_apps"`
	PowerAutomate       types.Bool `tfsdk:"power_automate"`
	PowerPages          types.Bool `tfsdk:"power_pages"`
	CopilotStudio       types.Bool `tfsdk:"copilot_studio"`
	Windows365ForAgents types.Bool `tfsdk:"windows_365_for_agents"`
}

// billingPolicyMeters lists the pay-as-you-go meters by API product category.
// The API rejects a request unless all entitlements of a product category are sent together with the same state.
var billingPolicyMeters = []struct {
	productCategory string
	entitlementIds  []string
}{
	{productCategory: "Dataverse", entitlementIds: []string{"File", "Log", "Database"}},
	{productCategory: "PowerApps", entitlementIds: []string{"AppPass"}},
	{productCategory: "PowerAutomate", entitlementIds: []string{"CloudFlowRuns", "PAUnattendedRPA", "PAAttendedRPA"}},
	{productCategory: "PowerPages", entitlementIds: []string{"PowerPagesAuthenticated", "PowerPagesAnonymous"}},
	{productCategory: "PowerVirtualAgent", entitlementIds: []string{"MCSMessages"}},
	{productCategory: "W365", entitlementIds: []string{"W365APAYGO"}},
}

var billingPolicyMetersAttrTypes = map[string]attr.Type{
	"dataverse":              types.BoolType,
	"power_apps":             types.BoolType,
	"power_automate":         types.BoolType,
	"power_pages":            types.BoolType,
	"copilot_studio":         types.BoolType,
	"windows_365_for_agents": types.BoolType,
}

func (m *BillingPolicyMetersResourceModel) byProductCategory() map[string]*types.Bool {
	return map[string]*types.Bool{
		"Dataverse":         &m.Dataverse,
		"PowerApps":         &m.PowerApps,
		"PowerAutomate":     &m.PowerAutomate,
		"PowerPages":        &m.PowerPages,
		"PowerVirtualAgent": &m.CopilotStudio,
		"W365":              &m.Windows365ForAgents,
	}
}

func productCategoryOfEntitlement(entitlement PayGoEntitlementDto) string {
	for _, meter := range billingPolicyMeters {
		if meter.productCategory == string(entitlement.ProductCategory) {
			return meter.productCategory
		}
	}
	for _, meter := range billingPolicyMeters {
		for _, id := range meter.entitlementIds {
			if id == entitlement.EntitlementId {
				return meter.productCategory
			}
		}
	}
	return string(entitlement.ProductCategory)
}

func metersFromObject(ctx context.Context, meters types.Object) (BillingPolicyMetersResourceModel, diag.Diagnostics) {
	model := BillingPolicyMetersResourceModel{}
	diags := meters.As(ctx, &model, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
	return model, diags
}

// convertPayGoEntitlementsDtoToMetersModel reports a meter as enabled when its entitlements are enabled.
// Meters that the API does not return are reported as disabled.
func convertPayGoEntitlementsDtoToMetersModel(ctx context.Context, entitlements []PayGoEntitlementDto) (types.Object, diag.Diagnostics) {
	model := BillingPolicyMetersResourceModel{}
	states := model.byProductCategory()
	for _, state := range states {
		*state = types.BoolValue(false)
	}
	for _, entitlement := range entitlements {
		if state, ok := states[productCategoryOfEntitlement(entitlement)]; ok && entitlement.PayAsYouGoState {
			*state = types.BoolValue(true)
		}
	}
	return types.ObjectValueFrom(ctx, billingPolicyMetersAttrTypes, model)
}

// convertMetersModelToPayGoEntitlementsDto builds the entitlements to send to the API.
// When current entitlements are provided (update), the full current list is returned and only meters with a known
// value in the model are changed, because the API disables every meter that is missing from an update request.
func convertMetersModelToPayGoEntitlementsDto(model BillingPolicyMetersResourceModel, current []PayGoEntitlementDto) []PayGoEntitlementDto {
	desired := model.byProductCategory()
	entitlements := []PayGoEntitlementDto{}
	handled := map[string]bool{}

	for _, entitlement := range current {
		category := productCategoryOfEntitlement(entitlement)
		state := entitlement.PayAsYouGoState
		if value, ok := desired[category]; ok && !value.IsNull() && !value.IsUnknown() {
			state = value.ValueBool()
		}
		entitlements = append(entitlements, PayGoEntitlementDto{
			EntitlementId:   entitlement.EntitlementId,
			ProductCategory: PayGoProductCategory(category),
			PayAsYouGoState: state,
		})
		handled[category] = true
	}

	for _, meter := range billingPolicyMeters {
		value := desired[meter.productCategory]
		if handled[meter.productCategory] || value.IsNull() || value.IsUnknown() {
			continue
		}
		for _, id := range meter.entitlementIds {
			entitlements = append(entitlements, PayGoEntitlementDto{
				EntitlementId:   id,
				ProductCategory: PayGoProductCategory(meter.productCategory),
				PayAsYouGoState: value.ValueBool(),
			})
		}
	}

	return entitlements
}
