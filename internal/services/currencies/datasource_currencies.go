// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package currencies

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
	_ datasource.DataSource              = &DataSource{}
	_ datasource.DataSourceWithConfigure = &DataSource{}
)

func NewCurrenciesDataSource() datasource.DataSource {
	return &DataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "currencies",
		},
	}
}

type DataSource struct {
	helpers.TypeInfo
	CurrenciesClient client
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches the list of available Dynamics 365 currencies. For more information see [Power Platform Currencies](https://learn.microsoft.com/power-platform/admin/manage-transactions-with-multiple-currencies)",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"location": schema.StringAttribute{
				MarkdownDescription: "Location of the currencies",
				Required:            true,
			},
			"currencies": schema.ListNestedAttribute{
				MarkdownDescription: "List of available currencies",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Unique identifier of the currency",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the currency",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the currency",
							Computed:            true,
						},
						"code": schema.StringAttribute{
							MarkdownDescription: "Code of the location",
							Computed:            true,
						},
						"symbol": schema.StringAttribute{
							MarkdownDescription: "Symbol of the currency",
							Computed:            true,
						},
						"is_tenant_default": schema.BoolAttribute{
							MarkdownDescription: "Is the currency the default for the tenant",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
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
	d.CurrenciesClient = newCurrenciesClient(client.Api)
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var state DataSourceModel
	resp.State.Get(ctx, &state)

	currencies, err := d.CurrenciesClient.GetCurrenciesByLocation(ctx, state.Location.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
		return
	}
	state.Location = types.StringValue(state.Location.ValueString())

	for _, location := range currencies.Value {
		state.Value = append(state.Value, DataModel{
			ID:              location.ID,
			Name:            location.Name,
			Type:            location.Type,
			Code:            location.Properties.Code,
			Symbol:          location.Properties.Symbol,
			IsTenantDefault: location.Properties.IsTenantDefault,
		})
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
