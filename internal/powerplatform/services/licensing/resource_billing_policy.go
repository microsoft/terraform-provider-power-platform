package powerplatform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
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
		Description:         "Manages a Power Platform Billing Policy. \n\nA Power Platform billing policy is a mechanism that allows you to manage the costs associated with your Power Platform usage. It's linked to an Azure subscription and is used to set up pay-as-you-go billing for an environment.",
		MarkdownDescription: "Manages a Power Platform Billing Policy. \n\nA Power Platform billing policy is a mechanism that allows you to manage the costs associated with your Power Platform usage. It's linked to an Azure subscription and is used to set up pay-as-you-go billing for an environment.\n\nAdditional Resources:\n\n* [What is a billing policy](https://learn.microsoft.com/en-us/power-platform/admin/pay-as-you-go-overview#what-is-a-billing-policy)\n* [Power Platform Billing Policy API](https://learn.microsoft.com/en-us/rest/api/power-platform/licensing/billing-policy/get-billing-policy)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The id of the billing policy",
				MarkdownDescription: "The id of the billing policy",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Description:         "The status of the billing policy",
				MarkdownDescription: "The status of the billing policy (Enabled, Disabled)",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Enabled", "Disabled"),
				},
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
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"resource_group": schema.StringAttribute{
						Description:         "The resource group of the billing instrument",
						MarkdownDescription: "The resource group of the billing instrument",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"subscription_id": schema.StringAttribute{
						Description:         "The subscription id of the billing instrument",
						MarkdownDescription: "The subscription id of the billing instrument",
						Required:            true,
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
	}

	if plan.Status.IsUnknown() {
		billingPolicyToCreate.Status = "Enabled"
	} else {
		billingPolicyToCreate.Status = plan.Status.ValueString()
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

	tflog.Debug(ctx, fmt.Sprintf("READ %s_%s with Id: %s", r.ProviderTypeName, r.TypeName, billing.Id))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE END: %s", r.ProviderTypeName))
}

func (r *BillingPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *BillingPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

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
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.ProviderTypeName), err.Error())
			return
		}

		plan.Id = types.StringValue(policy.Id)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *BillingPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *BillingPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.LicensingClient.DeleteBillingPolicy(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *BillingPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
