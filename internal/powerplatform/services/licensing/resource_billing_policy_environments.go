package powerplatform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/clients"
)

var _ resource.Resource = &BillingPolicyEnvironmentsResource{}
var _ resource.ResourceWithImportState = &BillingPolicyEnvironmentsResource{}

func NewBillingPolicyEnvironmentsResource() resource.Resource {
	return &BillingPolicyEnvironmentsResource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_billing_policy_environments",
	}
}

type BillingPolicyEnvironmentsResource struct {
	LicensingClient  LicensingClient
	ProviderTypeName string
	TypeName         string
}

type BillingPolicyEnvironmentsResourceModel struct {
	BillingPolicyId string   `tfsdk:"billing_policy_id"`
	Environments    []string `tfsdk:"environments"`
}

func (r *BillingPolicyEnvironmentsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

func (r *BillingPolicyEnvironmentsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Power Platform Billing Policy Environments",
		MarkdownDescription: "Power Platform Billing Policy Environments",
		Attributes: map[string]schema.Attribute{
			"billing_policy_id": schema.StringAttribute{
				Computed:            true,
				Description:         "The id of the billing policy",
				MarkdownDescription: "The id of the billing policy",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// "environments": schema.ListAttribute{
			// 	Description:         "The environments associated with the billing policy",
			// 	MarkdownDescription: "The environments associated with the billing policy",
			// },
		},
	}
}

func (r *BillingPolicyEnvironmentsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BillingPolicyEnvironmentsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

}

func (r *BillingPolicyEnvironmentsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

func (r *BillingPolicyEnvironmentsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (r *BillingPolicyEnvironmentsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (r *BillingPolicyEnvironmentsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
