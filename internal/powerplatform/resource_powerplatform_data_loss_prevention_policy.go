package powerplatform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	bapi "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api/bapi"
	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/models"
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
	BapiApiClient    bapi.BapiClientInterface
	ProviderTypeName string
	TypeName         string
}

type DataLossPreventionPolicyResourceModel struct {
	Id                                types.String `tfsdk:"id"`
	DisplayName                       types.String `tfsdk:"display_name"`
	DefaultConnectorsClassification   types.String `tfsdk:"default_connectors_classification"`
	EnvironmentType                   types.String `tfsdk:"environment_type"`
	CreatedBy                         types.String `tfsdk:"created_by"`
	CreatedTime                       types.String `tfsdk:"created_time"`
	LastModifiedBy                    types.String `tfsdk:"last_modified_by"`
	LastModifiedTime                  types.String `tfsdk:"last_modified_time"`
	Environments                      types.Set    `tfsdk:"environments"`
	NonBusinessConfidentialConnectors types.Set    `tfsdk:"non_business_connectors"`
	BusinessGeneralConnectors         types.Set    `tfsdk:"business_connectors"`
	BlockedConnectors                 types.Set    `tfsdk:"blocked_connectors"`
	CustomConnectorsPatterns          types.Set    `tfsdk:"custom_connectors_patterns"`
}

type DataLossPreventionPolicyResourceCustomConnectorPattern struct {
	Order          types.Int64  `tfsdk:"order"`
	HostUrlPattern types.String `tfsdk:"host_url_pattern"`
	DataGroup      types.String `tfsdk:"data_group"`
}

var customConnectorPatternSetObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"order":            types.Int64Type,
		"host_url_pattern": types.StringType,
		"data_group":       types.StringType,
	},
}

type DataLossPreventionPolicyResourceEnvironmentsModel struct {
	Name types.String `tfsdk:"name"`
}

var environmentSetObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name": types.StringType,
	},
}

type DataLossPreventionPolicyResourceConnectorModel struct {
	Id                        types.String                                                 `tfsdk:"id"`
	DefaultActionRuleBehavior types.String                                                 `tfsdk:"default_action_rule_behavior"`
	ActionRules               []DataLossPreventionPolicyResourceConnectorActionRuleModel   `tfsdk:"action_rules"`
	EndpointRules             []DataLossPreventionPolicyResourceConnectorEndpointRuleModel `tfsdk:"endpoint_rules"`
}

type DataLossPreventionPolicyResourceConnectorEndpointRuleModel struct {
	Order    types.Int64  `tfsdk:"order"`
	Behavior types.String `tfsdk:"behavior"`
	Endpoint types.String `tfsdk:"endpoint"`
}

type DataLossPreventionPolicyResourceConnectorActionRuleModel struct {
	ActionId types.String `tfsdk:"action_id"`
	Behavior types.String `tfsdk:"behavior"`
}

var connectorSetObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":                           types.StringType,
		"default_action_rule_behavior": types.StringType,
		"action_rules":                 types.ListType{ElemType: actionRuleListObjectType},
		"endpoint_rules":               types.ListType{ElemType: endpointRuleListObjectType},
	},
}

var endpointRuleListObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"order":    types.Int64Type,
		"behavior": types.StringType,
		"endpoint": types.StringType,
	},
}

var actionRuleListObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"action_id": types.StringType,
		"behavior":  types.StringType,
	},
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

	client, ok := req.ProviderData.(*PowerPlatformProvider).BapiApi.Client.(bapi.BapiClientInterface)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.BapiApiClient = client
}

func (r *DataLossPreventionPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *DataLossPreventionPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := r.BapiApiClient.GetPolicy(ctx, state.Id.ValueString())
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
	state.Environments = ConvertToAttrValueEnvironments(policy.Environments)
	state.CustomConnectorsPatterns = ConvertToAttrValueCustomConnectorUrlPatternsDefinition(policy.CustomConnectorUrlPatternsDefinition)
	state.BusinessGeneralConnectors = ConvertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
	state.NonBusinessConfidentialConnectors = ConvertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
	state.BlockedConnectors = ConvertToAttrValueConnectorsGroup("Blocked", policy.ConnectorGroups)

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

	policyToCreate := models.DlpPolicyModel{
		DefaultConnectorsClassification:      plan.DefaultConnectorsClassification.ValueString(),
		DisplayName:                          plan.DisplayName.ValueString(),
		EnvironmentType:                      plan.EnvironmentType.ValueString(),
		Environments:                         []models.DlpEnvironmentDto{},
		ConnectorGroups:                      []models.DlpConnectorGroupsModel{},
		CustomConnectorUrlPatternsDefinition: []models.DlpConnectorUrlPatternsDefinitionDto{},
	}

	policyToCreate.Environments = ConvertToDlpEnvironment(ctx, resp.Diagnostics, plan.Environments)
	policyToCreate.CustomConnectorUrlPatternsDefinition = ConvertToDlpCustomConnectorUrlPatternsDefinition(ctx, resp.Diagnostics, plan.CustomConnectorsPatterns)
	policyToCreate.ConnectorGroups = make([]models.DlpConnectorGroupsModel, 0)
	policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, ConvertToDlpConnectorGroup(ctx, resp.Diagnostics, "General", plan.BusinessGeneralConnectors))
	policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, ConvertToDlpConnectorGroup(ctx, resp.Diagnostics, "Confidential", plan.NonBusinessConfidentialConnectors))
	policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, ConvertToDlpConnectorGroup(ctx, resp.Diagnostics, "Blocked", plan.BlockedConnectors))

	policy, err_client := r.BapiApiClient.CreatePolicy(ctx, policyToCreate)
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
	plan.Environments = ConvertToAttrValueEnvironments(policy.Environments)
	plan.CustomConnectorsPatterns = ConvertToAttrValueCustomConnectorUrlPatternsDefinition(policy.CustomConnectorUrlPatternsDefinition)
	plan.BusinessGeneralConnectors = ConvertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
	plan.NonBusinessConfidentialConnectors = ConvertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
	plan.BlockedConnectors = ConvertToAttrValueConnectorsGroup("Blocked", policy.ConnectorGroups)

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

	policyToUpdate := models.DlpPolicyModel{
		Name:                            plan.Id.ValueString(),
		DisplayName:                     plan.DisplayName.ValueString(),
		EnvironmentType:                 plan.EnvironmentType.ValueString(),
		DefaultConnectorsClassification: plan.DefaultConnectorsClassification.ValueString(),
		Environments:                    []models.DlpEnvironmentDto{},
		ConnectorGroups:                 []models.DlpConnectorGroupsModel{},
	}

	policyToUpdate.Environments = ConvertToDlpEnvironment(ctx, resp.Diagnostics, plan.Environments)
	policyToUpdate.CustomConnectorUrlPatternsDefinition = ConvertToDlpCustomConnectorUrlPatternsDefinition(ctx, resp.Diagnostics, plan.CustomConnectorsPatterns)
	policyToUpdate.ConnectorGroups = make([]models.DlpConnectorGroupsModel, 0)
	policyToUpdate.ConnectorGroups = append(policyToUpdate.ConnectorGroups, ConvertToDlpConnectorGroup(ctx, resp.Diagnostics, "General", plan.BusinessGeneralConnectors))
	policyToUpdate.ConnectorGroups = append(policyToUpdate.ConnectorGroups, ConvertToDlpConnectorGroup(ctx, resp.Diagnostics, "Confidential", plan.NonBusinessConfidentialConnectors))
	policyToUpdate.ConnectorGroups = append(policyToUpdate.ConnectorGroups, ConvertToDlpConnectorGroup(ctx, resp.Diagnostics, "Blocked", plan.BlockedConnectors))

	policy, err_client := r.BapiApiClient.UpdatePolicy(ctx, plan.Id.ValueString(), policyToUpdate)
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
	plan.Environments = ConvertToAttrValueEnvironments(policy.Environments)
	plan.CustomConnectorsPatterns = ConvertToAttrValueCustomConnectorUrlPatternsDefinition(policy.CustomConnectorUrlPatternsDefinition)
	plan.BusinessGeneralConnectors = ConvertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
	plan.NonBusinessConfidentialConnectors = ConvertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
	plan.BlockedConnectors = ConvertToAttrValueConnectorsGroup("Blocked", policy.ConnectorGroups)

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

	err := r.BapiApiClient.DeletePolicy(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.TypeName), err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE END: %s", r.TypeName))
}

func (r *DataLossPreventionPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("Id"), req, resp)
}

func ConvertConnectorRuleClassificationValues(value string) string {
	if value == "Business" {
		return "General"
	} else if value == "NonBusiness" {
		return "Confidential"
	} else if value == "General" {
		return "Business"
	} else if value == "Confidential" {
		return "NonBusiness"
	} else {
		return value
	}
}

func ConvertToAttrValueConnectorsGroup(classification string, connectorsGroup []models.DlpConnectorGroupsModel) basetypes.SetValue {
	var connectorValues []attr.Value
	for _, conn := range connectorsGroup {
		if conn.Classification == classification {
			return types.SetValueMust(connectorSetObjectType, ConvertToAttrValueConnectors(conn, connectorValues))
		}
	}
	return types.SetValueMust(connectorSetObjectType, []attr.Value{})
}

func ConvertToAttrValueCustomConnectorUrlPatternsDefinition(urlPatterns []models.DlpConnectorUrlPatternsDefinitionDto) basetypes.SetValue {
	var connUrlPattern []attr.Value
	for _, connectorUrlPattern := range urlPatterns {
		for _, rules := range connectorUrlPattern.Rules {
			connUrlPattern = append(connUrlPattern, types.ObjectValueMust(
				map[string]attr.Type{
					"order":            types.Int64Type,
					"host_url_pattern": types.StringType,
					"data_group":       types.StringType,
				},
				map[string]attr.Value{
					"order":            types.Int64Value(rules.Order),
					"host_url_pattern": types.StringValue(rules.Pattern),
					"data_group":       types.StringValue(ConvertConnectorRuleClassificationValues(rules.ConnectorRuleClassification)),
				},
			))
		}
	}
	if len(urlPatterns) == 0 {
		return types.SetValueMust(customConnectorPatternSetObjectType, []attr.Value{})
	} else {
		return types.SetValueMust(customConnectorPatternSetObjectType, connUrlPattern)
	}
}

func ConvertToAttrValueEnvironments(environments []models.DlpEnvironmentDto) basetypes.SetValue {
	var env []attr.Value
	for _, environment := range environments {
		env = append(env, types.ObjectValueMust(
			map[string]attr.Type{
				"name": types.StringType,
			},
			map[string]attr.Value{
				"name": types.StringValue(environment.Name),
			},
		))
	}

	if len(environments) == 0 {
		return types.SetValueMust(environmentSetObjectType, []attr.Value{})
	} else {
		return types.SetValueMust(environmentSetObjectType, env)
	}
}

func ConvertToAttrValueConnectors(connectorsGroup models.DlpConnectorGroupsModel, connectors []attr.Value) []attr.Value {
	for _, connector := range connectorsGroup.Connectors {
		connectors = append(connectors, types.ObjectValueMust(
			map[string]attr.Type{
				//"name":                         types.StringType,
				"id":                           types.StringType,
				"default_action_rule_behavior": types.StringType,
				"action_rules": types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"action_id": types.StringType,
							"behavior":  types.StringType,
						},
					}},
				"endpoint_rules": types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"order":    types.Int64Type,
							"behavior": types.StringType,
							"endpoint": types.StringType,
						},
					}},
			},
			map[string]attr.Value{
				//"name":                         types.StringValue(connector.Name),
				"id":                           types.StringValue(connector.Id),
				"default_action_rule_behavior": types.StringValue(connector.DefaultActionRuleBehavior),
				"action_rules":                 types.ListValueMust(actionRuleListObjectType, ConvertToAtrValueActionRule(connector)),
				"endpoint_rules":               types.ListValueMust(endpointRuleListObjectType, ConvertToAtrValueEndpointRule(connector)),
			},
		))
	}
	return connectors
}

func ConvertToDlpConnectorGroup(ctx context.Context, diag diag.Diagnostics, classification string, connectorsAttr basetypes.SetValue) models.DlpConnectorGroupsModel {
	var connectors []DataLossPreventionPolicyResourceConnectorModel
	err := connectorsAttr.ElementsAs(ctx, &connectors, true)
	if err != nil {
		diag.AddError("Client error when converting DlpConnectorGroups", "")
	}

	connectorGroup := models.DlpConnectorGroupsModel{
		Classification: classification,
		Connectors:     make([]models.DlpConnectorModel, 0),
	}

	for _, connector := range connectors {
		defaultAction := "Allow"

		if connector.DefaultActionRuleBehavior.ValueString() != "" {
			defaultAction = connector.DefaultActionRuleBehavior.ValueString()
		}

		connectorGroup.Connectors = append(connectorGroup.Connectors, models.DlpConnectorModel{
			Id:   connector.Id.ValueString(),
			Type: "Microsoft.PowerApps/apis",

			DefaultActionRuleBehavior: defaultAction,
			ActionRules:               ConvertToDlpActionRule(connector),
			EndpointRules:             ConvertToDlpEndpointRule(connector),
		})
	}
	return connectorGroup
}

func ConvertToDlpEnvironment(ctx context.Context, diag diag.Diagnostics, environmentsAttr basetypes.SetValue) []models.DlpEnvironmentDto {
	var envs []DataLossPreventionPolicyResourceEnvironmentsModel
	err := environmentsAttr.ElementsAs(ctx, &envs, true)
	if err != nil {
		diag.AddError("Client error when converting DlpEnvironment", "")
	}

	environments := make([]models.DlpEnvironmentDto, 0)
	for _, environment := range envs {
		environments = append(environments, models.DlpEnvironmentDto{
			Name: environment.Name.ValueString(),
			Id:   "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/" + environment.Name.ValueString(),
			Type: "Microsoft.BusinessAppPlatform/scopes/environments",
		})
	}
	return environments
}

func ConvertToDlpCustomConnectorUrlPatternsDefinition(ctx context.Context, diag diag.Diagnostics, connectorPatternsAttr basetypes.SetValue) []models.DlpConnectorUrlPatternsDefinitionDto {
	var customConnectorsPatterns []DataLossPreventionPolicyResourceCustomConnectorPattern
	err := connectorPatternsAttr.ElementsAs(ctx, &customConnectorsPatterns, true)
	if err != nil {
		diag.AddError("Client error when converting DlpCustomConnectorUrlPatternsDefinition", "")
	}

	customConnectorUrlPatternsDefinition := make([]models.DlpConnectorUrlPatternsDefinitionDto, 0)
	for _, customConnectorPattern := range customConnectorsPatterns {
		urlPattern := models.DlpConnectorUrlPatternsDefinitionDto{
			Rules: []models.DlpConnectorUrlPatternsRuleDto{},
		}
		urlPattern.Rules = append(urlPattern.Rules, models.DlpConnectorUrlPatternsRuleDto{
			Order:                       customConnectorPattern.Order.ValueInt64(),
			ConnectorRuleClassification: ConvertConnectorRuleClassificationValues(customConnectorPattern.DataGroup.ValueString()),
			Pattern:                     customConnectorPattern.HostUrlPattern.ValueString(),
		})
		customConnectorUrlPatternsDefinition = append(customConnectorUrlPatternsDefinition, urlPattern)
	}
	return customConnectorUrlPatternsDefinition
}

func ConvertToDlpActionRule(connector DataLossPreventionPolicyResourceConnectorModel) []models.DlpActionRuleDto {
	var actionRules []models.DlpActionRuleDto
	for _, actionRule := range connector.ActionRules {
		actionRules = append(actionRules, models.DlpActionRuleDto{
			ActionId: actionRule.ActionId.ValueString(),
			Behavior: actionRule.Behavior.ValueString(),
		})
	}
	return actionRules
}

func ConvertToDlpEndpointRule(connector DataLossPreventionPolicyResourceConnectorModel) []models.DlpEndpointRuleDto {
	var endpointRules []models.DlpEndpointRuleDto
	for _, endpointRule := range connector.EndpointRules {
		endpointRules = append(endpointRules, models.DlpEndpointRuleDto{
			Order:    endpointRule.Order.ValueInt64(),
			Behavior: endpointRule.Behavior.ValueString(),
			Endpoint: endpointRule.Endpoint.ValueString(),
		})
	}
	return endpointRules
}

func ConvertToAtrValueActionRule(connector models.DlpConnectorModel) []attr.Value {
	var actionRules []attr.Value
	for _, actionRule := range connector.ActionRules {
		actionRules = append(actionRules, types.ObjectValueMust(
			map[string]attr.Type{
				"action_id": types.StringType,
				"behavior":  types.StringType,
			},
			map[string]attr.Value{
				"action_id": types.StringValue(actionRule.ActionId),
				"behavior":  types.StringValue(actionRule.Behavior),
			},
		))
	}
	return actionRules
}

func ConvertToAtrValueEndpointRule(connector models.DlpConnectorModel) []attr.Value {
	var endpointRules []attr.Value
	for _, endpointRule := range connector.EndpointRules {
		endpointRules = append(endpointRules, types.ObjectValueMust(
			map[string]attr.Type{
				"order":    types.Int64Type,
				"behavior": types.StringType,
				"endpoint": types.StringType,
			},
			map[string]attr.Value{
				"order":    types.Int64Value(endpointRule.Order),
				"behavior": types.StringValue(endpointRule.Behavior),
				"endpoint": types.StringValue(endpointRule.Endpoint),
			},
		))
	}
	return endpointRules
}
