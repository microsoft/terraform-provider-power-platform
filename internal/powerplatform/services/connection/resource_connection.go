// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"encoding/json"
	"fmt"

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

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
)

var _ resource.Resource = &ConnectionResource{}
var _ resource.ResourceWithImportState = &ConnectionResource{}

func NewConnectionResource() resource.Resource {
	return &ConnectionResource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_connection",
	}
}

type ConnectionResource struct {
	ConnectionsClient ConnectionsClient
	ProviderTypeName  string
	TypeName          string
}

type ConnectionResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	EnvironmentId           types.String `tfsdk:"environment_id"`
	DisplayName             types.String `tfsdk:"display_name"`
	Status                  types.Set    `tfsdk:"status"`
	ConnectionParameters    types.String `tfsdk:"connection_parameters"`
	ConnectionParametersSet types.String `tfsdk:"connection_parameters_set"`
}

func (r *ConnectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

func (r *ConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a [Connection](https://learn.microsoft.com/en-us/power-apps/maker/canvas-apps/add-manage-connections). A connection in Power Platform serves as a means to integrate external data sources and services with your Power Platform apps, flows, and other solutions. It acts as a bridge, facilitating secure communication between your solutions and various external systems.",
		Attributes: map[string]schema.Attribute{
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

func (d *ConnectionResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("connection_parameters"),
			path.MatchRoot("connection_parameters_set"),
		),
	}
}

func (r *ConnectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.ConnectionsClient = NewConnectionsClient(clientApi)
}

func (r *ConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *ConnectionResourceModel

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	connectionToCreate := ConnectionToCreateDto{
		Properties: ConnectionToCreatePropertiesDto{
			DisplayName: plan.DisplayName.ValueString(),
			Environment: ConnectionToCreateEnvironmentDto{
				Name: plan.EnvironmentId.ValueString(),
				Id:   fmt.Sprintf("/providers/Microsoft.PowerApps/environments/%s", plan.EnvironmentId.ValueString()),
			},
		},
	}

	if !plan.ConnectionParameters.IsNull() && plan.ConnectionParameters.ValueString() != "" {
		var params map[string]interface{} = nil
		err := json.Unmarshal([]byte(plan.ConnectionParameters.ValueString()), &params)
		if err != nil {
			resp.Diagnostics.AddError("Failed to convert connection parameters", err.Error())
			return
		}
		connectionToCreate.Properties.ConnectionParameters = params
	}
	if !plan.ConnectionParametersSet.IsNull() && plan.ConnectionParametersSet.ValueString() != "" {
		var params map[string]interface{} = nil
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

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *ConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *ConnectionResourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	connection, err := r.ConnectionsClient.GetConnection(ctx, state.EnvironmentId.ValueString(), state.Name.ValueString(), state.Id.ValueString())
	if err != nil {
		if powerplatform_helpers.Code(err) == powerplatform_helpers.ERROR_OBJECT_NOT_FOUND {
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
			return
		}
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
	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE END: %s", r.TypeName))
}

func (r *ConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *ConnectionResourceModel

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state *ConnectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var connParams map[string]interface{} = nil
	if !plan.ConnectionParameters.IsNull() && plan.ConnectionParameters.ValueString() != "" {

		err := json.Unmarshal([]byte(plan.ConnectionParameters.ValueString()), &connParams)
		if err != nil {
			resp.Diagnostics.AddError("Failed to convert connection parameters", err.Error())
			return
		}
	}

	var connParamsSet map[string]interface{} = nil
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

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *ConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *ConnectionResourceModel

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.ConnectionsClient.DeleteConnection(ctx, state.EnvironmentId.ValueString(), state.Name.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *ConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
