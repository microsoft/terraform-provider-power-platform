// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant_isolation_policy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/tenant"
)

var _ resource.Resource = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}
var _ resource.ResourceWithValidateConfig = &Resource{}

// NewTenantIsolationPolicyResource creates a new tenant isolation policy resource.
func NewTenantIsolationPolicyResource() resource.Resource {
	return &Resource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "tenant_isolation_policy",
		},
	}
}

// Metadata returns the resource type name.
func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

// Schema defines the schema for the resource.
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Power Platform tenant isolation policy. Tenant isolation can be used to block external tenants " +
			"from establishing connections into your tenant (inbound isolation) as well as block your tenant from " +
			"establishing connections to external tenants (outbound isolation). " +
			"Learn more: https://docs.microsoft.com/en-us/power-platform/admin/cross-tenant-restrictions",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{}),
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the tenant isolation policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_disabled": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether the tenant isolation policy is disabled.",
			},
			"allowed_tenants": schema.SetNestedAttribute{
				Required:            true,
				MarkdownDescription: "List of tenants that are allowed to connect with your tenant.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tenant_id": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "ID of the tenant that is allowed to connect.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"inbound": schema.BoolAttribute{
							Required:            true,
							MarkdownDescription: "Whether inbound connections from this tenant are allowed.",
						},
						"outbound": schema.BoolAttribute{
							Required:            true,
							MarkdownDescription: "Whether outbound connections to this tenant are allowed.",
						},
					},
				},
			},
		},
	}
}

// Configure configures the resource with the provider client.
func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
		return
	}

	client, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	tenantClient := tenant.NewTenantClient(client.Api)
	r.Client = NewTenantIsolationPolicyClient(client.Api, tenantClient)
}

// ValidateConfig validates the resource configuration.
func (r *Resource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var data TenantIsolationPolicyResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that allowed tenants have at least one direction (inbound or outbound) enabled
	var modelTenants []AllowedTenantModel
	resp.Diagnostics.Append(data.AllowedTenants.ElementsAs(ctx, &modelTenants, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, allowedTenant := range modelTenants {
		if !allowedTenant.Inbound.ValueBool() && !allowedTenant.Outbound.ValueBool() {
			resp.Diagnostics.AddAttributeError(
				path.Root("allowed_tenants"),
				"Invalid tenant connection configuration",
				fmt.Sprintf("At least one of 'inbound' or 'outbound' must be true for tenant %s", allowedTenant.TenantId.ValueString()),
			)
		}
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan TenantIsolationPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current tenant ID
	tenantInfo, err := r.Client.TenantApi.GetTenant(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving tenant information",
			fmt.Sprintf("Could not retrieve tenant information: %s", err.Error()),
		)
		return
	}

	// Convert the Terraform model to the API model
	policyDto, diags := convertToDto(ctx, tenantInfo.TenantId, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the policy
	policy, err := r.Client.createOrUpdateTenantIsolationPolicy(ctx, tenantInfo.TenantId, *policyDto)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating tenant isolation policy",
			fmt.Sprintf("Could not create tenant isolation policy: %s", err.Error()),
		)
		return
	}

	// Convert the API response back to the Terraform model
	state, diags := convertFromDto(ctx, policy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Id = types.StringValue(tenantInfo.TenantId)
	state.Timeouts = plan.Timeouts

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Created tenant isolation policy with ID %s", state.Id.ValueString()))
}

// Read refreshes the Terraform state with the latest data.
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state TenantIsolationPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Store timeouts to preserve them
	timeoutsVal := state.Timeouts

	tenantId := state.Id.ValueString()
	if tenantId == "" {
		resp.Diagnostics.AddError(
			"Missing tenant ID",
			"The tenant ID is unexpectedly missing from state. This is a provider error.",
		)
		return
	}

	// Get the current policy
	policy, err := r.Client.getTenantIsolationPolicy(ctx, tenantId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading tenant isolation policy",
			fmt.Sprintf("Could not read tenant isolation policy: %s", err.Error()),
		)
		return
	}

	// If the policy doesn't exist, remove from state
	if policy == nil {
		// follows best practice of removing the resource from state
		// https://developer.hashicorp.com/terraform/plugin/framework/resources/read#recommendations
		tflog.Debug(ctx, fmt.Sprintf("Removing tenant isolation policy with ID %s from state because it no longer exists", tenantId))
		resp.State.RemoveResource(ctx)
		return
	}

	// Convert from DTO and preserve timeouts
	updatedState, diags := convertFromDto(ctx, policy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedState.Id = types.StringValue(tenantId)
	updatedState.Timeouts = timeoutsVal

	// Set the refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read tenant isolation policy with ID %s", updatedState.Id.ValueString()))
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan TenantIsolationPolicyResourceModel
	var state TenantIsolationPolicyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get timeout values from plan since it represents desired configuration
	timeoutsVal := plan.Timeouts

	tenantId := state.Id.ValueString()
	// We can be confident the ID exists in state for an update operation
	if tenantId == "" {
		resp.Diagnostics.AddError(
			"Missing tenant ID",
			"The tenant ID is unexpectedly missing from state. This is a provider error.",
		)
		return
	}

	// Convert the Terraform model to the API model
	policyDto, diags := convertToDto(ctx, tenantId, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the policy
	updatedPolicy, err := r.Client.createOrUpdateTenantIsolationPolicy(ctx, tenantId, *policyDto)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating tenant isolation policy",
			fmt.Sprintf("Could not update tenant isolation policy: %s", err.Error()),
		)
		return
	}

	// Convert from DTO and preserve timeouts
	updatedState, diags := convertFromDto(ctx, updatedPolicy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedState.Id = types.StringValue(tenantId)
	updatedState.Timeouts = timeoutsVal

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updated tenant isolation policy with ID %s", updatedState.Id.ValueString()))
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state TenantIsolationPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := state.Id.ValueString()
	if tenantId == "" {
		resp.Diagnostics.AddError(
			"Missing tenant ID",
			"The tenant ID is unexpectedly missing from state. This is likely a provider bug.",
		)
		return
	}

	// To delete the policy, we update it with an empty policy
	emptyPolicy := TenantIsolationPolicyDto{
		Properties: TenantIsolationPolicyPropertiesDto{
			TenantId:       tenantId,
			IsDisabled:     types.BoolValue(false).ValueBoolPointer(),
			AllowedTenants: []AllowedTenantDto{},
		},
	}

	_, err := r.Client.createOrUpdateTenantIsolationPolicy(ctx, tenantId, emptyPolicy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting tenant isolation policy",
			fmt.Sprintf("Could not delete tenant isolation policy: %s", err.Error()),
		)
		return
	}

	// Remove the resource from state
	resp.State.RemoveResource(ctx)

	tflog.Info(ctx, fmt.Sprintf("Deleted tenant isolation policy with ID %s", tenantId))
}

// ImportState imports the resource into Terraform from an existing tenant isolation policy.
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// The import ID can either be empty (use current tenant) or a specific tenant ID
	tenantId := req.ID
	if tenantId == "" {
		// Get the current tenant ID if not provided
		tenantInfo, err := r.Client.TenantApi.GetTenant(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error retrieving tenant information",
				fmt.Sprintf("Could not retrieve tenant information: %s", err.Error()),
			)
			return
		}
		tenantId = tenantInfo.TenantId
	}

	// Set the ID into the state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), tenantId)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Imported tenant isolation policy with ID %s", tenantId))
}

// Helper functions

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
