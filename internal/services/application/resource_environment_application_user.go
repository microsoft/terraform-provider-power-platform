// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package application

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ resource.Resource = &EnvironmentApplicationUserResource{}
var _ resource.ResourceWithImportState = &EnvironmentApplicationUserResource{}

func NewEnvironmentApplicationUserResource() resource.Resource {
	return &EnvironmentApplicationUserResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment_application_user",
		},
	}
}

func (r *EnvironmentApplicationUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *EnvironmentApplicationUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates a Dataverse application user for a Microsoft Entra service principal within an environment using the Dataverse `systemusers` API. " +
			"The user is created as a non-interactive application user (`accessmode = 4`). Security roles are requested by name, resolved within the selected business unit, and assigned through Dataverse role associations.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Composite ID `{environment_id}/{application_id}`.",
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
			"application_id": schema.StringAttribute{
				MarkdownDescription: "Service principal application (client) ID.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"system_user_id": schema.StringAttribute{
				MarkdownDescription: "Dataverse system user ID for the application user.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"business_unit_id": schema.StringAttribute{
				MarkdownDescription: "Business unit ID for the application user. Defaults to the root business unit when omitted.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"security_roles": schema.SetAttribute{
				MarkdownDescription: "Security role names to assign to the application user. Role names are resolved within the selected business unit and must match exactly.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
			"resolved_security_roles": schema.ListNestedAttribute{
				MarkdownDescription: "Resolved security-role assignments for this environment-specific application user.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Resolved security role name.",
							Computed:            true,
						},
						"role_id": schema.StringAttribute{
							MarkdownDescription: "Environment-specific Dataverse role ID.",
							Computed:            true,
						},
						"business_unit_id": schema.StringAttribute{
							MarkdownDescription: "Business unit ID for the resolved role.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (r *EnvironmentApplicationUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EnvironmentApplicationUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan EnvironmentApplicationUserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dvExists, err := r.ApplicationClient.DataverseExists(ctx, plan.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()),
			err.Error(),
		)
		return
	}
	if !dvExists {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Environment '%s' does not have Dataverse", plan.EnvironmentId.ValueString()),
			"Environment application users can only be added to environments with Dataverse.",
		)
		return
	}

	businessUnitID := plan.BusinessUnitId.ValueString()
	if businessUnitID == "" {
		businessUnitID, err = r.ApplicationClient.GetRootBusinessUnitId(ctx, plan.EnvironmentId.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to resolve root business unit for environment '%s'", plan.EnvironmentId.ValueString()),
				err.Error(),
			)
			return
		}
	}

	user, err := r.ApplicationClient.CreateScopedApplicationUser(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString(), businessUnitID)
	if err != nil {
		r.addCreateFailureDiagnostics(ctx, resp, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString(), err)
		return
	}

	if len(plan.SecurityRoles) > 0 {
		resolvedRoles, err := r.ApplicationClient.ResolveSecurityRoleNames(ctx, plan.EnvironmentId.ValueString(), businessUnitID, plan.SecurityRoles)
		if err != nil {
			r.addCreateFailureWithCleanupDiagnostics(ctx, resp, plan.EnvironmentId.ValueString(), user.SystemUserId, fmt.Sprintf("Failed to resolve security roles for application user '%s'", plan.ApplicationId.ValueString()), err)
			return
		}

		roleIDs := make([]string, 0, len(resolvedRoles))
		for _, role := range resolvedRoles {
			roleIDs = append(roleIDs, role.RoleId)
		}

		user, err = r.ApplicationClient.AddApplicationUserSecurityRoles(ctx, plan.EnvironmentId.ValueString(), user.SystemUserId, roleIDs)
		if err != nil {
			r.addCreateFailureWithCleanupDiagnostics(ctx, resp, plan.EnvironmentId.ValueString(), user.SystemUserId, fmt.Sprintf("Failed to assign security roles to application user '%s'", plan.ApplicationId.ValueString()), err)
			return
		}
	}

	plan.Id = types.StringValue(fmt.Sprintf("%s/%s", plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString()))
	plan.SystemUserId = types.StringValue(user.SystemUserId)
	plan.BusinessUnitId, plan.SecurityRoles, plan.ResolvedSecurityRoles = convertApplicationUserToModel(user)

	stateDiags := resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		r.addCreateFailureWithCleanupDiagnostics(ctx, resp, plan.EnvironmentId.ValueString(), user.SystemUserId, fmt.Sprintf("Failed to persist state for application user '%s'", plan.ApplicationId.ValueString()), errors.New("terraform state update failed after remote application user creation"))
		return
	}
}

func (r *EnvironmentApplicationUserResource) addCreateFailureDiagnostics(ctx context.Context, resp *resource.CreateResponse, environmentID, applicationID string, err error) {
	exists, existsErr := r.ApplicationClient.ApplicationUserExists(ctx, environmentID, applicationID)
	if existsErr == nil && exists {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Application user '%s' already exists in environment '%s'", applicationID, environmentID),
			fmt.Sprintf("Terraform cannot adopt existing application users during create. Import it instead with '%s/%s'. Original create error: %s", environmentID, applicationID, err.Error()),
		)
		return
	}

	detail := err.Error()
	if existsErr != nil {
		detail = fmt.Sprintf("%s\n\nAn additional existence check failed: %s", detail, existsErr.Error())
	}

	resp.Diagnostics.AddError(
		fmt.Sprintf("Failed to create application user '%s' in environment '%s'", applicationID, environmentID),
		detail,
	)
}

func (r *EnvironmentApplicationUserResource) addCreateFailureWithCleanupDiagnostics(ctx context.Context, resp *resource.CreateResponse, environmentID, systemUserID, summary string, err error) {
	cleanupErr := r.cleanupCreatedApplicationUser(ctx, environmentID, systemUserID)
	if cleanupErr != nil {
		resp.Diagnostics.AddError(
			summary,
			fmt.Sprintf("%s\n\nRollback failed after create. The remote application user may still exist and must be removed manually. Cleanup error: %s", err.Error(), cleanupErr.Error()),
		)
		return
	}

	resp.Diagnostics.AddError(
		summary,
		fmt.Sprintf("%s\n\nRollback succeeded. The partially created remote application user was removed.", err.Error()),
	)
}

func (r *EnvironmentApplicationUserResource) cleanupCreatedApplicationUser(ctx context.Context, environmentID, systemUserID string) error {
	if err := r.ApplicationClient.DeactivateSystemUser(ctx, environmentID, systemUserID); err != nil {
		return fmt.Errorf("failed to deactivate created application user '%s': %w", systemUserID, err)
	}

	if err := r.ApplicationClient.DeleteSystemUser(ctx, environmentID, systemUserID); err != nil {
		return fmt.Errorf("failed to delete created application user '%s': %w", systemUserID, err)
	}

	return nil
}

func (r *EnvironmentApplicationUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state EnvironmentApplicationUserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.Id.IsNull() && state.Id.ValueString() != "" {
		idParts := strings.Split(state.Id.ValueString(), "/")
		if len(idParts) == 2 {
			if state.EnvironmentId.IsNull() || state.EnvironmentId.ValueString() == "" {
				state.EnvironmentId = types.StringValue(idParts[0])
			}
			if state.ApplicationId.IsNull() || state.ApplicationId.ValueString() == "" {
				state.ApplicationId = types.StringValue(idParts[1])
			}
		}
	}

	user, err := r.ApplicationClient.GetApplicationUser(ctx, state.EnvironmentId.ValueString(), state.ApplicationId.ValueString())
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to read application user '%s' in environment '%s'", state.ApplicationId.ValueString(), state.EnvironmentId.ValueString()),
			err.Error(),
		)
		return
	}

	state.SystemUserId = types.StringValue(user.SystemUserId)
	state.BusinessUnitId, state.SecurityRoles, state.ResolvedSecurityRoles = convertApplicationUserToModel(user)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EnvironmentApplicationUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan EnvironmentApplicationUserResourceModel
	var state EnvironmentApplicationUserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	currentUser, err := r.ApplicationClient.GetApplicationUser(ctx, state.EnvironmentId.ValueString(), state.ApplicationId.ValueString())
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to read application user '%s' for update", state.ApplicationId.ValueString()),
			err.Error(),
		)
		return
	}

	businessUnitID := plan.BusinessUnitId.ValueString()
	if businessUnitID == "" {
		businessUnitID = currentUser.BusinessUnitId
	}

	resolvedDesiredRoles, err := r.ApplicationClient.ResolveSecurityRoleNames(ctx, state.EnvironmentId.ValueString(), businessUnitID, plan.SecurityRoles)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to resolve security roles for application user '%s'", state.ApplicationId.ValueString()),
			err.Error(),
		)
		return
	}

	desiredRoleIDs := make(map[string]applicationSecurityRoleDto, len(resolvedDesiredRoles))
	for _, role := range resolvedDesiredRoles {
		desiredRoleIDs[role.RoleId] = role
	}

	currentRoleIDs := make(map[string]applicationSecurityRoleDto, len(currentUser.SecurityRoles))
	for _, role := range currentUser.SecurityRoles {
		currentRoleIDs[role.RoleId] = role
	}

	rolesToAdd := make([]string, 0)
	for roleID := range desiredRoleIDs {
		if _, exists := currentRoleIDs[roleID]; !exists {
			rolesToAdd = append(rolesToAdd, roleID)
		}
	}

	rolesToRemove := make([]string, 0)
	for roleID := range currentRoleIDs {
		if _, exists := desiredRoleIDs[roleID]; !exists {
			rolesToRemove = append(rolesToRemove, roleID)
		}
	}

	if len(rolesToAdd) > 0 {
		currentUser, err = r.ApplicationClient.AddApplicationUserSecurityRoles(ctx, state.EnvironmentId.ValueString(), currentUser.SystemUserId, rolesToAdd)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to add security roles for application user '%s'", state.ApplicationId.ValueString()),
				err.Error(),
			)
			return
		}
	}

	if len(rolesToRemove) > 0 {
		currentUser, err = r.ApplicationClient.RemoveApplicationUserSecurityRoles(ctx, state.EnvironmentId.ValueString(), currentUser.SystemUserId, rolesToRemove)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to remove security roles for application user '%s'", state.ApplicationId.ValueString()),
				err.Error(),
			)
			return
		}
	}

	plan.SystemUserId = types.StringValue(currentUser.SystemUserId)
	plan.BusinessUnitId, plan.SecurityRoles, plan.ResolvedSecurityRoles = convertApplicationUserToModel(currentUser)
	plan.Id = state.Id

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EnvironmentApplicationUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state EnvironmentApplicationUserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	systemUserID := state.SystemUserId.ValueString()
	if systemUserID == "" {
		resolvedSystemUserID, err := r.ApplicationClient.GetApplicationUserSystemId(ctx, state.EnvironmentId.ValueString(), state.ApplicationId.ValueString())
		if err != nil {
			if errors.Is(err, customerrors.ErrObjectNotFound) {
				return
			}
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to get system user ID for application user '%s' in environment '%s'", state.ApplicationId.ValueString(), state.EnvironmentId.ValueString()),
				err.Error(),
			)
			return
		}
		systemUserID = resolvedSystemUserID
	}

	if err := r.ApplicationClient.DeactivateSystemUser(ctx, state.EnvironmentId.ValueString(), systemUserID); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to deactivate system user for application '%s' in environment '%s'", state.ApplicationId.ValueString(), state.EnvironmentId.ValueString()),
			err.Error(),
		)
		return
	}

	if err := r.ApplicationClient.DeleteSystemUser(ctx, state.EnvironmentId.ValueString(), systemUserID); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete system user for application '%s' in environment '%s'", state.ApplicationId.ValueString(), state.EnvironmentId.ValueString()),
			err.Error(),
		)
		return
	}
}

func (r *EnvironmentApplicationUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	idParts := strings.Split(req.ID, "/")
	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in format 'environment_id/application_id', got '%s'", req.ID),
		)
		return
	}

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("application_id"), idParts[1])...)
}
