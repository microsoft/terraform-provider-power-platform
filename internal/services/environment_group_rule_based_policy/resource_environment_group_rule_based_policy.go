// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_based_policy

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ resource.Resource = &environmentGroupRuleBasedPolicyResource{}
var _ resource.ResourceWithImportState = &environmentGroupRuleBasedPolicyResource{}
var _ resource.ResourceWithValidateConfig = &environmentGroupRuleBasedPolicyResource{}

func NewEnvironmentGroupRuleBasedPolicyResource() resource.Resource {
	return &environmentGroupRuleBasedPolicyResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment_group_rule_based_policy",
		},
	}
}

func (r *environmentGroupRuleBasedPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *environmentGroupRuleBasedPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages rule-based policies for environment groups using the Power Platform governance API. " +
			"See [Power Platform documentation](https://learn.microsoft.com/power-platform/admin/environment-groups) for more information.\n\n" +
			"~> **Note:** This resource is available as **preview**.\n\n" +
			"~> **Important:** If a rule-based policy is already assigned to the specified environment group, this resource will **adopt** the existing policy " +
			"and update it with the desired rule sets rather than creating a new one. The Power Platform API only allows a single policy assignment per environment group. " +
			"On `destroy`, the managed rule sets are removed from the policy using the [Remove Rule](https://learn.microsoft.com/en-us/rest/api/power-platform/governance/rule-based-policies/remove-rule-from-rule-based-policy) endpoint, " +
			"but the policy itself is not deleted.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the rule-based policy",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_group_id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the environment group to assign the policy to",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rule_sets": schema.SingleNestedAttribute{
				MarkdownDescription: "Rule sets to apply to the environment group",
				Required:            true,
				Attributes: map[string]schema.Attribute{
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
							"configuration":               cspConfigurationSchema("model-driven apps", true),
							"configuration_for_canvas":    cspConfigurationSchema("canvas apps", true),
							"configuration_for_code_apps": cspConfigurationSchema("code-first apps", false),
						},
					},
					"advanced_connector_policies": schema.SingleNestedAttribute{
						MarkdownDescription: "Manages which connectors are allowed and what actions they can perform in environments within this group (API rule set ID: `ConnectorManagement`).",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"allowed_connectors": schema.ListNestedAttribute{
								MarkdownDescription: "List of connectors that are allowed in the environment group",
								Required:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"connector_id": schema.StringAttribute{
											MarkdownDescription: "Short connector identifier (e.g., `shared_commondataservice`). The provider automatically prepends `/providers/Microsoft.PowerApps/apis/`.",
											Required:            true,
										},
										"actions_mode": schema.StringAttribute{
											MarkdownDescription: "Controls which actions are allowed for this connector. Use `all_allowed` to permit all actions, or `some_allowed` to restrict to specific actions listed in `allowed_actions`.",
											Required:            true,
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

func cspConfigurationSchema(appType string, includeStrictCsp bool) schema.SingleNestedAttribute {
	attrs := map[string]schema.Attribute{
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

	if includeStrictCsp {
		attrs["strict_csp"] = schema.BoolAttribute{
			MarkdownDescription: fmt.Sprintf("When `true`, enables strict Content Security Policy enforcement for %s. This contributes to the computed `ContentSecurityPolicyOptions` value sent to the API.", appType),
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		}
	}

	return schema.SingleNestedAttribute{
		MarkdownDescription: fmt.Sprintf("CSP directive configuration for %s. Each attribute is a list of allowed source URIs.", appType),
		Optional:            true,
		Attributes:          attrs,
	}
}

func (r *environmentGroupRuleBasedPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
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
	r.RuleBasedPolicyClient = NewRuleBasedPolicyClient(providerClient.Api)
}

func (r *environmentGroupRuleBasedPolicyResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config environmentGroupRuleBasedPolicyResourceModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &config)...); resp.Diagnostics.HasError() {
		return
	}

	if config.RuleSets.IsNull() || config.RuleSets.IsUnknown() {
		return
	}

	ruleSetsAttrs := config.RuleSets.Attributes()
	hasAtLeastOne := false
	for _, attr := range ruleSetsAttrs {
		if objAttr, ok := attr.(types.Object); ok && !objAttr.IsNull() && !objAttr.IsUnknown() {
			hasAtLeastOne = true
			break
		}
	}

	if !hasAtLeastOne {
		resp.Diagnostics.AddAttributeError(
			path.Root("rule_sets"),
			"At least one rule set is required",
			"The rule_sets block must contain at least one of: advanced_connector_policies_only, content_security_policy, or advanced_connector_policies.",
		)
	}
}

func (r *environmentGroupRuleBasedPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan environmentGroupRuleBasedPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyRequest, err := convertModelToDto(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Conversion Error", err.Error())
		return
	}

	groupId := plan.EnvironmentGroupId.ValueString()

	// Check if a policy is already assigned to this environment group
	tflog.Debug(ctx, fmt.Sprintf("Checking for existing policy assignments on environment group %s", groupId))
	assignments, err := r.RuleBasedPolicyClient.ListAssignmentsByEnvironmentGroup(ctx, groupId)
	if err != nil {
		resp.Diagnostics.AddError("Failed to check existing policy assignments", err.Error())
		return
	}

	var policyId string
	if len(assignments) > 0 {
		// A policy already exists for this environment group — adopt it and update with desired state
		policyId = assignments[0].PolicyId
		tflog.Debug(ctx, fmt.Sprintf("Found existing policy %s assigned to environment group %s. Adopting and updating with desired state.", policyId, groupId))

		// Get the current policy to preserve its name and unmanaged rule sets
		existingPolicy, getErr := r.RuleBasedPolicyClient.GetPolicy(ctx, policyId)
		if getErr != nil {
			resp.Diagnostics.AddError("Failed to get existing rule-based policy", getErr.Error())
			return
		}
		policyRequest.Name = existingPolicy.Name
		policyRequest.RuleSets = mergeRuleSets(existingPolicy.RuleSets, policyRequest.RuleSets)

		_, err = r.RuleBasedPolicyClient.UpdatePolicy(ctx, policyId, policyRequest)
		if err != nil {
			resp.Diagnostics.AddError("Failed to update existing rule-based policy", err.Error())
			return
		}
	} else {
		// No existing policy — create a new one and assign it
		tflog.Debug(ctx, "No existing policy found. Creating new rule-based policy.")
		createdPolicy, createErr := r.RuleBasedPolicyClient.CreatePolicy(ctx, policyRequest)
		if createErr != nil {
			resp.Diagnostics.AddError("Failed to create rule-based policy", createErr.Error())
			return
		}

		policyId = createdPolicy.Id
		tflog.Debug(ctx, fmt.Sprintf("Created rule-based policy with ID: %s", policyId))

		tflog.Debug(ctx, fmt.Sprintf("Assigning policy %s to environment group %s", policyId, groupId))
		_, err = r.RuleBasedPolicyClient.CreateEnvironmentGroupAssignment(ctx, policyId, groupId)
		if err != nil {
			resp.Diagnostics.AddError("Failed to assign policy to environment group", err.Error())
			return
		}
	}

	plan.Id = types.StringValue(policyId)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *environmentGroupRuleBasedPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state environmentGroupRuleBasedPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading rule-based policy %s", state.Id.ValueString()))
	policy, err := r.RuleBasedPolicyClient.GetPolicy(ctx, state.Id.ValueString())
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to get rule-based policy", err.Error())
		return
	}

	newState, err := convertDtoToModel(*policy, state.EnvironmentGroupId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to convert rule-based policy dto to model", err.Error())
		return
	}
	newState.Timeouts = state.Timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *environmentGroupRuleBasedPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan environmentGroupRuleBasedPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state environmentGroupRuleBasedPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyRequest, err := convertModelToDto(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Conversion Error", err.Error())
		return
	}

	policyId := state.Id.ValueString()

	// Get current policy to preserve its name and unmanaged rule sets
	existingPolicy, getErr := r.RuleBasedPolicyClient.GetPolicy(ctx, policyId)
	if getErr != nil {
		resp.Diagnostics.AddError("Failed to get existing rule-based policy for update", getErr.Error())
		return
	}
	policyRequest.Name = existingPolicy.Name
	policyRequest.RuleSets = mergeRuleSets(existingPolicy.RuleSets, policyRequest.RuleSets)

	tflog.Debug(ctx, fmt.Sprintf("Updating rule-based policy %s", policyId))
	updatedPolicy, err := r.RuleBasedPolicyClient.UpdatePolicy(ctx, policyId, policyRequest)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update rule-based policy", err.Error())
		return
	}

	newState, err := convertDtoToModel(*updatedPolicy, plan.EnvironmentGroupId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to convert rule-based policy dto to model", err.Error())
		return
	}
	newState.Timeouts = plan.Timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *environmentGroupRuleBasedPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state environmentGroupRuleBasedPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyId := state.Id.ValueString()

	// Get the current policy to obtain its name (required by the removeRule endpoint)
	existingPolicy, err := r.RuleBasedPolicyClient.GetPolicy(ctx, policyId)
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			return
		}
		resp.Diagnostics.AddError("Failed to get rule-based policy for deletion", err.Error())
		return
	}

	// Remove only managed rule sets that actually exist in the policy
	existingRuleIds := make(map[string]bool)
	for _, rs := range existingPolicy.RuleSets {
		existingRuleIds[rs.Id] = true
	}

	for ruleSetId, version := range managedRuleSetIds {
		if !existingRuleIds[ruleSetId] {
			tflog.Debug(ctx, fmt.Sprintf("Skipping removal of rule %s — not present in policy %s", ruleSetId, policyId))
			continue
		}
		tflog.Debug(ctx, fmt.Sprintf("Removing rule %s from policy %s", ruleSetId, policyId))
		err := r.RuleBasedPolicyClient.RemoveRuleFromPolicy(ctx, policyId, existingPolicy.Name, ruleSetId, version)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Failed to remove rule %s from policy", ruleSetId), err.Error())
			return
		}
	}
}

func (r *environmentGroupRuleBasedPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	environmentGroupId := req.ID
	tflog.Debug(ctx, fmt.Sprintf("Importing rule-based policy for environment group %s", environmentGroupId))

	assignments, err := r.RuleBasedPolicyClient.ListAssignmentsByEnvironmentGroup(ctx, environmentGroupId)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list assignments for environment group", err.Error())
		return
	}

	if len(assignments) == 0 {
		resp.Diagnostics.AddError(
			"No rule-based policy found",
			fmt.Sprintf("No rule-based policy assignment found for environment group %s", environmentGroupId),
		)
		return
	}

	policyId := assignments[0].PolicyId
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), policyId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_group_id"), environmentGroupId)...)
}
