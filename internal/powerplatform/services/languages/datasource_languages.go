// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package languages

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

var (
	_ datasource.DataSource              = &LanguagesDataSource{}
	_ datasource.DataSourceWithConfigure = &LanguagesDataSource{}
)

func NewLanguagesDataSource() datasource.DataSource {
	return &LanguagesDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_languages",
	}
}

type LanguagesDataSource struct {
	LanguagesClient  LanguagesClient
	ProviderTypeName string
	TypeName         string
}

type LanguagesDataSourceModel struct {
	Id       types.Int64         `tfsdk:"id"`
	Location types.String        `tfsdk:"location"`
	Value    []LanguageDataModel `tfsdk:"languages"`
}

type LanguageDataModel struct {
	Name            string `tfsdk:"name"`
	ID              string `tfsdk:"id"`
	DisplayName     string `tfsdk:"display_name"`
	LocalizedName   string `tfsdk:"localized_name"`
	LocaleID        int64  `tfsdk:"locale_id"`
	IsTenantDefault bool   `tfsdk:"is_tenant_default"`
}

func (d *LanguagesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *LanguagesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of Dynamics 365 languages",
		MarkdownDescription: "Fetches the list of Dynamics 365 languages. For more information see [Power Platform Enable Languages](https://learn.microsoft.com/power-platform/admin/enable-languages)",
		Attributes: map[string]schema.Attribute{
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

func (d *LanguagesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.LanguagesClient = NewLanguagesClient(clientApi)
}

func (d *LanguagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan LanguagesDataSourceModel
	resp.State.Get(ctx, &plan)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE LANGUAGES START: %s", d.ProviderTypeName))

	languages, err := d.LanguagesClient.GetLanguagesByLocation(ctx, plan.Location.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	plan.Id = types.Int64Value(int64(len(languages.Value)))
	plan.Location = types.StringValue(plan.Location.ValueString())

	for _, language := range languages.Value {
		plan.Value = append(plan.Value, LanguageDataModel{
			ID:              language.ID,
			Name:            language.Name,
			DisplayName:     language.Properties.DisplayName,
			LocalizedName:   language.Properties.LocalizedName,
			LocaleID:        language.Properties.LocaleID,
			IsTenantDefault: language.Properties.IsTenantDefault,
		})
	}

	diags := resp.State.Set(ctx, &plan)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE LANGUAGES END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
