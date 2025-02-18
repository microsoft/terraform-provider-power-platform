// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package connection

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ resource.Resource = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}

func NewConnectionResource() resource.Resource {
	return &Resource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "connection",
		},
	}
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a [Connection](https://learn.microsoft.com/en-us/power-apps/maker/canvas-apps/add-manage-connections). A connection in Power Platform serves as a means to integrate external data sources and services with your Power Platform apps, flows, and other solutions. It acts as a bridge, facilitating secure communication between your solutions and various external systems.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique connection id",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Environment id where the connection is to be created",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the connection. This can be found using `powerplatform_connectors` data source by using the `name` attribute",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the connection",
				Required:            true,
			},
			"status": schema.SetAttribute{
				MarkdownDescription: "List of connection statuses",
				ElementType:         types.StringType,
				Computed:            true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"connection_parameters": schema.StringAttribute{
				MarkdownDescription: "Connection parameters. Json string containing the authentication connection parameters (if connection is interactive, leave blank), (for example)[https://learn.microsoft.com/en-us/power-automate/desktop-flows/alm/alm-connection#create-a-connection-using-your-service-principal]. Depending on required authentication parameters of a given connector, the connection parameters can vary. Due to how connection parameters and served by the platform, not all values are retrieved. If you don't want the connection to requried in-place-update all the time, consider using `ignore_changes` in the resource block.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connection_parameters_set": schema.StringAttribute{
				MarkdownDescription: "Set of connection parameters. Json string containing the authentication connection parameters (if connection is interactive, leave blank), (for example)[https://learn.microsoft.com/en-us/power-automate/desktop-flows/alm/alm-connection#create-a-connection-using-your-service-principal]. Depending on required authentication parameters of a given connector, the connection parameters can vary. Due to how connection parameters and served by the platform, not all values are retrieved. If you don't want the connection to requried in-place-update all the time, consider using `ignore_changes` in the resource block.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (d *Resource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("connection_parameters"),
			path.MatchRoot("connection_parameters_set"),
		),
	}
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.ConnectionsClient = newConnectionsClient(client.Api)
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var plan *ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	connectionToCreate := createDto{
		Properties: createPropertiesDto{
			DisplayName: plan.DisplayName.ValueString(),
			Environment: createEnvironmentDto{
				Name: plan.EnvironmentId.ValueString(),
				Id:   fmt.Sprintf("/providers/Microsoft.PowerApps/environments/%s", plan.EnvironmentId.ValueString()),
			},
		},
	}

	if !plan.ConnectionParameters.IsNull() && plan.ConnectionParameters.ValueString() != "" {
		var params map[string]any
		err := json.Unmarshal([]byte(plan.ConnectionParameters.ValueString()), &params)
		if err != nil {
			resp.Diagnostics.AddError("Failed to convert connection parameters", err.Error())
			return
		}
		connectionToCreate.Properties.ConnectionParameters = params
	}
	if !plan.ConnectionParametersSet.IsNull() && plan.ConnectionParametersSet.ValueString() != "" {
		var params map[string]any
		err := json.Unmarshal([]byte(plan.ConnectionParametersSet.ValueString()), &params)
		if err != nil {
			resp.Diagnostics.AddError("Failed to convert connection parameters set", err.Error())
			return
		}
		connectionToCreate.Properties.ConnectionParametersSet = params
	}

	connection, err := r.ConnectionsClient.CreateConnection(ctx, plan.EnvironmentId.ValueString(), plan.Name.ValueString(), connectionToCreate)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create connection", err.Error())
		return
	}

	conectionState := ConvertFromConnectionDto(*connection)
	plan.Id = types.String(conectionState.Id)
	statuses := []attr.Value{}
	for _, status := range conectionState.Status {
		statuses = append(statuses, types.StringValue(status))
	}
	plan.Status = types.SetValueMust(types.StringType, statuses)
	plan.DisplayName = types.String(conectionState.DisplayName)
	plan.Name = types.String(conectionState.Name)
	if conectionState.ConnectionParameters == types.StringNull() {
		plan.ConnectionParameters = types.StringValue("")
	}

	if conectionState.ConnectionParametersSet == types.StringNull() {
		plan.ConnectionParametersSet = types.StringValue("")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connection, err := r.ConnectionsClient.GetConnection(ctx, state.EnvironmentId.ValueString(), state.Name.ValueString(), state.Id.ValueString())
	if err != nil {
		if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	conectionState := ConvertFromConnectionDto(*connection)
	state.Id = types.String(conectionState.Id)
	state.DisplayName = types.String(conectionState.DisplayName)
	statuses := []attr.Value{}
	for _, status := range conectionState.Status {
		statuses = append(statuses, types.StringValue(status))
	}
	state.Status = types.SetValueMust(types.StringType, statuses)
	state.Name = types.String(conectionState.Name)
	state.ConnectionParameters = types.String(conectionState.ConnectionParameters)
	state.ConnectionParametersSet = types.String(conectionState.ConnectionParametersSet)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state *ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var connParams map[string]any
	if !plan.ConnectionParameters.IsNull() && plan.ConnectionParameters.ValueString() != "" {
		err := json.Unmarshal([]byte(plan.ConnectionParameters.ValueString()), &connParams)
		if err != nil {
			resp.Diagnostics.AddError("Failed to convert connection parameters", err.Error())
			return
		}
	}

	var connParamsSet map[string]any
	if !plan.ConnectionParametersSet.IsNull() && plan.ConnectionParametersSet.ValueString() != "" {
		err := json.Unmarshal([]byte(plan.ConnectionParametersSet.ValueString()), &connParamsSet)
		if err != nil {
			resp.Diagnostics.AddError("Failed to convert connection parameters set", err.Error())
			return
		}
	}

	connection, err := r.ConnectionsClient.UpdateConnection(ctx, plan.EnvironmentId.ValueString(), plan.Name.ValueString(), plan.Id.ValueString(), plan.DisplayName.ValueString(), connParams, connParamsSet)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	conectionState := ConvertFromConnectionDto(*connection)
	plan.Id = types.String(conectionState.Id)
	plan.DisplayName = types.String(conectionState.DisplayName)
	plan.Name = types.String(conectionState.Name)
	statuses := []attr.Value{}
	for _, status := range conectionState.Status {
		statuses = append(statuses, types.StringValue(status))
	}
	plan.Status = types.SetValueMust(types.StringType, statuses)

	if conectionState.ConnectionParameters == types.StringNull() {
		plan.ConnectionParameters = types.StringValue("")
	}
	if conectionState.ConnectionParametersSet == types.StringNull() {
		plan.ConnectionParametersSet = types.StringValue("")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.ConnectionsClient.DeleteConnection(ctx, state.EnvironmentId.ValueString(), state.Name.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
