// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

var (
	_ datasource.DataSource              = &DataSourceTenant{}
	_ datasource.DataSourceWithConfigure = &DataSourceTenant{}
)

type DataSourceTenant struct {
	TenantClient     ClientTenant
	ProviderTypeName string
	TypeName         string
}

type ModelTenantDataSource struct {
	TenantId                         types.String `tfsdk:"tenant_id"`
	State                            types.String `tfsdk:"state"`
	Location                         types.String `tfsdk:"location"`
	AadCountryGeo                    types.String `tfsdk:"aad_country_geo"`
	DataStorageGeo                   types.String `tfsdk:"data_storage_geo"`
	DefaultEnvironmentGeo            types.String `tfsdk:"default_environment_geo"`
	AadDataBoundary                  types.String `tfsdk:"aad_data_boundary"`
	FedRAMPHighCertificationRequired types.Bool   `tfsdk:"fed_ramp_high_certification_required"`
}

func NewTenantDataSource() datasource.DataSource {
	return &DataSourceTenant{
		ProviderTypeName: "powerplatform",
		TypeName:         "_tenant",
	}
}

func (d *DataSourceTenant) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (_ *DataSourceTenant) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the client configuration for the given tenant.",
		MarkdownDescription: "Fetches the client configuration for the given tenant.",
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Description:         "Tenant ID of the application.",
				MarkdownDescription: "Tenant ID of the application.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				Description:         "State of the tenant.",
				MarkdownDescription: "State of the tenant.",
				Computed:            true,
			},
			"location": schema.StringAttribute{
				Description:         "Location of the tenant.",
				MarkdownDescription: "Location of the tenant.",
				Computed:            true,
			},
			"aad_country_geo": schema.StringAttribute{
				Description:         "AAD country geo.",
				MarkdownDescription: "AAD country geo.",
				Computed:            true,
			},
			"data_storage_geo": schema.StringAttribute{
				Description:         "Data storage geo.",
				MarkdownDescription: "Data storage geo.",
				Computed:            true,
			},
			"default_environment_geo": schema.StringAttribute{
				Description:         "Default environment geo.",
				MarkdownDescription: "Default environment geo.",
				Computed:            true,
			},
			"aad_data_boundary": schema.StringAttribute{
				Description:         "AAD data boundary.",
				MarkdownDescription: "AAD data boundary.",
				Computed:            true,
			},
			"fed_ramp_high_certification_required": schema.BoolAttribute{
				Description:         "FedRAMP high certification required.",
				MarkdownDescription: "FedRAMP high certification required.",
				Computed:            true,
			},
		},
	}
}

func (d *DataSourceTenant) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	dto, err := d.TenantClient.GetTenant(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to fetch tenant", fmt.Sprintf("Failed to fetch tenant: %v", err))
		return
	}

	state := ModelTenantDataSource{
		TenantId:                         types.StringValue(dto.TenantId),
		State:                            types.StringValue(dto.State),
		Location:                         types.StringValue(dto.Location),
		AadCountryGeo:                    types.StringValue(dto.AadCountryGeo),
		DataStorageGeo:                   types.StringValue(dto.DataStorageGeo),
		DefaultEnvironmentGeo:            types.StringValue(dto.DefaultEnvironmentGeo),
		AadDataBoundary:                  types.StringValue(dto.AadDataBoundary),
		FedRAMPHighCertificationRequired: types.BoolValue(dto.FedRAMPHighCertificationRequired),
	}

	resp.State.Set(ctx, &state)
}

func (d *DataSourceTenant) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := req.ProviderData.(*api.ProviderClient).Api
	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.TenantClient = NewTenantClient(client)
}
