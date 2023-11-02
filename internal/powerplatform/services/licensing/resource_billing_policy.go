package powerplatform

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
				MarkdownDescription: "The status of the billing policy",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"Enabled", "Disabled"}...),
				},
				Required: true,
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

	clientBapi := req.ProviderData.(*clients.ProviderClient).LicensingApi.Client

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

	bill := BillingPolicyCreateDto{
		BillingInstrument: BillingInstrumentDto{
			Id:               "",
			Location:         "europe",
			ResourceGroup:    "tmp",
			SubscriptionId:   "2bc1f261-7e26-490c-9fd5-b7ca72032ad3",
			SubscriptionName: "Visual Studio Enterprise",
			//Tags:             []string{},
		},
		Location: "europe",
		Name:     "afssadadsadad",
		PowerAppsPolicy: PolicyDto{
			PayAsYouGoState: "Enabled",
		},
		PowerAutomatePolicy: PowerAutomatePolicyCreateDto{
			PayAsYouGoState: "Enabled",
		},
		StoragePolicy: PolicyDto{
			PayAsYouGoState: "Enabled",
		},
		TenantType: "TenantOwned",
	}

	_, err := r.LicensingClient.CreateBillingPolicy(ctx, bill)
	if err != nil {
		resp.Diagnostics.AddError("Error creating billing policy", err.Error())
		return
	}

	//tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE END: %s", return1.createdBy.id))

	// dto := powerplatform_bapi.BillingPolicyDto{
	// 	Id:       plan.Id.ValueString(),
	// 	Name:     plan.Name.ValueString(),
	// 	Location: plan.Location.ValueString(),
	// 	Status:   plan.Status.ValueString(),
	// 	BillingInstrument: powerplatform_bapi.BillingInstrumentDto{
	// 		Id:             plan.BillingInstrument.Id.ValueString(),
	// 		ResourceGroup:  plan.BillingInstrument.ResourceGroup.ValueString(),
	// 		SubscriptionId: plan.BillingInstrument.SubscriptionId.ValueString(),
	// 	},
	// }

	//r.ApiClient.Execute().(ctx, dto)

}

// Read
func (r *BillingPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

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
