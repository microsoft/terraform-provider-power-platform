// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant_isolation_policy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

// Resource represents the tenant isolation policy resource.
type Resource struct {
	helpers.TypeInfo
	Client Client
}

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

// convertToDto converts the Terraform model to the API DTO.
func convertToDto(ctx context.Context, tenantId string, model *TenantIsolationPolicyResourceModel) (*TenantIsolationPolicyDto, diag.Diagnostics) {
	var diags diag.Diagnostics
	var tenantsModel []AllowedTenantModel
	diags.Append(model.AllowedTenants.ElementsAs(ctx, &tenantsModel, false)...)
	if diags.HasError() {
		return nil, diags
	}

	// Convert AllowedTenants to DTO
	dtoTenants := make([]AllowedTenantDto, 0, len(tenantsModel))
	for _, allowedTenant := range tenantsModel {
		inbound := allowedTenant.Inbound.ValueBool()
		outbound := allowedTenant.Outbound.ValueBool()
		dtoTenants = append(dtoTenants, AllowedTenantDto{
			TenantId: allowedTenant.TenantId.ValueString(),
			Direction: DirectionDto{
				Inbound:  &inbound,
				Outbound: &outbound,
			},
		})
	}

	return &TenantIsolationPolicyDto{
		Properties: TenantIsolationPolicyPropertiesDto{
			TenantId:       tenantId,
			IsDisabled:     model.IsDisabled.ValueBoolPointer(),
			AllowedTenants: dtoTenants,
		},
	}, diags
}

// convertFromDto converts the API DTO to the Terraform model.
func convertFromDto(ctx context.Context, dto *TenantIsolationPolicyDto) (TenantIsolationPolicyResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	if dto == nil {
		return TenantIsolationPolicyResourceModel{}, diags
	}

	// Set defaults in case Properties is nil
	tenantId := ""
	var isDisabled *bool
	var allowedTenants []AllowedTenantDto

	if dto.Properties.TenantId != "" {
		tenantId = dto.Properties.TenantId
	}
	if dto.Properties.IsDisabled != nil {
		isDisabled = dto.Properties.IsDisabled
	}
	if dto.Properties.AllowedTenants != nil {
		allowedTenants = dto.Properties.AllowedTenants
	}

	// Convert AllowedTenants to model objects
	modelTenants := convertAllowedTenantsFromDto(allowedTenants)

	// Create allowed_tenants as types.Set with explicit value mapping
	var allowedTenantsSet types.Set

	// If we have tenants, create the set with values
	if len(modelTenants) > 0 {
		// Convert model tenants to object values
		elements := make([]attr.Value, 0, len(modelTenants))
		for _, modelTenant := range modelTenants {
			elements = append(elements, createObjectFromAllowedTenant(modelTenant))
		}

		// Create set with proper element type and values
		var setDiags diag.Diagnostics
		allowedTenantsSet, setDiags = types.SetValue(allowedTenantsObjectType(ctx), elements)
		diags.Append(setDiags...)
		if diags.HasError() {
			return TenantIsolationPolicyResourceModel{}, diags
		}
	} else {
		// Create empty set with proper element type
		var setDiags diag.Diagnostics
		allowedTenantsSet, setDiags = types.SetValue(allowedTenantsObjectType(ctx), []attr.Value{})
		diags.Append(setDiags...)
		if diags.HasError() {
			return TenantIsolationPolicyResourceModel{}, diags
		}
	}

	return TenantIsolationPolicyResourceModel{
		Id:             types.StringValue(tenantId),
		IsDisabled:     types.BoolValue(isDisabled != nil && *isDisabled),
		AllowedTenants: allowedTenantsSet,
	}, diags
}

// createObjectFromAllowedTenant creates an object value from an AllowedTenantModel.
func createObjectFromAllowedTenant(allowedTenant AllowedTenantModel) attr.Value {
	return types.ObjectValueMust(
		map[string]attr.Type{
			"tenant_id": types.StringType,
			"inbound":   types.BoolType,
			"outbound":  types.BoolType,
		},
		map[string]attr.Value{
			"tenant_id": allowedTenant.TenantId,
			"inbound":   allowedTenant.Inbound,
			"outbound":  allowedTenant.Outbound,
		},
	)
}

// convertAllowedTenantsFromDto converts the API DTO to the Terraform model.
func convertAllowedTenantsFromDto(dtoTenants []AllowedTenantDto) []AllowedTenantModel {
	if dtoTenants == nil {
		return []AllowedTenantModel{}
	}

	modelTenants := make([]AllowedTenantModel, 0, len(dtoTenants))
	for _, dtoTenant := range dtoTenants {
		inbound := false
		outbound := false

		if dtoTenant.Direction.Inbound != nil {
			inbound = *dtoTenant.Direction.Inbound
		}
		if dtoTenant.Direction.Outbound != nil {
			outbound = *dtoTenant.Direction.Outbound
		}

		// Skip tenants with empty IDs
		if dtoTenant.TenantId == "" {
			continue
		}

		// Create a consistent model from the DTO with all fields explicitly set
		modelTenants = append(modelTenants, AllowedTenantModel{
			TenantId: types.StringValue(dtoTenant.TenantId),
			Inbound:  types.BoolValue(inbound),
			Outbound: types.BoolValue(outbound),
		})
	}
	return modelTenants
}

// Define the object type for AllowedTenants set.
func allowedTenantsObjectType(ctx context.Context) types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"tenant_id": types.StringType,
			"inbound":   types.BoolType,
			"outbound":  types.BoolType,
		},
	}
}
