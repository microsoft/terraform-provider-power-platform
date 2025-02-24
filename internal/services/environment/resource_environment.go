// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
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
	"github.com/microsoft/terraform-provider-power-platform/internal/validators"
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

	policyAttributeSchema := map[string]schema.Attribute{
		"type": schema.StringAttribute{
			MarkdownDescription: "Type of the policy according to [schema definition](https://learn.microsoft.com/en-us/azure/templates/microsoft.powerplatform/enterprisepolicies?pivots=deployment-language-terraform#enterprisepolicies-2)",
			Computed:            true,
		},
		"id": schema.StringAttribute{
			MarkdownDescription: "Id (guid)",
			Computed:            true,
		},
		"location": schema.StringAttribute{
			MarkdownDescription: "Location of the policy",
			Computed:            true,
		},
		"system_id": schema.StringAttribute{
			MarkdownDescription: "System id (guid)",
			Computed:            true,
		},
		"status": schema.StringAttribute{
			MarkdownDescription: "Link status of the policy",
			Computed:            true,
		},
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource manages a PowerPlatform environment.",

		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Environment id (guid)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_group_id": schema.StringAttribute{
				MarkdownDescription: "Environment group id (guid) that the environment belongs to. See [Environment groups](https://learn.microsoft.com/en-us/power-platform/admin/environment-groups) for more information.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(helpers.GuidOrEmptyValueRegex), "environment_group_id must be a valid environment group id guid"),
					stringvalidator.AlsoRequires(path.Root("dataverse").Expression()),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the environment",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cadence": schema.StringAttribute{
				MarkdownDescription: "Cadence of updates for the environment (Frequent, Moderate). For more information check [here](https://learn.microsoft.com/en-us/power-platform/admin/create-environment#setting-an-environment-refresh-cadence).",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"release_cycle": schema.StringAttribute{
				MarkdownDescription: "Gives you the ability to create environments that are updated first. This allows you to experience and validate scenarios that are important to you before any updates reach your business-critical applications. See [more](https://learn.microsoft.com/en-us/power-platform/admin/early-release).",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(ReleaseCycleTypes...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "Location of the environment (europe, unitedstates etc.). Can be queried using the `powerplatform_locations` data source. The region of your Entra tenant may [limit the available locations for Power Platform](https://learn.microsoft.com/power-platform/admin/regions-overview#who-can-create-environments-in-these-regions). Changing this property after environment creation will result in a destroy and recreation of the environment (you can use the [`prevent_destroy` lifecycle metatdata](https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#prevent_destroy) as an added safeguard to prevent accidental deletion of environments).",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"azure_region": schema.StringAttribute{
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
				MarkdownDescription: "Type of the environment (Sandbox, Production etc.)",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(EnvironmentTypes...),
					validators.OtherFieldRequiredWhenValueOf(path.Root("owner_id").Expression(), nil, regexp.MustCompile(EnvironmentTypesDeveloperOnlyRegex), "owner_id must be set when environment_type is `Developer`"),
				},
			},
			"owner_id": schema.StringAttribute{
				MarkdownDescription: "Entra ID  user id (guid) of the environment owner when creating developer environment",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Root("dataverse").AtName("security_group_id").Expression()),
					stringvalidator.AlsoRequires(path.Root("dataverse").Expression()),
					validators.OtherFieldRequiredWhenValueOf(path.Root("environment_type").Expression(), regexp.MustCompile(EnvironmentTypesDeveloperOnlyRegex), nil, "owner_id can be used only when environment_type is `Developer`"),
				},
			},
			"allow_bing_search": schema.BoolAttribute{
				MarkdownDescription: "Allow Bing search in the environment",
				Optional:            true,
				Computed:            true,
			},
			"allow_moving_data_across_regions": schema.BoolAttribute{
				MarkdownDescription: "Allow moving data across regions",
				Optional:            true,
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name",
				Required:            true,
			},
			"billing_policy_id": &schema.StringAttribute{
				MarkdownDescription: "Billing policy id (guid) for pay-as-you-go environments using Azure subscription billing",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enterprise_policies": schema.SetNestedAttribute{
				MarkdownDescription: "Enterprise policies for the environment. See [Enterprise policies](https://learn.microsoft.com/en-us/power-platform/admin/enterprise-policies) for more details.",
				Computed:            true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: policyAttributeSchema,
				},
			},
			"dataverse": schema.SingleNestedAttribute{
				MarkdownDescription: "Dataverse environment details",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					modifiers.RequireReplaceObjectToEmptyModifier(),
				},
				Attributes: map[string]schema.Attribute{
					"unique_name": schema.StringAttribute{
						MarkdownDescription: "Unique name of the Dataverse environment",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"administration_mode_enabled": schema.BoolAttribute{
						MarkdownDescription: "Select to enable administration mode for the environment. See [Admin mode](https://learn.microsoft.com/en-us/power-platform/admin/admin-mode) for more information. ",
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"background_operation_enabled": schema.BoolAttribute{
						MarkdownDescription: "Indicates if background operation is enabled",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"currency_code": schema.StringAttribute{
						MarkdownDescription: "Currency name",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.RequireReplaceStringFromNonEmptyPlanModifier(),
						},
					},
					"url": schema.StringAttribute{
						MarkdownDescription: "Url of the environment",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
							modifiers.SetStringAttributeUnknownOnlyIfSecondAttributeChange(path.Root("dataverse").AtName("domain")),
						},
					},
					"domain": schema.StringAttribute{
						MarkdownDescription: "Domain name of the environment",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"organization_id": schema.StringAttribute{
						MarkdownDescription: "Organization id (guid)",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"security_group_id": schema.StringAttribute{
						MarkdownDescription: "Security group id (guid). For an empty security group, set this property to `0000000-0000-0000-0000-000000000000`",
						Optional:            true,
						Validators: []validator.String{
							validators.MakeFieldRequiredWhenOtherFieldDoesNotHaveValue(path.Root("environment_type").Expression(), regexp.MustCompile(EnvironmentTypesExceptDeveloperRegex), "dataverse.security_group_id is required for all environment_type values except `Developer`"),
						},
					},
					"language_code": schema.Int64Attribute{
						MarkdownDescription: "Language LCID (integer)",
						Required:            true,
						PlanModifiers: []planmodifier.Int64{
							modifiers.RequireReplaceIntAttributePlanModifier(),
						},
					},
					"version": schema.StringAttribute{
						MarkdownDescription: "Version of the environment",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"templates": schema.ListAttribute{
						MarkdownDescription: "The selected instance provisioning template (if any). See [ERP-based template](https://learn.microsoft.com/en-us/power-platform/admin/unified-experience/tutorial-deploy-new-environment-with-erp-template?tabs=PPAC) for more information.",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"template_metadata": schema.StringAttribute{
						MarkdownDescription: "Additional D365 environment template metadata (if any)",
						Optional:            true,
					},
					"linked_app_type": schema.StringAttribute{
						MarkdownDescription: "The type of the linked D365 application",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"linked_app_id": schema.StringAttribute{
						MarkdownDescription: "The GUID of the linked D365 application",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"linked_app_url": schema.StringAttribute{
						MarkdownDescription: "The URL of the linked D365 application",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
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

	envToCreate, err := convertCreateEnvironmentDtoFromSourceModel(ctx, plan, r)

	if err != nil {
		resp.Diagnostics.AddError("Error when converting source model to create environment dto", err.Error())
	}

	err = r.EnvironmentClient.LocationValidator(ctx, envToCreate.Location, envToCreate.Properties.AzureRegion)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Location validation failed for %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	/* ... check if non public cloud ... need PR 572 on main first*/
	/* if cloud is non public {
		allow_bing_search = false && allow_moving_data_across_regions = true => exception ("moving data can't be enabled for non public cloud")
	}*/
	if plan.Location.ValueString() == "unitedstates" && plan.AllowMovingDataAcrossRegions.ValueBool() {
		resp.Diagnostics.AddError(fmt.Sprintf("Location validation failed for %s_%s", r.ProviderTypeName, r.TypeName), "Moving data across regions is not supported in the unitedstates location")
		return
	}
	if plan.Location.ValueString() != "unitedstates" && plan.AllowBingSearch.ValueBool() && !plan.AllowMovingDataAcrossRegions.ValueBool() {
		resp.Diagnostics.AddError(fmt.Sprintf("Location validation failed for %s_%s", r.ProviderTypeName, r.TypeName), "To enable AI generative features, moving data across regions must be enabled")
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

	if !plan.AllowBingSearch.IsNull() && !plan.AllowBingSearch.IsUnknown() {
		err := r.updateEnvironmentAiFeatures(ctx, envDto.Name, plan.AllowBingSearch.ValueBool(), plan.AllowMovingDataAcrossRegions.ValueBoolPointer())
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
			return
		}

		envDto, err = r.EnvironmentClient.GetEnvironment(ctx, envDto.Name)
		if err != nil {
			if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
			return
		}
	}

	var currencyCode string
	var templateMetadata *createTemplateMetadataDto
	var templates []string
	if envToCreate.Properties.LinkedEnvironmentMetadata != nil {
		currencyCode = envToCreate.Properties.LinkedEnvironmentMetadata.Currency.Code

		// because BAPI does not retrieve template info after create, we have to rewrite it
		templateMetadata = envToCreate.Properties.LinkedEnvironmentMetadata.TemplateMetadata
		templates = envToCreate.Properties.LinkedEnvironmentMetadata.Templates
	}

	newState, err := convertSourceModelFromEnvironmentDto(*envDto, &currencyCode, plan.OwnerId.ValueStringPointer(), templateMetadata, templates, plan.Timeouts, *r.EnvironmentClient.Api.Config)

	if !plan.AzureRegion.IsNull() && plan.AzureRegion.ValueString() != "" && (plan.AzureRegion.ValueString() != newState.AzureRegion.ValueString()) {
		resp.Diagnostics.AddAttributeError(path.Root("azure_region"), fmt.Sprintf("Provisioning environment in azure region '%s' failed", plan.AzureRegion.ValueString()), "Provisioning environment in azure region was not successful, please try other region in that location or try again later")
		return
	}
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

	var templateMetadata *createTemplateMetadataDto
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
	newState, err := convertSourceModelFromEnvironmentDto(*envDto, &currencyCode, state.OwnerId.ValueStringPointer(), templateMetadata, templates, state.Timeouts, *r.EnvironmentClient.Api.Config)

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

	envProp := EnviromentPropertiesDto{
		DisplayName:     plan.DisplayName.ValueString(),
		EnvironmentSku:  plan.EnvironmentType.ValueString(),
		BingChatEnabled: plan.AllowBingSearch.ValueBool(),
	}

	environmentDto := EnvironmentDto{
		Id:         plan.Id.ValueString(),
		Name:       plan.DisplayName.ValueString(),
		Type:       plan.EnvironmentType.ValueString(),
		Location:   plan.Location.ValueString(),
		Properties: &envProp,
	}

	err := r.updateEnvironmentType(ctx, plan, state)
	if err != nil {
		resp.Diagnostics.AddError("Error when updating environment type", err.Error())
		return
	}
	updateDescription(plan, &environmentDto)
	updateCadence(plan, &environmentDto)
	err = r.updateAllowBingSearch(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error when updating allow bing search", err.Error())
		return
	}
	updateEnvironmentGroupId(plan, &environmentDto)
	updateBillingPolicyId(plan, &environmentDto)

	currencyCode, err := r.updateDataverse(ctx, plan, state, &environmentDto)
	if err != nil {
		resp.Diagnostics.AddError("Error when updating dataverse", err.Error())
		return
	}

	err = r.removeBillingPolicy(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("Error when removing billing policy", err.Error())
	}
	err = r.addBillingPolicy(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error when adding billing policy", err.Error())
	}

	envDto, err := r.EnvironmentClient.UpdateEnvironment(ctx, plan.Id.ValueString(), environmentDto)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.ProviderTypeName), err.Error())
		return
	}

	var templateMetadata *createTemplateMetadataDto
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

	newState, err := convertSourceModelFromEnvironmentDto(*envDto, &currencyCode, state.OwnerId.ValueStringPointer(), templateMetadata, templates, plan.Timeouts, *r.EnvironmentClient.Api.Config)
	if err != nil {
		resp.Diagnostics.AddError("Error when converting environment to source model", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *Resource) updateEnvironmentAiFeatures(ctx context.Context, environmentId string, allowBingSearch bool, allowMovingData *bool) error {
	featuresDto := GenerativeAiFeaturesDto{
		Properties: GenerativeAiFeaturesPropertiesDto{
			BingChatEnabled: allowBingSearch,
		},
	}

	if allowMovingData != nil {
		featuresDto.Properties.CopilotPolicies = &CopilotPoliciesDto{
			CrossGeoCopilotDataMovementEnabled: allowMovingData,
		}
	}

	err := r.EnvironmentClient.UpdateEnvironmentAiFeatures(ctx, environmentId, featuresDto)
	if err != nil {
		return err
	}

	return nil
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
		DomainName: dataverseSourcePlanModel.Domain.ValueString(),
	}

	if plan.EnvironmentType.ValueString() != EnvironmentTypesDeveloper {
		environmentDto.Properties.LinkedEnvironmentMetadata.SecurityGroupId = types.StringNull().ValueString()
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

func (r *Resource) updateEnvironmentType(ctx context.Context, plan *SourceModel, state *SourceModel) error {
	if plan.EnvironmentType.ValueString() != state.EnvironmentType.ValueString() {
		err := r.EnvironmentClient.ModifyEnvironmentType(ctx, plan.Id.ValueString(), plan.EnvironmentType.ValueString())
		if err != nil {
			return fmt.Errorf("Error when updating environment_type: %s", err.Error())
		}
	}
	return nil
}

func updateDescription(plan *SourceModel, environmentDto *EnvironmentDto) {
	if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
		environmentDto.Properties.Description = plan.Description.ValueString()
	}
}

func updateCadence(plan *SourceModel, environmentDto *EnvironmentDto) {
	if !plan.Cadence.IsNull() && plan.Cadence.ValueString() != "" {
		environmentDto.Properties.UpdateCadence = &UpdateCadenceDto{
			Id: plan.Cadence.ValueString(),
		}
	}
}

func (r *Resource) updateAllowBingSearch(ctx context.Context, plan *SourceModel) error {
	if !plan.AllowBingSearch.IsNull() && !plan.AllowBingSearch.IsUnknown() {
		err := r.updateEnvironmentAiFeatures(ctx, plan.Id.ValueString(), plan.AllowBingSearch.ValueBool(), plan.AllowMovingDataAcrossRegions.ValueBoolPointer())
		if err != nil {
			return err
		}
	}
	return nil
}

func updateEnvironmentGroupId(plan *SourceModel, environmentDto *EnvironmentDto) {
	if !plan.EnvironmentGroupId.IsNull() && !plan.EnvironmentGroupId.IsUnknown() {
		envGroupId := constants.ZERO_UUID
		if plan.EnvironmentGroupId.ValueString() != "" && plan.EnvironmentGroupId.ValueString() != constants.ZERO_UUID {
			envGroupId = plan.EnvironmentGroupId.ValueString()
		}
		environmentDto.Properties.ParentEnvironmentGroup = &ParentEnvironmentGroupDto{
			Id: envGroupId,
		}
	}
}

func updateBillingPolicyId(plan *SourceModel, environmentDto *EnvironmentDto) {
	if !plan.BillingPolicyId.IsNull() && plan.BillingPolicyId.ValueString() != "" {
		environmentDto.Properties.BillingPolicy = &BillingPolicyDto{
			Id: plan.BillingPolicyId.ValueString(),
		}
	}
}

func (r *Resource) updateDataverse(ctx context.Context, plan *SourceModel, state *SourceModel, environmentDto *EnvironmentDto) (string, error) {
	var currencyCode string
	if !isDataverseEnvironmentEmpty(ctx, state) && !isDataverseEnvironmentEmpty(ctx, plan) {
		currencyCode = updateExistingDataverse(ctx, plan, environmentDto, state)
	} else if isDataverseEnvironmentEmpty(ctx, state) && !isDataverseEnvironmentEmpty(ctx, plan) {
		code, err := addDataverse(ctx, plan, r)
		if err != nil {
			return "", err
		}
		currencyCode = code
	}
	return currencyCode, nil
}

func (r *Resource) removeBillingPolicy(ctx context.Context, state *SourceModel) error {
	if !state.BillingPolicyId.IsNull() && !state.BillingPolicyId.IsUnknown() && state.BillingPolicyId.ValueString() != "" {
		tflog.Debug(ctx, fmt.Sprintf("Removing environment %s from billing policy %s", state.Id.ValueString(), state.BillingPolicyId.ValueString()))
		err := r.LicensingClient.RemoveEnvironmentsToBillingPolicy(ctx, state.BillingPolicyId.ValueString(), []string{state.Id.ValueString()})
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Resource) addBillingPolicy(ctx context.Context, plan *SourceModel) error {
	if !plan.BillingPolicyId.IsNull() && !plan.BillingPolicyId.IsUnknown() && plan.BillingPolicyId.ValueString() != "" {
		tflog.Debug(ctx, fmt.Sprintf("Adding environment %s to billing policy %s", plan.Id.ValueString(), plan.BillingPolicyId.ValueString()))
		err := r.LicensingClient.AddEnvironmentsToBillingPolicy(ctx, plan.BillingPolicyId.ValueString(), []string{plan.Id.ValueString()})
		if err != nil {
			return err
		}
	}
	return nil
}
