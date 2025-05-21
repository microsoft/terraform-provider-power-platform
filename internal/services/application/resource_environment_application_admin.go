// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package application

import (
	"context"
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

func NewEnvironmentApplicationAdminResource() resource.Resource {
	return &EnvironmentApplicationAdminResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment_application_admin",
		},
	}
}

// EnvironmentApplicationAdminResource defines the resource implementation.
type EnvironmentApplicationAdminResource struct {
	helpers.TypeInfo
	ApplicationClient client
}

// EnvironmentApplicationAdminResourceModel describes the resource data model.
type EnvironmentApplicationAdminResourceModel struct {
	EnvironmentId types.String   `tfsdk:"environment_id"`
	ApplicationId types.String   `tfsdk:"application_id"`
	Id            types.String   `tfsdk:"id"`
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
}

func (r *EnvironmentApplicationAdminResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *EnvironmentApplicationAdminResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		MarkdownDescription: "Ensures a Microsoft Entra **service principal** exists in a Dataverse environment as an **application user** with the **System Administrator** role. " +
			"*Required for imported environments.* Environments created by the SP already include it. " +
			"**Deletion is a no-op** — Dataverse currently exposes no API to remove application users. If you must revoke access, delete it manually in PPAC or via the Dataverse Web API. " +
			"**Reference**: [Create a Dataverse application user (preview)](https://learn.microsoft.com/en-us/power-platform/admin/create-dataverseapplicationuser)",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
			}),
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Dataverse environment ID.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"application_id": schema.StringAttribute{
				MarkdownDescription: "Service‑principal *application_id* (client ID).",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Composite ID `{environment_id}/{application_id}`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *EnvironmentApplicationAdminResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
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

func (r *EnvironmentApplicationAdminResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan EnvironmentApplicationAdminResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Dataverse exists in the environment
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
			"Environment application admin can only be added to environments with Dataverse.",
		)
		return
	}

	// Add the application user to the environment
	err = r.ApplicationClient.AddApplicationUser(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to add application user '%s' to environment '%s'", plan.ApplicationId.ValueString(), plan.EnvironmentId.ValueString()),
			err.Error(),
		)
		return
	}

	// Create composite ID
	plan.Id = types.StringValue(fmt.Sprintf("%s/%s", plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EnvironmentApplicationAdminResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state EnvironmentApplicationAdminResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the composite ID if it exists
	if !state.Id.IsNull() && state.Id.ValueString() != "" {
		idParts := strings.Split(state.Id.ValueString(), "/")
		if len(idParts) == 2 {
			// If ID is set but component parts aren't, set them
			if state.EnvironmentId.IsNull() || state.EnvironmentId.ValueString() == "" {
				state.EnvironmentId = types.StringValue(idParts[0])
			}
			if state.ApplicationId.IsNull() || state.ApplicationId.ValueString() == "" {
				state.ApplicationId = types.StringValue(idParts[1])
			}
		}
	}

	// Check if the application user exists in the environment
	exists, err := r.ApplicationClient.ApplicationUserExists(ctx, state.EnvironmentId.ValueString(), state.ApplicationId.ValueString())
	if err != nil {
		if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
			// Environment doesn't exist or we don't have access, remove resource from state
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to check if application user '%s' exists in environment '%s'", state.ApplicationId.ValueString(), state.EnvironmentId.ValueString()),
			err.Error(),
		)
		return
	}

	if !exists {
		// Application user doesn't exist, remove resource from state
		tflog.Debug(ctx, fmt.Sprintf("Application user '%s' not found in environment '%s', removing from state", state.ApplicationId.ValueString(), state.EnvironmentId.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	// Application user exists, set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EnvironmentApplicationAdminResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Update is not supported, all changes require replacement
	var plan EnvironmentApplicationAdminResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EnvironmentApplicationAdminResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state EnvironmentApplicationAdminResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deactivating and removing application user for app ID '%s' in environment '%s'",
		state.ApplicationId.ValueString(), state.EnvironmentId.ValueString()))

	// Get the system user ID for the application user
	systemUserId, err := r.ApplicationClient.GetApplicationUserSystemId(ctx, state.EnvironmentId.ValueString(), state.ApplicationId.ValueString())
	if err != nil {
		if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
			// Application user doesn't exist, nothing to delete
			tflog.Info(ctx, fmt.Sprintf("Application user '%s' not found in environment '%s', nothing to delete",
				state.ApplicationId.ValueString(), state.EnvironmentId.ValueString()))
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get system user ID for application user '%s' in environment '%s'",
				state.ApplicationId.ValueString(), state.EnvironmentId.ValueString()),
			err.Error(),
		)
		return
	}

	// First deactivate the system user
	err = r.ApplicationClient.DeactivateSystemUser(ctx, state.EnvironmentId.ValueString(), systemUserId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to deactivate system user for application '%s' in environment '%s'",
				state.ApplicationId.ValueString(), state.EnvironmentId.ValueString()),
			err.Error(),
		)
		return
	}

	// Now delete the system user
	err = r.ApplicationClient.DeleteSystemUser(ctx, state.EnvironmentId.ValueString(), systemUserId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete system user for application '%s' in environment '%s'",
				state.ApplicationId.ValueString(), state.EnvironmentId.ValueString()),
			err.Error(),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully removed application user '%s' from environment '%s'",
		state.ApplicationId.ValueString(), state.EnvironmentId.ValueString()))
}

func (r *EnvironmentApplicationAdminResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Import ID format: {environment_id}/{application_id}
	idParts := strings.Split(req.ID, "/")
	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in format 'environment_id/application_id', got '%s'", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("application_id"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
