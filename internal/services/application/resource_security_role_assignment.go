// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package application

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ resource.Resource = &SecurityRoleAssignmentResource{}
var _ resource.ResourceWithImportState = &SecurityRoleAssignmentResource{}

type SecurityRoleAssignmentResource struct {
	helpers.TypeInfo
	ApplicationClient client
}

type SecurityRoleAssignmentResourceModel struct {
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
	Id               types.String   `tfsdk:"id"`
	EnvironmentId    types.String   `tfsdk:"environment_id"`
	SystemUserId     types.String   `tfsdk:"system_user_id"`
	BusinessUnitId   types.String   `tfsdk:"business_unit_id"`
	SecurityRoleName types.String   `tfsdk:"security_role_name"`
	RoleId           types.String   `tfsdk:"role_id"`
}

func NewSecurityRoleAssignmentResource() resource.Resource {
	return &SecurityRoleAssignmentResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "security_role_assignment",
		},
	}
}

func (r *SecurityRoleAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *SecurityRoleAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Assigns a single Dataverse security role, resolved by name within the target business unit, to a principal (system user). The assignment is managed independently of the principal's lifecycle.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Composite ID `{environment_id}/{system_user_id}/{security_role_name}`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Dataverse environment ID.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"system_user_id": schema.StringAttribute{
				MarkdownDescription: "Dataverse `systemuserid` of the user or application user the security role is assigned to. This is a Dataverse row id, not a Microsoft Entra object id, and `powerplatform_application_user` exposes it as `system_user_id`.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"business_unit_id": schema.StringAttribute{
				MarkdownDescription: "Business unit ID used to resolve the requested security role name. Defaults to the principal's current business unit.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"security_role_name": schema.StringAttribute{
				MarkdownDescription: "Dataverse security role name to assign.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role_id": schema.StringAttribute{
				MarkdownDescription: "Resolved Dataverse role ID for the assigned security role.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *SecurityRoleAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.ApplicationClient = newApplicationClient(client.Api)
}

func (r *SecurityRoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan SecurityRoleAssignmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resolved, err := r.resolveRequestedRole(
		ctx,
		plan.EnvironmentId.ValueString(),
		plan.SystemUserId.ValueString(),
		plan.BusinessUnitId.ValueString(),
		plan.SecurityRoleName.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to assign security role '%s' to principal '%s'", plan.SecurityRoleName.ValueString(), plan.SystemUserId.ValueString()),
			err.Error(),
		)
		return
	}

	if !principalHasRole(resolved.principal, resolved.role.RoleId) {
		resolved.principal, err = r.ApplicationClient.AddPrincipalSecurityRoles(ctx, plan.EnvironmentId.ValueString(), resolved.principal.SystemUserId, []string{resolved.role.RoleId})
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to assign security role '%s' to principal '%s'", plan.SecurityRoleName.ValueString(), plan.SystemUserId.ValueString()),
				err.Error(),
			)
			return
		}
	}

	plan.Id = types.StringValue(fmt.Sprintf("%s/%s/%s", plan.EnvironmentId.ValueString(), plan.SystemUserId.ValueString(), plan.SecurityRoleName.ValueString()))
	plan.BusinessUnitId = types.StringValue(resolved.businessUnitID)
	plan.RoleId = types.StringValue(resolved.role.RoleId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SecurityRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state SecurityRoleAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resolved, err := r.resolveRequestedRole(
		ctx,
		state.EnvironmentId.ValueString(),
		state.SystemUserId.ValueString(),
		state.BusinessUnitId.ValueString(),
		state.SecurityRoleName.ValueString(),
	)
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to read security role assignment '%s' for principal '%s'", state.SecurityRoleName.ValueString(), state.SystemUserId.ValueString()),
			err.Error(),
		)
		return
	}

	if !principalHasRole(resolved.principal, resolved.role.RoleId) {
		resp.State.RemoveResource(ctx)
		return
	}

	state.BusinessUnitId = types.StringValue(resolved.businessUnitID)
	state.RoleId = types.StringValue(resolved.role.RoleId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SecurityRoleAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan SecurityRoleAssignmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resolved, err := r.resolveRequestedRole(
		ctx,
		plan.EnvironmentId.ValueString(),
		plan.SystemUserId.ValueString(),
		plan.BusinessUnitId.ValueString(),
		plan.SecurityRoleName.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to refresh security role assignment '%s' for principal '%s'", plan.SecurityRoleName.ValueString(), plan.SystemUserId.ValueString()),
			err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(fmt.Sprintf("%s/%s/%s", plan.EnvironmentId.ValueString(), plan.SystemUserId.ValueString(), plan.SecurityRoleName.ValueString()))
	plan.BusinessUnitId = types.StringValue(resolved.businessUnitID)
	plan.RoleId = types.StringValue(resolved.role.RoleId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SecurityRoleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state SecurityRoleAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resolved, err := r.resolveRequestedRole(
		ctx,
		state.EnvironmentId.ValueString(),
		state.SystemUserId.ValueString(),
		state.BusinessUnitId.ValueString(),
		state.SecurityRoleName.ValueString(),
	)
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to read security role assignment '%s' for principal '%s'", state.SecurityRoleName.ValueString(), state.SystemUserId.ValueString()),
			err.Error(),
		)
		return
	}

	if !principalHasRole(resolved.principal, resolved.role.RoleId) {
		return
	}

	if _, err = r.ApplicationClient.RemovePrincipalSecurityRoles(ctx, state.EnvironmentId.ValueString(), resolved.principal.SystemUserId, []string{resolved.role.RoleId}); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to remove security role '%s' from principal '%s'", state.SecurityRoleName.ValueString(), state.SystemUserId.ValueString()),
			err.Error(),
		)
	}
}

func (r *SecurityRoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	idParts := strings.SplitN(req.ID, "/", 3)
	if len(idParts) != 3 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in format 'environment_id/system_user_id/security_role_name', got '%s'", req.ID),
		)
		return
	}

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("system_user_id"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("security_role_name"), idParts[2])...)
}

type resolvedRoleAssignment struct {
	principal      *applicationUserDto
	role           applicationSecurityRoleDto
	businessUnitID string
}

func (r *SecurityRoleAssignmentResource) resolveRequestedRole(ctx context.Context, environmentID, principalID, requestedBusinessUnitID, securityRoleName string) (*resolvedRoleAssignment, error) {
	dvExists, err := r.ApplicationClient.DataverseExists(ctx, environmentID)
	if err != nil {
		return nil, err
	}
	if !dvExists {
		return nil, fmt.Errorf("environment '%s' does not have Dataverse", environmentID)
	}

	currentPrincipal, err := r.ApplicationClient.GetPrincipalBySystemUserId(ctx, environmentID, principalID)
	if err != nil {
		return nil, err
	}

	businessUnitID := requestedBusinessUnitID
	if businessUnitID == "" {
		businessUnitID = currentPrincipal.BusinessUnitId
	}

	resolvedRoles, err := r.ApplicationClient.ResolveSecurityRoleNames(ctx, environmentID, businessUnitID, []string{securityRoleName})
	if err != nil {
		return nil, err
	}
	if len(resolvedRoles) != 1 {
		return nil, fmt.Errorf("expected exactly one resolved security role for '%s', got %d", securityRoleName, len(resolvedRoles))
	}

	return &resolvedRoleAssignment{
		principal:      currentPrincipal,
		role:           resolvedRoles[0],
		businessUnitID: businessUnitID,
	}, nil
}

func principalHasRole(principal *applicationUserDto, roleID string) bool {
	for _, role := range principal.SecurityRoles {
		if role.RoleId == roleID {
			return true
		}
	}

	return false
}
