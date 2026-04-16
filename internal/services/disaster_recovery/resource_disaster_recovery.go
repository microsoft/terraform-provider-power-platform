// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package disaster_recovery

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ resource.Resource = &DisasterRecoveryResource{}
var _ resource.ResourceWithImportState = &DisasterRecoveryResource{}

func NewDisasterRecoveryResource() resource.Resource {
	return &DisasterRecoveryResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment_disaster_recovery",
		},
	}
}

func (r *DisasterRecoveryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *DisasterRecoveryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		return
	}

	providerClient, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	if providerClient.Api == nil {
		resp.Diagnostics.AddError(
			"Nil Api client",
			"ProviderData contained a *api.ProviderClient but with nil Api. Please check provider initialization and credentials.",
		)
		return
	}
	r.client = newDisasterRecoveryClient(providerClient.Api)
}

func (r *DisasterRecoveryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages [Disaster Recovery](https://learn.microsoft.com/en-us/power-platform/admin/business-continuity-disaster-recovery?tabs=new#cross-region-self-service-disaster-recovery) for a Power Platform environment. " +
			"This resource enables or disables disaster recovery for an existing production environment. " +
			"The environment must be a managed environment with a billing policy attached. " +
			"The paired/secondary region is determined by the environment's configuration and is not user-selectable through this resource.",

		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the disaster recovery resource (same as the environment ID)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier (GUID) of the environment to enable disaster recovery for. The environment must be a production managed environment.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether disaster recovery is enabled for the environment. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
		},
	}
}

func (r *DisasterRecoveryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan DisasterRecoveryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	environmentId := plan.EnvironmentId.ValueString()
	enabled := plan.Enabled.ValueBool()

	if enabled {
		err := r.client.EnableDisasterRecovery(ctx, environmentId)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when enabling disaster recovery for environment '%s'", environmentId), err.Error())
			return
		}
	} else {
		err := r.client.DisableDisasterRecovery(ctx, environmentId)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when disabling disaster recovery for environment '%s'", environmentId), err.Error())
			return
		}
	}

	env, err := r.client.GetDisasterRecoveryState(ctx, environmentId)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading environment '%s' after disaster recovery operation", environmentId), err.Error())
		return
	}

	model := convertDtoToModel(environmentId, env)
	model.Timeouts = plan.Timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *DisasterRecoveryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state DisasterRecoveryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	environmentId := state.EnvironmentId.ValueString()

	env, err := r.client.GetDisasterRecoveryState(ctx, environmentId)
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
		return
	}

	model := convertDtoToModel(environmentId, env)
	model.Timeouts = state.Timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *DisasterRecoveryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan DisasterRecoveryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	environmentId := plan.EnvironmentId.ValueString()
	enabled := plan.Enabled.ValueBool()

	if enabled {
		err := r.client.EnableDisasterRecovery(ctx, environmentId)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when enabling disaster recovery for environment '%s'", environmentId), err.Error())
			return
		}
	} else {
		err := r.client.DisableDisasterRecovery(ctx, environmentId)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when disabling disaster recovery for environment '%s'", environmentId), err.Error())
			return
		}
	}

	env, err := r.client.GetDisasterRecoveryState(ctx, environmentId)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading environment '%s' after disaster recovery update", environmentId), err.Error())
		return
	}

	model := convertDtoToModel(environmentId, env)
	model.Timeouts = plan.Timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *DisasterRecoveryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state DisasterRecoveryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	environmentId := state.EnvironmentId.ValueString()

	env, err := r.client.GetDisasterRecoveryState(ctx, environmentId)
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading environment '%s' before disabling disaster recovery", environmentId), err.Error())
		return
	}

	drEnabled := env.Properties != nil && env.Properties.States != nil && env.Properties.States.DisasterRecovery != nil && env.Properties.States.DisasterRecovery.Id == "Enabled"
	if !drEnabled {
		tflog.Debug(ctx, fmt.Sprintf("Disaster recovery is already disabled for environment '%s', skipping disable call", environmentId))
		return
	}

	err = r.client.DisableDisasterRecovery(ctx, environmentId)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when disabling disaster recovery for environment '%s'", environmentId), err.Error())
		return
	}
}

func (r *DisasterRecoveryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), req.ID)...)
}
