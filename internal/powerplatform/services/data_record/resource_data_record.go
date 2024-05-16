// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

func NewDataRecordResource() resource.Resource {
	return &DataRecordResource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_data_record",
	}
}

type DataRecordResource struct {
	DataRecordClient DataRecordClient
	ProviderTypeName string
	TypeName         string
}

type DataRecordResourceModel struct {
	Id               types.String  `tfsdk:"id"`
	EnvironmentId    types.String  `tfsdk:"environment_id"`
	TableLogicalName types.String  `tfsdk:"table_logical_name"`
	Columns          types.Dynamic `tfsdk:"columns"`
}

func (r *DataRecordResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

func (r *DataRecordResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "PowerPlatform Data Record Resource",
		MarkdownDescription: "Resource for managing PowerPlatform Data Record",
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
			"table_logical_name": schema.StringAttribute{
				Description: "Name of the data record table",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"columns": schema.DynamicAttribute{
				Description: "Columns of the data record table",
				Required:    true,
			},
		},
	}
}

func (r *DataRecordResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.DataRecordClient = NewDataRecordClient(clientApi)
}

func (r *DataRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *DataRecordResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	var plan DataRecordResourceModel
	resp.State.Get(ctx, &plan)

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	plan.Id = types.StringValue(plan.Id.ValueString())
	plan.EnvironmentId = types.StringValue(plan.EnvironmentId.ValueString())
	plan.TableLogicalName = types.StringValue(plan.TableLogicalName.ValueString())
	plan.Columns = types.DynamicValue(plan.Columns)

	var mapColumns map[string]interface{}

	jsonColumns, _ := json.Marshal(plan.Columns.String())
	unquotedJsonColumns, _ := strconv.Unquote(string(jsonColumns))
	json.Unmarshal([]byte(unquotedJsonColumns), &mapColumns)

	dr, err := r.DataRecordClient.ApplyDataRecord(ctx, plan.Id.ValueString(), plan.EnvironmentId.ValueString(), plan.TableLogicalName.ValueString(), mapColumns)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.ProviderTypeName), err.Error())
		return
	}

	plan.Id = types.StringValue(dr.Id)

	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID %s", plan.TableLogicalName.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *DataRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *DataRecordResourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.DataRecordClient.GetDataRecord(ctx, state.Id.ValueString(), state.EnvironmentId.ValueString(), state.TableLogicalName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.ProviderTypeName), err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("READ: %s_data_record with table_name %s", r.ProviderTypeName, state.TableLogicalName.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE END: %s", r.ProviderTypeName))
}

func (r *DataRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *DataRecordResourceModel

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state *DataRecordResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	plan.Id = types.StringValue(plan.Id.ValueString())
	plan.EnvironmentId = types.StringValue(plan.EnvironmentId.ValueString())
	plan.TableLogicalName = types.StringValue(plan.TableLogicalName.ValueString())
	plan.Columns = types.DynamicValue(plan.Columns)

	var mapColumns map[string]interface{}

	jsonColumns, _ := json.Marshal(plan.Columns.String())
	unquotedJsonColumns, _ := strconv.Unquote(string(jsonColumns))
	json.Unmarshal([]byte(unquotedJsonColumns), &mapColumns)

	dr, err := r.DataRecordClient.ApplyDataRecord(ctx, state.Id.ValueString(), plan.EnvironmentId.ValueString(), plan.TableLogicalName.ValueString(), mapColumns)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.ProviderTypeName), err.Error())
		return
	}

	plan.Id = types.StringValue(dr.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *DataRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *DataRecordResourceModel

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var mapColumns map[string]interface{}

	jsonColumns, _ := json.Marshal(state.Columns.String())
	unquotedJsonColumns, _ := strconv.Unquote(string(jsonColumns))
	json.Unmarshal([]byte(unquotedJsonColumns), &mapColumns)

	err := r.DataRecordClient.DeleteDataRecord(ctx, state.Id.ValueString(), state.EnvironmentId.ValueString(), state.TableLogicalName.ValueString(), mapColumns)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *DataRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
