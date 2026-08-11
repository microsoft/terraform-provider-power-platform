// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environmentvariable

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
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
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ resource.Resource = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}

func NewEnvironmentVariableResource() resource.Resource {
	return &Resource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment_variable",
		},
	}
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the current value for an existing Power Platform environment variable definition. This resource does not create the definition itself. The definition must already exist in the target environment, typically because it was deployed by a solution. Deleting this resource removes the current value record so the environment falls back to the definition's default value, if one exists.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Read:   true,
				Delete: true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Composite resource id in the format `{environment_id}/{schema_name}`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Environment id where the environment variable definition already exists.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"schema_name": schema.StringAttribute{
				MarkdownDescription: "Schema name of the existing environment variable definition to manage.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Current value to set for the environment variable in the target environment.",
				Required:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(2000),
				},
			},
			"environment_variable_definition_id": schema.StringAttribute{
				MarkdownDescription: "Dataverse id of the environment variable definition.",
				Computed:            true,
			},
			"environment_variable_value_id": schema.StringAttribute{
				MarkdownDescription: "Dataverse id of the current environment variable value record.",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the environment variable definition.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the environment variable definition.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the environment variable definition.",
				Computed:            true,
			},
			"value_schema": schema.StringAttribute{
				MarkdownDescription: "Value schema declared on the environment variable definition, when present.",
				Computed:            true,
			},
			"secret_store": schema.StringAttribute{
				MarkdownDescription: "Secret store configured on the environment variable definition.",
				Computed:            true,
			},
		},
	}
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
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

	r.Client = newEnvironmentVariableClient(client.Api)
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.applyEnvironmentVariableValue(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create environment variable value", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resolved, err := r.Client.GetEnvironmentVariable(ctx, state.EnvironmentId.ValueString(), state.SchemaName.ValueString())
	if err != nil {
		// The definition being gone and the parent environment being gone (deleted out of band)
		// both mean the value no longer exists — remove it from state instead of failing refresh.
		if isEnvironmentVariableDefinitionNotFound(err) || errors.Is(err, customerrors.ErrObjectNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read environment variable value", err.Error())
		return
	}
	if resolved.Value == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	newState := buildStateFromResolved(state.EnvironmentId.ValueString(), state.Value.ValueString(), state.Timeouts, resolved)
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.applyEnvironmentVariableValue(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update environment variable value", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.Client.DeleteEnvironmentVariableValue(ctx, state.EnvironmentId.ValueString(), state.SchemaName.ValueString())
	if err != nil {
		// A missing definition or a parent environment deleted out of band both mean the
		// value is already gone; treat the delete as complete.
		if isEnvironmentVariableDefinitionNotFound(err) || errors.Is(err, customerrors.ErrObjectNotFound) {
			return
		}
		resp.Diagnostics.AddError("Failed to delete environment variable value", err.Error())
		return
	}
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	environmentID, schemaName, err := parseEnvironmentVariableResourceID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import id", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), environmentVariableResourceID(environmentID, schemaName))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), environmentID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("schema_name"), schemaName)...)
}

func (r *Resource) applyEnvironmentVariableValue(ctx context.Context, plan *ResourceModel) (*ResourceModel, error) {
	dataverseExists, err := r.Client.DataverseExists(ctx, plan.EnvironmentId.ValueString())
	if err != nil {
		return nil, err
	}
	if !dataverseExists {
		return nil, fmt.Errorf("environment %s does not have Dataverse linked", plan.EnvironmentId.ValueString())
	}

	resolved, err := r.Client.UpsertEnvironmentVariableValue(ctx, plan.EnvironmentId.ValueString(), plan.SchemaName.ValueString(), plan.Value.ValueString())
	if err != nil {
		return nil, err
	}
	if resolved.Value == nil {
		return nil, fmt.Errorf("environment variable '%s' did not return a current value after apply", plan.SchemaName.ValueString())
	}

	return buildStateFromResolved(plan.EnvironmentId.ValueString(), plan.Value.ValueString(), plan.Timeouts, resolved), nil
}

func buildStateFromResolved(environmentID, configuredValue string, timeoutValue timeouts.Value, resolved *resolvedEnvironmentVariableDto) *ResourceModel {
	state := &ResourceModel{
		Timeouts:                      timeoutValue,
		Id:                            types.StringValue(environmentVariableResourceID(environmentID, resolved.Definition.SchemaName)),
		EnvironmentId:                 types.StringValue(environmentID),
		SchemaName:                    types.StringValue(resolved.Definition.SchemaName),
		Value:                         types.StringValue(configuredValue),
		EnvironmentVariableDefinition: types.StringValue(resolved.Definition.EnvironmentVariableDefinitionId),
		DisplayName:                   types.StringValue(resolved.Definition.DisplayName),
		Description:                   types.StringValue(resolved.Definition.Description),
		Type:                          types.StringValue(environmentVariableTypeLabel(resolved.Definition.Type)),
		ValueSchema:                   types.StringValue(resolved.Definition.ValueSchema),
		SecretStore:                   types.StringValue(environmentVariableSecretStoreLabel(resolved.Definition.SecretStore)),
	}

	if resolved.Value != nil {
		state.EnvironmentVariableValue = types.StringValue(resolved.Value.EnvironmentVariableValueId)
	} else {
		state.EnvironmentVariableValue = types.StringNull()
	}

	if resolved.Definition.Description == "" {
		state.Description = types.StringNull()
	}
	if resolved.Definition.ValueSchema == "" {
		state.ValueSchema = types.StringNull()
	}

	return state
}
