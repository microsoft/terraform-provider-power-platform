// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant_settings

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customtypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ resource.Resource = &TenantSettingsResource{}
var _ resource.ResourceWithImportState = &TenantSettingsResource{}
var _ resource.ResourceWithModifyPlan = &TenantSettingsResource{}

func NewTenantSettingsResource() resource.Resource {
	return &TenantSettingsResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "tenant_settings",
		},
	}
}

func (r *TenantSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *TenantSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Power Platform Tenant Settings. Power Platform Tenant Settings are configuration options that apply to the entire tenant. They control various aspects of Power Platform features and behaviors, such as security, data protection, licensing, and more. These settings apply to all environments within your tenant. See [Tenant Settings Overview](https://learn.microsoft.com/power-platform/admin/tenant-settings) for more details.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Read:   true,
				Delete: true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the Power Platform Tenant",
				Computed:            true, Required: false, Optional: false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"walk_me_opt_out": schema.BoolAttribute{
				MarkdownDescription: "Walk Me Opt Out",
				Optional:            true,
			},
			"disable_newsletter_sendout": schema.BoolAttribute{
				MarkdownDescription: "Disable Newsletter Sendout",
				Optional:            true,
			},
			"disable_environment_creation_by_non_admin_users": schema.BoolAttribute{
				MarkdownDescription: "Disable Environment Creation By Non Admin Users. See [Control environment creation](https://learn.microsoft.com/power-platform/admin/control-environment-creation) for more details.",
				Optional:            true,
			},
			"disable_portals_creation_by_non_admin_users": schema.BoolAttribute{
				MarkdownDescription: "Disable Portals Creation By Non Admin Users",
				Optional:            true,
			},
			"disable_trial_environment_creation_by_non_admin_users": schema.BoolAttribute{
				MarkdownDescription: "Disable Trial Environment Creation By Non Admin Users. See [Control environment creation](https://learn.microsoft.com/power-platform/admin/control-environment-creation) for more details.",
				Optional:            true,
			},
			"disable_capacity_allocation_by_environment_admins": schema.BoolAttribute{
				MarkdownDescription: "Disable Capacity Allocation By Environment Admins. See [Add-on capacity management](https://learn.microsoft.com/power-platform/admin/capacity-add-on#control-who-can-allocate-add-on-capacity) for more details.",
				Optional:            true,
			},
			"disable_support_tickets_visible_by_all_users": schema.BoolAttribute{
				MarkdownDescription: "Disable Support Tickets Visible By All Users",
				Optional:            true,
			},
			"power_platform": schema.SingleNestedAttribute{
				MarkdownDescription: "Power Platform",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"product_feedback": schema.SingleNestedAttribute{
						MarkdownDescription: "Product Feedback",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"disable_microsoft_surveys_send": schema.BoolAttribute{
								MarkdownDescription: "Disable letting Microsoft send surveys",
								Optional:            true,
							},
							"disable_user_survey_feedback": schema.BoolAttribute{
								MarkdownDescription: "Disable users to choose to provide survey feedback",
								Optional:            true,
							},
							"disable_attachments": schema.BoolAttribute{
								MarkdownDescription: "Disable screenshots and attachments in feedback",
								Optional:            true,
							},
							"disable_microsoft_follow_up": schema.BoolAttribute{
								MarkdownDescription: "Disable letting Microsoft follow up on feedback",
								Optional:            true,
							},
						},
					},
					"search": schema.SingleNestedAttribute{
						MarkdownDescription: "Search",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"disable_docs_search": schema.BoolAttribute{
								MarkdownDescription: "Disable Docs Search",
								Optional:            true,
							},
							"disable_community_search": schema.BoolAttribute{
								MarkdownDescription: "Disable Community Search",
								Optional:            true,
							},
							"disable_bing_video_search": schema.BoolAttribute{
								MarkdownDescription: "Disable Bing Video Search",
								Optional:            true,
							},
						},
					},
					"teams_integration": schema.SingleNestedAttribute{
						Description: "Teams Integration",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"share_with_colleagues_user_limit": schema.Int64Attribute{
								MarkdownDescription: "Share With Colleagues User Limit",
								Optional:            true,
							},
						},
					},
					"power_apps": schema.SingleNestedAttribute{
						Description:   "Power Apps",
						Optional:      true,
						PlanModifiers: []planmodifier.Object{
							// objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"disable_share_with_everyone": schema.BoolAttribute{
								MarkdownDescription: "Disable Share With Everyone",
								Optional:            true,
							},
							"enable_guests_to_make": schema.BoolAttribute{
								MarkdownDescription: "Enable Guests To Make",
								Optional:            true,
							},
							"disable_maker_match": schema.BoolAttribute{
								MarkdownDescription: "Disable Maker Match",
								Optional:            true,
							},
							"disable_unused_license_assignment": schema.BoolAttribute{
								MarkdownDescription: "Disable Unused License Assignment",
								Optional:            true,
							},
							"disable_create_from_image": schema.BoolAttribute{
								DeprecationMessage:  "[DEPRECATED] This attribute is deprecated and will be removed in a future release.",
								MarkdownDescription: "[DEPRECATED] Disable Create From Image",
								Optional:            true,
							},
							"disable_create_from_figma": schema.BoolAttribute{
								DeprecationMessage:  "[DEPRECATED] This attribute is deprecated and will be removed in a future release.",
								MarkdownDescription: "[DEPRECATED] Disable Create From Figma",
								Optional:            true,
							},
							"disable_connection_sharing_with_everyone": schema.BoolAttribute{
								MarkdownDescription: "Disable Connection Sharing With Everyone",
								Optional:            true,
							},
						},
					},
					"power_automate": schema.SingleNestedAttribute{
						Description: "Power Automate",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"disable_copilot": schema.BoolAttribute{
								MarkdownDescription: "Disable Copilot",
								Optional:            true,
							},
						},
					},
					"environments": schema.SingleNestedAttribute{
						Description: "Environments",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"disable_preferred_data_location_for_teams_environment": schema.BoolAttribute{
								MarkdownDescription: "Disable Preferred Data Location For Teams Environment",
								Optional:            true,
							},
						},
					},
					"governance": schema.SingleNestedAttribute{
						MarkdownDescription: "Governance",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"disable_admin_digest": schema.BoolAttribute{
								MarkdownDescription: "Disable Admin Digest",
								Optional:            true,
							},
							"disable_developer_environment_creation_by_non_admin_users": schema.BoolAttribute{
								MarkdownDescription: "Disable Developer Environment Creation By Non Admin Users",
								Optional:            true,
							},
							"enable_default_environment_routing": schema.BoolAttribute{
								MarkdownDescription: "Enable Default Environment Routing",
								Optional:            true,
							},
							"environment_routing_all_makers": schema.BoolAttribute{
								MarkdownDescription: "Select who can be routed to a new personal developer environment. (All Makers = true, New Makers = false)",
								Optional:            true,
							},
							"environment_routing_target_environment_group_id": schema.StringAttribute{
								MarkdownDescription: "Assign newly created personal developer environments to a specific environment group",
								Optional:            true,
								CustomType:          customtypes.UUIDType{},
							},
							"environment_routing_target_security_group_id": schema.StringAttribute{
								MarkdownDescription: "Restrict routing to members of the following security group. (00000000-0000-0000-0000-000000000000 allows all users)",
								Optional:            true,
								CustomType:          customtypes.UUIDType{},
							},
							"policy": schema.SingleNestedAttribute{
								MarkdownDescription: "Policy",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"enable_desktop_flow_data_policy_management": schema.BoolAttribute{
										MarkdownDescription: "Enable Desktop Flow Data Policy Management",
										Optional:            true,
									},
								},
							},
						},
					},
					"licensing": schema.SingleNestedAttribute{
						MarkdownDescription: "Licensing",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"disable_billing_policy_creation_by_non_admin_users": schema.BoolAttribute{
								MarkdownDescription: "Disable Billing Policy Creation By Non Admin Users",
								Optional:            true,
							},
							"enable_tenant_capacity_report_for_environment_admins": schema.BoolAttribute{
								MarkdownDescription: "Enable Tenant Capacity Report For Environment Admins",
								Optional:            true,
							},
							"storage_capacity_consumption_warning_threshold": schema.Int64Attribute{
								MarkdownDescription: "Storage Capacity Consumption Warning Threshold",
								Optional:            true,
							},
							"enable_tenant_licensing_report_for_environment_admins": schema.BoolAttribute{
								MarkdownDescription: "Enable Tenant Licensing Report For Environment Admins",
								Optional:            true,
							},
							"disable_use_of_unassigned_ai_builder_credits": schema.BoolAttribute{
								MarkdownDescription: "Disable Use Of Unassigned AI Builder Credits",
								Optional:            true,
							},
						},
					},
					"power_pages": schema.SingleNestedAttribute{
						MarkdownDescription: "Power Pages",
						Optional:            true,
						Attributes:          map[string]schema.Attribute{},
					},
					"champions": schema.SingleNestedAttribute{
						MarkdownDescription: "Champions",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"disable_champions_invitation_reachout": schema.BoolAttribute{
								MarkdownDescription: "Disable Champions Invitation Reachout",
								Optional:            true,
							},
							"disable_skills_match_invitation_reachout": schema.BoolAttribute{
								MarkdownDescription: "Disable Skills Match Invitation Reachout",
								Optional:            true,
							},
						},
					},
					"intelligence": schema.SingleNestedAttribute{
						MarkdownDescription: "Intelligence",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"disable_copilot": schema.BoolAttribute{
								MarkdownDescription: "Disable Copilot",
								Optional:            true,
							},
							"enable_open_ai_bot_publishing": schema.BoolAttribute{
								MarkdownDescription: "Enable Open AI Bot Publishing",
								Optional:            true,
							},
						},
					},
					"model_experimentation": schema.SingleNestedAttribute{
						MarkdownDescription: "Model Experimentation",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"enable_model_data_sharing": schema.BoolAttribute{
								MarkdownDescription: "Enable Model Data Sharing",
								Optional:            true,
							},
							"disable_data_logging": schema.BoolAttribute{
								MarkdownDescription: "Disable Data Logging",
								Optional:            true,
							},
						},
					},
					"catalog_settings": schema.SingleNestedAttribute{
						MarkdownDescription: "Catalog Settings",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"power_catalog_audience_setting": schema.StringAttribute{
								MarkdownDescription: "Power Catalog Audience Setting",
								Optional:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("SpecificAdmins", "All"),
								},
							},
						},
					},
					"user_management_settings": schema.SingleNestedAttribute{
						MarkdownDescription: "User Management Settings",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"enable_delete_disabled_user_in_all_environments": schema.BoolAttribute{
								MarkdownDescription: "Enable Delete Disabled User In All Environments",
								Optional:            true,
							},
						},
					},
				},
			},
		},
	}
}

func (r *TenantSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
		return
	}

	providerClient, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.TenantSettingClient = newTenantSettingsClient(providerClient.Api)
}

func (r *TenantSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var plan TenantSettingsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save the original tenant settings in private state
	originalSettings, erro := r.TenantSettingClient.GetTenantSettings(ctx)
	if erro != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tenant Settings in Create",
			fmt.Sprintf("Could not read existing tenant settings during resource creation: %s", erro.Error()),
		)
		return
	}

	jsonSettings, errj := json.Marshal(originalSettings)
	if errj != nil {
		resp.Diagnostics.AddError(
			"Unable to Marshal Tenant Settings in Create",
			fmt.Sprintf("Could not marshal original tenant settings to JSON for backup: %s", errj.Error()),
		)
		return
	}
	resp.Private.SetKey(ctx, "original_settings", jsonSettings)

	// Get the plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update tenant settings via the API
	plannedSettingsDto, err := convertFromTenantSettingsModel(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Convert Tenant Settings Model to DTO in Create", err.Error())
		return
	}

	tenantSettingsDto, err := r.TenantSettingClient.UpdateTenantSettings(ctx, plannedSettingsDto)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Tenant Settings",
			fmt.Sprintf("Could not create tenant settings with API: %s", err.Error()),
		)
		return
	}

	stateDto, err := applyCorrections(ctx, plannedSettingsDto, *tenantSettingsDto)
	if err != nil {
		resp.Diagnostics.AddError("Error applying corrections to tenant settings", err.Error())
		return
	}

	state, _, err := convertFromTenantSettingsDto[TenantSettingsResourceModel](*stateDto, plan.Timeouts)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Convert Tenant Settings DTO to Model in Create", err.Error())
		return
	}
	state.Id = plan.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID %s", plan.Id.ValueString()))
}

func (r *TenantSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var state TenantSettingsResourceModel

	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantSettings, err := r.TenantSettingClient.GetTenantSettings(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tenant Settings in Read",
			fmt.Sprintf("Could not read current tenant settings during resource read: %s", err.Error()),
		)
		return
	}

	tenant, errt := r.TenantSettingClient.GetTenant(ctx)
	if errt != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tenant Information in Read",
			fmt.Sprintf("Could not read tenant information during resource read: %s", errt.Error()),
		)
		return
	}

	var configuredSettings TenantSettingsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &configuredSettings)...)
	oldStateDto, err := convertFromTenantSettingsModel(ctx, configuredSettings)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Convert Tenant Settings Model to DTO in Read", err.Error())
		return
	}
	newStateDto, err := applyCorrections(ctx, oldStateDto, *tenantSettings)
	if err != nil {
		resp.Diagnostics.AddError("Error applying corrections to tenant settings", err.Error())
		return
	}
	newState, _, err := convertFromTenantSettingsDto[TenantSettingsResourceModel](*newStateDto, state.Timeouts)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Convert Tenant Settings DTO to Model in Read", err.Error())
		return
	}
	newState.Id = types.StringValue(tenant.TenantId)

	tflog.Debug(ctx, fmt.Sprintf("READ: %s with id %s", r.FullTypeName(), newState.Id.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *TenantSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan TenantSettingsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plannedDto, err := convertFromTenantSettingsModel(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Convert Tenant Settings Model to DTO in Update", err.Error())
		return
	}
	// Preprocessing updates is unfortunately needed because Terraform can not treat a zeroed UUID as a null value.
	// This captures the case where a UUID is changed from known to zeroed/null.  Zeroed UUIDs come back as null from the API.
	// The plannedDto remembers what the user intended, and the preprocessedDto is what we will send to the API.
	preprocessedDto, err := convertFromTenantSettingsModel(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Convert Tenant Settings Model to DTO for Preprocessing in Update", err.Error())
		return
	}

	needsProcessing := func(p path.Path) bool {
		var attrPlan customtypes.UUID
		var attrState customtypes.UUID

		diag := req.State.GetAttribute(ctx, p, &attrState)
		diag2 := req.Plan.GetAttribute(ctx, p, &attrPlan)
		if !diag.HasError() && !diag2.HasError() {
			return !attrState.IsNull() && attrPlan.IsNull()
		}

		return false
	}

	if needsProcessing(path.Root("power_platform").AtName("governance").AtName("environment_routing_target_security_group_id")) {
		preprocessedDto.PowerPlatform.Governance.EnvironmentRoutingTargetSecurityGroupId = types.StringValue(constants.ZERO_UUID).ValueStringPointer()
	}

	if needsProcessing(path.Root("power_platform").AtName("governance").AtName("environment_routing_target_environment_group_id")) {
		preprocessedDto.PowerPlatform.Governance.EnvironmentRoutingTargetEnvironmentGroupId = types.StringValue(constants.ZERO_UUID).ValueStringPointer()
	}

	// send preprocessedDto to the API
	updatedSettingsDto, err := r.TenantSettingClient.UpdateTenantSettings(ctx, preprocessedDto)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update Tenant Settings",
			fmt.Sprintf("Could not update tenant settings with API: %s", err.Error()),
		)
		return
	}

	// need to make corrections from what the API returns to match what terraform expects
	filteredDto, err := applyCorrections(ctx, plannedDto, *updatedSettingsDto)
	if err != nil {
		resp.Diagnostics.AddError("Error applying corrections to tenant settings", err.Error())
		return
	}

	newState, _, err := convertFromTenantSettingsDto[TenantSettingsResourceModel](*filteredDto, plan.Timeouts)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Convert Tenant Settings DTO to Model in Update", err.Error())
		return
	}
	newState.Id = plan.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *TenantSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TenantSettingsResourceModel
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddWarning("Tenant Settings Cannot Be Deleted", "Tenant Settings cannot be permanently deleted in Power Platform. Deleting this resource will attempt to restore settings to their previous values and remove this configuration from Terraform state.")

	stateDto, err := convertFromTenantSettingsModel(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Convert Tenant Settings Model to DTO in Delete", err.Error())
		return
	}

	// restore to previous state
	previousBytes, diags := req.Private.GetKey(ctx, "original_settings")
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	var originalSettings tenantSettingsDto
	err2 := json.Unmarshal(previousBytes, &originalSettings)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Unable to Unmarshal Original Settings in Delete",
			fmt.Sprintf("Could not unmarshal backup of original tenant settings during resource deletion: %s", err2.Error()),
		)
		return
	}

	correctedDto, err := applyCorrections(ctx, stateDto, originalSettings)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error applying corrections", fmt.Sprintf("Error applying corrections: %s", err.Error()),
		)
		return
	}

	_, e := r.TenantSettingClient.UpdateTenantSettings(ctx, *correctedDto)
	if e != nil {
		resp.Diagnostics.AddError(
			"Unable to Restore Tenant Settings in Delete",
			fmt.Sprintf("Could not restore original tenant settings during resource deletion: %s", e.Error()),
		)
		return
	}
}

func (r *TenantSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *TenantSettingsResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan TenantSettingsResourceModel
	if !req.Plan.Raw.IsNull() {
		// this is create
		req.Plan.Get(ctx, &plan)
		if plan.Id.IsUnknown() || plan.Id.IsNull() {
			tenant, errt := r.TenantSettingClient.GetTenant(ctx)
			if errt != nil {
				resp.Diagnostics.AddError(
					"Unable to Read Tenant Information in ModifyPlan",
					fmt.Sprintf("Could not read tenant information during plan modification: %s", errt.Error()),
				)
				return
			}
			plan.Id = types.StringValue(tenant.TenantId)
			resp.Plan.Set(ctx, &plan)
		}
	}
}
