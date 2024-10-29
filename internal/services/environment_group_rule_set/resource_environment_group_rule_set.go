// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_set

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/numbervalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/tenant"
)

var _ resource.Resource = &environmentGroupRuleSetResource{}
var _ resource.ResourceWithImportState = &environmentGroupRuleSetResource{}
var _ resource.ResourceWithValidateConfig = &environmentGroupRuleSetResource{}

func NewEnvironmentGroupRuleSetResource() resource.Resource {
	return &environmentGroupRuleSetResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment_group_rule_set",
		},
	}
}

func (r *environmentGroupRuleSetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *environmentGroupRuleSetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	maxSharingRange := []*big.Float{}
	for i := -1; i < 100; i++ {
		maxSharingRange = append(maxSharingRange, big.NewFloat(float64(i)))
	}

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique id of the environment group ruleset",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_group_id": schema.StringAttribute{
				MarkdownDescription: "Unique id of the environment group",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rules": schema.SingleNestedAttribute{
				MarkdownDescription: "Rules for the environment group",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"sharing_controls": schema.SingleNestedAttribute{
						// type: Sharing -> Sharing controls for Canvas apps
						// CanShareWithSecurityGroups: noLimit, excludeSharingToSecurityGroups
						// IsGroupSharingDisabled: true, false
						// MaximumShareLimit: (-1..99)
						//
						// modes:
						// noLimit, false, -1
						// excludeSharingToSecurityGroups, true, (-1.....99)
						MarkdownDescription: "Sharing controls",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"share_mode": schema.StringAttribute{
								MarkdownDescription: "Share mode for canvas apps: `No limit`, `Exclude sharing with security groups`",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("no limit", "exclude sharing with security groups"),
								},
							},
							"share_max_limit": schema.NumberAttribute{
								MarkdownDescription: "Maximum total of individual who can be shared to: (-1..99). If `share_mode` is `No limit`, this value must be -1.",
								Optional:            true,
								PlanModifiers: []planmodifier.Number{
									numberplanmodifier.UseStateForUnknown(),
								},
								Validators: []validator.Number{
									// validation for -1..99
									numbervalidator.OneOf(maxSharingRange...),
								},
							},
						},
					},
					"usage_insights": schema.SingleNestedAttribute{
						// type: AdminDigest -> Usage Insights
						// IncludeOnHomePageInsights, ExcludeEnvironmentFromAnalysis
						// false, true (when unchecked Include Insights)
						// false, false (when checked Include Insights)
						MarkdownDescription: "Usage Insights",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"insights_enabled": schema.BoolAttribute{
								MarkdownDescription: "Inculde insights for all Managed Environment in this group in weekly email digest.",
								Required:            true,
							},
						},
					},
					"maker_welcome_content": schema.SingleNestedAttribute{
						// type: MakerOnboarding -> Maker welcome content
						// makerOnboardingUrl, makerOnboardingMarkdown, makerOnboardingTimestamp
						MarkdownDescription: "Maker Welcome Content",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"maker_onboarding_url": schema.StringAttribute{
								MarkdownDescription: "Maker onboarding URL",
								Required:            true,
							},
							"maker_onboarding_markdown": schema.StringAttribute{
								MarkdownDescription: "Maker onboarding markdown",
								Required:            true,
							},
						},
					},
					"solution_checker_enforcement": schema.SingleNestedAttribute{
						// SolutionChecker -> Solution checker enforcement
						// solutionCheckerMode, suppressValidationEmails(checkbox), solutionCheckerRuleOverrides
						// none/warm/block, false, ""
						// warm, true, ""
						MarkdownDescription: "Solution Checker Enforcement",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"solution_checker_mode": schema.StringAttribute{
								MarkdownDescription: "Solution checker enforceemnt mode: none, warm, block",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("none", "warn", "block"),
								},
							},
							"send_emails_enabled": schema.BoolAttribute{
								MarkdownDescription: "Send emails only when solution is blocked, if unchecked you'll also get emails when there are warnings",
								Required:            true,
							},
						},
					},
					"backup_retention": schema.SingleNestedAttribute{
						// Lifecycle -> Backup retention
						// RetentionPeriod: 14.00:00:00 / 7 / 21 / 28
						MarkdownDescription: "Backup Retention",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"period_in_days": schema.Int32Attribute{
								MarkdownDescription: "Backup retention period in days: 7, 14, 21, 28",
								Required:            true,
								Validators: []validator.Int32{
									int32validator.OneOf(7, 14, 21, 28),
								},
							},
						},
					},
					"ai_generated_descriptions": schema.SingleNestedAttribute{
						// Copilot -> AI generative description
						// DisableAiGeneratedDescriptions (checkbox) //Enable AI generated description
						// false (when checked as true)
						MarkdownDescription: "AI Generated Descriptions",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"ai_description_enabled": schema.BoolAttribute{
								MarkdownDescription: "Enable AI generated description",
								Required:            true,
							},
						},
					},
					"ai_generative_settings": schema.SingleNestedAttribute{
						// GenerativeAISettings -> AI generative settings
						// crossGeoCopilotDataMovementEnabled // Move data across regions enabled
						// bingChatEnabled //Bing Seach enbaled
						MarkdownDescription: "AI Generative Settings",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"move_data_across_regions_enabled": schema.BoolAttribute{
								MarkdownDescription: "Agree to move data across regions",
								Required:            true,
							},
							"bing_search_enabled": schema.BoolAttribute{
								MarkdownDescription: "Agree to enable Bing search features",
								Required:            true,
							},
						},
					},
				},
			},
		},
	}
}

func (r *environmentGroupRuleSetResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config environmentGroupRuleSetResourceModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &config)...); resp.Diagnostics.HasError() {
		return
	}

	if config.Rules.IsNull() || config.Rules.IsUnknown() {
		return
	}

	sharingControlsObj := config.Rules.Attributes()["sharing_controls"]
	if !sharingControlsObj.IsNull() && !sharingControlsObj.IsUnknown() {
		var sharingControl environmentGroupRuleSetSharingControlsModel
		sharingControlsObj.(basetypes.ObjectValue).As(ctx, &sharingControl, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if sharingControl.ShareMode.ValueString() == "no limit" {
			if !sharingControl.ShareMaxLimit.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("rules"),
					"sharing_controls validation error",
					"'share_max_limit' must be null when 'share_mode' is 'no limit'",
				)
			}
		} else {
			if sharingControl.ShareMaxLimit.IsNull() || sharingControl.ShareMaxLimit.Equal(basetypes.NewFloat64Value(-1)) {
				resp.Diagnostics.AddAttributeError(
					path.Root("rules"),
					"sharing_controls validation error",
					"'share_max_limit' must be a value between 0 and 99 when 'share_mode' is 'exclude sharing with security groups'",
				)
			}
		}
	}
}

func (r *environmentGroupRuleSetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
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
	r.EnvironmentGroupRuleSetClient = NewEnvironmentGroupRuleSetClient(client, tenant.NewTenantClient(client))
}

func (r *environmentGroupRuleSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	state := environmentGroupRuleSetResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleSetDto, err := r.EnvironmentGroupRuleSetClient.GetEnvironmentGroupRuleSet(ctx, state.EnvironmentGroupId.ValueString())
	if err != nil {
		if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to get environment group ruleset", err.Error())
		return
	}
	newState, err := convertEnvironmentGroupRuleSetDtoToModel(*ruleSetDto)
	newState.Timeouts = state.Timeouts
	if err != nil {
		resp.Diagnostics.AddError("Failed to convert environment group ruleset dto to model", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *environmentGroupRuleSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *environmentGroupRuleSetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plannedRuleSetDto := convertEnvironmentGroupRuleSetResourceModelToDto(ctx, *plan)
	createdRuleSetDto, err := r.EnvironmentGroupRuleSetClient.CreateEnvironmentGroupRuleSet(ctx, plan.EnvironmentGroupId.ValueString(), plannedRuleSetDto)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create environment group ruleset", err.Error())
		return
	}

	plan.Id = types.StringPointerValue(createdRuleSetDto.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *environmentGroupRuleSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *environmentGroupRuleSetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *environmentGroupRuleSetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plannedRuleSetDto := convertEnvironmentGroupRuleSetResourceModelToDto(ctx, *plan)
	updatedRuleSetDto, err := r.EnvironmentGroupRuleSetClient.UpdateEnvironmentGroupRuleSet(ctx, state.Id.ValueString(), plannedRuleSetDto)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update environment group ruleset", err.Error())
		return
	}

	newState, err := convertEnvironmentGroupRuleSetDtoToModel(*updatedRuleSetDto)
	newState.Timeouts = plan.Timeouts
	if err != nil {
		resp.Diagnostics.AddError("Failed to convert environment group ruleset dto to model", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *environmentGroupRuleSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *environmentGroupRuleSetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.EnvironmentGroupRuleSetClient.DeleteEnvironmentGroupRuleSet(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}
}

func (r *environmentGroupRuleSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("environment_group_id"), req, resp)
}
