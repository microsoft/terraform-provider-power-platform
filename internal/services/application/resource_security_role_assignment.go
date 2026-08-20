// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package application

import (
	"context"
	"errors"
	"fmt"
	"regexp"
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
	"github.com/microsoft/terraform-provider-power-platform/internal/customtypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ resource.Resource = &SecurityRoleAssignmentResource{}
var _ resource.ResourceWithImportState = &SecurityRoleAssignmentResource{}

type SecurityRoleAssignmentResource struct {
	helpers.TypeInfo
	ApplicationClient client
}

type SecurityRoleAssignmentResourceModel struct {
	Timeouts         timeouts.Value   `tfsdk:"timeouts"`
	Id               types.String     `tfsdk:"id"`
	EnvironmentId    customtypes.UUID `tfsdk:"environment_id"`
	SystemUserId     customtypes.UUID `tfsdk:"system_user_id"`
	TeamId           customtypes.UUID `tfsdk:"team_id"`
	BusinessUnitId   customtypes.UUID `tfsdk:"business_unit_id"`
	SecurityRoleName types.String     `tfsdk:"security_role_name"`
	SecurityRoleId   customtypes.UUID `tfsdk:"security_role_id"`
}

// holder is the Dataverse principal this assignment targets, chosen by whichever id is set.
func (m SecurityRoleAssignmentResourceModel) holder() roleHolder {
	if helpers.IsKnown(m.TeamId) {
		return teamRoleHolder(m.TeamId.ValueString())
	}
	return systemUserRoleHolder(m.SystemUserId.ValueString())
}

// compositeId identifies the assignment as {environment}/{entity set}/{principal}/{role id}. The
// entity set is included because a system user id and a team id are rows in different tables, and the
// role is identified by its immutable id rather than its name: Dataverse permits renaming custom
// roles while associations are made to role records by id, so a name in the identity would dangle on
// rename and cannot distinguish same-named roles in different business units.
func (m SecurityRoleAssignmentResourceModel) compositeId() string {
	h := m.holder()
	return fmt.Sprintf("%s/%s/%s/%s", m.EnvironmentId.ValueString(), h.entitySet(), h.id(), m.SecurityRoleId.ValueString())
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
		MarkdownDescription: "Assigns a single Dataverse security role to a principal: a user, an application user, or a team. The role name is resolved to its id within the target business unit at create time; from then on the assignment is anchored on the immutable role id, so renaming the role does not orphan it. If the association already exists it is adopted, and destroying the resource removes the association.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Delete: true,
				Read:   true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Composite ID `{environment_id}/{entity_set}/{principal_id}/{security_role_id}`, where entity set is `systemusers` or `teams`. The role is identified by its immutable id, not its name, so renaming a role does not orphan the assignment.",
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
				CustomType: customtypes.UUIDType{},
			},
			"system_user_id": schema.StringAttribute{
				MarkdownDescription: "Dataverse `systemuserid` of the user or application user the security role is assigned to. This is a Dataverse row id, not a Microsoft Entra object id, and `powerplatform_application_user` exposes it as `system_user_id`. Exactly one of `system_user_id` or `team_id` must be set.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				CustomType: customtypes.UUIDType{},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("system_user_id"), path.MatchRoot("team_id")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"team_id": schema.StringAttribute{
				MarkdownDescription: "Dataverse `teamid` of the team the security role is assigned to. Dataverse keeps teams in their own table with their own role association, so this is a different id from `system_user_id`. Only owner teams and Microsoft Entra group teams can hold security roles; an access team is refused. Exactly one of `system_user_id` or `team_id` must be set.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				CustomType: customtypes.UUIDType{},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"business_unit_id": schema.StringAttribute{
				MarkdownDescription: "Business unit the role belongs to. With `security_role_name` it scopes the name resolution and defaults to the principal's current business unit; with `security_role_id` it is computed from the role row, and a configured value must agree with it.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				CustomType: customtypes.UUIDType{},
			},
			"security_role_name": schema.StringAttribute{
				MarkdownDescription: "Dataverse security role name to assign, resolved to its id within the target business unit at create time. Exactly one of `security_role_name` or `security_role_id` must be set; the name is filled in from the live role when the id is used. Role names are not unique across business units, so an ambiguous name is refused: use `security_role_id` to pin one.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("security_role_name"), path.MatchRoot("security_role_id")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"security_role_id": schema.StringAttribute{
				MarkdownDescription: "Dataverse role ID of the security role to assign, computed when `security_role_name` is used. This id, not the role name, anchors the assignment from then on, so renaming the role does not affect it, and it is the way to pin one of several same-named roles in different business units.",
				Optional:            true,
				Computed:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
		plan.SecurityRoleId.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to assign security role %s to %s", plan.roleSelector(), plan.holder()),
			err.Error(),
		)
		return
	}

	if findAssignedRole(resolved.principal, resolved.role.RoleId) == nil {
		if err := r.ApplicationClient.AddPrincipalSecurityRoles(ctx, plan.EnvironmentId.ValueString(), plan.holder(), []string{resolved.role.RoleId}); err != nil {
			// The POST may have committed despite the failure. The association is idempotent, so
			// reading the relationship back settles it: if the role is there, the grant is live and
			// must be recorded in state rather than left as an active privilege outside it.
			principal, readErr := r.ApplicationClient.GetRoleHolder(ctx, plan.EnvironmentId.ValueString(), plan.holder())
			if readErr != nil || findAssignedRole(principal, resolved.role.RoleId) == nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Failed to assign security role %s to %s", plan.roleSelector(), plan.holder()),
					err.Error(),
				)
				return
			}
			tflog.Debug(ctx, fmt.Sprintf("Associating security role '%s' with %s failed but the association exists; recording it", resolved.role.RoleId, plan.holder()))
		}
	}

	plan.BusinessUnitId = customtypes.NewUUIDValue(resolved.businessUnitID)
	plan.SecurityRoleId = customtypes.NewUUIDValue(resolved.role.RoleId)
	if !helpers.IsKnown(plan.SecurityRoleName) {
		plan.SecurityRoleName = types.StringValue(resolved.role.Name)
	}
	plan.Id = types.StringValue(plan.compositeId())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// roleSelector names the configured role for error messages, whichever selector was used.
func (m SecurityRoleAssignmentResourceModel) roleSelector() string {
	if helpers.IsKnown(m.SecurityRoleName) {
		return fmt.Sprintf("'%s'", m.SecurityRoleName.ValueString())
	}
	return fmt.Sprintf("'%s'", m.SecurityRoleId.ValueString())
}

func (r *SecurityRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state SecurityRoleAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The assignment is anchored on the stored role id, never re-resolved by name. Renaming the
	// role therefore does not orphan the assignment, and same-named roles in other business units
	// cannot be confused with it.
	principal, err := r.ApplicationClient.GetRoleHolder(ctx, state.EnvironmentId.ValueString(), state.holder())
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to read security role assignment for %s", state.holder()),
			err.Error(),
		)
		return
	}

	assigned := findAssignedRole(principal, state.SecurityRoleId.ValueString())
	if assigned == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.BusinessUnitId = customtypes.NewUUIDValue(assigned.BusinessUnitId)
	// The name is a create-time selector, so an existing value is left alone even if the role has
	// been renamed since; the id keeps the assignment attached regardless. It is only filled in
	// when absent, which is the import path.
	if !helpers.IsKnown(state.SecurityRoleName) {
		state.SecurityRoleName = types.StringValue(assigned.Name)
	}
	state.Id = types.StringValue(state.compositeId())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SecurityRoleAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Every attribute requires replacement, so there is no in-place update.
}

func (r *SecurityRoleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state SecurityRoleAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	principal, err := r.ApplicationClient.GetRoleHolder(ctx, state.EnvironmentId.ValueString(), state.holder())
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			// The principal (or its environment) is gone, and the association went with it.
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to read security role assignment for %s", state.holder()),
			err.Error(),
		)
		return
	}

	if findAssignedRole(principal, state.SecurityRoleId.ValueString()) == nil {
		// Already unassigned; destruction is idempotent.
		return
	}

	if err = r.ApplicationClient.RemovePrincipalSecurityRoles(ctx, state.EnvironmentId.ValueString(), state.holder(), []string{state.SecurityRoleId.ValueString()}); err != nil {
		// The principal vanishing between the read above and the removal is the destruction this
		// wanted; only a genuine failure should stop the destroy.
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to remove security role '%s' from %s", state.SecurityRoleId.ValueString(), state.holder()),
			err.Error(),
		)
	}
}

// ImportState accepts {environment_id}/{systemusers|teams}/{principal_id}/{security_role_id}. The role is
// identified by its immutable id: names can be renamed and duplicated across business units, so an
// imported name could attach to the wrong role. Read fills security_role_name in from the live role.
func (r *SecurityRoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	idParts := strings.Split(req.ID, "/")
	if len(idParts) != 4 || (idParts[1] != "systemusers" && idParts[1] != "teams") || idParts[0] == "" || idParts[2] == "" || idParts[3] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in format 'environment_id/systemusers/{system_user_id}/{security_role_id}' or 'environment_id/teams/{team_id}/{security_role_id}', got '%s'", req.ID),
		)
		return
	}

	principalAttribute := "system_user_id"
	if idParts[1] == "teams" {
		principalAttribute = "team_id"
	}

	// The guid attributes carry semantic equality, so casing never causes drift; the segments are
	// still validated so a malformed id fails at import rather than at the first refresh.
	guidRegex := regexp.MustCompile(helpers.GuidRegex)
	for _, part := range []string{idParts[0], idParts[2], idParts[3]} {
		if !guidRegex.MatchString(part) {
			resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf("'%s' is not a guid", part))
			return
		}
	}

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(principalAttribute), idParts[2])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("security_role_id"), idParts[3])...)
}

type resolvedRoleAssignment struct {
	principal      *roleHolderDto
	role           applicationSecurityRoleDto
	businessUnitID string
}

// resolveRequestedRole resolves the configured role selector, a NAME or an ID, to the role row for
// Create. This is the only place the name is used; every later operation works from the resolved id.
func (r *SecurityRoleAssignmentResource) resolveRequestedRole(ctx context.Context, environmentID string, holder roleHolder, requestedBusinessUnitID, securityRoleName, securityRoleID string) (*resolvedRoleAssignment, error) {
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

	if securityRoleID != "" {
		// An id needs no business unit to resolve: the role row carries its own, which becomes the
		// computed business_unit_id. A configured business unit still has to agree with it.
		role, err := r.ApplicationClient.GetSecurityRoleById(ctx, environmentID, securityRoleID)
		if err != nil {
			return nil, err
		}
		if requestedBusinessUnitID != "" && !strings.EqualFold(role.BusinessUnitId, requestedBusinessUnitID) {
			return nil, fmt.Errorf("security role '%s' belongs to business unit '%s', not the configured business unit '%s'", securityRoleID, role.BusinessUnitId, requestedBusinessUnitID)
		}
		return &resolvedRoleAssignment{
			principal:      currentPrincipal,
			role:           *role,
			businessUnitID: role.BusinessUnitId,
		}, nil
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

// findAssignedRole returns the holder's assignment of the given role id, or nil when it is absent.
// findAssignedRole compares guids case-insensitively: Dataverse renders them lowercase, but an
// imported id may not be.
func findAssignedRole(principal *roleHolderDto, roleID string) *applicationSecurityRoleDto {
	for i := range principal.SecurityRoles {
		if strings.EqualFold(principal.SecurityRoles[i].RoleId, roleID) {
			return &principal.SecurityRoles[i]
		}
	}
	return nil
}
