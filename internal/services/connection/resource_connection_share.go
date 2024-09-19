// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package connection

import (
	"context"
	"fmt"

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
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ resource.Resource = &ShareResource{}
var _ resource.ResourceWithImportState = &ShareResource{}

func NewConnectionShareResource() resource.Resource {
	return &ShareResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "connection_share",
		},
	}
}

type ShareResource struct {
	helpers.TypeInfo
	ConnectionsClient ConnectionsClient
}

type ShareResourceModel struct {
	Timeouts      timeouts.Value              `tfsdk:"timeouts"`
	Id            types.String                `tfsdk:"id"`
	EnvironmentId types.String                `tfsdk:"environment_id"`
	ConnectorName types.String                `tfsdk:"connector_name"`
	ConnectionId  types.String                `tfsdk:"connection_id"`
	RoleName      types.String                `tfsdk:"role_name"`
	Principal     SharePrincipalResourceModel `tfsdk:"principal"`
}

type SharePrincipalResourceModel struct {
	EntraObjectId types.String `tfsdk:"entra_object_id"`
	DisplayName   types.String `tfsdk:"display_name"`
}

func (r *ShareResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *ShareResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the connection share",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the environment",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"connector_name": schema.StringAttribute{
				MarkdownDescription: "Name of the connector",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"connection_id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the connection",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role_name": schema.StringAttribute{
				MarkdownDescription: "Name of the role to assign to the principal",
				Required:            true,
				Computed:            false,
				Validators: []validator.String{
					stringvalidator.OneOf("CanView", "CanViewWithShare", "CanEdit"),
				},
			},
			"principal": schema.SingleNestedAttribute{
				MarkdownDescription: "Principal to share the connection with",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"entra_object_id": schema.StringAttribute{
						MarkdownDescription: "Entra Object Id of the principal",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"display_name": schema.StringAttribute{
						MarkdownDescription: "Display name of the principal",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
		},
	}
}

func (r *ShareResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.ConnectionsClient = NewConnectionsClient(clientApi)
}

func (r *ShareResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var plan *ShareResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.ConnectionsClient.ShareConnection(ctx, plan.EnvironmentId.ValueString(), plan.ConnectorName.ValueString(), plan.ConnectionId.ValueString(), plan.RoleName.ValueString(), plan.Principal.EntraObjectId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error sharing connection", err.Error())
		return
	}

	share, err := r.ConnectionsClient.GetConnectionShare(ctx, plan.EnvironmentId.ValueString(), plan.ConnectorName.ValueString(), plan.ConnectionId.ValueString(), plan.Principal.EntraObjectId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error getting connection share", err.Error())
		return
	}
	if share == nil {
		resp.Diagnostics.AddError("Error getting connection share", "Connection share not found")
	}

	state := ConvertFromConnectionResourceSharesDto(plan, share)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ShareResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var state *ShareResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	share, err := r.ConnectionsClient.GetConnectionShare(ctx, state.EnvironmentId.ValueString(), state.ConnectorName.ValueString(), state.ConnectionId.ValueString(), state.Principal.EntraObjectId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error getting connection share", err.Error())
		return
	}
	if share == nil {
		resp.Diagnostics.AddError("Error getting connection share", "Connection share not found")
	}

	newState := ConvertFromConnectionResourceSharesDto(state, share)

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *ShareResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *ShareResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state *ShareResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	share := ShareConnectionRequestDto{
		Put: []ShareConnectionRequestPutDto{
			{
				Properties: ShareConnectionRequestPutPropertiesDto{
					RoleName:     plan.RoleName.ValueString(),
					Capabilities: []any{},
					Principal: ShareConnectionRequestPutPropertiesPrincipalDto{
						Id:       plan.Principal.EntraObjectId.ValueString(),
						Type:     "ServicePrincipal",
						TenantId: nil,
					},
					NotifyShareTargetOption: "Notify",
				},
			},
		},
		Delete: []ShareConnectionRequestDeleteDto{},
	}

	err := r.ConnectionsClient.UpdateConnectionShare(ctx, plan.EnvironmentId.ValueString(), plan.ConnectorName.ValueString(), plan.ConnectionId.ValueString(), share)
	if err != nil {
		resp.Diagnostics.AddError("Error updating connection share", err.Error())
		return
	}

	newShare, err := r.ConnectionsClient.GetConnectionShare(ctx, plan.EnvironmentId.ValueString(), plan.ConnectorName.ValueString(), plan.ConnectionId.ValueString(), plan.Principal.EntraObjectId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error getting connection share", err.Error())
		return
	}
	if newShare == nil {
		resp.Diagnostics.AddError("Error getting connection share", "Connection share not found")
	}

	newState := ConvertFromConnectionResourceSharesDto(plan, newShare)

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *ShareResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *ShareResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.ConnectionsClient.DeleteConnectionShare(ctx, state.EnvironmentId.ValueString(), state.ConnectorName.ValueString(), state.ConnectionId.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting connection share", err.Error())
		return
	}
}

func (r *ShareResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func ConvertFromConnectionResourceSharesDto(oldPlan *ShareResourceModel, connection *ShareConnectionResponseDto) ShareResourceModel {
	share := ShareResourceModel{
		Timeouts:      oldPlan.Timeouts,
		EnvironmentId: oldPlan.EnvironmentId,
		ConnectorName: oldPlan.ConnectorName,
		ConnectionId:  oldPlan.ConnectionId,
		Id:            types.StringValue(connection.Name),
		RoleName:      types.StringValue(connection.Properties.RoleName),
		Principal: SharePrincipalResourceModel{
			EntraObjectId: types.StringValue(connection.Properties.Principal["id"].(string)),
			DisplayName:   types.StringValue(connection.Properties.Principal["displayName"].(string)),
		},
	}
	return share
}
