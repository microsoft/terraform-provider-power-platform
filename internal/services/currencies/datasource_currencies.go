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
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

var (
	_ datasource.DataSource              = &DataSource{}
	_ datasource.DataSourceWithConfigure = &DataSource{}
)

type DataSourceModel struct {
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	Id       types.Int64    `tfsdk:"id"`
	Location types.String   `tfsdk:"location"`
	Value    []DataModel    `tfsdk:"currencies"`
}

type DataModel struct {
	ID              string `tfsdk:"id"`
	Name            string `tfsdk:"name"`
	Type            string `tfsdk:"type"`
	Code            string `tfsdk:"code"`
	Symbol          string `tfsdk:"symbol"`
	IsTenantDefault bool   `tfsdk:"is_tenant_default"`
}

func NewCurrenciesDataSource() datasource.DataSource {
	return &DataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_currencies",
	}
}

type DataSource struct {
	CurrenciesClient Client
	ProviderTypeName string
	TypeName         string
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *DataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of available Dynamics 365 currencies",
		MarkdownDescription: "Fetches the list of available Dynamics 365 currencies. For more information see [Power Platform Currencies](https://learn.microsoft.com/power-platform/admin/manage-transactions-with-multiple-currencies)",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"id": schema.Int64Attribute{
				Description: "Id of the read operation",
				Optional:    true,
			},
			"location": schema.StringAttribute{
				Description: "Location of the currencies",
				Required:    true,
			},
			"currencies": schema.ListNestedAttribute{
				Description:         "List of available currencies",
				MarkdownDescription: "List of available currencies",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier of the currency",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the currency",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Type of the currency",
							Computed:    true,
						},
						"code": schema.StringAttribute{
							Description: "Code of the location",
							Computed:    true,
						},
						"symbol": schema.StringAttribute{
							Description: "Symbol of the currency",
							Computed:    true,
						},
						"is_tenant_default": schema.BoolAttribute{
							Description: "Is the currency the default for the tenant",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.CurrenciesClient = NewCurrenciesClient(clientApi)
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state DataSourceModel
	resp.State.Get(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE CURRENCIES START: %s", d.ProviderTypeName))

	timeout, diags := state.Timeouts.Read(ctx, constants.DEFAULT_RESOURCE_OPERATION_TIMEOUT_IN_MINUTES)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	currencies, err := d.CurrenciesClient.GetCurrenciesByLocation(ctx, state.Location.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}
	state.Id = types.Int64Value(int64(len(currencies.Value)))
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

	diags = resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE CURRENCIES END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
