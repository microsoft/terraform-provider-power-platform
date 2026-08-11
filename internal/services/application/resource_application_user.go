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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/modifiers"
)

var _ resource.Resource = &ApplicationUserResource{}
var _ resource.ResourceWithImportState = &ApplicationUserResource{}

func NewApplicationUserResource() resource.Resource {
	return &ApplicationUserResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "application_user",
		},
	}
}

func (r *ApplicationUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *ApplicationUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the unique active Dataverse application user for a Microsoft Entra service principal within an environment using the Dataverse `systemusers` API. " +
			"Create adopts an existing active user with the same environment and application id, or creates a non-interactive application user (`accessmode = 4`) when absent. " +
			"Adoption fails on multiple active matches or an explicitly configured business-unit mismatch, and reconciles the requested disabled state. Dataverse security-role assignment is managed separately. " +
			"Adoption takes ownership: destroying this resource deactivates and permanently deletes the application user even if it was adopted rather than created by Terraform. To abandon an application user without deleting it, remove the resource from state (for example with `terraform state rm`) instead of destroying it.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique Dataverse application user ID (guid). This is the Dataverse principal ID for the application user.",
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
			"disabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the Dataverse application user is disabled. Defaults to `false` and is actively reconciled by the resource.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					modifiers.UseDefaultBoolValueWhenConfigNull(false),
				},
			},
		},
	}
}

func (r *ApplicationUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApplicationUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan ApplicationUserResourceModel
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

	user, created, err := r.resolveApplicationUserForCreate(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to create or adopt application user '%s' in environment '%s'", plan.ApplicationId.ValueString(), plan.EnvironmentId.ValueString()),
			err.Error(),
		)
		return
	}

	desiredDisabled := plan.Disabled.ValueBool()
	if user.IsDisabled != desiredDisabled {
		user, err = r.ApplicationClient.SetApplicationUserDisabledState(
			ctx,
			plan.EnvironmentId.ValueString(),
			user.SystemUserId,
			plan.ApplicationId.ValueString(),
			desiredDisabled,
		)
		if err != nil {
			if created {
				r.addCreateFailureWithCleanupDiagnostics(
					ctx,
					resp,
					plan.EnvironmentId.ValueString(),
					user.SystemUserId,
					fmt.Sprintf("Failed to reconcile disabled state for application user '%s'", plan.ApplicationId.ValueString()),
					err,
				)
			} else {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Failed to reconcile disabled state for adopted application user '%s'", plan.ApplicationId.ValueString()),
					fmt.Sprintf("%s\n\nThe pre-existing application user was not deleted.", err.Error()),
				)
			}
			return
		}
	}

	plan.Id = types.StringValue(user.SystemUserId)
	plan.SystemUserId = types.StringValue(user.SystemUserId)
	plan.BusinessUnitId = types.StringValue(user.BusinessUnitId)
	plan.Disabled = types.BoolValue(user.IsDisabled)

	stateDiags := resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		if created {
			r.addCreateFailureWithCleanupDiagnostics(ctx, resp, plan.EnvironmentId.ValueString(), user.SystemUserId, fmt.Sprintf("Failed to persist state for application user '%s'", plan.ApplicationId.ValueString()), errors.New("terraform state update failed after remote application user creation"))
		}
		return
	}
}

func (r *ApplicationUserResource) resolveApplicationUserForCreate(ctx context.Context, plan *ApplicationUserResourceModel) (*applicationUserDto, bool, error) {
	environmentID := plan.EnvironmentId.ValueString()
	applicationID := plan.ApplicationId.ValueString()
	existing, err := r.ApplicationClient.GetApplicationUser(ctx, environmentID, applicationID)
	if err == nil {
		if configuredBusinessUnitID := plan.BusinessUnitId.ValueString(); configuredBusinessUnitID != "" && !strings.EqualFold(configuredBusinessUnitID, existing.BusinessUnitId) {
			return nil, false, fmt.Errorf(
				"existing application user '%s' belongs to business unit '%s', but business_unit_id is configured as '%s'",
				existing.SystemUserId,
				existing.BusinessUnitId,
				configuredBusinessUnitID)
		}
		return existing, false, nil
	}
	if !errors.Is(err, customerrors.ErrObjectNotFound) {
		return nil, false, fmt.Errorf("failed to inspect existing application users: %w", err)
	}

	configuredBusinessUnitID := plan.BusinessUnitId.ValueString()
	businessUnitID := configuredBusinessUnitID
	if businessUnitID == "" {
		businessUnitID, err = r.ApplicationClient.GetRootBusinessUnitId(ctx, environmentID)
		if err != nil {
			return nil, false, fmt.Errorf("failed to resolve root business unit: %w", err)
		}
	}

	created, createErr := r.ApplicationClient.CreateScopedApplicationUser(ctx, environmentID, applicationID, businessUnitID)
	if createErr == nil {
		return created, true, nil
	}

	// A concurrent graph or a pre-existing row discovered only after Dataverse's
	// uniqueness response can still converge by the same semantic identity. The
	// read remains ambiguity-aware and never adopts deleted rows.
	existing, lookupErr := r.ApplicationClient.GetApplicationUser(ctx, environmentID, applicationID)
	if lookupErr == nil {
		if configuredBusinessUnitID != "" && !strings.EqualFold(configuredBusinessUnitID, existing.BusinessUnitId) {
			return nil, false, fmt.Errorf(
				"create failed and the existing application user '%s' belongs to business unit '%s', not requested business unit '%s': %w",
				existing.SystemUserId,
				existing.BusinessUnitId,
				configuredBusinessUnitID,
				createErr)
		}
		return existing, false, nil
	}
	if !errors.Is(lookupErr, customerrors.ErrObjectNotFound) {
		return nil, false, fmt.Errorf("%w\n\nAn additional adoption lookup failed: %s", createErr, lookupErr.Error())
	}
	return nil, false, createErr
}

func (r *ApplicationUserResource) addCreateFailureWithCleanupDiagnostics(ctx context.Context, resp *resource.CreateResponse, environmentID, systemUserID, summary string, err error) {
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

func (r *ApplicationUserResource) cleanupCreatedApplicationUser(ctx context.Context, environmentID, systemUserID string) error {
	if err := r.ApplicationClient.DeactivateSystemUser(ctx, environmentID, systemUserID); err != nil {
		return fmt.Errorf("failed to deactivate created application user '%s': %w", systemUserID, err)
	}

	if err := r.ApplicationClient.PermanentlyDeleteSystemUser(ctx, environmentID, systemUserID); err != nil {
		return fmt.Errorf("failed to delete created application user '%s': %w", systemUserID, err)
	}

	return nil
}

func (r *ApplicationUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state ApplicationUserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
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

	state.Id = types.StringValue(user.SystemUserId)
	state.SystemUserId = types.StringValue(user.SystemUserId)
	state.BusinessUnitId = types.StringValue(user.BusinessUnitId)
	state.Disabled = types.BoolValue(user.IsDisabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ApplicationUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan ApplicationUserResourceModel
	var state ApplicationUserResourceModel
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

	if currentUser.IsDisabled != plan.Disabled.ValueBool() {
		currentUser, err = r.ApplicationClient.SetApplicationUserDisabledState(
			ctx,
			state.EnvironmentId.ValueString(),
			currentUser.SystemUserId,
			state.ApplicationId.ValueString(),
			plan.Disabled.ValueBool(),
		)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to update disabled state for application user '%s'", state.ApplicationId.ValueString()),
				err.Error(),
			)
			return
		}
	}

	plan.SystemUserId = types.StringValue(currentUser.SystemUserId)
	plan.BusinessUnitId = types.StringValue(currentUser.BusinessUnitId)
	plan.Id = types.StringValue(currentUser.SystemUserId)
	plan.Disabled = types.BoolValue(currentUser.IsDisabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ApplicationUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state ApplicationUserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	systemUserID := state.SystemUserId.ValueString()
	if systemUserID == "" {
		systemUserID = state.Id.ValueString()
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
	}

	if err := r.ApplicationClient.DeactivateSystemUser(ctx, state.EnvironmentId.ValueString(), systemUserID); err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to deactivate system user for application '%s' in environment '%s'", state.ApplicationId.ValueString(), state.EnvironmentId.ValueString()),
			err.Error(),
		)
		return
	}

	if err := r.ApplicationClient.PermanentlyDeleteSystemUser(ctx, state.EnvironmentId.ValueString(), systemUserID); err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete system user for application '%s' in environment '%s'", state.ApplicationId.ValueString(), state.EnvironmentId.ValueString()),
			err.Error(),
		)
		return
	}
}

func (r *ApplicationUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

	user, err := r.ApplicationClient.GetApplicationUser(ctx, idParts[0], idParts[1])
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to import application user",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), user.SystemUserId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("application_id"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("system_user_id"), user.SystemUserId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("business_unit_id"), user.BusinessUnitId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("disabled"), user.IsDisabled)...)
}
