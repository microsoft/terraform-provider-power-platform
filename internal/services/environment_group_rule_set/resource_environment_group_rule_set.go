// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_set

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/numbervalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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
		MarkdownDescription: "Allows the creation of environment group rulesets. See [Power Platform documentation](https://learn.microsoft.com/power-platform/admin/environment-groups) for more information on the available rules that can be applied to an environment group.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			// id and policy_id become null when the last rule backed by their respective API is
			// removed, so they cannot be carried forward from prior state.
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique id of the environment group ruleset",
				Computed:            true,
			},
			"policy_id": schema.StringAttribute{
				MarkdownDescription: "Unique id of the rule-based policy backing `maker_welcome_content`, `advanced_connector_policies_only`, `content_security_policy` and `advanced_connector_policies`. Null when none of those rules are configured.",
				Computed:            true,
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
					"advanced_connector_policies_only": schema.SingleNestedAttribute{
						MarkdownDescription: "Controls whether only advanced connector policies are applied to environments in this group. " +
							"When enabled, environments in the group will only use advanced connector policies.",
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Enable advanced connector policies only",
								Required:            true,
							},
						},
					},
					"content_security_policy": schema.SingleNestedAttribute{
						MarkdownDescription: "Configures the [Content Security Policy (CSP)](https://learn.microsoft.com/power-platform/admin/content-security-policy) for Power Apps in this environment group.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Enable Content Security Policy for model-driven apps",
								Required:            true,
							},
							"enabled_for_canvas": schema.BoolAttribute{
								MarkdownDescription: "Enable Content Security Policy for canvas apps",
								Required:            true,
							},
							"enabled_for_code_apps": schema.BoolAttribute{
								MarkdownDescription: "Enable Content Security Policy for code-first apps",
								Required:            true,
							},
							"report_uri": schema.StringAttribute{
								MarkdownDescription: "URI to send CSP violation reports to",
								Optional:            true,
							},
							"reporting_endpoint": schema.StringAttribute{
								MarkdownDescription: "Reporting endpoint for CSP violations",
								Optional:            true,
							},
							"configuration":               cspConfigurationSchema("model-driven apps"),
							"configuration_for_canvas":    cspConfigurationSchema("canvas apps"),
							"configuration_for_code_apps": cspCodeAppsConfigurationSchema("code-first apps"),
						},
					},
					"advanced_connector_policies": schema.SingleNestedAttribute{
						MarkdownDescription: "Manages which connectors are allowed and what actions they can perform in environments within this group.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"allowed_connectors": schema.ListNestedAttribute{
								MarkdownDescription: "List of connectors that are allowed in the environment group",
								Required:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"connector_id": schema.StringAttribute{
											MarkdownDescription: "Short connector identifier (e.g., `shared_commondataservice`). The provider automatically prepends `" + CONNECTOR_API_PREFIX + "`.",
											Required:            true,
										},
										"actions_mode": schema.StringAttribute{
											MarkdownDescription: "Controls which actions are allowed for this connector. Use `all_allowed` to permit all actions, or `some_allowed` to restrict to specific actions listed in `allowed_actions`.",
											Required:            true,
											Validators: []validator.String{
												stringvalidator.OneOf(actionsModeAllAllowed, actionsModeSomeAllowed),
											},
										},
										"allowed_actions": schema.ListAttribute{
											MarkdownDescription: "List of specific action names allowed for this connector. Only used when `actions_mode` is `some_allowed`.",
											Optional:            true,
											ElementType:         types.StringType,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func cspDirectiveSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"img_src": schema.ListAttribute{
			MarkdownDescription: "Allowed sources for images (`img-src` directive)",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"style_src": schema.ListAttribute{
			MarkdownDescription: "Allowed sources for stylesheets (`style-src` directive)",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"form_action": schema.ListAttribute{
			MarkdownDescription: "Allowed targets for form submissions (`form-action` directive)",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"frame_src": schema.ListAttribute{
			MarkdownDescription: "Allowed sources for frames (`frame-src` directive)",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"connect_src": schema.ListAttribute{
			MarkdownDescription: "Allowed sources for fetch, XHR, WebSocket (`connect-src` directive)",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"font_src": schema.ListAttribute{
			MarkdownDescription: "Allowed sources for fonts (`font-src` directive)",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"script_src": schema.ListAttribute{
			MarkdownDescription: "Allowed sources for scripts (`script-src` directive)",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"frame_ancestor": schema.ListAttribute{
			MarkdownDescription: "Allowed sources that can embed this content (`frame-ancestors` directive)",
			Optional:            true,
			ElementType:         types.StringType,
		},
	}
}

func cspConfigurationSchema(appType string) schema.SingleNestedAttribute {
	attrs := cspDirectiveSchemaAttributes()
	attrs["strict_csp"] = schema.BoolAttribute{
		MarkdownDescription: fmt.Sprintf("When `true`, enables strict Content Security Policy enforcement for %s.", appType),
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
	}

	return schema.SingleNestedAttribute{
		MarkdownDescription: fmt.Sprintf("CSP directive configuration for %s. Each attribute is a list of allowed source URIs.", appType),
		Optional:            true,
		Attributes:          attrs,
	}
}

func cspCodeAppsConfigurationSchema(appType string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: fmt.Sprintf("CSP directive configuration for %s. Each attribute is a list of allowed source URIs.", appType),
		Optional:            true,
		Attributes:          cspDirectiveSchemaAttributes(),
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

	providerClient, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.EnvironmentGroupRuleSetClient = NewEnvironmentGroupRuleSetClient(providerClient.Api, tenant.NewTenantClient(providerClient.Api))
}

// readRuleSet returns the legacy rule set for the group, or nil when the group has none.
func (r *environmentGroupRuleSetResource) readRuleSet(ctx context.Context, environmentGroupId string) (*EnvironmentGroupRuleSetValueSetDto, error) {
	ruleSet, err := r.EnvironmentGroupRuleSetClient.GetEnvironmentGroupRuleSet(ctx, environmentGroupId)
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return ruleSet, nil
}

// readPolicy returns the rule-based policy assigned to the group, or nil when none is assigned.
func (r *environmentGroupRuleSetResource) readPolicy(ctx context.Context, environmentGroupId string) (*ruleBasedPolicyDto, error) {
	assignments, err := r.EnvironmentGroupRuleSetClient.listAssignmentsByEnvironmentGroup(ctx, environmentGroupId)
	if err != nil {
		return nil, err
	}
	if len(assignments) == 0 {
		return nil, nil
	}

	policy, err := r.EnvironmentGroupRuleSetClient.getPolicy(ctx, assignments[0].PolicyId)
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return policy, nil
}

// applyPolicyRules creates or adopts the rule-based policy for the group and returns its id.
func (r *environmentGroupRuleSetResource) applyPolicyRules(ctx context.Context, model environmentGroupRuleSetResourceModel) (string, error) {
	request, err := convertRulesModelToPolicyDto(ctx, model)
	if err != nil {
		return "", err
	}

	environmentGroupId := model.EnvironmentGroupId.ValueString()
	assignments, err := r.EnvironmentGroupRuleSetClient.listAssignmentsByEnvironmentGroup(ctx, environmentGroupId)
	if err != nil {
		return "", fmt.Errorf("failed to check existing policy assignments: %w", err)
	}

	if len(assignments) > 0 {
		// The API allows a single policy per environment group, so adopt the assigned one.
		policyId := assignments[0].PolicyId
		tflog.Debug(ctx, fmt.Sprintf("Adopting rule-based policy %s already assigned to environment group %s", policyId, environmentGroupId))

		existingPolicy, err := r.EnvironmentGroupRuleSetClient.getPolicy(ctx, policyId)
		if err != nil {
			return "", fmt.Errorf("failed to get existing rule-based policy: %w", err)
		}

		// Rule sets can only be dropped through the removeRule endpoint; PUT rejects removals.
		desired := make(map[string]bool, len(request.RuleSets))
		for _, rs := range request.RuleSets {
			desired[rs.Id] = true
		}
		for _, rs := range existingPolicy.RuleSets {
			if _, managed := managedPolicyRuleSetIds[rs.Id]; !managed || desired[rs.Id] {
				continue
			}
			tflog.Debug(ctx, fmt.Sprintf("Removing rule %s from policy %s", rs.Id, policyId))
			if err := r.EnvironmentGroupRuleSetClient.removeRuleFromPolicy(ctx, policyId, existingPolicy.Name, rs.Id, rs.Version); err != nil {
				return "", fmt.Errorf("failed to remove rule %s from rule-based policy: %w", rs.Id, err)
			}
		}

		request.Name = existingPolicy.Name
		request.RuleSets = mergePolicyRuleSets(existingPolicy.RuleSets, request.RuleSets)
		if _, err := r.EnvironmentGroupRuleSetClient.updatePolicy(ctx, policyId, request); err != nil {
			return "", fmt.Errorf("failed to update rule-based policy: %w", err)
		}
		return policyId, nil
	}

	createdPolicy, err := r.EnvironmentGroupRuleSetClient.createPolicy(ctx, request)
	if err != nil {
		return "", fmt.Errorf("failed to create rule-based policy: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Assigning rule-based policy %s to environment group %s", createdPolicy.Id, environmentGroupId))
	if _, err := r.EnvironmentGroupRuleSetClient.createEnvironmentGroupAssignment(ctx, createdPolicy.Id, environmentGroupId); err != nil {
		return "", fmt.Errorf("failed to assign rule-based policy to environment group: %w", err)
	}

	return createdPolicy.Id, nil
}

// releasePolicy deletes the policy when this resource created it, otherwise it only strips
// the rule sets this resource manages so pre-existing configuration survives.
func (r *environmentGroupRuleSetResource) releasePolicy(ctx context.Context, policyId, environmentGroupId string) error {
	existingPolicy, err := r.EnvironmentGroupRuleSetClient.getPolicy(ctx, policyId)
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			return nil
		}
		return err
	}

	if existingPolicy.Name == terraformManagedPolicyName(environmentGroupId) {
		// A policy cannot be deleted while it still has assignments.
		tflog.Debug(ctx, fmt.Sprintf("Deleting rule-based policy %s created for environment group %s", policyId, environmentGroupId))
		if err := r.EnvironmentGroupRuleSetClient.deleteEnvironmentGroupAssignment(ctx, policyId, environmentGroupId); err != nil {
			return err
		}
		return r.EnvironmentGroupRuleSetClient.deletePolicy(ctx, policyId)
	}

	present := make(map[string]bool, len(existingPolicy.RuleSets))
	for _, rs := range existingPolicy.RuleSets {
		present[rs.Id] = true
	}

	for ruleSetId, version := range managedPolicyRuleSetIds {
		if !present[ruleSetId] {
			continue
		}
		tflog.Debug(ctx, fmt.Sprintf("Removing rule %s from adopted policy %s", ruleSetId, policyId))
		if err := r.EnvironmentGroupRuleSetClient.removeRuleFromPolicy(ctx, policyId, existingPolicy.Name, ruleSetId, version); err != nil {
			return err
		}
	}

	return nil
}

func (r *environmentGroupRuleSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	state := environmentGroupRuleSetResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	environmentGroupId := state.EnvironmentGroupId.ValueString()

	ruleSet, err := r.readRuleSet(ctx, environmentGroupId)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get environment group ruleset", err.Error())
		return
	}

	policy, err := r.readPolicy(ctx, environmentGroupId)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get environment group rule-based policy", err.Error())
		return
	}

	if ruleSet == nil && policy == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	newState, err := convertEnvironmentGroupRuleSetDtoToModel(environmentGroupId, ruleSet, policy)
	if err != nil {
		resp.Diagnostics.AddError("Failed to convert environment group ruleset dto to model", err.Error())
		return
	}
	newState.Timeouts = state.Timeouts

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

	plan.Id = types.StringNull()
	plan.PolicyId = types.StringNull()

	if hasLegacyRules(plan.Rules) {
		plannedRuleSetDto, err := convertEnvironmentGroupRuleSetResourceModelToDto(ctx, *plan)
		if err != nil {
			resp.Diagnostics.AddError("Conversion Error", err.Error())
			return
		}
		createdRuleSetDto, err := r.EnvironmentGroupRuleSetClient.CreateEnvironmentGroupRuleSet(ctx, plan.EnvironmentGroupId.ValueString(), plannedRuleSetDto)
		if err != nil {
			resp.Diagnostics.AddError("Failed to create environment group ruleset", err.Error())
			return
		}
		plan.Id = types.StringPointerValue(createdRuleSetDto.Id)
	}

	if hasPolicyRules(plan.Rules) {
		policyId, err := r.applyPolicyRules(ctx, *plan)
		if err != nil {
			resp.Diagnostics.AddError("Failed to apply environment group rule-based policy", err.Error())
			return
		}
		plan.PolicyId = types.StringValue(policyId)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, *plan)...)
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

	plan.Id = types.StringNull()
	plan.PolicyId = types.StringNull()

	if hasLegacyRules(plan.Rules) {
		// The legacy API expects the existing rule set id in the payload when updating.
		plan.Id = state.Id
		plannedRuleSetDto, err := convertEnvironmentGroupRuleSetResourceModelToDto(ctx, *plan)
		if err != nil {
			resp.Diagnostics.AddError("Conversion Error", err.Error())
			return
		}

		if state.Id.IsNull() {
			createdRuleSetDto, err := r.EnvironmentGroupRuleSetClient.CreateEnvironmentGroupRuleSet(ctx, plan.EnvironmentGroupId.ValueString(), plannedRuleSetDto)
			if err != nil {
				resp.Diagnostics.AddError("Failed to create environment group ruleset", err.Error())
				return
			}
			plan.Id = types.StringPointerValue(createdRuleSetDto.Id)
		} else {
			if _, err := r.EnvironmentGroupRuleSetClient.UpdateEnvironmentGroupRuleSet(ctx, state.Id.ValueString(), plannedRuleSetDto); err != nil {
				resp.Diagnostics.AddError("Failed to update environment group ruleset", err.Error())
				return
			}
		}
	} else if !state.Id.IsNull() {
		if err := r.EnvironmentGroupRuleSetClient.DeleteEnvironmentGroupRuleSet(ctx, state.Id.ValueString()); err != nil {
			resp.Diagnostics.AddError("Failed to delete environment group ruleset", err.Error())
			return
		}
	}

	if hasPolicyRules(plan.Rules) {
		policyId, err := r.applyPolicyRules(ctx, *plan)
		if err != nil {
			resp.Diagnostics.AddError("Failed to apply environment group rule-based policy", err.Error())
			return
		}
		plan.PolicyId = types.StringValue(policyId)
	} else if !state.PolicyId.IsNull() {
		if err := r.releasePolicy(ctx, state.PolicyId.ValueString(), plan.EnvironmentGroupId.ValueString()); err != nil {
			resp.Diagnostics.AddError("Failed to remove rules from environment group rule-based policy", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, *plan)...)
}

func (r *environmentGroupRuleSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *environmentGroupRuleSetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.Id.IsNull() {
		if err := r.EnvironmentGroupRuleSetClient.DeleteEnvironmentGroupRuleSet(ctx, state.Id.ValueString()); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
			return
		}
	}

	if !state.PolicyId.IsNull() {
		if err := r.releasePolicy(ctx, state.PolicyId.ValueString(), state.EnvironmentGroupId.ValueString()); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
			return
		}
	}
}

func (r *environmentGroupRuleSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("environment_group_id"), req, resp)
}
