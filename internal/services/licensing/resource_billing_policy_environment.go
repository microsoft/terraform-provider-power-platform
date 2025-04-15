// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package licensing

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
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
)

var _ resource.Resource = &BillingPolicyEnvironmentResource{}
var _ resource.ResourceWithImportState = &BillingPolicyEnvironmentResource{}

func NewBillingPolicyEnvironmentResource() resource.Resource {
	return &BillingPolicyEnvironmentResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "billing_policy_environment",
		},
	}
}

func (r *BillingPolicyEnvironmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *BillingPolicyEnvironmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource allows you to manage the environments associated with a [billing policy](https://learn.microsoft.com/power-platform/admin/pay-as-you-go-overview#what-is-a-billing-policy). A billing policy is a set of rules that define how a tenant is billed for usage of Power Platform services. A billing policy is associated with a billing instrument, which is a subscription and resource group that is used to pay for usage of Power Platform services.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"billing_policy_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The id of the billing policy",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environments": schema.SetAttribute{
				MarkdownDescription: "The environments associated with the billing policy",
				ElementType:         types.StringType,
				Required:            true,
			},
		},
	}
}

func (r *BillingPolicyEnvironmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.LicensingClient = NewLicensingClient(client.Api)
}

func (r *BillingPolicyEnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var plan *BillingPolicyEnvironmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	environments, err := r.LicensingClient.GetEnvironmentsForBillingPolicy(ctx, plan.BillingPolicyId)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
		return
	}

	if len(environments) > 0 {
		err = r.LicensingClient.RemoveEnvironmentsToBillingPolicy(ctx, plan.BillingPolicyId, environments)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
			return
		}
	}

	err = r.LicensingClient.AddEnvironmentsToBillingPolicy(ctx, plan.BillingPolicyId, plan.Environments)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
		return
	}

	environments, err = r.LicensingClient.GetEnvironmentsForBillingPolicy(ctx, plan.BillingPolicyId)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
		return
	}

	plan.Environments = environments

	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID %s", plan.BillingPolicyId))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BillingPolicyEnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var state *BillingPolicyEnvironmentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	environments, err := r.LicensingClient.GetEnvironmentsForBillingPolicy(ctx, state.BillingPolicyId)
	if err != nil {
		if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
		return
	}

	state.Environments = environments
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("READ %s with Id: %s", r.FullTypeName(), state.BillingPolicyId))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BillingPolicyEnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *BillingPolicyEnvironmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state *BillingPolicyEnvironmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	environments, err := r.LicensingClient.GetEnvironmentsForBillingPolicy(ctx, plan.BillingPolicyId)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
		return
	}

	err = r.LicensingClient.RemoveEnvironmentsToBillingPolicy(ctx, plan.BillingPolicyId, environments)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
		return
	}

	err = r.LicensingClient.AddEnvironmentsToBillingPolicy(ctx, plan.BillingPolicyId, plan.Environments)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
		return
	}

	environments, err = r.LicensingClient.GetEnvironmentsForBillingPolicy(ctx, plan.BillingPolicyId)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
		return
	}

	plan.Environments = environments

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BillingPolicyEnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *BillingPolicyEnvironmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.LicensingClient.RemoveEnvironmentsToBillingPolicy(ctx, state.BillingPolicyId, state.Environments)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
		return
	}
}

func (r *BillingPolicyEnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("billing_policy_id"), req, resp)
}
