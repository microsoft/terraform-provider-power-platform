// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	EnvironmentTypes = []string{"Sandbox", "Production", "Trial", "Developer"}
)

type EnvironmentDto struct {
	Id         string                   `json:"id"`
	Type       string                   `json:"type"`
	Location   string                   `json:"location"`
	Name       string                   `json:"name"`
	Properties EnvironmentPropertiesDto `json:"properties"`
}

type EnvironmentPropertiesDto struct {
	DatabaseType              string                        `json:"databaseType"`
	DisplayName               string                        `json:"displayName"`
	EnvironmentSku            string                        `json:"environmentSku"`
	LinkedAppMetadata         *LinkedAppMetadataDto         `json:"linkedAppMetadata,omitempty"`
	LinkedEnvironmentMetadata *LinkedEnvironmentMetadataDto `json:"linkedEnvironmentMetadata,omitempty"`
	States                    StatesEnvironmentDto          `json:"states"`
	TenantID                  string                        `json:"tenantId"`
	GovernanceConfiguration   GovernanceConfigurationDto    `json:"governanceConfiguration"`
	BillingPolicy             *BillingPolicyDto             `json:"billingPolicy,omitempty"`
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
	//MakerOnboardingTimestamp       time.Time `json:"makerOnboardingTimestamp"`
	MakerOnboardingMarkdown string `json:"makerOnboardingMarkdown"`
}

type LinkedEnvironmentMetadataDto struct {
	DomainName       string                            `json:"domainName,omitempty"`
	InstanceURL      string                            `json:"instanceUrl"`
	BaseLanguage     int                               `json:"baseLanguage"`
	SecurityGroupId  string                            `json:"securityGroupId"`
	ResourceId       string                            `json:"resourceId"`
	Version          string                            `json:"version"`
	Templates        []string                          `json:"template,omitempty"`
	TemplateMetadata EnvironmentCreateTemplateMetadata `json:"templateMetadata,omitempty"`
}

type LinkedAppMetadataDto struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

type StatesEnvironmentDto struct {
	Management StatesManagementEnvironmentDto `json:"management"`
}

type StatesManagementEnvironmentDto struct {
	Id string `json:"id"`
}

type EnvironmentDtoArray struct {
	Value []EnvironmentDto `json:"value"`
}

type EnvironmentCreateDto struct {
	Location   string                         `json:"location"`
	Properties EnvironmentCreatePropertiesDto `json:"properties"`
}

type EnvironmentCreatePropertiesDto struct {
	BillingPolicy             BillingPolicyDto                             `json:"billingPolicy,omitempty"`
	DataBaseType              string                                       `json:"databaseType,omitempty"`
	DisplayName               string                                       `json:"displayName"`
	EnvironmentSku            string                                       `json:"environmentSku"`
	LinkedEnvironmentMetadata *EnvironmentCreateLinkEnvironmentMetadataDto `json:"linkedEnvironmentMetadata,omitempty"`
}

type EnvironmentCreateLinkEnvironmentMetadataDto struct {
	BaseLanguage     int                               `json:"baseLanguage"`
	DomainName       string                            `json:"domainName,omitempty"`
	Currency         EnvironmentCreateCurrency         `json:"currency"`
	SecurityGroupId  string                            `json:"securityGroupId,omitempty"`
	Templates        []string                          `json:"templates,omitempty"`
	TemplateMetadata EnvironmentCreateTemplateMetadata `json:"templateMetadata,omitempty"`
}
type EnvironmentCreateCurrency struct {
	Code string `json:"code"`
}

type EnvironmentCreateTemplateMetadata struct {
	PostProvisioningPackages []EnvironmentCreatePostProvisioningPackages `json:"PostProvisioningPackages,omitempty"`
}

type EnvironmentCreatePostProvisioningPackages struct {
	ApplicationUniqueName string `json:"applicationUniqueName,omitempty"`
	Parameters            string `json:"parameters,omitempty"`
}

type EnvironmentCreateLinkedAppMetadataDto struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

type EnvironmentDeleteDto struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type EnvironmentLifecycleCreatedDto struct {
	Name       string                                   `json:"name"`
	Properties EnvironmentLifecycleCreatedPropertiesDto `json:"properties"`
}

type EnvironmentLifecycleCreatedPropertiesDto struct {
	ProvisioningState string `json:"provisioningState"`
}

type OrganizationSettingsArrayDto struct {
	Value []OrganizationSettingsDto `json:"value"`
}

type OrganizationSettingsDto struct {
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

type TransactionCurrencyArrayDto struct {
	Value []TransactionCurrencyDto `json:"value"`
}

type ValidateEnvironmentDetailsDto struct {
	DomainName          string `json:"domainName"`
	EnvironmentLocation string `json:"environmentLocation"`
}

func ConvertFromEnvironmentDto(environmentDto EnvironmentDto, currencyCode string) EnvironmentDataSourceModel {
	model := EnvironmentDataSourceModel{
		EnvironmentId:   types.StringValue(environmentDto.Name),
		DisplayName:     types.StringValue(environmentDto.Properties.DisplayName),
		Location:        types.StringValue(environmentDto.Location),
		EnvironmentType: types.StringValue(environmentDto.Properties.EnvironmentSku),
		//BillingPolicyId: types.StringValue(environmentDto.Properties.BillingPolicy.Id),
	}

	if environmentDto.Properties.BillingPolicy != nil {
		model.BillingPolicyId = types.StringValue(environmentDto.Properties.BillingPolicy.Id)
	} else {
		model.BillingPolicyId = types.StringValue("")
	}

	attrTypesDataverseObject := map[string]attr.Type{
		"url":               types.StringType,
		"domain":            types.StringType,
		"organization_id":   types.StringType,
		"security_group_id": types.StringType,
		"language_code":     types.Int64Type,
		"version":           types.StringType,
		"linked_app_type":   types.StringType,
		"linked_app_id":     types.StringType,
		"linked_app_url":    types.StringType,
		"currency_code":     types.StringType,
	}

	attrValuesProductProperties := map[string]attr.Value{}

	if environmentDto.Properties.LinkedAppMetadata != nil {
		attrValuesProductProperties["linked_app_type"] = types.StringValue(environmentDto.Properties.LinkedAppMetadata.Type)
		attrValuesProductProperties["linked_app_id"] = types.StringValue(environmentDto.Properties.LinkedAppMetadata.Id)
		attrValuesProductProperties["linked_app_url"] = types.StringValue(environmentDto.Properties.LinkedAppMetadata.Url)
	} else {
		attrValuesProductProperties["linked_app_type"] = types.StringValue("")
		attrValuesProductProperties["linked_app_id"] = types.StringValue("")
		attrValuesProductProperties["linked_app_url"] = types.StringValue("")
	}

	if environmentDto.Properties.LinkedEnvironmentMetadata != nil {
		attrValuesProductProperties["url"] = types.StringValue(environmentDto.Properties.LinkedEnvironmentMetadata.InstanceURL)
		attrValuesProductProperties["domain"] = types.StringValue(environmentDto.Properties.LinkedEnvironmentMetadata.DomainName)
		attrValuesProductProperties["organization_id"] = types.StringValue(environmentDto.Properties.LinkedEnvironmentMetadata.ResourceId)
		attrValuesProductProperties["security_group_id"] = types.StringValue(environmentDto.Properties.LinkedEnvironmentMetadata.SecurityGroupId)
		attrValuesProductProperties["language_code"] = types.Int64Value(int64(environmentDto.Properties.LinkedEnvironmentMetadata.BaseLanguage))
		attrValuesProductProperties["version"] = types.StringValue(environmentDto.Properties.LinkedEnvironmentMetadata.Version)
		attrValuesProductProperties["currency_code"] = types.StringValue(currencyCode)
	} else {
		attrValuesProductProperties["url"] = types.StringValue("")
		attrValuesProductProperties["domain"] = types.StringValue("")
		attrValuesProductProperties["organization_id"] = types.StringValue("")
		attrValuesProductProperties["security_group_id"] = types.StringValue("")
		attrValuesProductProperties["language_code"] = types.Int64Null()
		attrValuesProductProperties["version"] = types.StringValue("")
		attrValuesProductProperties["currency_code"] = types.StringValue("")
	}
	model.Dataverse = types.ObjectValueMust(attrTypesDataverseObject, attrValuesProductProperties)

	return model
}
