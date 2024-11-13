// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package enterprise_policy

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ resource.Resource = &Resource{}

func NewEnterpisePolicyResource() resource.Resource {
	return &Resource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "enterprise_policy",
		},
	}
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
		return
	}

	clientApi := req.ProviderData.(*api.ProviderClient).Api
	if clientApi == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.EnterprisePolicyClient = newEnterprisePolicyClient(clientApi)
	tflog.Debug(ctx, "Successfully created clients")
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Enterprise Policy environment assignment",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the resource",
				Computed:            true,
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Environment id",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"system_id": schema.StringAttribute{
				MarkdownDescription: "Policy SystemId value in following format `/regions/<location>/providers/Microsoft.PowerPlatform/enterprisePolicies/<policyid>`",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy_type": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("Policy type [%s, %s]", NETWORK_INJECTION_POLICY_TYPE, ENCRYPTION_POLICY_TYPE),
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(NETWORK_INJECTION_POLICY_TYPE, ENCRYPTION_POLICY_TYPE),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *sourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.EnterprisePolicyClient.LinkEnterprisePolicy(ctx, plan.EnvironmentId.ValueString(), plan.PolicyType.ValueString(), plan.SystemId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	newState := sourceModel{
		Timeouts:      plan.Timeouts,
		Id:            types.StringValue(fmt.Sprintf("%s_%s", plan.EnvironmentId.ValueString(), strings.ToLower(plan.PolicyType.ValueString()))),
		SystemId:      plan.SystemId,
		PolicyType:    plan.PolicyType,
		EnvironmentId: plan.EnvironmentId,
	}

	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID %s", plan.Id.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *sourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	env, err := r.EnterprisePolicyClient.EnvironmentClient.GetEnvironment(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading environment in %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}
	if env.Properties.EnterprisePolicies != nil {
		if state.PolicyType.ValueString() == NETWORK_INJECTION_POLICY_TYPE && env.Properties.EnterprisePolicies.Vnets != nil {
			state.SystemId = types.StringValue(env.Properties.EnterprisePolicies.Vnets.SystemId)
		} else if state.PolicyType.ValueString() == ENCRYPTION_POLICY_TYPE && env.Properties.EnterprisePolicies.CustomerManagedKeys != nil {
			state.SystemId = types.StringValue(env.Properties.EnterprisePolicies.CustomerManagedKeys.SystemId)
		}
	} else {
		state.SystemId = types.StringNull()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *sourceModel
	var state *sourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// There is nothing to update in this resource, as any attribute change would require a delete and create.
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *sourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.EnterprisePolicyClient.UnLinkEnterprisePolicy(ctx, state.EnvironmentId.ValueString(), state.PolicyType.ValueString(), state.SystemId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}
}
