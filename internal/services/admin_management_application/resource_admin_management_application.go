// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package admin_management_application

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customtypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

func NewAdminManagementApplicationResource() resource.Resource {
	return &AdminManagementApplicationResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "admin_management_application",
		},
	}
}

type AdminManagementApplicationResource struct {
	helpers.TypeInfo
	AdminManagementApplicationClient AdminManagementApplicationClient
}

type AdminManagementApplicationResourceModel struct {
	Timeouts timeouts.Value        `tfsdk:"timeouts"`
	Id       customtypes.UUIDValue `tfsdk:"id"`
}

func (r *AdminManagementApplicationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *AdminManagementApplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		Description:         "Power Platform Admin Management Application",
		MarkdownDescription: "This resource allows you to register a service principal as an administrator for Power Platform.\n\nThis resource implements the process documented here [Registering an admin management application](https://learn.microsoft.com/power-platform/admin/powerplatform-api-create-service-principal). A service principal can't register itself—by design, the application must be registered by an administrator.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Delete: true,
				Read:   true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Client id for the service principal",
				Description:         "Client id for the service principal",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *AdminManagementApplicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		resp.Diagnostics.AddError("Failed to configure %s because provider data is nil", r.TypeName)
		return
	}

	clientApi := req.ProviderData.(*api.ProviderClient).Api
	if clientApi == nil {
		resp.Diagnostics.AddError("Failed to configure AdminManagementApplicationResource", "Failed to configure AdminManagementApplicationResource")
		return
	}

	r.AdminManagementApplicationClient = NewAdminManagementApplicationClient(clientApi)
}

func (r *AdminManagementApplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AdminManagementApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state AdminManagementApplicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	adminApp, err := r.AdminManagementApplicationClient.GetAdminApplication(ctx, state.Id.ValueString())
	if err != nil {
		return
	}

	newState := AdminManagementApplicationResourceModel{
		Id:       customtypes.NewUUIDValue(adminApp.ClientId),
		Timeouts: state.Timeouts,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *AdminManagementApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan AdminManagementApplicationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	adminApp, err := r.AdminManagementApplicationClient.RegisterAdminApplication(ctx, plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to register admin application", fmt.Sprintf("Failed to register admin application: %v", err))
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(ctx, &AdminManagementApplicationResourceModel{
			Id:       customtypes.NewUUIDValue(adminApp.ClientId),
			Timeouts: plan.Timeouts,
		})...)
}

func (r *AdminManagementApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state AdminManagementApplicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.AdminManagementApplicationClient.UnregisterAdminApplication(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to unregister admin application", fmt.Sprintf("Failed to unregister admin application: %v", err))
		return
	}
}

func (r *AdminManagementApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.Diagnostics.AddError("Update not supported", "Update not supported")
}
