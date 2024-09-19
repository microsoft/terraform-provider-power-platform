// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package languages

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

func NewLanguagesDataSource() datasource.DataSource {
	return &DataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "languages",
		},
	}
}

type DataSource struct {
	helpers.TypeInfo
	LanguagesClient Client
}

type DataSourceModel struct {
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	Id       types.Int64    `tfsdk:"id"`
	Location types.String   `tfsdk:"location"`
	Value    []DataModel    `tfsdk:"languages"`
}

type DataModel struct {
	Name            string `tfsdk:"name"`
	ID              string `tfsdk:"id"`
	DisplayName     string `tfsdk:"display_name"`
	LocalizedName   string `tfsdk:"localized_name"`
	LocaleID        int64  `tfsdk:"locale_id"`
	IsTenantDefault bool   `tfsdk:"is_tenant_default"`
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

func (d *DataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of Dynamics 365 languages",
		MarkdownDescription: "Fetches the list of Dynamics 365 languages. For more information see [Power Platform Enable Languages](https://learn.microsoft.com/power-platform/admin/enable-languages)",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"id": schema.Int64Attribute{
				Description:         "Id of the read operation",
				MarkdownDescription: "Id of the read operation",
				Optional:            true,
			},
			"location": schema.StringAttribute{
				Description: "Location of the languages",
				Required:    true,
			},
			"languages": schema.ListNestedAttribute{
				Description:         "List of available languages",
				MarkdownDescription: "List of available languages",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the location",
							Computed:    true,
						},
						"id": schema.StringAttribute{
							Description: "Unique identifier of the location",
							Computed:    true,
						},
						"display_name": schema.StringAttribute{
							Description: "Display name of the location",
							Computed:    true,
						},
						"localized_name": schema.StringAttribute{
							Description: "Localized name of the location",
							Computed:    true,
						},
						"locale_id": schema.Int64Attribute{
							Description: "Locale ID of the location",
							Computed:    true,
						},
						"is_tenant_default": schema.BoolAttribute{
							Description: "Is the location the default for the tenant",
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
		resp.Diagnostics.AddError("Failed to configure %s because provider data is nil", d.TypeName)
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
	d.LanguagesClient = NewLanguagesClient(clientApi)
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	var state DataSourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE LANGUAGES START: %s", d.ProviderTypeName))

	languages, err := d.LanguagesClient.GetLanguagesByLocation(ctx, state.Location.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	state.Id = types.Int64Value(int64(len(languages.Value)))
	state.Location = types.StringValue(state.Location.ValueString())

	for _, language := range languages.Value {
		state.Value = append(state.Value, DataModel{
			ID:              language.ID,
			Name:            language.Name,
			DisplayName:     language.Properties.DisplayName,
			LocalizedName:   language.Properties.LocalizedName,
			LocaleID:        language.Properties.LocaleID,
			IsTenantDefault: language.Properties.IsTenantDefault,
		})
	}

	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE LANGUAGES END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
