// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package git_integration

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

const (
	gitProviderAzureDevOps                = "AzureDevOps"
	scopeEnvironment                      = "Environment"
	scopeSolution                         = "Solution"
	rootPartitionID                       = "00000000-0000-0000-0000-000000000000"
	rootFolderPath                        = "dataverse"
	commonDataServicesDefaultSolutionID   = "00000001-0000-0000-0001-00000000009b"
	activeSolutionID                      = "fd140aae-4df4-11dd-bd17-0019b9312238"
	defaultSolutionID                     = "fd140aaf-4df4-11dd-bd17-0019b9312238"
	commonDataServicesDefaultSolutionName = "common data services default solution"
	defaultSolutionName                   = "default solution"
)

type EnvironmentGitIntegrationResource struct {
	helpers.TypeInfo
	GitIntegrationClient client
}

type EnvironmentGitIntegrationResourceModel struct {
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
	ID               types.String   `tfsdk:"id"`
	EnvironmentID    types.String   `tfsdk:"environment_id"`
	GitProvider      types.String   `tfsdk:"git_provider"`
	Scope            types.String   `tfsdk:"scope"`
	OrganizationName types.String   `tfsdk:"organization_name"`
	ProjectName      types.String   `tfsdk:"project_name"`
	RepositoryName   types.String   `tfsdk:"repository_name"`
}

type SolutionGitBranchResource struct {
	helpers.TypeInfo
	GitIntegrationClient client
}

type SolutionGitBranchResourceModel struct {
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
	ID                 types.String   `tfsdk:"id"`
	EnvironmentID      types.String   `tfsdk:"environment_id"`
	GitIntegrationID   types.String   `tfsdk:"git_integration_id"`
	SolutionID         types.String   `tfsdk:"solution_id"`
	BranchName         types.String   `tfsdk:"branch_name"`
	UpstreamBranchName types.String   `tfsdk:"upstream_branch_name"`
	RootFolderPath     types.String   `tfsdk:"root_folder_path"`
}

func gitProviderToInt(value string) int {
	switch strings.TrimSpace(value) {
	case "", gitProviderAzureDevOps:
		return 0
	default:
		return -1
	}
}

func gitProviderFromInt(value int) string {
	switch value {
	case 0:
		return gitProviderAzureDevOps
	default:
		return ""
	}
}

func convertSourceControlConfigurationDtoToModel(environmentID, scope string, dto sourceControlConfigurationDto) EnvironmentGitIntegrationResourceModel {
	projectName := types.StringNull()
	if dto.ProjectName != "" {
		projectName = types.StringValue(dto.ProjectName)
	}

	return EnvironmentGitIntegrationResourceModel{
		ID:               types.StringValue(dto.ID),
		EnvironmentID:    types.StringValue(environmentID),
		GitProvider:      types.StringValue(gitProviderFromInt(dto.GitProvider)),
		Scope:            types.StringValue(scope),
		OrganizationName: types.StringValue(dto.OrganizationName),
		ProjectName:      projectName,
		RepositoryName:   types.StringValue(dto.RepositoryName),
	}
}

func convertSourceControlBranchConfigurationDtoToModel(environmentID string, dto sourceControlBranchConfigurationDto) SolutionGitBranchResourceModel {
	return SolutionGitBranchResourceModel{
		ID:                 types.StringValue(dto.ID),
		EnvironmentID:      types.StringValue(environmentID),
		GitIntegrationID:   types.StringValue(dto.SourceControlConfiguration),
		SolutionID:         types.StringValue(buildSolutionReference(environmentID, dto.PartitionID)),
		BranchName:         types.StringValue(dto.BranchName),
		UpstreamBranchName: types.StringValue(dto.UpstreamBranchName),
		RootFolderPath:     types.StringValue(dto.RootFolderPath),
	}
}
