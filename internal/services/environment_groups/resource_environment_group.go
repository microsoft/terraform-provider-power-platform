// Licensed under the MIT license.
// Copyright (c) Microsoft Corporation.

package environment_groups

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment_group_rule_set"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/tenant"
)

var _ resource.Resource = &EnvironmentGroupResource{}
var _ resource.ResourceWithImportState = &EnvironmentGroupResource{}

func NewEnvironmentGroupResource() resource.Resource {
	return &EnvironmentGroupResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment_group",
		},
	}
}

func (r *EnvironmentGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *EnvironmentGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource manages an [Environment Group](https://learn.microsoft.com/en-us/power-platform/admin/environment-groups).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique id of the environment group",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the environment group",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Display name of the environment group",
				Required:            true,
			},
		},
	}
}

func (r *EnvironmentGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.EnvironmentGroupClient = newEnvironmentGroupClient(client, tenant.NewTenantClient(client), environment_group_rule_set.NewEnvironmentGroupRuleSetClient(client, tenant.NewTenantClient(client)))
}

// Read function for EnvironmentGroupResource.
func (r *EnvironmentGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	state := EnvironmentGroupResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	environmentGroup, err := r.EnvironmentGroupClient.GetEnvironmentGroup(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	if environmentGroup == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Id = types.StringValue(environmentGroup.Id)
	state.DisplayName = types.StringValue(environmentGroup.DisplayName)
	state.Description = types.StringValue(environmentGroup.Description)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Create function.
func (r *EnvironmentGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *EnvironmentGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	environmentGroupToCreate := environmentGroupDto{
		DisplayName: plan.DisplayName.ValueString(),
		Description: plan.Description.ValueString(),
	}

	eg, err := r.EnvironmentGroupClient.CreateEnvironmentGroup(ctx, environmentGroupToCreate)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	state := EnvironmentGroupResourceModel{}
	state.Id = types.StringValue(eg.Id)
	state.DisplayName = types.StringValue(eg.DisplayName)
	state.Description = types.StringValue(eg.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update function.
func (r *EnvironmentGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *EnvironmentGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *EnvironmentGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	environmentGroupToUpdate := environmentGroupDto{
		Id:          plan.Id.ValueString(),
		DisplayName: plan.DisplayName.ValueString(),
		Description: plan.Description.ValueString(),
	}

	eg, err := r.EnvironmentGroupClient.UpdateEnvironmentGroup(ctx, state.Id.ValueString(), environmentGroupToUpdate)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	state.Id = types.StringValue(eg.Id)
	state.DisplayName = types.StringValue(eg.DisplayName)
	state.Description = types.StringValue(eg.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete function.
func (r *EnvironmentGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *EnvironmentGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.EnvironmentGroupClient.DeleteEnvironmentGroup(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	if customerrors.Code(err) == customerrors.ERROR_ENVIRONMENT_URL_NOT_FOUND || customerrors.Code(err) == customerrors.ERROR_POLICY_ASSIGNED_TO_ENV_GROUP {
		envs, err := r.EnvironmentGroupClient.GetEnvironmentsInEnvironmentGroup(ctx, state.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
			return
		}

		if len(envs) > 0 {
			tflog.Debug(ctx, fmt.Sprintf("Environment group %s_%s has %d environments. Removing them.", r.ProviderTypeName, r.TypeName, len(envs)))
			for _, env := range envs {
				err := r.EnvironmentGroupClient.RemoveEnvironmentFromEnvironmentGroup(ctx, state.Id.ValueString(), env.Name)
				if err != nil {
					resp.Diagnostics.AddError("error when removing environment", err.Error())
					return
				}
			}
		}

		ruleSet, err := r.EnvironmentGroupClient.RuleSetApi.GetEnvironmentGroupRuleSet(ctx, state.Id.ValueString())
		if err != nil && customerrors.Code(err) != customerrors.ERROR_OBJECT_NOT_FOUND {
			resp.Diagnostics.AddError("Failed to get environment group ruleset", err.Error())
			return
		}

		if customerrors.Code(err) != customerrors.ERROR_OBJECT_NOT_FOUND && ruleSet != nil && len(ruleSet.Parameters) > 0 {
			tflog.Debug(ctx, fmt.Sprintf("Environment group %s_%s has %d rule sets. Deleting them.", r.ProviderTypeName, r.TypeName, len(ruleSet.Parameters)))
			err := r.EnvironmentGroupClient.RuleSetApi.DeleteEnvironmentGroupRuleSet(ctx, *ruleSet.Id)
			if err != nil {
				resp.Diagnostics.AddError("error when deleting rule set", err.Error())
				return
			}
		}

		err = r.EnvironmentGroupClient.DeleteEnvironmentGroup(ctx, state.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
			return
		}
	}
}

// ImportState function.
func (r *EnvironmentGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
