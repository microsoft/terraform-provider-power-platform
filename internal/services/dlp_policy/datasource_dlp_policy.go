// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package dlp_policy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var (
	_ datasource.DataSource              = &DataLossPreventionPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &DataLossPreventionPolicyDataSource{}
)

func NewDataLossPreventionPolicyDataSource() datasource.DataSource {
	return &DataLossPreventionPolicyDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "data_loss_prevention_policies",
		},
	}
}

func (d *DataLossPreventionPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *DataLossPreventionPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	connectorSchema := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the connector",
				Optional:            true,
			},
			"default_action_rule_behavior": schema.StringAttribute{
				MarkdownDescription: "Default action rule behavior for the connector (\"Allow\", \"Block\")",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Allow", "Block", ""),
				},
			},
			"action_rules": schema.ListNestedAttribute{
				MarkdownDescription: "Action rules for the connector",
				Optional:            true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"action_id": schema.StringAttribute{
							MarkdownDescription: "ID of the action rule",
							Required:            true,
						},
						"behavior": schema.StringAttribute{
							MarkdownDescription: "Behavior of the action rule (\"Allow\", \"Block\")",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("Allow", "Block"),
							},
						},
					},
				},
			},
			"endpoint_rules": schema.ListNestedAttribute{
				MarkdownDescription: "Endpoint rules for the connector",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"order": schema.Int64Attribute{
							MarkdownDescription: "Order of the endpoint rule",
							Required:            true,
						},
						"behavior": schema.StringAttribute{
							MarkdownDescription: "Behavior of the endpoint rule (\"Allow\", \"Deny\")",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("Allow", "Deny"),
							},
						},
						"endpoint": schema.StringAttribute{
							MarkdownDescription: "Endpoint of the endpoint rule",
							Required:            true,
						},
					},
				},
			},
		},
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches the list of Data Loss Prevention Policies in a Power Platform tenant. See [Manage data loss prevention policies](https://learn.microsoft.com/power-platform/admin/prevent-data-loss) for more information.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: false,
				Update: false,
				Delete: false,
				Read:   false,
			}),
			"policies": schema.ListNestedAttribute{
				MarkdownDescription: "List of Data Loss Prevention Policies",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Unique name of the policy",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "Display name of the policy",
							Computed:            true,
						},
						"created_by": schema.StringAttribute{
							MarkdownDescription: "User who created the policy",
							Computed:            true,
						},
						"created_time": schema.StringAttribute{
							MarkdownDescription: "Time when the policy was created",
							Computed:            true,
						},
						"last_modified_by": schema.StringAttribute{
							MarkdownDescription: "User who last modified the policy",
							Computed:            true,
						},
						"last_modified_time": schema.StringAttribute{
							MarkdownDescription: "Time when the policy was last modified",
							Computed:            true,
						},
						"environment_type": schema.StringAttribute{
							MarkdownDescription: "Default environment handling for the policy (\"AllEnvironments\", \"ExceptEnvironments\", \"OnlyEnvironments\")",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("AllEnvironments", "ExceptEnvironments", "OnlyEnvironments"),
							},
						},
						"default_connectors_classification": schema.StringAttribute{
							MarkdownDescription: "Default classification for connectors (\"General\", \"Confidential\", \"Blocked\")",
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("General", "Confidential", "Blocked"),
							},
						},
						"environments": schema.SetAttribute{
							Description:         "Environment to which the policy is applied",
							MarkdownDescription: "Environment to which the policy is applied",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"business_connectors": schema.SetNestedAttribute{
							MarkdownDescription: "Connectors for sensitive data",
							Computed:            true,
							NestedObject:        connectorSchema,
						},
						"non_business_connectors": schema.SetNestedAttribute{
							MarkdownDescription: "Connectors for non-sensitive data",
							Computed:            true,
							NestedObject:        connectorSchema,
						},
						"blocked_connectors": schema.SetNestedAttribute{
							MarkdownDescription: "Blocked connectors can’t be used where this policy is applied.",
							Computed:            true,
							NestedObject:        connectorSchema,
						},
						"custom_connectors_patterns": schema.SetNestedAttribute{
							MarkdownDescription: "Custom connectors patterns",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"order": schema.Int64Attribute{
										MarkdownDescription: "Order of the connector",
										Computed:            true,
									},
									"host_url_pattern": schema.StringAttribute{
										MarkdownDescription: "Pattern of the connector",
										Computed:            true,
									},
									"data_group": schema.StringAttribute{
										MarkdownDescription: "Data group of the connector (\"Business\", \"NonBusiness\", \"Blocked\", \"Ignore\")",
										Computed:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("Business", "NonBusiness", "Blocked", "Ignore"),
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

func (d *DataLossPreventionPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
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

	d.DlpPolicyClient = newDlpPolicyClient(client)
}

func (d *DataLossPreventionPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var state policiesListDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE POLICIES START: %s_%s", d.ProviderTypeName, d.TypeName))

	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policies, err := d.DlpPolicyClient.GetPolicies(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s_%s", d.ProviderTypeName, d.TypeName), err.Error())
		return
	}

	for _, policy := range policies {
		policyModel := dataLossPreventionPolicyDatasourceModel{}
		policyModel.Id = types.StringValue(policy.Name)
		policyModel.DefaultConnectorsClassification = types.StringValue(policy.DefaultConnectorsClassification)
		policyModel.DisplayName = types.StringValue(policy.DisplayName)
		policyModel.CreatedBy = types.StringValue(policy.CreatedBy)
		policyModel.CreatedTime = types.StringValue(policy.CreatedTime)
		policyModel.LastModifiedBy = types.StringValue(policy.LastModifiedBy)
		policyModel.LastModifiedTime = types.StringValue(policy.LastModifiedTime)
		policyModel.EnvironmentType = types.StringValue(policy.EnvironmentType)
		policyModel.Environments = convertToAttrValueEnvironments(policy.Environments)
		policyModel.CustomConnectorsPatterns = convertToAttrValueCustomConnectorUrlPatternsDefinition(policy.CustomConnectorUrlPatternsDefinition)
		policyModel.BusinessGeneralConnectors = convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
		policyModel.NonBusinessConfidentialConnectors = convertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
		policyModel.BlockedConnectors = convertToAttrValueConnectorsGroup("Blocked", policy.ConnectorGroups)
		state.Policies = append(state.Policies, policyModel)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
