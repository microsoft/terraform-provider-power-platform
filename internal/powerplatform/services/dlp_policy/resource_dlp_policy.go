package powerplatform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	clients "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/clients"
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
	DlpPolicyClient  DlpPolicyClient
	ProviderTypeName string
	TypeName         string
}

func (r *DataLossPreventionPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

func (r *DataLossPreventionPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

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
		MarkdownDescription: "Data Loss Prevention Policy",
		Description:         "Data Loss Prevention Policy",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
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
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("General", "Confidential", "Blocked"),
				},
			},
			"environments": schema.SetNestedAttribute{
				MarkdownDescription: "Environment to which the policy is applied",
				Description:         "Environment to which the policy is applied",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Unique Identifier of the environment",
							Description:         "Unique Identifier of the environment",
							Required:            true,
						},
					},
				},
			},
			"business_connectors": schema.SetNestedAttribute{
				MarkdownDescription: "Connectors for sensitive data",
				Description:         "Connectors for sensitive data",
				Required:            true,
				NestedObject:        connectorSchema,
			},
			"non_business_connectors": schema.SetNestedAttribute{
				MarkdownDescription: "Connectors for non-sensitive data",
				Description:         "Connectors for non-sensitive data",
				Required:            true,
				NestedObject:        connectorSchema,
			},
			"blocked_connectors": schema.SetNestedAttribute{
				MarkdownDescription: "Blocked connectors can’t be used where this policy is applied.",
				Description:         "Blocked connectors can’t be used where this policy is applied.",
				Required:            true,
				NestedObject:        connectorSchema,
			},
			"custom_connectors_patterns": schema.SetNestedAttribute{
				MarkdownDescription: "Custom connectors patterns",
				Description:         "Custom connectors patterns",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"order": schema.Int64Attribute{
							MarkdownDescription: "Order of the connector",
							Description:         "Order of the connector",
							Required:            true,
						},
						"host_url_pattern": schema.StringAttribute{
							MarkdownDescription: "Pattern of the connector",
							Description:         "Pattern of the connector",
							Required:            true,
						},
						"data_group": schema.StringAttribute{
							MarkdownDescription: "Data group of the connector (\"Business\", \"NonBusiness\", \"Blocked\", \"Ignore\")",
							Description:         "Data group of the connector (\"Business\", \"NonBusiness\", \"Blocked\", \"Ignore\")",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("Business", "NonBusiness", "Blocked", "Ignore"),
							},
						},
					},
				},
			},
		},
	}
}

func (r *DataLossPreventionPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.DlpPolicyClient = NewDlpPolicyClient(client)
}

func (r *DataLossPreventionPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *DataLossPreventionPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := r.DlpPolicyClient.GetPolicy(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.TypeName), err.Error())
		resp.State.RemoveResource(ctx)
		return
	}

	state.Id = types.StringValue(policy.Name)
	state.DefaultConnectorsClassification = types.StringValue(policy.DefaultConnectorsClassification)
	state.DisplayName = types.StringValue(policy.DisplayName)
	state.CreatedBy = types.StringValue(policy.CreatedBy)
	state.CreatedTime = types.StringValue(policy.CreatedTime)
	state.LastModifiedBy = types.StringValue(policy.LastModifiedBy)
	state.LastModifiedTime = types.StringValue(policy.LastModifiedTime)
	state.EnvironmentType = types.StringValue(policy.EnvironmentType)
	state.Environments = convertToAttrValueEnvironments(policy.Environments)
	state.CustomConnectorsPatterns = convertToAttrValueCustomConnectorUrlPatternsDefinition(policy.CustomConnectorUrlPatternsDefinition)
	state.BusinessGeneralConnectors = convertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
	state.NonBusinessConfidentialConnectors = convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
	state.BlockedConnectors = convertToAttrValueConnectorsGroup("Blocked", policy.ConnectorGroups)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE END: %s", r.TypeName))
}

func (r *DataLossPreventionPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *DataLossPreventionPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE START: %s", r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	policyToCreate := DlpPolicyModelDto{
		DefaultConnectorsClassification:      plan.DefaultConnectorsClassification.ValueString(),
		DisplayName:                          plan.DisplayName.ValueString(),
		EnvironmentType:                      plan.EnvironmentType.ValueString(),
		Environments:                         []DlpEnvironmentDto{},
		ConnectorGroups:                      []DlpConnectorGroupsModelDto{},
		CustomConnectorUrlPatternsDefinition: []DlpConnectorUrlPatternsDefinitionDto{},
	}

	policyToCreate.Environments = convertToDlpEnvironment(ctx, resp.Diagnostics, plan.Environments)
	policyToCreate.CustomConnectorUrlPatternsDefinition = convertToDlpCustomConnectorUrlPatternsDefinition(ctx, resp.Diagnostics, plan.CustomConnectorsPatterns)
	policyToCreate.ConnectorGroups = make([]DlpConnectorGroupsModelDto, 0)
	policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, convertToDlpConnectorGroup(ctx, resp.Diagnostics, "General", plan.BusinessGeneralConnectors))
	policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, convertToDlpConnectorGroup(ctx, resp.Diagnostics, "Confidential", plan.NonBusinessConfidentialConnectors))
	policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, convertToDlpConnectorGroup(ctx, resp.Diagnostics, "Blocked", plan.BlockedConnectors))

	policy, err_client := r.DlpPolicyClient.CreatePolicy(ctx, policyToCreate)
	if err_client != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s_%s", r.ProviderTypeName, r.TypeName), err_client.Error())
		return
	}

	plan.Id = types.StringValue(policy.Name)
	plan.DefaultConnectorsClassification = types.StringValue(policy.DefaultConnectorsClassification)
	plan.DisplayName = types.StringValue(policy.DisplayName)
	plan.CreatedBy = types.StringValue(policy.CreatedBy)
	plan.CreatedTime = types.StringValue(policy.CreatedTime)
	plan.LastModifiedBy = types.StringValue(policy.LastModifiedBy)
	plan.LastModifiedTime = types.StringValue(policy.LastModifiedTime)
	plan.EnvironmentType = types.StringValue(policy.EnvironmentType)
	plan.Environments = convertToAttrValueEnvironments(policy.Environments)
	plan.CustomConnectorsPatterns = convertToAttrValueCustomConnectorUrlPatternsDefinition(policy.CustomConnectorUrlPatternsDefinition)
	plan.BusinessGeneralConnectors = convertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
	plan.NonBusinessConfidentialConnectors = convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
	plan.BlockedConnectors = convertToAttrValueConnectorsGroup("Blocked", policy.ConnectorGroups)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE END: %s", r.TypeName))
}

func (r *DataLossPreventionPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE START: %s", r.TypeName))

	var plan *DataLossPreventionPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state *DataLossPreventionPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	policyToUpdate := DlpPolicyModelDto{
		Name:                            plan.Id.ValueString(),
		DisplayName:                     plan.DisplayName.ValueString(),
		EnvironmentType:                 plan.EnvironmentType.ValueString(),
		DefaultConnectorsClassification: plan.DefaultConnectorsClassification.ValueString(),
		Environments:                    []DlpEnvironmentDto{},
		ConnectorGroups:                 []DlpConnectorGroupsModelDto{},
	}

	policyToUpdate.Environments = convertToDlpEnvironment(ctx, resp.Diagnostics, plan.Environments)
	policyToUpdate.CustomConnectorUrlPatternsDefinition = convertToDlpCustomConnectorUrlPatternsDefinition(ctx, resp.Diagnostics, plan.CustomConnectorsPatterns)
	policyToUpdate.ConnectorGroups = make([]DlpConnectorGroupsModelDto, 0)
	policyToUpdate.ConnectorGroups = append(policyToUpdate.ConnectorGroups, convertToDlpConnectorGroup(ctx, resp.Diagnostics, "General", plan.BusinessGeneralConnectors))
	policyToUpdate.ConnectorGroups = append(policyToUpdate.ConnectorGroups, convertToDlpConnectorGroup(ctx, resp.Diagnostics, "Confidential", plan.NonBusinessConfidentialConnectors))
	policyToUpdate.ConnectorGroups = append(policyToUpdate.ConnectorGroups, convertToDlpConnectorGroup(ctx, resp.Diagnostics, "Blocked", plan.BlockedConnectors))

	policy, err_client := r.DlpPolicyClient.UpdatePolicy(ctx, plan.Id.ValueString(), policyToUpdate)
	if err_client != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.TypeName), err_client.Error())
		return
	}

	plan.Id = types.StringValue(policy.Name)
	plan.DefaultConnectorsClassification = types.StringValue(policy.DefaultConnectorsClassification)
	plan.DisplayName = types.StringValue(policy.DisplayName)
	plan.CreatedBy = types.StringValue(policy.CreatedBy)
	plan.CreatedTime = types.StringValue(policy.CreatedTime)
	plan.LastModifiedBy = types.StringValue(policy.LastModifiedBy)
	plan.LastModifiedTime = types.StringValue(policy.LastModifiedTime)
	plan.EnvironmentType = types.StringValue(policy.EnvironmentType)
	plan.Environments = convertToAttrValueEnvironments(policy.Environments)
	plan.CustomConnectorsPatterns = convertToAttrValueCustomConnectorUrlPatternsDefinition(policy.CustomConnectorUrlPatternsDefinition)
	plan.BusinessGeneralConnectors = convertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
	plan.NonBusinessConfidentialConnectors = convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
	plan.BlockedConnectors = convertToAttrValueConnectorsGroup("Blocked", policy.ConnectorGroups)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE END: %s", r.TypeName))
}

func (r *DataLossPreventionPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE START: %s", r.TypeName))

	var state *DataLossPreventionPolicyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.DlpPolicyClient.DeletePolicy(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.TypeName), err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE END: %s", r.TypeName))
}

func (r *DataLossPreventionPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("Id"), req, resp)
}
