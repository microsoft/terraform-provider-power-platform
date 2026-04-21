// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package git_integration

type gitOrganizationDto struct {
	OrganizationName string `json:"organizationname,omitempty"`
}

type gitOrganizationArrayDto struct {
	Value []gitOrganizationDto `json:"value"`
}

type gitProjectDto struct {
	OrganizationName string `json:"organizationname,omitempty"`
	ProjectName      string `json:"projectname,omitempty"`
}

type gitProjectArrayDto struct {
	Value []gitProjectDto `json:"value"`
}

type gitRepositoryDto struct {
	OrganizationName string `json:"organizationname,omitempty"`
	ProjectName      string `json:"projectname,omitempty"`
	RepositoryName   string `json:"repositoryname,omitempty"`
	DefaultBranch    string `json:"defaultbranch,omitempty"`
}

type gitRepositoryArrayDto struct {
	Value []gitRepositoryDto `json:"value"`
}

type gitBranchDto struct {
	OrganizationName   string `json:"organizationname,omitempty"`
	ProjectName        string `json:"projectname,omitempty"`
	RepositoryName     string `json:"repositoryname,omitempty"`
	BranchName         string `json:"branchname,omitempty"`
	UpstreamBranchName string `json:"upstreambranchname,omitempty"`
}

type gitBranchArrayDto struct {
	Value []gitBranchDto `json:"value"`
}

type unmanagedSolutionDto struct {
	ID                                 string `json:"solutionid,omitempty"`
	UniqueName                         string `json:"uniquename,omitempty"`
	DisplayName                        string `json:"friendlyname,omitempty"`
	IsManaged                          bool   `json:"ismanaged,omitempty"`
	IsVisible                          bool   `json:"isvisible,omitempty"`
	EnabledForSourceControlIntegration bool   `json:"enabledforsourcecontrolintegration,omitempty"`
	Version                            string `json:"version,omitempty"`
}

type unmanagedSolutionArrayDto struct {
	Value []unmanagedSolutionDto `json:"value"`
}

type organizationSettingsDto struct {
	OrganizationID   string `json:"organizationid,omitempty"`
	OrgDbOrgSettings string `json:"orgdborgsettings,omitempty"`
}

type organizationSettingsArrayDto struct {
	Value []organizationSettingsDto `json:"value"`
}

type sourceControlConfigurationDto struct {
	ID               string `json:"sourcecontrolconfigurationid,omitempty"`
	Name             string `json:"name,omitempty"`
	OrganizationName string `json:"organizationname,omitempty"`
	ProjectName      string `json:"projectname,omitempty"`
	RepositoryName   string `json:"repositoryname,omitempty"`
	GitProvider      int    `json:"gitprovider,omitempty"`
}

type sourceControlConfigurationArrayDto struct {
	Value []sourceControlConfigurationDto `json:"value"`
}

type sourceControlBranchConfigurationDto struct {
	ID                         string `json:"sourcecontrolbranchconfigurationid,omitempty"`
	Name                       string `json:"name,omitempty"`
	PartitionID                string `json:"partitionid,omitempty"`
	BranchName                 string `json:"branchname,omitempty"`
	UpstreamBranchName         string `json:"upstreambranchname,omitempty"`
	RootFolderPath             string `json:"rootfolderpath,omitempty"`
	BranchSyncedCommitID       string `json:"branchsyncedcommitid,omitempty"`
	UpstreamBranchSyncedCommit string `json:"upstreambranchsyncedcommitid,omitempty"`
	StatusCode                 int    `json:"statuscode,omitempty"`
	SourceControlConfiguration string `json:"_sourcecontrolconfigurationid_value,omitempty"`
}

type sourceControlBranchConfigurationArrayDto struct {
	Value []sourceControlBranchConfigurationDto `json:"value"`
}

type createSourceControlConfigurationDto struct {
	ID               string `json:"sourcecontrolconfigurationid,omitempty"`
	Name             string `json:"name,omitempty"`
	OrganizationName string `json:"organizationname,omitempty"`
	ProjectName      string `json:"projectname,omitempty"`
	RepositoryName   string `json:"repositoryname,omitempty"`
	GitProvider      int    `json:"gitprovider,omitempty"`
}

type updateSourceControlConfigurationDto struct {
	Name             string `json:"name,omitempty"`
	OrganizationName string `json:"organizationname,omitempty"`
	ProjectName      string `json:"projectname,omitempty"`
	RepositoryName   string `json:"repositoryname,omitempty"`
	GitProvider      int    `json:"gitprovider,omitempty"`
}

type createSourceControlBranchConfigurationDto struct {
	ID                               string `json:"sourcecontrolbranchconfigurationid,omitempty"`
	Name                             string `json:"name,omitempty"`
	PartitionID                      string `json:"partitionid,omitempty"`
	BranchName                       string `json:"branchname,omitempty"`
	UpstreamBranchName               string `json:"upstreambranchname,omitempty"`
	RootFolderPath                   string `json:"rootfolderpath,omitempty"`
	SourceControlConfigurationBindID string `json:"sourcecontrolconfigurationid@odata.bind,omitempty"`
}

type updateSourceControlBranchConfigurationDto struct {
	Name               string `json:"name,omitempty"`
	BranchName         string `json:"branchname,omitempty"`
	UpstreamBranchName string `json:"upstreambranchname,omitempty"`
	RootFolderPath     string `json:"rootfolderpath,omitempty"`
}

type disableSourceControlBranchConfigurationDto struct {
	StatusCode int `json:"statuscode,omitempty"`
}

type preValidateGitComponentsRequestDto struct {
	SolutionUniqueName string `json:"SolutionUniqueName"`
}

type preValidateGitComponentsResponseDto struct {
	ValidationMessages string `json:"ValidationMessages,omitempty"`
}

type updateSolutionSourceControlIntegrationDto struct {
	EnabledForSourceControlIntegration bool `json:"enabledforsourcecontrolintegration,omitempty"`
}
