// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package git_integration

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
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
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ resource.Resource = &EnvironmentGitIntegrationResource{}
var _ resource.ResourceWithValidateConfig = &EnvironmentGitIntegrationResource{}
var _ resource.ResourceWithImportState = &EnvironmentGitIntegrationResource{}

func NewEnvironmentGitIntegrationResource() resource.Resource {
	return &EnvironmentGitIntegrationResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment_git_integration",
		},
	}
}

func (r *EnvironmentGitIntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *EnvironmentGitIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the environment-level Dataverse Git repository binding. This maps to the documented `sourcecontrolconfiguration` Dataverse table and stores the repository connection metadata for an environment.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the Dataverse source control configuration.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Environment ID of the Dataverse environment where the Git repository binding will be created.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"git_provider": schema.StringAttribute{
				MarkdownDescription: "Git provider for the repository binding. Supported value is `AzureDevOps`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(gitProviderAzureDevOps),
				},
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "Source control integration scope for the environment. Use `Solution` for solution-level branch bindings and `Environment` for an environment-level binding. In `Environment` scope, the provider manages the root branch binding and proactively enables eligible visible unmanaged solutions in the environment while excluding platform-owned default solutions.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(scopeEnvironment, scopeSolution),
				},
			},
			"organization_name": schema.StringAttribute{
				MarkdownDescription: "Organization or owner name for the configured Git provider.",
				Required:            true,
			},
			"project_name": schema.StringAttribute{
				MarkdownDescription: "Project name for the Azure DevOps repository binding.",
				Required:            true,
			},
			"repository_name": schema.StringAttribute{
				MarkdownDescription: "Repository name that the environment will bind to.",
				Required:            true,
			},
		},
	}
}

func (r *EnvironmentGitIntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EnvironmentGitIntegrationResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var projectName types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("project_name"), &projectName)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if projectName.IsUnknown() {
		return
	}

	if projectName.IsNull() || projectName.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("project_name"),
			"Missing project_name for AzureDevOps",
			"The `project_name` attribute is required when `git_provider` is `AzureDevOps`.",
		)
	}
}

func (r *EnvironmentGitIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan EnvironmentGitIntegrationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createDTO := createSourceControlConfigurationDto{
		ID:               uuid.NewString(),
		Name:             "",
		OrganizationName: plan.OrganizationName.ValueString(),
		ProjectName:      plan.ProjectName.ValueString(),
		RepositoryName:   plan.RepositoryName.ValueString(),
		GitProvider:      gitProviderToInt(plan.GitProvider.ValueString()),
	}

	r.validateRemoteConfiguration(ctx, plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.GitIntegrationClient.SetSourceControlIntegrationScope(ctx, plan.EnvironmentID.ValueString(), plan.Scope.ValueString()); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when setting scope for %s", r.FullTypeName()), err.Error())
		return
	}

	created, err := r.GitIntegrationClient.CreateEnvironmentGitIntegration(ctx, plan.EnvironmentID.ValueString(), createDTO)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err.Error())
		return
	}

	created, err = r.GitIntegrationClient.WaitForEnvironmentGitIntegrationReady(ctx, plan.EnvironmentID.ValueString(), created.ID)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when waiting for %s to stabilize", r.FullTypeName()), err.Error())
		return
	}

	if plan.Scope.ValueString() == scopeEnvironment {
		if err := r.GitIntegrationClient.EnsureSolutionScopeRootBranch(ctx, plan.EnvironmentID.ValueString(), created.ID, created.OrganizationName, created.ProjectName, created.RepositoryName); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating the root Git binding for %s", r.FullTypeName()), err.Error())
			return
		}

		if err := r.GitIntegrationClient.EnableEnvironmentScopeSolutions(ctx, plan.EnvironmentID.ValueString()); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when enabling environment-scoped solutions for %s", r.FullTypeName()), err.Error())
			return
		}
	}

	state := convertSourceControlConfigurationDtoToModel(plan.EnvironmentID.ValueString(), plan.Scope.ValueString(), *created)
	state.Timeouts = plan.Timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EnvironmentGitIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state EnvironmentGitIntegrationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dto, err := r.GitIntegrationClient.GetEnvironmentGitIntegration(ctx, state.EnvironmentID.ValueString(), state.ID.ValueString())
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
		return
	}

	scope, err := r.GitIntegrationClient.GetSourceControlIntegrationScope(ctx, state.EnvironmentID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading scope for %s", r.FullTypeName()), err.Error())
		return
	}

	newState := convertSourceControlConfigurationDtoToModel(state.EnvironmentID.ValueString(), scope, *dto)
	newState.Timeouts = state.Timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *EnvironmentGitIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan EnvironmentGitIntegrationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state EnvironmentGitIntegrationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateDTO := updateSourceControlConfigurationDto{
		Name:             plan.RepositoryName.ValueString(),
		OrganizationName: plan.OrganizationName.ValueString(),
		ProjectName:      plan.ProjectName.ValueString(),
		RepositoryName:   plan.RepositoryName.ValueString(),
		GitProvider:      gitProviderToInt(plan.GitProvider.ValueString()),
	}

	r.validateRemoteConfiguration(ctx, plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.GitIntegrationClient.SetSourceControlIntegrationScope(ctx, plan.EnvironmentID.ValueString(), plan.Scope.ValueString()); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when setting scope for %s", r.FullTypeName()), err.Error())
		return
	}

	updated, err := r.GitIntegrationClient.UpdateEnvironmentGitIntegration(ctx, plan.EnvironmentID.ValueString(), state.ID.ValueString(), updateDTO)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
		return
	}

	updated, err = r.GitIntegrationClient.WaitForEnvironmentGitIntegrationReady(ctx, plan.EnvironmentID.ValueString(), updated.ID)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when waiting for %s to stabilize", r.FullTypeName()), err.Error())
		return
	}

	if plan.Scope.ValueString() == scopeEnvironment {
		if err := r.GitIntegrationClient.EnsureSolutionScopeRootBranch(ctx, plan.EnvironmentID.ValueString(), updated.ID, updated.OrganizationName, updated.ProjectName, updated.RepositoryName); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating the root Git binding for %s", r.FullTypeName()), err.Error())
			return
		}

		if err := r.GitIntegrationClient.EnableEnvironmentScopeSolutions(ctx, plan.EnvironmentID.ValueString()); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when enabling environment-scoped solutions for %s", r.FullTypeName()), err.Error())
			return
		}
	}

	newState := convertSourceControlConfigurationDtoToModel(plan.EnvironmentID.ValueString(), plan.Scope.ValueString(), *updated)
	newState.Timeouts = plan.Timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *EnvironmentGitIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state EnvironmentGitIntegrationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.GitIntegrationClient.DeleteEnvironmentGitIntegration(ctx, state.EnvironmentID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
	}
}

func (r *EnvironmentGitIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	configurations, err := r.GitIntegrationClient.ListEnvironmentGitIntegrations(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error when importing environment git integration",
			err.Error(),
		)
		return
	}

	if len(configurations) == 0 {
		resp.Diagnostics.AddError(
			"Environment git integration not found",
			fmt.Sprintf("No Dataverse source control configuration was found in environment '%s'. Import expects the environment ID of an environment that already has a Git integration.", req.ID),
		)
		return
	}

	if len(configurations) > 1 {
		resp.Diagnostics.AddError(
			"Multiple environment git integrations found",
			fmt.Sprintf("Expected exactly one Dataverse source control configuration in environment '%s', found %d.", req.ID, len(configurations)),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), configurations[0].ID)...)
}
