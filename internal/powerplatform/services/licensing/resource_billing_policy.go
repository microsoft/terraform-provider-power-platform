package powerplatform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/clients"
)

var _ resource.Resource = &BillingPolicyResource{}
var _ resource.ResourceWithImportState = &BillingPolicyResource{}

func NewBillingPolicyResource() resource.Resource {
	return &BillingPolicyResource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_billing_policy",
	}
}

type BillingPolicyResource struct {
	LicensingClient  LicensingClient
	ProviderTypeName string
	TypeName         string
}

type BillingPolicyResourceModel struct {
	Id                types.String                   `tfsdk:"id"`
	Name              types.String                   `tfsdk:"name"`
	Location          types.String                   `tfsdk:"location"`
	Status            types.String                   `tfsdk:"status"`
	BillingInstrument BillingInstrumentResourceModel `tfsdk:"billing_instrument"`
}

type BillingInstrumentResourceModel struct {
	Id             types.String `tfsdk:"id"`
	ResourceGroup  types.String `tfsdk:"resource_group"`
	SubscriptionId types.String `tfsdk:"subscription_id"`
}

// Metadata
func (r *BillingPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

// Schema
func (r *BillingPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Power Platform Billing Policy",
		MarkdownDescription: "[Power Platform Billing Policy](https://learn.microsoft.com/en-us/rest/api/power-platform/licensing/billing-policy/get-billing-policy#billingpolicyresponsemodel)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The id of the billing policy",
				MarkdownDescription: "The id of the billing policy",
			},
			"name": schema.StringAttribute{
				Description:         "The name of the billing policy",
				MarkdownDescription: "The name of the billing policy",
				Required:            true,
			},
			"location": schema.StringAttribute{
				Description:         "The location of the billing policy",
				MarkdownDescription: "The location of the billing policy",
				Required:            true,
			},
			"status": schema.StringAttribute{
				Description:         "The status of the billing policy",
				MarkdownDescription: "The status of the billing policy (Enabled, Disabled)",
				Computed:            true,
			},
			"billing_instrument": schema.SingleNestedAttribute{
				Description:         "The billing instrument of the billing policy",
				MarkdownDescription: "The billing instrument of the billing policy",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "The id of the billing instrument",
						MarkdownDescription: "The id of the billing instrument",
					},
					"resource_group": schema.StringAttribute{
						Description:         "The resource group of the billing instrument",
						MarkdownDescription: "The resource group of the billing instrument",
						Required:            true,
					},
					"subscription_id": schema.StringAttribute{
						Description:         "The subscription id of the billing instrument",
						MarkdownDescription: "The subscription id of the billing instrument",
						Required:            true,
					},
				},
			},
		},
	}
}

func (r *BillingPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientBapi := req.ProviderData.(*clients.ProviderClient).PowerPlatformApi.Client

	if clientBapi == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.LicensingClient = NewLicensingClient(clientBapi)
}

func (r *BillingPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *BillingPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	billingPolicyToCreate := BillingPolicyCreateDto{
		BillingInstrument: BillingInstrumentDto{
			ResourceGroup:  plan.BillingInstrument.ResourceGroup.ValueString(),
			SubscriptionId: plan.BillingInstrument.SubscriptionId.ValueString(),
		},
		Location: plan.Location.ValueString(),
		Name:     plan.Name.ValueString(),
		Status:   "Enabled",
	}

	policy, err := r.LicensingClient.CreateBillingPolicy(ctx, billingPolicyToCreate)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
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

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *BillingPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *BillingPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	billing, err := r.LicensingClient.GetBillingPolicy(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	state.Id = types.StringValue(billing.Id)
	state.Name = types.StringValue(billing.Name)
	state.Location = types.StringValue(billing.Location)
	state.Status = types.StringValue(billing.Status)
	state.BillingInstrument.Id = types.StringValue(billing.BillingInstrument.Id)
	state.BillingInstrument.ResourceGroup = types.StringValue(billing.BillingInstrument.ResourceGroup)
	state.BillingInstrument.SubscriptionId = types.StringValue(billing.BillingInstrument.SubscriptionId)

	//TODO move to separate function
	ctx = tflog.SetField(ctx, "id", state.Id.ValueString())
	ctx = tflog.SetField(ctx, "name", state.Name.ValueString())
	ctx = tflog.SetField(ctx, "location", state.Location.ValueString())
	ctx = tflog.SetField(ctx, "status", state.Status.ValueString())
	ctx = tflog.SetField(ctx, "billing_instrument_id", state.BillingInstrument.Id.ValueString())
	ctx = tflog.SetField(ctx, "billing_instrument_resource_group", state.BillingInstrument.ResourceGroup.ValueString())
	ctx = tflog.SetField(ctx, "billing_instrument_subscription_id", state.BillingInstrument.SubscriptionId.ValueString())

	resp.Diagnostics.AddError(fmt.Sprintf("READ %s_%s with Id: %s", r.ProviderTypeName, r.TypeName, state.Id.ValueString()), err.Error())

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE END: %s", r.ProviderTypeName))
}

// Update
func (r *BillingPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

// Delete
func (r *BillingPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (r *BillingPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
