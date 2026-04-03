// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package git_integration

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ resource.Resource = &SolutionGitBranchResource{}
var _ resource.ResourceWithImportState = &SolutionGitBranchResource{}

func NewSolutionGitBranchResource() resource.Resource {
	return &SolutionGitBranchResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "solution_git_branch",
		},
	}
}

func (r *SolutionGitBranchResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *SolutionGitBranchResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a solution-level Dataverse Git branch binding. This maps to the documented `sourcecontrolbranchconfiguration` Dataverse table and links a solution partition to a branch and folder beneath an environment Git integration.\n\nKnown limitation: the underlying Power Platform Git integration bootstrap currently requires delegated user principal authentication with Azure DevOps access. Service principal, app-only, and OIDC pipeline identities are not currently supported by the backing Dataverse Git integration flow.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the Dataverse source control branch configuration.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Environment ID of the Dataverse environment where the branch binding exists.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"git_integration_id": schema.StringAttribute{
				MarkdownDescription: "ID of the parent `powerplatform_environment_git_integration` resource.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"solution_id": schema.StringAttribute{
				MarkdownDescription: "ID of the existing `powerplatform_solution` resource to bind to the Git branch. This must use the provider solution ID format for the same environment.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"branch_name": schema.StringAttribute{
				MarkdownDescription: "Branch name to bind the solution partition to.",
				Required:            true,
			},
			"upstream_branch_name": schema.StringAttribute{
				MarkdownDescription: "Upstream branch name. When omitted, the provider will use the same value as `branch_name`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"root_folder_path": schema.StringAttribute{
				MarkdownDescription: "Repository folder path that stores the solution's files.",
				Required:            true,
			},
		},
	}
}

func (r *SolutionGitBranchResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		return
	}

	providerClient, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.GitIntegrationClient = newGitIntegrationClient(providerClient.Api)
}

func (r *SolutionGitBranchResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan SolutionGitBranchResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upstreamBranchName := plan.UpstreamBranchName.ValueString()
	if upstreamBranchName == "" {
		upstreamBranchName = plan.BranchName.ValueString()
	}

	solutionID := r.validateRemoteConfiguration(ctx, plan, "", &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	solutionDetails, err := r.GitIntegrationClient.GetUnmanagedSolutionByID(ctx, plan.EnvironmentID.ValueString(), solutionID)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when resolving solution metadata for %s", r.FullTypeName()), err.Error())
		return
	}

	createDTO := createSourceControlBranchConfigurationDto{
		Name:                             "",
		PartitionID:                      solutionID,
		BranchName:                       plan.BranchName.ValueString(),
		UpstreamBranchName:               upstreamBranchName,
		RootFolderPath:                   plan.RootFolderPath.ValueString(),
		SourceControlConfigurationBindID: buildSourceControlConfigurationBindPath(plan.GitIntegrationID.ValueString()),
	}

	created, err := r.GitIntegrationClient.CreateSolutionGitBranch(ctx, plan.EnvironmentID.ValueString(), solutionDetails.UniqueName, createDTO)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err.Error())
		return
	}

	state := convertSourceControlBranchConfigurationDtoToModel(plan.EnvironmentID.ValueString(), *created)
	state.Timeouts = plan.Timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SolutionGitBranchResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state SolutionGitBranchResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	solutionID, err := normalizeSolutionID(state.EnvironmentID.ValueString(), state.SolutionID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
		return
	}

	dto, err := r.GitIntegrationClient.FindSolutionGitBranchByPartition(ctx, state.EnvironmentID.ValueString(), state.GitIntegrationID.ValueString(), solutionID)
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
		return
	}

	newState := convertSourceControlBranchConfigurationDtoToModel(state.EnvironmentID.ValueString(), *dto)
	newState.Timeouts = state.Timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *SolutionGitBranchResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan SolutionGitBranchResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state SolutionGitBranchResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upstreamBranchName := plan.UpstreamBranchName.ValueString()
	if upstreamBranchName == "" {
		upstreamBranchName = plan.BranchName.ValueString()
	}

	solutionID := r.validateRemoteConfiguration(ctx, plan, state.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	updateDTO := updateSourceControlBranchConfigurationDto{
		Name:               solutionID,
		BranchName:         plan.BranchName.ValueString(),
		UpstreamBranchName: upstreamBranchName,
		RootFolderPath:     plan.RootFolderPath.ValueString(),
	}

	updated, err := r.GitIntegrationClient.UpdateSolutionGitBranch(ctx, plan.EnvironmentID.ValueString(), state.ID.ValueString(), plan.GitIntegrationID.ValueString(), solutionID, updateDTO)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
		return
	}

	newState := convertSourceControlBranchConfigurationDtoToModel(plan.EnvironmentID.ValueString(), *updated)
	newState.Timeouts = plan.Timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *SolutionGitBranchResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state SolutionGitBranchResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	solutionID, err := normalizeSolutionID(state.EnvironmentID.ValueString(), state.SolutionID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
		return
	}

	if err := r.GitIntegrationClient.DeleteSolutionGitBranch(ctx, state.EnvironmentID.ValueString(), state.ID.ValueString(), state.GitIntegrationID.ValueString(), solutionID); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
	}
}

func (r *SolutionGitBranchResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	idParts := strings.SplitN(req.ID, "/", 3)
	if len(idParts) != 3 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in format 'environment_id/git_integration_id/solution_id', got '%s'", req.ID),
		)
		return
	}

	solutionID := idParts[2]
	if !strings.Contains(solutionID, "_") {
		solutionID = buildSolutionReference(idParts[0], solutionID)
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("git_integration_id"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("solution_id"), solutionID)...)
}
