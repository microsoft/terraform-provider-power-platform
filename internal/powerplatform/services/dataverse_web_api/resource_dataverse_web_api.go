// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

func NewDataverseWebApiResource() resource.Resource {
	return &DataverseWebApiResource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_dataverse_web_api",
	}
}

type DataverseWebApiResource struct {
	DataRecordClient WebApiClient
	ProviderTypeName string
	TypeName         string
}

type DataverseWebApiResourceModel struct {
	Id            types.String                      `tfsdk:"id"`
	EnvironmentId types.String                      `tfsdk:"environment_id"`
	Create        *DataverseWebApiOperationResource `tfsdk:"create"`
	Update        *DataverseWebApiOperationResource `tfsdk:"update"`
	Delete        *DataverseWebApiOperationResource `tfsdk:"delete"`
	Read          *DataverseWebApiOperationResource `tfsdk:"read"`
	Output        types.Object                      `tfsdk:"output"`
}

type DataverseWebApiOperationResource struct {
	Method  types.String                             `tfsdk:"method"`
	Url     types.String                             `tfsdk:"url"`
	Body    types.String                             `tfsdk:"body"`
	Headers []DataverseWebApiOperationHeaderResource `tfsdk:"headers"`
}

type DataverseWebApiOperationHeaderResource struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func (r *DataverseWebApiResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

func (r *DataverseWebApiResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource to execute web api requests",
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
			"create": r.buildOperationSchema("Create operation", true),
			"update": r.buildOperationSchema("Update operation", false),
			"delete": r.buildOperationSchema("Delete operation", false),
			"read":   r.buildOperationSchema("Read operation", false),
			"output": schema.SingleNestedAttribute{
				MarkdownDescription: "Response body after executing the web api request",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"body": schema.StringAttribute{
						MarkdownDescription: "Response body after executing the web api request",
						Computed:            true,
					},
					"status": schema.Int64Attribute{
						MarkdownDescription: "Response status code after executing the web api request",
						Computed:            true,
					},
				},
			},
		},
	}
}

func (r *DataverseWebApiResource) buildOperationSchema(description string, isRequired bool) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: description,
		Required:            isRequired,
		Optional:            !isRequired,
		Attributes: map[string]schema.Attribute{
			"method": schema.StringAttribute{
				MarkdownDescription: "HTTP method",
				Required:            true,
				Optional:            false,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL of the web api",
				Required:            true,
				Optional:            false,
			},
			"body": schema.StringAttribute{
				MarkdownDescription: "Body of the request",
				Required:            false,
				Optional:            true,
			},
			"headers": schema.ListNestedAttribute{
				MarkdownDescription: "Headers of the request",
				Required:            false,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Header name",
							Required:            true,
							Optional:            false,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Header value",
							Required:            true,
							Optional:            false,
						},
					},
				},
			},
		},
	}
}

func (r *DataverseWebApiResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.DataRecordClient = NewWebApiClient(clientApi)

}

func (r *DataverseWebApiResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state DataverseWebApiResourceModel
	resp.State.Get(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	state.Id = types.StringValue(strconv.Itoa(int(time.Now().UnixMilli())))
	output := r.DataRecordClient.SendOperation(ctx, state.EnvironmentId.ValueString(), state.Create)
	state.Output = output

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE END: %s", r.TypeName))
}

func (r *DataverseWebApiResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *DataverseWebApiResourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if state.Read != nil {
		output := r.DataRecordClient.SendOperation(ctx, state.EnvironmentId.ValueString(), state.Read)
		state.Output = output
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE END: %s", r.TypeName))
}

func (r *DataverseWebApiResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *DataverseWebApiResourceModel

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE START: %s", r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Update != nil {
		output := r.DataRecordClient.SendOperation(ctx, plan.EnvironmentId.ValueString(), plan.Update)
		plan.Output = output
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE END: %s", r.TypeName))
}

func (r *DataverseWebApiResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *DataverseWebApiResourceModel
	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE START: %s", r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if state.Delete != nil {
		r.DataRecordClient.SendOperation(ctx, state.EnvironmentId.ValueString(), state.Delete)
	}

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE END: %s", r.TypeName))

}
