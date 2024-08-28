// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/customtypes"
	modifiers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/modifiers"
)

var _ resource.Resource = &TenantSettingsResource{}
var _ resource.ResourceWithImportState = &TenantSettingsResource{}
var _ resource.ResourceWithModifyPlan = &TenantSettingsResource{}

const ZERO_UUID = "00000000-0000-0000-0000-000000000000"

func NewTenantSettingsResource() resource.Resource {
	return &TenantSettingsResource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_tenant_settings",
	}
}

type TenantSettingsResource struct {
	TenantSettingClient TenantSettingsClient
	ProviderTypeName    string
	TypeName            string
}

func (r *TenantSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

func (r *TenantSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Power Platform Tenant Settings.",
		MarkdownDescription: "Manages Power Platform Tenant Settings. Power Platform Tenant Settings are configuration options that apply to the entire tenant. They control various aspects of Power Platform features and behaviors, such as security, data protection, licensing, and more. These settings apply to all environments within your tenant. See [Tenant Settings Overview](https://learn.microsoft.com/power-platform/admin/tenant-settings) for more details.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Tenant ID",
				MarkdownDescription: "Id of the Power Platform Tenant",
				Computed:            true, Required: false, Optional: false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"walk_me_opt_out": schema.BoolAttribute{
				Description: "Walk Me Opt Out",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					// boolplanmodifier.UseStateForUnknown(),
					modifiers.RestoreOriginalBoolModifier(),
				},
			},
			"disable_nps_comments_reachout": schema.BoolAttribute{
				Description:   "Disable NPS Comments Reachout",
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{
					// boolplanmodifier.UseStateForUnknown(),
				},
			},
			"disable_newsletter_sendout": schema.BoolAttribute{
				Description:   "Disable Newsletter Sendout",
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{
					// boolplanmodifier.UseStateForUnknown(),
				},
			},
			"disable_environment_creation_by_non_admin_users": schema.BoolAttribute{
				Description:         "Disable Environment Creation By Non Admin Users",
				MarkdownDescription: "Disable Environment Creation By Non Admin Users. See [Control environment creation](https://learn.microsoft.com/power-platform/admin/control-environment-creation) for more details.",
				Optional:            true,
				PlanModifiers:       []planmodifier.Bool{
					// boolplanmodifier.UseStateForUnknown(),
				},
			},
			"disable_portals_creation_by_non_admin_users": schema.BoolAttribute{
				Description:   "Disable Portals Creation By Non Admin Users",
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{
					// boolplanmodifier.UseStateForUnknown(),
				},
			},
			"disable_survey_feedback": schema.BoolAttribute{
				Description:   "Disable Survey Feedback",
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{
					// boolplanmodifier.UseStateForUnknown(),
				},
			},
			"disable_trial_environment_creation_by_non_admin_users": schema.BoolAttribute{
				Description:         "Disable Trial Environment Creation By Non Admin Users",
				MarkdownDescription: "Disable Trial Environment Creation By Non Admin Users. See [Control environment creation](https://learn.microsoft.com/power-platform/admin/control-environment-creation) for more details.",
				Optional:            true,
				PlanModifiers:       []planmodifier.Bool{
					// boolplanmodifier.UseStateForUnknown(),
				},
			},
			"disable_capacity_allocation_by_environment_admins": schema.BoolAttribute{
				Description:         "Disable Capacity Allocation By Environment Admins",
				MarkdownDescription: "Disable Capacity Allocation By Environment Admins. See [Add-on capacity management](https://learn.microsoft.com/power-platform/admin/capacity-add-on#control-who-can-allocate-add-on-capacity) for more details.",
				Optional:            true,
				PlanModifiers:       []planmodifier.Bool{
					// boolplanmodifier.UseStateForUnknown(),
				},
			},
			"disable_support_tickets_visible_by_all_users": schema.BoolAttribute{
				Description: "Disable Support Tickets Visible By All Users",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"power_platform": schema.SingleNestedAttribute{
				Description:   "Power Platform",
				Optional:      true,
				PlanModifiers: []planmodifier.Object{
					// objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"search": schema.SingleNestedAttribute{
						Description:   "Search",
						Optional:      true,
						PlanModifiers: []planmodifier.Object{
							// objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"disable_docs_search": schema.BoolAttribute{
								Description:   "Disable Docs Search",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
							"disable_community_search": schema.BoolAttribute{
								Description:   "Disable Community Search",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
							"disable_bing_video_search": schema.BoolAttribute{
								Description:   "Disable Bing Video Search",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"teams_integration": schema.SingleNestedAttribute{
						Description:   "Teams Integration",
						Optional:      true,
						PlanModifiers: []planmodifier.Object{
							// objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"share_with_colleagues_user_limit": schema.Int64Attribute{
								Description:   "Share With Colleagues User Limit",
								Optional:      true,
								PlanModifiers: []planmodifier.Int64{
									// int64planmodifier.UseStateForUnknown(),
								},
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
								Description:   "Disable Share With Everyone",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_guests_to_make": schema.BoolAttribute{
								Description:   "Enable Guests To Make",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
							"disable_maker_match": schema.BoolAttribute{
								Description:   "Disable Maker Match",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
							"disable_unused_license_assignment": schema.BoolAttribute{
								Description:   "Disable Unused License Assignment",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
							"disable_create_from_image": schema.BoolAttribute{
								Description:   "Disable Create From Image",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
							"disable_create_from_figma": schema.BoolAttribute{
								Description:   "Disable Create From Figma",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
							"disable_connection_sharing_with_everyone": schema.BoolAttribute{
								Description:   "Disable Connection Sharing With Everyone",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"power_automate": schema.SingleNestedAttribute{
						Description:   "Power Automate",
						Optional:      true,
						PlanModifiers: []planmodifier.Object{
							// objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"disable_copilot": schema.BoolAttribute{
								Description:   "Disable Copilot",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"environments": schema.SingleNestedAttribute{
						Description:   "Environments",
						Optional:      true,
						PlanModifiers: []planmodifier.Object{
							// objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"disable_preferred_data_location_for_teams_environment": schema.BoolAttribute{
								Description:   "Disable Preferred Data Location For Teams Environment",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"governance": schema.SingleNestedAttribute{
						Description:   "Governance",
						Optional:      true,
						PlanModifiers: []planmodifier.Object{
							// objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"disable_admin_digest": schema.BoolAttribute{
								Description:   "Disable Admin Digest",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
							"disable_developer_environment_creation_by_non_admin_users": schema.BoolAttribute{
								Description:   "Disable Developer Environment Creation By Non Admin Users",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_default_environment_routing": schema.BoolAttribute{
								Description: "Enable Default Environment Routing",
								Optional:    true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
									modifiers.RestoreOriginalBoolModifier(),
								},
							},
							"environment_routing_all_makers": schema.BoolAttribute{
								Description: "Select who can be routed to a new personal developer environment. (All Makers = true, New Makers = false)",
								Optional:    true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
									modifiers.RestoreOriginalBoolModifier(),
								},
							},
							"environment_routing_target_environment_group_id": schema.StringAttribute{
								Description:   "Assign newly created personal developer environments to a specific environment group",
								Optional:      true,
								CustomType:    customtypes.UUIDType{},
								PlanModifiers: []planmodifier.String{},
							},
							"environment_routing_target_security_group_id": schema.StringAttribute{
								Description:   "Restrict routing to members of the following security group. (00000000-0000-0000-0000-000000000000 allows all users)",
								Optional:      true,
								CustomType:    customtypes.UUIDType{},
								PlanModifiers: []planmodifier.String{},
							},
							"policy": schema.SingleNestedAttribute{
								Description: "Policy",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"enable_desktop_flow_data_policy_management": schema.BoolAttribute{
										Description: "Enable Desktop Flow Data Policy Management",
										Optional:    true,
									},
								},
							},
						},
					},
					"licensing": schema.SingleNestedAttribute{
						Description:   "Licensing",
						Optional:      true,
						PlanModifiers: []planmodifier.Object{
							// objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"disable_billing_policy_creation_by_non_admin_users": schema.BoolAttribute{
								Description:   "Disable Billing Policy Creation By Non Admin Users",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_tenant_capacity_report_for_environment_admins": schema.BoolAttribute{
								Description:   "Enable Tenant Capacity Report For Environment Admins",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
							"storage_capacity_consumption_warning_threshold": schema.Int64Attribute{
								Description:   "Storage Capacity Consumption Warning Threshold",
								Optional:      true,
								PlanModifiers: []planmodifier.Int64{
									// int64planmodifier.UseStateForUnknown(),
								},
							},
							"enable_tenant_licensing_report_for_environment_admins": schema.BoolAttribute{
								Description:   "Enable Tenant Licensing Report For Environment Admins",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
							"disable_use_of_unassigned_ai_builder_credits": schema.BoolAttribute{
								Description:   "Disable Use Of Unassigned AI Builder Credits",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"power_pages": schema.SingleNestedAttribute{
						Description:   "Power Pages",
						Optional:      true,
						PlanModifiers: []planmodifier.Object{
							// objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{},
					},
					"champions": schema.SingleNestedAttribute{
						Description:   "Champions",
						Optional:      true,
						PlanModifiers: []planmodifier.Object{
							// objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"disable_champions_invitation_reachout": schema.BoolAttribute{
								Description:   "Disable Champions Invitation Reachout",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
							"disable_skills_match_invitation_reachout": schema.BoolAttribute{
								Description:   "Disable Skills Match Invitation Reachout",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"intelligence": schema.SingleNestedAttribute{
						Description:   "Intelligence",
						Optional:      true,
						PlanModifiers: []planmodifier.Object{
							// objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"disable_copilot": schema.BoolAttribute{
								Description:   "Disable Copilot",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_open_ai_bot_publishing": schema.BoolAttribute{
								Description:   "Enable Open AI Bot Publishing",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"model_experimentation": schema.SingleNestedAttribute{
						Description:   "Model Experimentation",
						Optional:      true,
						PlanModifiers: []planmodifier.Object{
							// objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"enable_model_data_sharing": schema.BoolAttribute{
								Description:   "Enable Model Data Sharing",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
							"disable_data_logging": schema.BoolAttribute{
								Description:   "Disable Data Logging",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"catalog_settings": schema.SingleNestedAttribute{
						Description:   "Catalog Settings",
						Optional:      true,
						PlanModifiers: []planmodifier.Object{
							// objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"power_catalog_audience_setting": schema.StringAttribute{
								Description:   "Power Catalog Audience Setting",
								Optional:      true,
								PlanModifiers: []planmodifier.String{
									// stringplanmodifier.UseStateForUnknown(),
								},
								Validators: []validator.String{
									stringvalidator.OneOf("SpecificAdmin", "All"),
								},
							},
						},
					},
					"user_management_settings": schema.SingleNestedAttribute{
						Description:   "User Management Settings",
						Optional:      true,
						PlanModifiers: []planmodifier.Object{
							// objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"enable_delete_disabled_user_in_all_environments": schema.BoolAttribute{
								Description:   "Enable Delete Disabled User In All Environments",
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{
									// boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *TenantSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := req.ProviderData.(*api.ProviderClient).Api

	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.TenantSettingClient = NewTenantSettingsClient(client)
}

func (r *TenantSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TenantSettingsSourceModel

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE START: %s", r.ProviderTypeName))

	// Save the original tenant settings in private state
	originalSettings, erro := r.TenantSettingClient.GetTenantSettings(ctx)
	if erro != nil {
		resp.Diagnostics.AddError(
			"Error reading tenant settings", fmt.Sprintf("Error reading tenant settings: %s", erro.Error()),
		)
		return
	}

	jsonSettings, errj := json.Marshal(originalSettings)
	if errj != nil {
		resp.Diagnostics.AddError(
			"Error marshalling tenant settings", fmt.Sprintf("Error marshalling tenant settings: %s", errj.Error()),
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
	plannedSettingsDto := ConvertFromTenantSettingsModel(ctx, plan)
	tenantSettingsDto, err := r.TenantSettingClient.UpdateTenantSettings(ctx, plannedSettingsDto)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating tenant settings", fmt.Sprintf("Error creating tenant settings: %s", err.Error()),
		)
		return
	}

	stateDto := applyCorrections(plannedSettingsDto, *tenantSettingsDto)

	state, _ := ConvertFromTenantSettingsDto(*stateDto)
	state.Id = plan.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID %s", plan.Id.ValueString()))
	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *TenantSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", r.ProviderTypeName))

	if resp.Diagnostics.HasError() {
		return
	}

	tenantSettings, err := r.TenantSettingClient.GetTenantSettings(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading tenant settings", fmt.Sprintf("Error reading tenant settings: %s", err.Error()),
		)
		return
	}

	tenant, errt := r.TenantSettingClient.GetTenant(ctx)
	if errt != nil {
		resp.Diagnostics.AddError(
			"Error reading tenant", fmt.Sprintf("Error reading tenant: %s", errt.Error()),
		)
		return
	}

	var configuredSettings TenantSettingsSourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &configuredSettings)...)
	oldStateDto := ConvertFromTenantSettingsModel(ctx, configuredSettings)
	newStateDto := applyCorrections(oldStateDto, *tenantSettings)
	newState, _ := ConvertFromTenantSettingsDto(*newStateDto)
	newState.Id = types.StringValue(tenant.TenantId)

	tflog.Debug(ctx, fmt.Sprintf("READ: %s_tenant_settings with id %s", r.ProviderTypeName, newState.Id.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE END: %s", r.ProviderTypeName))
}

func (r *TenantSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TenantSettingsSourceModel

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plannedDto := ConvertFromTenantSettingsModel(ctx, plan)

	// Preprocessing updates is unfortunately needed because Terraform can not treat a zeroed UUID as a null value.
	// This captures the case where a UUID is changed from known to zeroed/null.  Zeroed UUIDs come back as null from the API.
	// The plannedDto remembers what the user intended, and the preprocessedDto is what we will send to the API.
	preprocessedDto := ConvertFromTenantSettingsModel(ctx, plan)

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
		preprocessedDto.PowerPlatform.Governance.EnvironmentRoutingTargetSecurityGroupId = types.StringValue(ZERO_UUID).ValueStringPointer()
	}

	if needsProcessing(path.Root("power_platform").AtName("governance").AtName("environment_routing_target_environment_group_id")) {
		preprocessedDto.PowerPlatform.Governance.EnvironmentRoutingTargetEnvironmentGroupId = types.StringValue(ZERO_UUID).ValueStringPointer()
	}

	// send preprocessedDto to the API
	updatedSettingsDto, err := r.TenantSettingClient.UpdateTenantSettings(ctx, preprocessedDto)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating tenant settings", fmt.Sprintf("Error updating tenant settings: %s", err.Error()),
		)
		return
	}

	// need to make corrections from what the API returns to match what terraform expects
	filteredDto := applyCorrections(plannedDto, *updatedSettingsDto)

	newState, _ := ConvertFromTenantSettingsDto(*filteredDto)
	newState.Id = plan.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *TenantSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.AddWarning("Tenant Settings are not deleted", "Tenant Settings may not be deleted in Power Platform.  Deleting this resource will attempt to restore settings to their previous values and remove this configuration from Terraform state.")

	var state TenantSettingsSourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateDto := ConvertFromTenantSettingsModel(ctx, state)

	// restore to previous state
	previousBytes, err := req.Private.GetKey(ctx, "original_settings")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading original settings", fmt.Sprintf("Error reading original settings: %s", err.Errors()),
		)
		return
	}

	var originalSettings TenantSettingsDto
	err2 := json.Unmarshal(previousBytes, &originalSettings)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error unmarshalling original settings", fmt.Sprintf("Error unmarshalling original settings: %s", err2.Error()),
		)
		return
	}

	correctedDto := applyCorrections(stateDto, originalSettings)
	if correctedDto == nil {
		resp.Diagnostics.AddError(
			"Error applying corrections", fmt.Sprintf("Error applying corrections"),
		)
		return
	}

	_, e := r.TenantSettingClient.UpdateTenantSettings(ctx, *correctedDto)
	if e != nil {
		resp.Diagnostics.AddError(
			"Error deleting tenant settings", fmt.Sprintf("Error deleting tenant settings: %s", e.Error()),
		)
		return
	}
}

func (r *TenantSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *TenantSettingsResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan TenantSettingsSourceModel
	if !req.Plan.Raw.IsNull() {
		//this is create
		req.Plan.Get(ctx, &plan)
		if plan.Id.IsUnknown() || plan.Id.IsNull() {
			tenant, errt := r.TenantSettingClient.GetTenant(ctx)
			if errt != nil {
				resp.Diagnostics.AddError(
					"Error reading tenant", fmt.Sprintf("Error reading tenant: %s", errt.Error()),
				)
				return
			}
			plan.Id = types.StringValue(tenant.TenantId)
			resp.Plan.Set(ctx, &plan)
		}
	}
}
