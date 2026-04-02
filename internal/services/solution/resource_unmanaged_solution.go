// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package solution

import (
	"context"
	"errors"
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
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ resource.Resource = &UnmanagedSolutionResource{}
var _ resource.ResourceWithImportState = &UnmanagedSolutionResource{}

func NewUnmanagedSolutionResource() resource.Resource {
	return &UnmanagedSolutionResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "unmanaged_solution",
		},
	}
}

func (r *UnmanagedSolutionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *UnmanagedSolutionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing unmanaged Dataverse solutions as first-class solution records without coupling lifecycle to solution ZIP import operations.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the unmanaged solution.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Id of the Dataverse-enabled environment containing the unmanaged solution.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"uniquename": schema.StringAttribute{
				MarkdownDescription: "Unique name of the solution. This is the stable solution identity in Dataverse.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the solution.",
				Required:            true,
			},
			"publisher_id": schema.StringAttribute{
				MarkdownDescription: "Existing Dataverse publisher id that owns the solution.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description of the unmanaged solution.",
				Optional:            true,
			},
		},
	}
}

func (r *UnmanagedSolutionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
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

	r.SolutionClient = NewSolutionClient(client.Api)
}

func (r *UnmanagedSolutionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan UnmanagedSolutionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dvExists, err := r.SolutionClient.DataverseExists(ctx, plan.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), err.Error())
		return
	}
	if !dvExists {
		resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), "")
		return
	}

	solution, err := r.SolutionClient.CreateUnmanagedSolution(
		ctx,
		plan.EnvironmentId.ValueString(),
		plan.UniqueName.ValueString(),
		plan.DisplayName.ValueString(),
		plan.PublisherId.ValueString(),
		plan.Description.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err.Error())
		return
	}
	if err := validateUnmanagedSolution(solution, r.FullTypeName()); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err.Error())
		return
	}

	setUnmanagedSolutionState(&plan, solution)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UnmanagedSolutionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state UnmanagedSolutionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	solution, err := r.SolutionClient.GetSolutionById(ctx, state.EnvironmentId.ValueString(), state.Id.ValueString())
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
		return
	}
	if err := validateUnmanagedSolution(solution, r.FullTypeName()); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
		return
	}

	setUnmanagedSolutionState(&state, solution)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UnmanagedSolutionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan UnmanagedSolutionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state UnmanagedSolutionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	solution, err := r.SolutionClient.UpdateUnmanagedSolution(
		ctx,
		state.EnvironmentId.ValueString(),
		state.Id.ValueString(),
		plan.DisplayName.ValueString(),
		plan.Description.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
		return
	}
	if err := validateUnmanagedSolution(solution, r.FullTypeName()); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
		return
	}

	plan.EnvironmentId = state.EnvironmentId
	plan.UniqueName = state.UniqueName
	plan.PublisherId = state.PublisherId
	setUnmanagedSolutionState(&plan, solution)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UnmanagedSolutionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state UnmanagedSolutionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.EnvironmentId.IsNull() || state.Id.IsNull() {
		return
	}

	err := r.SolutionClient.DeleteSolution(ctx, state.EnvironmentId.ValueString(), state.Id.ValueString())
	if err != nil && !errors.Is(err, customerrors.ErrObjectNotFound) {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
	}
}

func (r *UnmanagedSolutionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	idParts := splitSolutionCompositeID(req.ID)
	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in format 'environment_id_solution_id', got '%s'", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), idParts[0])...)
}

func setUnmanagedSolutionState(model *UnmanagedSolutionResourceModel, solution *SolutionDto) {
	model.Id = types.StringValue(solution.Id)
	model.UniqueName = types.StringValue(solution.Name)
	model.DisplayName = types.StringValue(solution.DisplayName)
	model.PublisherId = types.StringValue(solution.PublisherId)

	if solution.Description == "" {
		model.Description = types.StringNull()
	} else {
		model.Description = types.StringValue(solution.Description)
	}
}

func validateUnmanagedSolution(solution *SolutionDto, typeName string) error {
	if solution != nil && solution.IsManaged {
		return fmt.Errorf("solution '%s' is managed and cannot be used with %s", solution.Name, typeName)
	}

	return nil
}

func splitSolutionCompositeID(id string) []string {
	return strings.Split(id, "_")
}
