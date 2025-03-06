// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package solution_checker_rules

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var (
	_ datasource.DataSource              = &DataSource{}
	_ datasource.DataSourceWithConfigure = &DataSource{}
)

func NewSolutionCheckerRulesDataSource() datasource.DataSource {
	return &DataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "solution_checker_rules",
		},
	}
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	d.ProviderTypeName = req.ProviderTypeName
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of solution checker rules for a Power Platform environment. Solution checker helps identify potential issues in solutions by analyzing components against a set of best practice rules. This data source can be used to retrieve the available rules for configuration in managed environments.",
		MarkdownDescription: "Fetches the list of solution checker rules for a Power Platform environment. Solution checker helps identify potential issues in solutions by analyzing components against a set of best practice rules. This data source can be used to retrieve the available rules for configuration in managed environments.\n\nAdditional Resources:\n\n* [Managed Environment Solution Checker](https://learn.microsoft.com/en-us/power-platform/admin/managed-environment-solution-checker)\n",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"environment_id": schema.StringAttribute{
				Description:         "The ID of the environment to retrieve solution checker rules from",
				MarkdownDescription: "The ID of the environment to retrieve solution checker rules from",
				Required:            true,
			},
			"rules": schema.ListNestedAttribute{
				Description:         "List of solution checker rules",
				MarkdownDescription: "List of solution checker rules",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"code": schema.StringAttribute{
							Description:         "The unique code of the rule",
							MarkdownDescription: "The unique code of the rule",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							Description:         "A detailed description of the rule",
							MarkdownDescription: "A detailed description of the rule",
							Computed:            true,
						},
						"summary": schema.StringAttribute{
							Description:         "A brief summary of the rule",
							MarkdownDescription: "A brief summary of the rule",
							Computed:            true,
						},
						"how_to_fix": schema.StringAttribute{
							Description:         "Instructions on how to fix issues identified by the rule",
							MarkdownDescription: "Instructions on how to fix issues identified by the rule",
							Computed:            true,
						},
						"guidance_url": schema.StringAttribute{
							Description:         "URL to detailed guidance on addressing the issue",
							MarkdownDescription: "URL to detailed guidance on addressing the issue",
							Computed:            true,
						},
						"component_type": schema.Int64Attribute{
							Description:         "The type of component this rule applies to",
							MarkdownDescription: "The type of component this rule applies to",
							Computed:            true,
						},
						"component_type_description": schema.StringAttribute{
							Description:         "Description of the component type",
							MarkdownDescription: "Description of the component type",
							Computed:            true,
						},
						"primary_category": schema.Int64Attribute{
							Description:         "The primary category of the rule",
							MarkdownDescription: "The primary category of the rule",
							Computed:            true,
						},
						"primary_category_description": schema.StringAttribute{
							Description:         "Description of the primary category",
							MarkdownDescription: "Description of the primary category",
							Computed:            true,
						},
						"include": schema.BoolAttribute{
							Description:         "Whether the rule is included/enabled by default",
							MarkdownDescription: "Whether the rule is included/enabled by default",
							Computed:            true,
						},
						"severity": schema.Int64Attribute{
							Description:         "The severity level of the rule (1-5, with 5 being most severe)",
							MarkdownDescription: "The severity level of the rule (1-5, with 5 being most severe)",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig. It's ok.
		return
	}
	client, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.SolutionCheckerRulesClient = newSolutionCheckerRulesClient(client.Api)
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var state DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	environmentId := state.EnvironmentId.ValueString()
	rules, err := d.SolutionCheckerRulesClient.GetSolutionCheckerRules(ctx, environmentId)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	state.Rules = []RuleModel{}
	for _, rule := range rules {
		ruleModel := convertFromRuleDto(rule)
		state.Rules = append(state.Rules, ruleModel)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
