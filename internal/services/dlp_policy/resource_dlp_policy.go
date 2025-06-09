// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package dlp_policy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ resource.Resource = &DataLossPreventionPolicyResource{}
var _ resource.ResourceWithImportState = &DataLossPreventionPolicyResource{}
var _ resource.ResourceWithValidateConfig = &DataLossPreventionPolicyResource{}

func NewDataLossPreventionPolicyResource() resource.Resource {
	return &DataLossPreventionPolicyResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "data_loss_prevention_policy",
		},
	}
}

func (r *DataLossPreventionPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *DataLossPreventionPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	connectorSchema := schema.NestedAttributeObject{

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the connector",
				Optional:            true,
			},
			"default_action_rule_behavior": schema.StringAttribute{
				MarkdownDescription: "Default action rule behavior for the connector (\"Allow\", \"Block\", \"\")",
				Optional:            true,
				Computed:            true,
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
		MarkdownDescription: "This resource manages a Data Loss Prevention Policy. See [Data Loss Prevention](https://learn.microsoft.com/power-platform/admin/prevent-data-loss) for more information.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique name of the policy",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the policy",
				Required:            true,
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "User who created the policy",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_time": schema.StringAttribute{
				MarkdownDescription: "Time when the policy was created",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("General", "Confidential", "Blocked"),
				},
			},
			"environments": schema.SetAttribute{
				MarkdownDescription: "Environment to which the policy is applied",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"business_connectors": schema.SetNestedAttribute{
				MarkdownDescription: "Connectors for sensitive data",
				Required:            true,
				NestedObject:        connectorSchema,
			},
			"non_business_connectors": schema.SetNestedAttribute{
				MarkdownDescription: "Connectors for non-sensitive data",
				Required:            true,
				NestedObject:        connectorSchema,
			},
			"blocked_connectors": schema.SetNestedAttribute{
				MarkdownDescription: "Blocked connectors can’t be used where this policy is applied.",
				Required:            true,
				NestedObject:        connectorSchema,
			},
			"custom_connectors_patterns": schema.SetNestedAttribute{
				MarkdownDescription: "Custom connectors patterns",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"order": schema.Int64Attribute{
							MarkdownDescription: "Order of the connector",
							Required:            true,
						},
						"host_url_pattern": schema.StringAttribute{
							MarkdownDescription: "Pattern of the connector",
							Required:            true,
						},
						"data_group": schema.StringAttribute{
							MarkdownDescription: "Data group of the connector (\"Business\", \"NonBusiness\", \"Blocked\", \"Ignore\")",
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
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
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

	r.DlpPolicyClient = newDlpPolicyClient(client.Api)
}

func (r DataLossPreventionPolicyResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var config *dataLossPreventionPolicyResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var connectors []dlpConnectorModelDto
	conn, err := getConnectorGroup(ctx, config.BusinessGeneralConnectors)
	if err != nil {
		resp.Diagnostics.AddError("BusinessGeneralConnectors validation error", err.Error())
	}
	connectors = append(connectors, conn.Connectors...)

	conn, err = getConnectorGroup(ctx, config.NonBusinessConfidentialConnectors)
	if err != nil {
		resp.Diagnostics.AddError("NonBusinessConfidentialConnectors validation error", err.Error())
	}
	connectors = append(connectors, conn.Connectors...)

	conn, err = getConnectorGroup(ctx, config.BlockedConnectors)
	if err != nil {
		resp.Diagnostics.AddError("BlockedConnectors validation error", err.Error())
	}
	connectors = append(connectors, conn.Connectors...)

	for _, c := range connectors {
		if (c.DefaultActionRuleBehavior != "" && len(c.ActionRules) == 0) || (c.DefaultActionRuleBehavior == "" && len(c.ActionRules) > 0) {
			resp.Diagnostics.AddAttributeError(
				path.Empty(),
				"Incorrect attribute Configuration",
				"Expected 'default_action_rule_behavior' to be empty if 'action_rules' are empty.",
			)
		}
	}
}

func (r *DataLossPreventionPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var state *dataLossPreventionPolicyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := r.DlpPolicyClient.GetPolicy(ctx, state.Id.ValueString())
	if err != nil {
		if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
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
	state.BusinessGeneralConnectors = convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
	state.NonBusinessConfidentialConnectors = convertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
	state.BlockedConnectors = convertToAttrValueConnectorsGroup("Blocked", policy.ConnectorGroups)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DataLossPreventionPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *dataLossPreventionPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyToCreate := dlpPolicyModelDto{
		DefaultConnectorsClassification:      plan.DefaultConnectorsClassification.ValueString(),
		DisplayName:                          plan.DisplayName.ValueString(),
		EnvironmentType:                      plan.EnvironmentType.ValueString(),
		Environments:                         []dlpEnvironmentDto{},
		ConnectorGroups:                      []dlpConnectorGroupsModelDto{},
		CustomConnectorUrlPatternsDefinition: []dlpConnectorUrlPatternsDefinitionDto{},
	}

	policyToCreate.Environments = convertToDlpEnvironment(ctx, plan.Environments)

	customConnectorUrlPatterns, err := convertToDlpCustomConnectorUrlPatternsDefinition(ctx, resp.Diagnostics, plan.CustomConnectorsPatterns)
	if err != nil {
		return
	}
	policyToCreate.CustomConnectorUrlPatternsDefinition = customConnectorUrlPatterns

	policyToCreate.ConnectorGroups = make([]dlpConnectorGroupsModelDto, 0)

	err = r.buildConnectorGroups(ctx, resp.Diagnostics, &policyToCreate.ConnectorGroups, plan)
	if err != nil {
		return
	}

	policy, err_client := r.DlpPolicyClient.CreatePolicy(ctx, policyToCreate)
	if err_client != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err_client.Error())
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
	plan.BusinessGeneralConnectors = convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
	plan.NonBusinessConfidentialConnectors = convertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
	plan.BlockedConnectors = convertToAttrValueConnectorsGroup("Blocked", policy.ConnectorGroups)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DataLossPreventionPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *dataLossPreventionPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state *dataLossPreventionPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	policyToUpdate := dlpPolicyModelDto{
		Name:                            plan.Id.ValueString(),
		DisplayName:                     plan.DisplayName.ValueString(),
		EnvironmentType:                 plan.EnvironmentType.ValueString(),
		DefaultConnectorsClassification: plan.DefaultConnectorsClassification.ValueString(),
		Environments:                    []dlpEnvironmentDto{},
		ConnectorGroups:                 []dlpConnectorGroupsModelDto{},
	}

	policyToUpdate.Environments = convertToDlpEnvironment(ctx, plan.Environments)

	customConnectorUrlPatterns, err := convertToDlpCustomConnectorUrlPatternsDefinition(ctx, resp.Diagnostics, plan.CustomConnectorsPatterns)
	if err != nil {
		return
	}
	policyToUpdate.CustomConnectorUrlPatternsDefinition = customConnectorUrlPatterns

	policyToUpdate.ConnectorGroups = make([]dlpConnectorGroupsModelDto, 0)

	err = r.buildConnectorGroups(ctx, resp.Diagnostics, &policyToUpdate.ConnectorGroups, plan)
	if err != nil {
		return
	}

	policy, err_client := r.DlpPolicyClient.UpdatePolicy(ctx, policyToUpdate)
	if err_client != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err_client.Error())
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
	plan.BusinessGeneralConnectors = convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
	plan.NonBusinessConfidentialConnectors = convertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
	plan.BlockedConnectors = convertToAttrValueConnectorsGroup("Blocked", policy.ConnectorGroups)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DataLossPreventionPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *dataLossPreventionPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.DlpPolicyClient.DeletePolicy(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
		return
	}
}

// buildConnectorGroups is a helper function that builds connector groups for DLP policies.
// It takes the target slice and a set of connector group configurations and processes them.
func (r *DataLossPreventionPolicyResource) buildConnectorGroups(
	ctx context.Context,
	diagnostics diag.Diagnostics,
	connectorGroups *[]dlpConnectorGroupsModelDto,
	plan *dataLossPreventionPolicyResourceModel,
) error {
	// Define the connector group configurations
	configs := []struct {
		classification string
		connectors     types.Set
	}{
		{"Confidential", plan.BusinessGeneralConnectors},
		{"General", plan.NonBusinessConfidentialConnectors},
		{"Blocked", plan.BlockedConnectors},
	}

	// Process each configuration
	for _, config := range configs {
		connectorGroup, err := convertToDlpConnectorGroup(ctx, diagnostics, config.classification, config.connectors)
		if err != nil {
			return err
		}
		*connectorGroups = append(*connectorGroups, connectorGroup)
	}

	return nil
}

func (r *DataLossPreventionPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
