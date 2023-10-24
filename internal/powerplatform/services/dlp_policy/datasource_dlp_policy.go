package powerplatform

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	clients "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/clients"
)

var (
	_ datasource.DataSource              = &DataLossPreventionPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &DataLossPreventionPolicyDataSource{}
)

func NewDataLossPreventionPolicyDataSource() datasource.DataSource {
	return &DataLossPreventionPolicyDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_data_loss_prevention_policies",
	}
}

type DataLossPreventionPolicyDataSource struct {
	DlpPolicyClient  DlpPolicyClient
	ProviderTypeName string
	TypeName         string
}

func (d *DataLossPreventionPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *DataLossPreventionPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	connectorSchema := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the connector",
				Description:         "ID of the connector",
				Optional:            true,
			},
			"default_action_rule_behavior": schema.StringAttribute{
				MarkdownDescription: "Default action rule behavior for the connector (\"Allow\", \"Block\")",
				Description:         "Default action rule behavior for the connector (\"Allow\", \"Block\")",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Allow", "Block", ""),
				},
			},
			"action_rules": schema.ListNestedAttribute{
				MarkdownDescription: "Action rules for the connector",
				Description:         "Action rules for the connector",
				Optional:            true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"action_id": schema.StringAttribute{
							MarkdownDescription: "ID of the action rule",
							Description:         "ID of the action rule",
							Required:            true,
						},
						"behavior": schema.StringAttribute{
							MarkdownDescription: "Behavior of the action rule (\"Allow\", \"Block\")",
							Description:         "Behavior of the action rule (\"Allow\", \"Block\")",
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
				Description:         "Endpoint rules for the connector",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"order": schema.Int64Attribute{
							MarkdownDescription: "Order of the endpoint rule",
							Description:         "Order of the endpoint rule",
							Required:            true,
						},
						"behavior": schema.StringAttribute{
							MarkdownDescription: "Behavior of the endpoint rule (\"Allow\", \"Deny\")",
							Description:         "Behavior of the endpoint rule (\"Allow\", \"Deny\")",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("Allow", "Deny"),
							},
						},
						"endpoint": schema.StringAttribute{
							MarkdownDescription: "Endpoint of the endpoint rule",
							Description:         "Endpoint of the endpoint rule",
							Required:            true,
						},
					},
				},
			},
		},
	}

	resp.Schema = schema.Schema{
		Description:         "Fetches the list of available connectors in a Power Platform tenant",
		MarkdownDescription: "Fetches the list of available connectors in a Power Platform tenant",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"policies": schema.ListNestedAttribute{
				Description:         "List of Data Loss Prevention Policies",
				MarkdownDescription: "List of Data Loss Prevention Policies",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Unique name of the policy",
							Description:         "Unique name of the policy",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "Display name of the policy",
							Description:         "The display name of the policy",
							Computed:            true,
						},
						"created_by": schema.StringAttribute{
							MarkdownDescription: "User who created the policy",
							Description:         "User who created the policy",
							Computed:            true,
						},
						"created_time": schema.StringAttribute{
							MarkdownDescription: "Time when the policy was created",
							Description:         "Time when the policy was created",
							Computed:            true,
						},
						"last_modified_by": schema.StringAttribute{
							MarkdownDescription: "User who last modified the policy",
							Description:         "User who last modified the policy",
							Computed:            true,
						},
						"last_modified_time": schema.StringAttribute{
							MarkdownDescription: "Time when the policy was last modified",
							Description:         "Time when the policy was last modified",
							Computed:            true,
						},
						"environment_type": schema.StringAttribute{
							MarkdownDescription: "Default environment handling for the policy (\"AllEnvironments\", \"ExceptEnvironments\", \"OnlyEnvironments\")",
							Description:         "Default environment handling for the policy (\"AllEnvironments\", \"ExceptEnvironments\", \"OnlyEnvironments\")",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("AllEnvironments", "ExceptEnvironments", "OnlyEnvironments"),
							},
						},
						"default_connectors_classification": schema.StringAttribute{
							MarkdownDescription: "Default classification for connectors (\"General\", \"Confidential\", \"Blocked\")",
							Description:         "Default classification for connectors (\"General\", \"Confidential\", \"Blocked\")",
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("General", "Confidential", "Blocked"),
							},
						},
						"environments": schema.SetNestedAttribute{
							MarkdownDescription: "Environment to which the policy is applied",
							Description:         "Environment to which the policy is applied",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										MarkdownDescription: "Unique Identifier of the environment",
										Description:         "Unique Identifier of the environment",
										Computed:            true,
									},
								},
							},
						},
						"business_connectors": schema.SetNestedAttribute{
							MarkdownDescription: "Connectors for sensitive data",
							Description:         "Connectors for sensitive data",
							Computed:            true,
							NestedObject:        connectorSchema,
						},
						"non_business_connectors": schema.SetNestedAttribute{
							MarkdownDescription: "Connectors for non-sensitive data",
							Description:         "Connectors for non-sensitive data",
							Computed:            true,
							NestedObject:        connectorSchema,
						},
						"blocked_connectors": schema.SetNestedAttribute{
							MarkdownDescription: "Blocked connectors can’t be used where this policy is applied.",
							Description:         "Blocked connectors can’t be used where this policy is applied.",
							Computed:            true,
							NestedObject:        connectorSchema,
						},
						"custom_connectors_patterns": schema.SetNestedAttribute{
							MarkdownDescription: "Custom connectors patterns",
							Description:         "Custom connectors patterns",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"order": schema.Int64Attribute{
										MarkdownDescription: "Order of the connector",
										Description:         "Order of the connector",
										Computed:            true,
									},
									"host_url_pattern": schema.StringAttribute{
										MarkdownDescription: "Pattern of the connector",
										Description:         "Pattern of the connector",
										Computed:            true,
									},
									"data_group": schema.StringAttribute{
										MarkdownDescription: "Data group of the connector (\"Business\", \"NonBusiness\", \"Blocked\", \"Ignore\")",
										Description:         "Data group of the connector (\"Business\", \"NonBusiness\", \"Blocked\", \"Ignore\")",
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

func (d *DataLossPreventionPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := req.ProviderData.(*clients.ProviderClient).BapiApi.Client

	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.DlpPolicyClient = NewDlpPolicyClient(client)
}

func (d *DataLossPreventionPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state PoliciesListDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE POLICIES START: %s_%s", d.ProviderTypeName, d.TypeName))

	policies, err := d.DlpPolicyClient.GetPolicies(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s_%s", d.ProviderTypeName, d.TypeName), err.Error())
		return
	}

	state.Id = types.StringValue(fmt.Sprint((time.Now().Unix())))

	for _, policy := range policies {

		policyModel := DataLossPreventionPolicyResourceModel{}
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
		policyModel.BusinessGeneralConnectors = convertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
		policyModel.NonBusinessConfidentialConnectors = convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
		policyModel.BlockedConnectors = convertToAttrValueConnectorsGroup("Blocked", policy.ConnectorGroups)
		state.Policies = append(state.Policies, policyModel)
	}

	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE POLICIES END: %s_%s", d.ProviderTypeName, d.TypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
