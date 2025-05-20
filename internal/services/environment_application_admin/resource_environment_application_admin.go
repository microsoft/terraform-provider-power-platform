// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_application_admin

import (
	"context"
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
	"github.com/microsoft/terraform-provider-power-platform/internal/customtypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

// NewEnvironmentApplicationAdminResource creates a new instance of the resource.
func NewEnvironmentApplicationAdminResource() resource.Resource {
	return &EnvironmentApplicationAdminResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment_application_admin",
		},
	}
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
		MarkdownDescription: "Ensures a Microsoft Entra **service principal** exists in a Dataverse environment as an **application user** with the **System Administrator** role.\n\n" +
			"*Required for imported environments.* Environments created by the SP already include it.\n\n" +
			"**Deletion is a no‑op** — Dataverse currently exposes no API to remove application users. If you must revoke access, delete it manually in PPAC or via the Dataverse Web API.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
			}),
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Dataverse environment ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(36, 36),
				},
			},
			"application_id": schema.StringAttribute{
				MarkdownDescription: "Service‑principal *application_id* (client ID).",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(36, 36),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Composite ID `{environment_id}/{application_id}`.",
				Computed:            true,
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

	r.EnvironmentApplicationAdminClient = newEnvironmentApplicationAdminClient(client.Api)
}

func (r *EnvironmentApplicationAdminResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Import format: {environment_id}/{application_id}
	idParts := strings.Split(req.ID, "/")
	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Error parsing import ID",
			fmt.Sprintf("Expected format: {environment_id}/{application_id}, got: %s", req.ID),
		)
		return
	}

	environmentId := idParts[0]
	applicationId := idParts[1]

	// Set the attributes in the state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), environmentId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("application_id"), applicationId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func (r *EnvironmentApplicationAdminResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state EnvironmentApplicationAdminResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Query Dataverse to check if application user exists.
	exists, err := r.EnvironmentApplicationAdminClient.GetApplicationUser(
		ctx,
		state.EnvironmentId.ValueString(),
		state.ApplicationId.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError("Client error when reading application user", err.Error())
		return
	}

	// If the application user doesn't exist, remove resource from state.
	if !exists {
		tflog.Debug(ctx, "Application user not found in Dataverse, removing from state", map[string]any{
			"environment_id": state.EnvironmentId.ValueString(),
			"application_id": state.ApplicationId.ValueString(),
		})
		resp.State.RemoveResource(ctx)
		return
	}

	// Application user exists, keep state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EnvironmentApplicationAdminResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan EnvironmentApplicationAdminResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Add the application user to the environment
	err := r.EnvironmentApplicationAdminClient.AddApplicationUser(
		ctx,
		plan.EnvironmentId.ValueString(),
		plan.ApplicationId.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError("Failed to add service principal as application user", err.Error())
		return
	}

	// Set the composite ID
	compositeId := fmt.Sprintf("%s/%s", plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString())
	plan.Id = types.StringValue(compositeId)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EnvironmentApplicationAdminResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Dataverse API does not expose a way to remove application users
	// Document this as a no-op in the resource description
	tflog.Info(ctx, "Delete is a no-op for environment_application_admin resource", map[string]any{
		"message": "Dataverse does not provide an API to remove application users. The user must be removed manually if needed.",
	})
}

func (r *EnvironmentApplicationAdminResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// All attributes require replacement, so Update should never be called
	resp.Diagnostics.AddError("Update not supported", "Update operation should not be triggered as all attributes have ForceNew set")
}
