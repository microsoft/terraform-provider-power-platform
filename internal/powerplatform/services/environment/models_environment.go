// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
	AzureRegion               string                        `json:"azureRegion,omitempty"`
	DatabaseType              string                        `json:"databaseType"`
	DisplayName               string                        `json:"displayName"`
	EnvironmentSku            string                        `json:"environmentSku"`
	LinkedAppMetadata         *LinkedAppMetadataDto         `json:"linkedAppMetadata,omitempty"`
	LinkedEnvironmentMetadata *LinkedEnvironmentMetadataDto `json:"linkedEnvironmentMetadata,omitempty"`
	States                    StatesEnvironmentDto          `json:"states"`
	TenantID                  string                        `json:"tenantId"`
	GovernanceConfiguration   GovernanceConfigurationDto    `json:"governanceConfiguration"`
	BillingPolicy             *BillingPolicyDto             `json:"billingPolicy,omitempty"`
	ProvisioningState         string                        `json:"provisioningState,omitempty"`
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
	DomainName       string                             `json:"domainName,omitempty"`
	InstanceURL      string                             `json:"instanceUrl"`
	BaseLanguage     int                                `json:"baseLanguage"`
	SecurityGroupId  string                             `json:"securityGroupId"`
	ResourceId       string                             `json:"resourceId"`
	Version          string                             `json:"version"`
	Templates        []string                           `json:"template,omitempty"`
	TemplateMetadata *EnvironmentCreateTemplateMetadata `json:"templateMetadata,omitempty"`
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
	AzureRegion               string                                       `json:"azureRegion,omitempty"`
	BillingPolicy             BillingPolicyDto                             `json:"billingPolicy,omitempty"`
	DataBaseType              string                                       `json:"databaseType,omitempty"`
	DisplayName               string                                       `json:"displayName"`
	EnvironmentSku            string                                       `json:"environmentSku"`
	LinkedEnvironmentMetadata *EnvironmentCreateLinkEnvironmentMetadataDto `json:"linkedEnvironmentMetadata,omitempty"`
}

type EnvironmentCreateLinkEnvironmentMetadataDto struct {
	BaseLanguage     int                                `json:"baseLanguage"`
	DomainName       string                             `json:"domainName,omitempty"`
	Currency         EnvironmentCreateCurrency          `json:"currency"`
	SecurityGroupId  string                             `json:"securityGroupId,omitempty"`
	Templates        []string                           `json:"templates,omitempty"`
	TemplateMetadata *EnvironmentCreateTemplateMetadata `json:"templateMetadata,omitempty"`
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

type EnvironmentsListDataSourceModel struct {
	Environments []EnvironmentSourceModel `tfsdk:"environments"`
	Id           types.Int64              `tfsdk:"id"`
}

type EnvironmentSourceModel struct {
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
	Id              types.String   `tfsdk:"id"`
	Location        types.String   `tfsdk:"location"`
	AzureRegion     types.String   `tfsdk:"azure_region"`
	DisplayName     types.String   `tfsdk:"display_name"`
	EnvironmentType types.String   `tfsdk:"environment_type"`
	BillingPolicyId types.String   `tfsdk:"billing_policy_id"`

	Dataverse types.Object `tfsdk:"dataverse"`
}

type DataverseSourceModel struct {
	Url              types.String `tfsdk:"url"`
	Domain           types.String `tfsdk:"domain"`
	OrganizationId   types.String `tfsdk:"organization_id"`
	SecurityGroupId  types.String `tfsdk:"security_group_id"`
	LanguageName     types.Int64  `tfsdk:"language_code"`
	Version          types.String `tfsdk:"version"`
	LinkedAppType    types.String `tfsdk:"linked_app_type"`
	LinkedAppId      types.String `tfsdk:"linked_app_id"`
	LinkedAppURL     types.String `tfsdk:"linked_app_url"`
	CurrencyCode     types.String `tfsdk:"currency_code"`
	Templates        []string     `tfsdk:"templates"`
	TemplateMetadata types.String `tfsdk:"template_metadata"`
}

func ConvertUpdateEnvironmentDtoFromSourceModel(ctx context.Context, environmentSource EnvironmentSourceModel) (*EnvironmentDto, error) {
	environmentDto := EnvironmentDto{
		Id:       environmentSource.Id.ValueString(),
		Name:     environmentSource.DisplayName.ValueString(),
		Type:     environmentSource.EnvironmentType.ValueString(),
		Location: environmentSource.Location.ValueString(),
		Properties: EnvironmentPropertiesDto{
			DisplayName:    environmentSource.DisplayName.ValueString(),
			EnvironmentSku: environmentSource.EnvironmentType.ValueString(),
		},
	}

	if !environmentSource.BillingPolicyId.IsNull() && environmentSource.BillingPolicyId.ValueString() != "" {
		environmentDto.Properties.BillingPolicy = &BillingPolicyDto{
			Id: environmentSource.BillingPolicyId.ValueString(),
		}
	}

	if !environmentSource.Dataverse.IsNull() && !environmentSource.Dataverse.IsUnknown() {

		var dataverseSourceModel DataverseSourceModel
		environmentSource.Dataverse.As(ctx, &dataverseSourceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		environmentDto.Properties.LinkedEnvironmentMetadata = &LinkedEnvironmentMetadataDto{
			SecurityGroupId: dataverseSourceModel.SecurityGroupId.ValueString(),
			DomainName:      dataverseSourceModel.Domain.ValueString(),
		}

	}

	return &environmentDto, nil
}

func IsDataverseEnvironmentEmpty(ctx context.Context, environment *EnvironmentSourceModel) bool {
	var dataverseSourceModel DataverseSourceModel
	environment.Dataverse.As(ctx, &dataverseSourceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

	return dataverseSourceModel.CurrencyCode.IsNull() || dataverseSourceModel.CurrencyCode.ValueString() == ""
}

func ConvertCreateEnvironmentDtoFromSourceModel(ctx context.Context, environmentSource EnvironmentSourceModel) (*EnvironmentCreateDto, error) {
	environmentDto := &EnvironmentCreateDto{
		Location: environmentSource.Location.ValueString(),
		Properties: EnvironmentCreatePropertiesDto{
			DisplayName:    environmentSource.DisplayName.ValueString(),
			EnvironmentSku: environmentSource.EnvironmentType.ValueString(),
		},
	}

	if !environmentSource.AzureRegion.IsNull() && environmentSource.AzureRegion.ValueString() != "" {
		environmentDto.Properties.AzureRegion = environmentSource.AzureRegion.ValueString()
	}

	if !environmentSource.BillingPolicyId.IsNull() && environmentSource.BillingPolicyId.ValueString() != "" {
		environmentDto.Properties.BillingPolicy = BillingPolicyDto{
			Id: environmentSource.BillingPolicyId.ValueString(),
		}
	}

	if !environmentSource.Dataverse.IsNull() && !environmentSource.Dataverse.IsUnknown() {
		var dataverseSourceModel DataverseSourceModel
		environmentSource.Dataverse.As(ctx, &dataverseSourceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		environmentDto.Properties.DataBaseType = "CommonDataService"
		linkedMetadata, err := ConvertEnvironmentCreateLinkEnvironmentMetadataDtoFromDataverseSourceModel(ctx, environmentSource.Dataverse)
		if err != nil {
			return nil, err
		}
		environmentDto.Properties.LinkedEnvironmentMetadata = linkedMetadata

	}
	return environmentDto, nil
}

func ConvertEnvironmentCreateLinkEnvironmentMetadataDtoFromDataverseSourceModel(ctx context.Context, dataverse types.Object) (*EnvironmentCreateLinkEnvironmentMetadataDto, error) {
	if !dataverse.IsNull() && !dataverse.IsUnknown() {
		var dataverseSourceModel DataverseSourceModel
		dataverse.As(ctx, &dataverseSourceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		var templateMetadataObject *EnvironmentCreateTemplateMetadata
		if dataverseSourceModel.TemplateMetadata.ValueString() != "" {
			err := json.Unmarshal([]byte(dataverseSourceModel.TemplateMetadata.ValueString()), &templateMetadataObject)
			if err != nil {
				return nil, fmt.Errorf("error when unmarshalling template metadata %s; internal error: %v", dataverseSourceModel.TemplateMetadata.ValueString(), err)
			}
			if len(templateMetadataObject.PostProvisioningPackages) == 0 {
				templateMetadataObject = nil
			}
		}

		linkedEnvironmentMetadata := &EnvironmentCreateLinkEnvironmentMetadataDto{
			BaseLanguage:    int(dataverseSourceModel.LanguageName.ValueInt64()),
			SecurityGroupId: dataverseSourceModel.SecurityGroupId.ValueString(),
			Currency: EnvironmentCreateCurrency{
				Code: dataverseSourceModel.CurrencyCode.ValueString(),
			},
			Templates:        dataverseSourceModel.Templates,
			TemplateMetadata: templateMetadataObject,
		}

		if !dataverseSourceModel.Domain.IsNull() && dataverseSourceModel.Domain.ValueString() != "" {
			linkedEnvironmentMetadata.DomainName = dataverseSourceModel.Domain.ValueString()
		} else {
			linkedEnvironmentMetadata.DomainName = ""
		}
		return linkedEnvironmentMetadata, nil
	}
	return nil, fmt.Errorf("dataverse object is null or unknown")
}

func ConvertSourceModelFromEnvironmentDto(environmentDto EnvironmentDto, currencyCode *string, templateMetadata *EnvironmentCreateTemplateMetadata, templates []string, timeouts timeouts.Value) (*EnvironmentSourceModel, error) {
	model := &EnvironmentSourceModel{
		Timeouts:        timeouts,
		Id:              types.StringValue(environmentDto.Name),
		DisplayName:     types.StringValue(environmentDto.Properties.DisplayName),
		Location:        types.StringValue(environmentDto.Location),
		AzureRegion:     types.StringValue(environmentDto.Properties.AzureRegion),
		EnvironmentType: types.StringValue(environmentDto.Properties.EnvironmentSku),
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
		"templates":         types.ListType{ElemType: types.StringType},
		"template_metadata": types.StringType,
	}

	attrValuesProductProperties := map[string]attr.Value{}
	model.Dataverse = types.ObjectNull(attrTypesDataverseObject)

	if environmentDto.Properties.LinkedAppMetadata != nil {
		attrValuesProductProperties["linked_app_type"] = types.StringValue(environmentDto.Properties.LinkedAppMetadata.Type)
		attrValuesProductProperties["linked_app_id"] = types.StringValue(environmentDto.Properties.LinkedAppMetadata.Id)
		attrValuesProductProperties["linked_app_url"] = types.StringValue(environmentDto.Properties.LinkedAppMetadata.Url)
	} else {
		attrValuesProductProperties["linked_app_type"] = types.StringNull()
		attrValuesProductProperties["linked_app_id"] = types.StringNull()
		attrValuesProductProperties["linked_app_url"] = types.StringNull()
	}

	if environmentDto.Properties.LinkedEnvironmentMetadata != nil {
		attrValuesProductProperties["url"] = types.StringValue(environmentDto.Properties.LinkedEnvironmentMetadata.InstanceURL)
		attrValuesProductProperties["domain"] = types.StringValue(environmentDto.Properties.LinkedEnvironmentMetadata.DomainName)
		attrValuesProductProperties["organization_id"] = types.StringValue(environmentDto.Properties.LinkedEnvironmentMetadata.ResourceId)
		attrValuesProductProperties["security_group_id"] = types.StringValue(environmentDto.Properties.LinkedEnvironmentMetadata.SecurityGroupId)
		attrValuesProductProperties["language_code"] = types.Int64Value(int64(environmentDto.Properties.LinkedEnvironmentMetadata.BaseLanguage))
		attrValuesProductProperties["version"] = types.StringValue(environmentDto.Properties.LinkedEnvironmentMetadata.Version)
		if currencyCode != nil && *currencyCode != "" {
			attrValuesProductProperties["currency_code"] = types.StringValue(*currencyCode)
		} else {
			attrValuesProductProperties["currency_code"] = types.StringNull()
		}
		if environmentDto.Properties.LinkedEnvironmentMetadata.Templates != nil {
			var templ []attr.Value
			for _, t := range environmentDto.Properties.LinkedEnvironmentMetadata.Templates {
				templ = append(templ, types.StringValue(t))
			}
			v, _ := types.ListValue(types.StringType, templ)

			attrValuesProductProperties["templates"] = v
		} else if templates != nil {
			var templ []attr.Value
			for _, t := range templates {
				templ = append(templ, types.StringValue(t))
			}
			v, _ := types.ListValue(types.StringType, templ)
			attrValuesProductProperties["templates"] = v
		} else {
			attrValuesProductProperties["templates"] = types.ListNull(types.StringType)
		}

		if environmentDto.Properties.LinkedEnvironmentMetadata.TemplateMetadata != nil && environmentDto.Properties.LinkedEnvironmentMetadata.TemplateMetadata.PostProvisioningPackages != nil {
			b, err := json.Marshal(environmentDto.Properties.LinkedEnvironmentMetadata.TemplateMetadata)
			if err != nil {
				return nil, err
			}
			attrValuesProductProperties["template_metadata"] = types.StringValue(string(b))
		} else if templateMetadata != nil {
			b, err := json.Marshal(templateMetadata)
			if err != nil {
				return nil, err
			}
			attrValuesProductProperties["template_metadata"] = types.StringValue(string(b))
		} else {
			attrValuesProductProperties["template_metadata"] = types.StringNull()
		}
	} else {
		attrValuesProductProperties["url"] = types.StringNull()
		attrValuesProductProperties["domain"] = types.StringNull()
		attrValuesProductProperties["organization_id"] = types.StringNull()
		attrValuesProductProperties["security_group_id"] = types.StringNull()
		attrValuesProductProperties["language_code"] = types.Int64Null()
		attrValuesProductProperties["version"] = types.StringNull()
		attrValuesProductProperties["currency_code"] = types.StringNull()
		attrValuesProductProperties["template_metadata"] = types.StringNull()
		attrValuesProductProperties["templates"] = types.ListNull(types.StringType)
	}
	model.Dataverse = types.ObjectValueMust(attrTypesDataverseObject, attrValuesProductProperties)
	return model, nil
}

type LocationArrayDto struct {
	Value []LocationDto `json:"value"`
}

type LocationDto struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	Properties struct {
		DisplayName                            string   `json:"displayName"`
		Code                                   string   `json:"code"`
		IsDefault                              bool     `json:"isDefault"`
		IsDisabled                             bool     `json:"isDisabled"`
		CanProvisionDatabase                   bool     `json:"canProvisionDatabase"`
		CanProvisionCustomerEngagementDatabase bool     `json:"canProvisionCustomerEngagementDatabase"`
		AzureRegions                           []string `json:"azureRegions"`
	} `json:"properties"`
}

type WhoAmIDto struct {
	BusinessUnitId string `json:"BusinessUnitId"`
	UserId         string `json:"UserId"`
	OrganizationId string `json:"OrganizationId"`
}
