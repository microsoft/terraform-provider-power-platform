// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant_isolation_policy

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
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

// Resource represents the tenant isolation policy resource.
type Resource struct {
	helpers.TypeInfo
	Client Client
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
		Description: "Manages a Power Platform tenant isolation policy. Tenant isolation can be used to block external tenants " +
			"from establishing connections into your tenant (inbound isolation) as well as block your tenant from " +
			"establishing connections to external tenants (outbound isolation). " +
			"Learn more: https://docs.microsoft.com/en-us/power-platform/admin/cross-tenant-restrictions",

		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{}),
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the tenant isolation policy.",
				MarkdownDescription: "The ID of the tenant isolation policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_disabled": schema.BoolAttribute{
				Required:            true,
				Description:         "Whether the tenant isolation policy is disabled.",
				MarkdownDescription: "Whether the tenant isolation policy is disabled.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_tenants": schema.SetNestedAttribute{
				Required:            true,
				Description:         "List of tenants that are allowed to connect with your tenant.",
				MarkdownDescription: "List of tenants that are allowed to connect with your tenant.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tenant_id": schema.StringAttribute{
							Required:            true,
							Description:         "ID of the tenant that is allowed to connect.",
							MarkdownDescription: "ID of the tenant that is allowed to connect.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"inbound": schema.BoolAttribute{
							Required:            true,
							Description:         "Whether inbound connections from this tenant are allowed.",
							MarkdownDescription: "Whether inbound connections from this tenant are allowed.",
						},
						"outbound": schema.BoolAttribute{
							Required:            true,
							Description:         "Whether outbound connections to this tenant are allowed.",
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
	for _, allowedTenant := range data.AllowedTenants {
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
	policyDto := convertToDto(tenantInfo.TenantId, &plan)

	// Create the policy
	policy, err := r.Client.CreateOrUpdateTenantIsolationPolicy(ctx, tenantInfo.TenantId, *policyDto)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating tenant isolation policy",
			fmt.Sprintf("Could not create tenant isolation policy: %s", err.Error()),
		)
		return
	}

	// Convert the API response back to the Terraform model
	state := convertFromDto(policy)
	state.Id = types.StringValue(tenantInfo.TenantId)
	state.Timeouts = plan.Timeouts

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Created tenant isolation policy with ID %s", state.Id.ValueString()))
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
		// Get the current tenant ID if not set
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

	// Get the current policy
	policy, err := r.Client.GetTenantIsolationPolicy(ctx, tenantId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading tenant isolation policy",
			fmt.Sprintf("Could not read tenant isolation policy: %s", err.Error()),
		)
		return
	}

	// If the policy doesn't exist, remove from state
	if policy == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Convert from DTO and preserve timeouts
	updatedState := convertFromDto(policy)
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
	if tenantId == "" {
		// Get the current tenant ID if not set
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

	// Convert the Terraform model to the API model
	policyDto := convertToDto(tenantId, &plan)

	// Update the policy
	updatedPolicy, err := r.Client.CreateOrUpdateTenantIsolationPolicy(ctx, tenantId, *policyDto)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating tenant isolation policy",
			fmt.Sprintf("Could not update tenant isolation policy: %s", err.Error()),
		)
		return
	}

	// Convert from DTO and preserve timeouts
	updatedState := convertFromDto(updatedPolicy)
	updatedState.Id = types.StringValue(tenantId)
	updatedState.Timeouts = timeoutsVal

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updated tenant isolation policy with ID %s", updatedState.Id.ValueString()))
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
		// Get the current tenant ID if not set
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

	// To delete the policy, we update it with an empty policy
	emptyPolicy := TenantIsolationPolicyDto{
		Properties: TenantIsolationPolicyPropertiesDto{
			TenantId:       state.Id.ValueString(),
			AllowedTenants: []AllowedTenantDto{},
		},
	}

	_, err := r.Client.CreateOrUpdateTenantIsolationPolicy(ctx, tenantId, emptyPolicy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting tenant isolation policy",
			fmt.Sprintf("Could not delete tenant isolation policy: %s", err.Error()),
		)
		return
	}

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
func convertToDto(tenantId string, model *TenantIsolationPolicyResourceModel) *TenantIsolationPolicyDto {
	return &TenantIsolationPolicyDto{
		Properties: TenantIsolationPolicyPropertiesDto{
			TenantId:       tenantId,
			IsDisabled:     model.IsDisabled.ValueBoolPointer(),
			AllowedTenants: convertAllowedTenantsToDto(model.AllowedTenants),
		},
	}
}

func convertAllowedTenantsToDto(modelTenants []AllowedTenantModel) []AllowedTenantDto {
	dtoTenants := make([]AllowedTenantDto, 0, len(modelTenants))
	for _, modelTenant := range modelTenants {
		inbound := modelTenant.Inbound.ValueBool()
		outbound := modelTenant.Outbound.ValueBool()

		dtoTenants = append(dtoTenants, AllowedTenantDto{
			TenantId: modelTenant.TenantId.ValueString(),
			Direction: DirectionDto{
				Inbound:  &inbound,
				Outbound: &outbound,
			},
		})
	}
	return dtoTenants
}

// convertFromDto converts the API DTO to the Terraform model.
func convertFromDto(dto *TenantIsolationPolicyDto) TenantIsolationPolicyResourceModel {
	if dto == nil {
		return TenantIsolationPolicyResourceModel{}
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

	return TenantIsolationPolicyResourceModel{
		Id:             types.StringValue(tenantId),
		IsDisabled:     types.BoolValue(isDisabled != nil && *isDisabled),
		AllowedTenants: convertAllowedTenantsFromDto(allowedTenants),
	}
}

func convertAllowedTenantsFromDto(dtoTenants []AllowedTenantDto) []AllowedTenantModel {
	if dtoTenants == nil {
		return []AllowedTenantModel{}
	}

	// Sort tenants by ID for consistent ordering
	sort.Slice(dtoTenants, func(i, j int) bool {
		return dtoTenants[i].TenantId < dtoTenants[j].TenantId
	})

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

		modelTenants = append(modelTenants, AllowedTenantModel{
			TenantId: types.StringValue(dtoTenant.TenantId),
			Inbound:  types.BoolValue(inbound),
			Outbound: types.BoolValue(outbound),
		})
	}
	return modelTenants
}
