// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

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
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

var _ resource.Resource = &BillingPolicyEnvironmentResource{}
var _ resource.ResourceWithImportState = &BillingPolicyEnvironmentResource{}

func NewBillingPolicyEnvironmentResource() resource.Resource {
	return &BillingPolicyEnvironmentResource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_billing_policy_environment",
	}
}

type BillingPolicyEnvironmentResource struct {
	LicensingClient  LicensingClient
	ProviderTypeName string
	TypeName         string
}

type BillingPolicyEnvironmentResourceModel struct {
	BillingPolicyId string   `tfsdk:"billing_policy_id"`
	Environments    []string `tfsdk:"environments"`
}

func (r *BillingPolicyEnvironmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

func (r *BillingPolicyEnvironmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This resource allows you to manage the environments associated with a Billing Policy",
		MarkdownDescription: "This resource allows you to manage the environments associated with a [billing policy](https://learn.microsoft.com/en-us/power-platform/admin/pay-as-you-go-overview#what-is-a-billing-policy). A billing policy is a set of rules that define how a tenant is billed for usage of Power Platform services. A billing policy is associated with a billing instrument, which is a subscription and resource group that is used to pay for usage of Power Platform services.",
		Attributes: map[string]schema.Attribute{
			"billing_policy_id": schema.StringAttribute{
				Required:            true,
				Description:         "The id of the billing policy",
				MarkdownDescription: "The id of the billing policy",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environments": schema.SetAttribute{
				Description:         "The environments associated with the billing policy",
				MarkdownDescription: "The environments associated with the billing policy",
				ElementType:         types.StringType,
				Required:            true,
			},
		},
	}
}

func (r *BillingPolicyEnvironmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
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
	r.LicensingClient = NewLicensingClient(clientApi)
}

func (r *BillingPolicyEnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *BillingPolicyEnvironmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	environments, err := r.LicensingClient.GetEnvironmentsForBillingPolicy(ctx, plan.BillingPolicyId)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.ProviderTypeName), err.Error())
		return
	}

	if len(environments) > 0 {
		err = r.LicensingClient.RemoveEnvironmentsToBillingPolicy(ctx, plan.BillingPolicyId, environments)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.ProviderTypeName), err.Error())
			return
		}
	}

	err = r.LicensingClient.AddEnvironmentsToBillingPolicy(ctx, plan.BillingPolicyId, plan.Environments)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.ProviderTypeName), err.Error())
		return
	}

	environments, err = r.LicensingClient.GetEnvironmentsForBillingPolicy(ctx, plan.BillingPolicyId)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.ProviderTypeName), err.Error())
		return
	}

	plan.Environments = environments

	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID %s", plan.BillingPolicyId))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *BillingPolicyEnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *BillingPolicyEnvironmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	environments, err := r.LicensingClient.GetEnvironmentsForBillingPolicy(ctx, state.BillingPolicyId)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	resp.State.SetAttribute(ctx, path.Root("environments"), environments)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("READ %s_%s with Id: %s", r.ProviderTypeName, r.TypeName, state.BillingPolicyId))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE END: %s", r.ProviderTypeName))
}

func (r *BillingPolicyEnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *BillingPolicyEnvironmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state *BillingPolicyEnvironmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	environments, err := r.LicensingClient.GetEnvironmentsForBillingPolicy(ctx, plan.BillingPolicyId)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.ProviderTypeName), err.Error())
		return
	}

	err = r.LicensingClient.RemoveEnvironmentsToBillingPolicy(ctx, plan.BillingPolicyId, environments)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.ProviderTypeName), err.Error())
		return
	}

	err = r.LicensingClient.AddEnvironmentsToBillingPolicy(ctx, plan.BillingPolicyId, plan.Environments)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.ProviderTypeName), err.Error())
		return
	}

	environments, err = r.LicensingClient.GetEnvironmentsForBillingPolicy(ctx, plan.BillingPolicyId)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.ProviderTypeName), err.Error())
		return
	}

	plan.Environments = environments

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *BillingPolicyEnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *BillingPolicyEnvironmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.LicensingClient.RemoveEnvironmentsToBillingPolicy(ctx, state.BillingPolicyId, state.Environments)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *BillingPolicyEnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}