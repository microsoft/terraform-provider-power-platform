package powerplatform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	powerplatform "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

var _ resource.Resource = &DataLossPreventionPolicyResource{}
var _ resource.ResourceWithImportState = &DataLossPreventionPolicyResource{}

func NewDataLossPreventionPolicyResource() resource.Resource {
	return &DataLossPreventionPolicyResource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_data_loss_prevention_policy",
	}
}

type DataLossPreventionPolicyResource struct {
	PowerPlatformApiClient powerplatform.ClientInterface
	ProviderTypeName       string
	TypeName               string
}

type DataLossPreventionPolicyResourceModel struct {
	Name                            types.String                                           `tfsdk:"name"`
	DisplayName                     types.String                                           `tfsdk:"display_name"`
	CreatedBy                       types.String                                           `tfsdk:"created_by"`
	CreatedTime                     types.String                                           `tfsdk:"created_time"`
	LastModifiedBy                  types.String                                           `tfsdk:"last_modified_by"`
	LastModifiedTime                types.String                                           `tfsdk:"last_modified_time"`
	ETag                            types.String                                           `tfsdk:"e_tag"`
	EnvironmentType                 types.String                                           `tfsdk:"environment_type"`
	DefaultConnectorsClassification types.String                                           `tfsdk:"default_connectors_classification"`
	Environments                    []DataLossPreventionPolicyResourceEnvironmentsModel    `tfsdk:"environments"`
	ConnectorGroups                 []DataLossPreventionPolicyResourceConnectorGroupsModel `tfsdk:"connector_groups"`
}

type DataLossPreventionPolicyResourceEnvironmentsModel struct {
	EnvironmentName types.String `tfsdk:"environment_name"`
}

type DataLossPreventionPolicyResourceConnectorGroupsModel struct {
	Classification types.String                                      `tfsdk:"classification"`
	Connectors     []DataLossPreventionPolicyResourceConnectorsModel `tfsdk:"connectors"`
}

type DataLossPreventionPolicyResourceConnectorsModel struct {
	Name types.String `tfsdk:"name"`
	Id   types.String `tfsdk:"id"`
}

func (r *DataLossPreventionPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

func (r *DataLossPreventionPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data Loss Prevention Policy",
		Description:         "Data Loss Prevention Policy",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Unique name of the policy",
				Description:         "Unique name of the policy",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the policy",
				Description:         "The display name of the policy",
				Required:            true,
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "User who created the policy",
				Description:         "User who created the policy",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_time": schema.StringAttribute{
				MarkdownDescription: "Time when the policy was created",
				Description:         "Time when the policy was created",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
			"e_tag": schema.StringAttribute{
				MarkdownDescription: "ETag of the policy",
				Description:         "ETag of the policy",
				Computed:            true,
			},
			"environment_type": schema.StringAttribute{
				MarkdownDescription: "Default environment handling for the policy (\"AllEnvironments\", \"ExceptEnvironments\", \"OnlyEnvironments\")",
				Description:         "Default environment handling for the policy (\"AllEnvironments\", \"ExceptEnvironments\", \"OnlyEnvironments\")",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("AllEnvironments", "ExceptEnvironments", "OnlyEnvironments"),
				},
			},
			"default_connectors_classification": schema.StringAttribute{
				MarkdownDescription: "Default classification for connectors (\"General\", \"Confidential\", \"Blocked\")",
				Description:         "Default classification for connectors (\"General\", \"Confidential\", \"Blocked\")",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("General", "Confidential", "Blocked"),
				},
			},
			"environments": schema.ListNestedAttribute{
				MarkdownDescription: "Environment to which the policy is applied",
				Description:         "Environment to which the policy is applied",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"environment_name": schema.StringAttribute{
							MarkdownDescription: "Name of the environment",
							Description:         "Name of the environment",
							Required:            true,
							//BUG if we don't set this to reqired, it will produce inconsistent state
							// When applying changes to powerplatform_data_loss_prevention_policy.my_policy,
							// │ provider "provider[\"github.com/microsoft/terraform-provider-power-platform\"]" produced an unexpected
							// │ new value: .environments: new element 1 has appeared.
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"connector_groups": schema.ListNestedAttribute{
				MarkdownDescription: "Connector groups to which the policy is applied",
				Description:         "Connector groups to which the policy is applied",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"classification": schema.StringAttribute{
							MarkdownDescription: "Classification of the connector group (\"General\", \"Confidential\", \"Blocked\")",
							Description:         "Classification of the connector group (\"General\", \"Confidential\", \"Blocked\")",
							Required:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Validators: []validator.String{
								stringvalidator.OneOf("General", "Confidential", "Blocked"),
							},
						},
						"connectors": schema.ListNestedAttribute{
							MarkdownDescription: "Connectors in the group",
							Description:         "Connectors in the group",
							Required:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										MarkdownDescription: "Name of the connector",
										Description:         "Name of the connector",
										Required:            true,
										// PlanModifiers: []planmodifier.String{
										// 	stringplanmodifier.RequiresReplace(),
										// },
									},
									"id": schema.StringAttribute{
										MarkdownDescription: "ID of the connector",
										Description:         "ID of the connector",
										Required:            true,
										// PlanModifiers: []planmodifier.String{
										// 	stringplanmodifier.RequiresReplace(),
										// },
									},
								},
							},
							Validators: []validator.List{
								listvalidator.UniqueValues(),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeBetween(1, 3),
					listvalidator.UniqueValues(),
				},
			},
		},
	}
}

func (r *DataLossPreventionPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*powerplatform.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.PowerPlatformApiClient = client
}

func (r *DataLossPreventionPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *DataLossPreventionPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	//TODO handle  404 NOT FOUND responses
	policy, err := r.PowerPlatformApiClient.GetPolicy(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.ProviderTypeName), err.Error())
		return
	}

	if policy == nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.ProviderTypeName), "Policy not found")
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(policy.Name)
	state.DefaultConnectorsClassification = types.StringValue(policy.DefaultConnectorsClassification)
	state.DisplayName = types.StringValue(policy.DisplayName)
	state.CreatedBy = types.StringValue(policy.CreatedBy)
	state.CreatedTime = types.StringValue(policy.CreatedTime)
	state.LastModifiedBy = types.StringValue(policy.LastModifiedBy)
	state.LastModifiedTime = types.StringValue(policy.LastModifiedTime)
	state.ETag = types.StringValue(policy.ETag)
	state.EnvironmentType = types.StringValue(policy.EnvironmentType)

	for index, env := range policy.Environments {
		state.Environments[index] = DataLossPreventionPolicyResourceEnvironmentsModel{
			EnvironmentName: types.StringValue(env.Name),
		}
	}

	for indexConnGroup, connectorGroup := range policy.ConnectorGroups {
		if len(policy.ConnectorGroups) > 0 {

			state.ConnectorGroups[indexConnGroup] = DataLossPreventionPolicyResourceConnectorGroupsModel{
				Classification: types.StringValue(connectorGroup.Classification),
				Connectors:     make([]DataLossPreventionPolicyResourceConnectorsModel, 0),
			}

			for _, connector := range policy.ConnectorGroups[indexConnGroup].Connectors {
				if len(policy.ConnectorGroups[indexConnGroup].Connectors) > 0 {

					state.ConnectorGroups[indexConnGroup].Connectors = append(state.ConnectorGroups[indexConnGroup].Connectors, DataLossPreventionPolicyResourceConnectorsModel{
						Id:   types.StringValue(connector.Id),
						Name: types.StringValue(connector.Name),
					})
				}
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE END: %s", r.ProviderTypeName))

}

func (r *DataLossPreventionPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//BUG when adding new policy, api returns set of default connectors, this is sees by provider as inconsistency and it fails
	var plan *DataLossPreventionPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	policyToCreate := powerplatform.DlpPolicy{
		DefaultConnectorsClassification: plan.DefaultConnectorsClassification.ValueString(),
		DisplayName:                     plan.DisplayName.ValueString(),
		EnvironmentType:                 plan.EnvironmentType.ValueString(),
		Environments:                    []powerplatform.DlpEnvironment{},
		ConnectorGroups:                 []powerplatform.DlpConnectorGroups{},
	}

	policyToCreate.Environments = make([]powerplatform.DlpEnvironment, 0)
	for _, environment := range plan.Environments {
		policyToCreate.Environments = append(policyToCreate.Environments, powerplatform.DlpEnvironment{
			Name: environment.EnvironmentName.ValueString(),
		})
	}

	policyToCreate.ConnectorGroups = make([]powerplatform.DlpConnectorGroups, 0)
	for index, connectorGroup := range plan.ConnectorGroups {
		policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, powerplatform.DlpConnectorGroups{
			Classification: connectorGroup.Classification.ValueString(),
			Connectors:     make([]powerplatform.DlpConnector, 0),
		})

		for _, connector := range connectorGroup.Connectors {
			policyToCreate.ConnectorGroups[index].Connectors = append(policyToCreate.ConnectorGroups[index].Connectors, powerplatform.DlpConnector{
				Id:   connector.Id.ValueString(),
				Name: connector.Name.ValueString(),
			})
		}
	}

	policy, err := r.PowerPlatformApiClient.CreatePolicy(ctx, policyToCreate)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.ProviderTypeName), err.Error())
		return
	}

	plan.Name = types.StringValue(policy.Name)
	plan.DefaultConnectorsClassification = types.StringValue(policy.DefaultConnectorsClassification)
	plan.DisplayName = types.StringValue(policy.DisplayName)
	plan.CreatedBy = types.StringValue(policy.CreatedBy)
	plan.CreatedTime = types.StringValue(policy.CreatedTime)
	plan.LastModifiedBy = types.StringValue(policy.LastModifiedBy)
	plan.LastModifiedTime = types.StringValue(policy.LastModifiedTime)
	plan.ETag = types.StringValue(policy.ETag)
	plan.EnvironmentType = types.StringValue(policy.EnvironmentType)

	for index, env := range policy.Environments {
		plan.Environments[index] = DataLossPreventionPolicyResourceEnvironmentsModel{
			EnvironmentName: types.StringValue(env.Name),
		}
	}

	for indexConnGroup, connectorGroup := range policy.ConnectorGroups {
		if len(policy.ConnectorGroups) > 0 {

			plan.ConnectorGroups[indexConnGroup] = DataLossPreventionPolicyResourceConnectorGroupsModel{
				Classification: types.StringValue(connectorGroup.Classification),
				Connectors:     make([]DataLossPreventionPolicyResourceConnectorsModel, 0),
			}

			for _, connector := range policy.ConnectorGroups[indexConnGroup].Connectors {
				if len(policy.ConnectorGroups[indexConnGroup].Connectors) > 0 {

					plan.ConnectorGroups[indexConnGroup].Connectors = append(plan.ConnectorGroups[indexConnGroup].Connectors, DataLossPreventionPolicyResourceConnectorsModel{
						Id:   types.StringValue(connector.Id),
						Name: types.StringValue(connector.Name),
					})
				}
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *DataLossPreventionPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	//BUG when updating data policy in UI, all the connectors are added and they provides sees that as a change to delete/add the policy
	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE START: %s", r.ProviderTypeName))

	var plan *DataLossPreventionPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state *DataLossPreventionPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	policyToUpdate := powerplatform.DlpPolicy{
		Name:                            plan.Name.ValueString(),
		DisplayName:                     plan.DisplayName.ValueString(),
		EnvironmentType:                 plan.EnvironmentType.ValueString(),
		DefaultConnectorsClassification: plan.DefaultConnectorsClassification.ValueString(),
		Environments:                    []powerplatform.DlpEnvironment{},
		ConnectorGroups:                 []powerplatform.DlpConnectorGroups{},
	}
	policyToUpdate.Environments = make([]powerplatform.DlpEnvironment, 0)
	policyToUpdate.ConnectorGroups = make([]powerplatform.DlpConnectorGroups, 0)

	policyToUpdate.Environments = make([]powerplatform.DlpEnvironment, 0)
	for _, environment := range plan.Environments {
		policyToUpdate.Environments = append(policyToUpdate.Environments, powerplatform.DlpEnvironment{
			Name: environment.EnvironmentName.ValueString(),
		})
	}

	policyToUpdate.ConnectorGroups = make([]powerplatform.DlpConnectorGroups, 0)
	for index, connectorGroup := range plan.ConnectorGroups {
		policyToUpdate.ConnectorGroups = append(policyToUpdate.ConnectorGroups, powerplatform.DlpConnectorGroups{
			Classification: connectorGroup.Classification.ValueString(),
			Connectors:     make([]powerplatform.DlpConnector, 0),
		})

		for _, connector := range connectorGroup.Connectors {
			policyToUpdate.ConnectorGroups[index].Connectors = append(policyToUpdate.ConnectorGroups[index].Connectors, powerplatform.DlpConnector{
				Id:   connector.Id.ValueString(),
				Name: connector.Name.ValueString(),
			})
		}
	}

	policy, err := r.PowerPlatformApiClient.UpdatePolicy(ctx, plan.Name.ValueString(), policyToUpdate)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.ProviderTypeName), err.Error())
		return
	}

	plan.Name = types.StringValue(policy.Name)
	plan.DefaultConnectorsClassification = types.StringValue(policy.DefaultConnectorsClassification)
	plan.DisplayName = types.StringValue(policy.DisplayName)
	plan.CreatedBy = types.StringValue(policy.CreatedBy)
	plan.CreatedTime = types.StringValue(policy.CreatedTime)
	plan.LastModifiedBy = types.StringValue(policy.LastModifiedBy)
	plan.LastModifiedTime = types.StringValue(policy.LastModifiedTime)
	plan.ETag = types.StringValue(policy.ETag)
	plan.EnvironmentType = types.StringValue(policy.EnvironmentType)

	for index, env := range policy.Environments {
		plan.Environments[index] = DataLossPreventionPolicyResourceEnvironmentsModel{
			EnvironmentName: types.StringValue(env.Name),
		}
	}

	for indexConnGroup, connectorGroup := range policy.ConnectorGroups {
		if len(policy.ConnectorGroups) > 0 {

			plan.ConnectorGroups[indexConnGroup] = DataLossPreventionPolicyResourceConnectorGroupsModel{
				Classification: types.StringValue(connectorGroup.Classification),
				Connectors:     make([]DataLossPreventionPolicyResourceConnectorsModel, 0),
			}

			for _, connector := range policy.ConnectorGroups[indexConnGroup].Connectors {
				if len(policy.ConnectorGroups[indexConnGroup].Connectors) > 0 {

					plan.ConnectorGroups[indexConnGroup].Connectors = append(plan.ConnectorGroups[indexConnGroup].Connectors, DataLossPreventionPolicyResourceConnectorsModel{
						Id:   types.StringValue(connector.Id),
						Name: types.StringValue(connector.Name),
					})
				}
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *DataLossPreventionPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE START: %s", r.ProviderTypeName))

	var state *DataLossPreventionPolicyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.PowerPlatformApiClient.DeletePolicy(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.ProviderTypeName), err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *DataLossPreventionPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
