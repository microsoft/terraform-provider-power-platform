package billing_policy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	powerplatform_bapi "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi"
)

var _ resource.Resource = &BillingPolicyResource{}
var _ resource.ResourceWithImportState = &BillingPolicyResource{}

type BillingPolicyResource struct {
	ApiClient        powerplatform_bapi.ApiClientInterface
	ProviderTypeName string
	TypeName         string
}

func NewBillingPolicyResource() resource.Resource {
	return &BillingPolicyResource{
		ProviderTypeName: "powerplatform",
		TypeName:         "billing_policy",
	}
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
			},
			"location": schema.StringAttribute{
				Description:         "The location of the billing policy",
				MarkdownDescription: "The location of the billing policy",
			},
			"status": schema.StringAttribute{
				Description:         "The status of the billing policy",
				MarkdownDescription: "The status of the billing policy",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"Enabled", "Disabled"}...),
				},
			},
			"billing_instrument": schema.SingleNestedAttribute{
				Description:         "The billing instrument of the billing policy",
				MarkdownDescription: "The billing instrument of the billing policy",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "The id of the billing instrument",
						MarkdownDescription: "The id of the billing instrument",
					},
					"resource_group": schema.StringAttribute{
						Description:         "The resource group of the billing instrument",
						MarkdownDescription: "The resource group of the billing instrument",
					},
					"subscription_id": schema.StringAttribute{
						Description:         "The subscription id of the billing instrument",
						MarkdownDescription: "The subscription id of the billing instrument",
					},
				},
			},
		},
	}
}

// Configure
func (r *BillingPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	provider, ok := req.ProviderData.(*PowerPlatformProvider)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider type", fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}

	r.ApiClient = provider.bapiClient
}

// ImportState
func (r *BillingPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create
func (r *BillingPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *BillingPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dto := powerplatform_bapi.BillingPolicyDto{
		Id:       plan.Id.ValueString(),
		Name:     plan.Name.ValueString(),
		Location: plan.Location.ValueString(),
		Status:   plan.Status.ValueString(),
		BillingInstrument: powerplatform_bapi.BillingInstrumentDto{
			Id:             plan.BillingInstrument.Id.ValueString(),
			ResourceGroup:  plan.BillingInstrument.ResourceGroup.ValueString(),
			SubscriptionId: plan.BillingInstrument.SubscriptionId.ValueString(),
		},
	}

}

// Read
func (r *BillingPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}
