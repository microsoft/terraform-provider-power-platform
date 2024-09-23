// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
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

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/modifiers"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/licensing"
)

var _ resource.Resource = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}

func NewEnvironmentResource() resource.Resource {
	return &Resource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment",
		},
	}
}

type Resource struct {
	helpers.TypeInfo
	EnvironmentClient Client
	LicensingClient   licensing.Client
}

// Metadata returns the full name of the resource type.
func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource manages a PowerPlatform environment",
		Description:         "This resource manages a PowerPlatform environment",

		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique environment id (guid)",
				Description:         "Unique environment id (guid)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_group_id": schema.StringAttribute{
				MarkdownDescription: "Unique environment group id (guid) that the environment belongs to. See [Environment groups](https://learn.microsoft.com/en-us/power-platform/admin/environment-groups) for more information.",
				Computed:            true,
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the environment",
				Optional:            true,
				Computed:            true,
			},
			"cadence": schema.StringAttribute{
				MarkdownDescription: "Cadence of updates for the environment (Frequent, Moderate). For more information check [here](https://learn.microsoft.com/en-us/power-platform/admin/create-environment#setting-an-environment-refresh-cadence).",
				Optional:            true,
				Computed:            true,
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

					"administration_mode_enabled": schema.BoolAttribute{
						MarkdownDescription: "Select to enable administration mode for the environment. See [Admin mode](https://learn.microsoft.com/en-us/power-platform/admin/admin-mode) for more information. ",
						Computed:            true,
						Optional:            true,
					},
					"background_operation_enabled": schema.BoolAttribute{
						MarkdownDescription: "Indicates if background operation is enabled",
						Optional:            true,
						Computed:            true,
					},
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

func (d *Resource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.RequiredTogether(
			path.Root("dataverse").AtName("administration_mode_enabled").Expression(),
			path.Root("dataverse").AtName("background_operation_enabled").Expression(),
		),
	}
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
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
	tflog.Debug(ctx, "Successfully created clients")
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *SourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	envToCreate, err := convertCreateEnvironmentDtoFromSourceModel(ctx, *plan)
	if err != nil {
		resp.Diagnostics.AddError("Error when converting source model to create environment dto", err.Error())
	}

	err = r.EnvironmentClient.LocationValidator(ctx, envToCreate.Location, envToCreate.Properties.AzureRegion)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Location validation failed for %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	// If it's dataverse environment, validate the currency and language code
	if envToCreate.Properties.LinkedEnvironmentMetadata != nil {
		err = languageCodeValidator(ctx, r.EnvironmentClient.Api, envToCreate.Location, fmt.Sprintf("%d", envToCreate.Properties.LinkedEnvironmentMetadata.BaseLanguage))
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Language code validation failed for %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
			return
		}

		err = currencyCodeValidator(ctx, r.EnvironmentClient.Api, envToCreate.Location, envToCreate.Properties.LinkedEnvironmentMetadata.Currency.Code)
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
	var templateMetadata *CreateTemplateMetadata
	var templates []string
	if envToCreate.Properties.LinkedEnvironmentMetadata != nil {
		currencyCode = envToCreate.Properties.LinkedEnvironmentMetadata.Currency.Code

		// because BAPI does not retrieve template info after create, we have to rewrite it
		templateMetadata = envToCreate.Properties.LinkedEnvironmentMetadata.TemplateMetadata
		templates = envToCreate.Properties.LinkedEnvironmentMetadata.Templates
	}

	newState, err := convertSourceModelFromEnvironmentDto(*envDto, &currencyCode, templateMetadata, templates, plan.Timeouts)
	if err != nil {
		resp.Diagnostics.AddError("Error when converting environment to source model", err.Error())
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID %s", plan.Id.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

// Read reads the resource state from the remote system. If the resource does not exist, the state should be removed from the state store.
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *SourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	envDto, err := r.EnvironmentClient.GetEnvironment(ctx, state.Id.ValueString())
	if err != nil {
		if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	currencyCode := ""
	defaultCurrency, err := r.EnvironmentClient.GetDefaultCurrencyForEnvironment(ctx, envDto.Name)
	if err != nil {
		if customerrors.Code(err) != customerrors.ERROR_ENVIRONMENT_URL_NOT_FOUND {
			// This is only a warning because you may have BAPI access to the environment but not WebAPI access to dataverse to get currency.
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

	var templateMetadata *CreateTemplateMetadata
	var templates []string
	if !state.Dataverse.IsNull() && !state.Dataverse.IsUnknown() {
		dv, err := convertEnvironmentCreateLinkEnvironmentMetadataDtoFromDataverseSourceModel(ctx, state.Dataverse)
		if err != nil {
			resp.Diagnostics.AddError("Error when converting dataverse source model to create link environment metadata", err.Error())
			return
		}
		if dv != nil {
			templateMetadata = dv.TemplateMetadata
			templates = dv.Templates
		}
	}
	newState, err := convertSourceModelFromEnvironmentDto(*envDto, &currencyCode, templateMetadata, templates, state.Timeouts)
	if err != nil {
		resp.Diagnostics.AddError("Error when converting environment to source model", err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("READ: %s_environment with id %s", r.ProviderTypeName, state.Id.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *SourceModel
	var state *SourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	environmentDto := EnvironmentDto{
		Id:       plan.Id.ValueString(),
		Name:     plan.DisplayName.ValueString(),
		Type:     plan.EnvironmentType.ValueString(),
		Location: plan.Location.ValueString(),
		Properties: EnviromentPropertiesDto{
			DisplayName:    plan.DisplayName.ValueString(),
			EnvironmentSku: plan.EnvironmentType.ValueString(),
		},
	}

	if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
		environmentDto.Properties.Description = plan.Description.ValueString()
	}

	if !plan.Cadence.IsNull() && plan.Cadence.ValueString() != "" {
		environmentDto.Properties.UpdateCadence = &UpdateCadenceDto{
			Id: plan.Cadence.ValueString(),
		}
	}

	if !plan.EnvironmentGroupId.IsNull() && !plan.EnvironmentGroupId.IsUnknown() {
		envGroupId := constants.ZERO_UUID
		if plan.EnvironmentGroupId.ValueString() != "" && plan.EnvironmentGroupId.ValueString() != constants.ZERO_UUID {
			envGroupId = plan.EnvironmentGroupId.ValueString()
		}
		environmentDto.Properties.ParentEnvironmentGroup = &ParentEnvironmentGroupDto{
			Id: envGroupId,
		}
	}

	if !plan.BillingPolicyId.IsNull() && plan.BillingPolicyId.ValueString() != "" {
		environmentDto.Properties.BillingPolicy = &BillingPolicyDto{
			Id: plan.BillingPolicyId.ValueString(),
		}
	}

	var currencyCode string
	if !isDataverseEnvironmentEmpty(ctx, state) && !isDataverseEnvironmentEmpty(ctx, plan) {
		currencyCode = updateExistingDataverse(ctx, plan, &environmentDto, state)
	} else if isDataverseEnvironmentEmpty(ctx, state) && !isDataverseEnvironmentEmpty(ctx, plan) {
		code, err := addDataverse(ctx, plan, r)
		if err != nil {
			resp.Diagnostics.AddError("Error when creating new dataverse environment", err.Error())
			return
		}
		currencyCode = code
	}

	if !state.BillingPolicyId.IsNull() && !state.BillingPolicyId.IsUnknown() && state.BillingPolicyId.ValueString() != "" {
		tflog.Debug(ctx, fmt.Sprintf("Removing environment %s from billing policy %s", state.Id.ValueString(), state.BillingPolicyId.ValueString()))
		err := r.LicensingClient.RemoveEnvironmentsToBillingPolicy(ctx, state.BillingPolicyId.ValueString(), []string{state.Id.ValueString()})
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error when removing environment %s from billing policy %s", state.Id.ValueString(), state.BillingPolicyId.ValueString()), err.Error())
			return
		}
	}

	if !plan.BillingPolicyId.IsNull() && !plan.BillingPolicyId.IsUnknown() && plan.BillingPolicyId.ValueString() != "" {
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

	var templateMetadata *CreateTemplateMetadata
	var templates []string
	if !state.Dataverse.IsNull() && !state.Dataverse.IsUnknown() {
		dv, err := convertEnvironmentCreateLinkEnvironmentMetadataDtoFromDataverseSourceModel(ctx, state.Dataverse)
		if err != nil {
			resp.Diagnostics.AddError("Error when converting dataverse source model to create link environment metadata", err.Error())
			return
		}
		if dv != nil {
			templateMetadata = dv.TemplateMetadata
			templates = dv.Templates
		}
	}

	newPlan, err := convertSourceModelFromEnvironmentDto(*envDto, &currencyCode, templateMetadata, templates, plan.Timeouts)
	if err != nil {
		resp.Diagnostics.AddError("Error when converting environment to source model", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newPlan)...)
}

func addDataverse(ctx context.Context, plan *SourceModel, r *Resource) (string, error) {
	linkedMetadataDto, err := convertEnvironmentCreateLinkEnvironmentMetadataDtoFromDataverseSourceModel(ctx, plan.Dataverse)
	if err != nil {
		return "", fmt.Errorf("Error when converting dataverse source model to create link environment metadata dto: %s", err.Error())
	}

	_, err = r.EnvironmentClient.AddDataverseToEnvironment(ctx, plan.Id.ValueString(), *linkedMetadataDto)
	if err != nil {
		return "", fmt.Errorf("Error when adding dataverse to environment %s: %s", plan.Id.ValueString(), err.Error())
	}
	return linkedMetadataDto.Currency.Code, nil
}

func updateExistingDataverse(ctx context.Context, plan *SourceModel, environmentDto *EnvironmentDto, state *SourceModel) string {
	var dataverseSourcePlanModel DataverseSourceModel
	plan.Dataverse.As(ctx, &dataverseSourcePlanModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

	environmentDto.Properties.LinkedEnvironmentMetadata = &LinkedEnvironmentMetadataDto{
		SecurityGroupId: dataverseSourcePlanModel.SecurityGroupId.ValueString(),
		DomainName:      dataverseSourcePlanModel.Domain.ValueString(),
	}

	if !dataverseSourcePlanModel.AdministrationMode.IsNull() && !dataverseSourcePlanModel.AdministrationMode.IsUnknown() {
		if dataverseSourcePlanModel.AdministrationMode.ValueBool() {
			environmentDto.Properties.States = &StatesEnvironmentDto{
				Runtime: &RuntimeEnvironmentDto{
					Id: "AdminMode",
				},
			}
		} else {
			environmentDto.Properties.States = &StatesEnvironmentDto{
				Runtime: &RuntimeEnvironmentDto{
					Id: "Enabled",
				},
			}
		}
	}

	if !dataverseSourcePlanModel.BackgroundOperation.IsNull() && !dataverseSourcePlanModel.BackgroundOperation.IsUnknown() {
		if dataverseSourcePlanModel.BackgroundOperation.ValueBool() {
			environmentDto.Properties.LinkedEnvironmentMetadata.BackgroundOperationsState = "Enabled"
		} else {
			environmentDto.Properties.LinkedEnvironmentMetadata.BackgroundOperationsState = "Disabled"
		}
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

	return dataverseSourcePlanModel.CurrencyCode.ValueString()
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *SourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.EnvironmentClient.DeleteEnvironment(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
