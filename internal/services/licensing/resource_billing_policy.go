// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package licensing

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

var _ resource.Resource = &BillingPolicyResource{}
var _ resource.ResourceWithImportState = &BillingPolicyResource{}

func NewBillingPolicyResource() resource.Resource {
	return &BillingPolicyResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "billing_policy",
		},
	}
}

func (r *BillingPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *BillingPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Power Platform Billing Policy. \n\nA Power Platform billing policy is a mechanism that allows you to manage the costs associated with your Power Platform usage. It's linked to an Azure subscription and is used to set up pay-as-you-go billing for an environment.\n\nAdditional Resources:\n\n* [What is a billing policy](https://learn.microsoft.com/power-platform/admin/pay-as-you-go-overview#what-is-a-billing-policy)\n* [Power Platform Billing Policy API](https://learn.microsoft.com/rest/api/power-platform/licensing/billing-policy/get-billing-policy)",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The id of the billing policy",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the billing policy",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "The location of the billing policy",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the billing policy (Enabled, Disabled)",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Enabled", "Disabled"),
				},
			},
			"billing_instrument": schema.SingleNestedAttribute{
				MarkdownDescription: "The billing instrument of the billing policy",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The id of the billing instrument",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"resource_group": schema.StringAttribute{
						MarkdownDescription: "The resource group of the billing instrument",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"subscription_id": schema.StringAttribute{
						MarkdownDescription: "The subscription id of the billing instrument",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
			},
		},
	}
}

func (r *BillingPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BillingPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *BillingPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	billingPolicyToCreate := billingPolicyCreateDto{
		BillingInstrument: BillingInstrumentDto{
			ResourceGroup:  plan.BillingInstrument.ResourceGroup.ValueString(),
			SubscriptionId: plan.BillingInstrument.SubscriptionId.ValueString(),
		},
		Location: plan.Location.ValueString(),
		Name:     plan.Name.ValueString(),
	}

	if plan.Status.IsUnknown() {
		billingPolicyToCreate.Status = "Enabled"
	} else {
		billingPolicyToCreate.Status = plan.Status.ValueString()
	}

	policy, err := r.LicensingClient.CreateBillingPolicy(ctx, billingPolicyToCreate)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err.Error())
		return
	}

	plan.Id = types.StringValue(policy.Id)
	plan.Name = types.StringValue(policy.Name)
	plan.Location = types.StringValue(policy.Location)
	plan.Status = types.StringValue(policy.Status)
	plan.BillingInstrument.Id = types.StringValue(policy.BillingInstrument.Id)
	plan.BillingInstrument.ResourceGroup = types.StringValue(policy.BillingInstrument.ResourceGroup)
	plan.BillingInstrument.SubscriptionId = types.StringValue(policy.BillingInstrument.SubscriptionId)

	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID %s", plan.Id.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BillingPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *BillingPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	billing, err := r.LicensingClient.GetBillingPolicy(ctx, state.Id.ValueString())
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
		return
	}

	state.Id = types.StringValue(billing.Id)
	state.Name = types.StringValue(billing.Name)
	state.Location = types.StringValue(billing.Location)
	state.Status = types.StringValue(billing.Status)
	state.BillingInstrument.Id = types.StringValue(billing.BillingInstrument.Id)
	state.BillingInstrument.ResourceGroup = types.StringValue(billing.BillingInstrument.ResourceGroup)
	state.BillingInstrument.SubscriptionId = types.StringValue(billing.BillingInstrument.SubscriptionId)

	tflog.Debug(ctx, fmt.Sprintf("READ %s with Id: %s", r.FullTypeName(), billing.Id))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BillingPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *BillingPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *BillingPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Name.ValueString() != state.Name.ValueString() ||
		plan.Status.ValueString() != state.Status.ValueString() {
		policyToUpdate := BillingPolicyUpdateDto{
			Name:   plan.Name.ValueString(),
			Status: plan.Status.ValueString(),
		}

		policy, err := r.LicensingClient.UpdateBillingPolicy(ctx, plan.Id.ValueString(), policyToUpdate)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
			return
		}

		plan.Id = types.StringValue(policy.Id)
		plan.Name = types.StringValue(policy.Name)
		plan.Location = types.StringValue(policy.Location)
		plan.Status = types.StringValue(policy.Status)
		plan.BillingInstrument.Id = types.StringValue(policy.BillingInstrument.Id)
		plan.BillingInstrument.ResourceGroup = types.StringValue(policy.BillingInstrument.ResourceGroup)
		plan.BillingInstrument.SubscriptionId = types.StringValue(policy.BillingInstrument.SubscriptionId)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BillingPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *BillingPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.LicensingClient.DeleteBillingPolicy(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
		return
	}
}

func (r *BillingPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
