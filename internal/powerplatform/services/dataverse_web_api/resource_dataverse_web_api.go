// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	Method        types.String `tfsdk:"method"`
	Url           types.String `tfsdk:"url"`
	Body          types.String `tfsdk:"body"`
	Output        types.String `tfsdk:"output"`
}

func (r *DataverseWebApiResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

func (r *DataverseWebApiResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "",
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
			"method": schema.StringAttribute{
				MarkdownDescription: "HTTP method",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("GET", "POST", "PUT", "PATCH", "DELETE"),
				},
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL of the web api",
				Required:            true,
			},
			"body": schema.StringAttribute{
				MarkdownDescription: "Body of the request",
				Required:            false,
			},
			"output": schema.StringAttribute{
				MarkdownDescription: "Response body after executing the web api request",
				Computed:            true,
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

	state.Id = types.StringValue(state.Id.ValueString())
	state.EnvironmentId = types.StringValue(state.EnvironmentId.ValueString())
	state.Method = types.StringValue(state.Method.ValueString())
	state.Url = types.StringValue(state.Url.ValueString())
	state.Body = types.StringValue(state.Body.ValueString())

	_, err := r.DataRecordClient.ExecuteWebApiRequest(ctx, state.EnvironmentId, state.Url, state.Method, state.Body)
	if err != nil {
		resp.Diagnostics.AddError("Failed to execute web api request", err.Error())
	}

}

func (r *DataverseWebApiResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *DataverseWebApiResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *DataverseWebApiResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
