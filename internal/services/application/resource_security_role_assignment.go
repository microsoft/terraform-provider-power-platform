// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package application

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	TeamId           types.String   `tfsdk:"team_id"`
	BusinessUnitId   types.String   `tfsdk:"business_unit_id"`
	SecurityRoleName types.String   `tfsdk:"security_role_name"`
	RoleId           types.String   `tfsdk:"role_id"`
}

// holder is the Dataverse principal this assignment targets, chosen by whichever id is set.
func (m SecurityRoleAssignmentResourceModel) holder() roleHolder {
	if helpers.IsKnown(m.TeamId) {
		return teamRoleHolder(m.TeamId.ValueString())
	}
	return systemUserRoleHolder(m.SystemUserId.ValueString())
}

// compositeId identifies the assignment as {environment}/{entity set}/{principal}/{role}. The entity
// set is included because a system user id and a team id are rows in different tables, so without it
// an imported id would be ambiguous.
func (m SecurityRoleAssignmentResourceModel) compositeId() string {
	h := m.holder()
	return fmt.Sprintf("%s/%s/%s/%s", m.EnvironmentId.ValueString(), h.entitySet(), h.id(), m.SecurityRoleName.ValueString())
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
				MarkdownDescription: "Composite ID `{environment_id}/{entity_set}/{principal_id}/{security_role_name}`, where entity set is `systemusers` or `teams`.",
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
				MarkdownDescription: "Dataverse `systemuserid` of the user or application user the security role is assigned to. This is a Dataverse row id, not a Microsoft Entra object id, and `powerplatform_application_user` exposes it as `system_user_id`. Exactly one of `system_user_id` or `team_id` must be set.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("system_user_id"), path.MatchRoot("team_id")),
				},
			},
			"team_id": schema.StringAttribute{
				MarkdownDescription: "Dataverse `teamid` of the team the security role is assigned to. Dataverse keeps teams in their own table with their own role association, so this is a different id from `system_user_id`. Exactly one of `system_user_id` or `team_id` must be set.",
				Optional:            true,
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
		plan.holder(),
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
		resolved.principal, err = r.ApplicationClient.AddPrincipalSecurityRoles(ctx, plan.EnvironmentId.ValueString(), plan.holder(), []string{resolved.role.RoleId})
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to assign security role '%s' to principal '%s'", plan.SecurityRoleName.ValueString(), plan.SystemUserId.ValueString()),
				err.Error(),
			)
			return
		}
	}

	plan.Id = types.StringValue(plan.compositeId())
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
		state.holder(),
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
		plan.holder(),
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

	plan.Id = types.StringValue(plan.compositeId())
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
		state.holder(),
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

	if _, err = r.ApplicationClient.RemovePrincipalSecurityRoles(ctx, state.EnvironmentId.ValueString(), state.holder(), []string{resolved.role.RoleId}); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to remove security role '%s' from principal '%s'", state.SecurityRoleName.ValueString(), state.SystemUserId.ValueString()),
			err.Error(),
		)
	}
}

func (r *SecurityRoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	idParts := strings.SplitN(req.ID, "/", 4)
	if len(idParts) != 4 || (idParts[1] != "systemusers" && idParts[1] != "teams") {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in format 'environment_id/systemusers/{system_user_id}/security_role_name' or 'environment_id/teams/{team_id}/security_role_name', got '%s'", req.ID),
		)
		return
	}

	principalAttribute := "system_user_id"
	if idParts[1] == "teams" {
		principalAttribute = "team_id"
	}

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(principalAttribute), idParts[2])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("security_role_name"), idParts[3])...)
}

type resolvedRoleAssignment struct {
	principal      *roleHolderDto
	role           applicationSecurityRoleDto
	businessUnitID string
}

func (r *SecurityRoleAssignmentResource) resolveRequestedRole(ctx context.Context, environmentID string, holder roleHolder, requestedBusinessUnitID, securityRoleName string) (*resolvedRoleAssignment, error) {
	dvExists, err := r.ApplicationClient.DataverseExists(ctx, environmentID)
	if err != nil {
		return nil, err
	}
	if !dvExists {
		return nil, fmt.Errorf("environment '%s' does not have Dataverse", environmentID)
	}

	currentPrincipal, err := r.ApplicationClient.GetRoleHolder(ctx, environmentID, holder)
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

func principalHasRole(principal *roleHolderDto, roleID string) bool {
	for _, role := range principal.SecurityRoles {
		if role.RoleId == roleID {
			return true
		}
	}

	return false
}
