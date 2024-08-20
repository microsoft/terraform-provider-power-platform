// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	modifiers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/modifiers"
	licensing "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/licensing"
)

var _ resource.Resource = &EnvironmentResource{}
var _ resource.ResourceWithImportState = &EnvironmentResource{}

func NewEnvironmentResource() resource.Resource {
	return &EnvironmentResource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_environment",
	}
}

type EnvironmentResource struct {
	EnvironmentClient EnvironmentClient
	LicensingClient   licensing.LicensingClient
	ProviderTypeName  string
	TypeName          string
}

func (r *EnvironmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

func (r *EnvironmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource manages a PowerPlatform environment",
		Description:         "This resource manages a PowerPlatform environment",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				MarkdownDescription: "Unique environment id (guid)",
				Description:         "Unique environment id (guid)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"location": schema.StringAttribute{
				Description:         "Location of the environment (europe, unitedstates etc.). Can be queried using the `powerplatform_locations` data source. The region of your Entra tenant may [limit the available locations for Power Platform](https://learn.microsoft.com/power-platform/admin/regions-overview#who-can-create-environments-in-these-regions). Changing this property after environment creation will result in a destroy and recreation of the environment (you can use the [`prevent_destroy` lifecycle metatdata](https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#prevent_destroy) as an added safeguard to prevent accidental deletion of environments).",
				MarkdownDescription: "Location of the environment (europe, unitedstates etc.). Can be queried using the `powerplatform_locations` data source. The region of your Entra tenant may [limit the available locations for Power Platform](https://learn.microsoft.com/power-platform/admin/regions-overview#who-can-create-environments-in-these-regions). Changing this property after environment creation will result in a destroy and recreation of the environment (you can use the [`prevent_destroy` lifecycle metatdata](https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#prevent_destroy) as an added safeguard to prevent accidental deletion of environments).",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"azure_region": schema.StringAttribute{
				Description:         "Azure region of the environment (westeurope, eastus etc.). Can be queried using the `powerplatform_locations` data source. This property should only be set if absolutely necessary like when trying to create an environment in the same Azure region as Azure resources or Fabric capacity.  Changing this property after environment creation will result in a destroy and recreation of the environment (you can use the [`prevent_destroy` lifecycle metatdata](https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#prevent_destroy) as an added safeguard to prevent accidental deletion of environments).",
				MarkdownDescription: "Azure region of the environment (westeurope, eastus etc.). Can be queried using the `powerplatform_locations` data source. This property should only be set if absolutely necessary like when trying to create an environment in the same Azure region as Azure resources or Fabric capacity.  Changing this property after environment creation will result in a destroy and recreation of the environment (you can use the [`prevent_destroy` lifecycle metatdata](https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#prevent_destroy) as an added safeguard to prevent accidental deletion of environments).",
				Required:            false,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_type": schema.StringAttribute{
				Description:         "Type of the environment (Sandbox, Production etc.)",
				MarkdownDescription: "Type of the environment (Sandbox, Production etc.)",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(EnvironmentTypes...),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name",
				Description:         "Display name",
				Required:            true,
			},
			"billing_policy_id": &schema.StringAttribute{
				Description:         "Billing policy id (guid) for pay-as-you-go environments using Azure subscription billing",
				MarkdownDescription: "Billing policy id (guid) for pay-as-you-go environments using Azure subscription billing",
				Optional:            true,
				Computed:            true,
			},
			"dataverse": schema.SingleNestedAttribute{
				MarkdownDescription: "Dataverse environment details",
				Description:         "Dataverse environment details",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					modifiers.RequireReplaceObjectToEmptyModifier(),
				},
				Attributes: map[string]schema.Attribute{
					"currency_code": schema.StringAttribute{
						Description:         "Unique currency code",
						MarkdownDescription: "Unique currency name",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.RequireReplaceStringFromNonEmptyPlanModifier(),
						},
					},
					"url": schema.StringAttribute{
						Description:         "Url of the environment",
						MarkdownDescription: "Url of the environment",
						Computed:            true,
					},
					"domain": schema.StringAttribute{
						Description:         "Domain name of the environment",
						MarkdownDescription: "Domain name of the environment",
						Optional:            true,
						Computed:            true,
					},
					"organization_id": schema.StringAttribute{
						Description:         "Unique organization id (guid)",
						MarkdownDescription: "Unique organization id (guid)",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"security_group_id": schema.StringAttribute{
						Description:         "Unique security group id (guid).  For an empty security group, set this property to 0000000-0000-0000-0000-000000000000",
						MarkdownDescription: "Unique security group id (guid).  For an empty security group, set this property to `0000000-0000-0000-0000-000000000000`",
						Required:            true,
					},
					"language_code": schema.Int64Attribute{
						Description:         "Unique language LCID (integer)",
						MarkdownDescription: "Unique language LCID (integer)",
						Required:            true,
						PlanModifiers: []planmodifier.Int64{
							modifiers.RequireReplaceIntAttributePlanModifier(),
						},
					},
					"version": schema.StringAttribute{
						Description:         "Version of the environment",
						MarkdownDescription: "Version of the environment",
						Computed:            true,
					},
					"templates": schema.ListAttribute{
						Description:         "The selected instance provisioning template (if any)",
						MarkdownDescription: "The selected instance provisioning template (if any)",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"template_metadata": schema.StringAttribute{
						Description:         "Additional D365 environment template metadata (if any)",
						MarkdownDescription: "Additional D365 environment template metadata (if any)",
						Optional:            true,
					},
					"linked_app_type": schema.StringAttribute{
						Description:         "The type of the linked D365 application",
						MarkdownDescription: "The type of the linked D365 application",
						Computed:            true,
					},
					"linked_app_id": schema.StringAttribute{
						Description:         "The GUID of the linked D365 application",
						MarkdownDescription: "The GUID of the linked D365 application",
						Computed:            true,
					},
					"linked_app_url": schema.StringAttribute{
						Description:         "The URL of the linked D365 application",
						MarkdownDescription: "The URL of the linked D365 application",
						Computed:            true,
					},
				},
			},
		},
	}
}

func (r *EnvironmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientApi := req.ProviderData.(*api.ProviderClient).Api

	if clientApi == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.EnvironmentClient = NewEnvironmentClient(clientApi)
	r.LicensingClient = licensing.NewLicensingClient(clientApi)
}

func (r *EnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *EnvironmentSourceModel

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	envToCreate, err := ConvertCreateEnvironmentDtoFromSourceModel(ctx, *plan)
	if err != nil {
		resp.Diagnostics.AddError("Error when converting source model to create environment dto", err.Error())
	}

	err = locationValidator(r.EnvironmentClient.Api, envToCreate.Location, envToCreate.Properties.AzureRegion)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Location validation failed for %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	if envToCreate.Properties.LinkedEnvironmentMetadata != nil {
		err = languageCodeValidator(r.EnvironmentClient.Api, envToCreate.Location, fmt.Sprintf("%d", envToCreate.Properties.LinkedEnvironmentMetadata.BaseLanguage))
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Language code validation failed for %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
			return
		}

		err = currencyCodeValidator(r.EnvironmentClient.Api, envToCreate.Location, envToCreate.Properties.LinkedEnvironmentMetadata.Currency.Code)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Currency code validation failed for %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
			return
		}
	}

	envDto, err := r.EnvironmentClient.CreateEnvironment(ctx, *envToCreate)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	var currencyCode string
	var templateMetadata *EnvironmentCreateTemplateMetadata = nil
	var templates []string = nil
	if envToCreate.Properties.LinkedEnvironmentMetadata != nil {
		currencyCode = envToCreate.Properties.LinkedEnvironmentMetadata.Currency.Code

		//because BAPI does not retrieve template info after create, we have to rewrite it
		templateMetadata = envToCreate.Properties.LinkedEnvironmentMetadata.TemplateMetadata
		templates = envToCreate.Properties.LinkedEnvironmentMetadata.Templates
	}

	newPlan, err := ConvertSourceModelFromEnvironmentDto(*envDto, &currencyCode, templateMetadata, templates)
	if err != nil {
		resp.Diagnostics.AddError("Error when converting environment to source model", err.Error())
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID %s", plan.Id.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &newPlan)...)

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE END: %s", r.ProviderTypeName))

}

func (r *EnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *EnvironmentSourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	envDto, err := r.EnvironmentClient.GetEnvironment(ctx, state.Id.ValueString())
	if err != nil {
		if helpers.Code(err) == helpers.ERROR_OBJECT_NOT_FOUND {
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
			return
		}
	}

	currencyCode := ""
	defaultCurrency, err := r.EnvironmentClient.GetDefaultCurrencyForEnvironment(ctx, envDto.Name)
	if err != nil {
		if helpers.Code(err) != helpers.ERROR_ENVIRONMENT_URL_NOT_FOUND {
			resp.Diagnostics.AddWarning(fmt.Sprintf("Error when reading default currency for environment %s", envDto.Name), err.Error())
		}

		if !state.Dataverse.IsNull() && !state.Dataverse.IsUnknown() {
			var dataverseSourceModel DataverseSourceModel
			state.Dataverse.As(ctx, &dataverseSourceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
			currencyCode = dataverseSourceModel.CurrencyCode.ValueString()
		}
	} else {
		currencyCode = defaultCurrency.IsoCurrencyCode
	}

	var templateMetadata *EnvironmentCreateTemplateMetadata = nil
	var templates []string = nil
	if !state.Dataverse.IsNull() && !state.Dataverse.IsUnknown() {
		dv, err := ConvertEnvironmentCreateLinkEnvironmentMetadataDtoFromDataverseSourceModel(ctx, state.Dataverse)
		if err != nil {
			resp.Diagnostics.AddError("Error when converting dataverse source model to create link environment metadata", err.Error())
			return
		}
		if dv != nil {
			templateMetadata = dv.TemplateMetadata
			templates = dv.Templates
		}
	}
	newState, err := ConvertSourceModelFromEnvironmentDto(*envDto, &currencyCode, templateMetadata, templates)
	if err != nil {
		resp.Diagnostics.AddError("Error when converting environment to source model", err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("READ: %s_environment with id %s", r.ProviderTypeName, state.Id.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE END: %s", r.ProviderTypeName))
}

func (r *EnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *EnvironmentSourceModel

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state *EnvironmentSourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	environmentDto := EnvironmentDto{
		Id:       plan.Id.ValueString(),
		Name:     plan.DisplayName.ValueString(),
		Type:     plan.EnvironmentType.ValueString(),
		Location: plan.Location.ValueString(),
		Properties: EnvironmentPropertiesDto{
			DisplayName:    plan.DisplayName.ValueString(),
			EnvironmentSku: plan.EnvironmentType.ValueString(),
		},
	}

	if !plan.BillingPolicyId.IsNull() && plan.BillingPolicyId.ValueString() != "" {
		environmentDto.Properties.BillingPolicy = &BillingPolicyDto{
			Id: plan.BillingPolicyId.ValueString(),
		}
	}

	var currencyCode string
	if !IsDataverseEnvironmentEmpty(ctx, state) && IsDataverseEnvironmentEmpty(ctx, plan) {
		resp.Diagnostics.AddError("Cannot remove dataverse environment from environment", "Cannot remove dataverse environment from environment")
		return
	} else if !IsDataverseEnvironmentEmpty(ctx, state) && !IsDataverseEnvironmentEmpty(ctx, plan) {

		var dataverseSourcePlanModel DataverseSourceModel
		plan.Dataverse.As(ctx, &dataverseSourcePlanModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		currencyCode = dataverseSourcePlanModel.CurrencyCode.ValueString()

		environmentDto.Properties.LinkedEnvironmentMetadata = &LinkedEnvironmentMetadataDto{
			SecurityGroupId: dataverseSourcePlanModel.SecurityGroupId.ValueString(),
			DomainName:      dataverseSourcePlanModel.Domain.ValueString(),
		}

		var dataverseSourceStateModel DataverseSourceModel
		state.Dataverse.As(ctx, &dataverseSourceStateModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if dataverseSourceStateModel.Domain.ValueString() != dataverseSourcePlanModel.Domain.ValueString() && !dataverseSourcePlanModel.Domain.IsNull() && dataverseSourcePlanModel.Domain.ValueString() != "" {
			environmentDto.Properties.LinkedEnvironmentMetadata.DomainName = dataverseSourcePlanModel.Domain.ValueString()
		}

		if !dataverseSourcePlanModel.LinkedAppId.IsNull() && dataverseSourcePlanModel.LinkedAppId.ValueString() != "" {
			environmentDto.Properties.LinkedAppMetadata = &LinkedAppMetadataDto{
				Type: dataverseSourcePlanModel.LinkedAppType.ValueString(),
				Id:   dataverseSourcePlanModel.LinkedAppId.ValueString(),
				Url:  dataverseSourcePlanModel.LinkedAppURL.ValueString(),
			}
		} else {
			environmentDto.Properties.LinkedAppMetadata = nil
		}
	} else {

		linkedMetadataDto, err := ConvertEnvironmentCreateLinkEnvironmentMetadataDtoFromDataverseSourceModel(ctx, plan.Dataverse)
		if err != nil {
			resp.Diagnostics.AddError("Error when converting dataverse source model to create link environment metadata dto", err.Error())
			return
		}

		_, err = r.EnvironmentClient.AddDataverseToEnvironment(ctx, plan.Id.ValueString(), *linkedMetadataDto)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error when adding dataverse to environment %s", plan.Id.ValueString()), err.Error())
			return
		}
		currencyCode = linkedMetadataDto.Currency.Code
	}

	if !state.BillingPolicyId.IsNull() &&
		!state.BillingPolicyId.IsUnknown() &&
		state.BillingPolicyId.ValueString() != "" {

		tflog.Debug(ctx, fmt.Sprintf("Removing environment %s from billing policy %s", state.Id.ValueString(), state.BillingPolicyId.ValueString()))
		err := r.LicensingClient.RemoveEnvironmentsToBillingPolicy(ctx, state.BillingPolicyId.ValueString(), []string{state.Id.ValueString()})
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error when removing environment %s from billing policy %s", state.Id.ValueString(), state.BillingPolicyId.ValueString()), err.Error())
			return
		}
	}

	if !plan.BillingPolicyId.IsNull() &&
		!plan.BillingPolicyId.IsUnknown() &&
		plan.BillingPolicyId.ValueString() != "" {

		tflog.Debug(ctx, fmt.Sprintf("Adding environment %s to billing policy %s", plan.Id.ValueString(), plan.BillingPolicyId.ValueString()))
		err := r.LicensingClient.AddEnvironmentsToBillingPolicy(ctx, plan.BillingPolicyId.ValueString(), []string{plan.Id.ValueString()})
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error when adding environment %s to billing policy %s", plan.Id.ValueString(), plan.BillingPolicyId.ValueString()), err.Error())
			return
		}
	}

	envDto, err := r.EnvironmentClient.UpdateEnvironment(ctx, plan.Id.ValueString(), environmentDto)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.ProviderTypeName), err.Error())
		return
	}

	var templateMetadata *EnvironmentCreateTemplateMetadata = nil
	var templates []string = nil
	if !state.Dataverse.IsNull() && !state.Dataverse.IsUnknown() {
		dv, err := ConvertEnvironmentCreateLinkEnvironmentMetadataDtoFromDataverseSourceModel(ctx, state.Dataverse)
		if err != nil {
			resp.Diagnostics.AddError("Error when converting dataverse source model to create link environment metadata", err.Error())
			return
		}
		if dv != nil {
			templateMetadata = dv.TemplateMetadata
			templates = dv.Templates
		}
	}

	newPlan, err := ConvertSourceModelFromEnvironmentDto(*envDto, &currencyCode, templateMetadata, templates)
	if err != nil {
		resp.Diagnostics.AddError("Error when converting environment to source model", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newPlan)...)

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *EnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *EnvironmentSourceModel

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.EnvironmentClient.DeleteEnvironment(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *EnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
