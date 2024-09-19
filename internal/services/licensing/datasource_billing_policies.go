// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package licensing

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var (
	_ datasource.DataSource              = &BillingPoliciesDataSource{}
	_ datasource.DataSourceWithConfigure = &BillingPoliciesDataSource{}
)

func NewBillingPoliciesDataSource() datasource.DataSource {
	return &BillingPoliciesDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "billing_policies",
		},
	}
}

type BillingPoliciesDataSource struct {
	helpers.TypeInfo
	LicensingClient Client
}

type BillingPoliciesListDataSourceModel struct {
	Timeouts        timeouts.Value                 `tfsdk:"timeouts"`
	BillingPolicies []BillingPolicyDataSourceModel `tfsdk:"billing_policies"`
	Id              types.Int64                    `tfsdk:"id"`
}

type BillingPolicyDataSourceModel struct {
	Id                types.String                     `tfsdk:"id"`
	Name              types.String                     `tfsdk:"name"`
	Location          types.String                     `tfsdk:"location"`
	Status            types.String                     `tfsdk:"status"`
	BillingInstrument BillingInstrumentDataSourceModel `tfsdk:"billing_instrument"`
}

type BillingInstrumentDataSourceModel struct {
	Id             types.String `tfsdk:"id"`
	ResourceGroup  types.String `tfsdk:"resource_group"`
	SubscriptionId types.String `tfsdk:"subscription_id"`
}

func (d *BillingPoliciesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *BillingPoliciesDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of billing policies in a tenant",
		MarkdownDescription: "Fetches the list of [billing policies](https://learn.microsoft.com/power-platform/admin/pay-as-you-go-overview#what-is-a-billing-policy) in a tenant. A billing policy is a set of rules that define how a tenant is billed for usage of Power Platform services. A billing policy is associated with a billing instrument, which is a subscription and resource group that is used to pay for usage of Power Platform services.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"id": schema.Int64Attribute{
				Description:         "Id of the read operation",
				MarkdownDescription: "Id of the read operation",
				Computed:            true,
			},
			"billing_policies": schema.ListNestedAttribute{
				Description:         "Power Platform Billing Policy",
				MarkdownDescription: "[Power Platform Billing Policy](https://learn.microsoft.com/rest/api/power-platform/licensing/billing-policy/get-billing-policy#billingpolicyresponsemodel)",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
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
							Optional:            true,
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
				},
			},
		},
	}
}

func (d *BillingPoliciesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.LicensingClient = NewLicensingClient(clientApi)
}

func (d *BillingPoliciesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state BillingPoliciesListDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE BILLING POLICIES START: %s", d.ProviderTypeName))

	policies, err := d.LicensingClient.GetBillingPolicies(ctx)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	for _, policy := range policies {
		state.BillingPolicies = append(state.BillingPolicies, BillingPolicyDataSourceModel{
			Id:       types.StringValue(policy.Id),
			Name:     types.StringValue(policy.Name),
			Location: types.StringValue(policy.Location),
			Status:   types.StringValue(policy.Status),
			BillingInstrument: BillingInstrumentDataSourceModel{
				Id:             types.StringValue(policy.BillingInstrument.Id),
				ResourceGroup:  types.StringValue(policy.BillingInstrument.ResourceGroup),
				SubscriptionId: types.StringValue(policy.BillingInstrument.SubscriptionId),
			},
		})
	}

	state.Id = types.Int64Value(int64(len(policies)))
	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE BILLING POLICIES END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
