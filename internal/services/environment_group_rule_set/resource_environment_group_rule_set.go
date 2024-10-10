// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_set

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
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
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ resource.Resource = &environmentGroupRuleSetResource{}
var _ resource.ResourceWithImportState = &environmentGroupRuleSetResource{}

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
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			// "timeouts": timeouts.Attributes(ctx, timeouts.Opts{
			// 	Read: true,
			// }),
			// TODO env filter
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
				},
			},
			"rules": schema.SetNestedAttribute{
				MarkdownDescription: "Set of rules",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the rule",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("Sharing controls", "Usage insights", "Maker welcome content", "Solution checker enforcement", "Backup retention", "AI generated descriptions", "AI generative settings"),
							},
						},
						// "resource_type": schema.StringAttribute{
						// 	MarkdownDescription: "Resource type",
						// 	Optional:            true,
						// 	Computed:            true,
						// 	Validators: []validator.String{
						// 		stringvalidator.OneOf("App"),
						// 	},
						// },
						"values": schema.SingleNestedAttribute{
							MarkdownDescription: "Configuration values",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								// type: Sharing -> Sharing controls for Canvas apps
								// CanShareWithSecurityGroups: noLimit, excludeSharingToSecurityGroups
								// IsGroupSharingDisabled: true, false
								// MaximumShareLimit: (-1..99)
								//
								// modes:
								// noLimit, false, -1
								// excludeSharingToSecurityGroups, true, (-1.....99)
								// TODO check noLimit, TRUE, -1.
								"share_mode": schema.StringAttribute{
									// noLimit, true
									// excludeSharingToSecurityGroups, false
									MarkdownDescription: "To be used together with `Sharing controls`.\n\nShare mode for canvas apps: `No limit`, `Exclude sharing with security groups`",
									Optional:            true,
									Computed:            true,
									Validators: []validator.String{
										stringvalidator.OneOf("no limit", "exclude sharing with security groups"),
										stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("share_max_limit")),
										AlsoRequiresValueString(types.StringValue("Sharing controls"), path.MatchRelative().AtParent().AtParent().AtName("type")),
									},
								},
								"share_max_limit": schema.Int32Attribute{
									// (-1..99)
									// if noLimit then send -1
									MarkdownDescription: "To be used together with `Sharing controls`.\n\nMaximum total of individual who can be shared to: (-1..99). If `share_mode` is `No limit`, this value must be -1.",
									Optional:            true,
									Computed:            true,
									Validators: []validator.Int32{
										int32validator.Between(-1, 99),
										AlsoRequiresValueInt32(types.StringValue("Sharing controls"), path.MatchRelative().AtParent().AtParent().AtName("type")),
									},
								},

								// type: AdminDigest -> Usage Insights
								// IncludeOnHomePageInsights, ExcludeEnvironmentFromAnalysis
								// false, true (when unchecked Include Insights)
								// false, false (when checked Include Insights)
								"insights_enabled": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "To be used together with `Usage insights`.\n\nInculde insights for all Managed Environment in this group in weekly email digest",
									Validators: []validator.Bool{
										AlsoRequiresValueBool(types.StringValue("Usage insights"), path.MatchRelative().AtParent().AtParent().AtName("type")),
									},
								},

								// type: MakerOnboarding -> Maker welcome content
								// makerOnboardingUrl, makerOnboardingMarkdown, makerOnboardingTimestamp
								// send value or send "" not null
								"onboarding_url": schema.StringAttribute{
									MarkdownDescription: "To be used together with `Maker welcome content`.\n\nMaker onboarding url",
									Computed:            true,
									Optional:            true,
								},
								"onboarding_markdown": schema.StringAttribute{
									MarkdownDescription: "To be used together with `Maker welcome content`.\n\nMaker onboarding markdown",
									Computed:            true,
									Optional:            true,
								},

								// SolutionChecker -> Solution checker enforcement
								// solutionCheckerMode, suppressValidationEmails(checkbox), solutionCheckerRuleOverrides
								// none/warm/block, false, ""
								// warm, true, ""
								"solution_checker_mode": schema.StringAttribute{
									MarkdownDescription: "To be used together with `Solution checker enforcement`.\n\nSolution checker enforceemnt mode: None, Warm, Block",
									Computed:            true,
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.OneOf("none", "warm", "block"),
									},
								},
								"send_emails_enabled": schema.BoolAttribute{
									MarkdownDescription: "To be used together with `Solution checker enforcement`.\n\nSend emails only when solution is blocked, if unchecked you'll also get emails when there are warnings",
									Optional:            true,
								},

								// Lifecycle -> Backup retention
								// RetentionPeriod: 14.00:00:00 / 7 / 21 / 28
								"period_in_days": schema.Int32Attribute{
									MarkdownDescription: "To be used together with `Backup retention`.\n\nBackup retention period in days: 7, 14, 21, 28",
									Computed:            true,
									Optional:            true,
									Validators: []validator.Int32{
										int32validator.OneOf(7, 14, 21, 28),
									},
								},

								// Copilot -> AI generative description
								// DisableAiGeneratedDescriptions (checkbox) //Enable AI generated description
								// false (when checked as true)
								"ai_description_enabled": schema.BoolAttribute{
									MarkdownDescription: "To be used together with `AI generated descriptions`.\n\nEnable AI generated description",
									Computed:            true,
									Optional:            true,
								},

								// GenerativeAISettings -> AI generative settings
								// crossGeoCopilotDataMovementEnabled // Move data across regions enabled
								// bingChatEnabled //Bing Seach enbaled
								"move_data_across_regions_enabled": schema.BoolAttribute{
									MarkdownDescription: "To be used together with `AI generative settings`.\n\nAgree to move data across regions",
									Computed:            true,
									Optional:            true,
								},
								"bing_search_enabled": schema.BoolAttribute{
									MarkdownDescription: "To be used together with `AI generative settings`.\n\nAgree to enable Bing search features",
									Computed:            true,
									Optional:            true,
								},
							},
						},
					},
				},
			},
		},
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

	r.EnvironmentGroupRuleSetClient = newEnvironmentGroupRuleSetClient(client)
}

func (r *environmentGroupRuleSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	state := environmentGroupRuleSetResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	panic("read")

	// envGroupRuleSet, err := r.EnvironmentGroupRuleSetClient.GetEnvironmentGroupRuleSet(ctx, state.Id.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
	// 	return
	// }

	//TODO what happens when rules set is deleted??
	// if envGroupRuleSet == nil {
	// 	resp.State.RemoveResource(ctx)
	// 	return
	// }

	// state = environmentGroupRuleSetResourceModel{}

	// state.EnvironmentGroupId = types.StringValue("bd6b30f1-e31e-4cdd-b82b-689a4b674f2f")
	// state.Id = types.StringValue("1234")

	// var rules []attr.Value

	// rules = append(rules, types.ObjectValueMust(map[string]attr.Type{
	// 	"type":          types.StringType,
	// 	"resource_type": types.StringType,
	// },
	// 	map[string]attr.Value{
	// 		"type":          types.StringValue("Sharing"),
	// 		"resource_type": types.StringValue("App"),
	// 	}))

	// rules = append(rules, types.ObjectValueMust(map[string]attr.Type{
	// 	"type":          types.StringType,
	// 	"resource_type": types.StringType,
	// },
	// 	map[string]attr.Value{
	// 		"type":          types.StringValue("AdminDigest"),
	// 		"resource_type": types.StringNull(),
	// 	}))

	// state.Rules = types.SetValueMust(ruleSetObjectType, rules)

	//state.Id = types.StringValue(envGroupRuleSet.Id)
	//TODO: set rules

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *environmentGroupRuleSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *environmentGroupRuleSetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	panic("asdads")

	// state := environmentGroupRuleSetResourceModel{}

	// state.EnvironmentGroupId = types.StringValue("bd6b30f1-e31e-4cdd-b82b-689a4b674f2f")
	// state.Id = types.StringValue("1234")

	// var rules []attr.Value
	// for _, rule := range plan.Rules.Elements() {

	// 	ruleObj := rule.(types.Object)
	// 	t := ruleObj.Attributes()["type"]
	// 	rt := ruleObj.Attributes()["resource_type"]

	// 	rules = append(rules, types.ObjectValueMust(map[string]attr.Type{
	// 		"type":          types.StringType,
	// 		"resource_type": types.StringType,
	// 	},
	// 		map[string]attr.Value{
	// 			"type":          t,
	// 			"resource_type": rt,
	// 		}))
	// }

	// state.Rules = types.SetValueMust(ruleSetObjectType, rules)

	//resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *environmentGroupRuleSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *environmentGroupRuleSetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *environmentGroupRuleSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("environment_group_id"), req, resp)
}
