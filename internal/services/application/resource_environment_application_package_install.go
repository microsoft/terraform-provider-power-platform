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
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

func NewEnvironmentApplicationPackageInstallResource() resource.Resource {
	return &EnvironmentApplicationPackageInstallResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment_application_package_install",
		},
	}
}

type EnvironmentApplicationPackageInstallResource struct {
	helpers.TypeInfo
	ApplicationClient Client
}

type EnvironmentApplicationPackageInstallResourceModel struct {
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
	Id            types.String   `tfsdk:"id"`
	UniqueName    types.String   `tfsdk:"unique_name"`
	EnvironmentId types.String   `tfsdk:"environment_id"`
}

func (r *EnvironmentApplicationPackageInstallResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *EnvironmentApplicationPackageInstallResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		Description:         "PowerPlatform application",
		MarkdownDescription: "This resource allows you to install a Dynamics 365 application in an environment.\n\nThis is functionally equivalent to the 'Install' button in the Power Platform admin center or [`pac application install` in the Power Platform CLI](https://docs.microsoft.com/powerapps/developer/data-platform/powerapps-cli#pac-application-install).  This resource uses the [Install Application Package](https://learn.microsoft.com/rest/api/power-platform/appmanagement/applications/install-application-package) endpoint in the Power Platform API.\n\n~> This resource does not support updating or deleting applications.  The expected behavior is that the application is installed and remains installed until the environment is deleted.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique id (guid)",
				Description:         "Unique id (guid)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				Description: "Id of the Dynamics 365 environment",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"unique_name": schema.StringAttribute{
				Description: "Unique name of the application",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *EnvironmentApplicationPackageInstallResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	if req.ProviderData == nil {
		resp.Diagnostics.AddError("Failed to configure %s because provider data is nil", r.TypeName)
		return
	}

	clientApi := req.ProviderData.(*api.ProviderClient).Api
	if clientApi == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.ApplicationClient = NewApplicationClient(clientApi)
}

func (r *EnvironmentApplicationPackageInstallResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state EnvironmentApplicationPackageInstallResourceModel
	resp.State.Get(ctx, &state)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Id = types.StringValue(fmt.Sprintf("%s_%s", state.EnvironmentId.ValueString(), strings.ReplaceAll(strings.ToLower(state.UniqueName.ValueString()), " ", "_")))
	state.EnvironmentId = types.StringValue(state.EnvironmentId.ValueString())
	state.UniqueName = types.StringValue(state.UniqueName.ValueString())

	dvExits, err := r.ApplicationClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
	}

	if !dvExits {
		resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), "")
		return
	}

	applicationId, err := r.ApplicationClient.InstallApplicationInEnvironment(ctx, state.EnvironmentId.ValueString(), state.UniqueName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.ProviderTypeName), err.Error())
		return
	}

	state.Id = types.StringValue(applicationId)

	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID %s", state.UniqueName.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EnvironmentApplicationPackageInstallResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *EnvironmentApplicationPackageInstallResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("READ: %s_application with application_name %s", r.ProviderTypeName, state.UniqueName.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EnvironmentApplicationPackageInstallResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *EnvironmentApplicationPackageInstallResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state *EnvironmentApplicationPackageInstallResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, "No application have been updated, as this is the expected behavior")
}

func (r *EnvironmentApplicationPackageInstallResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *EnvironmentApplicationPackageInstallResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "No application have been uninstalled, as this is the expected behavior")
}

func (r *EnvironmentApplicationPackageInstallResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
