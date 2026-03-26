// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package git_integration

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
)

func (r *EnvironmentGitIntegrationResource) validateRemoteConfiguration(ctx context.Context, data EnvironmentGitIntegrationResourceModel, diags *diag.Diagnostics) {
	organizations, err := r.GitIntegrationClient.ListGitOrganizations(ctx, data.EnvironmentID.ValueString())
	if err != nil {
		diags.AddError("Client error when validating git organization", err.Error())
		return
	}

	if !containsOrganization(organizations, data.OrganizationName.ValueString()) {
		diags.AddAttributeError(
			path.Root("organization_name"),
			"Invalid organization_name",
			fmt.Sprintf("The Git organization `%s` was not returned by the Dataverse `gitorganizations` endpoint for environment `%s`.", data.OrganizationName.ValueString(), data.EnvironmentID.ValueString()),
		)
		return
	}

	projects, err := r.GitIntegrationClient.ListGitProjects(ctx, data.EnvironmentID.ValueString(), data.OrganizationName.ValueString())
	if err != nil {
		diags.AddError("Client error when validating git project", err.Error())
		return
	}

	if !containsProject(projects, data.ProjectName.ValueString()) {
		diags.AddAttributeError(
			path.Root("project_name"),
			"Invalid project_name",
			fmt.Sprintf("The Git project `%s` was not returned by the Dataverse `gitprojects` endpoint for organization `%s`.", data.ProjectName.ValueString(), data.OrganizationName.ValueString()),
		)
		return
	}

	repositories, err := r.GitIntegrationClient.ListGitRepositories(ctx, data.EnvironmentID.ValueString(), data.OrganizationName.ValueString(), data.ProjectName.ValueString())
	if err != nil {
		diags.AddError("Client error when validating git repository", err.Error())
		return
	}

	if !containsRepository(repositories, data.RepositoryName.ValueString()) {
		diags.AddAttributeError(
			path.Root("repository_name"),
			"Invalid repository_name",
			fmt.Sprintf("The Git repository `%s` was not returned by the Dataverse `gitrepositories` endpoint for the configured organization and project.", data.RepositoryName.ValueString()),
		)
	}
}

func (r *SolutionGitBranchResource) validateRemoteConfiguration(ctx context.Context, data SolutionGitBranchResourceModel, currentBranchID string, diags *diag.Diagnostics) string {
	solutionID, err := normalizeSolutionID(data.EnvironmentID.ValueString(), data.SolutionID.ValueString())
	if err != nil {
		diags.AddAttributeError(
			path.Root("solution_id"),
			"Invalid solution_id",
			err.Error(),
		)
		return ""
	}

	configuration, err := r.GitIntegrationClient.GetEnvironmentGitIntegration(ctx, data.EnvironmentID.ValueString(), data.GitIntegrationID.ValueString())
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			diags.AddAttributeError(
				path.Root("git_integration_id"),
				"Unknown git_integration_id",
				fmt.Sprintf("The Git integration `%s` was not found in environment `%s`.", data.GitIntegrationID.ValueString(), data.EnvironmentID.ValueString()),
			)
			return ""
		}

		diags.AddError("Client error when validating git integration", err.Error())
		return ""
	}

	if _, err := r.GitIntegrationClient.GetUnmanagedSolutionByID(ctx, data.EnvironmentID.ValueString(), solutionID); err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			diags.AddAttributeError(
				path.Root("solution_id"),
				"Unknown solution_id",
				fmt.Sprintf("The solution `%s` was not found as an unmanaged solution in environment `%s`.", data.SolutionID.ValueString(), data.EnvironmentID.ValueString()),
			)
			return ""
		}

		diags.AddError("Client error when validating solution", err.Error())
		return ""
	}

	scope, err := r.GitIntegrationClient.GetSourceControlIntegrationScope(ctx, data.EnvironmentID.ValueString())
	if err != nil {
		diags.AddError("Client error when validating source control integration scope", err.Error())
		return ""
	}

	if scope != scopeSolution {
		diags.AddAttributeError(
			path.Root("git_integration_id"),
			"Invalid git integration scope for solution binding",
			fmt.Sprintf("The parent Git integration `%s` uses scope `%s`. `powerplatform_solution_git_branch` requires the parent `powerplatform_environment_git_integration.scope` to be `Solution`.", data.GitIntegrationID.ValueString(), scope),
		)
		return ""
	}

	branches, err := r.GitIntegrationClient.ListGitBranches(ctx, data.EnvironmentID.ValueString(), configuration.OrganizationName, configuration.ProjectName, configuration.RepositoryName)
	if err != nil {
		diags.AddError("Client error when validating git branch", err.Error())
		return ""
	}

	if !containsBranch(branches, data.BranchName.ValueString()) {
		diags.AddAttributeError(
			path.Root("branch_name"),
			"Invalid branch_name",
			fmt.Sprintf("The branch `%s` was not returned by the Dataverse `gitbranches` endpoint for repository `%s`.", data.BranchName.ValueString(), configuration.RepositoryName),
		)
	}

	if !data.UpstreamBranchName.IsNull() && !data.UpstreamBranchName.IsUnknown() && data.UpstreamBranchName.ValueString() != "" && !containsBranch(branches, data.UpstreamBranchName.ValueString()) {
		diags.AddAttributeError(
			path.Root("upstream_branch_name"),
			"Invalid upstream_branch_name",
			fmt.Sprintf("The upstream branch `%s` was not returned by the Dataverse `gitbranches` endpoint for repository `%s`.", data.UpstreamBranchName.ValueString(), configuration.RepositoryName),
		)
	}

	existingBinding, err := r.GitIntegrationClient.FindSolutionGitBranchByPartition(ctx, data.EnvironmentID.ValueString(), data.GitIntegrationID.ValueString(), solutionID)
	if err != nil && !errors.Is(err, customerrors.ErrObjectNotFound) {
		diags.AddError("Client error when validating existing solution git branch", err.Error())
		return ""
	}
	if err == nil && existingBinding != nil && existingBinding.ID != currentBranchID {
		diags.AddAttributeError(
			path.Root("solution_id"),
			"Duplicate solution git branch binding",
			fmt.Sprintf("A Git branch binding already exists for solution `%s` under Git integration `%s`. Only one `powerplatform_solution_git_branch` is allowed per solution within the same environment Git integration.", data.SolutionID.ValueString(), data.GitIntegrationID.ValueString()),
		)
	}

	return solutionID
}

func normalizeSolutionID(environmentID, configuredValue string) (string, error) {
	if configuredValue == "" {
		return "", errors.New("the `solution_id` attribute must not be empty")
	}

	environmentPrefix, solutionID, found := strings.Cut(configuredValue, "_")
	if !found {
		return "", fmt.Errorf("the `solution_id` value `%s` must be the `id` exported by `powerplatform_solution` for environment `%s`", configuredValue, environmentID)
	}

	if environmentPrefix != environmentID {
		return "", fmt.Errorf("the `solution_id` value `%s` uses environment `%s`, but this resource is targeting environment `%s`", configuredValue, environmentPrefix, environmentID)
	}

	return solutionID, nil
}

func buildSolutionReference(environmentID, solutionID string) string {
	return fmt.Sprintf("%s_%s", environmentID, solutionID)
}

func containsOrganization(values []gitOrganizationDto, organizationName string) bool {
	for _, item := range values {
		if strings.EqualFold(item.OrganizationName, organizationName) {
			return true
		}
	}

	return false
}

func containsProject(values []gitProjectDto, projectName string) bool {
	for _, item := range values {
		if strings.EqualFold(item.ProjectName, projectName) {
			return true
		}
	}

	return false
}

func containsRepository(values []gitRepositoryDto, repositoryName string) bool {
	for _, item := range values {
		if strings.EqualFold(item.RepositoryName, repositoryName) {
			return true
		}
	}

	return false
}

func containsBranch(values []gitBranchDto, branchName string) bool {
	for _, item := range values {
		if strings.EqualFold(item.BranchName, branchName) {
			return true
		}
	}

	return false
}
