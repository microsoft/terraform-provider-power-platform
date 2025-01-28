// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment

import (
	"time"
)

var (
	EnvironmentTypes = []string{"Sandbox", "Production", "Trial", "Developer", "Default"}
)

type EnvironmentDto struct {
	Id         string                   `json:"id"`
	Type       string                   `json:"type"`
	Location   string                   `json:"location"`
	Name       string                   `json:"name"`
	Properties *EnviromentPropertiesDto `json:"properties"`
}

type EnviromentPropertiesDto struct {
	AzureRegion               string                            `json:"azureRegion,omitempty"`
	DatabaseType              string                            `json:"databaseType"`
	DisplayName               string                            `json:"displayName"`
	EnvironmentSku            string                            `json:"environmentSku"`
	LinkedAppMetadata         *LinkedAppMetadataDto             `json:"linkedAppMetadata,omitempty"`
	RuntimeEndpoints          *RuntimeEndpointsDto              `json:"runtimeEndpoints,omitempty"`
	LinkedEnvironmentMetadata *LinkedEnvironmentMetadataDto     `json:"linkedEnvironmentMetadata,omitempty"`
	States                    *StatesEnvironmentDto             `json:"states"`
	TenantId                  string                            `json:"tenantId"`
	GovernanceConfiguration   GovernanceConfigurationDto        `json:"governanceConfiguration"`
	BillingPolicy             *BillingPolicyDto                 `json:"billingPolicy,omitempty"`
	ProvisioningState         string                            `json:"provisioningState,omitempty"`
	Description               string                            `json:"description,omitempty"`
	UpdateCadence             *UpdateCadenceDto                 `json:"updateCadence,omitempty"`
	ParentEnvironmentGroup    *ParentEnvironmentGroupDto        `json:"parentEnvironmentGroup,omitempty"`
	EnterprisePolicies        *EnvironmentEnterprisePoliciesDto `json:"enterprisePolicies,omitempty"`
}

type EnvironmentEnterprisePoliciesDto struct {
	Vnets               *EnterprisePolicyDto `json:"vnets,omitempty"`
	CustomerManagedKeys *EnterprisePolicyDto `json:"customerManagedKeys,omitempty"`
}

type EnterprisePolicyDto struct {
	PolicyId   string `json:"policyId"`
	Location   string `json:"location"`
	Id         string `json:"id"`
	SystemId   string `json:"systemId"`
	LinkStatus string `json:"linkStatus"`
}

type ParentEnvironmentGroupDto struct {
	Id string `json:"id"`
}

type UpdateCadenceDto struct {
	Id string `json:"id"`
}

type BillingPolicyDto struct {
	Id string `json:"id"`
}

type GovernanceConfigurationDto struct {
	ProtectionLevel string       `json:"protectionLevel"`
	Settings        *SettingsDto `json:"settings,omitempty"`
}

type SettingsDto struct {
	ExtendedSettings ExtendedSettingsDto `json:"extendedSettings"`
}

type ExtendedSettingsDto struct {
	ExcludeEnvironmentFromAnalysis string `json:"excludeEnvironmentFromAnalysis"`
	IsGroupSharingDisabled         string `json:"isGroupSharingDisabled"`
	MaxLimitUserSharing            string `json:"maxLimitUserSharing"`
	DisableAiGeneratedDescriptions string `json:"disableAiGeneratedDescriptions"`
	IncludeOnHomepageInsights      string `json:"includeOnHomepageInsights"`
	LimitSharingMode               string `json:"limitSharingMode"`
	SolutionCheckerMode            string `json:"solutionCheckerMode"`
	SuppressValidationEmails       string `json:"suppressValidationEmails"`
	SolutionCheckerRuleOverrides   string `json:"solutionCheckerRuleOverrides"`
	MakerOnboardingUrl             string `json:"makerOnboardingUrl"`
	MakerOnboardingMarkdown        string `json:"makerOnboardingMarkdown"`
}

type LinkedEnvironmentMetadataDto struct {
	BackgroundOperationsState string                     `json:"backgroundOperationsState,omitempty"`
	DomainName                string                     `json:"domainName,omitempty"`
	InstanceURL               string                     `json:"instanceUrl"`
	BaseLanguage              int                        `json:"baseLanguage"`
	SecurityGroupId           string                     `json:"securityGroupId"`
	ResourceId                string                     `json:"resourceId"`
	Version                   string                     `json:"version"`
	Templates                 []string                   `json:"template,omitempty"`
	TemplateMetadata          *createTemplateMetadataDto `json:"templateMetadata,omitempty"`
	UniqueName                string                     `json:"uniqueName"`
}

type RuntimeEndpointsDto struct {
	BusinessAppPlatform string `json:"microsoft.BusinessAppPlatform"`
	CommonDataModel     string `json:"microsoft.CommonDataModel"`
	PowerApps           string `json:"microsoft.PowerApps"`
	PowerAppsAdvisor    string `json:"microsoft.PowerAppsAdvisor"`
	PowerVirtualAgents  string `json:"microsoft.PowerVirtualAgents"`
	ApiManagement       string `json:"microsoft.ApiManagement"`
	Flow                string `json:"microsoft.Flow"`
}

type LinkedAppMetadataDto struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

type StatesEnvironmentDto struct {
	Management StatesManagementEnvironmentDto `json:"management"`
	Runtime    *RuntimeEnvironmentDto         `json:"runtime,omitempty"`
}

type RuntimeEnvironmentDto struct {
	Id string `json:"id"`
}

type StatesManagementEnvironmentDto struct {
	Id string `json:"id"`
}

type environmentArrayDto struct {
	Value []EnvironmentDto `json:"value"`
}

type environmentCreateDto struct {
	Location   string                         `json:"location"`
	Properties environmentCreatePropertiesDto `json:"properties"`
}

type modifySkuDto struct {
	EnvironmentSku string `json:"environmentSku"`
}

type environmentCreatePropertiesDto struct {
	AzureRegion               string                            `json:"azureRegion,omitempty"`
	BillingPolicy             BillingPolicyDto                  `json:"billingPolicy,omitempty"`
	DataBaseType              string                            `json:"databaseType,omitempty"`
	DisplayName               string                            `json:"displayName"`
	Description               string                            `json:"description,omitempty"`
	UpdateCadence             *UpdateCadenceDto                 `json:"updateCadence,omitempty"`
	EnvironmentSku            string                            `json:"environmentSku"`
	LinkedEnvironmentMetadata *createLinkEnvironmentMetadataDto `json:"linkedEnvironmentMetadata,omitempty"`
	ParentEnvironmentGroup    *ParentEnvironmentGroupDto        `json:"parentEnvironmentGroup,omitempty"`
}

type createLinkEnvironmentMetadataDto struct {
	BaseLanguage     int                        `json:"baseLanguage"`
	DomainName       string                     `json:"domainName,omitempty"`
	Currency         createCurrencyDto          `json:"currency"`
	SecurityGroupId  string                     `json:"securityGroupId,omitempty"`
	Templates        []string                   `json:"templates,omitempty"`
	TemplateMetadata *createTemplateMetadataDto `json:"templateMetadata,omitempty"`
}
type createCurrencyDto struct {
	Code string `json:"code"`
}

type createTemplateMetadataDto struct {
	PostProvisioningPackages []createPostProvisioningPackagesDto `json:"PostProvisioningPackages,omitempty"`
}

type createPostProvisioningPackagesDto struct {
	ApplicationUniqueName string `json:"applicationUniqueName,omitempty"`
	Parameters            string `json:"parameters,omitempty"`
}

type enironmentDeleteDto struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type lifecycleCreatedDto struct {
	Name       string                        `json:"name"`
	Properties lifecycleCreatedPropertiesDto `json:"properties"`
}

type lifecycleCreatedPropertiesDto struct {
	ProvisioningState string `json:"provisioningState"`
}

type organizationSettingsArrayDto struct {
	Value []organizationSettingsDto `json:"value"`
}

type organizationSettingsDto struct {
	ODataEtag      string    `json:"@odata.etag"`
	CreatedOn      time.Time `json:"createdon"`
	BaseCurrencyId string    `json:"_basecurrencyid_value"`
}

type TransactionCurrencyDto struct {
	OrganizationValue     string  `json:"_organizationid_value"`
	CurrencyName          string  `json:"currencyname"`
	CurrencySymbol        string  `json:"currencysymbol"`
	IsoCurrencyCode       string  `json:"isocurrencycode"`
	CreatedOn             string  `json:"createdon"`
	CurrencyPrecision     int     `json:"currencyprecision"`
	ExchangeRate          float32 `json:"exchangerate"`
	TransactionCurrencyId string  `json:"transactioncurrencyid"`
}

type transactionCurrencyArrayDto struct {
	Value []TransactionCurrencyDto `json:"value"`
}

type validateCreateEnvironmentDetailsDto struct {
	DomainName          string `json:"domainName"`
	EnvironmentLocation string `json:"environmentLocation"`
}

type validateUpdateEnvironmentDetailsDto struct {
	DomainName      string `json:"domainName"`
	EnvironmentName string `json:"environmentName"`
}

type LocationArrayDto struct {
	Value []LocationDto `json:"value"`
}

type LocationDto struct {
	ID         string                `json:"id"`
	Type       string                `json:"type"`
	Name       string                `json:"name"`
	Properties LocationPropertiesDto `json:"properties"`
}

type LocationPropertiesDto struct {
	DisplayName                            string   `json:"displayName"`
	Code                                   string   `json:"code"`
	IsDisabled                             bool     `json:"isDisabled"`
	IsDefault                              bool     `json:"isDefault"`
	CanProvisionDatabase                   bool     `json:"canProvisionDatabase"`
	CanProvisionCustomerEngagementDatabase bool     `json:"canProvisionCustomerEngagementDatabase"`
	AzureRegions                           []string `json:"azureRegions"`
}
