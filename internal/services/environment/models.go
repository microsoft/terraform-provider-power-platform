// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/licensing"
)

type EnvironmentsDataSource struct {
	helpers.TypeInfo
	EnvironmentClient Client
}

type Resource struct {
	helpers.TypeInfo
	EnvironmentClient Client
	LicensingClient   licensing.Client
}

type ListDataSourceModel struct {
	Timeouts     timeouts.Value `tfsdk:"timeouts"`
	Environments []SourceModel  `tfsdk:"environments"`
}

type SourceModel struct {
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
	Id                 types.String   `tfsdk:"id"`
	Location           types.String   `tfsdk:"location"`
	AzureRegion        types.String   `tfsdk:"azure_region"`
	DisplayName        types.String   `tfsdk:"display_name"`
	EnvironmentType    types.String   `tfsdk:"environment_type"`
	BillingPolicyId    types.String   `tfsdk:"billing_policy_id"`
	Description        types.String   `tfsdk:"description"`
	Cadence            types.String   `tfsdk:"cadence"`
	EnvironmentGroupId types.String   `tfsdk:"environment_group_id"`
	OwnerId            types.String   `tfsdk:"owner_id"`

	EnterprisePolicies basetypes.SetValue `tfsdk:"enterprise_policies"`

	Dataverse types.Object `tfsdk:"dataverse"`
}

type EnterprisePoliciesModel struct {
	Type     types.String `tfsdk:"type"`
	Id       types.String `tfsdk:"id"`
	Location types.String `tfsdk:"location"`
	SystemId types.String `tfsdk:"system_id"`
	Status   types.String `tfsdk:"status"`
}

type DataverseSourceModel struct {
	Url                 types.String `tfsdk:"url"`
	Domain              types.String `tfsdk:"domain"`
	OrganizationId      types.String `tfsdk:"organization_id"`
	SecurityGroupId     types.String `tfsdk:"security_group_id"`
	LanguageName        types.Int64  `tfsdk:"language_code"`
	Version             types.String `tfsdk:"version"`
	LinkedAppType       types.String `tfsdk:"linked_app_type"`
	LinkedAppId         types.String `tfsdk:"linked_app_id"`
	LinkedAppURL        types.String `tfsdk:"linked_app_url"`
	CurrencyCode        types.String `tfsdk:"currency_code"`
	Templates           []string     `tfsdk:"templates"`
	TemplateMetadata    types.String `tfsdk:"template_metadata"`
	AdministrationMode  types.Bool   `tfsdk:"administration_mode_enabled"`
	BackgroundOperation types.Bool   `tfsdk:"background_operation_enabled"`
	UniqueName          types.String `tfsdk:"unique_name"`
}

func isDataverseEnvironmentEmpty(ctx context.Context, environment *SourceModel) bool {
	var dataverseSourceModel DataverseSourceModel
	environment.Dataverse.As(ctx, &dataverseSourceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

	return dataverseSourceModel.CurrencyCode.IsNull() || dataverseSourceModel.CurrencyCode.ValueString() == ""
}

func convertCreateEnvironmentDtoFromSourceModel(ctx context.Context, environmentSource SourceModel) (*environmentCreateDto, error) {
	environmentDto := &environmentCreateDto{
		Location: environmentSource.Location.ValueString(),
		Properties: environmentCreatePropertiesDto{
			DisplayName:    environmentSource.DisplayName.ValueString(),
			EnvironmentSku: environmentSource.EnvironmentType.ValueString(),
		},
	}

	if !environmentSource.Description.IsNull() && environmentSource.Description.ValueString() != "" {
		environmentDto.Properties.Description = environmentSource.Description.ValueString()
	}

	if !environmentSource.Cadence.IsNull() && environmentSource.Cadence.ValueString() != "" {
		environmentDto.Properties.UpdateCadence = &UpdateCadenceDto{
			Id: environmentSource.Cadence.ValueString(),
		}
	}

	if !environmentSource.AzureRegion.IsNull() && environmentSource.AzureRegion.ValueString() != "" {
		environmentDto.Properties.AzureRegion = environmentSource.AzureRegion.ValueString()
	}

	if !environmentSource.BillingPolicyId.IsNull() && environmentSource.BillingPolicyId.ValueString() != "" {
		environmentDto.Properties.BillingPolicy = BillingPolicyDto{
			Id: environmentSource.BillingPolicyId.ValueString(),
		}
	}

	if !environmentSource.EnvironmentGroupId.IsNull() && !environmentSource.EnvironmentGroupId.IsUnknown() {
		environmentDto.Properties.ParentEnvironmentGroup = &ParentEnvironmentGroupDto{Id: environmentSource.EnvironmentGroupId.ValueString()}
	}

	if !environmentSource.Dataverse.IsNull() && !environmentSource.Dataverse.IsUnknown() {
		var dataverseSourceModel DataverseSourceModel
		environmentSource.Dataverse.As(ctx, &dataverseSourceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		environmentDto.Properties.DataBaseType = "CommonDataService"
		linkedMetadata, err := convertEnvironmentCreateLinkEnvironmentMetadataDtoFromDataverseSourceModel(ctx, environmentSource.Dataverse)
		if err != nil {
			return nil, err
		}
		environmentDto.Properties.LinkedEnvironmentMetadata = linkedMetadata
	}
	return environmentDto, nil
}

func convertEnvironmentCreateLinkEnvironmentMetadataDtoFromDataverseSourceModel(ctx context.Context, dataverse types.Object) (*createLinkEnvironmentMetadataDto, error) {
	if !dataverse.IsNull() && !dataverse.IsUnknown() {
		var dataverseSourceModel DataverseSourceModel
		dataverse.As(ctx, &dataverseSourceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		var templateMetadataObject *createTemplateMetadataDto
		if dataverseSourceModel.TemplateMetadata.ValueString() != "" {
			err := json.Unmarshal([]byte(dataverseSourceModel.TemplateMetadata.ValueString()), &templateMetadataObject)
			if err != nil {
				return nil, fmt.Errorf("error when unmarshalling template metadata %s; internal error: %v", dataverseSourceModel.TemplateMetadata.ValueString(), err)
			}
			if len(templateMetadataObject.PostProvisioningPackages) == 0 {
				templateMetadataObject = nil
			}
		}

		linkedEnvironmentMetadata := &createLinkEnvironmentMetadataDto{
			BaseLanguage:    int(dataverseSourceModel.LanguageName.ValueInt64()),
			SecurityGroupId: dataverseSourceModel.SecurityGroupId.ValueString(),
			Currency: createCurrencyDto{
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

func convertSourceModelFromEnvironmentDto(environmentDto EnvironmentDto, currencyCode *string, templateMetadata *createTemplateMetadataDto, templates []string, timeout timeouts.Value) (*SourceModel, error) {
	model := &SourceModel{
		Timeouts:        timeout,
		Description:     types.StringValue(environmentDto.Properties.Description),
		Id:              types.StringValue(environmentDto.Name),
		DisplayName:     types.StringValue(environmentDto.Properties.DisplayName),
		Location:        types.StringValue(environmentDto.Location),
		AzureRegion:     types.StringValue(environmentDto.Properties.AzureRegion),
		EnvironmentType: types.StringValue(environmentDto.Properties.EnvironmentSku),
		Cadence:         types.StringValue(environmentDto.Properties.UpdateCadence.Id),
	}

	convertBillingPolicyModelFromDto(environmentDto, model)
	convertEnvironmentGroupFromDto(environmentDto, model)
	convertEnterprisePolicyModelFromDto(environmentDto, model)

	attrTypesDataverseObject := map[string]attr.Type{
		"url":                          types.StringType,
		"domain":                       types.StringType,
		"organization_id":              types.StringType,
		"security_group_id":            types.StringType,
		"language_code":                types.Int64Type,
		"version":                      types.StringType,
		"linked_app_type":              types.StringType,
		"linked_app_id":                types.StringType,
		"linked_app_url":               types.StringType,
		"currency_code":                types.StringType,
		"templates":                    types.ListType{ElemType: types.StringType},
		"template_metadata":            types.StringType,
		"administration_mode_enabled":  types.BoolType,
		"background_operation_enabled": types.BoolType,
		"unique_name":                  types.StringType,
	}

	attrValuesProductProperties := map[string]attr.Value{}
	model.Dataverse = types.ObjectNull(attrTypesDataverseObject)

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
		attrValuesProductProperties["unique_name"] = types.StringValue(environmentDto.Properties.LinkedEnvironmentMetadata.UniqueName)
		if environmentDto.Properties.States != nil && environmentDto.Properties.States.Runtime != nil && environmentDto.Properties.States.Runtime.Id == "AdminMode" {
			attrValuesProductProperties["administration_mode_enabled"] = types.BoolValue(true)
		} else {
			attrValuesProductProperties["administration_mode_enabled"] = types.BoolValue(false)
		}
		if environmentDto.Properties.LinkedEnvironmentMetadata.BackgroundOperationsState == "Enabled" {
			attrValuesProductProperties["background_operation_enabled"] = types.BoolValue(true)
		} else {
			attrValuesProductProperties["background_operation_enabled"] = types.BoolValue(false)
		}

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
		model.Dataverse = types.ObjectValueMust(attrTypesDataverseObject, attrValuesProductProperties)
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
		attrValuesProductProperties["background_operation_enabled"] = types.BoolNull()
		attrValuesProductProperties["administration_mode_enabled"] = types.BoolNull()
		attrValuesProductProperties["environment_group_id"] = types.StringNull()
		attrValuesProductProperties["unique_name"] = types.StringNull()
	}
	return model, nil
}

func convertEnvironmentGroupFromDto(environmentDto EnvironmentDto, model *SourceModel) {
	if environmentDto.Properties.ParentEnvironmentGroup != nil {
		model.EnvironmentGroupId = types.StringValue(environmentDto.Properties.ParentEnvironmentGroup.Id)
	} else {
		model.EnvironmentGroupId = types.StringValue("")
	}
}

func convertBillingPolicyModelFromDto(environmentDto EnvironmentDto, model *SourceModel) {
	if environmentDto.Properties.BillingPolicy != nil {
		model.BillingPolicyId = types.StringValue(environmentDto.Properties.BillingPolicy.Id)
	} else {
		model.BillingPolicyId = types.StringValue("")
	}
}

func convertEnterprisePolicyModelFromDto(environmentDto EnvironmentDto, model *SourceModel) {
	enterprisePolicyAttrType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":      types.StringType,
			"id":        types.StringType,
			"location":  types.StringType,
			"system_id": types.StringType,
			"status":    types.StringType,
		},
	}
	if environmentDto.Properties.EnterprisePolicies != nil {
		if environmentDto.Properties.EnterprisePolicies.Vnets != nil {
			model.EnterprisePolicies = types.SetValueMust(enterprisePolicyAttrType, []attr.Value{
				types.ObjectValueMust(
					map[string]attr.Type{
						"type":      types.StringType,
						"id":        types.StringType,
						"location":  types.StringType,
						"system_id": types.StringType,
						"status":    types.StringType,
					},
					map[string]attr.Value{
						"type":      types.StringValue("NetworkInjection"),
						"id":        types.StringValue(environmentDto.Properties.EnterprisePolicies.Vnets.Id),
						"location":  types.StringValue(environmentDto.Properties.EnterprisePolicies.Vnets.Location),
						"system_id": types.StringValue(environmentDto.Properties.EnterprisePolicies.Vnets.SystemId),
						"status":    types.StringValue(environmentDto.Properties.EnterprisePolicies.Vnets.LinkStatus),
					},
				),
			})
		}
		if environmentDto.Properties.EnterprisePolicies.CustomerManagedKeys != nil {
			model.EnterprisePolicies = types.SetValueMust(enterprisePolicyAttrType, []attr.Value{
				types.ObjectValueMust(
					map[string]attr.Type{
						"type":      types.StringType,
						"id":        types.StringType,
						"location":  types.StringType,
						"system_id": types.StringType,
						"status":    types.StringType,
					},
					map[string]attr.Value{
						"type":      types.StringValue("Encryption"),
						"id":        types.StringValue(environmentDto.Properties.EnterprisePolicies.CustomerManagedKeys.Id),
						"location":  types.StringValue(environmentDto.Properties.EnterprisePolicies.CustomerManagedKeys.Location),
						"system_id": types.StringValue(environmentDto.Properties.EnterprisePolicies.CustomerManagedKeys.SystemId),
						"status":    types.StringValue(environmentDto.Properties.EnterprisePolicies.CustomerManagedKeys.LinkStatus),
					},
				),
			})
		}
	} else {
		model.EnterprisePolicies = types.SetValueMust(enterprisePolicyAttrType, []attr.Value{})
	}
}
