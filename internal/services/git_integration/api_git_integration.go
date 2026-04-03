// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package git_integration

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/solution"
)

const (
	sourceControlBranchConfigurationStatusActive   = 0
	sourceControlBranchConfigurationStatusInactive = 1
	stabilizedReadCountRequired                    = 2
)

func newGitIntegrationClient(apiClient *api.Client) client {
	return client{
		Api:               apiClient,
		EnvironmentClient: environment.NewEnvironmentClient(apiClient),
		SolutionClient:    solution.NewSolutionClient(apiClient),
	}
}

type client struct {
	Api               *api.Client
	EnvironmentClient environment.Client
	SolutionClient    solution.Client
}

func (c *client) GetEnvironmentGitIntegration(ctx context.Context, environmentID, configurationID string) (*sourceControlConfigurationDto, error) {
	environmentHost, err := c.EnvironmentClient.GetEnvironmentHostById(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	apiURL := helpers.BuildDataverseApiUrl(environmentHost, fmt.Sprintf("/api/data/v9.0/sourcecontrolconfigurations(%s)", configurationID), nil)

	var dto sourceControlConfigurationDto
	resp, err := c.Api.Execute(ctx, nil, http.MethodGet, apiURL, nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &dto)
	if err != nil {
		return nil, err
	}
	if err := c.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if resp.HttpResponse.StatusCode == http.StatusNotFound {
		return nil, customerrors.WrapIntoProviderError(nil, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), "source control configuration not found")
	}

	return &dto, nil
}

func (c *client) ListEnvironmentGitIntegrations(ctx context.Context, environmentID string) ([]sourceControlConfigurationDto, error) {
	environmentHost, err := c.EnvironmentClient.GetEnvironmentHostById(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	var dto sourceControlConfigurationArrayDto
	resp, err := c.Api.Execute(ctx, nil, http.MethodGet, helpers.BuildDataverseApiUrl(environmentHost, "/api/data/v9.0/sourcecontrolconfigurations", nil), nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &dto)
	if err != nil {
		return nil, err
	}
	if err := c.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if err := c.Api.HandleNotFoundResponse(resp); err != nil {
		return nil, err
	}

	return dto.Value, nil
}

func (c *client) CreateEnvironmentGitIntegration(ctx context.Context, environmentID string, dto createSourceControlConfigurationDto) (*sourceControlConfigurationDto, error) {
	environmentHost, err := c.EnvironmentClient.GetEnvironmentHostById(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	apiURL := helpers.BuildDataverseApiUrl(environmentHost, "/api/data/v9.0/sourcecontrolconfigurations", nil)
	resp, err := c.Api.Execute(ctx, nil, http.MethodPost, apiURL, nil, dto, []int{http.StatusNoContent, http.StatusCreated, http.StatusForbidden}, nil)
	if err != nil {
		return nil, err
	}
	if err := c.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}

	for {
		configuration, err := c.GetEnvironmentGitIntegration(ctx, environmentID, dto.ID)
		if err == nil {
			return configuration, nil
		}
		if !errors.Is(err, customerrors.ErrObjectNotFound) {
			return nil, err
		}
		if err := c.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
			return nil, err
		}
	}
}

func (c *client) WaitForEnvironmentGitIntegrationReady(ctx context.Context, environmentID, configurationID string) (*sourceControlConfigurationDto, error) {
	stableReads := 0

	for {
		configuration, err := c.GetEnvironmentGitIntegration(ctx, environmentID, configurationID)
		if err != nil {
			if !errors.Is(err, customerrors.ErrObjectNotFound) {
				return nil, err
			}
			stableReads = 0
		} else {
			stableReads++
			if stableReads >= stabilizedReadCountRequired {
				return configuration, nil
			}
		}

		if err := c.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
			return nil, err
		}
	}
}

func (c *client) EnsureSolutionScopeRootBranch(ctx context.Context, environmentID, configurationID, organizationName, projectName, repositoryName string) error {
	defaultBranch, err := c.GetGitRepositoryDefaultBranch(ctx, environmentID, organizationName, projectName, repositoryName)
	if err != nil {
		return err
	}

	existingBranch, err := c.lookupAnySolutionGitBranchByPartition(ctx, environmentID, configurationID, rootPartitionID)
	if err == nil {
		if existingBranch.StatusCode == sourceControlBranchConfigurationStatusActive &&
			strings.EqualFold(existingBranch.BranchName, defaultBranch) &&
			strings.EqualFold(existingBranch.UpstreamBranchName, defaultBranch) &&
			strings.EqualFold(existingBranch.RootFolderPath, rootFolderPath) {
			_, err = c.waitForBranchConfigurationState(ctx, environmentID, configurationID, rootPartitionID, defaultBranch, defaultBranch, rootFolderPath)
			return err
		}

		if existingBranch.StatusCode == sourceControlBranchConfigurationStatusActive {
			_, err = c.UpdateSolutionGitBranch(ctx, environmentID, existingBranch.ID, configurationID, rootPartitionID, updateSourceControlBranchConfigurationDto{
				BranchName:         defaultBranch,
				UpstreamBranchName: defaultBranch,
				RootFolderPath:     rootFolderPath,
			})
			return err
		}
	}
	if err != nil && !errors.Is(err, customerrors.ErrObjectNotFound) {
		return err
	}

	_, err = c.CreateBranchConfiguration(ctx, environmentID, createSourceControlBranchConfigurationDto{
		ID:                               "",
		PartitionID:                      rootPartitionID,
		BranchName:                       defaultBranch,
		UpstreamBranchName:               defaultBranch,
		RootFolderPath:                   rootFolderPath,
		SourceControlConfigurationBindID: fmt.Sprintf("/sourcecontrolconfigurations(%s)", configurationID),
	}, configurationID)
	return err
}

func (c *client) UpdateEnvironmentGitIntegration(ctx context.Context, environmentID, configurationID string, dto updateSourceControlConfigurationDto) (*sourceControlConfigurationDto, error) {
	environmentHost, err := c.EnvironmentClient.GetEnvironmentHostById(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	apiURL := helpers.BuildDataverseApiUrl(environmentHost, fmt.Sprintf("/api/data/v9.0/sourcecontrolconfigurations(%s)", configurationID), nil)
	resp, err := c.Api.Execute(ctx, nil, http.MethodPatch, apiURL, nil, dto, []int{http.StatusNoContent, http.StatusForbidden, http.StatusNotFound}, nil)
	if err != nil {
		return nil, err
	}
	if err := c.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if resp.HttpResponse.StatusCode == http.StatusNotFound {
		return nil, customerrors.WrapIntoProviderError(nil, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), "source control configuration not found")
	}

	return c.GetEnvironmentGitIntegration(ctx, environmentID, configurationID)
}

func (c *client) DeleteEnvironmentGitIntegration(ctx context.Context, environmentID, configurationID string) error {
	_, err := c.lookupAnySolutionGitBranchByPartition(ctx, environmentID, configurationID, rootPartitionID)
	if err != nil && !errors.Is(err, customerrors.ErrObjectNotFound) {
		return err
	}
	if err == nil {
		if err := c.DeleteSolutionGitBranch(ctx, environmentID, configurationID, rootPartitionID); err != nil {
			return err
		}
	}

	environmentHost, err := c.EnvironmentClient.GetEnvironmentHostById(ctx, environmentID)
	if err != nil {
		return err
	}

	apiURL := helpers.BuildDataverseApiUrl(environmentHost, fmt.Sprintf("/api/data/v9.0/sourcecontrolconfigurations(%s)", configurationID), nil)
	resp, err := c.Api.Execute(ctx, nil, http.MethodDelete, apiURL, nil, nil, []int{http.StatusNoContent, http.StatusNotFound, http.StatusForbidden}, nil)
	if err != nil {
		// Dataverse implicitly removes the parent config after the last solution binding disconnects,
		// but a direct DELETE can still return the legacy "can't be deleted" error afterwards.
		if strings.Contains(err.Error(), "Existing source control configurations can't be deleted.") {
			_, readErr := c.GetEnvironmentGitIntegration(ctx, environmentID, configurationID)
			if errors.Is(readErr, customerrors.ErrObjectNotFound) {
				return nil
			}
		}
		return err
	}
	return c.Api.HandleForbiddenResponse(resp)
}

func (c *client) ListEnvironmentScopeSolutions(ctx context.Context, environmentID string) ([]unmanagedSolutionDto, error) {
	environmentHost, err := c.EnvironmentClient.GetEnvironmentHostById(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Add("$filter", "(ismanaged eq false and isvisible eq true)")
	values.Add("$select", "solutionid,uniquename,friendlyname,ismanaged,isvisible,enabledforsourcecontrolintegration,version")

	var solutions unmanagedSolutionArrayDto
	resp, err := c.Api.Execute(ctx, nil, http.MethodGet, helpers.BuildDataverseApiUrl(environmentHost, "/api/data/v9.2/solutions", values), nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &solutions)
	if err != nil {
		return nil, err
	}
	if err := c.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if err := c.Api.HandleNotFoundResponse(resp); err != nil {
		return nil, err
	}

	filtered := make([]unmanagedSolutionDto, 0, len(solutions.Value))
	for _, solutionRow := range solutions.Value {
		if isEnvironmentScopeCandidateSolution(solutionRow) {
			filtered = append(filtered, solutionRow)
		}
	}

	return filtered, nil
}

func (c *client) GetSolutionGitBranch(ctx context.Context, environmentID, branchID string) (*sourceControlBranchConfigurationDto, error) {
	environmentHost, err := c.EnvironmentClient.GetEnvironmentHostById(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	apiURL := helpers.BuildDataverseApiUrl(environmentHost, fmt.Sprintf("/api/data/v9.0/sourcecontrolbranchconfigurations(%s)", branchID), nil)

	var dto sourceControlBranchConfigurationDto
	resp, err := c.Api.Execute(ctx, nil, http.MethodGet, apiURL, nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &dto)
	if err != nil {
		return nil, err
	}
	if err := c.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if resp.HttpResponse.StatusCode == http.StatusNotFound {
		return nil, customerrors.WrapIntoProviderError(nil, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), "source control branch configuration not found")
	}

	return &dto, nil
}

func (c *client) CreateSolutionGitBranch(ctx context.Context, environmentID, solutionUniqueName string, dto createSourceControlBranchConfigurationDto) (*sourceControlBranchConfigurationDto, error) {
	configurationID := strings.TrimPrefix(strings.TrimSuffix(dto.SourceControlConfigurationBindID, ")"), "/sourcecontrolconfigurations(")
	if _, err := c.CreateBranchConfiguration(ctx, environmentID, dto, configurationID); err != nil {
		return nil, err
	}

	configuration, err := c.GetEnvironmentGitIntegration(ctx, environmentID, configurationID)
	if err != nil {
		return nil, err
	}

	if err := c.EnsureSolutionScopeRootBranch(ctx, environmentID, configurationID, configuration.OrganizationName, configuration.ProjectName, configuration.RepositoryName); err != nil {
		return nil, err
	}

	return c.waitForBranchConfigurationState(ctx, environmentID, configurationID, dto.PartitionID, dto.BranchName, dto.UpstreamBranchName, dto.RootFolderPath, solutionUniqueName)
}

func (c *client) CreateBranchConfiguration(ctx context.Context, environmentID string, dto createSourceControlBranchConfigurationDto, configurationID string) (*sourceControlBranchConfigurationDto, error) {
	environmentHost, err := c.EnvironmentClient.GetEnvironmentHostById(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	apiURL := helpers.BuildDataverseApiUrl(environmentHost, "/api/data/v9.0/sourcecontrolbranchconfigurations", nil)
	resp, err := c.Api.Execute(ctx, nil, http.MethodPost, apiURL, nil, dto, []int{http.StatusNoContent, http.StatusCreated, http.StatusForbidden}, nil)
	if err != nil {
		return nil, err
	}
	if err := c.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}

	return c.waitForBranchConfigurationState(ctx, environmentID, configurationID, dto.PartitionID, dto.BranchName, dto.UpstreamBranchName, dto.RootFolderPath)
}

func (c *client) PreValidateGitComponents(ctx context.Context, environmentID, solutionUniqueName string) (bool, error) {
	environmentHost, err := c.EnvironmentClient.GetEnvironmentHostById(ctx, environmentID)
	if err != nil {
		return false, err
	}

	apiURL := helpers.BuildDataverseApiUrl(environmentHost, "/api/data/v9.0/PreValidateGitComponents", nil)
	var response preValidateGitComponentsResponseDto
	resp, err := c.Api.Execute(ctx, nil, http.MethodPost, apiURL, nil, preValidateGitComponentsRequestDto{
		SolutionUniqueName: solutionUniqueName,
	}, []int{http.StatusOK, http.StatusBadRequest, http.StatusForbidden, http.StatusNotFound}, &response)
	if err != nil {
		return false, err
	}
	if err := c.Api.HandleForbiddenResponse(resp); err != nil {
		return false, err
	}
	if resp.HttpResponse.StatusCode == http.StatusNotFound || resp.HttpResponse.StatusCode == http.StatusBadRequest {
		return false, nil
	}

	return strings.TrimSpace(response.ValidationMessages) == "", nil
}

func (c *client) waitForBranchConfigurationState(ctx context.Context, environmentID, configurationID, partitionID, branchName, upstreamBranchName, rootFolder string, solutionUniqueName ...string) (*sourceControlBranchConfigurationDto, error) {
	stableReads := 0
	var uniqueName string
	if len(solutionUniqueName) > 0 {
		uniqueName = solutionUniqueName[0]
	}

	for {
		branch, err := c.FindSolutionGitBranchByPartition(ctx, environmentID, configurationID, partitionID)
		if err == nil {
			if strings.EqualFold(branch.PartitionID, partitionID) &&
				strings.EqualFold(branch.BranchName, branchName) &&
				strings.EqualFold(branch.UpstreamBranchName, upstreamBranchName) &&
				strings.EqualFold(branch.RootFolderPath, rootFolder) &&
				strings.EqualFold(branch.SourceControlConfiguration, configurationID) &&
				branch.StatusCode == sourceControlBranchConfigurationStatusActive &&
				branch.BranchSyncedCommitID != "" &&
				branch.UpstreamBranchSyncedCommit != "" {
				if uniqueName != "" {
					preValidated, err := c.PreValidateGitComponents(ctx, environmentID, uniqueName)
					if err != nil {
						return nil, err
					}
					if !preValidated {
						stableReads = 0
						goto sleep
					}
				}

				stableReads++
				if stableReads >= stabilizedReadCountRequired {
					return branch, nil
				}
			} else {
				stableReads = 0
			}
		} else if !errors.Is(err, customerrors.ErrObjectNotFound) {
			return nil, err
		}

	sleep:
		if err := c.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
			return nil, err
		}
	}
}

func (c *client) UpdateSolutionGitBranch(ctx context.Context, environmentID, branchID, configurationID, partitionID string, dto updateSourceControlBranchConfigurationDto) (*sourceControlBranchConfigurationDto, error) {
	environmentHost, err := c.EnvironmentClient.GetEnvironmentHostById(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	apiURL := helpers.BuildDataverseApiUrl(environmentHost, buildSourceControlBranchConfigurationCompositeKeyPath(branchID, partitionID), nil)
	headers := http.Header{}
	headers.Set("If-Match", "*")
	resp, err := c.Api.Execute(ctx, nil, http.MethodPatch, apiURL, headers, dto, []int{http.StatusNoContent, http.StatusForbidden, http.StatusNotFound}, nil)
	if err != nil {
		return nil, err
	}
	if err := c.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if resp.HttpResponse.StatusCode == http.StatusNotFound {
		return nil, customerrors.WrapIntoProviderError(nil, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), "source control branch configuration not found")
	}

	return c.FindSolutionGitBranchByPartition(ctx, environmentID, configurationID, partitionID)
}

func (c *client) DeleteSolutionGitBranch(ctx context.Context, environmentID, configurationID, partitionID string) error {
	existingBranch, err := c.lookupAnySolutionGitBranchByPartition(ctx, environmentID, configurationID, partitionID)
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			return nil
		}
		return err
	}

	if existingBranch.StatusCode == sourceControlBranchConfigurationStatusInactive {
		return c.waitForSolutionGitBranchRemoval(ctx, environmentID, existingBranch.ID, configurationID, partitionID)
	}

	environmentHost, err := c.EnvironmentClient.GetEnvironmentHostById(ctx, environmentID)
	if err != nil {
		return err
	}

	apiURL := helpers.BuildDataverseApiUrl(environmentHost, buildSourceControlBranchConfigurationCompositeKeyPath(existingBranch.ID, partitionID), nil)
	headers := http.Header{}
	headers.Set("If-Match", "*")
	resp, err := c.Api.Execute(ctx, nil, http.MethodPatch, apiURL, headers, disableSourceControlBranchConfigurationDto{
		StatusCode: sourceControlBranchConfigurationStatusInactive,
	}, []int{http.StatusNoContent, http.StatusNotFound, http.StatusForbidden}, nil)
	if err != nil {
		return err
	}
	if err := c.Api.HandleForbiddenResponse(resp); err != nil {
		return err
	}
	if resp.HttpResponse.StatusCode == http.StatusNotFound {
		return nil
	}

	return c.waitForSolutionGitBranchRemoval(ctx, environmentID, existingBranch.ID, configurationID, partitionID)
}

func (c *client) waitForSolutionGitBranchRemoval(ctx context.Context, environmentID, branchID, configurationID, partitionID string) error {
	for {
		branches, err := c.ListSourceControlBranchConfigurationsByPartition(ctx, environmentID, partitionID)
		if err != nil {
			return err
		}

		found := false
		for _, branch := range branches {
			if strings.EqualFold(branch.ID, branchID) || (strings.EqualFold(branch.SourceControlConfiguration, configurationID) && strings.EqualFold(branch.PartitionID, partitionID)) {
				found = true
				break
			}
		}
		if !found {
			return nil
		}

		if err := c.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
			return err
		}
	}
}

func (c *client) FindSolutionGitBranchByPartition(ctx context.Context, environmentID, configurationID, partitionID string) (*sourceControlBranchConfigurationDto, error) {
	return c.lookupActiveSolutionGitBranchByPartition(ctx, environmentID, configurationID, partitionID)
}

func (c *client) lookupActiveSolutionGitBranchByPartition(ctx context.Context, environmentID, configurationID, partitionID string) (*sourceControlBranchConfigurationDto, error) {
	branches, err := c.ListSourceControlBranchConfigurationsByPartition(ctx, environmentID, partitionID)
	if err != nil {
		return nil, err
	}

	for _, branch := range branches {
		if strings.EqualFold(branch.SourceControlConfiguration, configurationID) && branch.StatusCode == sourceControlBranchConfigurationStatusActive {
			return &branch, nil
		}
	}

	return nil, customerrors.WrapIntoProviderError(nil, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), "source control branch configuration not found")
}

func (c *client) lookupAnySolutionGitBranchByPartition(ctx context.Context, environmentID, configurationID, partitionID string) (*sourceControlBranchConfigurationDto, error) {
	branches, err := c.ListSourceControlBranchConfigurationsByPartition(ctx, environmentID, partitionID)
	if err != nil {
		return nil, err
	}

	for _, branch := range branches {
		if strings.EqualFold(branch.SourceControlConfiguration, configurationID) {
			return &branch, nil
		}
	}

	return nil, customerrors.WrapIntoProviderError(nil, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), "source control branch configuration not found")
}

func (c *client) ListSourceControlBranchConfigurations(ctx context.Context, environmentID, configurationID string) ([]sourceControlBranchConfigurationDto, error) {
	values := url.Values{}
	values.Add("$filter", fmt.Sprintf("(_sourcecontrolconfigurationid_value eq %s)", configurationID))

	return c.querySourceControlBranchConfigurations(ctx, environmentID, values)
}

func (c *client) ListSourceControlBranchConfigurationsByPartition(ctx context.Context, environmentID, partitionID string) ([]sourceControlBranchConfigurationDto, error) {
	values := url.Values{}
	values.Add("partitionId", partitionID)

	return c.querySourceControlBranchConfigurations(ctx, environmentID, values)
}

func (c *client) GetSolutionPartition(ctx context.Context, environmentID, solutionUniqueName string) (*solution.SolutionDto, error) {
	return c.SolutionClient.GetSolutionUniqueName(ctx, environmentID, solutionUniqueName)
}

func (c *client) ListGitOrganizations(ctx context.Context, environmentID string) ([]gitOrganizationDto, error) {
	environmentHost, err := c.EnvironmentClient.GetEnvironmentHostById(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	var organizations gitOrganizationArrayDto
	resp, err := c.Api.Execute(ctx, nil, http.MethodGet, helpers.BuildDataverseApiUrl(environmentHost, "/api/data/v9.0/gitorganizations", nil), nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &organizations)
	if err != nil {
		return nil, err
	}
	if err := c.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if err := c.Api.HandleNotFoundResponse(resp); err != nil {
		return nil, err
	}

	return organizations.Value, nil
}

func (c *client) ListGitProjects(ctx context.Context, environmentID, organizationName string) ([]gitProjectDto, error) {
	environmentHost, err := c.EnvironmentClient.GetEnvironmentHostById(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Add("$filter", fmt.Sprintf("(organizationname eq '%s')", escapeODataString(organizationName)))

	var projects gitProjectArrayDto
	resp, err := c.Api.Execute(ctx, nil, http.MethodGet, helpers.BuildDataverseApiUrl(environmentHost, "/api/data/v9.0/gitprojects", values), nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &projects)
	if err != nil {
		return nil, err
	}
	if err := c.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if err := c.Api.HandleNotFoundResponse(resp); err != nil {
		return nil, err
	}

	return projects.Value, nil
}

func (c *client) ListGitRepositories(ctx context.Context, environmentID, organizationName, projectName string) ([]gitRepositoryDto, error) {
	environmentHost, err := c.EnvironmentClient.GetEnvironmentHostById(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	filters := []string{fmt.Sprintf("(organizationname eq '%s'", escapeODataString(organizationName))}
	if strings.TrimSpace(projectName) != "" {
		filters = append(filters, fmt.Sprintf("projectname eq '%s'", escapeODataString(projectName)))
	}
	values.Add("$filter", strings.Join(filters, " and ")+")")

	var repositories gitRepositoryArrayDto
	resp, err := c.Api.Execute(ctx, nil, http.MethodGet, helpers.BuildDataverseApiUrl(environmentHost, "/api/data/v9.0/gitrepositories", values), nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &repositories)
	if err != nil {
		return nil, err
	}
	if err := c.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if err := c.Api.HandleNotFoundResponse(resp); err != nil {
		return nil, err
	}

	return repositories.Value, nil
}

func (c *client) GetGitRepositoryDefaultBranch(ctx context.Context, environmentID, organizationName, projectName, repositoryName string) (string, error) {
	repositories, err := c.ListGitRepositories(ctx, environmentID, organizationName, projectName)
	if err != nil {
		return "", err
	}

	for _, repository := range repositories {
		if strings.EqualFold(repository.RepositoryName, repositoryName) {
			defaultBranch := strings.TrimSpace(repository.DefaultBranch)
			defaultBranch = strings.TrimPrefix(defaultBranch, "refs/heads/")
			if defaultBranch == "" {
				return "main", nil
			}
			return defaultBranch, nil
		}
	}

	return "", customerrors.WrapIntoProviderError(fmt.Errorf("git repository '%s' not found", repositoryName), customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), "git repository not found")
}

func (c *client) ListGitBranches(ctx context.Context, environmentID, organizationName, projectName, repositoryName string) ([]gitBranchDto, error) {
	environmentHost, err := c.EnvironmentClient.GetEnvironmentHostById(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	filters := []string{fmt.Sprintf("(organizationname eq '%s'", escapeODataString(organizationName))}
	if strings.TrimSpace(projectName) != "" {
		filters = append(filters, fmt.Sprintf("projectname eq '%s'", escapeODataString(projectName)))
	}
	filters = append(filters, fmt.Sprintf("repositoryname eq '%s'", escapeODataString(repositoryName)))
	values.Add("$filter", strings.Join(filters, " and ")+")")

	var branches gitBranchArrayDto
	resp, err := c.Api.Execute(ctx, nil, http.MethodGet, helpers.BuildDataverseApiUrl(environmentHost, "/api/data/v9.0/gitbranches", values), nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &branches)
	if err != nil {
		return nil, err
	}
	if err := c.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if err := c.Api.HandleNotFoundResponse(resp); err != nil {
		return nil, err
	}

	return branches.Value, nil
}

func (c *client) GetUnmanagedSolutionByID(ctx context.Context, environmentID, solutionID string) (*unmanagedSolutionDto, error) {
	environmentHost, err := c.EnvironmentClient.GetEnvironmentHostById(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Add("$filter", fmt.Sprintf("(ismanaged eq false and solutionid eq %s)", solutionID))
	values.Add("$select", "solutionid,uniquename,friendlyname,ismanaged,isvisible,enabledforsourcecontrolintegration,version")

	var solutions unmanagedSolutionArrayDto
	resp, err := c.Api.Execute(ctx, nil, http.MethodGet, helpers.BuildDataverseApiUrl(environmentHost, "/api/data/v9.2/solutions", values), nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &solutions)
	if err != nil {
		return nil, err
	}
	if err := c.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if err := c.Api.HandleNotFoundResponse(resp); err != nil {
		return nil, err
	}
	if len(solutions.Value) == 0 {
		baseErr := fmt.Errorf("unmanaged solution with id '%s' not found", solutionID)
		return nil, customerrors.WrapIntoProviderError(baseErr, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), baseErr.Error())
	}

	return &solutions.Value[0], nil
}

func (c *client) EnableEnvironmentScopeSolutions(ctx context.Context, environmentID string) error {
	solutions, err := c.ListEnvironmentScopeSolutions(ctx, environmentID)
	if err != nil {
		return err
	}

	for _, solutionRow := range solutions {
		if solutionRow.EnabledForSourceControlIntegration {
			continue
		}

		if err := c.EnableSolutionSourceControlIntegration(ctx, environmentID, solutionRow.ID); err != nil {
			return err
		}

		if err := c.waitForEnvironmentScopeSolutionEnabled(ctx, environmentID, solutionRow.ID, solutionRow.UniqueName); err != nil {
			return err
		}
	}

	return nil
}

func (c *client) EnableSolutionSourceControlIntegration(ctx context.Context, environmentID, solutionID string) error {
	environmentHost, err := c.EnvironmentClient.GetEnvironmentHostById(ctx, environmentID)
	if err != nil {
		return err
	}

	apiURL := helpers.BuildDataverseApiUrl(environmentHost, fmt.Sprintf("/api/data/v9.0/solutions(%s)", solutionID), nil)
	resp, err := c.Api.Execute(ctx, nil, http.MethodPatch, apiURL, nil, updateSolutionSourceControlIntegrationDto{
		EnabledForSourceControlIntegration: true,
	}, []int{http.StatusNoContent, http.StatusForbidden, http.StatusNotFound}, nil)
	if err != nil {
		return err
	}
	if err := c.Api.HandleForbiddenResponse(resp); err != nil {
		return err
	}
	return c.Api.HandleNotFoundResponse(resp)
}

func (c *client) waitForEnvironmentScopeSolutionEnabled(ctx context.Context, environmentID, solutionID, solutionUniqueName string) error {
	stableReads := 0

	for {
		solutionRow, err := c.GetUnmanagedSolutionByID(ctx, environmentID, solutionID)
		if err != nil {
			return err
		}

		if solutionRow.EnabledForSourceControlIntegration {
			preValidated, err := c.PreValidateGitComponents(ctx, environmentID, solutionUniqueName)
			if err != nil {
				return err
			}
			if preValidated {
				stableReads++
				if stableReads >= stabilizedReadCountRequired {
					return nil
				}
			} else {
				stableReads = 0
			}
		} else {
			stableReads = 0
		}

		if err := c.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
			return err
		}
	}
}

func (c *client) GetSourceControlIntegrationScope(ctx context.Context, environmentID string) (string, error) {
	orgSettings, err := c.getOrganizationSettings(ctx, environmentID)
	if err != nil {
		return "", err
	}

	return sourceControlIntegrationScopeFromOrgDbValue(extractOrgDbOrgSettingValue(orgSettings.OrgDbOrgSettings, "SourceControlIntegrationScope")), nil
}

func (c *client) SetSourceControlIntegrationScope(ctx context.Context, environmentID, scope string) error {
	orgSettings, err := c.getOrganizationSettings(ctx, environmentID)
	if err != nil {
		return err
	}

	updatedOrgSettings := setOrgDbOrgSettingValue(orgSettings.OrgDbOrgSettings, "SourceControlIntegrationScope", sourceControlIntegrationScopeToOrgDbValue(scope))
	updateDTO := organizationSettingsDto{
		OrgDbOrgSettings: updatedOrgSettings,
	}

	environmentHost, err := c.EnvironmentClient.GetEnvironmentHostById(ctx, environmentID)
	if err != nil {
		return err
	}

	apiURL := helpers.BuildDataverseApiUrl(environmentHost, fmt.Sprintf("/api/data/v9.0/organizations(%s)", orgSettings.OrganizationID), nil)
	resp, err := c.Api.Execute(ctx, nil, http.MethodPatch, apiURL, nil, updateDTO, []int{http.StatusNoContent, http.StatusForbidden, http.StatusNotFound}, nil)
	if err != nil {
		return err
	}
	if err := c.Api.HandleForbiddenResponse(resp); err != nil {
		return err
	}
	return c.Api.HandleNotFoundResponse(resp)
}

func (c *client) getOrganizationSettings(ctx context.Context, environmentID string) (*organizationSettingsDto, error) {
	environmentHost, err := c.EnvironmentClient.GetEnvironmentHostById(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Add("$select", "organizationid,orgdborgsettings")

	var organizations organizationSettingsArrayDto
	resp, err := c.Api.Execute(ctx, nil, http.MethodGet, helpers.BuildDataverseApiUrl(environmentHost, "/api/data/v9.0/organizations", values), nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &organizations)
	if err != nil {
		return nil, err
	}
	if err := c.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if err := c.Api.HandleNotFoundResponse(resp); err != nil {
		return nil, err
	}
	if len(organizations.Value) == 0 {
		return nil, customerrors.WrapIntoProviderError(nil, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), "organization settings not found")
	}

	return &organizations.Value[0], nil
}

func (c *client) querySourceControlBranchConfigurations(ctx context.Context, environmentID string, values url.Values) ([]sourceControlBranchConfigurationDto, error) {
	environmentHost, err := c.EnvironmentClient.GetEnvironmentHostById(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	var branches sourceControlBranchConfigurationArrayDto
	resp, err := c.Api.Execute(ctx, nil, http.MethodGet, helpers.BuildDataverseApiUrl(environmentHost, "/api/data/v9.0/sourcecontrolbranchconfigurations", values), nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &branches)
	if err != nil {
		return nil, err
	}
	if err := c.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if err := c.Api.HandleNotFoundResponse(resp); err != nil {
		return nil, err
	}

	return branches.Value, nil
}

func buildSourceControlConfigurationBindPath(configurationID string) string {
	return fmt.Sprintf("/sourcecontrolconfigurations(%s)", configurationID)
}

func buildSourceControlBranchConfigurationCompositeKeyPath(branchID, partitionID string) string {
	return fmt.Sprintf("/api/data/v9.0/sourcecontrolbranchconfigurations(sourcecontrolbranchconfigurationid=%s,partitionid='%s')", branchID, partitionID)
}

func escapeODataString(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func sourceControlIntegrationScopeToOrgDbValue(scope string) string {
	switch scope {
	case scopeEnvironment:
		return "EnvironmentScope"
	default:
		return "SolutionScope"
	}
}

func sourceControlIntegrationScopeFromOrgDbValue(value string) string {
	switch strings.TrimSpace(value) {
	case "EnvironmentScope":
		return scopeEnvironment
	case "SolutionScope":
		return scopeSolution
	default:
		return ""
	}
}

func extractOrgDbOrgSettingValue(orgSettingsXML, name string) string {
	if orgSettingsXML == "" {
		return ""
	}

	pattern := regexp.MustCompile(fmt.Sprintf(`<%s>([^<]*)</%s>`, regexp.QuoteMeta(name), regexp.QuoteMeta(name)))
	matches := pattern.FindStringSubmatch(orgSettingsXML)
	if len(matches) != 2 {
		return ""
	}

	return matches[1]
}

func setOrgDbOrgSettingValue(orgSettingsXML, name, value string) string {
	if strings.TrimSpace(orgSettingsXML) == "" {
		return fmt.Sprintf("<OrgSettings><%s>%s</%s></OrgSettings>", name, value, name)
	}

	pattern := regexp.MustCompile(fmt.Sprintf(`<%s>[^<]*</%s>`, regexp.QuoteMeta(name), regexp.QuoteMeta(name)))
	replacement := fmt.Sprintf("<%s>%s</%s>", name, value, name)
	if pattern.MatchString(orgSettingsXML) {
		return pattern.ReplaceAllString(orgSettingsXML, replacement)
	}

	return strings.Replace(orgSettingsXML, "</OrgSettings>", replacement+"</OrgSettings>", 1)
}

func isEnvironmentScopeCandidateSolution(solutionRow unmanagedSolutionDto) bool {
	if !solutionRow.IsVisible || solutionRow.IsManaged {
		return false
	}

	switch strings.ToLower(strings.TrimSpace(solutionRow.ID)) {
	case commonDataServicesDefaultSolutionID, activeSolutionID, defaultSolutionID:
		return false
	}

	switch strings.ToLower(strings.TrimSpace(solutionRow.DisplayName)) {
	case commonDataServicesDefaultSolutionName, defaultSolutionName:
		return false
	}

	return true
}
