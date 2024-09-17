// Licensed under the MIT license.
// Copyright (c) Microsoft Corporation.

package environment_groups

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/api"
)

var _ resource.Resource = &EnvironmentGroupResource{}
var _ resource.ResourceWithImportState = &EnvironmentGroupResource{}

func NewEnvironmentGroupResource() resource.Resource {
	return &EnvironmentGroupResource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_environment_group",
	}
}

type EnvironmentGroupResource struct {
	EnvironmentGroupClient EnvironmentGroupClient
	ProviderTypeName       string
	TypeName               string
}

type EnvironmentGroupResourceModel struct {
	Id          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
}

func (r *EnvironmentGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

func (r *EnvironmentGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource manages an [Environment Group](https://learn.microsoft.com/en-us/power-platform/admin/environment-groups).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique id of the environment group",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the environment group",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Display name of the environment group",
				Required:            true,
			},
		},
	}
}

func (r *EnvironmentGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := req.ProviderData.(*api.ProviderClient).Api
	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.EnvironmentGroupClient = NewEnvironmentGroupClient(client)
}

// Read function.
func (r *EnvironmentGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", r.TypeName))

	state := EnvironmentGroupResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	environmentGroup, err := r.EnvironmentGroupClient.GetEnvironmentGroup(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	if environmentGroup == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Id = types.StringValue(environmentGroup.Id)
	state.DisplayName = types.StringValue(environmentGroup.DisplayName)
	state.Description = types.StringValue(environmentGroup.Description)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE END: %s", r.TypeName))
}

// Create function.
func (r *EnvironmentGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE START: %s", r.TypeName))

	var plan *EnvironmentGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	environmentGroupToCreate := EnvironmentGroupDto{
		DisplayName: plan.DisplayName.ValueString(),
		Description: plan.Description.ValueString(),
	}

	eg, err := r.EnvironmentGroupClient.CreateEnvironmentGroup(ctx, environmentGroupToCreate)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	state := EnvironmentGroupResourceModel{}
	state.Id = types.StringValue(eg.Id)
	state.DisplayName = types.StringValue(eg.DisplayName)
	state.Description = types.StringValue(eg.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE END: %s", r.TypeName))
}

// Update function.
func (r *EnvironmentGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE START: %s", r.TypeName))

	var plan *EnvironmentGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *EnvironmentGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	environmentGroupToUpdate := EnvironmentGroupDto{
		Id:          plan.Id.ValueString(),
		DisplayName: plan.DisplayName.ValueString(),
		Description: plan.Description.ValueString(),
	}

	eg, err := r.EnvironmentGroupClient.UpdateEnvironmentGroup(ctx, state.Id.ValueString(), environmentGroupToUpdate)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	state.Id = types.StringValue(eg.Id)
	state.DisplayName = types.StringValue(eg.DisplayName)
	state.Description = types.StringValue(eg.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE END: %s", r.TypeName))
}

// Delete function
func (r *EnvironmentGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE START: %s", r.TypeName))

	var state *EnvironmentGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.EnvironmentGroupClient.DeleteEnvironmentGroup(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE END: %s", r.TypeName))
}

// ImportState function
func (r *EnvironmentGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, fmt.Sprintf("IMPORT STATE RESOURCE START: %s", r.TypeName))

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	tflog.Debug(ctx, fmt.Sprintf("IMPORT STATE RESOURCE END: %s", r.TypeName))
}
