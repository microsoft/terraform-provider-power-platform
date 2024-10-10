// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_set

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

var _ resource.Resource = &EnvironmentGroupRuleSetResource{}
var _ resource.ResourceWithImportState = &EnvironmentGroupRuleSetResource{}

func NewEnvironmentGroupRuleSetResource() resource.Resource {
	return &EnvironmentGroupRuleSetResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment_group_rule_set",
		},
	}
}

func (r *EnvironmentGroupRuleSetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *EnvironmentGroupRuleSetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
								stringvalidator.OneOf("Sharing", "AdminDigest", "MakerOnboarding", "SolutionChecker", "Lifecycle", "Copilot", "GenerativeAISettings"),
							},
						},
						"resource_type": schema.StringAttribute{
							MarkdownDescription: "Resource type",
							Optional:            true,
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("App"),
							},
						},
					},
				},
			},
		},
	}
}

func (r *EnvironmentGroupRuleSetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EnvironmentGroupRuleSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	state := EnvironmentGroupRuleSetResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

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

	state = EnvironmentGroupRuleSetResourceModel{}

	state.EnvironmentGroupId = types.StringValue("bd6b30f1-e31e-4cdd-b82b-689a4b674f2f")
	state.Id = types.StringValue("1234")

	var rules []attr.Value

	rules = append(rules, types.ObjectValueMust(map[string]attr.Type{
		"type":          types.StringType,
		"resource_type": types.StringType,
	},
		map[string]attr.Value{
			"type":          types.StringValue("Sharing"),
			"resource_type": types.StringValue("App"),
		}))

	rules = append(rules, types.ObjectValueMust(map[string]attr.Type{
		"type":          types.StringType,
		"resource_type": types.StringType,
	},
		map[string]attr.Value{
			"type":          types.StringValue("AdminDigest"),
			"resource_type": types.StringNull(),
		}))

	state.Rules = types.SetValueMust(ruleSetObjectType, rules)

	//state.Id = types.StringValue(envGroupRuleSet.Id)
	//TODO: set rules

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

var ruleSetObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"type":          types.StringType,
		"resource_type": types.StringType,
	},
}

func (r *EnvironmentGroupRuleSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *EnvironmentGroupRuleSetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := EnvironmentGroupRuleSetResourceModel{}

	state.EnvironmentGroupId = types.StringValue("bd6b30f1-e31e-4cdd-b82b-689a4b674f2f")
	state.Id = types.StringValue("1234")

	var rules []attr.Value
	for _, rule := range plan.Rules.Elements() {

		ruleObj := rule.(types.Object)
		t := ruleObj.Attributes()["type"]
		rt := ruleObj.Attributes()["resource_type"]

		rules = append(rules, types.ObjectValueMust(map[string]attr.Type{
			"type":          types.StringType,
			"resource_type": types.StringType,
		},
			map[string]attr.Value{
				"type":          t,
				"resource_type": rt,
			}))
	}

	state.Rules = types.SetValueMust(ruleSetObjectType, rules)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EnvironmentGroupRuleSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *EnvironmentGroupRuleSetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *EnvironmentGroupRuleSetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EnvironmentGroupRuleSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *EnvironmentGroupRuleSetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EnvironmentGroupRuleSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("environment_group_id"), req, resp)
}
