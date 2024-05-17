// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

func NewEnvironmentApplicationPackageInstallResource() resource.Resource {
	return &EnvironmentApplicationPackageInstallResource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_environment_application_package_install",
	}
}

type EnvironmentApplicationPackageInstallResource struct {
	ApplicationClient ApplicationClient
	ProviderTypeName  string
	TypeName          string
}

type EnvironmentApplicationPackageInstallResourceModel struct {
	Id            types.String `tfsdk:"id"`
	UniqueName    types.String `tfsdk:"unique_name"`
	EnvironmentId types.String `tfsdk:"environment_id"`
}

func (r *EnvironmentApplicationPackageInstallResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

func (r *EnvironmentApplicationPackageInstallResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "PowerPlatform application",
		MarkdownDescription: "This resource allows you to install a Dynamics 365 application in an environment.\n\nThis is functionally equivalent to the 'Install' button in the Power Platform admin center or [`pac application install` in the Power Platform CLI](https://docs.microsoft.com/en-us/powerapps/developer/data-platform/powerapps-cli#pac-application-install).  This resource uses the [Install Application Package](https://docs.microsoft.com/en-us/rest/api/power-platform/appmanagement/applications/installapplicationpackage) endpoint in the Power Platform API.\n\n~> This resource does not support updating or deleting applications.  The expected behavior is that the application is installed and remains installed until the environment is deleted.",
		Attributes: map[string]schema.Attribute{
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
	if req.ProviderData == nil {
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
	var plan EnvironmentApplicationPackageInstallResourceModel
	resp.State.Get(ctx, &plan)

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	plan.Id = types.StringValue(fmt.Sprintf("%s_%s", plan.EnvironmentId.ValueString(), strings.ReplaceAll(strings.ToLower(plan.UniqueName.ValueString()), " ", "_")))
	plan.EnvironmentId = types.StringValue(plan.EnvironmentId.ValueString())
	plan.UniqueName = types.StringValue(plan.UniqueName.ValueString())

	dvExits, err := r.ApplicationClient.DataverseExists(ctx, plan.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), err.Error())
	}

	if !dvExits {
		resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), "")
		return
	}

	applicationId, err := r.ApplicationClient.InstallApplicationInEnvironment(ctx, plan.EnvironmentId.ValueString(), plan.UniqueName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.ProviderTypeName), err.Error())
		return
	}

	plan.Id = types.StringValue(applicationId)

	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID %s", plan.UniqueName.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *EnvironmentApplicationPackageInstallResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *EnvironmentApplicationPackageInstallResourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("READ: %s_application with application_name %s", r.ProviderTypeName, state.UniqueName.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE END: %s", r.ProviderTypeName))
}

func (r *EnvironmentApplicationPackageInstallResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *EnvironmentApplicationPackageInstallResourceModel

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state *EnvironmentApplicationPackageInstallResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE END: %s", r.ProviderTypeName))
	tflog.Debug(ctx, "No application have been updated, as this is the expected behavior")
}

func (r *EnvironmentApplicationPackageInstallResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *EnvironmentApplicationPackageInstallResourceModel

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE END: %s", r.ProviderTypeName))
	tflog.Debug(ctx, "No application have been uninstalled, as this is the expected behavior")
}

func (r *EnvironmentApplicationPackageInstallResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("application_name"), req, resp)
}
